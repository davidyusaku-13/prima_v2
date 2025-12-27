import { register, init, getLocaleFromNavigator, locale } from 'svelte-i18n';

const DEFAULT_LOCALE = 'en';

// Load translations from JSON files
import en from './locales/en.json';
import id from './locales/id.json';

register('en', () => Promise.resolve(en));
register('id', () => Promise.resolve(id));

init({
  fallbackLocale: DEFAULT_LOCALE,
  initialLocale: localStorage.getItem('locale') || getLocaleFromNavigator() || DEFAULT_LOCALE
});

export { locale };
