import { sseService } from '$lib/services/sse.js';
import { toastStore } from '$lib/stores/toast.svelte.js';
import { get } from 'svelte/store';
import { _ } from 'svelte-i18n';

/**
 * Delivery status store using Svelte 5 runes
 * Manages real-time delivery status updates via SSE
 *
 * CRITICAL: This is Vite + Svelte 5, NOT SvelteKit!
 * - Uses Svelte 5 runes ($state, $derived)
 * - No legacy reactive statements ($:)
 * - No SvelteKit imports
 */
class DeliveryStore {
    deliveryStatuses = $state({});
    connectionStatus = $state('disconnected');
    failedReminders = $state([]); // Array of failed reminder objects

    // Derived reactive count
    failedCount = $derived(this.failedReminders.length);

    constructor() {
        // Subscribe to SSE events
        sseService.on('delivery.status.updated', (data) => {
            this.updateStatus(data.reminder_id, data.status, data.timestamp);
        });

        sseService.on('connection.status', (status) => {
            this.connectionStatus = status;
        });

        // Handle delivery.failed event
        sseService.on('delivery.failed', (data) => {
            // Add to failed reminders list
            this.addFailedReminder(data);

            // Show toast notification with i18n
            const t = get(_);
            toastStore.add(
                t('reminder.failedNotification', { values: { patientName: data.patient_name } }),
                {
                    type: 'error',
                    action: {
                        label: t('reminder.viewDetails'),
                        onClick: () => {
                            // Navigate to patient detail by scrolling to patient
                            // Dispatch custom event that PatientsView will handle
                            window.dispatchEvent(new CustomEvent('navigate-to-patient', {
                                detail: { patientId: data.patient_id }
                            }));
                        }
                    },
                    duration: 5000
                }
            );
        });
    }

    /**
     * Update delivery status for a reminder
     * CRITICAL: Create new object reference to trigger Svelte 5 reactivity
     */
    updateStatus(reminderId, status, timestamp) {
        this.deliveryStatuses = {
            ...this.deliveryStatuses,
            [reminderId]: {
                status,
                timestamp,
                updatedAt: new Date().toISOString(),
            }
        };
    }

    /**
     * Get delivery status for a reminder
     */
    getStatus(reminderId) {
        return this.deliveryStatuses[reminderId]?.status || null;
    }

    /**
     * Initialize SSE connection
     */
    connect() {
        sseService.connect();
    }

    /**
     * Disconnect SSE connection
     */
    disconnect() {
        sseService.disconnect();
    }

    /**
     * Add a failed reminder to the list
     * CRITICAL: Create new array reference to trigger Svelte 5 reactivity
     */
    addFailedReminder(data) {
        this.failedReminders = [
            ...this.failedReminders,
            {
                reminderId: data.reminder_id,
                patientId: data.patient_id,
                patientName: data.patient_name,
                error: data.error,
                timestamp: data.timestamp || new Date().toISOString()
            }
        ];
    }

    /**
     * Remove a failed reminder from the list
     * CRITICAL: Create new array reference to trigger Svelte 5 reactivity
     */
    removeFailedReminder(reminderId) {
        this.failedReminders = this.failedReminders.filter(
            r => r.reminderId !== reminderId
        );
    }

    /**
     * Clear all failed reminders
     */
    clearFailedReminders() {
        this.failedReminders = [];
    }

    /**
     * Get all failed reminders
     */
    getFailedReminders() {
        return this.failedReminders;
    }
}

export const deliveryStore = new DeliveryStore();
