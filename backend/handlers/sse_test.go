package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/gin-gonic/gin"
)

func TestSSEHandler_BroadcastDeliveryStatusUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{}
	handler := NewSSEHandler(cfg, nil)

	// Test broadcasting with no clients connected
	handler.BroadcastDeliveryStatusUpdate("reminder-123", "delivered", time.Now().UTC().Format(time.RFC3339))

	// Verify no panic occurs when no clients are connected
	if handler.GetClientCount() != 0 {
		t.Errorf("Expected 0 clients, got %d", handler.GetClientCount())
	}
}

func TestSSEHandler_HandleDeliveryStatusSSE_NoAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{}
	handler := NewSSEHandler(cfg, nil)

	router := gin.New()
	router.GET("/sse", handler.HandleDeliveryStatusSSE)

	req := httptest.NewRequest("GET", "/sse", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should send error event when no user_id in context
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check for SSE headers (allow charset in content-type)
	contentType := w.Header().Get("Content-Type")
	if contentType != "text/event-stream" && contentType != "text/event-stream;charset=utf-8" {
		t.Errorf("Expected Content-Type text/event-stream, got %s", contentType)
	}
}

func TestSSEHandler_HandleDeliveryStatusSSE_WithAuth(t *testing.T) {
	// Skip this test as httptest.ResponseRecorder doesn't support CloseNotifier
	// SSE functionality is tested manually and in integration tests
	t.Skip("Skipping SSE connection test - httptest.ResponseRecorder doesn't support CloseNotifier")
}

func TestSSEEvent_JSONMarshaling(t *testing.T) {
	event := SSEEvent{
		Event: "delivery.status.updated",
		Data: map[string]string{
			"reminder_id": "test-123",
			"status":      "delivered",
			"timestamp":   "2025-12-30T06:00:00Z",
		},
	}

	data, err := json.Marshal(event.Data)
	if err != nil {
		t.Fatalf("Failed to marshal event data: %v", err)
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal event data: %v", err)
	}

	if result["reminder_id"] != "test-123" {
		t.Errorf("Expected reminder_id test-123, got %s", result["reminder_id"])
	}

	if result["status"] != "delivered" {
		t.Errorf("Expected status delivered, got %s", result["status"])
	}
}

func TestSSEHandler_GetClientCount(t *testing.T) {
	cfg := &config.Config{}
	handler := NewSSEHandler(cfg, nil)

	// Initially should have 0 clients
	if count := handler.GetClientCount(); count != 0 {
		t.Errorf("Expected 0 clients initially, got %d", count)
	}

	// Manually add a client for testing
	handler.mu.Lock()
	testChan := make(chan SSEEvent, 10)
	handler.clients[testChan] = "test-user"
	handler.mu.Unlock()

	if count := handler.GetClientCount(); count != 1 {
		t.Errorf("Expected 1 client after adding, got %d", count)
	}

	// Clean up
	handler.mu.Lock()
	delete(handler.clients, testChan)
	close(testChan)
	handler.mu.Unlock()
}
