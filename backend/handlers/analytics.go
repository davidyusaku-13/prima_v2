package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/utils"
)

// DeliveryAnalyticsStats represents the delivery statistics response
type DeliveryAnalyticsStats struct {
	TotalSent        int               `json:"totalSent"`
	SuccessRate      float64           `json:"successRate"`
	FailedLast7Days  int               `json:"failedLast7Days"`
	AvgDeliveryTime  string            `json:"avgDeliveryTime"` // Human readable, e.g., "2m 30s"
	Breakdown        map[string]int    `json:"breakdown"`       // status -> count
	Period           string            `json:"period"`          // today, 7d, 30d, all
	PeriodStartDate  string            `json:"periodStartDate"` // ISO 8601 UTC
	PeriodEndDate    string            `json:"periodEndDate"`   // ISO 8601 UTC
}

// AnalyticsHandler handles delivery analytics endpoints
type AnalyticsHandler struct {
	patientStore *models.PatientStore
}

// NewAnalyticsHandler creates a new analytics handler
func NewAnalyticsHandler(patientStore *models.PatientStore) *AnalyticsHandler {
	return &AnalyticsHandler{
		patientStore: patientStore,
	}
}

// parsePeriod parses the period query parameter and returns start/end dates
// period values: "today", "7d" (7 days), "30d" (30 days), "all"
func parsePeriod(period string) (startDate, endDate time.Time, err error) {
	now := time.Now().UTC()
	endDate = now

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	case "7d":
		startDate = now.AddDate(0, 0, -7)
	case "30d":
		startDate = now.AddDate(0, 0, -30)
	case "all", "":
		// Return zero time to indicate no filtering
		return time.Time{}, time.Time{}, nil
	default:
		return time.Time{}, time.Time{}, nil
	}

	return startDate, endDate, nil
}

