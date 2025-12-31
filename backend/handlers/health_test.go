package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/davidyusaku-13/prima_v2/handlers"
	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/services"
)

// Helper to create a test gin context
func createTestContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	return c, w
}

func TestGetHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock patient store and gowa client
	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	// Test basic health endpoint
	c, w := createTestContext("GET", "/api/health")
	healthHandler.GetHealth(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "Health check successful" {
		t.Errorf("Expected success message, got %v", response["message"])
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data field in response")
	}

	if data["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", data["status"])
	}

	if data["timestamp"] == "" {
		t.Error("Expected timestamp to be present")
	}
}

func TestGetHealthDetailed_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	// Test detailed health endpoint without admin role
	c, w := createTestContext("GET", "/api/health/detailed")
	c.Set("role", "volunteer")
	healthHandler.GetHealthDetailed(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestGetHealthDetailed_Admin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	// Test detailed health endpoint with admin role
	c, w := createTestContext("GET", "/api/health/detailed")
	c.Set("role", "admin")
	healthHandler.GetHealthDetailed(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "Detailed health status retrieved" {
		t.Errorf("Expected detailed health message, got %v", response["message"])
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data field in response")
	}

	// Check GOWA section
	gowa, ok := data["gowa"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected gowa field in data")
	}
	if gowa["endpoint"] != "http://localhost:3000" {
		t.Errorf("Expected endpoint 'http://localhost:3000', got %v", gowa["endpoint"])
	}

	// Check circuit breaker section
	circuitBreaker, ok := data["circuit_breaker"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected circuit_breaker field in data")
	}
	if circuitBreaker["state"] != "closed" {
		t.Errorf("Expected circuit breaker state 'closed', got %v", circuitBreaker["state"])
	}

	// Check queue section
	queue, ok := data["queue"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected queue field in data")
	}
	if queue["total"] != float64(0) {
		t.Errorf("Expected queue total 0, got %v", queue["total"])
	}
}

func TestGetHealthDetailed_Superadmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})

	// Add a patient with reminders for testing queue counts
	patientStore.Patients["patient1"] = &models.Patient{
		ID:   "patient1",
		Name: "Test Patient",
		Reminders: []*models.Reminder{
			{
				ID:             "reminder1",
				DeliveryStatus: models.DeliveryStatusPending,
			},
			{
				ID:             "reminder2",
				DeliveryStatus: models.DeliveryStatusFailed,
				RetryCount:     1,
			},
		},
	}

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	// Test detailed health endpoint with superadmin role
	c, w := createTestContext("GET", "/api/health/detailed")
	c.Set("role", "superadmin")
	healthHandler.GetHealthDetailed(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	queue := data["queue"].(map[string]interface{})

	// Should have 2 reminders in queue (1 pending + 1 retrying)
	if queue["total"] != float64(2) {
		t.Errorf("Expected queue total 2, got %v", queue["total"])
	}
}

func TestUpdateGOWAPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	// Update GOWA ping status
	healthHandler.UpdateGOWAPing(true)

	// Get health details to verify
	c, w := createTestContext("GET", "/api/health/detailed")
	c.Set("role", "admin")
	healthHandler.GetHealthDetailed(c)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	gowa := data["gowa"].(map[string]interface{})

	if gowa["connected"] != true {
		t.Errorf("Expected GOWA connected to be true")
	}
}

func TestGetQueueCounts_Empty(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	// Test with empty patient store
	c, w := createTestContext("GET", "/api/health/detailed")
	c.Set("role", "admin")
	healthHandler.GetHealthDetailed(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	queue := data["queue"].(map[string]interface{})

	if queue["total"] != float64(0) {
		t.Errorf("Expected queue total 0, got %v", queue["total"])
	}
}

func TestGetCircuitBreakerDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	details := gowaClient.GetCircuitBreakerDetails()

	if details.State != "closed" {
		t.Errorf("Expected state 'closed', got %v", details.State)
	}
	if details.FailureCount != 0 {
		t.Errorf("Expected failure count 0, got %v", details.FailureCount)
	}
	// When circuit breaker is closed, cooldown remaining should be 0
	if details.CooldownRemaining != 0 {
		t.Errorf("Expected cooldown remaining 0 for closed state, got %v", details.CooldownRemaining)
	}
}

