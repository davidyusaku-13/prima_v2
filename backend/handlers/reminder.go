package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/services"
	"github.com/davidyusaku-13/prima_v2/utils"
)

// Role constants
const (
	RoleVolunteer = "volunteer"
)

// IDGenerator is a function type for generating unique IDs
type IDGenerator func() string

// ReminderHandler handles reminder-related HTTP requests
type ReminderHandler struct {
	store        *models.PatientStore
	config       *config.Config
	gowaClient   *services.GOWAClient
	logger       *slog.Logger
	generateID   IDGenerator
	contentStore *ContentStore // Added for attachment validation and content lookup
	sseHandler   *SSEHandler   // SSE handler for broadcasting delivery status updates
}

// NewReminderHandler creates a new reminder handler
func NewReminderHandler(store *models.PatientStore, cfg *config.Config, gowaClient *services.GOWAClient, logger *slog.Logger, idGen IDGenerator, contentStore *ContentStore) *ReminderHandler {
	return &ReminderHandler{
		store:        store,
		config:       cfg,
		gowaClient:   gowaClient,
		logger:       logger,
		generateID:   idGen,
		contentStore: contentStore,
		sseHandler:   nil, // Will be set via SetSSEHandler
	}
}

// SetSSEHandler sets the SSE handler for broadcasting delivery status updates
func (h *ReminderHandler) SetSSEHandler(sseHandler *SSEHandler) {
	h.sseHandler = sseHandler
}

// CreateReminderRequest represents the request body for creating a reminder
type CreateReminderRequest struct {
	Title       string              `json:"title" binding:"required"`
	Description string              `json:"description"`
	DueDate     string              `json:"dueDate"`
	Priority    string              `json:"priority"`
	Recurrence  models.Recurrence   `json:"recurrence"`
	Attachments []models.Attachment `json:"attachments"`
}

// UpdateReminderRequest represents the request body for updating a reminder
type UpdateReminderRequest struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	DueDate     string              `json:"dueDate"`
	Priority    string              `json:"priority"`
	Recurrence  models.Recurrence   `json:"recurrence"`
	Attachments []models.Attachment `json:"attachments"`
}

// MaxAttachments is the maximum number of content attachments per reminder
const MaxAttachments = 3

// validateAttachments validates the attachments array
func (h *ReminderHandler) validateAttachments(attachments []models.Attachment) error {
	for i, att := range attachments {
		if att.Type != "article" && att.Type != "video" {
			return fmt.Errorf("attachment[%d]: type must be 'article' or 'video'", i)
		}
		if att.ID == "" {
			return fmt.Errorf("attachment[%d]: id is required", i)
		}
		if att.Title == "" {
			return fmt.Errorf("attachment[%d]: title is required", i)
		}
		if len(att.Title) > 200 {
			return fmt.Errorf("attachment[%d]: title must not exceed 200 characters", i)
		}

		// Validate that attachment ID exists in content store
		if h.contentStore != nil {
			if att.Type == "article" {
				h.contentStore.Articles.Mu.RLock()
				if _, exists := h.contentStore.Articles.Articles[att.ID]; !exists {
					h.contentStore.Articles.Mu.RUnlock()
					return fmt.Errorf("attachment[%d]: article with id '%s' not found", i, att.ID)
				}
				h.contentStore.Articles.Mu.RUnlock()
			} else if att.Type == "video" {
				h.contentStore.Videos.Mu.RLock()
				if _, exists := h.contentStore.Videos.Videos[att.ID]; !exists {
					h.contentStore.Videos.Mu.RUnlock()
					return fmt.Errorf("attachment[%d]: video with id '%s' not found", i, att.ID)
				}
				h.contentStore.Videos.Mu.RUnlock()
			}
		}
	}
	return nil
}