// GetDeliveryAnalytics returns delivery statistics for admin dashboard
func (h *AnalyticsHandler) GetDeliveryAnalytics(c *gin.Context) {
	// Check admin role
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	period := c.DefaultQuery("period", "all")
	startDate, endDate, err := parsePeriod(period)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period parameter"})
		return
	}

	// Track counts
	totalSent := 0
	failedLast7Days := 0
	breakdown := make(map[string]int)

	// For delivery time calculation
	var totalDeliveryDurations []time.Duration
	now := time.Now().UTC()

	// Collect all reminders
	h.patientStore.RLock()
	for _, patient := range h.patientStore.Patients {
		for _, reminder := range patient.Reminders {
			// Skip reminders that haven't been sent yet
			if reminder.DeliveryStatus == "" ||
				reminder.DeliveryStatus == models.DeliveryStatusPending ||
				reminder.DeliveryStatus == models.DeliveryStatusScheduled {
				continue
			}

			// Filter by date if period is specified
			if !startDate.IsZero() && reminder.MessageSentAt != "" {
				sentTime, parseErr := time.Parse(time.RFC3339, reminder.MessageSentAt)
				if parseErr != nil {
					continue
				}
				if sentTime.Before(startDate) || sentTime.After(endDate) {
					continue
				}

				// Count failed in last 7 days
				if reminder.DeliveryStatus == models.DeliveryStatusFailed {
					if sentTime.After(now.AddDate(0, 0, -7)) {
						failedLast7Days++
					}
				}
			}

			// Only count sent/delivered/read/failed as "sent"
			switch reminder.DeliveryStatus {
			case models.DeliveryStatusSent,
				models.DeliveryStatusDelivered,
				models.DeliveryStatusRead,
				models.DeliveryStatusFailed:
				totalSent++
			}

			// Track breakdown by status
			breakdown[reminder.DeliveryStatus]++

			// Calculate delivery time (sent -> delivered/read)
			if reminder.MessageSentAt != "" &&
				(reminder.DeliveredAt != "" || reminder.ReadAt != "") {
				sentTime, parseErr := time.Parse(time.RFC3339, reminder.MessageSentAt)
				if parseErr == nil {
					var deliveredTime time.Time
					if reminder.DeliveredAt != "" {
						deliveredTime, _ = time.Parse(time.RFC3339, reminder.DeliveredAt)
					} else if reminder.ReadAt != "" {
						deliveredTime, _ = time.Parse(time.RFC3339, reminder.ReadAt)
					}

					if !deliveredTime.IsZero() && deliveredTime.After(sentTime) {
						totalDeliveryDurations = append(totalDeliveryDurations, deliveredTime.Sub(sentTime))
					}
				}
			}
		}
	}
	h.patientStore.RUnlock()

	// Calculate success rate: (delivered + read) / (sent + delivered + read + failed) * 100
	deliveredCount := breakdown[models.DeliveryStatusDelivered]
	readCount := breakdown[models.DeliveryStatusRead]
	failedCount := breakdown[models.DeliveryStatusFailed]
	sentCount := breakdown[models.DeliveryStatusSent]

	deliverableCount := deliveredCount + readCount + sentCount + failedCount
	var successRate float64
	if deliverableCount > 0 {
		successRate = float64(deliveredCount+readCount) / float64(deliverableCount) * 100
	}

	// Calculate average delivery time
	var avgDeliveryTime string
	if len(totalDeliveryDurations) > 0 {
		var totalDuration time.Duration
		for _, d := range totalDeliveryDurations {
			totalDuration += d
		}
		avgDuration := totalDuration / time.Duration(len(totalDeliveryDurations))
		avgDeliveryTime = formatDuration(avgDuration)
	}

	// Ensure all expected statuses are in breakdown (even if 0)
	allStatuses := []string{
		models.DeliveryStatusPending,
		models.DeliveryStatusScheduled,
		models.DeliveryStatusQueued,
		models.DeliveryStatusSending,
		models.DeliveryStatusRetrying,
		models.DeliveryStatusSent,
		models.DeliveryStatusDelivered,
		models.DeliveryStatusRead,
		models.DeliveryStatusFailed,
		models.DeliveryStatusExpired,
	}
	for _, status := range allStatuses {
		if _, exists := breakdown[status]; !exists {
			breakdown[status] = 0
		}
	}

	// Sort breakdown keys for consistent output
	sortedStatuses := make([]string, 0, len(breakdown))
	for status := range breakdown {
		sortedStatuses = append(sortedStatuses, status)
	}
	slices.Sort(sortedStatuses)
	sortedBreakdown := make(map[string]int, len(breakdown))
	for _, status := range sortedStatuses {
		sortedBreakdown[status] = breakdown[status]
	}

	// Determine period start and end dates for response
	var periodStartDate, periodEndDate string
	if startDate.IsZero() {
		periodStartDate = ""
		periodEndDate = ""
	} else {
		periodStartDate = startDate.Format(time.RFC3339)
		periodEndDate = endDate.Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": DeliveryAnalyticsStats{
			TotalSent:       totalSent,
			SuccessRate:     successRate,
			FailedLast7Days: failedLast7Days,
			AvgDeliveryTime: avgDeliveryTime,
			Breakdown:       sortedBreakdown,
			Period:          period,
			PeriodStartDate: periodStartDate,
			PeriodEndDate:   periodEndDate,
		},
		"message": "success",
	})
}

// formatDuration formats a duration as a human readable string
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return "< 1s"
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// FailedDeliveryItem represents a single failed delivery in the list
type FailedDeliveryItem struct {
	ReminderID           string `json:"reminder_id"`
	PatientNameMasked    string `json:"patient_name_masked"`
	PhoneMasked          string `json:"phone_masked,omitempty"`
	VolunteerName        string `json:"volunteer_name"`
	ReminderTitle        string `json:"reminder_title"`
	FailureReason        string `json:"failure_reason"`
	FailureReasonCode    string `json:"failure_reason_code"`
	FailureTimestamp     string `json:"failure_timestamp"`
	RetryCount           int    `json:"retry_count"`
	DeliveryErrorMessage string `json:"delivery_error_message,omitempty"`
}

// FailedDeliveryPagination represents pagination data
type FailedDeliveryPagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// FailedDeliveryFilterCounts represents filter counts by reason
type FailedDeliveryFilterCounts struct {
	InvalidPhone    int `json:"invalid_phone"`
	GOWATimeout     int `json:"gowa_timeout"`
	MessageRejected int `json:"message_rejected"`
	Other           int `json:"other"`
}

