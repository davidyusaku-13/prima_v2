package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/davidyusaku-13/prima_v2/models"
)

func TestGetDeliveryAnalytics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create mock patient store with test data
	patientStore := models.NewPatientStore(func() {})

	// Add test patients with reminders
	patientStore.Patients = map[string]*models.Patient{
		"patient-1": {
			ID:   "patient-1",
			Name: "Test Patient 1",
			Reminders: []*models.Reminder{
				{
					ID:              "reminder-1",
					Title:           "Test Reminder 1",
					DeliveryStatus:  models.DeliveryStatusDelivered,
					MessageSentAt:   time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
					DeliveredAt:     time.Now().AddDate(0, 0, -1).Add(time.Minute * 2).Format(time.RFC3339),
				},
				{
					ID:              "reminder-2",
					Title:           "Test Reminder 2",
					DeliveryStatus:  models.DeliveryStatusFailed,
					MessageSentAt:   time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
				},
				{
					ID:              "reminder-3",
					Title:           "Test Reminder 3",
					DeliveryStatus:  models.DeliveryStatusPending,
				},
			},
		},
		"patient-2": {
			ID:   "patient-2",
			Name: "Test Patient 2",
			Reminders: []*models.Reminder{
				{
					ID:              "reminder-4",
					Title:           "Test Reminder 4",
					DeliveryStatus:  models.DeliveryStatusSent,
					MessageSentAt:   time.Now().AddDate(0, 0, -10).Format(time.RFC3339),
				},
				{
					ID:              "reminder-5",
					Title:           "Test Reminder 5",
					DeliveryStatus:  models.DeliveryStatusRead,
					MessageSentAt:   time.Now().AddDate(0, 0, -2).Format(time.RFC3339),
					ReadAt:          time.Now().AddDate(0, 0, -2).Add(time.Minute * 5).Format(time.RFC3339),
				},
			},
		},
	}

	handler := NewAnalyticsHandler(patientStore)

	tests := []struct {
		name           string
		role           string
		period         string
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "admin access with all period",
			role:           "admin",
			period:         "all",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "superadmin access with 7d period",
			role:           "superadmin",
			period:         "7d",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "volunteer denied access",
			role:           "volunteer",
			period:         "all",
			expectedStatus: http.StatusForbidden,
			expectError:    true,
		},
		{
			name:           "today period",
			role:           "admin",
			period:         "today",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "30d period",
			role:           "admin",
			period:         "30d",
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set role in context
			c.Set("role", tt.role)

			// Set up request
			c.Request = httptest.NewRequest("GET", "/analytics/delivery?period="+tt.period, nil)

			// Call handler
			handler.GetDeliveryAnalytics(c)

			// Check status
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				var response map[string]string
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to parse error response: %v", err)
				}
				if response["error"] == "" {
					t.Error("Expected error message in response")
				}
			} else {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to parse success response: %v", err)
				}
				if response["data"] == nil {
					t.Error("Expected data field in response")
				}
			}
		})
	}
}