// Create handles POST /patients/:id/reminders
func (h *ReminderHandler) Create(c *gin.Context) {
	patientID := c.Param("id")
	userID := c.GetString("userID")
	role := c.GetString("role")

	var req CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate attachments count
	if len(req.Attachments) > MaxAttachments {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      fmt.Sprintf("Maksimal %d konten yang dapat dilampirkan", MaxAttachments),
			"code":       "MAX_ATTACHMENTS_EXCEEDED",
			"max_count":  MaxAttachments,
			"actual":     len(req.Attachments),
		})
		return
	}

	// Validate attachments content
	if err := h.validateAttachments(req.Attachments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_ATTACHMENT",
		})
		return
	}

	h.store.Lock()
	patient, exists := h.store.GetPatient(patientID)
	if !exists {
		h.store.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	// Check access for volunteers
	if role == RoleVolunteer && patient.CreatedBy != userID {
		h.store.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	reminder := &models.Reminder{
		ID:             h.generateID(),
		Title:          req.Title,
		Description:    req.Description,
		DueDate:        req.DueDate,
		Priority:       req.Priority,
		Completed:      false,
		Recurrence:     req.Recurrence,
		Notified:       false,
		Attachments:    req.Attachments,
		DeliveryStatus: models.DeliveryStatusPending,
	}
	patient.Reminders = append(patient.Reminders, reminder)
	h.store.Unlock()

	h.store.SaveData()

	if h.logger != nil {
		h.logger.Info("Reminder created",
			"reminder_id", reminder.ID,
			"patient_id", patientID,
			"user_id", userID,
		)
	}

	c.JSON(http.StatusCreated, reminder)
}

// Update handles PUT /patients/:id/reminders/:reminderId
func (h *ReminderHandler) Update(c *gin.Context) {
	patientID := c.Param("id")
	reminderID := c.Param("reminderId")
	userID := c.GetString("userID")
	role := c.GetString("role")

	var req UpdateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate attachments count
	if len(req.Attachments) > MaxAttachments {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      fmt.Sprintf("Maksimal %d konten yang dapat dilampirkan", MaxAttachments),
			"code":       "MAX_ATTACHMENTS_EXCEEDED",
			"max_count":  MaxAttachments,
			"actual":     len(req.Attachments),
		})
		return
	}

	// Validate attachments content
	if err := h.validateAttachments(req.Attachments); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  "INVALID_ATTACHMENT",
		})
		return
	}

	h.store.Lock()
	patient, exists := h.store.GetPatient(patientID)
	if !exists {
		h.store.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	// Check access for volunteers
	if role == RoleVolunteer && patient.CreatedBy != userID {
		h.store.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	for _, r := range patient.Reminders {
		if r.ID == reminderID {
			if req.Title != "" {
				r.Title = req.Title
			}
			r.Description = req.Description
			r.DueDate = req.DueDate
			r.Priority = req.Priority
			r.Recurrence = req.Recurrence
			if req.Attachments != nil {
				r.Attachments = req.Attachments
			}
			if req.DueDate != "" && req.DueDate != r.DueDate {
				r.Notified = false
			}
			h.store.Unlock()
			h.store.SaveData()

			if h.logger != nil {
				h.logger.Info("Reminder updated",
					"reminder_id", reminderID,
					"patient_id", patientID,
					"user_id", userID,
				)
			}

			c.JSON(http.StatusOK, r)
			return
		}
	}
	h.store.Unlock()
	c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
}

// Toggle handles POST /patients/:id/reminders/:reminderId/toggle
func (h *ReminderHandler) Toggle(c *gin.Context) {
	patientID := c.Param("id")
	reminderID := c.Param("reminderId")
	userID := c.GetString("userID")
	role := c.GetString("role")

	h.store.Lock()
	patient, exists := h.store.GetPatient(patientID)
	if !exists {
		h.store.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	// Check access for volunteers
	if role == RoleVolunteer && patient.CreatedBy != userID {
		h.store.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	for _, r := range patient.Reminders {
		if r.ID == reminderID {
			r.Completed = !r.Completed
			if !r.Completed {
				r.Notified = false
			}
			h.store.Unlock()
			h.store.SaveData()

			if h.logger != nil {
				h.logger.Info("Reminder toggled",
					"reminder_id", reminderID,
					"patient_id", patientID,
					"completed", r.Completed,
				)
			}

			c.JSON(http.StatusOK, r)
			return
		}
	}
	h.store.Unlock()
	c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
}

// Delete handles DELETE /patients/:id/reminders/:reminderId
func (h *ReminderHandler) Delete(c *gin.Context) {
	patientID := c.Param("id")
	reminderID := c.Param("reminderId")
	userID := c.GetString("userID")
	role := c.GetString("role")

	h.store.Lock()
	patient, exists := h.store.GetPatient(patientID)
	if !exists {
		h.store.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	// Check access for volunteers
	if role == RoleVolunteer && patient.CreatedBy != userID {
		h.store.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	for i, r := range patient.Reminders {
		if r.ID == reminderID {
			patient.Reminders = append(patient.Reminders[:i], patient.Reminders[i+1:]...)
			h.store.Unlock()
			h.store.SaveData()

			if h.logger != nil {
				h.logger.Info("Reminder deleted",
					"reminder_id", reminderID,
					"patient_id", patientID,
					"user_id", userID,
				)
			}

			c.JSON(http.StatusOK, gin.H{"message": "reminder deleted"})
			return
		}
	}
	h.store.Unlock()
	c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
}

// Send handles POST /api/patients/:id/reminders/:reminderId/send
func (h *ReminderHandler) Send(c *gin.Context) {
	patientID := c.Param("id")
	reminderID := c.Param("reminderId")
	userID := c.GetString("userID")
	role := c.GetString("role")

	// 1. Get patient and validate access
	h.store.Lock()
	patient, exists := h.store.GetPatient(patientID)
	if !exists {
		h.store.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found", "code": "PATIENT_NOT_FOUND"})
		return
	}

	if role == RoleVolunteer && patient.CreatedBy != userID {
		h.store.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions", "code": "FORBIDDEN"})
		return
	}

	// 2. Find reminder
	var reminder *models.Reminder
	for _, r := range patient.Reminders {
		if r.ID == reminderID {
			reminder = r
			break
		}
	}
	if reminder == nil {
		h.store.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found", "code": "REMINDER_NOT_FOUND"})
		return
	}

	// 3. Validate phone number
	phoneResult := utils.ValidatePhoneNumber(patient.Phone)
	if !phoneResult.Valid {
		h.store.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nomor WhatsApp tidak valid",
			"code":  "INVALID_PHONE",
		})
		return
	}

	// 4. Check if already sending
	if reminder.DeliveryStatus == models.DeliveryStatusSending {
		h.store.Unlock()
		c.JSON(http.StatusConflict, gin.H{
			"error": "Reminder sedang dalam proses pengiriman",
			"code":  "ALREADY_SENDING",
		})
		return
	}

	// 5. Check quiet hours - schedule for later if in quiet hours
	now := time.Now()
	if utils.IsQuietHours(now, &h.config.QuietHours) {
		scheduledTime := utils.GetNextActiveTime(now, &h.config.QuietHours)
		reminder.DeliveryStatus = models.DeliveryStatusScheduled
		reminder.ScheduledDeliveryAt = scheduledTime.Format(time.RFC3339)
		h.store.Unlock()
		h.store.SaveData()

		if h.logger != nil {
			h.logger.Info("Reminder scheduled for quiet hours",
				"reminder_id", reminderID,
				"patient_id", patientID,
				"scheduled_at", reminder.ScheduledDeliveryAt,
			)
		}

		c.JSON(http.StatusOK, gin.H{
			"data":         reminder,
			"message":      "Reminder dijadwalkan untuk dikirim jam 06:00",
			"scheduled":    true,
			"scheduled_at": reminder.ScheduledDeliveryAt,
		})
		return
	}

	// 6. Update status to sending (optimistic) - active hours
	reminder.DeliveryStatus = models.DeliveryStatusSending
	// Capture sentAt timestamp before GOWA call for accuracy
	sentAt := time.Now().UTC()
	h.store.Unlock()
	h.store.SaveData()

	// 7. Format message
	message := h.formatReminderMessage(reminder, patient)

	// 8. Send via GOWA (outside lock)
	whatsappPhone := utils.FormatWhatsAppNumber(patient.Phone)
	response, err := h.gowaClient.SendMessage(whatsappPhone, message)

	// 9. Update status based on result
	h.store.Lock()
	if err != nil {
		// Check if circuit breaker is open - queue for retry (NFR-I2)
		if h.gowaClient.GetCircuitBreakerState() == "open" {
			reminder.DeliveryStatus = models.DeliveryStatusQueued
			reminder.DeliveryErrorMessage = "GOWA sedang tidak tersedia. Coba lagi nanti."
			reminder.RetryCount++

			h.store.Unlock()
			h.store.SaveData()

			if h.logger != nil {
				h.logger.Warn("Reminder queued for retry - circuit breaker open",
					"reminder_id", reminderID,
					"retry_count", reminder.RetryCount,
					"phone", utils.MaskPhone(whatsappPhone),
				)
			}

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":       reminder.DeliveryErrorMessage,
				"code":        "GOWA_UNAVAILABLE",
				"queued":      true,
				"retry_count": reminder.RetryCount,
			})
			return
		}

		// Check if error is retryable and we haven't exceeded max attempts
		if services.ShouldRetry(err) && reminder.RetryCount < h.config.Retry.MaxAttempts {
			// Schedule retry with exponential backoff
			retryDelay := services.GetRetryDelay(reminder.RetryCount, h.config.Retry.Delays)
			nextRetryTime := time.Now().UTC().Add(retryDelay)
			reminder.DeliveryStatus = models.DeliveryStatusRetrying
			reminder.ScheduledDeliveryAt = nextRetryTime.Format(time.RFC3339)
			reminder.DeliveryErrorMessage = err.Error()
			reminder.RetryCount++

			h.store.Unlock()
			h.store.SaveData()

			if h.logger != nil {
				h.logger.Info("Reminder scheduled for retry - transient failure",
					"reminder_id", reminderID,
					"patient_id", patientID,
					"retry_count", reminder.RetryCount,
					"next_retry_at", nextRetryTime.Format(time.RFC3339),
					"error", err.Error(),
				)
			}

			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":          "Pengiriman gagal, akan dicoba lagi",
				"code":           "RETRY_SCHEDULED",
				"retrying":       true,
				"retry_count":    reminder.RetryCount,
				"next_retry_at":  nextRetryTime.Format(time.RFC3339),
			})
			return
		}

		// Max retries exceeded or non-retryable error
		reminder.DeliveryStatus = models.DeliveryStatusFailed
		reminder.DeliveryErrorMessage = err.Error()

		h.store.Unlock()
		h.store.SaveData()

		if h.logger != nil {
			h.logger.Error("Failed to send reminder - max retries or non-retryable",
				"reminder_id", reminderID,
				"phone", utils.MaskPhone(whatsappPhone),
				"retry_count", reminder.RetryCount,
				"error", err.Error(),
			)
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":        reminder.DeliveryErrorMessage,
			"code":         "GOWA_ERROR",
			"retry_count":  reminder.RetryCount,
		})
		return
	}

	// Success - use captured timestamp for accuracy
	reminder.DeliveryStatus = models.DeliveryStatusSent
	reminder.GOWAMessageID = response.MessageID
	reminder.MessageSentAt = sentAt.Format(time.RFC3339)
	reminder.DeliveryErrorMessage = ""
	h.store.Unlock()
	h.store.SaveData()

	// Broadcast SSE event for real-time UI updates
	if h.sseHandler != nil {
		h.sseHandler.BroadcastDeliveryStatusUpdate(
			reminderID,
			string(models.DeliveryStatusSent),
			sentAt.Format(time.RFC3339),
		)
	}

	if h.logger != nil {
		h.logger.Info("Reminder sent successfully",
			"reminder_id", reminderID,
			"patient_id", patientID,
			"phone", utils.MaskPhone(whatsappPhone),
			"gowa_message_id", response.MessageID,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    reminder,
		"message": "Reminder berhasil dikirim",
	})
}

