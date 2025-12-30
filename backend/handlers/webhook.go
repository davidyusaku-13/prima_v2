package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/utils"
	"github.com/gin-gonic/gin"
)

// WebhookHandler handles GOWA webhook callbacks for delivery status updates
type WebhookHandler struct {
	patientStore *models.PatientStore
	config       *config.Config
	logger       *slog.Logger
	sseHandler   *SSEHandler // SSE handler for broadcasting updates
}

// processedWebhook tracks webhooks that have been processed for idempotency
type processedWebhook struct {
	messageID   string
	status      string
	processedAt time.Time
}

var (
	// webhookProcessor tracks processed webhooks for idempotency
	webhookProcessor = struct {
		mu       sync.RWMutex
		webhooks map[string]processedWebhook
	}{
		webhooks: make(map[string]processedWebhook),
	}
	webhookTTL = 24 * time.Hour
)

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(patientStore *models.PatientStore, cfg *config.Config, logger *slog.Logger) *WebhookHandler {
	return &WebhookHandler{
		patientStore: patientStore,
		config:       cfg,
		logger:       logger,
		sseHandler:   nil, // Will be set via SetSSEHandler
	}
}

// SetSSEHandler sets the SSE handler for broadcasting delivery status updates
func (h *WebhookHandler) SetSSEHandler(sseHandler *SSEHandler) {
	h.sseHandler = sseHandler
}

// GOWAPayload represents the webhook payload from GOWA
type GOWAPayload struct {
	Event   string      `json:"event"`
	Message MessageAck  `json:"message"`
}

// MessageAck represents the message acknowledgment data
type MessageAck struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// WebhookResponse represents the response from webhook processing
type WebhookResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
}

// HandleGOWAWebhook processes incoming GOWA webhook callbacks
// POST /api/webhook/gowa
func (h *WebhookHandler) HandleGOWAWebhook(c *gin.Context) {
	// Read request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("Failed to read webhook body",
				"error", err.Error(),
			)
		}
		c.JSON(http.StatusBadRequest, WebhookResponse{
			Error: "Failed to read request body",
			Code:  "INVALID_REQUEST",
		})
		return
	}

	// Validate HMAC signature
	signature := c.GetHeader("X-Webhook-Signature")
	if signature == "" {
		if h.logger != nil {
			h.logger.Warn("Webhook received without signature header")
		}
		c.JSON(http.StatusUnauthorized, WebhookResponse{
			Error: "Missing webhook signature",
			Code:  "MISSING_SIGNATURE",
		})
		return
	}

	// Validate signature using HMAC-SHA256 (uses utils for secure validation)
	if !utils.ValidateWebhookSignature(body, signature, h.config.GOWA.WebhookSecret) {
		if h.logger != nil {
			h.logger.Warn("Webhook received with invalid signature",
				"signature_length", len(signature),
			)
		}
		c.JSON(http.StatusUnauthorized, WebhookResponse{
			Error: "Invalid webhook signature",
			Code:  "INVALID_SIGNATURE",
		})
		return
	}

	// Parse webhook payload
	var payload GOWAPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		if h.logger != nil {
			h.logger.Error("Failed to parse webhook payload",
				"error", err.Error(),
			)
		}
		c.JSON(http.StatusBadRequest, WebhookResponse{
			Error: "Invalid webhook payload",
			Code:  "INVALID_PAYLOAD",
		})
		return
	}

	// Check for idempotency - skip if already processed
	idempotencyKey := payload.Message.ID + ":" + payload.Message.Status
	if isWebhookProcessed(idempotencyKey) {
		if h.logger != nil {
			h.logger.Debug("Webhook already processed (idempotent skip)",
				"message_id", payload.Message.ID,
				"status", payload.Message.Status,
			)
		}
		c.JSON(http.StatusOK, WebhookResponse{
			Data:    map[string]string{"message_id": payload.Message.ID},
			Message: "Webhook already processed",
		})
		return
	}

	// Process based on event type
	switch payload.Event {
	case "message.ack":
		h.processMessageAck(c, &payload)
	default:
		if h.logger != nil {
			h.logger.Warn("Unknown webhook event type",
				"event", payload.Event,
				"message_id", payload.Message.ID,
			)
		}
		c.JSON(http.StatusOK, WebhookResponse{
			Data:    map[string]string{"message_id": payload.Message.ID},
			Message: fmt.Sprintf("Event type '%s' acknowledged but not processed", payload.Event),
		})
	}
}

