import { JSDOM } from 'jsdom';
import type { Environment } from 'vitest';

export default <Environment>{
  name: 'jsdom-custom',
  setup() {
    // Create jsdom instance
    const dom = new JSDOM('<!DOCTYPE html><html><body><div id="root"></div></body></html>', {
      url: 'http://localhost',
      pretendToBeVisual: true,
      runScripts: 'dangerously',
      resources: 'usable'
    });

    // Set up global objects
    globalThis.window = dom.window;
    globalThis.document = dom.window.document;
    globalThis.navigator = dom.window.navigator;
    globalThis.Element = dom.window.Element;
    globalThis.HTMLElement = dom.window.HTMLElement;
    globalThis.HTMLDivElement = dom.window.HTMLDivElement;
    globalThis.HTMLButtonElement = dom.window.HTMLButtonElement;
    globalThis.Node = dom.window.Node;
    globalThis.Event = dom.window.Event;
    globalThis.KeyboardEvent = dom.window.KeyboardEvent;
    globalThis.MouseEvent = dom.window.MouseEvent;
    globalThis.CustomEvent = dom.window.CustomEvent;
    globalThis.InputEvent = dom.window.InputEvent;
    globalThis.Text = dom.window.Text;
    globalThis.Comment = dom.window.Comment;
    globalThis.DocumentFragment = dom.window.DocumentFragment;
    globalThis.Request = dom.window.Request;
    globalThis.Response = dom.window.Response;
    globalThis.Headers = dom.window.Headers;
    globalThis.AbortController = dom.window.AbortController;
    globalThis.AbortSignal = dom.window.AbortSignal;

    // Element.prototype methods
    Element.prototype.getBoundingClientRect = () => ({
      left: 0, top: 0, right: 0, bottom: 0, width: 0, height: 0, x: 0, y: 0
    });
    Element.prototype.scrollIntoView = () => {};
    Element.prototype.click = () => {};
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

    return {
      teardown() {
        dom.window.close();
      }
    };
  }
};
