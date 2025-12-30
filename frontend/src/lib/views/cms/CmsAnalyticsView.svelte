<script>
  import { t } from 'svelte-i18n';
  import ContentAnalyticsWidget from '../../components/analytics/ContentAnalyticsWidget.svelte';

  /**
   * CmsAnalyticsView - Full analytics page for admin
   * Shows content attachment statistics and popular content
   */

  /** @type {Object} */
  let { token = null } = $props();
</script>

<div class="analytics-page">
  <div class="mb-6">
    <h1 class="text-2xl font-bold text-gray-900">
      {$t('analytics.title', { default: 'Analitik Konten' })}
    </h1>
    <p class="text-gray-600 mt-1">
      {$t('analytics.sectionTitle', { default: 'Statistik Lampiran Konten' })}
    </p>
  </div>

  <div class="max-w-2xl">
    <ContentAnalyticsWidget {token} />
  </div>
</div>
