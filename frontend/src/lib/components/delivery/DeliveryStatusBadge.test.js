/**
 * @vitest-environment happy-dom
 */
import { describe, it, expect, vi, afterEach, beforeEach } from 'vitest';
import { render, screen, cleanup, fireEvent } from '@testing-library/svelte';

// Mock svelte-i18n with proper store - define inline to avoid hoisting issues
vi.mock('svelte-i18n', async () => {
  const { readable } = await import('svelte/store');

  const translations = {
    'reminder.status.pending': 'Tertunda',
    'reminder.status.scheduled': 'Dijadwalkan',
    'reminder.status.queued': 'Dalam Antrian',
    'reminder.status.sending': 'Mengirim...',
    'reminder.status.sent': 'Terkirim',
    'reminder.status.delivered': 'Diterima',
    'reminder.status.read': 'Dibaca',
    'reminder.status.failed': 'Gagal',
    'reminder.status.expired': 'Kedaluwarsa',
    'reminder.status.cancelled': 'Dibatalkan',
    'reminder.status.retry': 'Coba Lagi'
  };

  const mockT = (key, options) => {
    if (key === 'reminder.status.retrying') {
      return `Mengirim ulang (${options?.values?.count || 0}/${options?.values?.max || 3})`;
    }
    return translations[key] || key;
  };

  // Create a callable store (works both as $t and t())
  const tStore = readable(mockT);
  const t = Object.assign(mockT, { subscribe: tStore.subscribe });

  return {
    t,
    locale: readable('id'),
    locales: readable(['id', 'en']),
    loading: readable(false),
    init: vi.fn(),
    getLocaleFromNavigator: vi.fn(() => 'id'),
    addMessages: vi.fn(),
    _: t
  };
});

import DeliveryStatusBadge from './DeliveryStatusBadge.svelte';

// Clean up after each test
afterEach(() => {
  vi.restoreAllMocks();
  cleanup();
});

