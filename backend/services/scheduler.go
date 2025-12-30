package services

import (
	"log/slog"
	"sync"
	"time"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/utils"
)

// SSEHandler interface for broadcasting delivery status updates
type SSEHandler interface {
	BroadcastDeliveryStatusUpdate(reminderID, status, timestamp string)
}

// ReminderScheduler handles automatic sending of scheduled reminders
type ReminderScheduler struct {
	store      *models.PatientStore
	gowaClient *GOWAClient
	config     *config.Config
	logger     *slog.Logger
	sseHandler SSEHandler // SSE handler for broadcasting delivery status updates
	stopCh     chan struct{}
	wg         sync.WaitGroup
	interval   time.Duration
}

// NewReminderScheduler creates a new reminder scheduler
func NewReminderScheduler(store *models.PatientStore, gowaClient *GOWAClient, cfg *config.Config, logger *slog.Logger) *ReminderScheduler {
	return &ReminderScheduler{
		store:      store,
		gowaClient: gowaClient,
		config:     cfg,
		logger:     logger,
		sseHandler: nil, // Will be set via SetSSEHandler
		stopCh:     make(chan struct{}),
		interval:   1 * time.Minute,
	}
}

// SetSSEHandler sets the SSE handler for broadcasting delivery status updates
func (s *ReminderScheduler) SetSSEHandler(sseHandler SSEHandler) {
	s.sseHandler = sseHandler
}

// Start begins the scheduler goroutine
func (s *ReminderScheduler) Start() {
	s.wg.Add(1)
	go s.run()
}

// Stop gracefully stops the scheduler
func (s *ReminderScheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()

	if s.logger != nil {
		s.logger.Info("Reminder scheduler stopped")
	}
}

// run is the main scheduler loop
func (s *ReminderScheduler) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Check immediately on start for any pending scheduled reminders
	s.processScheduledReminders()

	for {
		select {
		case <-ticker.C:
			s.processScheduledReminders()
		case <-s.stopCh:
			return
		}
	}
}