// buildContentAttachments builds ContentAttachment slice from reminder attachments
// Looks up article/video content from content store to get excerpts and URLs
// Sorts attachments: articles first, then videos (AC #3)
func (h *ReminderHandler) buildContentAttachments(reminder *models.Reminder) []utils.ContentAttachment {
	var attachments []utils.ContentAttachment

	for _, att := range reminder.Attachments {
		contentAtt := utils.ContentAttachment{
			Type:  att.Type,
			Title: att.Title,
		}

		if h.contentStore != nil {
			if att.Type == "article" {
				// Look up article for excerpt
				h.contentStore.Articles.Mu.RLock()
				if article, exists := h.contentStore.Articles.Articles[att.ID]; exists {
					contentAtt.Excerpt = article.Excerpt
					// Generate article URL from slug
					if article.Slug != "" {
						contentAtt.URL = fmt.Sprintf("https://prima.app/artikel/%s", article.Slug)
					}
				} else {
					// Article not found - use fallback text
					contentAtt.Excerpt = "Konten tidak tersedia"
				}
				h.contentStore.Articles.Mu.RUnlock()
			} else if att.Type == "video" {
				// Look up video for YouTube ID
				h.contentStore.Videos.Mu.RLock()
				if video, exists := h.contentStore.Videos.Videos[att.ID]; exists {
					// Generate YouTube URL from YouTube ID
					if video.YouTubeID != "" {
						contentAtt.URL = fmt.Sprintf("https://youtube.com/watch?v=%s", video.YouTubeID)
					}
				}
				h.contentStore.Videos.Mu.RUnlock()
			}
		}

		// Fallback to attachment URL if content store lookup didn't provide one
		if contentAtt.URL == "" && att.URL != "" {
			contentAtt.URL = att.URL
		}

		attachments = append(attachments, contentAtt)
	}

	// Sort attachments: articles first, then videos (AC #3 requirement)
	sort.Slice(attachments, func(i, j int) bool {
		if attachments[i].Type == attachments[j].Type {
			// Same type - maintain original order based on reminder.Attachments
			return false
		}
		return attachments[i].Type == "article"
	})

	return attachments
}

