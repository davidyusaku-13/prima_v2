<script>
  import { toastStore } from '$lib/stores/toast.svelte.js';
  import { _ as t } from 'svelte-i18n';

  /**
   * Toast notification component
   *
   * CRITICAL: This is Vite + Svelte 5, NOT SvelteKit!
   * - Uses Svelte 5 runes ($derived)
   * - No legacy reactive statements ($:)
   * - No SvelteKit imports
   */

  // Reactive toasts from store
  let toasts = $derived(toastStore.toasts);

  function handleClose(id) {
    toastStore.remove(id);
  }

  function handleAction(toast) {
    if (toast.action?.onClick) {
      toast.action.onClick();
    }
    toastStore.remove(toast.id);
  }
</script>

<!-- Toast container - fixed position top-right -->
<div class="fixed top-4 right-4 z-50 flex flex-col gap-2 max-w-sm pointer-events-none">
  {#each toasts as toast (toast.id)}
    <div
      class="bg-white dark:bg-gray-800 rounded-lg shadow-lg p-4 flex items-start gap-3 border-l-4 pointer-events-auto animate-slide-in"
      class:border-red-500={toast.type === 'error'}
      class:border-green-500={toast.type === 'success'}
      class:border-yellow-500={toast.type === 'warning'}
      class:border-blue-500={toast.type === 'info'}
      role="alert"
      aria-live="polite"
    >
      <!-- Icon based on type -->
      <div class="flex-shrink-0">
        {#if toast.type === 'error'}
          <svg class="w-5 h-5 text-red-500" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
          </svg>
        {:else if toast.type === 'success'}
          <svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
          </svg>
        {:else if toast.type === 'warning'}
          <svg class="w-5 h-5 text-yellow-500" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
            <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
          </svg>
        {:else}
          <svg class="w-5 h-5 text-blue-500" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
          </svg>
        {/if}
      </div>

      <!-- Message -->
      <div class="flex-1 text-sm text-gray-700 dark:text-gray-200">
        {toast.message}
      </div>

      <!-- Action button (if provided) -->
      {#if toast.action}
        <button
          onclick={() => handleAction(toast)}
          class="text-sm font-medium text-teal-600 hover:text-teal-700 dark:text-teal-400 focus:outline-none focus:underline"
        >
          {toast.action.label}
        </button>
      {/if}

      <!-- Close button -->
      <button
        onclick={() => handleClose(toast.id)}
        class="flex-shrink-0 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 focus:outline-none focus:ring-2 focus:ring-gray-400 rounded"
        aria-label={$t('common.close')}
      >
        <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
          <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/>
        </svg>
      </button>
    </div>
  {/each}
</div>

<style>
  @keyframes slide-in {
    from {
      transform: translateX(100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  .animate-slide-in {
    animation: slide-in 0.3s ease-out;
  }
</style>
