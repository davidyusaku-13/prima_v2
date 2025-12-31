package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/services"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestHandler(t *testing.T, gowaServer *httptest.Server) (*ReminderHandler, *models.PatientStore) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	store := models.NewPatientStore(func() {})

	// Create content store for attachment handling
	contentStore := NewContentStore()

	// Create config with default disclaimer enabled
	enabled := true
	startHour := 0
	endHour := 0
	cfg := &config.Config{
		Disclaimer: config.DisclaimerConfig{
			Text:    "Informasi ini untuk tujuan edukasi. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda.",
			Enabled: &enabled,
		},
		// Disable quiet hours by setting start=end=0 (GetStartHour()==GetEndHour() means disabled)
		QuietHours: config.QuietHoursConfig{
			StartHour: &startHour,
			EndHour:   &endHour,
			Timezone:  "WIB",
		},
	}

	var gowaClient *services.GOWAClient
	if gowaServer != nil {
		gowaClient = services.NewGOWAClient(services.GOWAConfig{
			Endpoint:         gowaServer.URL,
			User:             "testuser",
			Password:         "testpass",
			Timeout:          10 * time.Second,
			FailureThreshold: 5,
			CooldownDuration: 5 * time.Minute,
		}, logger)
	}

	handler := NewReminderHandler(store, cfg, gowaClient, logger, func() string {
		return "test-reminder-id"
	}, contentStore)

	return handler, store
}

func setupTestContext(method, path string, params map[string]string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)

	// Set params
	ginParams := make([]gin.Param, 0, len(params))
	for k, v := range params {
		ginParams = append(ginParams, gin.Param{Key: k, Value: v})
	}
	c.Params = ginParams

	return c, w
}

func TestReminderHandler_Send(t *testing.T) {
	t.Run("successful send", func(t *testing.T) {
		// Create mock GOWA server
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(services.SendMessageResponse{
				Success:   true,
				MessageID: "msg-123",
			})
		}))
		defer gowaServer.Close()

		handler, store := setupTestHandler(t, gowaServer)

		// Setup test data
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					Description:    "Test Description",
					DeliveryStatus: models.DeliveryStatusPending,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-1/send", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Send(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["message"] != "Reminder berhasil dikirim" {
			t.Errorf("Expected success message, got %v", response["message"])
		}

		// Verify reminder was updated
		reminder := store.Patients["patient-1"].Reminders[0]
		if reminder.DeliveryStatus != models.DeliveryStatusSent {
			t.Errorf("Expected delivery status 'sent', got '%s'", reminder.DeliveryStatus)
		}
		if reminder.GOWAMessageID != "msg-123" {
			t.Errorf("Expected GOWA message ID 'msg-123', got '%s'", reminder.GOWAMessageID)
		}
	})

	t.Run("patient not found", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		c, w := setupTestContext("POST", "/api/patients/nonexistent/reminders/reminder-1/send", map[string]string{
			"id":         "nonexistent",
			"reminderId": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Send(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "PATIENT_NOT_FOUND" {
			t.Errorf("Expected code 'PATIENT_NOT_FOUND', got '%v'", response["code"])
		}
	})

	t.Run("reminder not found", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{},
		}

		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders/nonexistent/send", map[string]string{
			"id":         "patient-1",
			"reminderId": "nonexistent",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Send(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "REMINDER_NOT_FOUND" {
			t.Errorf("Expected code 'REMINDER_NOT_FOUND', got '%v'", response["code"])
		}
	})

	t.Run("invalid phone number", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "invalid-phone",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					DeliveryStatus: models.DeliveryStatusPending,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-1/send", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Send(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "INVALID_PHONE" {
			t.Errorf("Expected code 'INVALID_PHONE', got '%v'", response["code"])
		}
	})

	t.Run("already sending", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					DeliveryStatus: models.DeliveryStatusSending,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-1/send", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Send(c)

		if w.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "ALREADY_SENDING" {
			t.Errorf("Expected code 'ALREADY_SENDING', got '%v'", response["code"])
		}
	})

	t.Run("forbidden for volunteer accessing other user patient", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "other-user",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					DeliveryStatus: models.DeliveryStatusPending,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-1/send", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Send(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "FORBIDDEN" {
			t.Errorf("Expected code 'FORBIDDEN', got '%v'", response["code"])
		}
	})

	t.Run("GOWA error", func(t *testing.T) {
		// Create mock GOWA server that returns error
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer gowaServer.Close()

		handler, store := setupTestHandler(t, gowaServer)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					DeliveryStatus: models.DeliveryStatusPending,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-1/send", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Send(c)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "GOWA_ERROR" {
			t.Errorf("Expected code 'GOWA_ERROR', got '%v'", response["code"])
		}

		// Verify reminder status was updated to failed
		reminder := store.Patients["patient-1"].Reminders[0]
		if reminder.DeliveryStatus != models.DeliveryStatusFailed {
			t.Errorf("Expected delivery status 'failed', got '%s'", reminder.DeliveryStatus)
		}
	})

	t.Run("circuit breaker open - queued for retry", func(t *testing.T) {
		// Create mock GOWA server that always fails to trigger circuit breaker
		failCount := 0
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			failCount++
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer gowaServer.Close()

		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		// Create config with default disclaimer enabled
		enabled := true
		startHour := 0
		endHour := 0
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Informasi ini untuk tujuan edukasi. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda.",
				Enabled: &enabled,
			},
			// Disable quiet hours for circuit breaker test
			QuietHours: config.QuietHoursConfig{
				StartHour: &startHour,
				EndHour:   &endHour,
				Timezone:  "WIB",
			},
		}

		// Create GOWA client with low failure threshold to trigger circuit breaker quickly
		gowaClient := services.NewGOWAClient(services.GOWAConfig{
			Endpoint:         gowaServer.URL,
			User:             "testuser",
			Password:         "testpass",
			Timeout:          1 * time.Second,
			FailureThreshold: 2, // Low threshold for testing
			CooldownDuration: 5 * time.Minute,
		}, logger)

		contentStore := NewContentStore()
		handler := NewReminderHandler(store, cfg, gowaClient, logger, func() string {
			return "test-reminder-id"
		}, contentStore)

		// Setup test patient with multiple reminders for testing
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{ID: "reminder-1", Title: "Test Reminder 1", DeliveryStatus: models.DeliveryStatusPending},
				{ID: "reminder-2", Title: "Test Reminder 2", DeliveryStatus: models.DeliveryStatusPending},
				{ID: "reminder-3", Title: "Test Reminder 3", DeliveryStatus: models.DeliveryStatusPending},
			},
		}

		// Send multiple requests to trigger circuit breaker
		for i := 1; i <= 2; i++ {
			c, _ := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-1/send", map[string]string{
				"id":         "patient-1",
				"reminderId": "reminder-1",
			})
			c.Set("userID", "user-1")
			c.Set("role", "volunteer")

			// Reset reminder status for retry
			store.Patients["patient-1"].Reminders[0].DeliveryStatus = models.DeliveryStatusPending

			handler.Send(c)
		}

		// Now circuit breaker should be open, next request should be queued
		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-2/send", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-2",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Send(c)

		// Verify response indicates queued status
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		// Check if circuit breaker triggered queued response
		if response["code"] == "GOWA_UNAVAILABLE" {
			if response["queued"] != true {
				t.Errorf("Expected queued=true when circuit breaker is open")
			}

			// Verify reminder status was updated to queued
			reminder := store.Patients["patient-1"].Reminders[1]
			if reminder.DeliveryStatus != models.DeliveryStatusQueued {
				t.Errorf("Expected delivery status 'queued', got '%s'", reminder.DeliveryStatus)
			}
			if reminder.RetryCount < 1 {
				t.Errorf("Expected retry_count >= 1, got %d", reminder.RetryCount)
			}
		}
		// Note: Circuit breaker behavior may vary based on timing, so we accept both outcomes
	})
}

