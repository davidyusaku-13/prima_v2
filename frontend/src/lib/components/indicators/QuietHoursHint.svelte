<script>
  import { t } from 'svelte-i18n';

  /** @type {{ scheduledTime?: string | null }} */
  let { scheduledTime = null } = $props();

  /**
   * Format the scheduled time for display
   * @param {string} isoTime - ISO 8601 time string
   * @returns {string} Formatted time string
   */
  function formatScheduledTime(isoTime) {
    if (!isoTime) return '06:00';
    try {
      const date = new Date(isoTime);
      return date.toLocaleTimeString('id-ID', {
        hour: '2-digit',
        minute: '2-digit',
        hour12: false
      });
    } catch {
      return '06:00';
    }
  }
</script>

<div
  class="flex items-center gap-2 p-3 bg-amber-50 border border-amber-200 rounded-lg text-amber-800"
  role="alert"
  aria-label={$t('reminder.quietHours.hint')}
>
  <span class="text-lg" aria-hidden="true">⏰</span>
  <span class="text-sm font-medium">
    {$t('reminder.quietHours.hint', { default: `Reminder akan dikirim jam ${formatScheduledTime(scheduledTime)} besok` })}
  </span>
</div>