// formatReminderMessage creates the WhatsApp message content with excerpts
func (h *ReminderHandler) formatReminderMessage(reminder *models.Reminder, patient *models.Patient) string {
	disclaimerEnabled := h.config.Disclaimer.Enabled != nil && *h.config.Disclaimer.Enabled

	// Build content attachments with excerpts
	contentAttachments := h.buildContentAttachments(reminder)

	return utils.FormatReminderMessageWithExcerpts(utils.ReminderMessageParams{
		PatientName:         patient.Name,
		ReminderTitle:       reminder.Title,
		ReminderDescription: reminder.Description,
		DisclaimerText:      h.config.Disclaimer.Text,
		DisclaimerEnabled:   disclaimerEnabled,
	}, contentAttachments)
}

// DefaultIDGenerator generates a unique ID using timestamp and random string
func DefaultIDGenerator() string {
	return time.Now().Format("20060102150405") + "-" + reminderRandomString(8)
}

// reminderRandomString generates a cryptographically secure random string of n characters
func reminderRandomString(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to time-based if crypto/rand fails (should never happen in practice)
		for i := range bytes {
			bytes[i] = byte(time.Now().UnixNano() + int64(i))
		}
	}
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := range bytes {
		bytes[i] = letters[int(bytes[i])%len(letters)]
	}
	return string(bytes)
}