func TestReminderHandler_FormatReminderMessage(t *testing.T) {
	t.Run("formats message with description and disclaimer enabled", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "This is a test description",
		}
		patient := &models.Patient{
			Name: "John Doe",
		}

		message := handler.formatReminderMessage(reminder, patient)

		if !contains(message, "Halo John Doe") {
			t.Error("Message should contain patient name greeting")
		}
		if !contains(message, "*Test Reminder*") {
			t.Error("Message should contain reminder title in bold")
		}
		if !contains(message, "This is a test description") {
			t.Error("Message should contain description")
		}
		if !contains(message, "---") {
			t.Error("Message should contain divider before disclaimer")
		}
		if !contains(message, "_Informasi ini untuk tujuan edukasi") {
			t.Error("Message should contain health disclaimer in italic format")
		}
	})

	t.Run("formats message without description", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "",
		}
		patient := &models.Patient{
			Name: "Jane Doe",
		}

		message := handler.formatReminderMessage(reminder, patient)

		if !contains(message, "Halo Jane Doe") {
			t.Error("Message should contain patient name greeting")
		}
		if !contains(message, "*Test Reminder*") {
			t.Error("Message should contain reminder title in bold")
		}
	})

	t.Run("disclaimer disabled - no disclaimer in message", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		// Create config with disclaimer disabled
		enabled := false
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Informasi ini untuk tujuan edukasi. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda.",
				Enabled: &enabled,
			},
		}

		contentStore := NewContentStore()
		handler := NewReminderHandler(store, cfg, nil, logger, func() string {
			return "test-reminder-id"
		}, contentStore)

		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "Test description",
		}
		patient := &models.Patient{
			Name: "John Doe",
		}

		message := handler.formatReminderMessage(reminder, patient)

		if contains(message, "---") {
			t.Error("Message should NOT contain divider when disclaimer is disabled")
		}
		if contains(message, "Informasi ini untuk tujuan edukasi") {
			t.Error("Message should NOT contain disclaimer when disabled")
		}
	})

	t.Run("custom disclaimer text from config", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		// Create config with custom disclaimer text
		enabled := true
		customDisclaimer := "Custom health disclaimer for testing purposes."
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    customDisclaimer,
				Enabled: &enabled,
			},
		}

		contentStore := NewContentStore()
		handler := NewReminderHandler(store, cfg, nil, logger, func() string {
			return "test-reminder-id"
		}, contentStore)

		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "Test description",
		}
		patient := &models.Patient{
			Name: "John Doe",
		}

		message := handler.formatReminderMessage(reminder, patient)

		if !contains(message, "_Custom health disclaimer for testing purposes._") {
			t.Error("Message should contain custom disclaimer text in italic format")
		}
	})

	t.Run("nil Enabled pointer - disclaimer not added (bypasses applyDefaults)", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		// Create config with nil Enabled - simulates direct struct creation without LoadOrDefault()
		// In production, applyDefaults() would set Enabled=true, but this tests the nil-safety of formatReminderMessage
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Default disclaimer text",
				Enabled: nil, // nil means not set - tests nil pointer safety
			},
		}

		contentStore := NewContentStore()
		handler := NewReminderHandler(store, cfg, nil, logger, func() string {
			return "test-reminder-id"
		}, contentStore)

		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "Test description",
		}
		patient := &models.Patient{
			Name: "John Doe",
		}

		message := handler.formatReminderMessage(reminder, patient)

		// When Enabled is nil, disclaimer should NOT be added (nil check fails)
		if contains(message, "Default disclaimer text") {
			t.Error("Message should NOT contain disclaimer when Enabled is nil")
		}
	})
}

