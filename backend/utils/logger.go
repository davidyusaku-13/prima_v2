package utils

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level  string
	Format string
	Output io.Writer
}

// NewLogger creates a new slog.Logger based on configuration
func NewLogger(cfg LoggerConfig) *slog.Logger {
	// Default output to stdout
	if cfg.Output == nil {
		cfg.Output = os.Stdout
	}

	// Parse log level
	var level slog.Level
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level: level,
	}

	// Create handler based on format
	var handler slog.Handler
	switch strings.ToLower(cfg.Format) {
	case "json":
		handler = slog.NewJSONHandler(cfg.Output, opts)
	case "text":
		handler = slog.NewTextHandler(cfg.Output, opts)
	default:
		handler = slog.NewJSONHandler(cfg.Output, opts)
	}

	return slog.New(handler)
}

// MaskPhone masks phone number for logging to comply with NFR-S6
// Input: "628123456789" → Output: "628***789"
// Input: "08123456789" → Output: "081***789"
func MaskPhone(phone string) string {
	// Remove any whitespace
	phone = strings.TrimSpace(phone)

	if len(phone) < 6 {
		return "***"
	}

	// Show first 3 and last 3 characters
	return phone[:3] + "***" + phone[len(phone)-3:]
}

// MaskEmail masks email for logging
// Input: "user@example.com" → Output: "u***@example.com"
func MaskEmail(email string) string {
	email = strings.TrimSpace(email)

	atIndex := strings.Index(email, "@")
	if atIndex <= 0 {
		return "***"
	}

	localPart := email[:atIndex]
	domain := email[atIndex:]

	if len(localPart) <= 1 {
		return localPart + "***" + domain
	}

	return localPart[:1] + "***" + domain
}

// DefaultLogger is the global logger instance
var DefaultLogger *slog.Logger

// InitDefaultLogger initializes the default logger with given configuration
func InitDefaultLogger(level, format string) {
	DefaultLogger = NewLogger(LoggerConfig{
		Level:  level,
		Format: format,
	})
}

// SetDefaultLogger sets a custom logger as the default
func SetDefaultLogger(logger *slog.Logger) {
	DefaultLogger = logger
}

func init() {
	// Initialize with default settings
	DefaultLogger = NewLogger(LoggerConfig{
		Level:  "info",
		Format: "json",
	})
}
