import { describe, it, expect, vi } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/svelte';
import Sidebar from './Sidebar.svelte';
import BottomNav from './BottomNav.svelte';

// Mock i18n
vi.mock('svelte-i18n', () => ({
  t: vi.fn((key) => {
    const translations = {
      'navigation.dashboard': 'Dashboard',
      'navigation.patients': 'Patients',
      'navigation.users': 'Users',
      'navigation.analytics': 'Analytics',
      'cms.dashboard': 'CMS Dashboard',
      'common.cms': 'CMS',
      'common.more': 'More',
      'berita.title': 'Health News',
      'video.title': 'Educational Videos'
    };
    return translations[key] || key;
  })
}));

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
      const { container } = render(BottomNav, {
        props: {
          user: { role: 'admin', fullName: 'Admin User' },
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
      const { container } = render(BottomNav, {
        props: {
          user: { role: 'volunteer', fullName: 'Volunteer User' },
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
      const { container } = render(BottomNav, {
        props: {
          user: { role: 'superadmin', fullName: 'Admin User' },
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
