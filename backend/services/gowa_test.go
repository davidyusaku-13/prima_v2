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

func TestGetRetryDelay(t *testing.T) {
	delays := []time.Duration{
		1 * time.Second,
		5 * time.Second,
		30 * time.Second,
	}

	tests := []struct {
		name     string
		attempt  int
		expected time.Duration
	}{
		{"first attempt (0)", 0, 1 * time.Second},
		{"first retry (1)", 1, 1 * time.Second},
		{"second retry (2)", 2, 5 * time.Second},
		{"third retry (3)", 3, 30 * time.Second},
		{"exceeding retries (4)", 4, 30 * time.Second},
		{"exceeding retries (10)", 10, 30 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRetryDelay(tt.attempt, delays)
			if result != tt.expected {
				t.Errorf("GetRetryDelay(%d) = %v, want %v", tt.attempt, result, tt.expected)
			}
		})
	}
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		shouldRetry bool
	}{
		// Non-retryable errors
		{
			name:        "nil error",
			err:         nil,
			shouldRetry: false,
		},
		{
			name:        "circuit breaker open",
			err:         &testError{msg: "circuit breaker is open"},
			shouldRetry: false,
		},
		{
			name:        "invalid phone",
			err:         &testError{msg: "invalid phone number"},
			shouldRetry: false,
		},
		{
			name:        "nomor error",
			err:         &testError{msg: "nomor HP tidak valid"},
			shouldRetry: false,
		},
		{
			name:        "400 Bad Request",
			err:         &testError{msg: "400 Bad Request"},
			shouldRetry: false,
		},
		{
			name:        "401 Unauthorized",
			err:         &testError{msg: "401 Unauthorized"},
			shouldRetry: false,
		},
		{
			name:        "403 Forbidden",
			err:         &testError{msg: "403 Forbidden"},
			shouldRetry: false,
		},
		{
			name:        "404 Not Found",
			err:         &testError{msg: "404 Not Found"},
			shouldRetry: false,
		},
		// Retryable errors
		{
			name:        "timeout",
			err:         &testError{msg: "connection timeout"},
			shouldRetry: true,
		},
		{
			name:        "context deadline exceeded",
			err:         &testError{msg: "context deadline exceeded"},
			shouldRetry: true,
		},
		{
			name:        "connection refused",
			err:         &testError{msg: "connection refused"},
			shouldRetry: true,
		},
		{
			name:        "connection reset",
			err:         &testError{msg: "connection reset by peer"},
			shouldRetry: true,
		},
		{
			name:        "no such host",
			err:         &testError{msg: "dial tcp: lookup gowa.example.com: no such host"},
			shouldRetry: true,
		},
		{
			name:        "network unreachable",
			err:         &testError{msg: "network is unreachable"},
			shouldRetry: true,
		},
		{
			name:        "i/o timeout",
			err:         &testError{msg: "read tcp: i/o timeout"},
			shouldRetry: true,
		},
		{
			name:        "EOF error",
			err:         &testError{msg: "unexpected EOF"},
			shouldRetry: true,
		},
		{
			name:        "500 Internal Server Error",
			err:         &testError{msg: "500 Internal Server Error"},
			shouldRetry: true,
		},
		{
			name:        "502 Bad Gateway",
			err:         &testError{msg: "502 Bad Gateway"},
			shouldRetry: true,
		},
		{
			name:        "503 Service Unavailable",
			err:         &testError{msg: "503 Service Unavailable"},
			shouldRetry: true,
		},
		{
			name:        "504 Gateway Timeout",
			err:         &testError{msg: "504 Gateway Timeout"},
			shouldRetry: true,
		},
		{
			name:        "server misbehaving",
			err:         &testError{msg: "server misbehaving"},
			shouldRetry: true,
		},
		// Unknown errors - default to retry
		{
			name:        "unknown error",
			err:         &testError{msg: "something unexpected happened"},
			shouldRetry: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldRetry(tt.err)
			if result != tt.shouldRetry {
				t.Errorf("ShouldRetry(%v) = %v, want %v", tt.err, result, tt.shouldRetry)
			}
		})
	}
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestContainsString(t *testing.T) {
	tests := []struct {
		s        string
		pattern  string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "hello", true},
		{"hello world", "lo wo", true},
		{"hello world", "xyz", false},
		{"", "test", false},
		{"test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.pattern, func(t *testing.T) {
			result := containsString(tt.s, tt.pattern)
			if result != tt.expected {
				t.Errorf("containsString(%q, %q) = %v, want %v", tt.s, tt.pattern, result, tt.expected)
			}
		})
	}
}
