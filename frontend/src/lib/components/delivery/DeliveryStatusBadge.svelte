<script>
  import { t } from 'svelte-i18n';

  /**
   * @typedef {Object} Props
   * @property {string} status - Delivery status (pending, sending, sent, delivered, read, failed, retrying, scheduled, queued)
   * @property {number} [retryCount] - Current retry count
   * @property {number} [maxAttempts] - Maximum retry attempts
   * @property {Function} [onRetry] - Callback for retry button (only shown for failed status)
   * @property {boolean} [isRetrying] - Whether retry is in progress
   */

  /** @type {Props} */
  let { status = 'pending', retryCount = 0, maxAttempts = 3, onRetry = null, isRetrying = false } = $props();

  /**
   * Get status configuration with colors and labels per AC requirements
   * AC specifies exact icons and colors for each status
   * @param {string} statusValue
   * @returns {{ label: string, color: string, icon: string, iconColor: string }}
   */
  function getStatusConfig(statusValue) {
    const configs = {
      pending: {
        label: $t('reminder.status.pending', { default: 'Tertunda' }),
        color: 'bg-gray-100 border-gray-300',
        icon: '🕐',  // AC: clock icon in gray
        iconColor: 'text-gray-500'
      },
      scheduled: {
        label: $t('reminder.status.scheduled', { default: 'Dijadwalkan' }),
        color: 'bg-blue-100 border-blue-300',
        icon: '📅',
        iconColor: 'text-blue-600'
      },
      queued: {
        label: $t('reminder.status.queued', { default: 'Dalam Antrian' }),
        color: 'bg-yellow-100 border-yellow-300',
        icon: '🕐',
        iconColor: 'text-yellow-600'
      },
      sending: {
        label: $t('reminder.status.sending', { default: 'Mengirim...' }),
        color: 'bg-gray-100 border-gray-300 animate-pulse',
        icon: '⏳',  // AC: spinner icon in gray (using hourglass as spinner representation)
        iconColor: 'text-gray-500'
      },
      retrying: {
        label: $t('reminder.status.retrying', { values: { count: retryCount, max: maxAttempts }, default: `Mengirim ulang (${retryCount}/${maxAttempts})` }),
        color: 'bg-orange-100 border-orange-300',
        icon: '🔄',
        iconColor: 'text-orange-600'
      },
      sent: {
        label: $t('reminder.status.sent', { default: 'Terkirim' }),
        color: 'bg-gray-100 border-gray-300',
        icon: '✓',  // AC: single checkmark (✓) in gray #64748b
        iconColor: 'text-[#64748b]'
      },
      delivered: {
        label: $t('reminder.status.delivered', { default: 'Diterima' }),
        color: 'bg-emerald-100 border-emerald-300',
        icon: '✓✓',  // AC: double checkmarks (✓✓) in WhatsApp green #25D366
        iconColor: 'text-[#25D366]'
      },
      read: {
        label: $t('reminder.status.read', { default: 'Dibaca' }),
        color: 'bg-blue-100 border-blue-300',
        icon: '✓✓',  // AC: double checkmarks (✓✓) in WhatsApp blue #53bdeb
        iconColor: 'text-[#53bdeb]'
      },
      failed: {
        label: $t('reminder.status.failed', { default: 'Gagal' }),
        color: 'bg-red-100 border-red-300',
        icon: '✕',  // AC: X icon in red #dc2626
        iconColor: 'text-[#dc2626]'
      },
      expired: {
        label: $t('reminder.status.expired', { default: 'Kedaluwarsa' }),
        color: 'bg-gray-100 border-gray-400',
        icon: '⏰',
        iconColor: 'text-gray-600'
      }
    };

    return configs[statusValue] || configs.pending;
  }

  let statusConfig = $derived(getStatusConfig(status));

  function handleRetry() {
    if (onRetry && !isRetrying) {
      onRetry();
    }
  }

  function handleKeydown(event) {
    if ((event.key === 'Enter' || event.key === ' ') && onRetry && !isRetrying) {
      event.preventDefault();
      onRetry();
    }
  }
</script>

<span
  class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border {statusConfig.color}"
  role="status"
  aria-label={statusConfig.label}
  aria-live="polite"
  aria-atomic="true"
>
  <span class="{statusConfig.iconColor} font-bold" aria-hidden="true">{statusConfig.icon}</span>
  <span class="sr-only">{statusConfig.label}</span>
  <span class="text-gray-700">{statusConfig.label}</span>
  {#if status === 'retrying' && retryCount > 0}
    <span class="ml-0.5 opacity-75 text-gray-600">
      ({retryCount}/{maxAttempts})
    </span>
  {/if}
  {#if status === 'failed' && onRetry}
    <button
      onclick={handleRetry}
      onkeydown={handleKeydown}
      disabled={isRetrying}
      class="ml-1 px-2 py-0.5 text-xs bg-red-200 text-red-700 rounded hover:bg-red-300 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-1"
      aria-label={$t('reminder.status.retry', { default: 'Coba Lagi' })}
      tabindex="0"
    >
      {#if isRetrying}
        <svg class="w-3 h-3 animate-spin inline" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      {:else}
        {$t('reminder.status.retry', { default: 'Coba Lagi' })}
      {/if}
    </button>
  {/if}
</span>