// GenerateSecureID generates a cryptographically secure unique ID
func GenerateSecureID() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based if crypto/rand fails
		return time.Now().Format("20060102150405") + "-" + reminderRandomString(8)
	}
	return time.Now().Format("20060102150405") + "-" + hex.EncodeToString(bytes)[:8]
}

// GetReminderStatus handles GET /api/reminders/:id/status
func (h *ReminderHandler) GetReminderStatus(c *gin.Context) {
	reminderID := c.Param("id")

	h.store.RLock()
	defer h.store.RUnlock()

	// Find the reminder across all patients
	var foundReminder *models.Reminder
	var patientID string

	for pid, patient := range h.store.Patients {
		for _, r := range patient.Reminders {
			if r.ID == reminderID {
				foundReminder = r
				patientID = pid
				break
			}
		}
		if foundReminder != nil {
			break
		}
	}

	if foundReminder == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found", "code": "REMINDER_NOT_FOUND"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":              foundReminder.ID,
			"title":           foundReminder.Title,
			"delivery_status": foundReminder.DeliveryStatus,
			"retry_count":     foundReminder.RetryCount,
			"max_attempts":    h.config.Retry.MaxAttempts,
			"error_message":   foundReminder.DeliveryErrorMessage,
			"scheduled_at":    foundReminder.ScheduledDeliveryAt,
			"sent_at":         foundReminder.MessageSentAt,
			"gowa_message_id": foundReminder.GOWAMessageID,
			"patient_id":      patientID,
		},
	})
}

