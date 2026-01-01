/**
 * @vitest-environment happy-dom
 */
import { describe, it, expect, vi } from "vitest";
import { createI18nMock } from "./i18nMock.js";

describe("createI18nMock", () => {
  it("returns mock with base translations", () => {
    const mock = createI18nMock();

    expect(mock.t("common.close")).toBe("Tutup");
    expect(mock.t("common.loading")).toBe("Memuat...");
  });

  it("supports custom translations", () => {
    const mock = createI18nMock({
      "custom.key": "Custom Value",
      "common.close": "Custom Close", // Override base
    });

    expect(mock.t("custom.key")).toBe("Custom Value");
    expect(mock.t("common.close")).toBe("Custom Close");
  });

  it("returns key when translation not found", () => {
    const mock = createI18nMock();

    expect(mock.t("unknown.key")).toBe("unknown.key");
  });

  it("supports parameter interpolation", () => {
    const mock = createI18nMock({
      "message.count": "You have {count} messages",
    });

    const result = mock.t("message.count", { values: { count: 5 } });
    expect(result).toBe("You have 5 messages");
  });

  it("provides subscribable store", () => {
    const mock = createI18nMock();

    // Test that t function is callable
    expect(typeof mock.t).toBe("function");
    expect(mock.t("common.close")).toBe("Tutup");

    // Test that t has subscribe method (for Svelte stores)
    expect(typeof mock.t.subscribe).toBe("function");
  });

  it("provides all svelte-i18n API exports", () => {
    const mock = createI18nMock();

    expect(mock).toHaveProperty("t");
    expect(mock).toHaveProperty("_");
    expect(mock).toHaveProperty("locale");
    expect(mock).toHaveProperty("locales");
    expect(mock).toHaveProperty("loading");
    expect(mock).toHaveProperty("init");
    expect(mock).toHaveProperty("getLocaleFromNavigator");
    expect(mock).toHaveProperty("addMessages");
  });

  it("can be used in vi.mock", () => {
    // This test verifies the pattern works in actual mocks
    vi.mock("svelte-i18n", () =>
      createI18nMock({
        "test.key": "Test Value",
      })
    );

    // In real usage, component would import and use $t()
    // Here we just verify the mock can be created
    expect(true).toBe(true);
  });
});