// processScheduledReminders finds and sends all due scheduled reminders
func (s *ReminderScheduler) processScheduledReminders() {
	now := time.Now().UTC()

	// Collect reminders to send (read lock)
	s.store.RLock()
	var toSend []struct {
		patientID  string
		patient    *models.Patient
		reminderID string
		reminder   *models.Reminder
	}

	reminderCount := 0
	for patientID, patient := range s.store.Patients {
		for _, reminder := range patient.Reminders {
			reminderCount++
			// Handle scheduled reminders (quiet hours)
			if reminder.DeliveryStatus == models.DeliveryStatusScheduled {
				if reminder.ScheduledDeliveryAt == "" {
					continue
				}

				scheduledTime, err := time.Parse(time.RFC3339, reminder.ScheduledDeliveryAt)
				if err != nil {
					if s.logger != nil {
						s.logger.Error("Failed to parse scheduled delivery time",
							"reminder_id", reminder.ID,
							"scheduled_at", reminder.ScheduledDeliveryAt,
							"error", err.Error(),
						)
					}
					continue
				}

				// Check if it's time to send
				if now.After(scheduledTime) || now.Equal(scheduledTime) {
					toSend = append(toSend, struct {
						patientID  string
						patient    *models.Patient
						reminderID string
						reminder   *models.Reminder
					}{
						patientID:  patientID,
						patient:    patient,
						reminderID: reminder.ID,
						reminder:   reminder,
					})
				}
				continue
			}

			// Handle retrying reminders
			if reminder.DeliveryStatus == models.DeliveryStatusRetrying {
				if reminder.ScheduledDeliveryAt == "" {
					continue
				}

				scheduledTime, err := time.Parse(time.RFC3339, reminder.ScheduledDeliveryAt)
				if err != nil {
					if s.logger != nil {
						s.logger.Error("Failed to parse retry scheduled time",
							"reminder_id", reminder.ID,
							"scheduled_at", reminder.ScheduledDeliveryAt,
							"error", err.Error(),
						)
					}
					continue
				}

				// Check if it's time to retry
				if now.After(scheduledTime) || now.Equal(scheduledTime) {
					toSend = append(toSend, struct {
						patientID  string
						patient    *models.Patient
						reminderID string
						reminder   *models.Reminder
					}{
						patientID:  patientID,
						patient:    patient,
						reminderID: reminder.ID,
						reminder:   reminder,
					})
				}
				continue
			}

			// Handle auto-send for due reminders (replaces old checkReminders goroutine)
			if (reminder.DeliveryStatus == "" || reminder.DeliveryStatus == models.DeliveryStatusPending) &&
				!reminder.Completed && !reminder.Notified && reminder.DueDate != "" {

				// Parse due date (support multiple formats for compatibility)
				var dueTime time.Time
				var err error

				// Try RFC3339 first (new format with timezone)
				dueTime, err = time.Parse(time.RFC3339, reminder.DueDate)
				if err != nil {
					// Try local time format (old format from checkReminders - no timezone)
					// IMPORTANT: Use ParseInLocation with time.Local to interpret as local time
					dueTime, err = time.ParseInLocation("2006-01-02T15:04", reminder.DueDate, time.Local)
					if err != nil {
						dueTime, err = time.ParseInLocation("2006-01-02T15:04:05", reminder.DueDate, time.Local)
						if err != nil {
							if s.logger != nil {
								s.logger.Error("Failed to parse due date",
									"reminder_id", reminder.ID,
									"due_date", reminder.DueDate,
									"error", err.Error(),
								)
							}
							continue
						}
					}
					// Convert to UTC for comparison with now (which is UTC)
					dueTime = dueTime.UTC()
				}

				// Check if due date has passed
				if now.After(dueTime) || now.Equal(dueTime) {
					// Only send if within reasonable window (24 hours) to avoid sending very old reminders
					if now.Before(dueTime.Add(24 * time.Hour)) {
						toSend = append(toSend, struct {
							patientID  string
							patient    *models.Patient
							reminderID string
							reminder   *models.Reminder
						}{
							patientID:  patientID,
							patient:    patient,
							reminderID: reminder.ID,
							reminder:   reminder,
						})
					} else {
						if s.logger != nil {
							s.logger.Warn("Reminder too old, skipping",
								"reminder_id", reminder.ID,
								"due_time", dueTime.Format(time.RFC3339),
							)
						}
					}
				}
				continue
			}
		}
	}
	s.store.RUnlock()

	// Send each reminder
	for _, item := range toSend {
		if item.reminder.DeliveryStatus == models.DeliveryStatusRetrying {
			s.processRetryReminder(item.patientID, item.patient, item.reminder)
		} else {
			s.sendScheduledReminder(item.patientID, item.patient, item.reminder)
		}
	}
}