// RetryReminder handles POST /api/reminders/:id/retry
func (h *ReminderHandler) RetryReminder(c *gin.Context) {
	reminderID := c.Param("id")
	userID := c.GetString("userID")
	role := c.GetString("role")

	// 1. Find reminder across all patients
	h.store.Lock()
	var patient *models.Patient
	var reminder *models.Reminder

	for _, p := range h.store.Patients {
		for _, r := range p.Reminders {
			if r.ID == reminderID {
				patient = p
				reminder = r
				break
			}
		}
		if reminder != nil {
			break
		}
	}

	if reminder == nil {
		h.store.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found", "code": "REMINDER_NOT_FOUND"})
		return
	}

	// 2. Check RBAC - volunteers can only retry their own reminders
	if role == RoleVolunteer && patient.CreatedBy != userID {
		h.store.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions", "code": "FORBIDDEN"})
		return
	}

	// 3. Validate reminder is in failed state
	if reminder.DeliveryStatus != models.DeliveryStatusFailed {
		h.store.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Only failed reminders can be retried",
			"code":  "INVALID_STATUS",
		})
		return
	}

	// 4. Validate phone number
	phoneResult := utils.ValidatePhoneNumber(patient.Phone)
	if !phoneResult.Valid {
		h.store.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nomor WhatsApp tidak valid",
			"code":  "INVALID_PHONE",
		})
		return
	}

	// 5. Check circuit breaker state
	if !h.gowaClient.IsAvailable() {
		// Queue reminder for retry when circuit breaker resets
		reminder.DeliveryStatus = models.DeliveryStatusQueued
		reminder.DeliveryErrorMessage = "GOWA sedang tidak tersedia. Akan dicoba lagi."
		h.store.Unlock()
		h.store.SaveData()

		if h.logger != nil {
			h.logger.Warn("Manual retry queued - circuit breaker open",
				"reminder_id", reminderID,
				"phone", utils.MaskPhone(utils.FormatWhatsAppNumber(patient.Phone)),
			)
		}

		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"reminder_id": reminderID,
				"status":      "queued",
			},
			"message": "GOWA sedang tidak tersedia. Reminder akan dikirim otomatis saat layanan kembali normal.",
		})
		return
	}

	// 6. Update status to sending (optimistic)
	reminder.DeliveryStatus = models.DeliveryStatusSending
	reminder.DeliveryErrorMessage = ""
	sentAt := time.Now().UTC()
	h.store.Unlock()
	h.store.SaveData()

	// 7. Format message
	message := h.formatReminderMessage(reminder, patient)

	// 8. Send via GOWA (outside lock)
	whatsappPhone := utils.FormatWhatsAppNumber(patient.Phone)
	response, err := h.gowaClient.SendMessage(whatsappPhone, message)

	// 9. Update status based on result
	h.store.Lock()
	if err != nil {
		// Retry failed
		reminder.DeliveryStatus = models.DeliveryStatusFailed
		reminder.DeliveryErrorMessage = err.Error()
		h.store.Unlock()
		h.store.SaveData()

		if h.logger != nil {
			h.logger.Error("Failed to retry reminder",
				"reminder_id", reminderID,
				"phone", utils.MaskPhone(whatsappPhone),
				"error", err.Error(),
			)
		}

		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": reminder.DeliveryErrorMessage,
			"code":  "GOWA_ERROR",
		})
		return
	}

	// Success - reset retry count
	reminder.DeliveryStatus = models.DeliveryStatusSent
	reminder.GOWAMessageID = response.MessageID
	reminder.MessageSentAt = sentAt.Format(time.RFC3339)
	reminder.DeliveryErrorMessage = ""
	reminder.RetryCount = 0 // Reset retry count on manual retry
	h.store.Unlock()
	h.store.SaveData()

	// Broadcast SSE event for real-time UI updates
	if h.sseHandler != nil {
		h.sseHandler.BroadcastDeliveryStatusUpdate(
			reminderID,
			string(models.DeliveryStatusSent),
			sentAt.Format(time.RFC3339),
		)
	}

	if h.logger != nil {
		h.logger.Info("Reminder retried successfully",
			"reminder_id", reminderID,
			"phone", utils.MaskPhone(whatsappPhone),
			"gowa_message_id", response.MessageID,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"reminder_id": reminderID,
			"status":      "sent",
			"message_id":  response.MessageID,
		},
		"message": "Reminder berhasil dikirim ulang",
	})
}

