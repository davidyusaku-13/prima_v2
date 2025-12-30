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

	// Verify disclaimer defaults
	expectedDisclaimerText := "Informasi ini untuk tujuan edukasi. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda."
	if cfg.Disclaimer.Text != expectedDisclaimerText {
		t.Errorf("Expected default disclaimer text '%s', got '%s'", expectedDisclaimerText, cfg.Disclaimer.Text)
	}
	if cfg.Disclaimer.Enabled == nil {
		t.Error("Expected disclaimer Enabled to be set, got nil")
	} else if !*cfg.Disclaimer.Enabled {
		t.Error("Expected disclaimer Enabled to be true by default")
	}
}

func TestApplyDefaults_DisclaimerPreservesExplicitValues(t *testing.T) {
	// Test that explicit values are NOT overwritten by defaults
	enabled := false
	cfg := &Config{
		Disclaimer: DisclaimerConfig{
			Text:    "Custom disclaimer",
			Enabled: &enabled,
		},
	}
	cfg.applyDefaults()

	// Custom text should be preserved
	if cfg.Disclaimer.Text != "Custom disclaimer" {
		t.Errorf("Expected custom disclaimer text to be preserved, got '%s'", cfg.Disclaimer.Text)
	}

	// Explicit false should be preserved (not overwritten to true)
	if cfg.Disclaimer.Enabled == nil {
		t.Error("Expected disclaimer Enabled to remain set")
	} else if *cfg.Disclaimer.Enabled != false {
		t.Error("Expected explicit Enabled=false to be preserved, not overwritten to true")
	}
}

func TestApplyDefaults_QuietHours(t *testing.T) {
	cfg := &Config{}
	cfg.applyDefaults()

	// Verify quiet hours defaults
	if cfg.QuietHours.GetStartHour() != 21 {
		t.Errorf("Expected default quiet hours start_hour 21, got %d", cfg.QuietHours.GetStartHour())
	}
	if cfg.QuietHours.GetEndHour() != 6 {
		t.Errorf("Expected default quiet hours end_hour 6, got %d", cfg.QuietHours.GetEndHour())
	}
	if cfg.QuietHours.Timezone != "WIB" {
		t.Errorf("Expected default quiet hours timezone 'WIB', got '%s'", cfg.QuietHours.Timezone)
	}
}

func TestApplyDefaults_QuietHoursPreservesExplicitValues(t *testing.T) {
	// Test that explicit values are NOT overwritten by defaults
	startHour := 22
	endHour := 7
	cfg := &Config{
		QuietHours: QuietHoursConfig{
			StartHour: &startHour,
			EndHour:   &endHour,
			Timezone:  "WITA",
		},
	}
	cfg.applyDefaults()

	// Custom values should be preserved
	if cfg.QuietHours.GetStartHour() != 22 {
		t.Errorf("Expected custom start_hour 22 to be preserved, got %d", cfg.QuietHours.GetStartHour())
	}
	if cfg.QuietHours.GetEndHour() != 7 {
		t.Errorf("Expected custom end_hour 7 to be preserved, got %d", cfg.QuietHours.GetEndHour())
	}
	if cfg.QuietHours.Timezone != "WITA" {
		t.Errorf("Expected custom timezone 'WITA' to be preserved, got '%s'", cfg.QuietHours.Timezone)
	}
}

func TestApplyDefaults_QuietHoursPreservesZeroValues(t *testing.T) {
	// Test that explicit zero values (midnight) are NOT overwritten by defaults
	startHour := 0  // Midnight
	endHour := 5
	cfg := &Config{
		QuietHours: QuietHoursConfig{
			StartHour: &startHour,
			EndHour:   &endHour,
			Timezone:  "WIB",
		},
	}
	cfg.applyDefaults()

	// Zero value should be preserved (not overwritten to 21)
	if cfg.QuietHours.GetStartHour() != 0 {
		t.Errorf("Expected start_hour 0 (midnight) to be preserved, got %d", cfg.QuietHours.GetStartHour())
	}
	if cfg.QuietHours.GetEndHour() != 5 {
		t.Errorf("Expected end_hour 5 to be preserved, got %d", cfg.QuietHours.GetEndHour())
	}
}

func TestLoadQuietHoursFromFile(t *testing.T) {
	content := `
quiet_hours:
  start_hour: 20
  end_hour: 5
  timezone: "WIT"
`
	tmpFile, err := os.CreateTemp("", "config-quiethours-*.yaml")
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

	if cfg.QuietHours.GetStartHour() != 20 {
		t.Errorf("Expected start_hour 20, got %d", cfg.QuietHours.GetStartHour())
	}
	if cfg.QuietHours.GetEndHour() != 5 {
		t.Errorf("Expected end_hour 5, got %d", cfg.QuietHours.GetEndHour())
	}
	if cfg.QuietHours.Timezone != "WIT" {
		t.Errorf("Expected timezone 'WIT', got '%s'", cfg.QuietHours.Timezone)
	}
}

func TestQuietHoursValidation_ValidConfig(t *testing.T) {
	startHour := 21
	endHour := 6
	cfg := &QuietHoursConfig{
		StartHour: &startHour,
		EndHour:   &endHour,
		Timezone:  "WIB",
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Expected valid config, got error: %v", err)
	}
}

func TestQuietHoursValidation_InvalidStartHour(t *testing.T) {
	startHour := 25 // Invalid: > 23
	endHour := 6
	cfg := &QuietHoursConfig{
		StartHour: &startHour,
		EndHour:   &endHour,
		Timezone:  "WIB",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected error for invalid start_hour, got nil")
	}
}

func TestQuietHoursValidation_InvalidEndHour(t *testing.T) {
	startHour := 21
	endHour := -1 // Invalid: < 0
	cfg := &QuietHoursConfig{
		StartHour: &startHour,
		EndHour:   &endHour,
		Timezone:  "WIB",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected error for invalid end_hour, got nil")
	}
}

func TestQuietHoursValidation_InvalidTimezone(t *testing.T) {
	startHour := 21
	endHour := 6
	cfg := &QuietHoursConfig{
		StartHour: &startHour,
		EndHour:   &endHour,
		Timezone:  "INVALID",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Expected error for invalid timezone, got nil")
	}
}

func TestQuietHoursValidation_EmptyTimezoneIsValid(t *testing.T) {
	// Empty timezone is valid (will use default WIB)
	startHour := 21
	endHour := 6
	cfg := &QuietHoursConfig{
		StartHour: &startHour,
		EndHour:   &endHour,
		Timezone:  "",
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Expected empty timezone to be valid, got error: %v", err)
	}
}
