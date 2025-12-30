import { describe, it, expect, beforeEach, vi } from 'vitest';
import { deliveryStore } from './delivery.svelte.js';
import { sseService } from '$lib/services/sse.js';

// Mock SSE service
vi.mock('$lib/services/sse.js', () => ({
  sseService: {
    on: vi.fn(),
    connect: vi.fn(),
    disconnect: vi.fn()
  }
}));

// Mock toast store
vi.mock('./toast.svelte.js', () => ({
  toastStore: {
    add: vi.fn()
  }
}));

describe('DeliveryStore', () => {
  beforeEach(() => {
    // Clear all mocks
    vi.clearAllMocks();

    // Reset store state
    deliveryStore.deliveryStatuses = {};
    deliveryStore.failedReminders = [];
    deliveryStore.connectionStatus = 'disconnected';
  });

  it('should initialize with empty state', () => {
    expect(deliveryStore.deliveryStatuses).toEqual({});
    expect(deliveryStore.failedReminders).toEqual([]);
    expect(deliveryStore.failedCount).toBe(0);
    expect(deliveryStore.connectionStatus).toBe('disconnected');
  });

  it('should update delivery status correctly', () => {
    const reminderId = 'reminder-123';
    const status = 'sent';
    const timestamp = '2025-12-30T10:00:00Z';

    deliveryStore.updateStatus(reminderId, status, timestamp);

    expect(deliveryStore.deliveryStatuses[reminderId]).toMatchObject({
      status,
      timestamp
    });
    expect(deliveryStore.deliveryStatuses[reminderId].updatedAt).toBeDefined();
  });

  it('should add failed reminder correctly', () => {
    const failedData = {
      reminder_id: 'reminder-123',
      patient_id: 'patient-456',
      patient_name: 'John Doe',
      error: 'Nomor tidak valid',
      timestamp: '2025-12-30T10:00:00Z'
    };

    deliveryStore.addFailedReminder(failedData);

    expect(deliveryStore.failedReminders).toHaveLength(1);
    expect(deliveryStore.failedReminders[0]).toMatchObject({
      reminderId: 'reminder-123',
      patientId: 'patient-456',
      patientName: 'John Doe',
      error: 'Nomor tidak valid',
      timestamp: '2025-12-30T10:00:00Z'
    });
  });

  it('should update failedCount reactively when adding failed reminders', () => {
    expect(deliveryStore.failedCount).toBe(0);

    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1',
      timestamp: '2025-12-30T10:00:00Z'
    });

    expect(deliveryStore.failedCount).toBe(1);

    deliveryStore.addFailedReminder({
      reminder_id: 'r2',
      patient_id: 'p2',
      patient_name: 'Patient 2',
      error: 'Error 2',
      timestamp: '2025-12-30T10:01:00Z'
    });

    expect(deliveryStore.failedCount).toBe(2);
  });

  it('should remove failed reminder correctly', () => {
    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    deliveryStore.addFailedReminder({
      reminder_id: 'r2',
      patient_id: 'p2',
      patient_name: 'Patient 2',
      error: 'Error 2'
    });

    expect(deliveryStore.failedCount).toBe(2);

    deliveryStore.removeFailedReminder('r1');

    expect(deliveryStore.failedCount).toBe(1);
    expect(deliveryStore.failedReminders[0].reminderId).toBe('r2');
  });

  it('should clear all failed reminders', () => {
    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    deliveryStore.addFailedReminder({
      reminder_id: 'r2',
      patient_id: 'p2',
      patient_name: 'Patient 2',
      error: 'Error 2'
    });

    expect(deliveryStore.failedCount).toBe(2);

    deliveryStore.clearFailedReminders();

    expect(deliveryStore.failedCount).toBe(0);
    expect(deliveryStore.failedReminders).toEqual([]);
  });

  it('should get failed reminders list', () => {
    const failed1 = {
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    };

    const failed2 = {
      reminder_id: 'r2',
      patient_id: 'p2',
      patient_name: 'Patient 2',
      error: 'Error 2'
    };

    deliveryStore.addFailedReminder(failed1);
    deliveryStore.addFailedReminder(failed2);

    const failedList = deliveryStore.getFailedReminders();

    expect(failedList).toHaveLength(2);
    expect(failedList[0].reminderId).toBe('r1');
    expect(failedList[1].reminderId).toBe('r2');
  });

  it('should get delivery status for a reminder', () => {
    deliveryStore.updateStatus('reminder-123', 'sent', '2025-12-30T10:00:00Z');

    const status = deliveryStore.getStatus('reminder-123');

    expect(status).toBe('sent');
  });

  it('should return null for non-existent reminder status', () => {
    const status = deliveryStore.getStatus('non-existent');

    expect(status).toBeNull();
  });

  it('should create new object reference when updating status (Svelte 5 reactivity)', () => {
    const originalRef = deliveryStore.deliveryStatuses;

    deliveryStore.updateStatus('reminder-123', 'sent', '2025-12-30T10:00:00Z');

    // Reference should change
    expect(deliveryStore.deliveryStatuses).not.toBe(originalRef);
  });

  it('should create new array reference when adding failed reminder (Svelte 5 reactivity)', () => {
    const originalRef = deliveryStore.failedReminders;

    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    // Reference should change
    expect(deliveryStore.failedReminders).not.toBe(originalRef);
  });

  it('should create new array reference when removing failed reminder (Svelte 5 reactivity)', () => {
    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    const originalRef = deliveryStore.failedReminders;

    deliveryStore.removeFailedReminder('r1');

    // Reference should change
    expect(deliveryStore.failedReminders).not.toBe(originalRef);
  });
});
