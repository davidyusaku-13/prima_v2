<script>
  import { t } from 'svelte-i18n';

  export let show = false;
  export let message = '';
  export let onClose = () => {};
  export let onConfirm = () => {};
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
      onclick={onClose}
      onkeydown={(e) => e.key === 'Escape' && onClose()}
      role="button"
      tabindex="0"
      aria-label="Close modal"
    ></div>
    <div
      class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-sm p-4 sm:p-6"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <div class="flex items-center gap-3 sm:gap-4 mb-4 sm:mb-6">
        <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center flex-shrink-0">
          <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <p class="text-slate-700">{message}</p>
      </div>
      <div class="flex flex-col sm:flex-row gap-3">
        <button
          onclick={onClose}
          class="flex-1 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1"
        >
          {$t('common.cancel')}
        </button>
        <button
          onclick={onConfirm}
          class="flex-1 px-4 py-2.5 bg-red-600 text-white font-medium rounded-xl hover:bg-red-700 transition-colors duration-200 order-1 sm:order-2"
        >
          {$t('common.delete')}
        </button>
      </div>
    </div>
  </div>
{/if}
