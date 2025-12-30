<script>
  import { t } from 'svelte-i18n';
  import { getFailedDeliveries, exportFailedDeliveries } from '$lib/utils/api.js';

  /**
   * FailedDeliveriesView - View for listing and filtering failed deliveries
   */

  /** @type {Object} */
  let { token = null } = $props();

  // Filter options
  const filterOptions = [
    { value: '', label: 'analytics.failed_deliveries.all_reasons' },
    { value: 'invalid_phone', label: 'analytics.failed_deliveries.invalid_phone' },
    { value: 'gowa_timeout', label: 'analytics.failed_deliveries.gowa_timeout' },
    { value: 'message_rejected', label: 'analytics.failed_deliveries.message_rejected' },
    { value: 'other', label: 'analytics.failed_deliveries.other_reason' }
  ];

  // State
  let status = $state('loading');
  let error = $state(null);
  let data = $state(null);
  let selectedFilter = $state('');
  let expandedItemId = $state(null);
  let expandedItemDetails = $state(null);
  let loadingDetails = $state(false);
  let exporting = $state(false);

  // Load failed deliveries
  async function loadFailedDeliveries() {
    if (!token) return;

    status = 'loading';
    error = null;

    try {
      data = await getFailedDeliveries(token, {
        page: 1,
        limit: 20,
        reason: selectedFilter
      });
      status = 'success';
    } catch (e) {
      error = e.message;
      status = 'error';
    }
  }

  // Reload when filter changes
  $effect(() => {
    if (selectedFilter !== undefined && token) {
      loadFailedDeliveries();
    }
  });

  // Export CSV
  async function handleExport() {
    if (!token || exporting) return;

    exporting = true;
    try {
      const blob = await exportFailedDeliveries(token, { reason: selectedFilter });

      // Download the file
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `failed-deliveries-${new Date().toISOString().split('T')[0]}.csv`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    } catch (e) {
      error = e.message;
    } finally {
      exporting = false;
    }
  }

  // Toggle expanded item
  function toggleExpand(reminderId) {
    expandedItemId = expandedItemId === reminderId ? null : reminderId;
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
</script>

<div class="failed-deliveries-view">
  <!-- Header -->
  <div class="mb-6">
    <h1 class="text-2xl font-bold text-gray-900 mb-2">
      {$t('analytics.failed_deliveries.list_title', { default: 'Daftar Pengiriman Gagal' })}
    </h1>

    <!-- Filters and Actions -->
    <div class="flex flex-wrap items-center gap-4 mt-4">
      <!-- Filter by reason -->
      <div class="flex items-center gap-2">
        <label class="text-sm text-gray-600">
          {$t('analytics.failed_deliveries.filter_by_reason', { default: 'Filter berdasarkan alasan' })}:
        </label>
        <select
          class="px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          bind:value={selectedFilter}
        >
          {#each filterOptions as option}
            <option value={option.value}>
              {$t(option.label, { default: option.value || 'Semua' })}
            </option>
          {/each}
        </select>
      </div>

      <!-- Export button -->
      <button
        class="px-4 py-2 bg-green-600 text-white rounded-lg text-sm font-medium hover:bg-green-700 transition-colors flex items-center gap-2"
        onclick={handleExport}
        disabled={exporting || status === 'loading'}
      >
        {#if exporting}
          <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
          <span>{$t('analytics.failed_deliveries.exporting', { default: 'Mengexport...' })}</span>
        {:else}
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          <span>{$t('analytics.failed_deliveries.export_csv', { default: 'Export CSV' })}</span>
        {/if}
      </button>
    </div>
  </div>

  <!-- Content -->
  {#if status === 'loading'}
    <div class="flex items-center justify-center py-16">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600"></div>
      <span class="ml-3 text-gray-600">{$t('analytics.failed_deliveries.loading', { default: 'Memuat data...' })}</span>
    </div>
  {:else if status === 'error'}
    <div class="bg-red-50 border border-red-200 rounded-lg p-6">
      <p class="text-red-700">{$t('analytics.failed_deliveries.error', { default: 'Gagal memuat data' })}: {error}</p>
      <button
        class="mt-3 text-sm text-red-600 hover:text-red-800 font-medium"
        onclick={loadFailedDeliveries}
      >
        {$t('analytics.retry', { default: 'Coba lagi' })}
      </button>
    </div>
  {:else if data}
    <!-- Filter counts -->
    {#if data.filter_counts}
      <div class="flex gap-2 mb-4">
        {#each filterOptions as option}
          {@const count = option.value === '' ? data.items?.length || 0 : data.filter_counts[option.value] || 0}
          <button
            class="px-3 py-1 rounded-full text-xs font-medium transition-colors
                   {selectedFilter === option.value
                     ? 'bg-blue-600 text-white'
                     : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
            onclick={() => selectedFilter = option.value}
          >
            {$t(option.label, { default: option.value || 'Semua' })} ({count})
          </button>
        {/each}
      </div>
    {/if}

    {#if data.items && data.items.length > 0}
      <!-- Table -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {$t('analytics.failed_deliveries.patient_name', { default: 'Nama Pasien' })}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {$t('analytics.failed_deliveries.reminder_title', { default: 'Judul Reminder' })}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {$t('analytics.failed_deliveries.failure_reason', { default: 'Alasan Gagal' })}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {$t('analytics.failed_deliveries.failure_timestamp', { default: 'Waktu Gagal' })}
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {$t('analytics.failed_deliveries.retry_count', { default: 'Jumlah Percobaan' })}
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                {$t('common.viewAll', { default: 'Aksi' })}
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            {#each data.items as item}
              <tr class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {item.patient_name_masked}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {item.reminder_title}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full {getReasonClass(item.failure_reason_code)}">
                    {item.failure_reason}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {formatDate(item.failure_timestamp)}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {item.retry_count}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button
                    class="text-blue-600 hover:text-blue-900"
                    onclick={() => toggleExpand(item.reminder_id)}
                  >
                    {$t('analytics.failed_deliveries.view_details', { default: 'Lihat Detail' })}
                  </button>
                </td>
              </tr>

              <!-- Expanded detail -->
              {#if expandedItemId === item.reminder_id}
                <tr class="bg-gray-50">
                  <td colspan="6" class="px-6 py-4">
                    <div class="bg-white rounded-lg border border-gray-200 p-4">
                      <h4 class="font-medium text-gray-900 mb-3">
                        {$t('analytics.failed_deliveries.error_message', { default: 'Pesan Error' })}
                      </h4>
                      <p class="text-sm text-gray-600 mb-4 font-mono bg-gray-100 p-3 rounded">
                        {item.delivery_error_message || '-'}
                      </p>

                      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
                        <div>
                          <span class="text-gray-500">{$t('analytics.failed_deliveries.volunteer_name', { default: 'Nama Volunteer' })}:</span>
                          <span class="ml-2 text-gray-900">{item.volunteer_name || '-'}</span>
                        </div>
                        <div>
                          <span class="text-gray-500">{$t('analytics.failed_deliveries.phone_number', { default: 'Nomor Telepon' })}:</span>
                          <span class="ml-2 text-gray-900">{item.phone_masked || '-'}</span>
                        </div>
                        <div>
                          <span class="text-gray-500">{$t('analytics.failed_deliveries.retry_count', { default: 'Jumlah Percobaan' })}:</span>
                          <span class="ml-2 text-gray-900">{item.retry_count}</span>
                        </div>
                      </div>
                    </div>
                  </td>
                </tr>
              {/if}
            {/each}
          </tbody>
        </table>
      </div>

      <!-- Pagination info -->
      {#if data.pagination}
        <div class="mt-4 text-sm text-gray-600 text-center">
          {$t('analytics.failed_deliveries.page_of', {
            default: `Halaman ${data.pagination.page} dari ${data.pagination.total_pages}`,
            values: { page: data.pagination.page, total: data.pagination.total_pages }
          })}
          <span class="mx-2">|</span>
          <span>{$t('common.total', { default: 'Total' })}: {data.pagination.total}</span>
        </div>
      {/if}
    {:else}
      <!-- Empty state -->
      <div class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">
          {$t('analytics.failed_deliveries.no_failed_deliveries', { default: 'Tidak ada pengiriman gagal' })}
        </h3>
        <p class="mt-1 text-sm text-gray-500">
          {$t('delivery.filter.empty.failed', { default: 'Tidak ada reminder yang gagal. Semua reminder terkirim dengan baik!' })}
        </p>
      </div>
    {/if}
  {/if}
</div>
