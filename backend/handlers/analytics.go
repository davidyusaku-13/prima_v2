package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/davidyusaku-13/prima_v2/models"
)

// DeliveryAnalyticsStats represents the delivery statistics response
type DeliveryAnalyticsStats struct {
	TotalSent        int               `json:"totalSent"`
	SuccessRate      float64           `json:"successRate"`
	FailedLast7Days  int               `json:"failedLast7Days"`
	AvgDeliveryTime  string            `json:"avgDeliveryTime"` // Human readable, e.g., "2m 30s"
	Breakdown        map[string]int    `json:"breakdown"`       // status -> count
	Period           string            `json:"period"`          // today, 7d, 30d, all
	PeriodStartDate  string            `json:"periodStartDate"` // ISO 8601 UTC
	PeriodEndDate    string            `json:"periodEndDate"`   // ISO 8601 UTC
}

// AnalyticsHandler handles delivery analytics endpoints
type AnalyticsHandler struct {
	patientStore *models.PatientStore
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(patientStore *models.PatientStore) *AnalyticsHandler {
	return &AnalyticsHandler{
		patientStore: patientStore,
	}
}

// parsePeriod parses the period query parameter and returns start/end dates
// period values: "today", "7d" (7 days), "30d" (30 days), "all"
func parsePeriod(period string) (startDate, endDate time.Time, err error) {
	now := time.Now().UTC()
	endDate = now

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	case "7d":
		startDate = now.AddDate(0, 0, -7)
	case "30d":
		startDate = now.AddDate(0, 0, -30)
	case "all", "":
		// Return zero time to indicate no filtering
		return time.Time{}, time.Time{}, nil
	default:
		return time.Time{}, time.Time{}, nil
	}

	return startDate, endDate, nil
}

// GetDeliveryAnalytics returns delivery statistics for admin dashboard
func (h *AnalyticsHandler) GetDeliveryAnalytics(c *gin.Context) {
	// Check admin role
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	period := c.DefaultQuery("period", "all")
	startDate, endDate, err := parsePeriod(period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period parameter"})
		return
	}

	// Track counts
	totalSent := 0
	failedLast7Days := 0
	breakdown := make(map[string]int)

	// For delivery time calculation
	var totalDeliveryDurations []time.Duration
	now := time.Now().UTC()

	// Collect all reminders
	h.patientStore.RLock()
	for _, patient := range h.patientStore.Patients {
		for _, reminder := range patient.Reminders {
			// Skip reminders that haven't been sent yet
			if reminder.DeliveryStatus == "" ||
				reminder.DeliveryStatus == models.DeliveryStatusPending ||
				reminder.DeliveryStatus == models.DeliveryStatusScheduled {
				continue
			}

			// Filter by date if period is specified
			if !startDate.IsZero() && reminder.MessageSentAt != "" {
				sentTime, parseErr := time.Parse(time.RFC3339, reminder.MessageSentAt)
				if parseErr != nil {
					continue
				}
				if sentTime.Before(startDate) || sentTime.After(endDate) {
					continue
				}

				// Count failed in last 7 days
				if reminder.DeliveryStatus == models.DeliveryStatusFailed {
					if sentTime.After(now.AddDate(0, 0, -7)) {
						failedLast7Days++
					}
				}
			}

			// Only count sent/delivered/read/failed as "sent"
			switch reminder.DeliveryStatus {
			case models.DeliveryStatusSent,
				models.DeliveryStatusDelivered,
				models.DeliveryStatusRead,
				models.DeliveryStatusFailed:
				totalSent++
			}

			// Track breakdown by status
			breakdown[reminder.DeliveryStatus]++

			// Calculate delivery time (sent -> delivered/read)
			if reminder.MessageSentAt != "" &&
				(reminder.DeliveredAt != "" || reminder.ReadAt != "") {
				sentTime, parseErr := time.Parse(time.RFC3339, reminder.MessageSentAt)
				if parseErr == nil {
					var deliveredTime time.Time
					if reminder.DeliveredAt != "" {
						deliveredTime, _ = time.Parse(time.RFC3339, reminder.DeliveredAt)
					} else if reminder.ReadAt != "" {
						deliveredTime, _ = time.Parse(time.RFC3339, reminder.ReadAt)
					}

					if !deliveredTime.IsZero() && deliveredTime.After(sentTime) {
						totalDeliveryDurations = append(totalDeliveryDurations, deliveredTime.Sub(sentTime))
					}
				}
			}
		}
	}
	h.patientStore.RUnlock()

	// Calculate success rate: (delivered + read) / (sent + delivered + read + failed) * 100
	deliveredCount := breakdown[models.DeliveryStatusDelivered]
	readCount := breakdown[models.DeliveryStatusRead]
	failedCount := breakdown[models.DeliveryStatusFailed]
	sentCount := breakdown[models.DeliveryStatusSent]

	deliverableCount := deliveredCount + readCount + sentCount + failedCount
	var successRate float64
	if deliverableCount > 0 {
		successRate = float64(deliveredCount+readCount) / float64(deliverableCount) * 100
	}

	// Calculate average delivery time
	var avgDeliveryTime string
	if len(totalDeliveryDurations) > 0 {
		var totalDuration time.Duration
		for _, d := range totalDeliveryDurations {
			totalDuration += d
		}
		avgDuration := totalDuration / time.Duration(len(totalDeliveryDurations))
		avgDeliveryTime = formatDuration(avgDuration)
	}

	// Ensure all expected statuses are in breakdown (even if 0)
	allStatuses := []string{
		models.DeliveryStatusPending,
		models.DeliveryStatusScheduled,
		models.DeliveryStatusQueued,
		models.DeliveryStatusSending,
		models.DeliveryStatusRetrying,
		models.DeliveryStatusSent,
		models.DeliveryStatusDelivered,
		models.DeliveryStatusRead,
		models.DeliveryStatusFailed,
		models.DeliveryStatusExpired,
	}
	for _, status := range allStatuses {
		if _, exists := breakdown[status]; !exists {
			breakdown[status] = 0
		}
	}

	// Sort breakdown keys for consistent output
	sortedStatuses := make([]string, 0, len(breakdown))
	for status := range breakdown {
		sortedStatuses = append(sortedStatuses, status)
	}
	slices.Sort(sortedStatuses)
	sortedBreakdown := make(map[string]int, len(breakdown))
	for _, status := range sortedStatuses {
		sortedBreakdown[status] = breakdown[status]
	}

	// Determine period start and end dates for response
	var periodStartDate, periodEndDate string
	if startDate.IsZero() {
		periodStartDate = ""
		periodEndDate = ""
	} else {
		periodStartDate = startDate.Format(time.RFC3339)
		periodEndDate = endDate.Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": DeliveryAnalyticsStats{
			TotalSent:       totalSent,
			SuccessRate:     successRate,
			FailedLast7Days: failedLast7Days,
			AvgDeliveryTime: avgDeliveryTime,
			Breakdown:       sortedBreakdown,
			Period:          period,
			PeriodStartDate: periodStartDate,
			PeriodEndDate:   periodEndDate,
		},
		"message": "success",
	})
}

// formatDuration formats a duration as a human readable string
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "< 1s"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
