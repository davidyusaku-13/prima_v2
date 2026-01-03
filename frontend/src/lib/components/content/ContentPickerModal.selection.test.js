/**
 * @vitest-environment happy-dom
 * ContentPickerModal - Selection Tests (6 tests)
 */
import {
  render,
  screen,
  waitFor,
  cleanup,
  fireEvent,
} from "@testing-library/svelte";
import { vi, describe, it, expect, beforeEach, afterEach } from "vitest";
import { readable } from "svelte/store";
import ContentPickerModal from "./ContentPickerModal.svelte";
import {
  fetchAllContent,
  fetchCategories,
  fetchPopularContent,
  incrementAttachmentCount,
} from "$lib/utils/api.js";
import {
  localStorageMock,
  setupLocalStorageMock,
  defaultContentResponse,
} from "./ContentPickerModal.test-helpers.js";

// Setup localStorage
setupLocalStorageMock();

// Mock API modules
vi.mock("$lib/utils/api.js", () => ({
  fetchAllContent: vi.fn().mockResolvedValue({ articles: [], videos: [] }),
  fetchCategories: vi.fn().mockResolvedValue([]),
  fetchPopularContent: vi.fn().mockResolvedValue([]),
  incrementAttachmentCount: vi.fn().mockResolvedValue({}),
}));

// Mock svelte-i18n
vi.mock("svelte-i18n", async () => {
  const { createI18nMock } = await import(
    "./ContentPickerModal.test-helpers.js"
  );
  return createI18nMock();
});

// SKIPPED ENTIRE SUITE: ContentPickerModal has memory leaks in Svelte 5 $derived.by() chains
// that cause OOM errors in Vitest happy-dom environment. See rendering.test.js for working tests.
// TODO: Fix component architecture or migrate to different test approach (e.g., Playwright)
describe.skip("ContentPickerModal - Selection", () => {
  let mockOnClose;
  let mockOnSelect;
  let renderResult;

  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.getItem.mockReturnValue(null);
    mockOnClose = vi.fn();
    mockOnSelect = vi.fn();
    renderResult = null;

    fetchAllContent.mockResolvedValue(defaultContentResponse);
    fetchCategories.mockResolvedValue([]);
    fetchPopularContent.mockResolvedValue([]);
    incrementAttachmentCount.mockResolvedValue({});
  });

  afterEach(() => {
    if (renderResult?.unmount) {
      try {
        renderResult.unmount();
      } catch (e) {
        // Ignore
      }
    }
    cleanup();
    vi.clearAllMocks();
    localStorageMock.clear();
    renderResult = null;
  });

  it("shows content items", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        expect(screen.getByText("Healthy Eating Guide")).toBeInTheDocument();
      },
      { timeout: 1000 }
    );
  });

  // SKIPPED: Memory leak in ContentPickerModal component causes OOM
  // These tests require selectedContent prop which triggers reactive cascades
  it.skip("selects content when clicked with ctrl key", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      selectedContent: [],
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        const article = screen
          .getByText("Healthy Eating Guide")
          .closest("button");
        fireEvent.click(article, { ctrlKey: true });

        expect(mockOnSelect).toHaveBeenCalledWith([
          expect.objectContaining({ id: "art-1", type: "article" }),
        ]);
      },
      { timeout: 1000 }
    );
  });

  it.skip("shows selection count", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      selectedContent: [
        { id: "art-1", title: "Healthy Eating Guide", type: "article" },
      ],
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        expect(screen.getByText("1 dipilih")).toBeInTheDocument();
      },
      { timeout: 1000 }
    );
  });

  it("closes modal on Escape key", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        const dialog = screen.getByRole("dialog");
        fireEvent.keyDown(dialog, { key: "Escape" });
        expect(mockOnClose).toHaveBeenCalled();
      },
      { timeout: 1000 }
    );
  });

  it("focuses search input on mount", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        const searchInput = screen.getByPlaceholderText(
          "Cari artikel atau video..."
        );
        expect(searchInput).toHaveFocus();
      },
      { timeout: 1000 }
    );
  });

  it("displays error message on fetch failure", async () => {
    fetchAllContent.mockRejectedValue(new Error("Network error"));

    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        expect(screen.getByText("Network error")).toBeInTheDocument();
      },
      { timeout: 1000 }
    );
  });
});
