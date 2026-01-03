/**
 * Shared test helpers for ContentPickerModal tests
 * Split into multiple files to prevent memory exhaustion
 */
import { vi } from "vitest";
import { readable } from "svelte/store";

// Mock localStorage
export const localStorageMock = {
  getItem: vi.fn(() => null),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};

// Setup localStorage mock on window
export function setupLocalStorageMock() {
  Object.defineProperty(window, "localStorage", { value: localStorageMock });
}

// i18n translations
export const translations = {
  "content.picker.title": "Pilih Konten Edukasi",
  "content.picker.search": "Cari artikel atau video...",
  "content.picker.tab.all": "Semua",
  "content.picker.tab.articles": "Artikel",
  "content.picker.tab.videos": "Video",
  "content.picker.articles": "Artikel",
  "content.picker.videos": "Video",
  "content.picker.attach": "Lampirkan",
  "content.picker.noResults": "Tidak ada hasil",
  "content.picker.empty": "Belum ada konten edukasi",
  "common.cancel": "Batal",
};

// Create i18n mock factory
export function createI18nMock() {
  const mockT = (key) => translations[key] || key;
  const tStore = readable(mockT);
  const t = Object.assign(mockT, { subscribe: tStore.subscribe });
  return {
    t,
    locale: readable("id"),
    locales: readable(["id", "en"]),
    loading: readable(false),
    init: vi.fn(),
    getLocaleFromNavigator: vi.fn(() => "id"),
    addMessages: vi.fn(),
    _: t,
  };
}

// Default mock data
export const defaultArticle = {
  id: "art-1",
  title: "Healthy Eating Guide",
  excerpt: "A guide",
};

export const defaultVideo = {
  id: "vid-1",
  title: "Morning Yoga",
  YouTubeID: "abc123",
  duration: "5:30",
};

export const defaultContentResponse = {
  articles: [defaultArticle],
  videos: [defaultVideo],
};
