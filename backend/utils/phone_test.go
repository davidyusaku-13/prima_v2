package utils

import (
	"testing"
)

func TestValidatePhoneNumber(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectValid    bool
		expectNorm     string
		expectWA       string
		expectError    string
	}{
		{
			name:        "Valid 08 format",
			input:       "08123456789",
			expectValid: true,
			expectNorm:  "628123456789",
			expectWA:    "628123456789@s.whatsapp.net",
		},
		{
			name:        "Valid 62 format",
			input:       "628123456789",
			expectValid: true,
			expectNorm:  "628123456789",
			expectWA:    "628123456789@s.whatsapp.net",
		},
		{
			name:        "Valid +62 format",
			input:       "+628123456789",
			expectValid: true,
			expectNorm:  "628123456789",
			expectWA:    "628123456789@s.whatsapp.net",
		},
		{
			name:        "Valid with spaces",
			input:       "0812 3456 789",
			expectValid: true,
			expectNorm:  "628123456789",
			expectWA:    "628123456789@s.whatsapp.net",
		},
		{
			name:        "Valid with dashes",
			input:       "0812-3456-789",
			expectValid: true,
			expectNorm:  "628123456789",
			expectWA:    "628123456789@s.whatsapp.net",
		},
		{
			name:        "Valid with parentheses",
			input:       "(0812) 3456789",
			expectValid: true,
			expectNorm:  "628123456789",
			expectWA:    "628123456789@s.whatsapp.net",
		},
		{
			name:        "Valid 8 format (without leading 0)",
			input:       "8123456789",
			expectValid: true,
			expectNorm:  "628123456789",
			expectWA:    "628123456789@s.whatsapp.net",
		},
		{
			name:        "Empty string",
			input:       "",
			expectValid: false,
			expectError: "phone number is empty",
		},
		{
			name:        "Too short",
			input:       "08123",
			expectValid: false,
			expectError: "phone number is too short",
		},
		{
			name:        "Too long",
			input:       "081234567890123456",
			expectValid: false,
			expectError: "phone number is too long",
		},
		{
			name:        "With letters",
			input:       "0812abc3456",
			expectValid: false,
			expectError: "phone number contains invalid characters",
		},
		{
			name:        "Whitespace only",
			input:       "   ",
			expectValid: false,
			expectError: "phone number is empty",
		},
		{
			name:        "Valid 13 digit number",
			input:       "0812345678901",
			expectValid: true,
			expectNorm:  "62812345678901",
			expectWA:    "62812345678901@s.whatsapp.net",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePhoneNumber(tt.input)

			if result.Valid != tt.expectValid {
				t.Errorf("ValidatePhoneNumber(%q).Valid = %v, want %v", tt.input, result.Valid, tt.expectValid)
			}

			if tt.expectValid {
				if result.Normalized != tt.expectNorm {
					t.Errorf("ValidatePhoneNumber(%q).Normalized = %q, want %q", tt.input, result.Normalized, tt.expectNorm)
				}
				if result.WhatsAppFormat != tt.expectWA {
					t.Errorf("ValidatePhoneNumber(%q).WhatsAppFormat = %q, want %q", tt.input, result.WhatsAppFormat, tt.expectWA)
				}
			} else {
				if result.ErrorMessage != tt.expectError {
					t.Errorf("ValidatePhoneNumber(%q).ErrorMessage = %q, want %q", tt.input, result.ErrorMessage, tt.expectError)
				}
			}
		})
	}
}

func TestFormatWhatsAppNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"08 format", "08123456789", "628123456789@s.whatsapp.net"},
		{"62 format", "628123456789", "628123456789@s.whatsapp.net"},
		{"+62 format", "+628123456789", "628123456789@s.whatsapp.net"},
		{"With spaces", "0812 3456 789", "628123456789@s.whatsapp.net"},
		{"With dashes", "0812-3456-789", "628123456789@s.whatsapp.net"},
		{"8 format", "8123456789", "628123456789@s.whatsapp.net"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatWhatsAppNumber(tt.input)
			if result != tt.expected {
				t.Errorf("FormatWhatsAppNumber(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidIndonesianMobile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Telkomsel 0812", "08123456789", true},
		{"Telkomsel 0813", "08133456789", true},
		{"Telkomsel 0852", "08523456789", true},
		{"Indosat 0814", "08143456789", true},
		{"XL 0817", "08173456789", true},
		{"Axis 0831", "08313456789", true},
		{"Three 0895", "08953456789", true},
		{"Smartfren 0881", "08813456789", true},
		{"Invalid prefix 0800", "08003456789", false},
		{"Invalid prefix 0801", "08013456789", false},
		{"Empty", "", false},
		{"Too short", "0812345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidIndonesianMobile(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidIndonesianMobile(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizePhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"08 format", "08123456789", "628123456789"},
		{"62 format", "628123456789", "628123456789"},
		{"+62 format", "+628123456789", "628123456789"},
		{"Invalid", "abc", ""},
		{"Empty", "", ""},
		{"Too short", "0812", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePhoneNumber(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizePhoneNumber(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCleanPhoneNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0812 3456 789", "08123456789"},
		{"0812-3456-789", "08123456789"},
		{"(0812) 3456789", "08123456789"},
		{"+628123456789", "628123456789"},
		{"  08123456789  ", "08123456789"},
		{"0812.3456.789", "08123456789"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := cleanPhoneNumber(tt.input)
			if result != tt.expected {
				t.Errorf("cleanPhoneNumber(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeToIndonesian(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"08123456789", "628123456789"},
		{"8123456789", "628123456789"},
		{"628123456789", "628123456789"},
		{"123456789", "123456789"}, // Non-Indonesian format returned as-is
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeToIndonesian(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeToIndonesian(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
