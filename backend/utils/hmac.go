package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// ValidateWebhookSignature validates the HMAC-SHA256 signature of a webhook payload.
// Returns true if the signature is valid, false otherwise.
// Uses constant-time comparison to prevent timing attacks.
func ValidateWebhookSignature(payload []byte, signature, secret string) bool {
	// Security: Reject empty secret - webhook should not be processed without a secret
	if secret == "" {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))

	// Use constant-time comparison to prevent timing attacks
	return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// GenerateWebhookSignature generates an HMAC-SHA256 signature for testing purposes.
func GenerateWebhookSignature(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
