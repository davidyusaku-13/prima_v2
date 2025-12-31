import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import Sidebar from './Sidebar.svelte';
import BottomNav from './BottomNav.svelte';

// Mock svelte-i18n properly with store and function
vi.mock('svelte-i18n', async () => {
  const { readable } = await import('svelte/store');

  const translations = {
    'navigation.dashboard': 'Dashboard',
    'navigation.patients': 'Patients',
    'navigation.users': 'Users',
    'navigation.analytics': 'Analytics',
    'cms.dashboard': 'CMS Dashboard',
    'common.cms': 'CMS',
    'common.more': 'More',
    'berita.title': 'Health News',
    'video.title': 'Educational Videos',
    'users.superadmin': 'Superadmin',
    'users.admin': 'Admin',
    'users.volunteer': 'Volunteer',
    'auth.logout': 'Logout',
    'navigation.volunteerDashboard': 'Volunteer Dashboard',
    'app.name': 'PRIMA'
  };

  const mockT = (key, options) => translations[key] || key;

  // Create a store that returns the function
  const tStore = readable(mockT);

  // t must be both a function and have subscribe for $t to work
  const t = Object.assign(mockT, { subscribe: tStore.subscribe });

  return {
    t,
    _: mockT,
    locale: readable('en'),
    locales: readable(['en', 'id']),
    loading: readable(false),
    init: vi.fn(),
    getLocaleFromNavigator: vi.fn(() => 'en'),
    addMessages: vi.fn()
  };
});

describe('Analytics Navigation', () => {
  describe('Sidebar', () => {
    it('should show analytics link for superadmin', async () => {
      const onNavigate = vi.fn();
      const { container } = render(Sidebar, {
        props: {
          user: { role: 'superadmin', username: 'admin' },
          currentView: 'dashboard',
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      const analyticsButton = container.querySelectorAll('button');
      const hasAnalytics = Array.from(analyticsButton).some(
        btn => btn.textContent?.includes('Analytics')
      );

      expect(hasAnalytics).toBe(true);
    });

    it('should show analytics link for admin', async () => {
      const onNavigate = vi.fn();
      const { container } = render(Sidebar, {
        props: {
          user: { role: 'admin', username: 'admin' },
          currentView: 'dashboard',
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      const analyticsButton = container.querySelectorAll('button');
      const hasAnalytics = Array.from(analyticsButton).some(
        btn => btn.textContent?.includes('Analytics')
      );

      expect(hasAnalytics).toBe(true);
    });

    it('should NOT show analytics link for volunteer', async () => {
      const onNavigate = vi.fn();
      const { container } = render(Sidebar, {
        props: {
          user: { role: 'volunteer', username: 'volunteer' },
          currentView: 'dashboard',
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      const analyticsButton = container.querySelectorAll('button');
      const hasAnalytics = Array.from(analyticsButton).some(
        btn => btn.textContent?.includes('Analytics')
      );

      expect(hasAnalytics).toBe(false);
    });

    it('should call onNavigate with "analytics" when analytics button clicked', async () => {
      const onNavigate = vi.fn();
      const { container } = render(Sidebar, {
        props: {
          user: { role: 'superadmin', username: 'admin' },
          currentView: 'dashboard',
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      const buttons = container.querySelectorAll('button');
      const analyticsButton = Array.from(buttons).find(
        btn => btn.textContent?.includes('Analytics')
      );

      await fireEvent.click(analyticsButton);

      expect(onNavigate).toHaveBeenCalledWith('analytics');
    });

    it('should highlight analytics button when currentView is analytics', async () => {
      const onNavigate = vi.fn();
      const { container } = render(Sidebar, {
        props: {
          user: { role: 'superadmin', username: 'admin' },
          currentView: 'analytics',
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      const buttons = container.querySelectorAll('button');
      const analyticsButton = Array.from(buttons).find(
        btn => btn.textContent?.includes('Analytics')
      );

      // Analytics button should have amber background for active state
      expect(analyticsButton?.className).toContain('bg-amber-50');
    });

    it('should handle null onNavigate gracefully', async () => {
      const { container } = render(Sidebar, {
        props: {
          user: { role: 'superadmin', username: 'admin' },
          currentView: 'dashboard',
          stats: { totalPatients: 5 },
          users: [],
          onNavigate: null
        }
      });

      // Should not throw error
      const buttons = container.querySelectorAll('button');
      const analyticsButton = Array.from(buttons).find(
        btn => btn.textContent?.includes('Analytics')
      );

      await fireEvent.click(analyticsButton);

      // Test passes if no error is thrown
      expect(true).toBe(true);
    });
  });

  describe('BottomNav', () => {
    it('should show analytics link for superadmin', async () => {
      const onNavigate = vi.fn();
      const { container } = render(BottomNav, {
        props: {
          user: { role: 'superadmin', fullName: 'Admin User' },
          currentView: 'analytics', // Must be analytics to show button
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      // Wait for rendering
      await waitFor(() => {
        const buttons = container.querySelectorAll('button');
        const hasAnalytics = Array.from(buttons).some(
          btn => btn.textContent?.includes('Analytics')
        );
        expect(hasAnalytics).toBe(true);
      });
    });

    it('should show analytics link for admin', async () => {
      const onNavigate = vi.fn();
      const { container } = render(BottomNav, {
        props: {
          user: { role: 'admin', fullName: 'Admin User' },
          currentView: 'analytics', // Must be analytics to show button
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      await waitFor(() => {
        const buttons = container.querySelectorAll('button');
        const hasAnalytics = Array.from(buttons).some(
          btn => btn.textContent?.includes('Analytics')
        );
        expect(hasAnalytics).toBe(true);
      });
    });

    it('should NOT show analytics link for volunteer', async () => {
      const onNavigate = vi.fn();
      const { container } = render(BottomNav, {
        props: {
          user: { role: 'volunteer', fullName: 'Volunteer User' },
          currentView: 'analytics',
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      const analyticsButton = container.querySelectorAll('button');
      const hasAnalytics = Array.from(analyticsButton).some(
        btn => btn.textContent?.includes('Analytics')
      );

      expect(hasAnalytics).toBe(false);
    });

    it('should call onNavigate with "analytics" when analytics button clicked', async () => {
      const onNavigate = vi.fn();
      const { container } = render(BottomNav, {
        props: {
          user: { role: 'superadmin', fullName: 'Admin User' },
          currentView: 'analytics', // Must be analytics to show button
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      // Wait for button to render
      await waitFor(() => {
        const buttons = container.querySelectorAll('button');
        const analyticsButton = Array.from(buttons).find(
          btn => btn.textContent?.includes('Analytics')
        );
        expect(analyticsButton).toBeTruthy();
        return analyticsButton;
      });

      const buttons = container.querySelectorAll('button');
      const analyticsButton = Array.from(buttons).find(
        btn => btn.textContent?.includes('Analytics')
      );

      await fireEvent.click(analyticsButton);

      expect(onNavigate).toHaveBeenCalledWith('analytics');
    });

    it('should highlight analytics button when currentView is analytics', async () => {
      const onNavigate = vi.fn();
      const { container } = render(BottomNav, {
        props: {
          user: { role: 'superadmin', fullName: 'Admin User' },
          currentView: 'analytics',
          stats: { totalPatients: 5 },
          users: [],
          onNavigate
        }
      });

      const buttons = container.querySelectorAll('button');
      const analyticsButton = Array.from(buttons).find(
        btn => btn.textContent?.includes('Analytics')
      );

      // Analytics button should have amber background for active state
      expect(analyticsButton?.className).toContain('bg-amber-100');
    });
  });
});
