<script>
  import { t } from 'svelte-i18n';

  /**
   * FailedDeliveryCard - Expandable card for displaying failed delivery details
   */

  /** @type {Object} */
  let { item = null } = $props();

  // Format date
  function formatDate(dateString) {
    if (!dateString) return '-';
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString('id-ID', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return dateString;
    }
  }

  // Get reason display class
  function getReasonClass(reasonCode) {
    switch (reasonCode) {
      case 'invalid_phone': return 'bg-red-100 text-red-800';
      case 'gowa_timeout': return 'bg-yellow-100 text-yellow-800';
      case 'message_rejected': return 'bg-orange-100 text-orange-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  }
</script>

{#if item}
  <div class="failed-delivery-card bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
    <!-- Main row -->
    <div class="p-4">
      <div class="flex items-start justify-between">
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 mb-1">
            <span class="text-sm font-medium text-gray-900 truncate">
              {item.patient_name_masked}
            </span>
            <span class="px-2 py-0.5 text-xs font-semibold rounded-full {getReasonClass(item.failure_reason_code)}">
              {item.failure_reason}
            </span>
          </div>
          <p class="text-sm text-gray-600 truncate">
            {item.reminder_title}
          </p>
        </div>
        <div class="ml-4 text-right text-sm text-gray-500">
          <p>{formatDate(item.failure_timestamp)}</p>
          <p class="text-xs text-gray-400">
            {$t('analytics.failed_deliveries.retry_count', { default: 'Percobaan' })}: {item.retry_count}
          </p>
        </div>
      </div>
    </div>
  </div>
{/if}
