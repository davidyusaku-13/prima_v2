package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
type Config struct {
	Server         ServerConfig         `yaml:"server"`
	GOWA           GOWAConfig           `yaml:"gowa"`
	CircuitBreaker CircuitBreakerConfig `yaml:"circuit_breaker"`
	Retry          RetryConfig          `yaml:"retry"`
	Logging        LoggingConfig        `yaml:"logging"`
	Disclaimer     DisclaimerConfig     `yaml:"disclaimer"`
	QuietHours     QuietHoursConfig     `yaml:"quiet_hours"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port       int    `yaml:"port"`
	CORSOrigin string `yaml:"cors_origin"`
}

// GOWAConfig holds GOWA service configuration
type GOWAConfig struct {
	Endpoint      string        `yaml:"endpoint"`
	User          string        `yaml:"user"`
	Password      string        `yaml:"password"`
	WebhookSecret string        `yaml:"webhook_secret"`
	Timeout       time.Duration `yaml:"timeout"`
}

// CircuitBreakerConfig holds circuit breaker settings
type CircuitBreakerConfig struct {
	FailureThreshold int           `yaml:"failure_threshold"`
	CooldownDuration time.Duration `yaml:"cooldown_duration"`
}

// RetryConfig holds retry settings
type RetryConfig struct {
	MaxAttempts int             `yaml:"max_attempts"`
	Delays      []time.Duration `yaml:"delays"`
}

// LoggingConfig holds logging settings
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// DisclaimerConfig holds health disclaimer settings
type DisclaimerConfig struct {
	Text    string `yaml:"text"`
	Enabled *bool  `yaml:"enabled"`
}

// QuietHoursConfig holds quiet hours settings for reminder delivery
type QuietHoursConfig struct {
	StartHour *int   `yaml:"start_hour"` // 21 (9 PM) - pointer to distinguish 0 from unset
	EndHour   *int   `yaml:"end_hour"`   // 6 (6 AM) - pointer to distinguish 0 from unset
	Timezone  string `yaml:"timezone"`   // "WIB" (UTC+7)
}

// GetStartHour returns the start hour value, defaulting to 21 if not set
func (q *QuietHoursConfig) GetStartHour() int {
	if q.StartHour == nil {
		return 21
	}
	return *q.StartHour
}

// GetEndHour returns the end hour value, defaulting to 6 if not set
func (q *QuietHoursConfig) GetEndHour() int {
	if q.EndHour == nil {
		return 6
	}
	return *q.EndHour
}

// Validate checks if the quiet hours configuration is valid
func (q *QuietHoursConfig) Validate() error {
	startHour := q.GetStartHour()
	endHour := q.GetEndHour()

	if startHour < 0 || startHour > 23 {
		return fmt.Errorf("quiet_hours.start_hour must be between 0 and 23, got %d", startHour)
	}
	if endHour < 0 || endHour > 23 {
		return fmt.Errorf("quiet_hours.end_hour must be between 0 and 23, got %d", endHour)
	}

	validTimezones := map[string]bool{"WIB": true, "WITA": true, "WIT": true}
	if q.Timezone != "" && !validTimezones[q.Timezone] {
		return fmt.Errorf("quiet_hours.timezone must be one of WIB, WITA, WIT, got %s", q.Timezone)
	}

	return nil
}

// ValidateCircuitBreaker checks if the circuit breaker configuration is valid
func (c *CircuitBreakerConfig) Validate() error {
	if c.FailureThreshold <= 0 {
		return fmt.Errorf("circuit_breaker.failure_threshold must be > 0, got %d", c.FailureThreshold)
	}
	if c.CooldownDuration <= 0 {
		return fmt.Errorf("circuit_breaker.cooldown_duration must be > 0, got %v", c.CooldownDuration)
	}
	return nil
}

// ValidateRetry checks if the retry configuration is valid
func (r *RetryConfig) Validate() error {
	if r.MaxAttempts <= 0 {
		return fmt.Errorf("retry.max_attempts must be > 0, got %d", r.MaxAttempts)
	}
	for i, delay := range r.Delays {
		if delay < 0 {
			return fmt.Errorf("retry.delays[%d] must be >= 0, got %v", i, delay)
		}
	}
	return nil
}

// Load reads configuration from a YAML file and returns a Config struct
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables in the config
	expandedData := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expandedData), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	cfg.applyDefaults()

	// Validate circuit breaker config
	if err := cfg.CircuitBreaker.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Validate retry config
	if err := cfg.Retry.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Validate quiet hours config
	if err := cfg.QuietHours.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// applyDefaults sets default values for unset configuration options
func (c *Config) applyDefaults() {
	// Server defaults
	if c.Server.Port == 0 {
		c.Server.Port = 8080
	}
	if c.Server.CORSOrigin == "" {
		c.Server.CORSOrigin = "http://localhost:5173"
	}

	// GOWA defaults
	if c.GOWA.Endpoint == "" {
		c.GOWA.Endpoint = "http://localhost:3000"
	}
	if c.GOWA.Timeout == 0 {
		c.GOWA.Timeout = 30 * time.Second
	}

	// Circuit breaker defaults
	if c.CircuitBreaker.FailureThreshold == 0 {
		c.CircuitBreaker.FailureThreshold = 5
	}
	if c.CircuitBreaker.CooldownDuration == 0 {
		c.CircuitBreaker.CooldownDuration = 5 * time.Minute
	}

	// Retry defaults
	if c.Retry.MaxAttempts == 0 {
		c.Retry.MaxAttempts = 5
	}
	if len(c.Retry.Delays) == 0 {
		c.Retry.Delays = []time.Duration{
			1 * time.Second,
			5 * time.Second,
			30 * time.Second,
			2 * time.Minute,
			10 * time.Minute,
		}
	}

	// Logging defaults
	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}
	if c.Logging.Format == "" {
		c.Logging.Format = "json"
	}

	// Disclaimer defaults
	if c.Disclaimer.Text == "" {
		c.Disclaimer.Text = "Informasi ini untuk tujuan edukasi. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda."
	}
	if c.Disclaimer.Enabled == nil {
		enabled := true
		c.Disclaimer.Enabled = &enabled
	}

	// Quiet hours defaults
	if c.QuietHours.StartHour == nil {
		startHour := 21 // 9 PM WIB
		c.QuietHours.StartHour = &startHour
	}
	if c.QuietHours.EndHour == nil {
		endHour := 6 // 6 AM WIB
		c.QuietHours.EndHour = &endHour
	}
	if c.QuietHours.Timezone == "" {
		c.QuietHours.Timezone = "WIB" // UTC+7
	}
}

// LoadOrDefault attempts to load config from path, returns default config if file doesn't exist
func LoadOrDefault(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		// Return default config
		cfg = &Config{}
		cfg.applyDefaults()
	}
	return cfg
}
