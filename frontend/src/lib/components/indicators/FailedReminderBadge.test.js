import { describe, it, expect, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import { deliveryStore } from '$lib/stores/delivery.svelte.js';
import FailedReminderBadge from './FailedReminderBadge.svelte';

// Mock i18n
vi.mock('svelte-i18n', () => ({
  _: {
    subscribe: (fn) => {
      fn((key, options) => {
        if (key === 'delivery.failedRemindersCount') {
          return `${options.values.count} failed reminders`;
        }
        return key;
      });
      return () => {};
    }
  }
}));

describe('FailedReminderBadge', () => {
  beforeEach(() => {
    // Reset store
    deliveryStore.failedReminders = [];
  });

  it('should not render when failedCount is 0', () => {
    const { container } = render(FailedReminderBadge);

    expect(container.querySelector('button')).toBeNull();
  });

  it('should render badge when failedCount > 0', () => {
    // Add a failed reminder
    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    const { container } = render(FailedReminderBadge);

    const button = container.querySelector('button');
    expect(button).toBeTruthy();
  });

  it('should display correct count', () => {
    // Add 3 failed reminders
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
    deliveryStore.addFailedReminder({
      reminder_id: 'r3',
      patient_id: 'p3',
      patient_name: 'Patient 3',
      error: 'Error 3'
    });

    const { container } = render(FailedReminderBadge);

    const countSpan = container.querySelector('span');
    expect(countSpan?.textContent).toBe('3');
  });

  it('should have red background styling', () => {
    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    const { container } = render(FailedReminderBadge);

    const button = container.querySelector('button');
    expect(button?.className).toContain('bg-red-600');
  });

  it('should have alert icon', () => {
    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    const { container } = render(FailedReminderBadge);

    const svg = container.querySelector('svg');
    expect(svg).toBeTruthy();
  });

  it('should dispatch show-failed-reminders event on click', () => {
    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    const eventSpy = vi.fn();
    window.addEventListener('show-failed-reminders', eventSpy);

    const { container } = render(FailedReminderBadge);

    const button = container.querySelector('button');
    button?.click();

    expect(eventSpy).toHaveBeenCalled();

    window.removeEventListener('show-failed-reminders', eventSpy);
  });

  it('should update count reactively when failed reminders change', () => {
    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    const { container } = render(FailedReminderBadge);

    let countSpan = container.querySelector('span');
    expect(countSpan?.textContent).toBe('1');

    // Add another failed reminder
    deliveryStore.addFailedReminder({
      reminder_id: 'r2',
      patient_id: 'p2',
      patient_name: 'Patient 2',
      error: 'Error 2'
    });

    // Count should update
    countSpan = container.querySelector('span');
    expect(countSpan?.textContent).toBe('2');
  });

  it('should hide when all failed reminders are cleared', () => {
    deliveryStore.addFailedReminder({
      reminder_id: 'r1',
      patient_id: 'p1',
      patient_name: 'Patient 1',
      error: 'Error 1'
    });

    const { container } = render(FailedReminderBadge);

    let button = container.querySelector('button');
    expect(button).toBeTruthy();

    // Clear failed reminders
    deliveryStore.clearFailedReminders();

    // Badge should be hidden
    button = container.querySelector('button');
    expect(button).toBeNull();
  });
});