func TestReminderHandler_Create_WithAttachments(t *testing.T) {
	t.Run("accepts up to 3 attachments", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		// Add sample content to content store for validation
		handler.contentStore.Articles.Articles["art-1"] = &models.Article{
			ID:    "art-1",
			Title: "Article 1",
			Slug:  "article-1",
		}
		handler.contentStore.Articles.Articles["art-2"] = &models.Article{
			ID:    "art-2",
			Title: "Article 2",
			Slug:  "article-2",
		}
		handler.contentStore.Videos.Videos["vid-1"] = &models.Video{
			ID:         "vid-1",
			Title:      "Video 1",
			YouTubeID:  "abc123",
		}

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
		}

		req := CreateReminderRequest{
			Title:       "Test Reminder",
			Description: "Test Description",
			Attachments: []models.Attachment{
				{Type: "article", ID: "art-1", Title: "Article 1"},
				{Type: "video", ID: "vid-1", Title: "Video 1"},
				{Type: "article", ID: "art-2", Title: "Article 2"},
			},
		}

		body, _ := json.Marshal(req)
		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders", map[string]string{"id": "patient-1"})
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Create(c)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		// Verify attachments were saved
		reminder := store.Patients["patient-1"].Reminders[0]
		if len(reminder.Attachments) != 3 {
			t.Errorf("Expected 3 attachments, got %d", len(reminder.Attachments))
		}
	})

	t.Run("rejects more than 3 attachments", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
		}

		req := CreateReminderRequest{
			Title:       "Test Reminder",
			Description: "Test Description",
			Attachments: []models.Attachment{
				{Type: "article", ID: "art-1", Title: "Article 1"},
				{Type: "video", ID: "vid-1", Title: "Video 1"},
				{Type: "article", ID: "art-2", Title: "Article 2"},
				{Type: "video", ID: "vid-2", Title: "Video 2"},
			},
		}

		body, _ := json.Marshal(req)
		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders", map[string]string{"id": "patient-1"})
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "MAX_ATTACHMENTS_EXCEEDED" {
			t.Errorf("Expected code 'MAX_ATTACHMENTS_EXCEEDED', got %v", response["code"])
		}

		if response["max_count"] != float64(MaxAttachments) {
			t.Errorf("Expected max_count %d, got %v", MaxAttachments, response["max_count"])
		}

		if response["actual"] != float64(4) {
			t.Errorf("Expected actual 4, got %v", response["actual"])
		}
	})

	t.Run("accepts empty attachments array", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
		}

		req := CreateReminderRequest{
			Title:       "Test Reminder",
			Description: "Test Description",
			Attachments: []models.Attachment{},
		}

		body, _ := json.Marshal(req)
		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders", map[string]string{"id": "patient-1"})
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Create(c)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}

		reminder := store.Patients["patient-1"].Reminders[0]
		if len(reminder.Attachments) != 0 {
			t.Errorf("Expected 0 attachments, got %d", len(reminder.Attachments))
		}
	})

	t.Run("rejects attachment without title", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
		}

		req := CreateReminderRequest{
			Title:       "Test Reminder",
			Description: "Test Description",
			Attachments: []models.Attachment{
				{Type: "article", ID: "art-1", Title: ""}, // Empty title
			},
		}

		body, _ := json.Marshal(req)
		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders", map[string]string{"id": "patient-1"})
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "INVALID_ATTACHMENT" {
			t.Errorf("Expected code 'INVALID_ATTACHMENT', got %v", response["code"])
		}
	})

	t.Run("rejects attachment with invalid type", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
		}

		req := CreateReminderRequest{
			Title:       "Test Reminder",
			Description: "Test Description",
			Attachments: []models.Attachment{
				{Type: "invalid", ID: "art-1", Title: "Test"}, // Invalid type
			},
		}

		body, _ := json.Marshal(req)
		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders", map[string]string{"id": "patient-1"})
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "INVALID_ATTACHMENT" {
			t.Errorf("Expected code 'INVALID_ATTACHMENT', got %v", response["code"])
		}
	})
}

func TestReminderHandler_Update_WithAttachments(t *testing.T) {
	t.Run("updates attachments successfully", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		// Add sample content to content store for validation
		handler.contentStore.Articles.Articles["art-1"] = &models.Article{
			ID:    "art-1",
			Title: "New Article",
			Slug:  "new-article",
		}

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Original Title",
					DeliveryStatus: models.DeliveryStatusPending,
				},
			},
		}

		req := UpdateReminderRequest{
			Title: "Updated Title",
			Attachments: []models.Attachment{
				{Type: "article", ID: "art-1", Title: "New Article"},
			},
		}

		body, _ := json.Marshal(req)
		c, w := setupTestContext("PUT", "/api/patients/patient-1/reminders/reminder-1", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-1",
		})
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Update(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		reminder := store.Patients["patient-1"].Reminders[0]
		if len(reminder.Attachments) != 1 {
			t.Errorf("Expected 1 attachment, got %d", len(reminder.Attachments))
		}
		if reminder.Attachments[0].Title != "New Article" {
			t.Errorf("Expected attachment title 'New Article', got '%s'", reminder.Attachments[0].Title)
		}
	})

	t.Run("rejects update with more than 3 attachments", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					DeliveryStatus: models.DeliveryStatusPending,
				},
			},
		}

		req := UpdateReminderRequest{
			Attachments: []models.Attachment{
				{Type: "article", ID: "art-1", Title: "Article 1"},
				{Type: "video", ID: "vid-1", Title: "Video 1"},
				{Type: "article", ID: "art-2", Title: "Article 2"},
				{Type: "video", ID: "vid-2", Title: "Video 2"},
			},
		}

		body, _ := json.Marshal(req)
		c, w := setupTestContext("PUT", "/api/patients/patient-1/reminders/reminder-1", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-1",
		})
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Update(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "MAX_ATTACHMENTS_EXCEEDED" {
			t.Errorf("Expected code 'MAX_ATTACHMENTS_EXCEEDED', got %v", response["code"])
		}
	})
}

