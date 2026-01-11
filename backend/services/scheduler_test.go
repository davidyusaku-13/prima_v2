package services

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/utils"
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

func TestReminderScheduler_ProcessScheduledReminders_WithAttachments(t *testing.T) {
	t.Run("sends reminder with article and video attachments", func(t *testing.T) {
		// Track the message sent to GOWA
		var sentMessage string
		var mu sync.Mutex

		// Create mock GOWA server that captures the message
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Phone   string `json:"phone"`
				Message string `json:"message"`
			}
			json.NewDecoder(r.Body).Decode(&req)

			mu.Lock()
			sentMessage = req.Message
			mu.Unlock()

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SendMessageResponse{
				Success:   true,
				MessageID: "msg-with-attachments-123",
			})
		}))
		defer gowaServer.Close()

		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		enabled := true
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Test health disclaimer",
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

		// Create mock article store
		articleStore := &models.ArticleStore{
			Articles: map[string]*models.Article{
				"article-1": {
					ID:      "article-1",
					Title:   "Test Article",
					Slug:    "test-article",
					Excerpt: "This is a test article excerpt",
				},
			},
		}

		// Create mock video store
		videoStore := &models.VideoStore{
			Videos: map[string]*models.Video{
				"video-1": {
					ID:        "video-1",
					Title:     "Test Video",
					YouTubeID: "abc123xyz",
				},
			},
		}

		scheduler := NewReminderScheduler(store, gowaClient, cfg, logger)
		scheduler.SetContentStores(articleStore, videoStore)

		// Add patient with reminder that has attachments
		dueTime := time.Now().UTC().Add(-1 * time.Minute) // Due 1 minute ago
		store.Lock()
		store.Patients["patient-1"] = &models.Patient{
			ID:    "patient-1",
			Name:  "John Doe",
			Phone: "+6281234567890",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder With Attachments",
					Description:    "This is a test reminder",
					DueDate:        dueTime.Format(time.RFC3339),
					DeliveryStatus: models.DeliveryStatusPending,
					Attachments: []models.Attachment{
						{ID: "article-1", Type: "article", Title: "Test Article"},
						{ID: "video-1", Type: "video", Title: "Test Video"},
					},
				},
			},
		}
		store.Unlock()

		// Process scheduled reminders
		scheduler.processScheduledReminders()

		// Wait a bit for async processing
		time.Sleep(100 * time.Millisecond)

		// Verify the message was sent with attachments
		mu.Lock()
		msg := sentMessage
		mu.Unlock()

		if msg == "" {
			t.Fatal("Expected message to be sent")
		}

		// Check message contains patient name
		if !strings.Contains(msg, "Halo John Doe") {
			t.Error("Message should contain patient name")
		}

		// Check message contains reminder title
		if !strings.Contains(msg, "*Test Reminder With Attachments*") {
			t.Error("Message should contain reminder title")
		}

		// Check message contains article content
		if !strings.Contains(msg, "📖 Test Article") {
			t.Error("Message should contain article title with emoji")
		}
		if !strings.Contains(msg, "This is a test article excerpt") {
			t.Error("Message should contain article excerpt")
		}
		if !strings.Contains(msg, "https://prima.app/artikel/test-article") {
			t.Error("Message should contain article URL")
		}

		// Check message contains video content
		if !strings.Contains(msg, "🎬 Test Video") {
			t.Error("Message should contain video title with emoji")
		}
		if !strings.Contains(msg, "https://youtube.com/watch?v=abc123xyz") {
			t.Error("Message should contain YouTube URL")
		}

		// Check message contains disclaimer
		if !strings.Contains(msg, "Test health disclaimer") {
			t.Error("Message should contain disclaimer")
		}
	})

	t.Run("sends reminder without attachments", func(t *testing.T) {
		var sentMessage string
		var mu sync.Mutex

		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Phone   string `json:"phone"`
				Message string `json:"message"`
			}
			json.NewDecoder(r.Body).Decode(&req)

			mu.Lock()
			sentMessage = req.Message
			mu.Unlock()

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SendMessageResponse{
				Success:   true,
				MessageID: "msg-no-attachments-123",
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

		// Add patient with reminder without attachments
		dueTime := time.Now().UTC().Add(-1 * time.Minute)
		store.Lock()
		store.Patients["patient-2"] = &models.Patient{
			ID:    "patient-2",
			Name:  "Jane Doe",
			Phone: "+6281234567891",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-2",
					Title:          "Simple Reminder",
					Description:    "No attachments here",
					DueDate:        dueTime.Format(time.RFC3339),
					DeliveryStatus: models.DeliveryStatusPending,
					Attachments:    nil,
				},
			},
		}
		store.Unlock()

		scheduler.processScheduledReminders()
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		msg := sentMessage
		mu.Unlock()

		if msg == "" {
			t.Fatal("Expected message to be sent")
		}

		// Check message contains basic content
		if !strings.Contains(msg, "Halo Jane Doe") {
			t.Error("Message should contain patient name")
		}
		if !strings.Contains(msg, "*Simple Reminder*") {
			t.Error("Message should contain reminder title")
		}

		// Check message does NOT contain "Konten Edukasi" section
		if strings.Contains(msg, "Konten Edukasi") {
			t.Error("Message without attachments should not contain 'Konten Edukasi' section")
		}
	})

	t.Run("handles missing content gracefully", func(t *testing.T) {
		var sentMessage string
		var mu sync.Mutex

		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var req struct {
				Phone   string `json:"phone"`
				Message string `json:"message"`
			}
			json.NewDecoder(r.Body).Decode(&req)

			mu.Lock()
			sentMessage = req.Message
			mu.Unlock()

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SendMessageResponse{
				Success:   true,
				MessageID: "msg-missing-content-123",
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

		// Create empty content stores (content was deleted)
		articleStore := &models.ArticleStore{
			Articles: map[string]*models.Article{},
		}
		videoStore := &models.VideoStore{
			Videos: map[string]*models.Video{},
		}

		scheduler := NewReminderScheduler(store, gowaClient, cfg, logger)
		scheduler.SetContentStores(articleStore, videoStore)

		// Add patient with reminder referencing deleted content
		dueTime := time.Now().UTC().Add(-1 * time.Minute)
		store.Lock()
		store.Patients["patient-3"] = &models.Patient{
			ID:    "patient-3",
			Name:  "Bob Smith",
			Phone: "+6281234567892",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-3",
					Title:          "Reminder With Deleted Content",
					Description:    "Content was deleted",
					DueDate:        dueTime.Format(time.RFC3339),
					DeliveryStatus: models.DeliveryStatusPending,
					Attachments: []models.Attachment{
						{ID: "deleted-article", Type: "article", Title: "Deleted Article"},
					},
				},
			},
		}
		store.Unlock()

		scheduler.processScheduledReminders()
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		msg := sentMessage
		mu.Unlock()

		if msg == "" {
			t.Fatal("Expected message to be sent even with missing content")
		}

		// Check message contains fallback text for missing content
		if !strings.Contains(msg, "Konten tidak tersedia") {
			t.Error("Message should contain fallback text for missing content")
		}
	})
}

