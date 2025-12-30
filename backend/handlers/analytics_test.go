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