// sendScheduledReminder sends a single scheduled reminder
func (s *ReminderScheduler) sendScheduledReminder(patientID string, patient *models.Patient, reminder *models.Reminder) {
	reminderID := reminder.ID

	// Validate phone number (no lock needed for reading)
	phoneResult := utils.ValidatePhoneNumber(patient.Phone)
	if !phoneResult.Valid {
		s.store.Lock()
		// Re-fetch reminder to ensure it still exists and is in correct state
		currentPatient, exists := s.store.Patients[patientID]
		if !exists {
			s.store.Unlock()
			return
		}
		currentReminder := findReminderByID(currentPatient, reminderID)
		if currentReminder == nil || currentReminder.DeliveryStatus != models.DeliveryStatusScheduled {
			s.store.Unlock()
			return
		}
		currentReminder.DeliveryStatus = models.DeliveryStatusFailed
		currentReminder.DeliveryErrorMessage = "Nomor WhatsApp tidak valid"
		s.store.Unlock()
		s.store.SaveData()

		if s.logger != nil {
			s.logger.Error("Scheduled reminder failed - invalid phone",
				"reminder_id", reminderID,
				"patient_id", patientID,
			)
		}
		return
	}

	// Update status to sending (with re-fetch)
	s.store.Lock()
	currentPatient, exists := s.store.Patients[patientID]
	if !exists {
		s.store.Unlock()
		return
	}
	currentReminder := findReminderByID(currentPatient, reminderID)
	// Allow both scheduled (quiet hours) and pending (auto-send) reminders
	if currentReminder == nil {
		s.store.Unlock()
		return
	}
	// Skip if already being processed or sent
	if currentReminder.DeliveryStatus != models.DeliveryStatusScheduled &&
		currentReminder.DeliveryStatus != models.DeliveryStatusPending &&
		currentReminder.DeliveryStatus != "" {
		s.store.Unlock()
		return
	}
	currentReminder.DeliveryStatus = models.DeliveryStatusSending
	// Capture current state for message formatting
	patientName := currentPatient.Name
	reminderTitle := currentReminder.Title
	reminderDescription := currentReminder.Description
	s.store.Unlock()
	s.store.SaveData()

	// Format message (using captured values, no lock needed)
	message := s.formatReminderMessageFromValues(reminderTitle, reminderDescription, patientName)

	// Send via GOWA (outside lock)
	whatsappPhone := utils.FormatWhatsAppNumber(patient.Phone)
	response, err := s.gowaClient.SendMessage(whatsappPhone, message)

	// Update status based on result (with re-fetch)
	s.store.Lock()
	currentPatient, exists = s.store.Patients[patientID]
	if !exists {
		s.store.Unlock()
		return
	}
	currentReminder = findReminderByID(currentPatient, reminderID)
	if currentReminder == nil {
		s.store.Unlock()
		return
	}

	if err != nil {
		// Check if circuit breaker is open - requeue
		if s.gowaClient.GetCircuitBreakerState() == "open" {
			currentReminder.DeliveryStatus = models.DeliveryStatusQueued
			currentReminder.DeliveryErrorMessage = "GOWA sedang tidak tersedia. Akan dicoba lagi."
			currentReminder.RetryCount++

			if s.logger != nil {
				s.logger.Warn("Scheduled reminder queued - circuit breaker open",
					"reminder_id", reminderID,
					"patient_id", patientID,
					"retry_count", currentReminder.RetryCount,
				)
			}
		} else {
			currentReminder.DeliveryStatus = models.DeliveryStatusFailed
			currentReminder.DeliveryErrorMessage = err.Error()

			if s.logger != nil {
				s.logger.Error("Scheduled reminder failed",
					"reminder_id", reminderID,
					"patient_id", patientID,
					"phone", utils.MaskPhone(whatsappPhone),
					"error", err.Error(),
				)
			}
		}
	} else {
		// Success
		currentReminder.DeliveryStatus = models.DeliveryStatusSent
		currentReminder.GOWAMessageID = response.MessageID
		sentAt := time.Now().UTC()
		currentReminder.MessageSentAt = sentAt.Format(time.RFC3339)
		currentReminder.DeliveryErrorMessage = ""
		currentReminder.ScheduledDeliveryAt = "" // Clear scheduled time

		// Broadcast SSE event for real-time UI updates (before unlock)
		if s.sseHandler != nil {
			s.sseHandler.BroadcastDeliveryStatusUpdate(
				reminderID,
				string(models.DeliveryStatusSent),
				sentAt.Format(time.RFC3339),
			)
		}
	}
	s.store.Unlock()
	s.store.SaveData()
}

// findReminderByID finds a reminder by ID in a patient's reminders
func findReminderByID(patient *models.Patient, reminderID string) *models.Reminder {
	for _, r := range patient.Reminders {
		if r.ID == reminderID {
			return r
		}
	}
	return nil
}

// formatReminderMessageFromValues creates the WhatsApp message content from values
func (s *ReminderScheduler) formatReminderMessageFromValues(title, description, patientName string) string {
	disclaimerEnabled := s.config.Disclaimer.Enabled != nil && *s.config.Disclaimer.Enabled

	return utils.FormatReminderMessage(utils.ReminderMessageParams{
		PatientName:         patientName,
		ReminderTitle:       title,
		ReminderDescription: description,
		DisclaimerText:      s.config.Disclaimer.Text,
		DisclaimerEnabled:   disclaimerEnabled,
	})
}