// TestBuildContentAttachments tests the shared utility function
func TestBuildContentAttachments(t *testing.T) {
	t.Run("builds attachments with article and video", func(t *testing.T) {
		articleStore := &models.ArticleStore{
			Articles: map[string]*models.Article{
				"art-1": {
					ID:      "art-1",
					Title:   "Article One",
					Slug:    "article-one",
					Excerpt: "Article one excerpt",
				},
			},
		}

		videoStore := &models.VideoStore{
			Videos: map[string]*models.Video{
				"vid-1": {
					ID:        "vid-1",
					Title:     "Video One",
					YouTubeID: "yt123",
				},
			},
		}

		attachments := []models.Attachment{
			{ID: "vid-1", Type: "video", Title: "Video One"},
			{ID: "art-1", Type: "article", Title: "Article One"},
		}

		result := utils.BuildContentAttachments(attachments, articleStore, videoStore)

		if len(result) != 2 {
			t.Fatalf("Expected 2 attachments, got %d", len(result))
		}

		// Articles should be sorted first
		if result[0].Type != "article" {
			t.Error("Articles should be sorted before videos")
		}
		if result[0].Excerpt != "Article one excerpt" {
			t.Error("Article should have excerpt from store")
		}
		if result[0].URL != "https://prima.app/artikel/article-one" {
			t.Errorf("Article URL incorrect: %s", result[0].URL)
		}

		if result[1].Type != "video" {
			t.Error("Video should be second")
		}
		if result[1].URL != "https://youtube.com/watch?v=yt123" {
			t.Errorf("Video URL incorrect: %s", result[1].URL)
		}
	})

	t.Run("handles nil stores gracefully", func(t *testing.T) {
		attachments := []models.Attachment{
			{ID: "art-1", Type: "article", Title: "Article", URL: "https://fallback.url"},
		}

		result := utils.BuildContentAttachments(attachments, nil, nil)

		if len(result) != 1 {
			t.Fatalf("Expected 1 attachment, got %d", len(result))
		}

		// Should use fallback URL when stores are nil
		if result[0].URL != "https://fallback.url" {
			t.Errorf("Should use fallback URL, got: %s", result[0].URL)
		}
	})
}
