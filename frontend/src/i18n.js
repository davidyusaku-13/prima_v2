import { register, init, getLocaleFromNavigator, locale } from "svelte-i18n";

const DEFAULT_LOCALE = "en";

// Load translations from JSON files
import en from "./locales/en.json";
import id from "./locales/id.json";

register("en", () => Promise.resolve(en));
register("id", () => Promise.resolve(id));

// Safe localStorage access - check if we're in browser context
const getInitialLocale = () => {
  if (typeof window !== 'undefined' && window.localStorage) {
    return localStorage.getItem("locale") || getLocaleFromNavigator() || DEFAULT_LOCALE;
  }
  return DEFAULT_LOCALE;
};

init({
  fallbackLocale: DEFAULT_LOCALE,
  initialLocale: getInitialLocale(),
});

export { locale };
