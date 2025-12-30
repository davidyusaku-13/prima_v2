package utils

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Indonesian format with 62", "628123456789", "628***789"},
		{"Indonesian format with 08", "08123456789", "081***789"},
		{"Short number", "12345", "***"},
		{"Very short", "123", "***"},
		{"Empty string", "", "***"},
		{"With spaces", " 628123456789 ", "628***789"},
		{"Longer number", "6281234567890123", "628***123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskPhone(tt.input)
			if result != tt.expected {
				t.Errorf("MaskPhone(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaskEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal email", "user@example.com", "u***@example.com"},
		{"Short local part", "a@example.com", "a***@example.com"},
		{"No @ symbol", "invalid-email", "***"},
		{"Empty string", "", "***"},
		{"With spaces", " user@example.com ", "u***@example.com"},
		{"Long local part", "verylongusername@example.com", "v***@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskEmail(tt.input)
			if result != tt.expected {
				t.Errorf("MaskEmail(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name      string
		level     string
		format    string
		logLevel  string // level to log at
	}{
		{"JSON info logger", "info", "json", "info"},
		{"Text debug logger", "debug", "text", "debug"},
		{"JSON warn logger", "warn", "json", "warn"},
		{"JSON error logger", "error", "json", "error"},
		{"Default level", "", "json", "info"},
		{"Default format", "info", "", "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(LoggerConfig{
				Level:  tt.level,
				Format: tt.format,
				Output: &buf,
			})

			if logger == nil {
				t.Error("NewLogger returned nil")
			}

			// Log at the appropriate level
			switch tt.logLevel {
			case "debug":
				logger.Debug("test message", "key", "value")
			case "info":
				logger.Info("test message", "key", "value")
			case "warn":
				logger.Warn("test message", "key", "value")
			case "error":
				logger.Error("test message", "key", "value")
			}

			output := buf.String()
			if output == "" {
				t.Error("Logger did not produce any output")
			}

			// Verify JSON format produces valid JSON
			if tt.format == "json" || tt.format == "" {
				var jsonOutput map[string]interface{}
				if err := json.Unmarshal([]byte(output), &jsonOutput); err != nil {
					t.Errorf("JSON logger output is not valid JSON: %v", err)
				}
			}
		})
	}
}

func TestLoggerLevels(t *testing.T) {
	// Test that debug messages are filtered at info level
	var buf bytes.Buffer
	logger := NewLogger(LoggerConfig{
		Level:  "info",
		Format: "json",
		Output: &buf,
	})

	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("Debug message should not appear at info level")
	}

	logger.Info("info message")
	if buf.Len() == 0 {
		t.Error("Info message should appear at info level")
	}
}

func TestInitDefaultLogger(t *testing.T) {
	// Save original logger
	originalLogger := DefaultLogger

	// Initialize with new settings
	InitDefaultLogger("debug", "text")

	if DefaultLogger == nil {
		t.Error("DefaultLogger should not be nil after initialization")
	}

	if DefaultLogger == originalLogger {
		t.Error("DefaultLogger should be a new instance after initialization")
	}

	// Restore original logger
	DefaultLogger = originalLogger
}

func TestSetDefaultLogger(t *testing.T) {
	// Save original logger
	originalLogger := DefaultLogger

	// Create and set a custom logger
	var buf bytes.Buffer
	customLogger := NewLogger(LoggerConfig{
		Level:  "debug",
		Format: "json",
		Output: &buf,
	})

	SetDefaultLogger(customLogger)

	if DefaultLogger != customLogger {
		t.Error("DefaultLogger should be the custom logger after SetDefaultLogger")
	}

	// Restore original logger
	DefaultLogger = originalLogger
}

func TestMaskPhoneEdgeCases(t *testing.T) {
	// Test with exactly 6 characters (boundary case)
	result := MaskPhone("123456")
	if result != "123***456" {
		t.Errorf("MaskPhone(\"123456\") = %q, want \"123***456\"", result)
	}

	// Test with 7 characters
	result = MaskPhone("1234567")
	if !strings.HasPrefix(result, "123") || !strings.HasSuffix(result, "567") {
		t.Errorf("MaskPhone(\"1234567\") = %q, should start with 123 and end with 567", result)
	}
}
