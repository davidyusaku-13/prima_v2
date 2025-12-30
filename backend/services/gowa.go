package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/davidyusaku-13/prima_v2/utils"
)

// CircuitBreaker implements the circuit breaker pattern for GOWA service
type CircuitBreaker struct {
	mu               sync.Mutex
	failures         int
	lastFailure      time.Time
	state            string // "closed", "open"
	threshold        int
	cooldownDuration time.Duration
	logger           *slog.Logger
}

// NewCircuitBreaker creates a new circuit breaker with the given configuration
func NewCircuitBreaker(threshold int, cooldownDuration time.Duration, logger *slog.Logger) *CircuitBreaker {
	return &CircuitBreaker{
		state:            "closed",
		threshold:        threshold,
		cooldownDuration: cooldownDuration,
		logger:           logger,
	}
}

// Allow checks if a request is allowed through the circuit breaker
func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == "open" {
		// Check if cooldown period has passed
		if time.Since(cb.lastFailure) > cb.cooldownDuration {
			cb.state = "closed"
			cb.failures = 0
			cb.logger.Info("Circuit breaker state transition",
				"from", "open",
				"to", "closed",
				"reason", "cooldown_expired",
			)
			return true
		}
		return false
	}
	return true
}

// RecordFailure records a failure and potentially opens the circuit
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.threshold {
		previousState := cb.state
		cb.state = "open"
		if previousState != "open" {
			cb.logger.Warn("Circuit breaker state transition",
				"from", previousState,
				"to", "open",
				"failures", cb.failures,
				"threshold", cb.threshold,
			)
		}
	}
}

// RecordSuccess records a successful request and resets the failure count
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() string {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

// Failures returns the current failure count
func (cb *CircuitBreaker) Failures() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.failures
}

// GOWAClient is a client for the GOWA WhatsApp gateway service
type GOWAClient struct {
	endpoint       string
	user           string
	password       string
	timeout        time.Duration
	circuitBreaker *CircuitBreaker
	httpClient     *http.Client
	logger         *slog.Logger
}

// GOWAConfig holds configuration for the GOWA client
type GOWAConfig struct {
	Endpoint         string
	User             string
	Password         string
	Timeout          time.Duration
	FailureThreshold int
	CooldownDuration time.Duration
}

// NewGOWAClient creates a new GOWA client with the given configuration
func NewGOWAClient(cfg GOWAConfig, logger *slog.Logger) *GOWAClient {
	if logger == nil {
		logger = utils.DefaultLogger
	}

	return &GOWAClient{
		endpoint: cfg.Endpoint,
		user:     cfg.User,
		password: cfg.Password,
		timeout:  cfg.Timeout,
		circuitBreaker: NewCircuitBreaker(
			cfg.FailureThreshold,
			cfg.CooldownDuration,
			logger,
		),
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		logger: logger,
	}
}

// NewGOWAClientFromConfig creates a GOWA client from application config
func NewGOWAClientFromConfig(cfg *config.Config, logger *slog.Logger) *GOWAClient {
	return NewGOWAClient(GOWAConfig{
		Endpoint:         cfg.GOWA.Endpoint,
		User:             cfg.GOWA.User,
		Password:         cfg.GOWA.Password,
		Timeout:          cfg.GOWA.Timeout,
		FailureThreshold: cfg.CircuitBreaker.FailureThreshold,
		CooldownDuration: cfg.CircuitBreaker.CooldownDuration,
	}, logger)
}

// SendMessageRequest represents a request to send a WhatsApp message
type SendMessageRequest struct {
	Phone   string `json:"phone"`
	Message string `json:"message"`
}

// SendMessageResponse represents the response from GOWA after sending a message
type SendMessageResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Error     string `json:"error,omitempty"`
}

