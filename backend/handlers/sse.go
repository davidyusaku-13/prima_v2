package handlers

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/gin-gonic/gin"
)

// SSEHandler handles Server-Sent Events for real-time delivery status updates
type SSEHandler struct {
	clients  map[chan SSEEvent]string // channel -> user_id
	mu       sync.RWMutex
	logger   *slog.Logger
	stopCh   chan struct{} // Signal to stop all SSE connections
	isClosed bool          // Track if handler is closed
}

// SSEEvent represents an event to be sent via SSE
type SSEEvent struct {
	Event string      `json:"-"`
	Data  interface{} `json:"data"`
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler(cfg *config.Config, logger *slog.Logger) *SSEHandler {
	return &SSEHandler{
		clients:  make(map[chan SSEEvent]string),
		logger:   logger,
		stopCh:   make(chan struct{}),
		isClosed: false,
	}
}

// HandleDeliveryStatusSSE handles SSE connections for delivery status updates
// GET /api/sse/delivery-status
func (h *SSEHandler) HandleDeliveryStatusSSE(c *gin.Context) {
	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no") // Disable nginx buffering

	// Get user from JWT (already authenticated by middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		// Send error event and close
		c.SSEvent("error", `{"error": "Unauthorized"}`)
		c.Writer.Flush()
		return
	}

	// Create client channel with buffer to prevent blocking
	clientChan := make(chan SSEEvent, 10)

	// Register client
	h.mu.Lock()
	h.clients[clientChan] = userID.(string)
	h.mu.Unlock()

	// Send initial connection event
	connectionData := map[string]string{
		"message":   "Connected to delivery status updates",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	dataJSON, _ := json.Marshal(connectionData)
	c.SSEvent("connection.established", string(dataJSON))
	c.Writer.Flush()

	// Use request context for disconnect detection (replaces deprecated CloseNotifier)
	ctx := c.Request.Context()

	// Listen for events or client disconnect
	for {
		select {
		case event := <-clientChan:
			// Send event to client
			data, err := json.Marshal(event.Data)
			if err != nil {
				if h.logger != nil {
					h.logger.Error("Failed to marshal SSE event data",
						"error", err.Error(),
						"event", event.Event,
					)
				}
				continue
			}
			c.SSEvent(event.Event, string(data))
			c.Writer.Flush()

		case <-ctx.Done():
			// Client disconnected
			h.mu.Lock()
			delete(h.clients, clientChan)
			h.mu.Unlock()
			close(clientChan)
			return

		case <-h.stopCh:
			// Server is shutting down, close connection gracefully
			h.mu.Lock()
			delete(h.clients, clientChan)
			h.mu.Unlock()
			close(clientChan)
			return
		}
	}
}

// BroadcastDeliveryStatusUpdate broadcasts delivery status update to all connected clients
func (h *SSEHandler) BroadcastDeliveryStatusUpdate(reminderID, status, timestamp string) {
	event := SSEEvent{
		Event: "delivery.status.updated",
		Data: map[string]string{
			"reminder_id": reminderID,
			"status":      status,
			"timestamp":   timestamp,
		},
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.clients) == 0 {
		// No clients connected, skip broadcasting
		return
	}

	// Broadcast to all connected clients
	for clientChan := range h.clients {
		select {
		case clientChan <- event:
			// Event sent successfully
		default:
			// Channel full, skip this client (client is slow)
			if h.logger != nil {
				h.logger.Warn("SSE client channel full, skipping event",
					"reminder_id", reminderID,
				)
			}
		}
	}
}

// BroadcastDeliveryFailed broadcasts delivery failure to all connected clients
func (h *SSEHandler) BroadcastDeliveryFailed(reminderID, patientID, patientName, errorMsg string) {
	event := SSEEvent{
		Event: "delivery.failed",
		Data: map[string]string{
			"reminder_id":  reminderID,
			"patient_id":   patientID,
			"patient_name": patientName,
			"error":        errorMsg,
			"timestamp":    time.Now().UTC().Format(time.RFC3339),
		},
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.clients) == 0 {
		// No clients connected, skip broadcasting
		return
	}

	// Broadcast to all connected clients
	for clientChan := range h.clients {
		select {
		case clientChan <- event:
			// Event sent successfully
		default:
			// Channel full, skip this client (client is slow)
			if h.logger != nil {
				h.logger.Warn("SSE client channel full, skipping delivery.failed event",
					"reminder_id", reminderID,
				)
			}
		}
	}
}

// GetClientCount returns the number of connected SSE clients
func (h *SSEHandler) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Shutdown gracefully closes all SSE connections
func (h *SSEHandler) Shutdown() {
	h.mu.Lock()
	if h.isClosed {
		h.mu.Unlock()
		return
	}
	h.isClosed = true
	h.mu.Unlock()

	// Signal all connections to close
	close(h.stopCh)

	// Wait a moment for connections to close gracefully
	time.Sleep(100 * time.Millisecond)
}