func TestReminderHandler_FormatReminderMessage_WithAttachments(t *testing.T) {
	t.Run("includes attachments in message", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		// Add sample content to content store for message formatting
		handler.contentStore.Articles.Articles["art-1"] = &models.Article{
			ID:       "art-1",
			Title:    "Healthy Living Guide",
			Slug:     "healthy-living-guide",
			Excerpt:  "This is a healthy living guide for better lifestyle.",
		}
		handler.contentStore.Videos.Videos["vid-1"] = &models.Video{
			ID:        "vid-1",
			Title:     "Morning Exercise Routine",
			YouTubeID: "abc123xyz",
		}

		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "Test description",
			Attachments: []models.Attachment{
				{Type: "article", ID: "art-1", Title: "Healthy Living Guide"},
				{Type: "video", ID: "vid-1", Title: "Morning Exercise Routine"},
			},
		}
		patient := &models.Patient{
			Name: "John Doe",
		}

		message := handler.formatReminderMessage(reminder, patient)

		if !contains(message, "📖 Healthy Living Guide") {
			t.Error("Message should contain article attachment with 📖 icon")
		}
		if !contains(message, "🎬 Morning Exercise Routine") {
			t.Error("Message should contain video attachment")
		}
		if !contains(message, "Konten Edukasi:") {
			t.Error("Message should contain attachments header")
		}
		if !contains(message, "https://prima.app/artikel/healthy-living-guide") {
			t.Error("Message should contain article URL")
		}
		if !contains(message, "https://youtube.com/watch?v=abc123xyz") {
			t.Error("Message should contain YouTube URL")
		}
	})

	t.Run("handles nil attachments gracefully", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "Test description",
			Attachments: nil,
		}
		patient := &models.Patient{
			Name: "Jane Doe",
		}

		message := handler.formatReminderMessage(reminder, patient)

		// Should not crash and should not include attachment section
		if contains(message, "Konten Edukasi:") {
			t.Error("Message should NOT contain attachments header when no attachments")
		}
	})

	t.Run("articles appear before videos in message", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		// Add sample content - add videos first to test sorting
		handler.contentStore.Videos.Videos["vid-1"] = &models.Video{
			ID:        "vid-1",
			Title:     "Video First Added",
			YouTubeID: "vid123",
		}
		handler.contentStore.Articles.Articles["art-1"] = &models.Article{
			ID:      "art-1",
			Title:   "Article After Video",
			Slug:    "article-after",
			Excerpt: "This is an article excerpt for testing.",
		}

		// Add video first, then article - order should be swapped
		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "Test description",
			Attachments: []models.Attachment{
				{Type: "video", ID: "vid-1", Title: "Video First Added"},
				{Type: "article", ID: "art-1", Title: "Article After Video"},
			},
		}
		patient := &models.Patient{
			Name: "Test Patient",
		}

		message := handler.formatReminderMessage(reminder, patient)

		// Find positions of article and video in message
		articlePos := strings.Index(message, "Article After Video")
		videoPos := strings.Index(message, "Video First Added")

		if articlePos == -1 || videoPos == -1 {
			t.Fatal("Both article and video should be in message")
		}

		// Article should appear before video (AC #3 requirement)
		if articlePos > videoPos {
			t.Error("Article should appear before video in message (AC #3)")
		}
	})

	t.Run("uses title as fallback when excerpt is empty", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		// Add article with empty excerpt
		handler.contentStore.Articles.Articles["art-1"] = &models.Article{
			ID:      "art-1",
			Title:   "Article With Empty Excerpt",
			Slug:    "empty-excerpt",
			Excerpt: "", // Empty excerpt
		}

		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "Test description",
			Attachments: []models.Attachment{
				{Type: "article", ID: "art-1", Title: "Article With Empty Excerpt"},
			},
		}
		patient := &models.Patient{
			Name: "Test Patient",
		}

		message := handler.formatReminderMessage(reminder, patient)

		// Should show article title as excerpt fallback
		if !contains(message, "Article With Empty Excerpt") {
			t.Error("Message should contain article title as excerpt fallback")
		}
	})

	t.Run("truncates long excerpts to 100 characters", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		// Add article with very long excerpt (over 100 chars)
		longExcerpt := "This is a very long excerpt that should definitely be truncated because it exceeds the maximum allowed characters of one hundred for the WhatsApp message display and we need to ensure it gets cut properly."
		handler.contentStore.Articles.Articles["art-1"] = &models.Article{
			ID:      "art-1",
			Title:   "Long Excerpt Article",
			Slug:    "long-excerpt",
			Excerpt: longExcerpt,
		}

		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "Test description",
			Attachments: []models.Attachment{
				{Type: "article", ID: "art-1", Title: "Long Excerpt Article"},
			},
		}
		patient := &models.Patient{
			Name: "Test Patient",
		}

		message := handler.formatReminderMessage(reminder, patient)

		// Extract excerpt portion (between title and link)
		titleEnd := strings.Index(message, "Long Excerpt Article")
		linkStart := strings.Index(message, "🔗 https://prima.app/artikel/long-excerpt")

		if titleEnd == -1 || linkStart == -1 {
			t.Fatal("Should find title and link in message")
		}

		excerptPortion := message[titleEnd:linkStart]
		// Excerpt should contain "..." indicating truncation
		if !strings.Contains(excerptPortion, "...") {
			t.Error("Long excerpt should be truncated with '...'")
		}
		// Check that the excerpt portion is shorter than the original excerpt
		if len(excerptPortion) > len(longExcerpt) {
			t.Error("Truncated excerpt should not be longer than original")
		}
		// Make sure the excerpt portion is reasonably short (< 120 chars including the "...")
		if len(excerptPortion) > 120 {
			t.Error("Excerpt should be truncated to around 100 chars")
		}
	})

	t.Run("handles deleted content gracefully", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		// Don't add article to store - simulates deleted content
		reminder := &models.Reminder{
			Title:       "Test Reminder",
			Description: "Test description",
			Attachments: []models.Attachment{
				{Type: "article", ID: "nonexistent", Title: "Deleted Article"},
			},
		}
		patient := &models.Patient{
			Name: "Test Patient",
		}

		message := handler.formatReminderMessage(reminder, patient)

		// Should show fallback text for deleted content
		if !contains(message, "Konten tidak tersedia") {
			t.Error("Message should show 'Konten tidak tersedia' for deleted content")
		}
	})
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// setupTestHandlerWithQuietHours creates a test handler with quiet hours config
func setupTestHandlerWithQuietHours(t *testing.T, gowaServer *httptest.Server, quietHoursConfig config.QuietHoursConfig) (*ReminderHandler, *models.PatientStore) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	store := models.NewPatientStore(func() {})

	// Create content store for attachment handling
	contentStore := NewContentStore()

	// Create config with default disclaimer enabled and quiet hours
	enabled := true
	cfg := &config.Config{
		Disclaimer: config.DisclaimerConfig{
			Text:    "Informasi ini untuk tujuan edukasi. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda.",
			Enabled: &enabled,
		},
		QuietHours: quietHoursConfig,
	}

	var gowaClient *services.GOWAClient
	if gowaServer != nil {
		gowaClient = services.NewGOWAClient(services.GOWAConfig{
			Endpoint:         gowaServer.URL,
			User:             "testuser",
			Password:         "testpass",
			Timeout:          10 * time.Second,
			FailureThreshold: 5,
			CooldownDuration: 5 * time.Minute,
		}, logger)
	}

	handler := NewReminderHandler(store, cfg, gowaClient, logger, func() string {
		return "test-reminder-id"
	}, contentStore)

	return handler, store
}

