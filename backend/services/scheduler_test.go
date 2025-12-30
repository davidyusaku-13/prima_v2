package services

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/davidyusaku-13/prima_v2/models"
)

func TestReminderScheduler_ProcessScheduledReminders(t *testing.T) {
	t.Run("sends due scheduled reminders", func(t *testing.T) {
		// Create mock GOWA server
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SendMessageResponse{
				Success:   true,
				MessageID: "msg-scheduled-123",
			})
		}))
		defer gowaServer.Close()

		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		enabled := true
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Test disclaimer",
				Enabled: &enabled,
			},
		}

		gowaClient := NewGOWAClient(GOWAConfig{
			Endpoint:         gowaServer.URL,
			User:             "testuser",
			Password:         "testpass",
			Timeout:          10 * time.Second,
			FailureThreshold: 5,
			CooldownDuration: 5 * time.Minute,
		}, logger)

		scheduler := NewReminderScheduler(store, gowaClient, cfg, logger)

		// Setup test data with a scheduled reminder that is due
		pastTime := time.Now().UTC().Add(-1 * time.Minute).Format(time.RFC3339)
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:                  "reminder-1",
					Title:               "Scheduled Reminder",
					Description:         "Test Description",
					DeliveryStatus:      models.DeliveryStatusScheduled,
					ScheduledDeliveryAt: pastTime,
				},
			},
		}

		// Process scheduled reminders
		scheduler.processScheduledReminders()

		// Verify reminder was sent
		reminder := store.Patients["patient-1"].Reminders[0]
		if reminder.DeliveryStatus != models.DeliveryStatusSent {
			t.Errorf("Expected delivery status 'sent', got '%s'", reminder.DeliveryStatus)
		}
		if reminder.GOWAMessageID != "msg-scheduled-123" {
			t.Errorf("Expected GOWA message ID 'msg-scheduled-123', got '%s'", reminder.GOWAMessageID)
		}
		if reminder.MessageSentAt == "" {
			t.Error("Expected MessageSentAt to be set")
		}
		if reminder.ScheduledDeliveryAt != "" {
			t.Error("Expected ScheduledDeliveryAt to be cleared after sending")
		}
	})

	t.Run("ignores non-scheduled reminders", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		enabled := true
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Test disclaimer",
				Enabled: &enabled,
			},
		}

		scheduler := NewReminderScheduler(store, nil, cfg, logger)

		// Setup test data with non-scheduled reminders
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Pending Reminder",
					DeliveryStatus: models.DeliveryStatusPending,
				},
				{
					ID:             "reminder-2",
					Title:          "Sent Reminder",
					DeliveryStatus: models.DeliveryStatusSent,
				},
				{
					ID:             "reminder-3",
					Title:          "Failed Reminder",
					DeliveryStatus: models.DeliveryStatusFailed,
				},
			},
		}

		// Process scheduled reminders - should not affect any
		scheduler.processScheduledReminders()

		// Verify statuses unchanged
		reminders := store.Patients["patient-1"].Reminders
		if reminders[0].DeliveryStatus != models.DeliveryStatusPending {
			t.Errorf("Expected pending reminder to remain pending, got '%s'", reminders[0].DeliveryStatus)
		}
		if reminders[1].DeliveryStatus != models.DeliveryStatusSent {
			t.Errorf("Expected sent reminder to remain sent, got '%s'", reminders[1].DeliveryStatus)
		}
		if reminders[2].DeliveryStatus != models.DeliveryStatusFailed {
			t.Errorf("Expected failed reminder to remain failed, got '%s'", reminders[2].DeliveryStatus)
		}
	})

	t.Run("does not send future scheduled reminders", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		enabled := true
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Test disclaimer",
				Enabled: &enabled,
			},
		}

		scheduler := NewReminderScheduler(store, nil, cfg, logger)

		// Setup test data with a scheduled reminder in the future
		futureTime := time.Now().UTC().Add(1 * time.Hour).Format(time.RFC3339)
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:                  "reminder-1",
					Title:               "Future Scheduled Reminder",
					DeliveryStatus:      models.DeliveryStatusScheduled,
					ScheduledDeliveryAt: futureTime,
				},
			},
		}

		// Process scheduled reminders - should not send future reminder
		scheduler.processScheduledReminders()

		// Verify reminder still scheduled
		reminder := store.Patients["patient-1"].Reminders[0]
		if reminder.DeliveryStatus != models.DeliveryStatusScheduled {
			t.Errorf("Expected future reminder to remain scheduled, got '%s'", reminder.DeliveryStatus)
		}
	})

	t.Run("handles invalid phone number", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		enabled := true
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Test disclaimer",
				Enabled: &enabled,
			},
		}

		scheduler := NewReminderScheduler(store, nil, cfg, logger)

		// Setup test data with invalid phone
		pastTime := time.Now().UTC().Add(-1 * time.Minute).Format(time.RFC3339)
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "invalid-phone",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:                  "reminder-1",
					Title:               "Scheduled Reminder",
					DeliveryStatus:      models.DeliveryStatusScheduled,
					ScheduledDeliveryAt: pastTime,
				},
			},
		}

		// Process scheduled reminders
		scheduler.processScheduledReminders()

		// Verify reminder failed due to invalid phone
		reminder := store.Patients["patient-1"].Reminders[0]
		if reminder.DeliveryStatus != models.DeliveryStatusFailed {
			t.Errorf("Expected delivery status 'failed' for invalid phone, got '%s'", reminder.DeliveryStatus)
		}
		if reminder.DeliveryErrorMessage == "" {
			t.Error("Expected error message to be set")
		}
	})
}

func TestReminderScheduler_StartStop(t *testing.T) {
	t.Run("starts and stops gracefully", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		enabled := true
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Test disclaimer",
				Enabled: &enabled,
			},
		}

		scheduler := NewReminderScheduler(store, nil, cfg, logger)
		scheduler.SetInterval(100 * time.Millisecond) // Short interval for testing

		// Start scheduler
		scheduler.Start()

		// Let it run briefly
		time.Sleep(150 * time.Millisecond)

		// Stop scheduler - should not hang
		done := make(chan struct{})
		go func() {
			scheduler.Stop()
			close(done)
		}()

		select {
		case <-done:
			// Success - stopped gracefully
		case <-time.After(2 * time.Second):
			t.Error("Scheduler did not stop within timeout")
		}
	})
}

func TestReminderScheduler_FormatReminderMessage(t *testing.T) {
	t.Run("formats message with disclaimer", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		enabled := true
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Test health disclaimer",
				Enabled: &enabled,
			},
		}

		scheduler := NewReminderScheduler(store, nil, cfg, logger)

		// Use the correct method formatReminderMessageFromValues
		message := scheduler.formatReminderMessageFromValues(
			"Test Reminder",  // title
			"Test description", // description
			"John Doe",        // patientName
		)

		if message == "" {
			t.Error("Expected non-empty message")
		}
		if !strings.Contains(message, "Halo John Doe") {
			t.Error("Message should contain patient name")
		}
		if !strings.Contains(message, "*Test Reminder*") {
			t.Error("Message should contain reminder title")
		}
		if !strings.Contains(message, "Test description") {
			t.Error("Message should contain description")
		}
		if !strings.Contains(message, "Test health disclaimer") {
			t.Error("Message should contain disclaimer")
		}
	})
}
