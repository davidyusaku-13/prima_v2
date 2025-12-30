// Setup happy-dom environment BEFORE any imports
import { Window } from 'happy-dom';

const window = new Window({
  url: 'http://localhost',
  pretendToBeVisual: true
});

const document = window.document;

// Set up global objects BEFORE any test code runs
globalThis.window = window;
globalThis.document = document;
globalThis.navigator = window.navigator;
globalThis.Element = window.Element;
globalThis.HTMLElement = window.HTMLElement;
globalThis.HTMLDivElement = window.HTMLDivElement;
globalThis.HTMLButtonElement = window.HTMLButtonElement;
globalThis.Node = window.Node;
globalThis.Event = window.Event;
globalThis.KeyboardEvent = window.KeyboardEvent;
globalThis.MouseEvent = window.MouseEvent;
globalThis.CustomEvent = window.CustomEvent;
globalThis.InputEvent = window.InputEvent;
globalThis.Text = window.Text;
globalThis.Comment = window.Comment;
globalThis.DocumentFragment = window.DocumentFragment;
globalThis.Request = window.Request;
globalThis.Response = window.Response;
globalThis.Headers = window.Headers;
globalThis.AbortController = window.AbortController;
globalThis.AbortSignal = window.AbortSignal;

import { beforeAll, afterAll, vi } from 'vitest';
import '@testing-library/jest-dom';

// Element.prototype methods
Element.prototype.getBoundingClientRect = () => ({
  left: 0,
  top: 0,
  right: 0,
  bottom: 0,
  width: 0,
  height: 0,
  x: 0,
  y: 0
});

Element.prototype.scrollIntoView = vi.fn();
Element.prototype.click = vi.fn();
Element.prototype.getClientRects = () => [];

// ResizeObserver mock
globalThis.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
};

// MutationObserver mock
globalThis.MutationObserver = class MutationObserver {
  observe() {}
  disconnect() {}
};

// window.matchMedia mock
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  configurable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn()
  }))
});

// window.scrollTo mock
window.scrollTo = vi.fn();

// requestAnimationFrame mock
window.requestAnimationFrame = vi.fn(callback => callback(1));
window.cancelAnimationFrame = vi.fn();

// localStorage mock
const localStorageMock = {
  store: {},
  getItem: vi.fn((key) => this.store[key] || null),
  setItem: vi.fn((key, value) => { this.store[key] = value; }),
  clear: vi.fn(() => { this.store = {}; }),
  removeItem: vi.fn((key) => { delete this.store[key]; })
};
Object.defineProperty(window, 'localStorage', {
  writable: true,
  configurable: true,
  value: localStorageMock
});

// sessionStorage mock
const sessionStorageMock = {
  store: {},
  getItem: vi.fn((key) => this.store[key] || null),
  setItem: vi.fn((key, value) => { this.store[key] = value; }),
  clear: vi.fn(() => { this.store = {}; }),
  removeItem: vi.fn((key) => { delete this.store[key]; })
};
Object.defineProperty(window, 'sessionStorage', {
  writable: true,
  configurable: true,
  value: sessionStorageMock
});

// Suppress console errors during tests
const originalError = console.error;
beforeAll(() => {
  console.error = (...args) => {
    // Suppress Svelte 5 hydration and lifecycle warnings
    if (
      args[0]?.includes?.('Hydration') ||
      args[0]?.includes?.('was passed to') ||
      args[0]?.includes?.('lifecycle_function_unavailable') ||
      args[0]?.includes?.('[vite]') ||
      args[0]?.includes?.('rune_outside_svelte') ||
      args[0]?.includes?.('Svelte error')
    ) {
      return;
    }
    originalError.call(console, ...args);
  };
});

afterAll(() => {
  console.error = originalError;
});
