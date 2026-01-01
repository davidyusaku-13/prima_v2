<script>
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';

  /**
   * @typedef {Object} Props
   * @property {boolean} show - Whether to show the modal
   * @property {Object} reminder - Reminder to cancel
   * @property {Function} onClose - Close handler
   * @property {Function} onConfirm - Confirm handler
   */

  /** @type {Props} */
  let {
    show = false,
    reminder = null,
    onClose = () => {},
    onConfirm = () => {}
  } = $props();

  let loading = $state(false);

  async function handleConfirm() {
    loading = true;
    try {
      await onConfirm();
    } finally {
      loading = false;
    }
  }

  // Handle ESC key globally when modal is open
  $effect(() => {
    if (!show) return;
    
    function handleKeyDown(e) {
      if (e.key === 'Escape' && !loading) {
        onClose();
      }
    }
    
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  });
</script>

{#if show && reminder}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
      onclick={onClose}
      role="button"
      tabindex="-1"
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
      <!-- Warning icon -->
      <div class="flex items-center gap-3 sm:gap-4 mb-4 sm:mb-6">
        <div class="w-12 h-12 bg-amber-100 rounded-full flex items-center justify-center flex-shrink-0">
          <svg class="w-6 h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <div>
          <p class="text-slate-700 font-medium">{$t('reminder.cancel.title') || 'Batalkan Reminder?'}</p>
          <p class="text-slate-500 text-sm mt-1">
            {$t('reminder.cancel.confirm') || 'Apakah Anda yakin ingin membatalkan reminder ini?'}
          </p>
        </div>
      </div>

      <!-- Reminder info -->
      <div class="bg-slate-50 rounded-lg p-3 mb-4">
        <p class="font-medium text-slate-900">{reminder.title}</p>
        {#if reminder.message || reminder.description}
          <p class="text-sm text-slate-500 mt-1 line-clamp-2">
            {reminder.message || reminder.description}
          </p>
        {/if}
        {#if reminder.scheduled_at || reminder.due_date}
          <div class="flex items-center gap-1.5 mt-2 text-xs text-slate-600">
            <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <span>{$t('reminder.cancel.scheduledTime') || 'Dijadwalkan'}: {new Date(reminder.scheduled_at || reminder.due_date).toLocaleString($locale || 'id', { dateStyle: 'medium', timeStyle: 'short' })}</span>
          </div>
        {/if}
      </div>

      <!-- Actions -->
      <div class="flex flex-col sm:flex-row gap-3">
        <button
          onclick={onClose}
          disabled={loading}
          class="flex-1 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1 disabled:opacity-50"
        >
          {$t('common.cancel')}
        </button>
        <button
          onclick={handleConfirm}
          disabled={loading}
          class="flex-1 px-4 py-2.5 bg-amber-600 text-white font-medium rounded-xl hover:bg-amber-700 transition-colors duration-200 order-1 sm:order-2 disabled:opacity-50 flex items-center justify-center gap-2"
        >
          {#if loading}
            <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{$t('reminder.cancel.cancelling') || 'Membatalkan...'}</span>
          {:else}
            <span>{$t('reminder.cancel.button') || 'Ya, Batalkan'}</span>
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