describe('DeliveryStatusBadge Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Status Configuration Mapping (AC Compliant)', () => {
    it('should render pending status with clock icon', async () => {
      render(DeliveryStatusBadge, { status: 'pending' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Tertunda');
      expect(status).toHaveTextContent('🕐');  // AC: clock icon
      expect(status).toHaveTextContent('Tertunda');
    });

    it('should render scheduled status', async () => {
      render(DeliveryStatusBadge, { status: 'scheduled' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Dijadwalkan');
      expect(status).toHaveTextContent('📅');
    });

    it('should render queued status', async () => {
      render(DeliveryStatusBadge, { status: 'queued' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Dalam Antrian');
    });

    it('should render sending status with spinner/hourglass and pulse animation', async () => {
      render(DeliveryStatusBadge, { status: 'sending' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Mengirim...');
      expect(status).toHaveTextContent('⏳');  // AC: spinner icon (hourglass representation)
      expect(status.className).toContain('animate-pulse');
    });

    it('should render sent status with single checkmark in gray', async () => {
      render(DeliveryStatusBadge, { status: 'sent' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Terkirim');
      expect(status).toHaveTextContent('✓');  // AC: single checkmark
      // Verify gray color class
      const iconSpan = status.querySelector('span[aria-hidden="true"]');
      expect(iconSpan.className).toContain('text-[#64748b]');
    });

    it('should render delivered status with double checkmarks in WhatsApp green', async () => {
      render(DeliveryStatusBadge, { status: 'delivered' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Diterima');
      expect(status).toHaveTextContent('✓✓');  // AC: double checkmarks
      // Verify WhatsApp green color
      const iconSpan = status.querySelector('span[aria-hidden="true"]');
      expect(iconSpan.className).toContain('text-[#25D366]');
    });

    it('should render read status with double checkmarks in WhatsApp blue', async () => {
      render(DeliveryStatusBadge, { status: 'read' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Dibaca');
      expect(status).toHaveTextContent('✓✓');  // AC: double checkmarks
      // Verify WhatsApp blue color
      const iconSpan = status.querySelector('span[aria-hidden="true"]');
      expect(iconSpan.className).toContain('text-[#53bdeb]');
    });

    it('should render failed status with X icon in red', async () => {
      render(DeliveryStatusBadge, { status: 'failed' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Gagal');
      expect(status).toHaveTextContent('✕');  // AC: X icon
      // Verify red color
      const iconSpan = status.querySelector('span[aria-hidden="true"]');
      expect(iconSpan.className).toContain('text-[#dc2626]');
    });

    it('should render expired status', async () => {
      render(DeliveryStatusBadge, { status: 'expired' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Kedaluwarsa');
      expect(status).toHaveTextContent('⏰');
    });

    it('should render cancelled status with X icon in amber', async () => {
      render(DeliveryStatusBadge, { status: 'cancelled' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Dibatalkan');
      expect(status).toHaveTextContent('✕');  // AC: X icon
      // Verify amber color
      const iconSpan = status.querySelector('span[aria-hidden="true"]');
      expect(iconSpan.className).toContain('text-amber-600');
    });

    it('should fallback to pending status for unknown status', async () => {
      render(DeliveryStatusBadge, { status: 'unknown' });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveAttribute('aria-label', 'Tertunda');
    });
  });

  describe('Retrying Status', () => {
    it('should render retrying status with retry count', async () => {
      render(DeliveryStatusBadge, { status: 'retrying', retryCount: 2, maxAttempts: 3 });

      const status = screen.getByRole('status');
      expect(status).toBeInTheDocument();
      expect(status).toHaveTextContent('🔄');
      expect(status).toHaveTextContent('(2/3)');
    });

    it('should use default maxAttempts of 3', async () => {
      render(DeliveryStatusBadge, { status: 'retrying', retryCount: 1 });

      const status = screen.getByRole('status');
      expect(status).toHaveTextContent('(1/3)');
    });
  });

  describe('Retry Button (AC Requirement)', () => {
    it('should show retry button when status is failed and onRetry is provided', async () => {
      const onRetry = vi.fn();
      render(DeliveryStatusBadge, { status: 'failed', onRetry });

      const retryButton = screen.getByRole('button', { name: /coba lagi/i });
      expect(retryButton).toBeInTheDocument();
    });

    it('should NOT show retry button when status is failed but onRetry is not provided', async () => {
      render(DeliveryStatusBadge, { status: 'failed' });

      const retryButton = screen.queryByRole('button');
      expect(retryButton).not.toBeInTheDocument();
    });

    it('should NOT show retry button for non-failed statuses', async () => {
      const onRetry = vi.fn();
      render(DeliveryStatusBadge, { status: 'delivered', onRetry });

      const retryButton = screen.queryByRole('button');
      expect(retryButton).not.toBeInTheDocument();
    });

    it('should call onRetry when retry button is clicked', async () => {
      const onRetry = vi.fn();
      render(DeliveryStatusBadge, { status: 'failed', onRetry });

      const retryButton = screen.getByRole('button', { name: /coba lagi/i });
      await fireEvent.click(retryButton);

      expect(onRetry).toHaveBeenCalledTimes(1);
    });

    it('should disable retry button when isRetrying is true', async () => {
      const onRetry = vi.fn();
      render(DeliveryStatusBadge, { status: 'failed', onRetry, isRetrying: true });

      const retryButton = screen.getByRole('button');
      expect(retryButton).toBeDisabled();
    });

    it('should show spinner when isRetrying is true', async () => {
      const onRetry = vi.fn();
      render(DeliveryStatusBadge, { status: 'failed', onRetry, isRetrying: true });

      const retryButton = screen.getByRole('button');
      const spinner = retryButton.querySelector('svg.animate-spin');
      expect(spinner).toBeInTheDocument();
    });

    it('should not call onRetry when button is disabled', async () => {
      const onRetry = vi.fn();
      render(DeliveryStatusBadge, { status: 'failed', onRetry, isRetrying: true });

      const retryButton = screen.getByRole('button');
      await fireEvent.click(retryButton);

      expect(onRetry).not.toHaveBeenCalled();
    });
  });

  describe('Accessibility', () => {
    it('should have role="status" for screen readers', async () => {
      render(DeliveryStatusBadge, { status: 'delivered' });

      expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('should have aria-label with status text', async () => {
      render(DeliveryStatusBadge, { status: 'delivered' });

      const status = screen.getByRole('status');
      expect(status).toHaveAttribute('aria-label', 'Diterima');
    });

    it('should have aria-live="polite" for real-time updates', async () => {
      render(DeliveryStatusBadge, { status: 'delivered' });

      const status = screen.getByRole('status');
      expect(status).toHaveAttribute('aria-live', 'polite');
    });

    it('should have aria-atomic="true" for complete announcements', async () => {
      render(DeliveryStatusBadge, { status: 'delivered' });

      const status = screen.getByRole('status');
      expect(status).toHaveAttribute('aria-atomic', 'true');
    });

    it('should have icon with aria-hidden="true"', async () => {
      render(DeliveryStatusBadge, { status: 'delivered' });

      const status = screen.getByRole('status');
      const iconSpan = status.querySelector('span[aria-hidden="true"]');
      expect(iconSpan).toBeInTheDocument();
    });

    it('should have screen-reader-only text for status label', async () => {
      render(DeliveryStatusBadge, { status: 'delivered' });

      const status = screen.getByRole('status');
      const srOnly = status.querySelector('.sr-only');
      expect(srOnly).toBeInTheDocument();
      expect(srOnly).toHaveTextContent('Diterima');
    });

    it('should have accessible retry button with aria-label', async () => {
      const onRetry = vi.fn();
      render(DeliveryStatusBadge, { status: 'failed', onRetry });

      const retryButton = screen.getByRole('button');
      expect(retryButton).toHaveAttribute('aria-label', 'Coba Lagi');
    });
  });

  describe('Visual Properties', () => {
    it('should have inline-flex layout', async () => {
      render(DeliveryStatusBadge, { status: 'pending' });

      const status = screen.getByRole('status');
      expect(status.className).toContain('inline-flex');
    });

    it('should have rounded-full class for badge shape', async () => {
      render(DeliveryStatusBadge, { status: 'pending' });

      const status = screen.getByRole('status');
      expect(status.className).toContain('rounded-full');
    });

    it('should have correct background color for delivered status', async () => {
      render(DeliveryStatusBadge, { status: 'delivered' });

      const status = screen.getByRole('status');
      expect(status.className).toContain('bg-emerald-100');
    });

    it('should have correct background color for failed status', async () => {
      render(DeliveryStatusBadge, { status: 'failed' });

      const status = screen.getByRole('status');
      expect(status.className).toContain('bg-red-100');
    });
  });

  describe('Props Defaults', () => {
    it('should default status to "pending" when not provided', async () => {
      render(DeliveryStatusBadge, {});

      const status = screen.getByRole('status');
      expect(status).toHaveAttribute('aria-label', 'Tertunda');
    });

    it('should default retryCount to 0', async () => {
      render(DeliveryStatusBadge, { status: 'retrying' });

      const status = screen.getByRole('status');
      expect(status).toHaveTextContent('(0/3)');
    });

    it('should default maxAttempts to 3', async () => {
      render(DeliveryStatusBadge, { status: 'retrying', retryCount: 1 });

      const status = screen.getByRole('status');
      expect(status).toHaveTextContent('(1/3)');
    });

    it('should default onRetry to null', async () => {
      render(DeliveryStatusBadge, { status: 'failed' });

      const retryButton = screen.queryByRole('button');
      expect(retryButton).not.toBeInTheDocument();
    });

    it('should default isRetrying to false', async () => {
      const onRetry = vi.fn();
      render(DeliveryStatusBadge, { status: 'failed', onRetry });

      const retryButton = screen.getByRole('button');
      expect(retryButton).not.toBeDisabled();
    });
  });
});
