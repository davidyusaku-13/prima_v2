<script>
  import { t } from 'svelte-i18n';

  export let show = false;
  export let user = null;
  export let locale = 'en';
  export let onSetLocale = () => {};
  export let onLogout = () => {};
  export let onClose = () => {};
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-end lg:items-center justify-center lg:p-4">
    <div
      class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
      onclick={onClose}
      onkeydown={(e) => e.key === 'Escape' && onClose()}
      role="button"
      tabindex="0"
      aria-label="Close modal"
    ></div>
    <div
      class="relative bg-white rounded-t-2xl lg:rounded-2xl shadow-xl w-full lg:max-w-sm p-4 pb-8 lg:p-6"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <!-- Drag handle -->
      <div class="lg:hidden flex justify-center mb-4">
        <div class="w-12 h-1.5 bg-slate-200 rounded-full"></div>
      </div>

      <!-- User info -->
      <div class="flex items-center gap-3 mb-6">
        <div class="w-12 h-12 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-lg">
          {user?.fullName?.charAt(0)?.toUpperCase() || user?.username?.charAt(0)?.toUpperCase() || 'U'}
        </div>
        <div class="flex-1 min-w-0">
          <p class="font-medium text-slate-900 truncate">{user?.fullName || user?.username}</p>
          <p class="text-sm text-slate-500 truncate">@{user?.username}</p>
        </div>
        <button onclick={onLogout} class="p-2 text-slate-400 hover:text-red-600 transition-colors" title={$t('auth.logout')}>
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
        </button>
      </div>

      <!-- Language switcher -->
      <div class="border-t border-slate-100 pt-4">
        <p class="text-sm font-medium text-slate-700 mb-2">{$t('common.language')}</p>
        <div class="flex gap-2">
          <button
            onclick={() => onSetLocale('en')}
            class="flex-1 py-2.5 rounded-xl font-medium transition-colors duration-200 {locale === 'en' ? 'bg-teal-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
          >
            English
          </button>
          <button
            onclick={() => onSetLocale('id')}
            class="flex-1 py-2.5 rounded-xl font-medium transition-colors duration-200 {locale === 'id' ? 'bg-teal-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
          >
            Bahasa Indonesia
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}
