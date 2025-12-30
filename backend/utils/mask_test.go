package utils

import "testing"

func TestMaskPatientName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal name", "Budiman", "Bu***n"},
		{"short name 2 chars", "Bu", "B***u"},
		{"short name 1 char", "B", "B***"},
		{"empty string", "", ""},
		{"whitespace trimmed", "  Budi  ", "Bu***i"},
		{"three chars", "Budi", "Bu***i"},
		{"four chars", "Budi", "Bu***i"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskPatientName(tt.input)
			if result != tt.expected {
				t.Errorf("MaskPatientName(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMaskPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal phone", "6281234567890", "628***890"},
		{"short phone 6 chars", "628123", "628123"},
		{"short phone 5 chars", "62812", "62812"},
		{"empty string", "", ""},
		{"whitespace trimmed", "  6281234567890  ", "628***890"},
		{"phone with plus", "+6281234567890", "+62***890"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskPhoneNumber(tt.input)
			if result != tt.expected {
				t.Errorf("MaskPhoneNumber(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func BenchmarkMaskPatientName(b *testing.B) {
	names := []string{"Budiman", "Ahmad Fauzi", "Sri Dewi Lestari", "Muhammad Akbar"}

	for i := 0; i < b.N; i++ {
		for _, name := range names {
			MaskPatientName(name)
		}
	}
}

func BenchmarkMaskPhoneNumber(b *testing.B) {
	phones := []string{"6281234567890", "6280987654321", "6285212345678"}

	for i := 0; i < b.N; i++ {
		for _, phone := range phones {
			MaskPhoneNumber(phone)
		}
	}
}
