<script>
  import { t } from 'svelte-i18n';
  import { getContentAnalytics } from '../utils/api.js';

  /**
   * ContentAnalyticsWidget - Displays content attachment statistics for admins
   * Shows top 10 most-attached content and lists all articles/videos with counts
   */

  /**
   * @typedef {Object} Props
   * @property {string} [token] - JWT token for API authentication
   */

  /** @type {Props} */
  let { token = null } = $props();

  let loading = $state(true);
  let error = $state(null);
  let analyticsData = $state(null);

  // Fetch analytics data on mount
  $effect(() => {
    if (token) {
      fetchAnalytics();
    }
  });

  async function fetchAnalytics() {
    loading = true;
    error = null;
    try {
      const data = await getContentAnalytics(token);
      analyticsData = data;
    } catch (err) {
      error = err.message;
      console.error('Failed to fetch content analytics:', err);
    } finally {
      loading = false;
    }
  }

  /**
   * Format attachment count with proper pluralization
   * @param {number} count
   * @returns {string}
   */
  function formatAttachmentCount(count) {
    if (count === 0) {
      return $t('analytics.attachmentCount_zero', { default: 'Belum pernah dilampirkan' });
    }
    if (count === 1) {
      return $t('analytics.attachmentCount_one', { values: { count: 1 }, default: 'Dilampirkan 1 kali' });
    }
    return $t('analytics.attachmentCount_other', { values: { count }, default: `Dilampirkan ${count} kali` });
  }

  /**
   * Get type icon
   * @param {string} type
   * @returns {string}
   */
  function getTypeIcon(type) {
    return type === 'article' ? '📄' : '🎬';
  }

  /**
   * Get type label
   * @param {string} type
   * @returns {string}
   */
  function getTypeLabel(type) {
    return type === 'article'
      ? $t('analytics.articleLabel', { default: 'Artikel' })
      : $t('analytics.videoLabel', { default: 'Video' });
  }
</script>

<div class="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
  <!-- Header -->
  <div class="px-4 py-3 bg-gray-50 border-b border-gray-200">
    <h3 class="text-sm font-semibold text-gray-700">
      {$t('analytics.topContent', { default: 'Top 10 Konten Populer' })}
    </h3>
  </div>

  <!-- Content -->
  <div class="p-4">
    {#if loading}
      <div class="flex items-center justify-center py-8 text-gray-500">
        <svg class="animate-spin h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span>{$t('analytics.loading', { default: 'Memuat data analitik...' })}</span>
      </div>
    {:else if error}
      <div class="text-center py-4">
        <p class="text-red-600 text-sm mb-2">{$t('analytics.error', { default: 'Gagal memuat data analitik' })}</p>
        <button
          onclick={fetchAnalytics}
          class="text-sm text-teal-600 hover:text-teal-700 font-medium"
        >
          {$t('common.refresh', { default: 'Coba Lagi' })}
        </button>
      </div>
    {:else if !analyticsData?.topContent?.length && !analyticsData?.articles?.length && !analyticsData?.videos?.length}
      <div class="text-center py-4 text-gray-500 text-sm">
        {$t('analytics.empty', { default: 'Belum ada data lampiran' })}
      </div>
    {:else}
      <!-- Top 10 Content List -->
      <div class="space-y-2">
        <h4 class="text-xs font-medium text-gray-500 uppercase tracking-wider mb-3">
          {$t('analytics.topContent', { default: 'Top 10 Konten Populer' })}
        </h4>
        {#each analyticsData.topContent || [] as item, index (item.id)}
          <div class="flex items-center gap-3 p-2 rounded-lg hover:bg-gray-50 transition-colors">
            <!-- Rank -->
            <span class="flex-shrink-0 w-6 h-6 flex items-center justify-center rounded-full bg-teal-100 text-teal-700 text-xs font-bold">
              {index + 1}
            </span>
            <!-- Type Icon -->
            <span class="flex-shrink-0 text-lg" aria-hidden="true">
              {getTypeIcon(item.type)}
            </span>
            <!-- Content Info -->
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 truncate" title={item.title}>
                {item.title}
              </p>
              <p class="text-xs text-gray-500">
                {getTypeLabel(item.type)} • {formatAttachmentCount(item.attachmentCount)}
              </p>
            </div>
            <!-- Count Badge -->
            <span class="flex-shrink-0 px-2 py-1 rounded-full bg-teal-50 text-teal-700 text-xs font-medium">
              {item.attachmentCount}
            </span>
          </div>
        {/each}

        {#if !analyticsData.topContent?.length}
          <p class="text-sm text-gray-500 text-center py-4">
            {$t('analytics.empty', { default: 'Belum ada data lampiran' })}
          </p>
        {/if}
      </div>

      <!-- All Articles -->
      {#if analyticsData.articles?.length > 0}
        <div class="mt-6 pt-4 border-t border-gray-100">
          <h4 class="text-xs font-medium text-gray-500 uppercase tracking-wider mb-3">
            {$t('analytics.articles', { default: 'Artikel' })} ({analyticsData.articles.length})
          </h4>
          <div class="space-y-1 max-h-48 overflow-y-auto">
            {#each analyticsData.articles as article (article.id)}
              <div class="flex items-center justify-between p-2 rounded hover:bg-gray-50">
                <div class="flex items-center gap-2 min-w-0">
                  <span aria-hidden="true">📄</span>
                  <span class="text-sm text-gray-700 truncate">{article.title}</span>
                </div>
                <span class="text-xs text-gray-500 flex-shrink-0 ml-2">
                  {formatAttachmentCount(article.attachmentCount)}
                </span>
              </div>
            {/each}
          </div>
        </div>
      {/if}

      <!-- All Videos -->
      {#if analyticsData.videos?.length > 0}
        <div class="mt-4 pt-4 border-t border-gray-100">
          <h4 class="text-xs font-medium text-gray-500 uppercase tracking-wider mb-3">
            {$t('analytics.videos', { default: 'Video Edukasi' })} ({analyticsData.videos.length})
          </h4>
          <div class="space-y-1 max-h-48 overflow-y-auto">
            {#each analyticsData.videos as video (video.id)}
              <div class="flex items-center justify-between p-2 rounded hover:bg-gray-50">
                <div class="flex items-center gap-2 min-w-0">
                  <span aria-hidden="true">🎬</span>
                  <span class="text-sm text-gray-700 truncate">{video.title}</span>
                </div>
                <span class="text-xs text-gray-500 flex-shrink-0 ml-2">
                  {formatAttachmentCount(video.attachmentCount)}
                </span>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>