// ReminderHistoryResponse represents a reminder in the history API response
type ReminderHistoryResponse struct {
	ID              string              `json:"id"`
	Title           string              `json:"title"`
	Message         string              `json:"message"`
	MessagePreview  string              `json:"message_preview"`
	ScheduledAt     string              `json:"scheduled_at"`
	DeliveryStatus  string              `json:"delivery_status"`
	DeliveryError   string              `json:"delivery_error,omitempty"`
	SentAt          string              `json:"sent_at,omitempty"`
	DeliveredAt     string              `json:"delivered_at,omitempty"`
	ReadAt          string              `json:"read_at,omitempty"`
	CancelledAt     string              `json:"cancelled_at,omitempty"`
	Attachments     []models.Attachment `json:"attachments"`
	AttachmentCount int                 `json:"attachment_count"`
}

// PaginationResponse represents pagination info
type PaginationResponse struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	HasMore bool `json:"has_more"`
}

// GetPatientReminders handles GET /api/patients/:id/reminders
// Supports history=true query parameter to return all reminders including cancelled
func (h *ReminderHandler) GetPatientReminders(c *gin.Context) {
	patientID := c.Param("id")
	userID := c.GetString("userID")
	role := c.GetString("role")
	history := c.Query("history") == "true"

	// Parse pagination params
	page := 1
	limit := 20
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	h.store.RLock()
	patient, exists := h.store.Patients[patientID]
	if !exists {
		h.store.RUnlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	// Check access for volunteers
	if role == RoleVolunteer && patient.CreatedBy != userID {
		h.store.RUnlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	// Build reminder list based on history flag
	var reminders []ReminderHistoryResponse
	allReminders := patient.Reminders

	for _, r := range allReminders {
		// In non-history mode, exclude cancelled reminders
		if !history && r.DeliveryStatus == models.DeliveryStatusCancelled {
			continue
		}

		// Build message preview (first 50 chars of description or empty)
		messagePreview := ""
		if len(r.Description) > 50 {
			messagePreview = r.Description[:50] + "..."
		} else {
			messagePreview = r.Description
		}

		reminderResp := ReminderHistoryResponse{
			ID:              r.ID,
			Title:           r.Title,
			Message:         r.Description,
			MessagePreview:  messagePreview,
			ScheduledAt:     r.ScheduledDeliveryAt,
			DeliveryStatus:  r.DeliveryStatus,
			DeliveryError:   r.DeliveryErrorMessage,
			SentAt:          r.MessageSentAt,
			DeliveredAt:     r.DeliveredAt,
			ReadAt:          r.ReadAt,
			CancelledAt:     r.CancelledAt,
			Attachments:     r.Attachments,
			AttachmentCount: len(r.Attachments),
		}
		reminders = append(reminders, reminderResp)
	}

	// Sort by scheduled_at descending (most recent first)
	sort.Slice(reminders, func(i, j int) bool {
		return reminders[i].ScheduledAt > reminders[j].ScheduledAt
	})

	total := len(reminders)
	hasMore := (page * limit) < total

	// Apply pagination
	startIdx := (page - 1) * limit
	if startIdx >= total {
		reminders = []ReminderHistoryResponse{}
		hasMore = false
	} else {
		endIdx := startIdx + limit
		if endIdx > total {
			endIdx = total
		}
		reminders = reminders[startIdx:endIdx]
	}

	h.store.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"data":    reminders,
		"message": "Success",
		"pagination": PaginationResponse{
			Page:    page,
			Limit:   limit,
			Total:   total,
			HasMore: hasMore,
		},
	})
}

