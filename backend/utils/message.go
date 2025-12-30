package utils

import (
	"fmt"
	"strings"
)

// MaxExcerptLength is the maximum length for article excerpts
const MaxExcerptLength = 100

// TruncateString truncates a string to maxLength and adds "..." if truncated
// Tries to truncate at word boundary when possible
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}

	// If maxLength is too small, just truncate
	if maxLength < 3 {
		return s[:maxLength] + "..."
	}

	// Try to find a space to truncate at word boundary
	truncated := strings.TrimRight(s[:maxLength], " ")
	if idx := strings.LastIndex(truncated, " "); idx > maxLength/2 {
		return strings.TrimRight(s[:idx], " ") + "..."
	}

	return strings.TrimRight(s[:maxLength], " ") + "..."
}

// ReminderMessageParams holds the parameters for formatting a reminder message
type ReminderMessageParams struct {
	PatientName         string
	ReminderTitle       string
	ReminderDescription string
	Attachments         []string // Pre-formatted attachment strings
	DisclaimerText      string
	DisclaimerEnabled   bool
}

// ContentAttachment represents an attachment with content details for message formatting
type ContentAttachment struct {
	Type     string // "article" or "video"
	Title    string
	Excerpt  string // Only for articles
	URL      string
	Slug     string // For article URL generation
	YouTubeID string // For video URL generation
}

// FormatReminderMessageWithExcerpts creates a WhatsApp message with content excerpts
// Articles show: title + excerpt (max 100 chars) + link
// Videos show: title + link only
func FormatReminderMessageWithExcerpts(params ReminderMessageParams, attachments []ContentAttachment) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Halo %s,\n\n", params.PatientName))
	sb.WriteString(fmt.Sprintf("*%s*\n\n", params.ReminderTitle))

	if params.ReminderDescription != "" {
		sb.WriteString(params.ReminderDescription)
		sb.WriteString("\n\n")
	}

	// Add attachments with excerpts if present
	if len(attachments) > 0 {
		sb.WriteString("---\n")
		sb.WriteString("Konten Edukasi:\n")

		for _, att := range attachments {
			if att.Type == "article" {
				// Format: 📖 Title\nExcerpt...\n🔗 link
				sb.WriteString(fmt.Sprintf("📖 %s\n", att.Title))

				// Add excerpt (truncated to 100 chars)
				// If excerpt is empty, use title as fallback per AC #1
				excerptText := att.Excerpt
				if excerptText == "" {
					excerptText = att.Title
				}
				truncated := TruncateString(excerptText, MaxExcerptLength)
				sb.WriteString(fmt.Sprintf("%s\n", truncated))

				// Add link
				if att.URL != "" {
					sb.WriteString(fmt.Sprintf("🔗 %s\n", att.URL))
				}
			} else if att.Type == "video" {
				// Format: 🎬 Title\n🔗 link
				sb.WriteString(fmt.Sprintf("🎬 %s\n", att.Title))

				if att.URL != "" {
					sb.WriteString(fmt.Sprintf("🔗 %s\n", att.URL))
				}
			}
			sb.WriteString("\n")
		}
	}

	// Add health disclaimer if enabled
	if params.DisclaimerEnabled && params.DisclaimerText != "" {
		sb.WriteString("---\n")
		sb.WriteString(fmt.Sprintf("_%s_", params.DisclaimerText))
	}

	return sb.String()
}

// FormatReminderMessage creates the WhatsApp message content for a reminder
func FormatReminderMessage(params ReminderMessageParams) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Halo %s,\n\n", params.PatientName))
	sb.WriteString(fmt.Sprintf("*%s*\n\n", params.ReminderTitle))

	if params.ReminderDescription != "" {
		sb.WriteString(params.ReminderDescription)
		sb.WriteString("\n\n")
	}

	// Add attachments if present
	if len(params.Attachments) > 0 {
		sb.WriteString("---\n")
		sb.WriteString("Konten Edukasi:\n")
		for _, attachment := range params.Attachments {
			sb.WriteString(attachment)
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Add health disclaimer if enabled
	if params.DisclaimerEnabled && params.DisclaimerText != "" {
		sb.WriteString("---\n")
		sb.WriteString(fmt.Sprintf("_%s_", params.DisclaimerText))
	}

	return sb.String()
}
