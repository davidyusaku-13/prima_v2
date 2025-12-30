/**
 * @vitest-environment happy-dom
 */
import { describe, it, expect, vi, afterEach } from 'vitest';
import { render, fireEvent, cleanup } from '@testing-library/svelte';

// Mock svelte-i18n with proper store
vi.mock('svelte-i18n', async () => {
  const { readable } = await import('svelte/store');

  const translations = {
    'delivery.filter.label': 'Filter Status Pengiriman',
    'delivery.filter.all': 'Semua',
    'delivery.filter.pending': 'Pending',
    'delivery.filter.sent': 'Terkirim',
    'delivery.filter.failed': 'Gagal'
  };

  const mockT = (key) => translations[key] || key;

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

import DeliveryStatusFilter from './DeliveryStatusFilter.svelte';

// Clean up after each test
afterEach(() => {
  cleanup();
});

describe('DeliveryStatusFilter', () => {
  it('renders all filter buttons', () => {
    const { getByText } = render(DeliveryStatusFilter, {
      props: {
        selectedFilter: 'all',
        counts: { all: 10, pending: 3, sent: 5, failed: 2 }
      }
    });

    expect(getByText(/Semua/i)).toBeTruthy();
    expect(getByText(/Pending/i)).toBeTruthy();
    expect(getByText(/Terkirim/i)).toBeTruthy();
    expect(getByText(/Gagal/i)).toBeTruthy();
  });

  it('displays count badges for each filter', () => {
    const { getByText } = render(DeliveryStatusFilter, {
      props: {
        selectedFilter: 'all',
        counts: { all: 10, pending: 3, sent: 5, failed: 2 }
      }
    });

    expect(getByText('10')).toBeTruthy();
    expect(getByText('3')).toBeTruthy();
    expect(getByText('5')).toBeTruthy();
    expect(getByText('2')).toBeTruthy();
  });

  it('applies active styling to selected filter', () => {
    const { container } = render(DeliveryStatusFilter, {
      props: {
        selectedFilter: 'failed',
        counts: { all: 10, pending: 3, sent: 5, failed: 2 }
      }
    });

    const buttons = container.querySelectorAll('button');
    const failedButton = Array.from(buttons).find(btn =>
      btn.textContent.includes('Gagal')
    );

    expect(failedButton.classList.contains('bg-teal-600')).toBe(true);
    expect(failedButton.classList.contains('text-white')).toBe(true);
  });

  it('calls onFilterChange callback when filter is clicked', async () => {
    const mockHandler = vi.fn();
    const { getByText } = render(DeliveryStatusFilter, {
      props: {
        selectedFilter: 'all',
        counts: { all: 10, pending: 3, sent: 5, failed: 2 },
        onFilterChange: mockHandler
      }
    });

    const failedButton = getByText(/Gagal/i);
    await fireEvent.click(failedButton);

    expect(mockHandler).toHaveBeenCalledTimes(1);
    expect(mockHandler).toHaveBeenCalledWith('failed');
  });

  it('has proper accessibility attributes', () => {
    const { container } = render(DeliveryStatusFilter, {
      props: {
        selectedFilter: 'all',
        counts: { all: 10, pending: 3, sent: 5, failed: 2 }
      }
    });

    const tablist = container.querySelector('[role="tablist"]');
    expect(tablist).toBeTruthy();

    const tabs = container.querySelectorAll('[role="tab"]');
    expect(tabs.length).toBe(4);

    tabs.forEach(tab => {
      expect(tab.hasAttribute('aria-selected')).toBe(true);
      expect(tab.hasAttribute('aria-label')).toBe(true);
    });
  });

  it('does not show count badge when count is 0', () => {
    const { container } = render(DeliveryStatusFilter, {
      props: {
        selectedFilter: 'all',
        counts: { all: 0, pending: 0, sent: 0, failed: 0 }
      }
    });

    const badges = container.querySelectorAll('.ml-2.px-2');
    expect(badges.length).toBe(0);
  });
});