// CancelReminder handles POST /api/reminders/:id/cancel
func (h *ReminderHandler) CancelReminder(c *gin.Context) {
	reminderID := c.Param("id")
	userID := c.GetString("userID")
	role := c.GetString("role")

	// Find reminder across all patients
	h.store.Lock()
	var patient *models.Patient
	var reminder *models.Reminder

	for _, p := range h.store.Patients {
		for _, r := range p.Reminders {
			if r.ID == reminderID {
				patient = p
				reminder = r
				break
			}
		}
		if reminder != nil {
			break
		}
	}

	if reminder == nil {
		h.store.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "Reminder not found", "code": "REMINDER_NOT_FOUND"})
		return
	}

	// Check RBAC - volunteers can only cancel their own reminders
	if role == RoleVolunteer && patient.CreatedBy != userID {
		h.store.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions", "code": "FORBIDDEN"})
		return
	}

	// Validate status - only pending or scheduled can be cancelled
	if reminder.DeliveryStatus != models.DeliveryStatusPending &&
		reminder.DeliveryStatus != models.DeliveryStatusScheduled {
		h.store.Unlock()
		c.JSON(http.StatusBadRequest, gin.H{
			"error":          "Reminder tidak dapat dibatalkan",
			"code":           "CANNOT_CANCEL",
			"current_status": reminder.DeliveryStatus,
		})
		return
	}

	// Cancel the reminder
	previousStatus := reminder.DeliveryStatus
	reminder.DeliveryStatus = models.DeliveryStatusCancelled
	reminder.CancelledAt = time.Now().UTC().Format(time.RFC3339)
	reminder.CancelledBy = userID

	h.store.Unlock()
	h.store.SaveData()

	// Log for audit
	if h.logger != nil {
		h.logger.Info("Reminder cancelled",
			"reminder_id", reminderID,
			"patient_id", patient.ID,
			"cancelled_by", userID,
			"previous_status", previousStatus,
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":           reminder.ID,
			"title":        reminder.Title,
			"status":       reminder.DeliveryStatus,
			"cancelled_at": reminder.CancelledAt,
		},
		"message": "Reminder berhasil dibatalkan",
	})
}