func TestParsePeriod(t *testing.T) {
	tests := []struct {
		period          string
		expectZeroStart bool
	}{
		{"all", true},
		{"", true},
		{"today", false},
		{"7d", false},
		{"30d", false},
	}

	for _, tt := range tests {
		t.Run(tt.period, func(t *testing.T) {
			startDate, endDate, err := parsePeriod(tt.period)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectZeroStart {
				if !startDate.IsZero() {
					t.Error("Expected zero start date for period:", tt.period)
				}
			} else {
				if startDate.IsZero() {
					t.Error("Expected non-zero start date for period:", tt.period)
				}
				if endDate.IsZero() {
					t.Error("Expected non-zero end date for period:", tt.period)
				}
				if !endDate.After(startDate) {
					t.Error("Expected end date to be after start date")
				}
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{time.Second * 30, "30s"},
		{time.Minute * 2 + time.Second * 30, "2m 30s"},
		{time.Hour*1 + time.Minute*30 + time.Second*15, "1h 30m 15s"},
		{time.Millisecond * 500, "< 1s"},
	}

	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			result := formatDuration(tt.input)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetDeliveryAnalytics_EmptyStore(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create empty patient store
	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	handler := NewAnalyticsHandler(patientStore)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "admin")
	c.Request = httptest.NewRequest("GET", "/analytics/delivery?period=all", nil)

	handler.GetDeliveryAnalytics(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	data := response["data"].(map[string]interface{})
	if data["totalSent"].(float64) != 0 {
		t.Error("Expected totalSent to be 0 for empty store")
	}
	if data["successRate"].(float64) != 0 {
		t.Error("Expected successRate to be 0 for empty store")
	}
}

func TestGetDeliveryAnalytics_InvalidPeriod(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	handler := NewAnalyticsHandler(patientStore)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "admin")
	c.Request = httptest.NewRequest("GET", "/analytics/delivery?period=invalid", nil)

	handler.GetDeliveryAnalytics(c)

	// Should default to "all" behavior (no error for invalid period)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestGetDeliveryAnalytics_AdminRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	handler := NewAnalyticsHandler(patientStore)

	roles := []string{"admin", "superadmin", "volunteer", "", "user"}

	for _, role := range roles {
		t.Run("role_"+role, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Set("role", role)
			c.Request = httptest.NewRequest("GET", "/analytics/delivery?period=all", nil)

			handler.GetDeliveryAnalytics(c)

			if role == "admin" || role == "superadmin" {
				if w.Code != http.StatusOK {
					t.Errorf("Expected 200 for role %s, got %d", role, w.Code)
				}
			} else {
				if w.Code != http.StatusForbidden {
					t.Errorf("Expected 403 for role %s, got %d", role, w.Code)
				}
			}
		})
	}
}

func TestGetDeliveryAnalytics_SuccessRateCalculation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	now := time.Now().UTC()

	// Create test data with known success/failure counts
	patientStore.Patients = map[string]*models.Patient{
		"patient-1": {
			ID:   "patient-1",
			Name: "Test",
			Reminders: []*models.Reminder{
				// 3 delivered, 1 read, 2 sent, 2 failed = 8 total, 4 success
				{
					ID:             "r1",
					DeliveryStatus: models.DeliveryStatusDelivered,
					MessageSentAt:  now.AddDate(0, 0, -1).Format(time.RFC3339),
				},
				{
					ID:             "r2",
					DeliveryStatus: models.DeliveryStatusDelivered,
					MessageSentAt:  now.AddDate(0, 0, -1).Format(time.RFC3339),
				},
				{
					ID:             "r3",
					DeliveryStatus: models.DeliveryStatusDelivered,
					MessageSentAt:  now.AddDate(0, 0, -1).Format(time.RFC3339),
				},
				{
					ID:             "r4",
					DeliveryStatus: models.DeliveryStatusRead,
					MessageSentAt:  now.AddDate(0, 0, -1).Format(time.RFC3339),
				},
				{
					ID:             "r5",
					DeliveryStatus: models.DeliveryStatusSent,
					MessageSentAt:  now.AddDate(0, 0, -1).Format(time.RFC3339),
				},
				{
					ID:             "r6",
					DeliveryStatus: models.DeliveryStatusSent,
					MessageSentAt:  now.AddDate(0, 0, -1).Format(time.RFC3339),
				},
				{
					ID:             "r7",
					DeliveryStatus: models.DeliveryStatusFailed,
					MessageSentAt:  now.AddDate(0, 0, -1).Format(time.RFC3339),
				},
				{
					ID:             "r8",
					DeliveryStatus: models.DeliveryStatusFailed,
					MessageSentAt:  now.AddDate(0, 0, -1).Format(time.RFC3339),
				},
			},
		},
	}

	handler := NewAnalyticsHandler(patientStore)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("role", "admin")
	c.Request = httptest.NewRequest("GET", "/analytics/delivery?period=all", nil)

	handler.GetDeliveryAnalytics(c)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	data := response["data"].(map[string]interface{})

	// Check total sent count (delivered + read + sent + failed = 8)
	totalSent := int(data["totalSent"].(float64))
	if totalSent != 8 {
		t.Errorf("Expected totalSent = 8, got %d", totalSent)
	}

	// Check success rate = (3 + 1) / (3 + 1 + 2 + 2) * 100 = 4/8 * 100 = 50%
	successRate := data["successRate"].(float64)
	if successRate != 50.0 {
		t.Errorf("Expected successRate = 50.0, got %.2f", successRate)
	}

	breakdown := data["breakdown"].(map[string]interface{})
	if breakdown["delivered"].(float64) != 3 {
		t.Error("Expected 3 delivered")
	}
	if breakdown["read"].(float64) != 1 {
		t.Error("Expected 1 read")
	}
	if breakdown["sent"].(float64) != 2 {
		t.Error("Expected 2 sent")
	}
	if breakdown["failed"].(float64) != 2 {
		t.Error("Expected 2 failed")
	}
}

func BenchmarkGetDeliveryAnalytics(b *testing.B) {
	gin.SetMode(gin.TestMode)

	// Create patient store with 100 patients, 10 reminders each
	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	for i := 0; i < 100; i++ {
		reminders := make([]*models.Reminder, 10)
		for j := 0; j < 10; j++ {
			statuses := []string{
				models.DeliveryStatusPending,
				models.DeliveryStatusSent,
				models.DeliveryStatusDelivered,
				models.DeliveryStatusFailed,
			}
			reminders[j] = &models.Reminder{
				ID:              "r-" + strings.Repeat("0", i/100) + "-" + string(rune(j)),
				DeliveryStatus:  statuses[j%len(statuses)],
				MessageSentAt:   time.Now().Format(time.RFC3339),
			}
		}
		patientStore.Patients["patient-"+string(rune(i))] = &models.Patient{
			ID:       "patient-" + string(rune(i)),
			Reminders: reminders,
		}
	}

	handler := NewAnalyticsHandler(patientStore)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "admin")
		c.Request = httptest.NewRequest("GET", "/analytics/delivery?period=all", nil)

		handler.GetDeliveryAnalytics(c)
	}
}