// FailedDeliveriesResponse represents the failed deliveries list response
type FailedDeliveriesResponse struct {
	Items        []FailedDeliveryItem      `json:"items"`
	Pagination   FailedDeliveryPagination  `json:"pagination"`
	FilterCounts FailedDeliveryFilterCounts `json:"filter_counts"`
}

// categorizeFailureReason categorizes an error message into a reason code
func categorizeFailureReason(errorMsg string) (reasonCode, displayText string) {
	errorMsg = strings.ToLower(errorMsg)

	if strings.Contains(errorMsg, "invalid") || strings.Contains(errorMsg, "nomor") {
		return "invalid_phone", "Nomor tidak valid"
	}
	if strings.Contains(errorMsg, "timeout") || strings.Contains(errorMsg, "connection") {
		return "gowa_timeout", "GOWA timeout"
	}
	if strings.Contains(errorMsg, "rejected") || strings.Contains(errorMsg, "ditolak") {
		return "message_rejected", "Pesan ditolak"
	}
	return "other", "Lainnya"
}

// GetFailedDeliveries returns a list of failed deliveries for admin review
func (h *AnalyticsHandler) GetFailedDeliveries(c *gin.Context) {
	// Check admin role
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	// Parse filter parameter
	filterReason := c.Query("reason")

	// Collect failed deliveries and count by category
	var failedDeliveries []FailedDeliveryItem
	filterCounts := FailedDeliveryFilterCounts{}

	h.patientStore.RLock()
	for _, patient := range h.patientStore.Patients {
		for _, reminder := range patient.Reminders {
			// Only include failed deliveries
			if reminder.DeliveryStatus != models.DeliveryStatusFailed {
				continue
			}

			// Categorize the failure reason
			reasonCode, _ := categorizeFailureReason(reminder.DeliveryErrorMessage)

			// Count for filter totals
			switch reasonCode {
			case "invalid_phone":
				filterCounts.InvalidPhone++
			case "gowa_timeout":
				filterCounts.GOWATimeout++
			case "message_rejected":
				filterCounts.MessageRejected++
			default:
				filterCounts.Other++
			}

			// Apply filter if specified
			if filterReason != "" && reasonCode != filterReason {
				continue
			}

			failedDeliveries = append(failedDeliveries, FailedDeliveryItem{
				ReminderID:           reminder.ID,
				PatientNameMasked:    utils.MaskPatientName(patient.Name),
				PhoneMasked:          utils.MaskPhoneNumber(patient.Phone),
				VolunteerName:        patient.CreatedBy,
				ReminderTitle:        reminder.Title,
				FailureReason:        reminder.DeliveryErrorMessage,
				FailureReasonCode:    reasonCode,
				FailureTimestamp:     reminder.MessageSentAt,
				RetryCount:           reminder.RetryCount,
				DeliveryErrorMessage: reminder.DeliveryErrorMessage,
			})
		}
	}
	h.patientStore.RUnlock()

	// Sort by failure timestamp (newest first)
	slices.SortFunc(failedDeliveries, func(a, b FailedDeliveryItem) int {
		if a.FailureTimestamp > b.FailureTimestamp {
			return -1
		}
		if a.FailureTimestamp < b.FailureTimestamp {
			return 1
		}
		return 0
	})

	// Calculate pagination
	total := len(failedDeliveries)
	totalPages := (total + limit - 1) / limit
	startIdx := (page - 1) * limit
	endIdx := startIdx + limit
	if endIdx > total {
		endIdx = total
	}

	// Get paginated slice
	var paginatedItems []FailedDeliveryItem
	if startIdx < total {
		if endIdx > total {
			endIdx = total
		}
		paginatedItems = failedDeliveries[startIdx:endIdx]
	}

	c.JSON(http.StatusOK, gin.H{
		"data": FailedDeliveriesResponse{
			Items:        paginatedItems,
			Pagination: FailedDeliveryPagination{
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: totalPages,
			},
			FilterCounts: filterCounts,
		},
		"message": "Failed deliveries retrieved successfully",
	})
}

