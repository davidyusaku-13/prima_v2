/**
 * @vitest-environment happy-dom
 */
/**
 * ContentPreviewPanel Tests
 *
 * Tests for the ContentPreviewPanel component that provides preview functionality
 * for articles and videos in the content picker modal.
 *
 * Uses @testing-library/svelte for component testing with Vitest.
 *
 * Test Coverage:
 * - Article preview rendering (title, excerpt, hero image, date)
 * - Video preview rendering (thumbnail, duration, channel, play icon)
 * - Keyboard navigation (Escape, attach button)
 * - Backdrop click handling
 * - Excerpt truncation (200 char limit)
 * - Accessibility (ARIA attributes, keyboard navigation)
 * - Date formatting (valid and missing dates)
 * - Image fallback handling
 * - Content type detection (article vs video)
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor, cleanup } from '@testing-library/svelte';
import { tick } from 'svelte';

// Mock svelte-i18n with proper store - make t both callable and subscribable
vi.mock('svelte-i18n', async () => {
  const { readable } = await import('svelte/store');

  const translations = {
    'content.picker.articles': 'Artikel',
    'content.picker.videos': 'Video',
    'content.preview.attach': 'Lampirkan',
    'content.preview.selected': 'Dipilih',
    'content.attribution.publishedOn': 'Diterbitkan pada',
    'content.attribution.author': 'Penulis',
    'content.attribution.channel': 'Channel',
    'content.attribution.uploadedOn': 'Diunggah pada',
    'content.disclaimer.text': 'Konten ini hanya untuk edukasi kesehatan.',
    'common.close': 'Tutup'
  };

  const mockT = (key, params) => translations[key] || key;

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

// Import component after mocking
import ContentPreviewPanel from './ContentPreviewPanel.svelte';

// Clean up after each test to prevent state leakage
afterEach(() => {
  vi.restoreAllMocks();
  cleanup();
});

// Test data
const mockArticle = {
  id: '1',
  title: 'Test Article Title',
  excerpt: 'This is a test excerpt that is longer than 200 characters to test the truncation functionality of the preview panel. It should be truncated properly when displayed in the UI.',
  publishedAt: '2024-01-15T10:00:00Z',
  heroImages: {
    hero16x9: '/uploads/articles/hero-16x9-1.jpg',
    hero1x1: '/uploads/articles/hero-1x1-1.jpg'
  }
};

const mockVideo = {
  id: '2',
  title: 'Test Video Title',
  YouTubeID: 'dQw4w9WgXcQ',
  thumbnailURL: 'https://img.youtube.com/vi/dQw4w9WgXcQ/maxresdefault.jpg',
  duration: '10:30',
  channelName: 'Test Channel'
};

describe('ContentPreviewPanel Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Article Preview', () => {
    it('should render article preview with title, excerpt, and date', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        expect(screen.getByText('Test Article Title')).toBeInTheDocument();
      });

      // Check truncated excerpt (should be ~200 chars with ...)
      const excerpt = screen.getByText(/This is a test excerpt/);
      expect(excerpt).toBeInTheDocument();

      // Check for date - "Diterbitkan pada" is in aria-label (Indonesian), date is visible text
      const dateElement = screen.getByLabelText(/Diterbitkan pada/);
      expect(dateElement).toBeInTheDocument();
    });

    it('should show hero image for article', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        const img = screen.getByAltText('Test Article Title');
        expect(img).toHaveAttribute('src', '/uploads/articles/hero-16x9-1.jpg');
      });
    });

    it('should show Attach button when not selected', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        expect(screen.getByText('Lampirkan')).toBeInTheDocument();
      });

      const attachButton = screen.getByText('Lampirkan').closest('button');
      expect(attachButton).toBeInTheDocument();
      expect(attachButton).not.toHaveTextContent(/Dipilih/);
    });

    it('should show Selected button when isSelected=true', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: true, onClose, onAttach }
      });

      await waitFor(() => {
        expect(screen.getByText('Dipilih')).toBeInTheDocument();
      });

      const selectedButton = screen.getByText('Dipilih').closest('button');
      expect(selectedButton).toBeInTheDocument();
    });
  });

  describe('Video Preview', () => {
    it('should render video preview with thumbnail, duration, and channel', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockVideo, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        expect(screen.getByText('Test Video Title')).toBeInTheDocument();
        expect(screen.getByText('Test Channel')).toBeInTheDocument();
        expect(screen.getByText('10:30')).toBeInTheDocument();
      });
    });

    it('should show play icon overlay on video thumbnail', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockVideo, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Check that play icon SVG is rendered
        const playIcon = screen.getByRole('img', { hidden: true });
        expect(playIcon).toBeInTheDocument();
      });
    });

    it('should display video type badge', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockVideo, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        expect(screen.getByText('Video')).toBeInTheDocument();
      });
    });
  });

  describe('Keyboard Navigation', () => {
    it('should close preview on Escape key', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      const { container } = render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Use container.querySelector because dialog is inside aria-hidden backdrop
        const dialog = container.querySelector('[role="dialog"]');
        expect(dialog).toBeInTheDocument();
      });

      // Press Escape - component listens on window via <svelte:window>
      await fireEvent.keyDown(window, { key: 'Escape', code: 'Escape' });

      expect(onClose).toHaveBeenCalledTimes(1);
    });

    it('should call onAttach when Attach button is clicked', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        expect(screen.getByText('Lampirkan')).toBeInTheDocument();
      });

      const attachButton = screen.getByText('Lampirkan').closest('button');
      await fireEvent.click(attachButton);

      expect(onAttach).toHaveBeenCalledTimes(1);
    });

    it('should call onClose when close button is clicked', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        expect(screen.getByLabelText(/Tutup pratinjau/i)).toBeInTheDocument();
      });

      const closeButton = screen.getByLabelText(/Tutup pratinjau/i);
      await fireEvent.click(closeButton);

      expect(onClose).toHaveBeenCalledTimes(1);
    });
  });

  describe('Backdrop Click', () => {
    it('should close when clicking outside the panel', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      const { container } = render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Use container.querySelector because dialog is inside aria-hidden backdrop
        const dialog = container.querySelector('[role="dialog"]');
        expect(dialog).toBeInTheDocument();
      });

      // The backdrop is the container div with role="presentation" and aria-hidden
      const backdrop = container.querySelector('[role="presentation"]');
      await fireEvent.click(backdrop);

      expect(onClose).toHaveBeenCalledTimes(1);
    });
  });

  describe('Excerpt Truncation', () => {
    it('should truncate long excerpts to 200 characters', async () => {
      const longExcerptArticle = {
        ...mockArticle,
        excerpt: 'A'.repeat(300)
      };

      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: longExcerptArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Use robust assertion: check length is 200 and ends with ellipsis
        const excerptText = screen.getByText((content) => {
          return content.length === 200 && content.endsWith('...');
        });
        expect(excerptText).toBeInTheDocument();
      });
    });

    it('should not truncate short excerpts', async () => {
      const shortExcerptArticle = {
        ...mockArticle,
        excerpt: 'Short excerpt'
      };

      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: shortExcerptArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        expect(screen.getByText('Short excerpt')).toBeInTheDocument();
      });
    });
  });

  describe('Accessibility', () => {
    it('should have proper ARIA dialog attributes', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      const { container } = render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Use container.querySelector because dialog is inside aria-hidden backdrop
        const dialog = container.querySelector('[role="dialog"]');
        expect(dialog).toBeInTheDocument();
        expect(dialog).toHaveAttribute('aria-modal', 'true');
        expect(dialog).toHaveAttribute('aria-labelledby', 'preview-title');
      });
    });

    it('should have close button with proper aria-label', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Close button has aria-label but no visible text
        const closeButton = screen.getByLabelText(/Tutup pratinjau/i);
        expect(closeButton).toBeInTheDocument();
      });
    });

    it('should have attach button with aria-pressed state', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: true, onClose, onAttach }
      });

      await waitFor(() => {
        // Button text is "Dipilih" (Indonesian for "Selected")
        const attachButton = screen.getByText('Dipilih').closest('button');
        expect(attachButton).toHaveAttribute('aria-pressed', 'true');
      });
    });
  });

  describe('Date Formatting', () => {
    it('should handle valid date', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Date should be formatted in Indonesian locale
        expect(screen.getByText(/2024/)).toBeInTheDocument();
        expect(screen.getByText(/Januari/)).toBeInTheDocument();
      });
    });

    it('should handle missing date gracefully', async () => {
      const noDateArticle = {
        ...mockArticle,
        publishedAt: null,
        createdAt: null
      };

      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: noDateArticle, isSelected: false, onClose, onAttach }
      });

      // Should not crash and should render without error
      await waitFor(() => {
        expect(screen.getByText('Test Article Title')).toBeInTheDocument();
      });
    });
  });

  describe('Image Fallback', () => {
    it('should show placeholder when hero image is missing', async () => {
      const noImageArticle = {
        ...mockArticle,
        heroImages: null
      };

      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: noImageArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Should still render the component without hero image
        expect(screen.getByText('Test Article Title')).toBeInTheDocument();
      });
    });

    it('should show placeholder when video thumbnail is missing', async () => {
      const noThumbnailVideo = {
        ...mockVideo,
        thumbnailURL: null
      };

      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: noThumbnailVideo, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Should still render video title
        expect(screen.getByText('Test Video Title')).toBeInTheDocument();
      });
    });
  });

  describe('Content Type Detection', () => {
    it('should detect video by YouTubeID', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockVideo, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Video badge should be shown (Indonesian: "Video")
        expect(screen.getByText('Video')).toBeInTheDocument();
        // Duration should be displayed
        expect(screen.getByText('10:30')).toBeInTheDocument();
      });
    });

    it('should show article badge for content without YouTubeID', async () => {
      const onClose = vi.fn();
      const onAttach = vi.fn();

      render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
      });

      await waitFor(() => {
        // Article badge should be shown (Indonesian: "Artikel")
        expect(screen.getByText('Artikel')).toBeInTheDocument();
      });
    });
  });
});