// Test CategorizeFailureReason tests the error categorization function
func TestCategorizeFailureReason(t *testing.T) {
	tests := []struct {
		name         string
		errorMsg     string
		expectedCode string
	}{
		{"invalid phone", "invalid phone number", "invalid_phone"},
		{"nomor invalid", "nomor tidak valid", "invalid_phone"},
		{"timeout", "connection timeout", "gowa_timeout"},
		{"timeout explicit", "GOWA timeout after 30s", "gowa_timeout"},
		{"rejected", "message rejected by server", "message_rejected"},
		{"ditolak", "pesAN ditolak", "message_rejected"},
		{"other", "some other error", "other"},
		{"empty", "", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reasonCode, _ := categorizeFailureReason(tt.errorMsg)
			if reasonCode != tt.expectedCode {
				t.Errorf("categorizeFailureReason(%q) = %s, expected %s", tt.errorMsg, reasonCode, tt.expectedCode)
			}
		})
	}
}

// Test GetFailedDeliveries tests the failed deliveries handler
func TestGetFailedDeliveries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now().UTC()
	patientStore := models.NewPatientStore(func() {})

	// Create test data with failed deliveries
	patientStore.Patients = map[string]*models.Patient{
		"patient-1": {
			ID:        "patient-1",
			Name:      "John Doe",
			Phone:     "6281234567890",
			CreatedBy: "volunteer-1",
			Reminders: []*models.Reminder{
				{
					ID:                   "reminder-1",
					Title:                "Test Reminder 1",
					DeliveryStatus:       models.DeliveryStatusFailed,
					DeliveryErrorMessage: "Connection timeout",
					MessageSentAt:        now.AddDate(0, 0, -1).Format(time.RFC3339),
					RetryCount:           3,
				},
				{
					ID:                   "reminder-2",
					Title:                "Test Reminder 2",
					DeliveryStatus:       models.DeliveryStatusDelivered,
					MessageSentAt:        now.AddDate(0, 0, -1).Format(time.RFC3339),
				},
			},
		},
		"patient-2": {
			ID:        "patient-2",
			Name:      "Jane Smith",
			Phone:     "6280987654321",
			CreatedBy: "volunteer-2",
			Reminders: []*models.Reminder{
				{
					ID:                   "reminder-3",
					Title:                "Test Reminder 3",
					DeliveryStatus:       models.DeliveryStatusFailed,
					DeliveryErrorMessage: "Invalid phone number",
					MessageSentAt:        now.AddDate(0, 0, -2).Format(time.RFC3339),
					RetryCount:           1,
				},
				{
					ID:                   "reminder-4",
					Title:                "Test Reminder 4",
					DeliveryStatus:       models.DeliveryStatusFailed,
					DeliveryErrorMessage: "Message rejected",
					MessageSentAt:        now.AddDate(0, 0, -3).Format(time.RFC3339),
					RetryCount:           2,
				},
			},
		},
	}

	handler := NewAnalyticsHandler(patientStore)

	t.Run("admin access", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "admin")
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries", nil)

		handler.GetFailedDeliveries(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		if len(items) != 3 {
			t.Errorf("Expected 3 failed deliveries, got %d", len(items))
		}

		// Check that patient names are masked
		firstItem := items[0].(map[string]interface{})
		if firstItem["patient_name_masked"] == "John Doe" {
			t.Error("Patient name should be masked")
		}
	})

	t.Run("superadmin access", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "superadmin")
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries", nil)

		handler.GetFailedDeliveries(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 for superadmin, got %d", w.Code)
		}
	})

	t.Run("volunteer denied", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "volunteer")
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries", nil)

		handler.GetFailedDeliveries(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 for volunteer, got %d", w.Code)
		}
	})

	t.Run("filter by reason", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "admin")
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries?reason=invalid_phone", nil)

		handler.GetFailedDeliveries(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		data := response["data"].(map[string]interface{})
		items := data["items"].([]interface{})
		if len(items) != 1 {
			t.Errorf("Expected 1 failed delivery for invalid_phone filter, got %d", len(items))
		}
	})

	t.Run("pagination", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "admin")
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries?page=1&limit=2", nil)

		handler.GetFailedDeliveries(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		data := response["data"].(map[string]interface{})
		pagination := data["pagination"].(map[string]interface{})
		if pagination["page"].(float64) != 1 {
			t.Error("Expected page 1")
		}
		if pagination["limit"].(float64) != 2 {
			t.Error("Expected limit 2")
		}
		if pagination["total"].(float64) != 3 {
			t.Errorf("Expected total 3, got %v", pagination["total"])
		}
		if pagination["total_pages"].(float64) != 2 {
			t.Errorf("Expected total_pages 2, got %v", pagination["total_pages"])
		}
	})
}