// SendMessage sends a WhatsApp message via GOWA
func (c *GOWAClient) SendMessage(phone, message string) (*SendMessageResponse, error) {
	// Check circuit breaker
	if !c.circuitBreaker.Allow() {
		c.logger.Warn("GOWA request blocked by circuit breaker",
			"phone", utils.MaskPhone(phone),
			"circuit_state", c.circuitBreaker.State(),
		)
		return nil, fmt.Errorf("circuit breaker is open, GOWA service temporarily unavailable")
	}

	// Prepare request
	endpoint := c.endpoint + "/send/message"
	payload := SendMessageRequest{
		Phone:   phone,
		Message: message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	auth := c.user + ":" + c.password
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))

	// Log request (with masked phone)
	c.logger.Debug("Sending GOWA request",
		"endpoint", endpoint,
		"phone", utils.MaskPhone(phone),
	)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.circuitBreaker.RecordFailure()
		c.logger.Error("GOWA request failed",
			"error", err.Error(),
			"phone", utils.MaskPhone(phone),
			"circuit_failures", c.circuitBreaker.Failures(),
		)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.circuitBreaker.RecordFailure()
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		c.circuitBreaker.RecordFailure()
		c.logger.Error("GOWA returned non-OK status",
			"status_code", resp.StatusCode,
			"phone", utils.MaskPhone(phone),
			"response", string(body),
			"circuit_failures", c.circuitBreaker.Failures(),
		)
		return nil, fmt.Errorf("GOWA returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result SendMessageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		// If we can't parse but got 200, consider it a success
		c.circuitBreaker.RecordSuccess()
		c.logger.Info("GOWA message sent (unparseable response)",
			"phone", utils.MaskPhone(phone),
			"status_code", resp.StatusCode,
		)
		return &SendMessageResponse{Success: true}, nil
	}

	// Record success
	c.circuitBreaker.RecordSuccess()
	c.logger.Info("GOWA message sent successfully",
		"phone", utils.MaskPhone(phone),
		"message_id", result.MessageID,
	)

	return &result, nil
}

// IsAvailable checks if the GOWA service is available (circuit breaker is closed)
func (c *GOWAClient) IsAvailable() bool {
	return c.circuitBreaker.Allow()
}

// GetCircuitBreakerState returns the current state of the circuit breaker
func (c *GOWAClient) GetCircuitBreakerState() string {
	return c.circuitBreaker.State()
}

// GetCircuitBreakerFailures returns the current failure count
func (c *GOWAClient) GetCircuitBreakerFailures() int {
	return c.circuitBreaker.Failures()
}

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts int
	Delays      []time.Duration
}

// GetRetryDelay returns the delay duration for a given attempt (1-indexed)
// Uses exponential backoff: 1s, 5s, 30s for attempts 1, 2, 3
func GetRetryDelay(attempt int, delays []time.Duration) time.Duration {
	if attempt <= 0 {
		return delays[0]
	}
	idx := attempt - 1
	if idx >= len(delays) {
		idx = len(delays) - 1
	}
	return delays[idx]
}

// ShouldRetry determines if an error is retryable
// Returns true for transient errors like network timeouts, connection issues
func ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Non-retryable errors
	nonRetryable := []string{
		"circuit breaker is open",
		"invalid phone",
		"nomor",
		"400",
		"401",
		"403",
		"404",
	}

	for _, pattern := range nonRetryable {
		if containsString(errStr, pattern) {
			return false
		}
	}

	// Retryable errors - network issues, timeouts, server errors
	retryable := []string{
		"timeout",
		"context deadline exceeded",
		"connection refused",
		"connection reset",
		"no such host",
		"network is unreachable",
		"i/o timeout",
		"EOF",
		"server misbehaving",
		"500",
		"502",
		"503",
		"504",
	}

	for _, pattern := range retryable {
		if containsString(errStr, pattern) {
			return true
		}
	}

	// Default to retry for unknown errors (conservative approach)
	return true
}

// containsString checks if a string contains a pattern (case-insensitive for some)
func containsString(s, pattern string) bool {
	return len(s) >= len(pattern) && (s == pattern || len(s) > len(pattern) && containsSubstring(s, pattern))
}

// containsSubstring is a simple substring check
func containsSubstring(s, pattern string) bool {
	for i := 0; i <= len(s)-len(pattern); i++ {
		if s[i:i+len(pattern)] == pattern {
			return true
		}
	}
	return false
}
