import { describe, it, expect, beforeEach, vi } from 'vitest';
import { toastStore } from './toast.svelte.js';

describe('ToastStore', () => {
  beforeEach(() => {
    // Clear toasts before each test
    toastStore.clear();
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('should add a toast with default options', () => {
    const id = toastStore.add('Test message');

    expect(toastStore.toasts).toHaveLength(1);
    expect(toastStore.toasts[0]).toMatchObject({
      id,
      message: 'Test message',
      type: 'info',
      action: null,
      duration: 5000
    });
  });

  it('should add a toast with custom options', () => {
    const action = { label: 'Click me', onClick: vi.fn() };
    const id = toastStore.add('Error message', {
      type: 'error',
      action,
      duration: 3000
    });

    expect(toastStore.toasts).toHaveLength(1);
    expect(toastStore.toasts[0]).toMatchObject({
      id,
      message: 'Error message',
      type: 'error',
      action,
      duration: 3000
    });
  });

  it('should auto-dismiss toast after duration', () => {
    toastStore.add('Test message', { duration: 5000 });

    expect(toastStore.toasts).toHaveLength(1);

    // Fast-forward time by 5000ms
    vi.advanceTimersByTime(5000);

    expect(toastStore.toasts).toHaveLength(0);
  });

  it('should not auto-dismiss toast with duration 0', () => {
    toastStore.add('Persistent message', { duration: 0 });

    expect(toastStore.toasts).toHaveLength(1);

    // Fast-forward time
    vi.advanceTimersByTime(10000);

    // Toast should still be there
    expect(toastStore.toasts).toHaveLength(1);
  });

  it('should remove toast by id', () => {
    const id1 = toastStore.add('Message 1');
    const id2 = toastStore.add('Message 2');

    expect(toastStore.toasts).toHaveLength(2);

    toastStore.remove(id1);

    expect(toastStore.toasts).toHaveLength(1);
    expect(toastStore.toasts[0].id).toBe(id2);
  });

  it('should clear all toasts', () => {
    toastStore.add('Message 1');
    toastStore.add('Message 2');
    toastStore.add('Message 3');

    expect(toastStore.toasts).toHaveLength(3);

    toastStore.clear();

    expect(toastStore.toasts).toHaveLength(0);
  });

  it('should limit max toasts to 5', () => {
    // Add 6 toasts
    for (let i = 0; i < 6; i++) {
      toastStore.add(`Message ${i}`);
    }

    // Should only have 5 toasts (oldest removed)
    expect(toastStore.toasts).toHaveLength(5);
    expect(toastStore.toasts[0].message).toBe('Message 1');
    expect(toastStore.toasts[4].message).toBe('Message 5');
  });

  it('should generate unique IDs for each toast', () => {
    const id1 = toastStore.add('Message 1');
    const id2 = toastStore.add('Message 2');
    const id3 = toastStore.add('Message 3');

    expect(id1).not.toBe(id2);
    expect(id2).not.toBe(id3);
    expect(id1).not.toBe(id3);
  });

  it('should handle multiple toasts with different durations', () => {
    toastStore.add('Fast toast', { duration: 1000 });
    toastStore.add('Slow toast', { duration: 5000 });

    expect(toastStore.toasts).toHaveLength(2);

    // Fast-forward by 1000ms
    vi.advanceTimersByTime(1000);

    // First toast should be removed
    expect(toastStore.toasts).toHaveLength(1);
    expect(toastStore.toasts[0].message).toBe('Slow toast');

    // Fast-forward by another 4000ms
    vi.advanceTimersByTime(4000);

    // Second toast should be removed
    expect(toastStore.toasts).toHaveLength(0);
  });
});