// processMessageAck processes message acknowledgment events
func (h *WebhookHandler) processMessageAck(c *gin.Context, payload *GOWAPayload) {
	messageID := payload.Message.ID
	newStatus := payload.Message.Status

	if h.logger != nil {
		h.logger.Info("Processing message acknowledgment",
			"message_id", messageID,
			"status", newStatus,
		)
	}

	// Find and update the reminder
	h.patientStore.Lock()
	defer h.patientStore.Unlock()

	var updatedReminder *models.Reminder
	var patientID string
	var patientName string
	var patientPhone string
	var previousStatus string

	// Search through all patients for the reminder with matching GOWA message ID
	for _, patient := range h.patientStore.Patients {
		for _, reminder := range patient.Reminders {
			if reminder.GOWAMessageID == messageID {
				// Capture previous status BEFORE updating
				previousStatus = string(reminder.DeliveryStatus)

				// Update delivery status based on acknowledgment status
				switch newStatus {
				case "delivered":
					reminder.DeliveryStatus = models.DeliveryStatusDelivered
					reminder.DeliveredAt = time.Now().UTC().Format(time.RFC3339)
				case "read":
					reminder.DeliveryStatus = models.DeliveryStatusRead
					reminder.ReadAt = time.Now().UTC().Format(time.RFC3339)
				case "failed":
					reminder.DeliveryStatus = models.DeliveryStatusFailed
					reminder.DeliveryErrorMessage = "Delivery failed according to GOWA webhook"
				default:
					// Log unknown status but don't update
					if h.logger != nil {
						h.logger.Warn("Unknown message status in webhook",
							"message_id", messageID,
							"status", newStatus,
						)
					}
					c.JSON(http.StatusOK, WebhookResponse{
						Data:    map[string]string{"message_id": messageID},
						Message: fmt.Sprintf("Status '%s' acknowledged but not processed", newStatus),
					})
					return
				}

				updatedReminder = reminder
				patientID = patient.ID
				patientName = patient.Name
				patientPhone = utils.MaskPhone(patient.Phone)
				break
			}
		}
		if updatedReminder != nil {
			break
		}
	}

	if updatedReminder == nil {
		if h.logger != nil {
			h.logger.Warn("Reminder not found for GOWA message ID",
				"message_id", messageID,
			)
		}
		c.JSON(http.StatusOK, WebhookResponse{
			Data:    map[string]string{"message_id": messageID},
			Message: "Message ID not found, may have been deleted",
		})
		return
	}

	// Mark webhook as processed (idempotency)
	markWebhookProcessed(messageID, newStatus)

	// Log the status update for audit purposes (FR36)
	if h.logger != nil {
		h.logger.Info("Reminder delivery status updated",
			"reminder_id", updatedReminder.ID,
			"patient_name", patientName,
			"patient_phone", patientPhone,
			"previous_status", previousStatus,
			"new_status", newStatus,
			"message_id", messageID,
		)
	}

	// Save data
	h.patientStore.SaveData()

	// Broadcast SSE event for real-time updates (if SSE handler is configured)
	if h.sseHandler != nil {
		h.sseHandler.BroadcastDeliveryStatusUpdate(
			updatedReminder.ID,
			string(updatedReminder.DeliveryStatus),
			time.Now().UTC().Format(time.RFC3339),
		)

		// Broadcast delivery.failed event if status is failed
		if newStatus == "failed" {
			h.sseHandler.BroadcastDeliveryFailed(
				updatedReminder.ID,
				patientID,
				patientName,
				updatedReminder.DeliveryErrorMessage,
			)
		}
	} else {
		// Log warning if SSE handler not configured (should be set via SetSSEHandler)
		if h.logger != nil {
			h.logger.Warn("SSE handler not configured, real-time updates unavailable",
				"reminder_id", updatedReminder.ID,
			)
		}
	}

	c.JSON(http.StatusOK, WebhookResponse{
		Data: map[string]interface{}{
			"message_id":      messageID,
			"reminder_id":     updatedReminder.ID,
			"delivery_status": newStatus,
		},
		Message: fmt.Sprintf("Reminder status updated to '%s'", newStatus),
	})
}

// isWebhookProcessed checks if a webhook has already been processed
func isWebhookProcessed(key string) bool {
	webhookProcessor.mu.RLock()
	defer webhookProcessor.mu.RUnlock()

	entry, exists := webhookProcessor.webhooks[key]
	if !exists {
		return false
	}

	// Check if entry has expired
	if time.Since(entry.processedAt) > webhookTTL {
		return false
	}

	return true
}

// markWebhookProcessed records a webhook as processed
func markWebhookProcessed(messageID, status string) {
	webhookProcessor.mu.Lock()
	defer webhookProcessor.mu.Unlock()

	key := messageID + ":" + status
	webhookProcessor.webhooks[key] = processedWebhook{
		messageID:   messageID,
		status:      status,
		processedAt: time.Now(),
	}

	// Clean up old entries periodically
	if len(webhookProcessor.webhooks) > 1000 {
		cleanupExpiredWebhooks()
	}
}

// cleanupExpiredWebhooks removes expired webhook entries
func cleanupExpiredWebhooks() {
	cutoff := time.Now().Add(-webhookTTL)
	for key, entry := range webhookProcessor.webhooks {
		if entry.processedAt.Before(cutoff) {
			delete(webhookProcessor.webhooks, key)
		}
	}
}
