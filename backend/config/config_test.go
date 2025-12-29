package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	content := `
server:
  port: 9090
  cors_origin: "http://example.com"

gowa:
  endpoint: "http://gowa.example.com"
  user: "testuser"
  password: "testpass"
  webhook_secret: "secret123"
  timeout: 60s

circuit_breaker:
  failure_threshold: 10
  cooldown_duration: 10m

retry:
  max_attempts: 3
  delays:
    - 2s
    - 10s

logging:
  level: "debug"
  format: "text"
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	// Test loading the config
	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify server config
	if cfg.Server.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Server.CORSOrigin != "http://example.com" {
		t.Errorf("Expected CORS origin 'http://example.com', got '%s'", cfg.Server.CORSOrigin)
	}

	// Verify GOWA config
	if cfg.GOWA.Endpoint != "http://gowa.example.com" {
		t.Errorf("Expected GOWA endpoint 'http://gowa.example.com', got '%s'", cfg.GOWA.Endpoint)
	}
	if cfg.GOWA.User != "testuser" {
		t.Errorf("Expected GOWA user 'testuser', got '%s'", cfg.GOWA.User)
	}
	if cfg.GOWA.Timeout != 60*time.Second {
		t.Errorf("Expected GOWA timeout 60s, got %v", cfg.GOWA.Timeout)
	}

	// Verify circuit breaker config
	if cfg.CircuitBreaker.FailureThreshold != 10 {
		t.Errorf("Expected failure threshold 10, got %d", cfg.CircuitBreaker.FailureThreshold)
	}
	if cfg.CircuitBreaker.CooldownDuration != 10*time.Minute {
		t.Errorf("Expected cooldown duration 10m, got %v", cfg.CircuitBreaker.CooldownDuration)
	}

	// Verify retry config
	if cfg.Retry.MaxAttempts != 3 {
		t.Errorf("Expected max attempts 3, got %d", cfg.Retry.MaxAttempts)
	}
	if len(cfg.Retry.Delays) != 2 {
		t.Errorf("Expected 2 retry delays, got %d", len(cfg.Retry.Delays))
	}

	// Verify logging config
	if cfg.Logging.Level != "debug" {
		t.Errorf("Expected logging level 'debug', got '%s'", cfg.Logging.Level)
	}
	if cfg.Logging.Format != "text" {
		t.Errorf("Expected logging format 'text', got '%s'", cfg.Logging.Format)
	}
}

func TestLoadOrDefault(t *testing.T) {
	// Test with non-existent file - should return defaults
	cfg := LoadOrDefault("non-existent-file.yaml")

	// Verify defaults are applied
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.CORSOrigin != "http://localhost:5173" {
		t.Errorf("Expected default CORS origin 'http://localhost:5173', got '%s'", cfg.Server.CORSOrigin)
	}
	if cfg.GOWA.Endpoint != "http://localhost:3000" {
		t.Errorf("Expected default GOWA endpoint 'http://localhost:3000', got '%s'", cfg.GOWA.Endpoint)
	}
	if cfg.GOWA.Timeout != 30*time.Second {
		t.Errorf("Expected default GOWA timeout 30s, got %v", cfg.GOWA.Timeout)
	}
	if cfg.CircuitBreaker.FailureThreshold != 5 {
		t.Errorf("Expected default failure threshold 5, got %d", cfg.CircuitBreaker.FailureThreshold)
	}
	if cfg.CircuitBreaker.CooldownDuration != 5*time.Minute {
		t.Errorf("Expected default cooldown duration 5m, got %v", cfg.CircuitBreaker.CooldownDuration)
	}
	if cfg.Retry.MaxAttempts != 5 {
		t.Errorf("Expected default max attempts 5, got %d", cfg.Retry.MaxAttempts)
	}
	if len(cfg.Retry.Delays) != 5 {
		t.Errorf("Expected 5 default retry delays, got %d", len(cfg.Retry.Delays))
	}
	if cfg.Logging.Level != "info" {
		t.Errorf("Expected default logging level 'info', got '%s'", cfg.Logging.Level)
	}
	if cfg.Logging.Format != "json" {
		t.Errorf("Expected default logging format 'json', got '%s'", cfg.Logging.Format)
	}
}

func TestLoadWithEnvExpansion(t *testing.T) {
	// Set environment variable
	os.Setenv("TEST_GOWA_SECRET", "env-secret-value")
	defer os.Unsetenv("TEST_GOWA_SECRET")

	content := `
gowa:
  webhook_secret: "${TEST_GOWA_SECRET}"
`
	tmpFile, err := os.CreateTemp("", "config-env-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.GOWA.WebhookSecret != "env-secret-value" {
		t.Errorf("Expected webhook secret 'env-secret-value', got '%s'", cfg.GOWA.WebhookSecret)
	}
}

func TestApplyDefaults(t *testing.T) {
	cfg := &Config{}
	cfg.applyDefaults()

	// Verify all defaults are applied
	if cfg.Server.Port != 8080 {
		t.Errorf("Expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.GOWA.Endpoint != "http://localhost:3000" {
		t.Errorf("Expected default GOWA endpoint, got '%s'", cfg.GOWA.Endpoint)
	}
	if cfg.CircuitBreaker.FailureThreshold != 5 {
		t.Errorf("Expected default failure threshold 5, got %d", cfg.CircuitBreaker.FailureThreshold)
	}
	if cfg.Retry.MaxAttempts != 5 {
		t.Errorf("Expected default max attempts 5, got %d", cfg.Retry.MaxAttempts)
	}
}