// processRetryReminder handles retrying a failed reminder
func (s *ReminderScheduler) processRetryReminder(patientID string, patient *models.Patient, reminder *models.Reminder) {
	reminderID := reminder.ID

	// Validate phone number
	phoneResult := utils.ValidatePhoneNumber(patient.Phone)
	if !phoneResult.Valid {
		s.store.Lock()
		currentPatient, exists := s.store.Patients[patientID]
		if !exists {
			s.store.Unlock()
			return
		}
		currentReminder := findReminderByID(currentPatient, reminderID)
		if currentReminder == nil || currentReminder.DeliveryStatus != models.DeliveryStatusRetrying {
			s.store.Unlock()
			return
		}
		currentReminder.DeliveryStatus = models.DeliveryStatusFailed
		currentReminder.DeliveryErrorMessage = "Nomor WhatsApp tidak valid"
		s.store.Unlock()
		s.store.SaveData()

		if s.logger != nil {
			s.logger.Error("Retry reminder failed - invalid phone",
				"reminder_id", reminderID,
				"patient_id", patientID,
				"retry_count", currentReminder.RetryCount,
			)
		}
		return
	}

	// Update status to sending
	s.store.Lock()
	currentPatient, exists := s.store.Patients[patientID]
	if !exists {
		s.store.Unlock()
		return
	}
	currentReminder := findReminderByID(currentPatient, reminderID)
	if currentReminder == nil || currentReminder.DeliveryStatus != models.DeliveryStatusRetrying {
		s.store.Unlock()
		return
	}
	currentReminder.DeliveryStatus = models.DeliveryStatusSending
	patientName := currentPatient.Name
	reminderTitle := currentReminder.Title
	reminderDescription := currentReminder.Description
	s.store.Unlock()

	// Format message
	message := s.formatReminderMessageFromValues(reminderTitle, reminderDescription, patientName)

	// Send via GOWA
	whatsappPhone := utils.FormatWhatsAppNumber(patient.Phone)
	response, err := s.gowaClient.SendMessage(whatsappPhone, message)

	// Update status based on result
	s.store.Lock()
	currentPatient, exists = s.store.Patients[patientID]
	if !exists {
		s.store.Unlock()
		return
	}
	currentReminder = findReminderByID(currentPatient, reminderID)
	if currentReminder == nil {
		s.store.Unlock()
		return
	}

	if err != nil {
		// Check if error is retryable
		if ShouldRetry(err) && currentReminder.RetryCount < s.config.Retry.MaxAttempts {
			// Schedule next retry
			retryDelay := GetRetryDelay(currentReminder.RetryCount, s.config.Retry.Delays)
			nextRetryTime := time.Now().UTC().Add(retryDelay)
			currentReminder.DeliveryStatus = models.DeliveryStatusRetrying
			currentReminder.ScheduledDeliveryAt = nextRetryTime.Format(time.RFC3339)
		} else {
			// Max retries exceeded or non-retryable error
			currentReminder.DeliveryStatus = models.DeliveryStatusFailed
			currentReminder.DeliveryErrorMessage = err.Error()

			if s.logger != nil {
				s.logger.Error("Retry reminder failed - max retries exceeded",
					"reminder_id", reminderID,
					"patient_id", patientID,
					"retry_count", currentReminder.RetryCount,
					"error", err.Error(),
				)
			}
		}
	} else {
		// Success
		currentReminder.DeliveryStatus = models.DeliveryStatusSent
		currentReminder.GOWAMessageID = response.MessageID
		sentAt := time.Now().UTC()
		currentReminder.MessageSentAt = sentAt.Format(time.RFC3339)
		currentReminder.DeliveryErrorMessage = ""
		currentReminder.ScheduledDeliveryAt = ""
		currentReminder.RetryCount = 0

		// Broadcast SSE event for real-time UI updates (before unlock)
		if s.sseHandler != nil {
			s.sseHandler.BroadcastDeliveryStatusUpdate(
				reminderID,
				string(models.DeliveryStatusSent),
				sentAt.Format(time.RFC3339),
			)
		}
	}
	s.store.Unlock()
	s.store.SaveData()
}

// SetInterval allows changing the check interval (useful for testing)
func (s *ReminderScheduler) SetInterval(interval time.Duration) {
	s.interval = interval
}
