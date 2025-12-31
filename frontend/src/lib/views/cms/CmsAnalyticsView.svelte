<script>
  import { t } from 'svelte-i18n';
  import ContentAnalyticsWidget from '../../components/analytics/ContentAnalyticsWidget.svelte';
  import DeliveryAnalyticsWidget from '../../components/analytics/DeliveryAnalyticsWidget.svelte';
  import SystemHealthWidget from '../../components/health/SystemHealthWidget.svelte';

  /**
   * CmsAnalyticsView - Full analytics page for admin
   * Shows content attachment statistics, delivery statistics, and system health
   */

  /** @type {Object} */
  let { token = null } = $props();
</script>

<div class="analytics-page">
  <div class="mb-6">
    <h1 class="text-2xl font-bold text-gray-900">
      {$t('analytics.title', { default: 'Analitik' })}
    </h1>
    <p class="text-gray-600 mt-1">
      {$t('analytics.sectionTitle', { default: 'Statistik Sistem' })}
    </p>
  </div>

  <!-- System Health Widget - Top Right -->
  <div class="max-w-md ml-auto mb-6">
    <SystemHealthWidget {token} />
  </div>

  <!-- Delivery Statistics Section -->
  <div class="mb-8">
    <h2 class="text-lg font-semibold text-gray-800 mb-4">
      {$t('analytics.deliverySectionTitle', { default: 'Statistik Pengiriman' })}
    </h2>
    <div class="max-w-4xl">
      <DeliveryAnalyticsWidget {token} />
    </div>
  </div>

  <!-- Content Analytics Section -->
  <div class="max-w-2xl">
    <h2 class="text-lg font-semibold text-gray-800 mb-4">
      {$t('analytics.contentSectionTitle', { default: 'Statistik Lampiran Konten' })}
    </h2>
    <ContentAnalyticsWidget {token} />
  </div>
</div>
