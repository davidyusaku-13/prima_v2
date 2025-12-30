<script>
  import { t } from 'svelte-i18n';
  import { getDeliveryAnalytics } from '$lib/utils/api.js';

  /**
   * DeliveryAnalyticsWidget - Shows delivery statistics dashboard
   * Displays summary cards and status breakdown for reminder delivery
   */

  /** @type {Object} */
  let { token = null } = $props();

  // Period options
  const periods = [
    { value: 'today', label: 'analytics.period.today' },
    { value: '7d', label: 'analytics.period.7days' },
    { value: '30d', label: 'analytics.period.30days' },
    { value: 'all', label: 'analytics.period.allTime' }
  ];

  // State
  let selectedPeriod = $state('7d');
  let status = $state('loading');
  let error = $state(null);
  let data = $state(null);

  // Load analytics data
  async function loadAnalytics() {
    if (!token) return;

    status = 'loading';
    error = null;

    try {
      data = await getDeliveryAnalytics(token, selectedPeriod);
      status = 'success';
    } catch (e) {
      error = e.message;
      status = 'error';
    }
  }

  // Reload when period changes
  $effect(() => {
    if (selectedPeriod && token) {
      loadAnalytics();
    }
  });

  // Format percentage
  function formatPercent(value) {
    return value.toFixed(1) + '%';
  }

  // Status display names
  const statusNames = {
    pending: 'analytics.status.pending',
    scheduled: 'analytics.status.scheduled',
    queued: 'analytics.status.queued',
    sending: 'analytics.status.sending',
    retrying: 'analytics.status.retrying',
    sent: 'analytics.status.sent',
    delivered: 'analytics.status.delivered',
    read: 'analytics.status.read',
    failed: 'analytics.status.failed',
    expired: 'analytics.status.expired'
  };

  // Calculate total from breakdown for accurate percentage
  function getBreakdownTotal(breakdown) {
    if (!breakdown) return 0;
    return Object.values(breakdown).reduce((sum, count) => sum + count, 0);
  }

  // Status colors
  const statusColors = {
    pending: 'bg-gray-100 text-gray-800',
    scheduled: 'bg-yellow-100 text-yellow-800',
    queued: 'bg-blue-100 text-blue-800',
    sending: 'bg-indigo-100 text-indigo-800',
    retrying: 'bg-orange-100 text-orange-800',
    sent: 'bg-cyan-100 text-cyan-800',
    delivered: 'bg-green-100 text-green-800',
    read: 'bg-emerald-100 text-emerald-800',
    failed: 'bg-red-100 text-red-800',
    expired: 'bg-gray-100 text-gray-600'
  };

  // Navigate to failed deliveries view
  let onNavigate = $props();

  function navigateToFailedDeliveries() {
    if (onNavigate) {
      onNavigate('failed-deliveries');
    } else if (typeof window !== 'undefined') {
      // Fallback to direct navigation
      window.location.href = '/analytics/failed-deliveries';
    }
  }
</script>

<div class="delivery-analytics">
  <!-- Period Filter -->
  <div class="mb-6">
    <div class="flex gap-2">
      {#each periods as period}
        <button
          class="px-4 py-2 rounded-lg text-sm font-medium transition-colors
                 {selectedPeriod === period.value
                   ? 'bg-blue-600 text-white'
                   : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
          onclick={() => selectedPeriod = period.value}
        >
          {$t(period.label, { default: period.value })}
        </button>
      {/each}
    </div>
  </div>

  {#if status === 'loading'}
    <div class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span class="ml-3 text-gray-600">{$t('analytics.loading', { default: 'Memuat data...' })}</span>
    </div>
  {:else if status === 'error'}
    <div class="bg-red-50 border border-red-200 rounded-lg p-4">
      <p class="text-red-700">{$t('analytics.error', { default: 'Gagal memuat data analitik' })}: {error}</p>
      <button
        class="mt-2 text-sm text-red-600 hover:text-red-800"
        onclick={loadAnalytics}
      >
        {$t('analytics.retry', { default: 'Coba lagi' })}
      </button>
    </div>
  {:else if data}
    <!-- Summary Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
      <!-- Total Sent -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4">
        <p class="text-sm text-gray-600 mb-1">
          {$t('analytics.totalSent', { default: 'Total Terkirim' })}
        </p>
        <p class="text-2xl font-bold text-gray-900">{data.totalSent || 0}</p>
      </div>

      <!-- Success Rate -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4">
        <p class="text-sm text-gray-600 mb-1">
          {$t('analytics.successRate', { default: 'Tingkat Berhasil' })}
        </p>
        <p class="text-2xl font-bold text-green-600">{formatPercent(data.successRate)}</p>
      </div>

      <!-- Failed Last 7 Days -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4">
        <p class="text-sm text-gray-600 mb-1">
          {$t('analytics.failedLast7Days', { default: 'Gagal (7 Hari)' })}
        </p>
        <div class="flex items-center gap-2">
          <p class="text-2xl font-bold text-red-600">{data.failedLast7Days || 0}</p>
          {#if data.failedLast7Days > 0}
            <button
              class="text-xs text-blue-600 hover:text-blue-800 font-medium"
              onclick={navigateToFailedDeliveries}
            >
              {$t('analytics.failed_deliveries.view_details', { default: 'Lihat Detail' })}
            </button>
          {/if}
        </div>
      </div>

      <!-- Avg Delivery Time -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-4">
        <p class="text-sm text-gray-600 mb-1">
          {$t('analytics.avgDeliveryTime', { default: 'Waktu Rata-rata' })}
        </p>
        <p class="text-2xl font-bold text-blue-600">{data.avgDeliveryTime || '-'}</p>
      </div>
    </div>

    <!-- Status Breakdown -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">
        {$t('analytics.breakdownTitle', { default: 'Rincian Status Pengiriman' })}
      </h3>

      <div class="space-y-3">
        {#each Object.entries(data.breakdown || {}) as [statusName, count]}
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <span class="px-2 py-1 rounded text-xs font-medium {statusColors[statusName] || 'bg-gray-100 text-gray-800'}">
                {$t(statusNames[statusName], { default: statusName })}
              </span>
            </div>
            <div class="flex items-center gap-2">
              <span class="text-gray-900 font-medium">{count || 0}</span>
              {#if getBreakdownTotal(data.breakdown) > 0}
                <span class="text-gray-500 text-sm">
                  ({((count / getBreakdownTotal(data.breakdown)) * 100).toFixed(1)}%)
                </span>
              {/if}
            </div>
          </div>
        {/each}
      </div>
    </div>
  {/if}
</div>
