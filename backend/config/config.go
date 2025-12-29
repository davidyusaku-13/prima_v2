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