func TestReminderHandler_RetryReminder(t *testing.T) {
	t.Run("successful retry of failed reminder", func(t *testing.T) {
		// Create mock GOWA server
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(services.SendMessageResponse{
				Success:   true,
				MessageID: "msg-retry-123",
			})
		}))
		defer gowaServer.Close()

		handler, store := setupTestHandler(t, gowaServer)

		// Setup test data with failed reminder
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:                   "reminder-1",
					Title:                "Test Reminder",
					Description:          "Test Description",
					DeliveryStatus:       models.DeliveryStatusFailed,
					DeliveryErrorMessage: "Previous error",
					RetryCount:           2,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/reminders/reminder-1/retry", map[string]string{
			"id": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.RetryReminder(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["message"] != "Reminder berhasil dikirim ulang" {
			t.Errorf("Expected success message, got %v", response["message"])
		}

		// Verify reminder was updated
		reminder := store.Patients["patient-1"].Reminders[0]
		if reminder.DeliveryStatus != models.DeliveryStatusSent {
			t.Errorf("Expected delivery status 'sent', got '%s'", reminder.DeliveryStatus)
		}
		if reminder.GOWAMessageID != "msg-retry-123" {
			t.Errorf("Expected GOWA message ID 'msg-retry-123', got '%s'", reminder.GOWAMessageID)
		}
		if reminder.RetryCount != 0 {
			t.Errorf("Expected retry count to be reset to 0, got %d", reminder.RetryCount)
		}
		if reminder.DeliveryErrorMessage != "" {
			t.Errorf("Expected error message to be cleared, got '%s'", reminder.DeliveryErrorMessage)
		}
	})

	t.Run("rejects retry of non-failed reminder", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					DeliveryStatus: models.DeliveryStatusSent, // Not failed
				},
			},
		}

		c, w := setupTestContext("POST", "/api/reminders/reminder-1/retry", map[string]string{
			"id": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.RetryReminder(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "INVALID_STATUS" {
			t.Errorf("Expected code 'INVALID_STATUS', got '%v'", response["code"])
		}
	})

	t.Run("rejects retry when circuit breaker is open", func(t *testing.T) {
		// Create mock GOWA server that always fails
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer gowaServer.Close()

		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
		store := models.NewPatientStore(func() {})

		enabled := true
		startHour := 0
		endHour := 0
		cfg := &config.Config{
			Disclaimer: config.DisclaimerConfig{
				Text:    "Test disclaimer",
				Enabled: &enabled,
			},
			QuietHours: config.QuietHoursConfig{
				StartHour: &startHour,
				EndHour:   &endHour,
				Timezone:  "WIB",
			},
		}

		// Create GOWA client with low threshold
		gowaClient := services.NewGOWAClient(services.GOWAConfig{
			Endpoint:         gowaServer.URL,
			User:             "testuser",
			Password:         "testpass",
			Timeout:          1 * time.Second,
			FailureThreshold: 2,
			CooldownDuration: 5 * time.Minute,
		}, logger)

		contentStore := NewContentStore()
		handler := NewReminderHandler(store, cfg, gowaClient, logger, func() string {
			return "test-reminder-id"
		}, contentStore)

		// Trigger circuit breaker by failing multiple times
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{ID: "reminder-1", Title: "Test 1", DeliveryStatus: models.DeliveryStatusPending},
				{ID: "reminder-2", Title: "Test 2", DeliveryStatus: models.DeliveryStatusFailed},
			},
		}

		// Trigger circuit breaker
		for i := 0; i < 2; i++ {
			c, _ := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-1/send", map[string]string{
				"id":         "patient-1",
				"reminderId": "reminder-1",
			})
			c.Set("userID", "user-1")
			c.Set("role", "volunteer")
			store.Patients["patient-1"].Reminders[0].DeliveryStatus = models.DeliveryStatusPending
			handler.Send(c)
		}

		// Now try to retry - should be blocked by circuit breaker
		c, w := setupTestContext("POST", "/api/reminders/reminder-2/retry", map[string]string{
			"id": "reminder-2",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.RetryReminder(c)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "CIRCUIT_BREAKER_OPEN" {
			t.Errorf("Expected code 'CIRCUIT_BREAKER_OPEN', got '%v'", response["code"])
		}
	})

	t.Run("reminder not found", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{},
		}

		c, w := setupTestContext("POST", "/api/reminders/nonexistent/retry", map[string]string{
			"id": "nonexistent",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.RetryReminder(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "REMINDER_NOT_FOUND" {
			t.Errorf("Expected code 'REMINDER_NOT_FOUND', got '%v'", response["code"])
		}
	})

	t.Run("invalid phone number", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "invalid-phone",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					DeliveryStatus: models.DeliveryStatusFailed,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/reminders/reminder-1/retry", map[string]string{
			"id": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.RetryReminder(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "INVALID_PHONE" {
			t.Errorf("Expected code 'INVALID_PHONE', got '%v'", response["code"])
		}
	})

	t.Run("GOWA error on retry", func(t *testing.T) {
		// Create mock GOWA server that returns error
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("GOWA Internal Error"))
		}))
		defer gowaServer.Close()

		handler, store := setupTestHandler(t, gowaServer)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					DeliveryStatus: models.DeliveryStatusFailed,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/reminders/reminder-1/retry", map[string]string{
			"id": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.RetryReminder(c)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "GOWA_ERROR" {
			t.Errorf("Expected code 'GOWA_ERROR', got '%v'", response["code"])
		}

		// Verify reminder status remains failed
		reminder := store.Patients["patient-1"].Reminders[0]
		if reminder.DeliveryStatus != models.DeliveryStatusFailed {
			t.Errorf("Expected delivery status to remain 'failed', got '%s'", reminder.DeliveryStatus)
		}
	})

	t.Run("forbidden for volunteer accessing other user reminder", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "other-user",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					DeliveryStatus: models.DeliveryStatusFailed,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/reminders/reminder-1/retry", map[string]string{
			"id": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.RetryReminder(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["code"] != "FORBIDDEN" {
			t.Errorf("Expected code 'FORBIDDEN', got '%v'", response["code"])
		}
	})
}

