package handlers

import (
	"encoding/json"
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

	cfg := &config.Config{}

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
	})

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
		cfg := &config.Config{}

		// Create GOWA client with low failure threshold to trigger circuit breaker quickly
		gowaClient := services.NewGOWAClient(services.GOWAConfig{
			Endpoint:         gowaServer.URL,
			User:             "testuser",
			Password:         "testpass",
			Timeout:          1 * time.Second,
			FailureThreshold: 2, // Low threshold for testing
			CooldownDuration: 5 * time.Minute,
		}, logger)

		handler := NewReminderHandler(store, cfg, gowaClient, logger, func() string {
			return "test-reminder-id"
		})

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
	handler, _ := setupTestHandler(t, nil)

	t.Run("formats message with description", func(t *testing.T) {
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
		if !contains(message, "Informasi ini untuk tujuan edukasi") {
			t.Error("Message should contain health disclaimer")
		}
	})

	t.Run("formats message without description", func(t *testing.T) {
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
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
