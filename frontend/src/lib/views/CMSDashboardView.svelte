<script>
  import { t } from 'svelte-i18n';
  import DashboardStats from '$lib/components/DashboardStats.svelte';
  import ActivityLog from '$lib/components/ActivityLog.svelte';
  import ArticleCard from '$lib/components/ArticleCard.svelte';
  import VideoCard from '$lib/components/VideoCard.svelte';
  import * as api from '$lib/utils/api.js';

  export let onNavigateToArticleEditor = () => {};
  export let onNavigateToVideoManager = () => {};
  export let onNavigateToArticle = () => {};

  let stats = { totalArticles: 0, totalVideos: 0, totalViews: 0, popularViews: 0 };
  let recentArticles = [];
  let recentVideos = [];
  let activities = [];
  let loading = true;

  async function loadDashboardData() {
    loading = true;
    try {
      const [dashboardStats, articles, videos, activityData] = await Promise.all([
        api.fetchDashboardStats(),
        api.fetchArticles(),
        api.fetchVideos(),
        api.fetchActivityLog()
      ]);
      // Transform backend stats to expected format
      stats = {
        totalArticles: dashboardStats.articles?.total || 0,
        totalVideos: dashboardStats.videos?.total || 0,
        totalViews: dashboardStats.total_views?.articles || 0,
        popularViews: dashboardStats.total_views?.articles || 0
      };
      recentArticles = articles.slice(0, 3);
      recentVideos = videos.slice(0, 3);
      activities = activityData.slice(0, 10);
    } catch (e) {
      console.error($t('common.errorLoading'), e);
    } finally {
      loading = false;
    }
  }

  function formatViews(views) {
    if (!views) return '0';
    if (views >= 1000000) return `${(views / 1000000).toFixed(1)}M`;
    if (views >= 1000) return `${(views / 1000).toFixed(1)}K`;
    return views.toString();
  }

  loadDashboardData();
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div>
      <h1 class="text-2xl font-bold text-slate-900">{$t('cms.dashboard')}</h1>
      <p class="text-slate-500 mt-1">{$t('cms.title')}</p>
    </div>
    <div class="flex gap-2">
      <button
        onclick={onNavigateToArticleEditor}
        class="px-4 py-2 bg-teal-600 text-white rounded-xl font-medium hover:bg-teal-700 transition-colors flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {$t('cms.addArticle')}
      </button>
      <button
        onclick={onNavigateToVideoManager}
        class="px-4 py-2 bg-blue-600 text-white rounded-xl font-medium hover:bg-blue-700 transition-colors flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {$t('cms.addVideo')}
      </button>
    </div>
  </div>

  <!-- Loading -->
  {#if loading}
    <div class="animate-pulse space-y-6">
      <div class="h-28 bg-slate-200 rounded-2xl"></div>
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 h-64 bg-slate-200 rounded-2xl"></div>
        <div class="h-64 bg-slate-200 rounded-2xl"></div>
      </div>
    </div>
  {:else}
    <!-- Stats -->
    <DashboardStats {stats} />

    <!-- Quick Actions & Activity -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Popular Content -->
      <div class="lg:col-span-2 bg-white rounded-2xl border border-slate-200">
        <div class="px-5 py-4 border-b border-slate-100">
          <h3 class="font-semibold text-slate-900">{$t('cms.stats.popularContent')}</h3>
        </div>
        <div class="p-5 grid grid-cols-1 sm:grid-cols-2 gap-4">
          {#if recentArticles.length > 0}
            {#each recentArticles as article}
              <ArticleCard
                {article}
                onClick={() => onNavigateToArticle(article)}
              />
            {/each}
          {:else}
            <div class="col-span-2 text-center py-8 text-slate-500">
              {$t('berita.noArticles')}
            </div>
          {/if}
        </div>
      </div>

      <!-- Activity Log -->
      <div class="lg:col-span-1">
        <ActivityLog {activities} />
      </div>
    </div>

    <!-- Recent Videos -->
    {#if recentVideos.length > 0}
      <div class="bg-white rounded-2xl border border-slate-200">
        <div class="px-5 py-4 border-b border-slate-100 flex items-center justify-between">
          <h3 class="font-semibold text-slate-900">{$t('cms.recentVideos')}</h3>
          <button
            onclick={onNavigateToVideoManager}
            class="text-sm text-teal-600 hover:text-teal-700 font-medium"
          >
            {$t('common.viewAll')}
          </button>
        </div>
        <div class="p-5 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each recentVideos as video}
            <VideoCard
              {video}
              onClick={() => {}}
            />
          {/each}
        </div>
      </div>
    {/if}
  {/if}
</div>
