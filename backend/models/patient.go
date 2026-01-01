package models

import (
	"sync"
)

// DeliveryStatus constants for reminder delivery tracking
// State machine transitions:
//   pending → scheduled (quiet hours) → sending → sent → delivered → read
//   pending → queued → sending → sent → delivered → read
//   sending → failed
//   scheduled → sending (at scheduled time)
//   sending → retrying → sending (on transient failure)
//   retrying → sent (on success)
//   retrying → failed (after max retries exhausted)
//   any → cancelled (user cancelled the reminder)
const (
	DeliveryStatusPending   = "pending"
	DeliveryStatusScheduled = "scheduled" // Queued for quiet hours delivery
	DeliveryStatusQueued    = "queued"    // Queued due to circuit breaker open
	DeliveryStatusSending   = "sending"
	DeliveryStatusRetrying  = "retrying"  // Waiting for retry after transient failure
	DeliveryStatusSent      = "sent"
	DeliveryStatusDelivered = "delivered"
	DeliveryStatusRead      = "read"
	DeliveryStatusFailed    = "failed"
	DeliveryStatusExpired   = "expired"
	DeliveryStatusCancelled = "cancelled" // Reminder was cancelled by user
)

// Recurrence represents reminder recurrence settings
type Recurrence struct {
	Frequency  string `json:"frequency"`
	Interval   int    `json:"interval"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	EndDate    string `json:"endDate,omitempty"`
}

// Attachment represents content attached to a reminder
type Attachment struct {
	Type  string `json:"type"`  // "article" or "video"
	ID    string `json:"id"`
	Title string `json:"title"` // Cached title for display in WhatsApp message
	URL   string `json:"url"`   // URL to the content (article slug or video watch page)
}

// Reminder represents a patient reminder with delivery tracking
type Reminder struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     string     `json:"dueDate,omitempty"`
	Priority    string     `json:"priority"`
	Completed   bool       `json:"completed"`
	Recurrence  Recurrence `json:"recurrence"`
	Notified    bool       `json:"notified"`

	// Content attachments
	Attachments []Attachment `json:"attachments,omitempty"`

	// Delivery tracking fields (verbose names per architecture)
	GOWAMessageID        string `json:"gowa_message_id,omitempty"`
	DeliveryStatus       string `json:"delivery_status,omitempty"`
	DeliveryErrorMessage string `json:"delivery_error_message,omitempty"`
	MessageSentAt        string `json:"message_sent_at,omitempty"`        // ISO 8601 UTC
	DeliveredAt          string `json:"delivered_at,omitempty"`           // ISO 8601 UTC
	ReadAt               string `json:"read_at,omitempty"`                // ISO 8601 UTC
	RetryCount           int    `json:"retry_count,omitempty"`            // Number of retry attempts
	ScheduledDeliveryAt  string `json:"scheduled_delivery_at,omitempty"` // ISO 8601 UTC - for quiet hours scheduling
	CancelledAt          string `json:"cancelled_at,omitempty"`           // ISO 8601 UTC - when reminder was cancelled
	CancelledBy          string `json:"cancelled_by,omitempty"`           // User ID who cancelled the reminder
}

// Patient represents a patient record
type Patient struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Phone     string      `json:"phone"`
	Email     string      `json:"email,omitempty"`
	Notes     string      `json:"notes,omitempty"`
	Reminders []*Reminder `json:"reminders,omitempty"`
	CreatedBy string      `json:"createdBy,omitempty"`
}

// PatientStore handles patient data persistence with thread-safe operations
type PatientStore struct {
	Mu       sync.RWMutex
	Patients map[string]*Patient
	SaveFunc func()
}

// NewPatientStore creates a new patient store
func NewPatientStore(saveFunc func()) *PatientStore {
	return &PatientStore{
		Patients: make(map[string]*Patient),
		SaveFunc: saveFunc,
	}
}

// GetPatient retrieves a patient by ID
func (s *PatientStore) GetPatient(id string) (*Patient, bool) {
	p, ok := s.Patients[id]
	return p, ok
}

// SaveData triggers the save function
func (s *PatientStore) SaveData() {
	if s.SaveFunc != nil {
		s.SaveFunc()
	}
}

// Lock acquires write lock
func (s *PatientStore) Lock() { s.Mu.Lock() }

// Unlock releases write lock
func (s *PatientStore) Unlock() { s.Mu.Unlock() }

// RLock acquires read lock
func (s *PatientStore) RLock() { s.Mu.RLock() }

// RUnlock releases read lock
func (s *PatientStore) RUnlock() { s.Mu.RUnlock() }
