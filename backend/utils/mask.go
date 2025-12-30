package utils

import (
	"strings"
)

// MaskPatientName masks a patient name for privacy
// "Budiman" -> "Bu***n"
// "Bu" -> "B***u"
func MaskPatientName(name string) string {
	if name == "" {
		return ""
	}

	name = strings.TrimSpace(name)

	if len(name) <= 2 {
		if len(name) == 1 {
			return name[0:1] + "***"
		}
		return name[0:1] + "***" + name[len(name)-1:]
	}

	return name[0:2] + "***" + name[len(name)-1:]
}

// MaskPhoneNumber masks a phone number for privacy
// "6281234567890" -> "628***890"
func MaskPhoneNumber(phone string) string {
	if phone == "" {
		return ""
	}

	phone = strings.TrimSpace(phone)

	if len(phone) <= 6 {
		return phone
	}

	return phone[0:3] + "***" + phone[len(phone)-3:]
}
