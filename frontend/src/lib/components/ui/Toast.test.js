import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import { toastStore } from '$lib/stores/toast.svelte.js';
import Toast from './Toast.svelte';

// Mock i18n
vi.mock('svelte-i18n', () => ({
  _: {
    subscribe: (fn) => {
      fn((key) => {
        if (key === 'common.close') return 'Close';
        return key;
      });
      return () => {};
    }
  }
}));

describe('Toast Component', () => {
  beforeEach(() => {
    toastStore.clear();
  });

  it('should render no toasts when store is empty', () => {
    const { container } = render(Toast);

    const toastElements = container.querySelectorAll('[role="alert"]');
    expect(toastElements.length).toBe(0);
  });

  it('should render toast when added to store', () => {
    toastStore.add('Test message');

    const { container } = render(Toast);

    const toastElements = container.querySelectorAll('[role="alert"]');
    expect(toastElements.length).toBe(1);
    expect(container.textContent).toContain('Test message');
  });

  it('should render multiple toasts', () => {
    toastStore.add('Message 1');
    toastStore.add('Message 2');
    toastStore.add('Message 3');

    const { container } = render(Toast);

    const toastElements = container.querySelectorAll('[role="alert"]');
    expect(toastElements.length).toBe(3);
  });

  it('should render error toast with red border', () => {
    toastStore.add('Error message', { type: 'error' });

    const { container } = render(Toast);

    const toast = container.querySelector('[role="alert"]');
    expect(toast?.className).toContain('border-red-500');
  });

  it('should render success toast with green border', () => {
    toastStore.add('Success message', { type: 'success' });

    const { container } = render(Toast);

    const toast = container.querySelector('[role="alert"]');
    expect(toast?.className).toContain('border-green-500');
  });

  it('should render warning toast with yellow border', () => {
    toastStore.add('Warning message', { type: 'warning' });

    const { container } = render(Toast);

    const toast = container.querySelector('[role="alert"]');
    expect(toast?.className).toContain('border-yellow-500');
  });

  it('should render info toast with blue border', () => {
    toastStore.add('Info message', { type: 'info' });

    const { container } = render(Toast);

    const toast = container.querySelector('[role="alert"]');
    expect(toast?.className).toContain('border-blue-500');
  });

  it('should render action button when provided', () => {
    const action = { label: 'Click me', onClick: vi.fn() };
    toastStore.add('Message with action', { action });

    const { container } = render(Toast);

    const actionButton = Array.from(container.querySelectorAll('button'))
      .find(btn => btn.textContent === 'Click me');
    expect(actionButton).toBeTruthy();
  });

  it('should not render action button when not provided', () => {
    toastStore.add('Message without action');

    const { container } = render(Toast);

    // Should only have close button
    const buttons = container.querySelectorAll('button');
    expect(buttons.length).toBe(1);
  });

  it('should call action onClick when action button clicked', () => {
    const onClick = vi.fn();
    const action = { label: 'Click me', onClick };
    toastStore.add('Message with action', { action });

    const { container } = render(Toast);

    const actionButton = Array.from(container.querySelectorAll('button'))
      .find(btn => btn.textContent === 'Click me');
    actionButton?.click();

    expect(onClick).toHaveBeenCalled();
  });

  it('should remove toast when close button clicked', () => {
    toastStore.add('Test message');

    const { container } = render(Toast);

    let toastElements = container.querySelectorAll('[role="alert"]');
    expect(toastElements.length).toBe(1);

    // Click close button (last button)
    const buttons = container.querySelectorAll('button');
    const closeButton = buttons[buttons.length - 1];
    closeButton?.click();

    // Toast should be removed
    toastElements = container.querySelectorAll('[role="alert"]');
    expect(toastElements.length).toBe(0);
  });

  it('should remove toast when action button clicked', () => {
    const action = { label: 'Click me', onClick: vi.fn() };
    toastStore.add('Message with action', { action });

    const { container } = render(Toast);

    let toastElements = container.querySelectorAll('[role="alert"]');
    expect(toastElements.length).toBe(1);

    // Click action button
    const actionButton = Array.from(container.querySelectorAll('button'))
      .find(btn => btn.textContent === 'Click me');
    actionButton?.click();

    // Toast should be removed
    toastElements = container.querySelectorAll('[role="alert"]');
    expect(toastElements.length).toBe(0);
  });

  it('should render error icon for error type', () => {
    toastStore.add('Error message', { type: 'error' });

    const { container } = render(Toast);

    const svg = container.querySelector('svg.text-red-500');
    expect(svg).toBeTruthy();
  });

  it('should render success icon for success type', () => {
    toastStore.add('Success message', { type: 'success' });

    const { container } = render(Toast);

    const svg = container.querySelector('svg.text-green-500');
    expect(svg).toBeTruthy();
  });

  it('should render warning icon for warning type', () => {
    toastStore.add('Warning message', { type: 'warning' });

    const { container } = render(Toast);

    const svg = container.querySelector('svg.text-yellow-500');
    expect(svg).toBeTruthy();
  });

  it('should render info icon for info type', () => {
    toastStore.add('Info message', { type: 'info' });

    const { container } = render(Toast);

    const svg = container.querySelector('svg.text-blue-500');
    expect(svg).toBeTruthy();
  });

  it('should have aria-live="polite" for accessibility', () => {
    toastStore.add('Test message');

    const { container } = render(Toast);

    const toast = container.querySelector('[role="alert"]');
    expect(toast?.getAttribute('aria-live')).toBe('polite');
  });

  it('should be positioned at top-right', () => {
    const { container } = render(Toast);

    const toastContainer = container.querySelector('.fixed');
    expect(toastContainer?.className).toContain('top-4');
    expect(toastContainer?.className).toContain('right-4');
  });
});
