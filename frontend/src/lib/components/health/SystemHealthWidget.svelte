<script>
  import { t } from 'svelte-i18n';
  import { getHealthDetailed } from '$lib/utils/api.js';

  /**
   * SystemHealthWidget - Shows system health status for admins
   * Displays GOWA connectivity, circuit breaker status, and queue size
   */

  /** @type {Object} */
  let { token = null } = $props();

  // State
  let status = $state('idle');
  let error = $state(null);
  let data = $state(null);
  let expanded = $state(false);

  // Load health data
  async function loadHealth() {
    if (!token) {
      status = 'idle';
      return;
    }

    status = 'loading';
    error = null;

    try {
      data = await getHealthDetailed(token);
      status = 'success';
    } catch (e) {
      error = e.message;
      status = 'error';
    }
  }

  // Auto-refresh every 60 seconds (matches AC requirement)
  $effect(() => {
    loadHealth();
    const interval = setInterval(loadHealth, 60000);
    return () => clearInterval(interval);
  });

  // Toggle expanded view
  function toggleExpanded() {
    expanded = !expanded;
  }

  // Format time for display
  function formatTime(isoString) {
    if (!isoString) return '-';
    try {
      const date = new Date(isoString);
      return date.toLocaleTimeString('id-ID', {
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return '-';
    }
  }

  // Get cooldown minutes
  function getCooldownMinutes(seconds) {
    if (!seconds || seconds <= 0) return 0;
    return Math.ceil(seconds / 60);
  }

  // Status indicator class
  function getStatusClass(gowaConnected, queueTotal) {
    if (!gowaConnected) return 'bg-red-100 text-red-800 border-red-200';
    if (queueTotal > 10) return 'bg-amber-100 text-amber-800 border-amber-200';
    return 'bg-green-100 text-green-800 border-green-200';
  }

  // Status icon
  function getStatusIcon(gowaConnected, queueTotal) {
    if (!gowaConnected) return '💔';
    if (queueTotal > 10) return '⚠️';
    return '💚';
  }
</script>

<div class="system-health-widget">
  <!-- Collapsed View -->
  <button
    class="w-full text-left p-4 rounded-xl border transition-all duration-200 hover:shadow-md {getStatusClass(data?.gowa?.connected, data?.queue?.total)}"
    onclick={toggleExpanded}
    aria-expanded={expanded}
  >
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <span class="text-xl">{getStatusIcon(data?.gowa?.connected, data?.queue?.total)}</span>
        <div>
          <p class="font-semibold">
            {$t('system_health.widget_title', { default: 'Kesehatan Sistem' })}
          </p>
          <p class="text-sm opacity-80">
            {#if status === 'loading'}
              {$t('common.loading', { default: 'Memuat...' })}
            {:else if status === 'error'}
              {error || $t('common.error', { default: 'Kesalahan' })}
            {:else if data}
              {data.gowa?.connected
                ? $t('system_health.gowa_connected', { default: 'GOWA: Terhubung' })
                : $t('system_health.gowa_disconnected', { default: 'GOWA: Tidak Tersedia' })}
            {:else}
              -
            {/if}
          </p>
        </div>
      </div>
      <div class="flex items-center gap-4">
        {#if data && data.queue}
          <span class="text-sm font-medium">
            {#if data.queue.total > 10}
              {$t('system_health.queue_warning', { default: 'Antrian: {count} reminder', values: { count: data.queue.total } })}
            {:else}
              {$t('system_health.queue_size', { default: 'Antrian: {count} reminder', values: { count: data.queue.total } })}
            {/if}
          </span>
        {/if}
        <span class="text-lg transition-transform duration-200 {expanded ? 'rotate-180' : ''}">▼</span>
      </div>
    </div>

    <!-- Last checked time (always visible) -->
    {#if data && data.timestamp}
      <p class="text-xs mt-2 opacity-70">
        {$t('system_health.last_checked', { default: 'Terakhir dicek: {time}', values: { time: formatTime(data.timestamp) } })}
      </p>
    {/if}
  </button>

  <!-- Expanded View -->
  {#if expanded && data}
    <div class="mt-2 p-4 bg-white rounded-xl border border-gray-200 shadow-sm">
      <!-- GOWA Status Detail -->
      <div class="mb-4 pb-4 border-b border-gray-100">
        <h4 class="font-medium text-gray-900 mb-2">GOWA</h4>
        <div class="grid grid-cols-2 gap-3 text-sm">
          <div>
            <p class="text-gray-500">{$t('system_health.detail_endpoint', { default: 'Endpoint' })}</p>
            <p class="font-mono text-gray-700">{data.gowa?.endpoint || '-'}</p>
          </div>
          <div>
            <p class="text-gray-500">{$t('system_health.detail_last_success', { default: 'Terakhir sukses' })}</p>
            <p class="text-gray-700">{formatTime(data.gowa?.last_ping)}</p>
          </div>
        </div>
      </div>

      <!-- Circuit Breaker Status -->
      <div class="mb-4 pb-4 border-b border-gray-100">
        <h4 class="font-medium text-gray-900 mb-2">{$t('system_health.detail_circuit_state', { default: 'Status circuit breaker' })}</h4>
        <div class="flex items-center gap-2 mb-2">
          {#if data.circuit_breaker?.state === 'closed'}
            <span class="w-3 h-3 rounded-full bg-green-500"></span>
            <span class="text-sm text-green-700">
              {$t('system_health.circuit_breaker_closed', { default: 'Circuit breaker tertutup' })}
            </span>
          {:else if data.circuit_breaker?.state === 'open'}
            <span class="w-3 h-3 rounded-full bg-red-500"></span>
            <span class="text-sm text-red-700">
              {$t('system_health.circuit_breaker_open', { default: 'Circuit breaker aktif, reset dalam {minutes} menit', values: { minutes: getCooldownMinutes(data.circuit_breaker?.cooldown_remaining) } })}
            </span>
          {:else}
            <span class="w-3 h-3 rounded-full bg-yellow-500"></span>
            <span class="text-sm text-yellow-700">
              {$t('system_health.circuit_breaker_half_open', { default: 'Circuit breaker dalam mode pengujian' })}
            </span>
          {/if}
        </div>
        <p class="text-xs text-gray-500">
          {$t('system_health.failure_count', { default: 'Jumlah kegagalan: {count}', values: { count: data.circuit_breaker?.failure_count || 0 } })}
        </p>
      </div>

      <!-- Queue Breakdown -->
      <div>
        <h4 class="font-medium text-gray-900 mb-2">
          {$t('system_health.detail_queue_breakdown', { default: 'Rincian antrian' })}
        </h4>
        <div class="space-y-2 text-sm">
          <div class="flex justify-between">
            <span class="text-gray-600">{$t('system_health.detail_scheduled', { default: 'Terjadwal' })}</span>
            <span class="font-medium">{data.queue?.scheduled || 0}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">{$t('system_health.detail_retrying', { default: 'Retry' })}</span>
            <span class="font-medium">{data.queue?.retrying || 0}</span>
          </div>
          <div class="flex justify-between">
            <span class="text-gray-600">{$t('system_health.detail_quiet_hours', { default: 'Quiet hours' })}</span>
            <span class="font-medium">{data.queue?.quiet_hours || 0}</span>
          </div>
        </div>
      </div>

      <!-- Refresh Button -->
      <button
        class="mt-4 w-full py-2 px-4 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg text-sm font-medium transition-colors flex items-center justify-center gap-2"
        onclick={loadHealth}
        disabled={status === 'loading'}
      >
        <span class="{status === 'loading' ? 'animate-spin' : ''}">↻</span>
        {$t('system_health.refresh', { default: 'Perbarui' })}
      </button>
    </div>
  {:else if expanded && status === 'error'}
    <div class="mt-2 p-4 bg-red-50 rounded-xl border border-red-200">
      <p class="text-red-700 text-sm">{error}</p>
      <button
        class="mt-2 text-sm text-red-600 hover:text-red-800 font-medium"
        onclick={loadHealth}
      >
        {$t('analytics.retry', { default: 'Coba lagi' })}
      </button>
    </div>
  {/if}
</div>