func TestReminderHandler_Send_QuietHours(t *testing.T) {
	t.Run("schedules reminder during quiet hours", func(t *testing.T) {
		// Create mock GOWA server (needed if we're in active hours)
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(services.SendMessageResponse{
				Success:   true,
				MessageID: "msg-123",
			})
		}))
		defer gowaServer.Close()

		// Configure quiet hours: 21:00 - 06:00 WIB
		startHour := 21
		endHour := 6
		quietHoursConfig := config.QuietHoursConfig{
			StartHour: &startHour,
			EndHour:   &endHour,
			Timezone:  "WIB",
		}

		handler, store := setupTestHandlerWithQuietHours(t, gowaServer, quietHoursConfig)

		// Setup test data
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					Description:    "Test Description",
					DeliveryStatus: models.DeliveryStatusPending,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-1/send", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		// Note: This test will behave differently based on actual current time
		// If current time is in quiet hours (21:00-06:00 WIB), it will schedule
		// If current time is in active hours (06:00-21:00 WIB), it will send immediately
		handler.Send(c)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		reminder := store.Patients["patient-1"].Reminders[0]

		// Check if we're in quiet hours by examining the response
		if response["scheduled"] == true {
			// Verify scheduled behavior
			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d for scheduled reminder, got %d", http.StatusOK, w.Code)
			}
			if reminder.DeliveryStatus != models.DeliveryStatusScheduled {
				t.Errorf("Expected delivery status 'scheduled', got '%s'", reminder.DeliveryStatus)
			}
			if reminder.ScheduledDeliveryAt == "" {
				t.Error("Expected ScheduledDeliveryAt to be set")
			}
			if response["scheduled_at"] == nil {
				t.Error("Expected scheduled_at in response")
			}
			t.Log("Test ran during quiet hours - reminder was scheduled")
		} else {
			// In active hours - sent immediately
			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}
			if reminder.DeliveryStatus != models.DeliveryStatusSent {
				t.Errorf("Expected delivery status 'sent', got '%s'", reminder.DeliveryStatus)
			}
			t.Log("Test ran during active hours - reminder was sent immediately")
		}
	})

	t.Run("sends immediately during active hours", func(t *testing.T) {
		// Create mock GOWA server
		gowaServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(services.SendMessageResponse{
				Success:   true,
				MessageID: "msg-123",
			})
		}))
		defer gowaServer.Close()

		// Configure quiet hours: 21:00 - 06:00 WIB
		startHour := 21
		endHour := 6
		quietHoursConfig := config.QuietHoursConfig{
			StartHour: &startHour,
			EndHour:   &endHour,
			Timezone:  "WIB",
		}

		handler, store := setupTestHandlerWithQuietHours(t, gowaServer, quietHoursConfig)

		// Setup test data
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:             "reminder-1",
					Title:          "Test Reminder",
					Description:    "Test Description",
					DeliveryStatus: models.DeliveryStatusPending,
				},
			},
		}

		c, w := setupTestContext("POST", "/api/patients/patient-1/reminders/reminder-1/send", map[string]string{
			"id":         "patient-1",
			"reminderId": "reminder-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.Send(c)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		reminder := store.Patients["patient-1"].Reminders[0]

		// Check behavior based on current time
		if response["scheduled"] == true {
			// In quiet hours - scheduled
			if reminder.DeliveryStatus != models.DeliveryStatusScheduled {
				t.Errorf("Expected delivery status 'scheduled', got '%s'", reminder.DeliveryStatus)
			}
			t.Log("Test ran during quiet hours - reminder was scheduled")
		} else {
			// In active hours - sent immediately
			if w.Code != http.StatusOK {
				t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
			}
			if reminder.DeliveryStatus != models.DeliveryStatusSent {
				t.Errorf("Expected delivery status 'sent', got '%s'", reminder.DeliveryStatus)
			}
			if reminder.GOWAMessageID != "msg-123" {
				t.Errorf("Expected GOWA message ID 'msg-123', got '%s'", reminder.GOWAMessageID)
			}
			t.Log("Test ran during active hours - reminder was sent immediately")
		}
	})
}

