/**
 * @vitest-environment happy-dom
 * ContentPickerModal - Basic Rendering Tests (5 tests)
 */
import { render, screen, waitFor, cleanup } from "@testing-library/svelte";
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

describe("ContentPickerModal - Rendering", () => {
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

  it("renders modal when show is true", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        expect(screen.getByRole("dialog")).toBeInTheDocument();
      },
      { timeout: 1000 }
    );
  });

  it("does not render when show is false", () => {
    renderResult = render(ContentPickerModal, {
      show: false,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
  });

  it("displays title", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        expect(screen.getByText("Pilih Konten Edukasi")).toBeInTheDocument();
      },
      { timeout: 1000 }
    );
  });

  it("fetches content on mount", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        expect(fetchAllContent).toHaveBeenCalledWith(null, "all");
      },
      { timeout: 1000 }
    );
  });

  it("fetches categories on mount", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        expect(fetchCategories).toHaveBeenCalled();
      },
      { timeout: 1000 }
    );
  });
});
