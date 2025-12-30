/**
 * Toast notification store using Svelte 5 runes
 * Manages toast notifications with auto-dismiss and action buttons
 *
 * CRITICAL: This is Vite + Svelte 5, NOT SvelteKit!
 * - Uses Svelte 5 runes ($state)
 * - No legacy reactive statements ($:)
 * - No SvelteKit imports
 */
class ToastStore {
    toasts = $state([]);
    nextId = 0;
    maxToasts = 5; // Limit max toasts displayed

    /**
     * Add a new toast notification
     * @param {string} message - The message to display
     * @param {Object} options - Toast options
     * @param {string} options.type - Toast type: 'success' | 'error' | 'warning' | 'info'
     * @param {Object} options.action - Action button: { label: string, onClick: function }
     * @param {number} options.duration - Auto-dismiss duration in ms (0 = no auto-dismiss)
     * @returns {number} Toast ID
     */
    add(message, options = {}) {
        const id = this.nextId++;
        const toast = {
            id,
            message,
            type: options.type || 'info',
            action: options.action || null,
            duration: options.duration !== undefined ? options.duration : 5000
        };

        // Remove oldest toast if limit reached
        if (this.toasts.length >= this.maxToasts) {
            // CRITICAL: Create new array reference for Svelte 5 reactivity
            this.toasts = this.toasts.slice(1);
        }

        // CRITICAL: Create new array reference to trigger Svelte 5 reactivity
        this.toasts = [...this.toasts, toast];

        // Auto-dismiss
        if (toast.duration > 0) {
            setTimeout(() => {
                this.remove(id);
            }, toast.duration);
        }

        return id;
    }

    /**
     * Remove a toast by ID
     * @param {number} id - Toast ID to remove
     */
    remove(id) {
        // CRITICAL: Create new array reference to trigger Svelte 5 reactivity
        this.toasts = this.toasts.filter(t => t.id !== id);
    }

    /**
     * Clear all toasts
     */
    clear() {
        this.toasts = [];
    }
}

export const toastStore = new ToastStore();
