package utils

import (
	"regexp"
	"strings"
)

// Indonesian phone number patterns
var (
	// Valid Indonesian mobile prefixes (after 62)
	validMobilePrefixes = []string{
		"811", "812", "813", "814", "815", "816", "817", "818", "819", // Telkomsel
		"821", "822", "823", "824", "825", "826", "827", "828", "829", // Telkomsel
		"851", "852", "853", "854", "855", "856", "857", "858", "859", // Telkomsel
		"831", "832", "833", "834", "835", "836", "837", "838", "839", // Axis
		"881", "882", "883", "884", "885", "886", "887", "888", "889", // Smartfren
		"895", "896", "897", "898", "899", // Three
	}

	// Regex for basic phone number validation (digits only, 10-15 chars)
	phoneDigitsRegex = regexp.MustCompile(`^\d{10,15}$`)
)

// PhoneValidationResult contains the result of phone validation
type PhoneValidationResult struct {
	Valid         bool   `json:"valid"`
	Normalized    string `json:"normalized,omitempty"`
	WhatsAppFormat string `json:"whatsapp_format,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
}

// ValidatePhoneNumber validates an Indonesian phone number
// Returns validation result with normalized format
func ValidatePhoneNumber(phone string) PhoneValidationResult {
	// Clean the phone number
	cleaned := cleanPhoneNumber(phone)

	if cleaned == "" {
		return PhoneValidationResult{
			Valid:        false,
			ErrorMessage: "phone number is empty",
		}
	}

	// Normalize to 62xxx format
	normalized := normalizeToIndonesian(cleaned)

	// Check minimum length (62 + 9 digits minimum)
	if len(normalized) < 11 {
		return PhoneValidationResult{
			Valid:        false,
			ErrorMessage: "phone number is too short",
		}
	}

	// Check maximum length (62 + 13 digits maximum)
	if len(normalized) > 15 {
		return PhoneValidationResult{
			Valid:        false,
			ErrorMessage: "phone number is too long",
		}
	}

	// Validate it's all digits
	if !phoneDigitsRegex.MatchString(normalized) {
		return PhoneValidationResult{
			Valid:        false,
			ErrorMessage: "phone number contains invalid characters",
		}
	}

	// Check if it starts with 62
	if !strings.HasPrefix(normalized, "62") {
		return PhoneValidationResult{
			Valid:        false,
			ErrorMessage: "phone number must be Indonesian format (starting with 62 or 08)",
		}
	}

	return PhoneValidationResult{
		Valid:          true,
		Normalized:     normalized,
		WhatsAppFormat: normalized + "@s.whatsapp.net",
	}
}

// cleanPhoneNumber removes all non-digit characters except leading +
func cleanPhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)

	// Remove common separators and whitespace
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, ".", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Handle + prefix
	if strings.HasPrefix(phone, "+") {
		phone = phone[1:]
	}

	return phone
}

// normalizeToIndonesian converts various Indonesian phone formats to 62xxx
func normalizeToIndonesian(phone string) string {
	// Already in 62xxx format
	if strings.HasPrefix(phone, "62") {
		return phone
	}

	// Convert 08xxx to 628xxx
	if strings.HasPrefix(phone, "08") {
		return "62" + phone[1:]
	}

	// Convert 8xxx to 628xxx (without leading 0)
	if strings.HasPrefix(phone, "8") {
		return "62" + phone
	}

	// Return as-is if doesn't match known patterns
	return phone
}

// FormatWhatsAppNumber formats a phone number for WhatsApp API
// Input: various formats (08xxx, +62xxx, 62xxx)
// Output: 62xxx@s.whatsapp.net
func FormatWhatsAppNumber(phone string) string {
	result := ValidatePhoneNumber(phone)
	if result.Valid {
		return result.WhatsAppFormat
	}

	// Fallback: try basic normalization even if validation fails
	cleaned := cleanPhoneNumber(phone)
	normalized := normalizeToIndonesian(cleaned)
	return normalized + "@s.whatsapp.net"
}

// IsValidIndonesianMobile checks if the phone number is a valid Indonesian mobile number
// This is a stricter check that validates the mobile prefix
func IsValidIndonesianMobile(phone string) bool {
	result := ValidatePhoneNumber(phone)
	if !result.Valid {
		return false
	}

	// Extract the prefix after 62 (first 3 digits)
	if len(result.Normalized) < 5 {
		return false
	}

	prefix := result.Normalized[2:5] // Get digits after "62"

	for _, validPrefix := range validMobilePrefixes {
		if prefix == validPrefix {
			return true
		}
	}

	return false
}

// NormalizePhoneNumber normalizes a phone number to 62xxx format
// Returns empty string if the phone number is invalid
func NormalizePhoneNumber(phone string) string {
	result := ValidatePhoneNumber(phone)
	if result.Valid {
		return result.Normalized
	}
	return ""
}