// ExportFailedDeliveries exports failed deliveries as CSV
func (h *AnalyticsHandler) ExportFailedDeliveries(c *gin.Context) {
	// Check admin role
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	// Parse filter parameter
	filterReason := c.Query("reason")

	// Collect failed deliveries
	var failedDeliveries []FailedDeliveryItem

	h.patientStore.RLock()
	for _, patient := range h.patientStore.Patients {
		for _, reminder := range patient.Reminders {
			if reminder.DeliveryStatus != models.DeliveryStatusFailed {
				continue
			}

			reasonCode, _ := categorizeFailureReason(reminder.DeliveryErrorMessage)

			if filterReason != "" && reasonCode != filterReason {
				continue
			}

			failedDeliveries = append(failedDeliveries, FailedDeliveryItem{
				ReminderID:           reminder.ID,
				PatientNameMasked:    utils.MaskPatientName(patient.Name),
				PhoneMasked:          utils.MaskPhoneNumber(patient.Phone),
				VolunteerName:        patient.CreatedBy,
				ReminderTitle:        reminder.Title,
				FailureReason:        reminder.DeliveryErrorMessage,
				FailureReasonCode:    reasonCode,
				FailureTimestamp:     reminder.MessageSentAt,
				RetryCount:           reminder.RetryCount,
				DeliveryErrorMessage: reminder.DeliveryErrorMessage,
			})
		}
	}
	h.patientStore.RUnlock()

	// Sort by failure timestamp (newest first)
	slices.SortFunc(failedDeliveries, func(a, b FailedDeliveryItem) int {
		if a.FailureTimestamp > b.FailureTimestamp {
			return -1
		}
		if a.FailureTimestamp < b.FailureTimestamp {
			return 1
		}
		return 0
	})

	// Generate CSV
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=failed-deliveries-"+time.Now().Format("2006-01-02")+".csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	writer.Write([]string{
		"Reminder ID",
		"Patient Name (Masked)",
		"Volunteer",
		"Reminder Title",
		"Failure Reason",
		"Failure Timestamp",
		"Retry Count",
		"Error Message",
	})

	// Write data rows
	for _, item := range failedDeliveries {
		writer.Write([]string{
			item.ReminderID,
			item.PatientNameMasked,
			item.VolunteerName,
			item.ReminderTitle,
			item.FailureReason,
			item.FailureTimestamp,
			strconv.Itoa(item.RetryCount),
			item.DeliveryErrorMessage,
		})
	}
}

// GetFailedDeliveryDetail returns detailed information about a single failed delivery
func (h *AnalyticsHandler) GetFailedDeliveryDetail(c *gin.Context) {
	// Check admin role
	role := c.GetString("role")
	if role != "admin" && role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
		return
	}

	reminderID := c.Param("id")
	if reminderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reminder ID is required"})
		return
	}

	var foundReminder *models.Reminder
	var foundPatient *models.Patient

	h.patientStore.RLock()
	for _, patient := range h.patientStore.Patients {
		for _, reminder := range patient.Reminders {
			if reminder.ID == reminderID {
				if reminder.DeliveryStatus == models.DeliveryStatusFailed {
					foundReminder = reminder
					foundPatient = patient
					break
				}
			}
		}
		if foundReminder != nil {
			break
		}
	}
	h.patientStore.RUnlock()

	if foundReminder == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "failed delivery not found"})
		return
	}

	reasonCode, _ := categorizeFailureReason(foundReminder.DeliveryErrorMessage)

	// Build response (note: retry_attempts is empty because model doesn't store individual retry attempts)
	response := gin.H{
		"data": gin.H{
			"reminder_id":           foundReminder.ID,
			"patient_name_masked":   utils.MaskPatientName(foundPatient.Name),
			"phone_masked":          utils.MaskPhoneNumber(foundPatient.Phone),
			"volunteer_name":        foundPatient.CreatedBy,
			"reminder_title":        foundReminder.Title,
			"delivery_error_message": foundReminder.DeliveryErrorMessage,
			"failure_reason":        foundReminder.DeliveryErrorMessage,
			"failure_reason_code":   reasonCode,
			"failure_timestamp":     foundReminder.MessageSentAt,
			"retry_count":           foundReminder.RetryCount,
			"retry_attempts":        []interface{}{}, // Empty - model doesn't store individual attempts
		},
		"message": "Failed delivery details retrieved",
	}

	c.JSON(http.StatusOK, response)
}
