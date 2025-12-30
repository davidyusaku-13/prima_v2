/**
 * SSE Client Service for real-time delivery status updates
 * Uses native browser EventSource API (NOT SvelteKit-specific)
 *
 * CRITICAL: This is Vite + Svelte 5, NOT SvelteKit!
 * - Uses native EventSource (browser API)
 * - No SvelteKit imports ($app/*)
 * - Token passed via query parameter (EventSource doesn't support headers)
 */
class SSEService {
    constructor() {
        this.eventSource = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 10;
        this.reconnectDelay = 1000; // Start with 1 second
        this.listeners = new Map();
        this.connectionStatus = 'disconnected'; // 'connected' | 'connecting' | 'disconnected'
    }

    /**
     * Connect to SSE endpoint with JWT authentication
     */
    connect() {
        if (this.eventSource) {
            return; // Already connected
        }

        // Get token from localStorage (where App.svelte stores it)
        const token = localStorage.getItem('token');
        if (!token) {
            console.error('No JWT token available for SSE connection');
            return;
        }

        this.connectionStatus = 'connecting';
        this.notifyStatusChange();

        // Pass JWT token via query parameter (EventSource doesn't support headers)
        const url = `http://localhost:8080/api/sse/delivery-status?token=${token}`;

        this.eventSource = new EventSource(url);

        // Connection established
        this.eventSource.addEventListener('connection.established', (e) => {
            console.log('SSE connection established:', e.data);
            this.connectionStatus = 'connected';
            this.reconnectAttempts = 0;
            this.reconnectDelay = 1000;
            this.notifyStatusChange();
        });

        // Delivery status updated
        this.eventSource.addEventListener('delivery.status.updated', (e) => {
            const data = JSON.parse(e.data);
            this.notifyListeners('delivery.status.updated', data);
        });

        // Connection error
        this.eventSource.onerror = (error) => {
            console.error('SSE connection error:', error);
            this.connectionStatus = 'disconnected';
            this.notifyStatusChange();

            // Close and attempt reconnect
            this.eventSource.close();
            this.eventSource = null;
            this.attemptReconnect();
        };
    }

    /**
     * Attempt to reconnect with exponential backoff
     */
    attemptReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('Max SSE reconnection attempts reached');
            return;
        }

        this.reconnectAttempts++;
        const delay = Math.min(this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1), 30000);

        console.log(`Attempting SSE reconnection in ${delay}ms (attempt ${this.reconnectAttempts})`);

        setTimeout(() => {
            this.connect();
        }, delay);
    }

    /**
     * Subscribe to SSE events
     */
    on(eventType, callback) {
        if (!this.listeners.has(eventType)) {
            this.listeners.set(eventType, []);
        }
        this.listeners.get(eventType).push(callback);
    }

    /**
     * Unsubscribe from SSE events
     */
    off(eventType, callback) {
        if (!this.listeners.has(eventType)) return;

        const callbacks = this.listeners.get(eventType);
        const index = callbacks.indexOf(callback);
        if (index > -1) {
            callbacks.splice(index, 1);
        }
    }

    /**
     * Notify all listeners of an event
     */
    notifyListeners(eventType, data) {
        if (!this.listeners.has(eventType)) return;

        this.listeners.get(eventType).forEach(callback => {
            try {
                callback(data);
            } catch (error) {
                console.error('Error in SSE listener:', error);
            }
        });
    }

    /**
     * Notify status change listeners
     */
    notifyStatusChange() {
        this.notifyListeners('connection.status', this.connectionStatus);
    }

    /**
     * Get current connection status
     */
    getStatus() {
        return this.connectionStatus;
    }

    /**
     * Disconnect from SSE endpoint
     */
    disconnect() {
        if (this.eventSource) {
            this.eventSource.close();
            this.eventSource = null;
        }
        this.connectionStatus = 'disconnected';
        this.notifyStatusChange();
    }
}

// Singleton instance
export const sseService = new SSEService();
