// Mock svelte-i18n module for testing
// This is used by Vite's resolve.alias to intercept svelte-i18n imports

import { vi } from 'vitest';

// Create a proper Svelte store for mocking i18n
const createI18nStore = (value) => {
  let subscribers = [];
  const subscribe = (fn) => {
    subscribers.push(fn);
    fn(value);
    return () => {
      subscribers = subscribers.filter(s => s !== fn);
    };
  };
  const set = (newValue) => {
    value = newValue;
    subscribers.forEach(fn => fn(value));
  };
  const update = (fn) => {
    set(fn(value));
  };
  return { subscribe, set, update };
};

// Translations for testing
const translations = {
  'content.picker.articles': 'Articles',
  'content.picker.videos': 'Videos',
  'content.preview.attach': 'Attach',
  'content.preview.selected': 'Selected',
  'berita.publishedOn': 'Published on',
  'common.close': 'Close',
  'content.disclaimer.text': 'Konten ini untuk tujuan edukasi kesehatan. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda.',
  'delivery.status.pending': 'Pending',
  'delivery.status.sending': 'Sending...',
  'delivery.status.sent': 'Sent',
  'delivery.status.delivered': 'Delivered',
  'delivery.status.read': 'Read',
  'delivery.status.failed': 'Failed',
  'delivery.retry': 'Retry',
  'common.status.pending': 'Tertunda',
  'reminder.status.scheduled': 'Dijadwalkan',
  'reminder.status.queued': 'Dalam Antrian',
  'reminder.status.sending': 'Mengirim...',
  'reminder.status.sent': 'Terkirim',
  'reminder.status.delivered': 'Diterima',
  'reminder.status.read': 'Dibaca',
  'reminder.status.failed': 'Gagal',
  'reminder.status.expired': 'Kedaluwarsa',
  'content.picker.title': 'Pilih Konten Edukasi',
  'content.picker.search': 'Cari artikel atau video...',
  'content.picker.tab.all': 'Semua',
  'content.picker.tab.articles': 'Artikel',
  'content.picker.tab.videos': 'Video',
  'content.picker.attach': 'Lampirkan',
  'content.picker.noResults': 'Tidak ada hasil',
  'content.picker.empty': 'Belum ada konten edukasi',
  'common.cancel': 'Batal'
};

// Create a mock t function that returns stores
const mockT = (key, options) => {
  const value = translations[key] || key;
  return createI18nStore(value);
};

export const t = mockT;
export const locale = { subscribe: vi.fn() };
export const defaultLocale = 'id';
export const locales = ['en', 'id'];
export const dictionary = {};
export const isLoading = { subscribe: vi.fn() };
export const load = vi.fn();
export const hydrate = vi.fn();
export const init = vi.fn();
export const getLocaleFromPath = vi.fn();
export const navigate = vi.fn();
export const addMessages = vi.fn();
export const setLocale = vi.fn();
export const setLocales = vi.fn();
