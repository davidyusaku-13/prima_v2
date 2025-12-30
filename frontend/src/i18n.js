import { register, init, getLocaleFromNavigator, locale } from "svelte-i18n";

const DEFAULT_LOCALE = "en";

// Load translations from JSON files
import en from "./locales/en.json";
import id from "./locales/id.json";

register("en", () => Promise.resolve(en));
register("id", () => Promise.resolve(id));

// Get initial locale synchronously
const getInitialLocale = () => {
  if (typeof window !== 'undefined' && window.localStorage) {
    const savedLocale = localStorage.getItem("locale");
    if (savedLocale) return savedLocale;
  }
  return getLocaleFromNavigator() || DEFAULT_LOCALE;
};

// Set locale synchronously BEFORE init
const initialLocale = getInitialLocale();
locale.set(initialLocale);

// Initialize i18n with the pre-set locale
init({
  fallbackLocale: DEFAULT_LOCALE,
  initialLocale: initialLocale,
});

export { locale };
