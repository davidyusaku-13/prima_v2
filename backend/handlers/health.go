package handlers

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/services"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	patientStore   *models.PatientStore
	gowaClient     *services.GOWAClient
	lastGOWAPing   time.Time
	gowaConnected  bool
	mu             struct {
		sync.RWMutex
	}
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(patientStore *models.PatientStore, gowaClient *services.GOWAClient) *HealthHandler {
	h := &HealthHandler{
		patientStore:  patientStore,
		gowaClient:    gowaClient,
		lastGOWAPing:  time.Time{},
		gowaConnected: false,
	}
	return h
}

// HealthStatus represents the basic health status response
type HealthStatus struct {
	Status     string `json:"status"`
	Timestamp  string `json:"timestamp"`
}

// DetailedHealthStatus represents the detailed health status response for admins
type DetailedHealthStatus struct {
	Status     string                  `json:"status"`
	Timestamp  string                  `json:"timestamp"`
	GOWA       GOWAHealthStatus        `json:"gowa"`
	CircuitBreaker CircuitBreakerStatus `json:"circuit_breaker"`
	Queue      QueueStatus             `json:"queue"`
}

// GOWAHealthStatus represents GOWA connectivity status
type GOWAHealthStatus struct {
	Connected  bool   `json:"connected"`
	LastPing   string `json:"last_ping,omitempty"`
	Endpoint   string `json:"endpoint"`
}

// CircuitBreakerStatus represents circuit breaker state
type CircuitBreakerStatus struct {
	State               string `json:"state"`
	FailureCount        int    `json:"failure_count"`
	CooldownRemaining   int    `json:"cooldown_remaining_seconds"`
}

// QueueStatus represents reminder queue status
type QueueStatus struct {
	Total      int `json:"total"`
	Scheduled  int `json:"scheduled"`
	Retrying   int `json:"retrying"`
	QuietHours int `json:"quiet_hours"`
}

// GetHealth returns basic health status (public endpoint)
func (h *HealthHandler) GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": HealthStatus{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
		"message": "Health check successful",
	})
}

// GetHealthDetailed returns detailed health status (admin only)
func (h *HealthHandler) GetHealthDetailed(c *gin.Context) {
	// Check admin role
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required", "code": "ADMIN_REQUIRED"})
		return
	}

	now := time.Now().UTC()

	// Get GOWA status
	h.mu.RLock()
	lastPing := h.lastGOWAPing
	gowaConnected := h.gowaConnected
	h.mu.RUnlock()

	gowaEndpoint := ""
	if h.gowaClient != nil {
		gowaEndpoint = h.gowaClient.GetEndpoint()
	}

	// Get circuit breaker state
	circuitState := "closed"
	failureCount := 0
	cooldownRemaining := int(0)

	if h.gowaClient != nil {
		details := h.gowaClient.GetCircuitBreakerDetails()
		circuitState = details.State
		failureCount = details.FailureCount
		cooldownRemaining = int(details.CooldownRemaining.Seconds())
	}

	// Get queue counts
	queueCounts := h.getQueueCounts()

	// Format last ping time
	var lastPingStr string
	if !lastPing.IsZero() {
		lastPingStr = lastPing.Format(time.RFC3339)
	}

	response := DetailedHealthStatus{
		Status:    "ok",
		Timestamp: now.Format(time.RFC3339),
		GOWA: GOWAHealthStatus{
			Connected:  gowaConnected,
			LastPing:   lastPingStr,
			Endpoint:   gowaEndpoint,
		},
		CircuitBreaker: CircuitBreakerStatus{
			State:             circuitState,
			FailureCount:      failureCount,
			CooldownRemaining: cooldownRemaining,
		},
		Queue: QueueStatus{
			Total:      queueCounts.Total,
			Scheduled:  queueCounts.Scheduled,
			Retrying:   queueCounts.Retrying,
			QuietHours: queueCounts.QuietHours,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    response,
		"message": "Detailed health status retrieved",
	})
}

// QueueCounts represents queue count results
type QueueCounts struct {
	Total      int
	Scheduled  int
	Retrying   int
	QuietHours int
}

// getQueueCounts returns counts of reminders by queue category
func (h *HealthHandler) getQueueCounts() QueueCounts {
	counts := QueueCounts{}

	if h.patientStore == nil {
		return counts
	}

	h.patientStore.RLock()
	defer h.patientStore.RUnlock()

	now := time.Now().UTC()

	for _, patient := range h.patientStore.Patients {
		for _, reminder := range patient.Reminders {
			switch reminder.DeliveryStatus {
			case models.DeliveryStatusPending:
				// Pending reminders - count as scheduled
				counts.Scheduled++
				counts.Total++

			case models.DeliveryStatusScheduled:
				// Scheduled for quiet hours delivery
				if reminder.ScheduledDeliveryAt != "" {
					scheduledTime, err := time.Parse(time.RFC3339, reminder.ScheduledDeliveryAt)
					if err == nil && scheduledTime.After(now) {
						counts.QuietHours++
						counts.Total++
						continue
					}
				}
				// If no scheduled time or in the past, count as scheduled
				counts.Scheduled++
				counts.Total++

			case models.DeliveryStatusRetrying:
				// Retrying after transient failure
				if reminder.RetryCount < 3 {
					counts.Retrying++
					counts.Total++
				}

			case models.DeliveryStatusFailed:
				// Failed but still has retries left
				if reminder.RetryCount < 3 {
					counts.Retrying++
					counts.Total++
				}
			}
		}
	}

	return counts
}

// UpdateGOWAPing updates the last GOWA ping time and connection status
func (h *HealthHandler) UpdateGOWAPing(connected bool) {
	h.mu.Lock()
	h.lastGOWAPing = time.Now().UTC()
	h.gowaConnected = connected
	h.mu.Unlock()
}
