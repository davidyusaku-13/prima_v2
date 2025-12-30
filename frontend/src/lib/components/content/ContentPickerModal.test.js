/**
 * @vitest-environment happy-dom
 */
import { render, screen, fireEvent, waitFor, cleanup } from '@testing-library/svelte';
import { vi, describe, it, expect, beforeEach, afterEach } from 'vitest';
import { readable } from 'svelte/store';
import ContentPickerModal from './ContentPickerModal.svelte';
import { fetchAllContent, fetchCategories, fetchPopularContent, incrementAttachmentCount } from '$lib/utils/api.js';

// Constants for debounce timing (must match component implementation)
const DEBOUNCE_MS = 300;
const TEST_BUFFER_MS = 200;

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(() => null),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn()
};
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

// Mock API modules with default return values
vi.mock('$lib/utils/api.js', () => ({
  fetchAllContent: vi.fn().mockResolvedValue({ articles: [], videos: [] }),
  fetchCategories: vi.fn().mockResolvedValue([]),
  fetchPopularContent: vi.fn().mockResolvedValue([]),
  incrementAttachmentCount: vi.fn().mockResolvedValue({})
}));

// Mock svelte-i18n with proper store - make t both callable and subscribable
vi.mock('svelte-i18n', async () => {
  const { readable } = await import('svelte/store');

  const translations = {
    'content.picker.title': 'Pilih Konten Edukasi',
    'content.picker.search': 'Cari artikel atau video...',
    'content.picker.tab.all': 'Semua',
    'content.picker.tab.articles': 'Artikel',
    'content.picker.tab.videos': 'Video',
    'content.picker.articles': 'Artikel',
    'content.picker.videos': 'Video',
    'content.picker.attach': 'Lampirkan',
    'content.picker.noResults': 'Tidak ada hasil',
    'content.picker.empty': 'Belum ada konten edukasi',
    'common.cancel': 'Batal'
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

describe('ContentPickerModal', () => {
  let mockOnClose;
  let mockOnSelect;

  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.getItem.mockReturnValue(null);
    mockOnClose = vi.fn();
    mockOnSelect = vi.fn();

    // Default mock implementations
    fetchAllContent.mockResolvedValue({
      articles: [
        { id: 'art-1', title: 'Healthy Eating Guide', excerpt: 'A guide to eating healthy' },
        { id: 'art-2', title: 'Exercise Tips', excerpt: 'Tips for better exercise' }
      ],
      videos: [
        { id: 'vid-1', title: 'Morning Yoga', YouTubeID: 'abc123', duration: '5:30' }
      ]
    });

    fetchCategories.mockResolvedValue([
      { id: 'cat-1', name: 'Health', type: 'article' },
      { id: 'cat-2', name: 'Nutrition', type: 'video' }
    ]);

    fetchPopularContent.mockResolvedValue([]);
    incrementAttachmentCount.mockResolvedValue({});
  });

  afterEach(() => {
    cleanup();
    vi.clearAllMocks();
  });

  describe('Basic Rendering', () => {
    it('renders modal when show is true', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByRole('dialog')).toBeInTheDocument();
      });
    });

    it('does not render when show is false', () => {
      render(ContentPickerModal, {
        show: false,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });

    it('displays title correctly', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByText('Pilih Konten Edukasi')).toBeInTheDocument();
      });
    });

    it('fetches content on mount', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(fetchAllContent).toHaveBeenCalledWith(null, 'all');
      });
    });

    it('fetches categories on mount', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(fetchCategories).toHaveBeenCalled();
      });
    });
  });

  describe('Search Functionality', () => {
    it('displays search input', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByPlaceholderText('Cari artikel atau video...')).toBeInTheDocument();
      });
    });

    it('has tabindex attribute on search input', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      const searchInput = await screen.findByPlaceholderText('Cari artikel atau video...');
      expect(searchInput).toHaveAttribute('tabindex', '0');
    });

    it('shows clear button when search query is not empty', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      const searchInput = await screen.findByPlaceholderText('Cari artikel atau video...');
      await fireEvent.input(searchInput, { target: { value: 'test' } });

      // Clear button should appear
      const clearButton = screen.getByLabelText('Hapus pencarian');
      expect(clearButton).toBeInTheDocument();
    });

    it('clears search when clear button is clicked', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      const searchInput = await screen.findByPlaceholderText('Cari artikel atau video...');
      await fireEvent.input(searchInput, { target: { value: 'test' } });

      const clearButton = screen.getByLabelText('Hapus pencarian');
      await fireEvent.click(clearButton);

      expect(searchInput).toHaveValue('');
    });
  });

  describe.skip('Category Tabs', () => {
    it('displays three tabs: All, Articles, Videos', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByRole('tablist')).toBeInTheDocument();
        expect(screen.getByRole('tab', { name: 'Semua' })).toBeInTheDocument();
        expect(screen.getByRole('tab', { name: 'Artikel' })).toBeInTheDocument();
        expect(screen.getByRole('tab', { name: 'Video' })).toBeInTheDocument();
      });
    });

    it('highlights active tab with teal styling', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        const activeTab = screen.getByRole('tab', { name: 'Semua' });
        expect(activeTab).toHaveClass('bg-teal-50');
      });
    });

    it('filters content when tab is clicked', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        // Click Articles tab
        const articlesTab = screen.getByRole('tab', { name: 'Artikel' });
        fireEvent.click(articlesTab);

        // Should show only articles
        expect(screen.getByText('Healthy Eating Guide')).toBeInTheDocument();
        expect(screen.queryByText('Morning Yoga')).not.toBeInTheDocument();
      });
    });
  });

  describe.skip('Content Selection', () => {
    it.skip('displays articles and videos', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByText('Healthy Eating Guide')).toBeInTheDocument();
        expect(screen.getByText('Morning Yoga')).toBeInTheDocument();
      });
    });

    it.skip('selects content when clicked', async () => {
      render(ContentPickerModal, {
        show: true,
        selectedContent: [],
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        const article = screen.getByText('Healthy Eating Guide').closest('button');
        fireEvent.click(article);

        expect(mockOnSelect).toHaveBeenCalledWith([
          expect.objectContaining({ id: 'art-1', type: 'article' })
        ]);
      });
    });

    it.skip('deselects content when clicked again', async () => {
      const selectedContent = [
        { id: 'art-1', title: 'Healthy Eating Guide', type: 'article' }
      ];

      render(ContentPickerModal, {
        show: true,
        selectedContent,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        const article = screen.getByText('Healthy Eating Guide').closest('button');
        fireEvent.click(article);

        expect(mockOnSelect).toHaveBeenCalledWith([]);
      });
    });

    it.skip('shows selection count', async () => {
      render(ContentPickerModal, {
        show: true,
        selectedContent: [
          { id: 'art-1', title: 'Healthy Eating Guide', type: 'article' }
        ],
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByText('1 dipilih')).toBeInTheDocument();
      });
    });

    it.skip('limits selection to 3 items', async () => {
      render(ContentPickerModal, {
        show: true,
        selectedContent: [
          { id: 'art-1', title: 'Article 1', type: 'article' },
          { id: 'art-2', title: 'Article 2', type: 'article' },
          { id: 'vid-1', title: 'Video 1', type: 'video' }
        ],
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        const newArticle = screen.getByText('Exercise Tips').closest('button');
        fireEvent.click(newArticle);

        // onSelect should not be called since limit is reached
        expect(mockOnSelect).not.toHaveBeenCalled();
      });
    });
  });

  describe.skip('Keyboard Navigation', () => {
    it('closes modal on Escape key', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        fireEvent.keyDown(window, { key: 'Escape' });
        expect(mockOnClose).toHaveBeenCalled();
      });
    });

    it('focuses search input on mount', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        const searchInput = screen.getByPlaceholderText('Cari artikel atau video...');
        expect(searchInput).toHaveFocus();
      });
    });
  });

  describe.skip('Error Handling', () => {
    it('displays error message on fetch failure', async () => {
      fetchAllContent.mockRejectedValue(new Error('Network error'));

      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByText('Network error')).toBeInTheDocument();
      });
    });

    it('shows refresh button on error', async () => {
      fetchAllContent.mockRejectedValue(new Error('Network error'));

      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Refresh' })).toBeInTheDocument();
      });
    });
  });

  describe.skip('Empty State', () => {
    it('shows empty state when no content', async () => {
      fetchAllContent.mockResolvedValue({ articles: [], videos: [] });

      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByText('Belum ada konten edukasi')).toBeInTheDocument();
      });
    });

    it('shows no results message when search has no matches', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(async () => {
        const searchInput = screen.getByPlaceholderText('Cari artikel atau video...');
        await fireEvent.input(searchInput, { target: { value: 'nonexistent content' } });

        // Wait for debounce + buffer for async filtering
        await new Promise(resolve => setTimeout(resolve, DEBOUNCE_MS + TEST_BUFFER_MS));

        expect(screen.getByText(/Tidak ada hasil/)).toBeInTheDocument();
      });
    });
  });

  describe.skip('Footer Actions', () => {
    it('displays Cancel button', async () => {
      render(ContentPickerModal, {
        show: true,
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Batal' })).toBeInTheDocument();
      });
    });

    it('displays Attach button with selection count', async () => {
      render(ContentPickerModal, {
        show: true,
        selectedContent: [
          { id: 'art-1', title: 'Healthy Eating Guide', type: 'article' }
        ],
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Lampirkan (1)' })).toBeInTheDocument();
      });
    });

    it('disables Attach button when no content selected', async () => {
      render(ContentPickerModal, {
        show: true,
        selectedContent: [],
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        const attachButton = screen.getByRole('button', { name: 'Lampirkan (0)' });
        expect(attachButton).toBeDisabled();
      });
    });

    it('calls onSelect and closes when Attach is clicked', async () => {
      render(ContentPickerModal, {
        show: true,
        selectedContent: [
          { id: 'art-1', title: 'Healthy Eating Guide', type: 'article' }
        ],
        onClose: mockOnClose,
        onSelect: mockOnSelect
      });

      await waitFor(() => {
        const attachButton = screen.getByRole('button', { name: 'Lampirkan (1)' });
        fireEvent.click(attachButton);

        expect(mockOnSelect).toHaveBeenCalledWith([
          expect.objectContaining({ id: 'art-1', type: 'article' })
        ]);
        expect(mockOnClose).toHaveBeenCalled();
      });
    });
  });
});
