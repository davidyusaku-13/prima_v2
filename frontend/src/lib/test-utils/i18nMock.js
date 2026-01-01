/**
 * Reusable i18n mock for testing Svelte 5 components with svelte-i18n
 *
 * Usage:
 * ```javascript
 * vi.mock('svelte-i18n', () => createI18nMock({
 *   'custom.key': 'Custom Translation'
 * }));
 * ```
 */
import { readable } from "svelte/store";

/**
 * Creates a mock for svelte-i18n with custom translations
 *
 * @param {Object} customTranslations - Optional custom translation key-value pairs
 * @returns {Object} Mock object compatible with svelte-i18n API
 */
export function createI18nMock(customTranslations = {}) {
  // Base translations used across many tests
  const baseTranslations = {
    "common.close": "Tutup",
    "common.cancel": "Batal",
    "common.save": "Simpan",
    "common.loading": "Memuat...",
    "common.refresh": "Segarkan",
    "common.errorLoading": "Gagal memuat data",
    "common.search": "Cari",
    "common.filter": "Filter",
    "common.noResults": "Tidak ada hasil",
  };

  // Merge custom translations with base (custom takes precedence)
  const translations = { ...baseTranslations, ...customTranslations };

  /**
   * Mock translation function
   * Supports both direct calls and parameter interpolation
   *
   * @param {string} key - Translation key
   * @param {Object} options - Optional parameters for interpolation
   * @returns {string} Translated text or key if not found
   */
  const mockT = (key, options) => {
    let text = translations[key] || key;

    // Simple parameter interpolation (basic support)
    if (options && options.values) {
      Object.entries(options.values).forEach(([param, value]) => {
        text = text.replace(`{${param}}`, value);
      });
    }

    return text;
  };

  // Create readable store with the translation function
  const tStore = readable(mockT);

  // Combine function and store subscription (required for Svelte reactivity)
  const t = Object.assign(mockT, { subscribe: tStore.subscribe });

  // Return full mock object matching svelte-i18n API
  return {
    t,
    _: mockT, // Alias for t
    locale: readable("id"),
    locales: readable(["id", "en"]),
    loading: readable(false),
    init: () => {},
    getLocaleFromNavigator: () => "id",
    addMessages: () => {},
    dictionary: readable(translations),
    waitLocale: () => Promise.resolve(),
  };
}