func TestGetHealthDetailed_CircuitBreakerOpen(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	// Manually set circuit breaker state to open
	gowaClient.SetCircuitBreakerStateForTest("open", 3, 4*time.Minute)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	c, w := createTestContext("GET", "/api/health/detailed")
	c.Set("role", "admin")
	healthHandler.GetHealthDetailed(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	circuitBreaker := data["circuit_breaker"].(map[string]interface{})

	if circuitBreaker["state"] != "open" {
		t.Errorf("Expected circuit breaker state 'open', got %v", circuitBreaker["state"])
	}
	if circuitBreaker["failure_count"] != float64(3) {
		t.Errorf("Expected failure count 3, got %v", circuitBreaker["failure_count"])
	}
	cooldownSecs, ok := circuitBreaker["cooldown_remaining_seconds"].(float64)
	if !ok {
		t.Errorf("Expected cooldown_remaining_seconds to be a number")
	}
	if cooldownSecs <= 0 {
		t.Errorf("Expected cooldown remaining > 0, got %v", cooldownSecs)
	}
}

func TestGetHealthDetailed_CircuitBreakerHalfOpen(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	patientStore.Patients = make(map[string]*models.Patient)

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	// Manually set circuit breaker state to half-open
	gowaClient.SetCircuitBreakerStateForTest("half-open", 0, 0)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	c, w := createTestContext("GET", "/api/health/detailed")
	c.Set("role", "admin")
	healthHandler.GetHealthDetailed(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	circuitBreaker := data["circuit_breaker"].(map[string]interface{})

	if circuitBreaker["state"] != "half-open" {
		t.Errorf("Expected circuit breaker state 'half-open', got %v", circuitBreaker["state"])
	}
}

func TestGetQueueCounts_WithQuietHours(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})
	now := time.Now().UTC()

	// Add patient with scheduled delivery for future (quiet hours)
	futureTime := now.Add(2 * time.Hour).Format(time.RFC3339)
	pastTime := now.Add(-1 * time.Hour).Format(time.RFC3339)

	patientStore.Patients["patient1"] = &models.Patient{
		ID:   "patient1",
		Name: "Test Patient",
		Reminders: []*models.Reminder{
			{
				ID:                  "reminder1",
				DeliveryStatus:      models.DeliveryStatusScheduled,
				ScheduledDeliveryAt: futureTime, // Future = quiet hours
			},
			{
				ID:             "reminder2",
				DeliveryStatus: models.DeliveryStatusScheduled,
				ScheduledDeliveryAt: pastTime, // Past = scheduled (not quiet hours)
			},
		},
	}

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	c, w := createTestContext("GET", "/api/health/detailed")
	c.Set("role", "admin")
	healthHandler.GetHealthDetailed(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	queue := data["queue"].(map[string]interface{})

	// Should have 2 reminders total
	if queue["total"] != float64(2) {
		t.Errorf("Expected queue total 2, got %v", queue["total"])
	}
	// 1 quiet hours (future scheduled)
	if queue["quiet_hours"] != float64(1) {
		t.Errorf("Expected quiet_hours 1, got %v", queue["quiet_hours"])
	}
	// 1 scheduled (past scheduled)
	if queue["scheduled"] != float64(1) {
		t.Errorf("Expected scheduled 1, got %v", queue["scheduled"])
	}
}

func TestGetQueueCounts_RetryingExhausted(t *testing.T) {
	gin.SetMode(gin.TestMode)

	patientStore := models.NewPatientStore(func() {})

	// Add patient with failed reminders - some with retries left, some exhausted
	patientStore.Patients["patient1"] = &models.Patient{
		ID:   "patient1",
		Name: "Test Patient",
		Reminders: []*models.Reminder{
			{
				ID:             "reminder1",
				DeliveryStatus: models.DeliveryStatusFailed,
				RetryCount:     1, // Has retries left
			},
			{
				ID:             "reminder2",
				DeliveryStatus: models.DeliveryStatusFailed,
				RetryCount:     3, // Exhausted retries (max is 3)
			},
		},
	}

	gowaClient := services.NewGOWAClient(services.GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "test",
		Password:         "test",
		Timeout:          5 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, nil)

	healthHandler := handlers.NewHealthHandler(patientStore, gowaClient)

	c, w := createTestContext("GET", "/api/health/detailed")
	c.Set("role", "admin")
	healthHandler.GetHealthDetailed(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	data := response["data"].(map[string]interface{})
	queue := data["queue"].(map[string]interface{})

	// Should only have 1 reminder (the one with retries left)
	if queue["total"] != float64(1) {
		t.Errorf("Expected queue total 1, got %v", queue["total"])
	}
	if queue["retrying"] != float64(1) {
		t.Errorf("Expected retrying 1, got %v", queue["retrying"])
	}
}