// Test ExportFailedDeliveries tests the CSV export handler
func TestExportFailedDeliveries(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now().UTC()
	patientStore := models.NewPatientStore(func() {})

	patientStore.Patients = map[string]*models.Patient{
		"patient-1": {
			ID:        "patient-1",
			Name:      "John Doe",
			Phone:     "6281234567890",
			CreatedBy: "volunteer-1",
			Reminders: []*models.Reminder{
				{
					ID:                   "reminder-1",
					Title:                "Test Reminder 1",
					DeliveryStatus:       models.DeliveryStatusFailed,
					DeliveryErrorMessage: "Connection timeout",
					MessageSentAt:        now.Format(time.RFC3339),
					RetryCount:           3,
				},
			},
		},
	}

	handler := NewAnalyticsHandler(patientStore)

	t.Run("export CSV", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "admin")
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries/export", nil)

		handler.ExportFailedDeliveries(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check content type
		contentType := w.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/csv") {
			t.Errorf("Expected text/csv content type, got %s", contentType)
		}

		// Check content disposition
		disposition := w.Header().Get("Content-Disposition")
		if !strings.Contains(disposition, "failed-deliveries-") {
			t.Error("Expected Content-Disposition with filename")
		}

		// Check CSV content
		body := w.Body.String()
		lines := strings.Split(body, "\n")
		if len(lines) < 2 {
			t.Error("Expected at least header and data row")
		}

		// Check header
		if !strings.Contains(lines[0], "Reminder ID") {
			t.Error("Expected header with 'Reminder ID'")
		}

		// Check data row has masked patient name
		if strings.Contains(lines[1], "John Doe") {
			t.Error("Patient name should be masked in CSV")
		}
	})

	t.Run("volunteer denied", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "volunteer")
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries/export", nil)

		handler.ExportFailedDeliveries(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 for volunteer, got %d", w.Code)
		}
	})
}

// Test GetFailedDeliveryDetail tests the detail handler
func TestGetFailedDeliveryDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	now := time.Now().UTC()
	patientStore := models.NewPatientStore(func() {})

	patientStore.Patients = map[string]*models.Patient{
		"patient-1": {
			ID:        "patient-1",
			Name:      "John Doe",
			Phone:     "6281234567890",
			CreatedBy: "volunteer-1",
			Reminders: []*models.Reminder{
				{
					ID:                   "reminder-1",
					Title:                "Test Reminder 1",
					DeliveryStatus:       models.DeliveryStatusFailed,
					DeliveryErrorMessage: "Connection timeout",
					MessageSentAt:        now.Format(time.RFC3339),
					RetryCount:           3,
				},
			},
		},
	}

	handler := NewAnalyticsHandler(patientStore)

	t.Run("get detail", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "admin")
		c.Params = gin.Params{{Key: "id", Value: "reminder-1"}}
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries/reminder-1", nil)

		handler.GetFailedDeliveryDetail(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Errorf("Failed to parse response: %v", err)
		}

		data := response["data"].(map[string]interface{})
		if data["reminder_id"] != "reminder-1" {
			t.Error("Expected reminder_id to be reminder-1")
		}

		// Check patient name and phone are masked
		if data["patient_name_masked"] == "John Doe" {
			t.Error("Patient name should be masked")
		}
		if data["phone_masked"] == "6281234567890" {
			t.Error("Phone should be masked")
		}
	})

	t.Run("not found", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "admin")
		c.Params = gin.Params{{Key: "id", Value: "non-existent"}}
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries/non-existent", nil)

		handler.GetFailedDeliveryDetail(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})

	t.Run("volunteer denied", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Set("role", "volunteer")
		c.Params = gin.Params{{Key: "id", Value: "reminder-1"}}
		c.Request = httptest.NewRequest("GET", "/analytics/failed-deliveries/reminder-1", nil)

		handler.GetFailedDeliveryDetail(c)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 for volunteer, got %d", w.Code)
		}
	})
}