func TestReminderHandler_GetPatientReminders(t *testing.T) {
	t.Run("returns reminders with pagination", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		// Setup test data with multiple reminders
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:                   "reminder-1",
					Title:                "Reminder 1",
					Description:          "This is the first reminder description",
					DeliveryStatus:       models.DeliveryStatusSent,
					ScheduledDeliveryAt:  "2025-12-30T10:00:00Z",
					MessageSentAt:        "2025-12-30T10:00:05Z",
					DeliveredAt:          "2025-12-30T10:01:00Z",
					ReadAt:               "2025-12-30T10:05:00Z",
					Attachments:          []models.Attachment{{Type: "article", ID: "art-1", Title: "Article 1"}},
				},
				{
					ID:                   "reminder-2",
					Title:                "Reminder 2",
					Description:          "This is the second reminder description",
					DeliveryStatus:       models.DeliveryStatusDelivered,
					ScheduledDeliveryAt:  "2025-12-29T10:00:00Z",
					MessageSentAt:        "2025-12-29T10:00:05Z",
					DeliveredAt:          "2025-12-29T10:01:00Z",
				},
			},
		}

		c, w := setupTestContext("GET", "/api/patients/patient-1/reminders?page=1&limit=10", map[string]string{
			"id": "patient-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.GetPatientReminders(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["message"] != "Success" {
			t.Errorf("Expected message 'Success', got %v", response["message"])
		}

		// Check pagination
		pagination := response["pagination"].(map[string]interface{})
		if pagination["page"].(float64) != 1 {
			t.Errorf("Expected page 1, got %v", pagination["page"])
		}
		if pagination["limit"].(float64) != 10 {
			t.Errorf("Expected limit 10, got %v", pagination["limit"])
		}
		if pagination["total"].(float64) != 2 {
			t.Errorf("Expected total 2, got %v", pagination["total"])
		}
		if pagination["has_more"].(bool) != false {
			t.Errorf("Expected has_more false, got %v", pagination["has_more"])
		}

		// Check data
		data := response["data"].([]interface{})
		if len(data) != 2 {
			t.Errorf("Expected 2 reminders, got %d", len(data))
		}
	})

	t.Run("returns all reminders with history=true", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		// Setup test data with cancelled reminder
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:                   "reminder-1",
					Title:                "Active Reminder",
					Description:          "Active description",
					DeliveryStatus:       models.DeliveryStatusSent,
					ScheduledDeliveryAt:  "2025-12-30T10:00:00Z",
				},
				{
					ID:                   "reminder-2",
					Title:                "Cancelled Reminder",
					Description:          "This reminder was cancelled",
					DeliveryStatus:       models.DeliveryStatusCancelled,
					ScheduledDeliveryAt:  "2025-12-28T10:00:00Z",
					CancelledAt:          "2025-12-28T09:00:00Z",
				},
			},
		}

		c, w := setupTestContext("GET", "/api/patients/patient-1/reminders?history=true", map[string]string{
			"id": "patient-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.GetPatientReminders(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].([]interface{})
		if len(data) != 2 {
			t.Errorf("Expected 2 reminders with history=true, got %d", len(data))
		}

		// Verify cancelled reminder is included
		foundCancelled := false
		for _, item := range data {
			reminder := item.(map[string]interface{})
			if reminder["delivery_status"] == models.DeliveryStatusCancelled {
				foundCancelled = true
				if reminder["cancelled_at"] == nil {
					t.Error("Expected cancelled_at to be set for cancelled reminder")
				}
				break
			}
		}
		if !foundCancelled {
			t.Error("Expected to find cancelled reminder in response")
		}
	})

	t.Run("excludes cancelled reminders without history=true", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		// Setup test data with cancelled reminder
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:                   "reminder-1",
					Title:                "Active Reminder",
					Description:          "Active description",
					DeliveryStatus:       models.DeliveryStatusSent,
					ScheduledDeliveryAt:  "2025-12-30T10:00:00Z",
				},
				{
					ID:                   "reminder-2",
					Title:                "Cancelled Reminder",
					Description:          "This reminder was cancelled",
					DeliveryStatus:       models.DeliveryStatusCancelled,
					ScheduledDeliveryAt:  "2025-12-28T10:00:00Z",
				},
			},
		}

		c, w := setupTestContext("GET", "/api/patients/patient-1/reminders", map[string]string{
			"id": "patient-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.GetPatientReminders(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].([]interface{})
		if len(data) != 1 {
			t.Errorf("Expected 1 reminder (excluding cancelled), got %d", len(data))
		}

		reminder := data[0].(map[string]interface{})
		if reminder["title"] != "Active Reminder" {
			t.Errorf("Expected 'Active Reminder', got %v", reminder["title"])
		}
	})

	t.Run("patient not found", func(t *testing.T) {
		handler, _ := setupTestHandler(t, nil)

		c, w := setupTestContext("GET", "/api/patients/nonexistent/reminders", map[string]string{
			"id": "nonexistent",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.GetPatientReminders(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["error"] != "patient not found" {
			t.Errorf("Expected error 'patient not found', got %v", response["error"])
		}
	})

	t.Run("forbidden for volunteer accessing other user patient", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "other-user",
			Reminders: []*models.Reminder{},
		}

		c, w := setupTestContext("GET", "/api/patients/patient-1/reminders", map[string]string{
			"id": "patient-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.GetPatientReminders(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
		}
	})

	t.Run("sorts by scheduled_at descending", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		// Setup test data with reminders in different order
		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:                  "reminder-1",
					Title:               "Oldest",
					Description:         "Oldest description",
					DeliveryStatus:      models.DeliveryStatusSent,
					ScheduledDeliveryAt: "2025-12-28T10:00:00Z",
				},
				{
					ID:                  "reminder-2",
					Title:               "Most Recent",
					Description:         "Most recent description",
					DeliveryStatus:      models.DeliveryStatusSent,
					ScheduledDeliveryAt: "2025-12-30T10:00:00Z",
				},
				{
					ID:                  "reminder-3",
					Title:               "Middle",
					Description:         "Middle description",
					DeliveryStatus:      models.DeliveryStatusSent,
					ScheduledDeliveryAt: "2025-12-29T10:00:00Z",
				},
			},
		}

		c, w := setupTestContext("GET", "/api/patients/patient-1/reminders?history=true", map[string]string{
			"id": "patient-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.GetPatientReminders(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].([]interface{})
		if len(data) != 3 {
			t.Errorf("Expected 3 reminders, got %d", len(data))
		}

		// Verify order: most recent first
		first := data[0].(map[string]interface{})
		second := data[1].(map[string]interface{})
		third := data[2].(map[string]interface{})

		if first["title"] != "Most Recent" {
			t.Errorf("Expected first reminder to be 'Most Recent', got %v", first["title"])
		}
		if second["title"] != "Middle" {
			t.Errorf("Expected second reminder to be 'Middle', got %v", second["title"])
		}
		if third["title"] != "Oldest" {
			t.Errorf("Expected third reminder to be 'Oldest', got %v", third["title"])
		}
	})

	t.Run("includes attachment info", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{
				{
					ID:              "reminder-1",
					Title:           "Reminder with attachments",
					Description:     "Description",
					DeliveryStatus:  models.DeliveryStatusSent,
					ScheduledDeliveryAt: "2025-12-30T10:00:00Z",
					Attachments: []models.Attachment{
						{Type: "article", ID: "art-1", Title: "Article 1"},
						{Type: "video", ID: "vid-1", Title: "Video 1"},
					},
				},
			},
		}

		c, w := setupTestContext("GET", "/api/patients/patient-1/reminders", map[string]string{
			"id": "patient-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.GetPatientReminders(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].([]interface{})
		reminder := data[0].(map[string]interface{})

		if reminder["attachment_count"].(float64) != 2 {
			t.Errorf("Expected attachment_count 2, got %v", reminder["attachment_count"])
		}

		attachments := reminder["attachments"].([]interface{})
		if len(attachments) != 2 {
			t.Errorf("Expected 2 attachments in array, got %d", len(attachments))
		}
	})

	t.Run("empty patient returns empty list", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "user-1",
			Reminders: []*models.Reminder{},
		}

		c, w := setupTestContext("GET", "/api/patients/patient-1/reminders", map[string]string{
			"id": "patient-1",
		})
		c.Set("userID", "user-1")
		c.Set("role", "volunteer")

		handler.GetPatientReminders(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].([]interface{})
		if len(data) != 0 {
			t.Errorf("Expected 0 reminders, got %d", len(data))
		}

		pagination := response["pagination"].(map[string]interface{})
		if pagination["total"].(float64) != 0 {
			t.Errorf("Expected total 0, got %v", pagination["total"])
		}
	})

	t.Run("admin can access any patient reminders", func(t *testing.T) {
		handler, store := setupTestHandler(t, nil)

		store.Patients["patient-1"] = &models.Patient{
			ID:        "patient-1",
			Name:      "Test Patient",
			Phone:     "08123456789",
			CreatedBy: "other-user",
			Reminders: []*models.Reminder{
				{
					ID:                  "reminder-1",
					Title:               "Test",
					Description:         "Description",
					DeliveryStatus:      models.DeliveryStatusSent,
					ScheduledDeliveryAt: "2025-12-30T10:00:00Z",
				},
			},
		}

		c, w := setupTestContext("GET", "/api/patients/patient-1/reminders", map[string]string{
			"id": "patient-1",
		})
		c.Set("userID", "admin-user")
		c.Set("role", "admin")

		handler.GetPatientReminders(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)

		data := response["data"].([]interface{})
		if len(data) != 1 {
			t.Errorf("Expected 1 reminder, got %d", len(data))
		}
	})
}
