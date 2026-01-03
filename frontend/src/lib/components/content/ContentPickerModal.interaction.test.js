/**
 * @vitest-environment happy-dom
 * ContentPickerModal - Search & Tab Interaction Tests (6 tests)
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

describe("ContentPickerModal - Search & Tabs", () => {
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

  // Search tests
  it("displays search input", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        expect(
          screen.getByPlaceholderText("Cari artikel atau video...")
        ).toBeInTheDocument();
      },
      { timeout: 1000 }
    );
  });

  it("shows clear button when search query is not empty", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    const searchInput = await screen.findByPlaceholderText(
      "Cari artikel atau video..."
    );
    await fireEvent.input(searchInput, { target: { value: "test" } });

    const clearButton = screen.getByLabelText("Hapus pencarian");
    expect(clearButton).toBeInTheDocument();
  });

  it("clears search when clear button is clicked", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    const searchInput = await screen.findByPlaceholderText(
      "Cari artikel atau video..."
    );
    await fireEvent.input(searchInput, { target: { value: "test" } });

    const clearButton = screen.getByLabelText("Hapus pencarian");
    await fireEvent.click(clearButton);

    expect(searchInput).toHaveValue("");
  });

  // Tab tests
  it("displays three tabs", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        expect(screen.getByRole("tablist")).toBeInTheDocument();
        expect(screen.getByRole("tab", { name: "Semua" })).toBeInTheDocument();
        expect(
          screen.getByRole("tab", { name: "Artikel" })
        ).toBeInTheDocument();
        expect(screen.getByRole("tab", { name: "Video" })).toBeInTheDocument();
      },
      { timeout: 1000 }
    );
  });

  it("highlights active tab", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        const activeTab = screen.getByRole("tab", { name: "Semua" });
        expect(activeTab).toHaveClass("bg-teal-50");
      },
      { timeout: 1000 }
    );
  });

  it("changes tab selection when clicked", async () => {
    renderResult = render(ContentPickerModal, {
      show: true,
      onClose: mockOnClose,
      onSelect: mockOnSelect,
    });

    await waitFor(
      () => {
        const articlesTab = screen.getByRole("tab", { name: "Artikel" });
        fireEvent.click(articlesTab);

        expect(articlesTab).toHaveAttribute("aria-selected", "true");
      },
      { timeout: 1000 }
    );
  });
});
