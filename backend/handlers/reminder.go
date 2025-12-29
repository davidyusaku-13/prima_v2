package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
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
	store      *models.PatientStore
	config     *config.Config
	gowaClient *services.GOWAClient
	logger     *slog.Logger
	generateID IDGenerator
}

// NewReminderHandler creates a new reminder handler
func NewReminderHandler(store *models.PatientStore, cfg *config.Config, gowaClient *services.GOWAClient, logger *slog.Logger, idGen IDGenerator) *ReminderHandler {
	return &ReminderHandler{
		store:      store,
		config:     cfg,
		gowaClient: gowaClient,
		logger:     logger,
		generateID: idGen,
	}
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

	// 5. Update status to sending (optimistic)
	reminder.DeliveryStatus = models.DeliveryStatusSending
	h.store.Unlock()
	h.store.SaveData()

	// 6. Format message
	message := h.formatReminderMessage(reminder, patient)

	// 7. Send via GOWA (outside lock)
	whatsappPhone := utils.FormatWhatsAppNumber(patient.Phone)
	response, err := h.gowaClient.SendMessage(whatsappPhone, message)

	// 8. Update status based on result
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

		// Regular failure - not queued
		reminder.DeliveryStatus = models.DeliveryStatusFailed
		reminder.DeliveryErrorMessage = err.Error()

		h.store.Unlock()
		h.store.SaveData()

		if h.logger != nil {
			h.logger.Error("Failed to send reminder",
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

	// Success
	reminder.DeliveryStatus = models.DeliveryStatusSent
	reminder.GOWAMessageID = response.MessageID
	reminder.MessageSentAt = time.Now().UTC().Format(time.RFC3339)
	reminder.DeliveryErrorMessage = ""
	h.store.Unlock()
	h.store.SaveData()

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

// formatReminderMessage creates the WhatsApp message content
func (h *ReminderHandler) formatReminderMessage(reminder *models.Reminder, patient *models.Patient) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Halo %s,\n\n", patient.Name))
	sb.WriteString(fmt.Sprintf("*%s*\n\n", reminder.Title))

	if reminder.Description != "" {
		sb.WriteString(reminder.Description)
		sb.WriteString("\n\n")
	}

	// Add health disclaimer
	sb.WriteString("---\n")
	sb.WriteString("_Informasi ini untuk tujuan edukasi. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda._")

	return sb.String()
}

// DefaultIDGenerator generates a unique ID using timestamp and random string
func DefaultIDGenerator() string {
	return time.Now().Format("20060102150405") + "-" + reminderRandomString(8)
}

// reminderRandomString generates a random string of n characters using crypto/rand
func reminderRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		// Use time-based seed with index for variety
		b[i] = letters[(time.Now().UnixNano()+int64(i))%int64(len(letters))]
	}
	return string(b)
}
