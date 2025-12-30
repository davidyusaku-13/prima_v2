<script>
  import { deliveryStore } from '$lib/stores/delivery.svelte.js';
  import { _ as t } from 'svelte-i18n';

  /**
   * FailedReminderBadge - Displays count of failed reminders
   *
   * CRITICAL: This is Vite + Svelte 5, NOT SvelteKit!
   * - Uses Svelte 5 runes ($derived)
   * - No legacy reactive statements ($:)
   * - No SvelteKit imports
   */

  // Reactive failed count from store
  let failedCount = $derived(deliveryStore.failedCount);

  function handleClick() {
    // Navigate to failed reminders view by dispatching custom event
    // PatientsView will listen to this event and filter by failed status
    window.dispatchEvent(new CustomEvent('show-failed-reminders'));
  }
</script>

{#if failedCount > 0}
  <button
    onclick={handleClick}
    class="relative inline-flex items-center px-3 py-1.5 rounded-full bg-red-600 text-white text-sm font-medium hover:bg-red-700 transition-colors focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
    aria-label={$t('delivery.failedRemindersCount', { values: { count: failedCount } })}
  >
    <!-- Alert icon -->
    <svg class="w-4 h-4 mr-1.5" fill="currentColor" viewBox="0 0 20 20" aria-hidden="true">
      <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
    </svg>

    <!-- Count -->
    <span>{failedCount}</span>
  </button>
{/if}
