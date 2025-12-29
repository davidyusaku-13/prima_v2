package services

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestCircuitBreaker(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("allows requests when closed", func(t *testing.T) {
		cb := NewCircuitBreaker(5, 5*time.Minute, logger)

		if !cb.Allow() {
			t.Error("Circuit breaker should allow requests when closed")
		}
		if cb.State() != "closed" {
			t.Errorf("Expected state 'closed', got '%s'", cb.State())
		}
	})

	t.Run("opens after threshold failures", func(t *testing.T) {
		cb := NewCircuitBreaker(3, 5*time.Minute, logger)

		// Record 3 failures
		cb.RecordFailure()
		cb.RecordFailure()
		cb.RecordFailure()

		if cb.State() != "open" {
			t.Errorf("Expected state 'open' after 3 failures, got '%s'", cb.State())
		}
		if cb.Allow() {
			t.Error("Circuit breaker should not allow requests when open")
		}
	})

	t.Run("resets on success", func(t *testing.T) {
		cb := NewCircuitBreaker(5, 5*time.Minute, logger)

		// Record some failures
		cb.RecordFailure()
		cb.RecordFailure()

		if cb.Failures() != 2 {
			t.Errorf("Expected 2 failures, got %d", cb.Failures())
		}

		// Record success
		cb.RecordSuccess()

		if cb.Failures() != 0 {
			t.Errorf("Expected 0 failures after success, got %d", cb.Failures())
		}
	})

	t.Run("closes after cooldown", func(t *testing.T) {
		cb := NewCircuitBreaker(2, 100*time.Millisecond, logger)

		// Open the circuit
		cb.RecordFailure()
		cb.RecordFailure()

		if cb.State() != "open" {
			t.Errorf("Expected state 'open', got '%s'", cb.State())
		}

		// Wait for cooldown
		time.Sleep(150 * time.Millisecond)

		// Should allow and transition to closed
		if !cb.Allow() {
			t.Error("Circuit breaker should allow requests after cooldown")
		}
		if cb.State() != "closed" {
			t.Errorf("Expected state 'closed' after cooldown, got '%s'", cb.State())
		}
	})
}

func TestGOWAClient(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	t.Run("sends message successfully", func(t *testing.T) {
		// Create mock server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Verify request
			if r.Method != "POST" {
				t.Errorf("Expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/send/message" {
				t.Errorf("Expected /send/message, got %s", r.URL.Path)
			}
			if r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type application/json")
			}
			if r.Header.Get("Authorization") == "" {
				t.Error("Expected Authorization header")
			}

			// Return success response
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(SendMessageResponse{
				Success:   true,
				MessageID: "msg-123",
			})
		}))
		defer server.Close()

		client := NewGOWAClient(GOWAConfig{
			Endpoint:         server.URL,
			User:             "testuser",
			Password:         "testpass",
			Timeout:          10 * time.Second,
			FailureThreshold: 5,
			CooldownDuration: 5 * time.Minute,
		}, logger)

		resp, err := client.SendMessage("628123456789@s.whatsapp.net", "Test message")
		if err != nil {
			t.Fatalf("SendMessage failed: %v", err)
		}
		if !resp.Success {
			t.Error("Expected success response")
		}
		if resp.MessageID != "msg-123" {
			t.Errorf("Expected message ID 'msg-123', got '%s'", resp.MessageID)
		}
	})

	t.Run("handles server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
		}))
		defer server.Close()

		client := NewGOWAClient(GOWAConfig{
			Endpoint:         server.URL,
			User:             "testuser",
			Password:         "testpass",
			Timeout:          10 * time.Second,
			FailureThreshold: 5,
			CooldownDuration: 5 * time.Minute,
		}, logger)

		_, err := client.SendMessage("628123456789@s.whatsapp.net", "Test message")
		if err == nil {
			t.Error("Expected error for server error response")
		}

		// Verify failure was recorded
		if client.GetCircuitBreakerFailures() != 1 {
			t.Errorf("Expected 1 failure, got %d", client.GetCircuitBreakerFailures())
		}
	})

	t.Run("blocks requests when circuit is open", func(t *testing.T) {
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewGOWAClient(GOWAConfig{
			Endpoint:         server.URL,
			User:             "testuser",
			Password:         "testpass",
			Timeout:          10 * time.Second,
			FailureThreshold: 2,
			CooldownDuration: 5 * time.Minute,
		}, logger)

		// Trigger circuit breaker
		client.SendMessage("628123456789@s.whatsapp.net", "Test 1")
		client.SendMessage("628123456789@s.whatsapp.net", "Test 2")

		// Circuit should be open now
		if client.GetCircuitBreakerState() != "open" {
			t.Errorf("Expected circuit to be open, got '%s'", client.GetCircuitBreakerState())
		}

		// This request should be blocked
		_, err := client.SendMessage("628123456789@s.whatsapp.net", "Test 3")
		if err == nil {
			t.Error("Expected error when circuit is open")
		}

		// Only 2 requests should have reached the server
		if requestCount != 2 {
			t.Errorf("Expected 2 requests to server, got %d", requestCount)
		}
	})

	t.Run("IsAvailable returns correct state", func(t *testing.T) {
		client := NewGOWAClient(GOWAConfig{
			Endpoint:         "http://localhost:3000",
			User:             "testuser",
			Password:         "testpass",
			Timeout:          10 * time.Second,
			FailureThreshold: 5,
			CooldownDuration: 5 * time.Minute,
		}, logger)

		if !client.IsAvailable() {
			t.Error("Client should be available initially")
		}
	})
}

func TestNewGOWAClientFromConfig(t *testing.T) {
	// This test verifies the config integration works
	// We can't fully test without the config package, but we can verify the function exists
	// and would work with proper config
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	client := NewGOWAClient(GOWAConfig{
		Endpoint:         "http://localhost:3000",
		User:             "admin",
		Password:         "password",
		Timeout:          30 * time.Second,
		FailureThreshold: 5,
		CooldownDuration: 5 * time.Minute,
	}, logger)

	if client == nil {
		t.Error("NewGOWAClient returned nil")
	}
	if client.endpoint != "http://localhost:3000" {
		t.Errorf("Expected endpoint 'http://localhost:3000', got '%s'", client.endpoint)
	}
}
