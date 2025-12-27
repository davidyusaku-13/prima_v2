<script>
  import { t } from 'svelte-i18n';

  export let stats = {
    totalArticles: 0,
    totalVideos: 0,
    totalViews: 0
  };
</script>

<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
  <!-- Total Articles -->
  <div class="bg-white rounded-2xl border border-slate-200 p-5">
    <div class="flex items-center gap-3">
      <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
        <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
        </svg>
      </div>
      <div>
        <p class="text-2xl font-bold text-slate-900">{stats.totalArticles || 0}</p>
        <p class="text-xs text-slate-500">{$t('cms.stats.totalArticles')}</p>
      </div>
    </div>
  </div>

  <!-- Total Videos -->
  <div class="bg-white rounded-2xl border border-slate-200 p-5">
    <div class="flex items-center gap-3">
      <div class="w-12 h-12 bg-red-100 rounded-xl flex items-center justify-center">
        <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
      </div>
      <div>
        <p class="text-2xl font-bold text-slate-900">{stats.totalVideos || 0}</p>
        <p class="text-xs text-slate-500">{$t('cms.stats.totalVideos')}</p>
      </div>
    </div>
  </div>

  <!-- Total Views -->
  <div class="bg-white rounded-2xl border border-slate-200 p-5">
    <div class="flex items-center gap-3">
      <div class="w-12 h-12 bg-amber-100 rounded-xl flex items-center justify-center">
        <svg class="w-6 h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
      </div>
      <div>
        <p class="text-2xl font-bold text-slate-900">{stats.totalViews?.toLocaleString() || 0}</p>
        <p class="text-xs text-slate-500">{$t('cms.stats.totalViews')}</p>
      </div>
    </div>
  </div>

  <!-- Popular Content -->
  <div class="bg-white rounded-2xl border border-slate-200 p-5">
    <div class="flex items-center gap-3">
      <div class="w-12 h-12 bg-teal-100 rounded-xl flex items-center justify-center">
        <svg class="w-6 h-6 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
        </svg>
      </div>
      <div>
        <p class="text-2xl font-bold text-slate-900">{(stats.popularViews || 0).toLocaleString()}</p>
        <p class="text-xs text-slate-500">{$t('cms.stats.popularContent')}</p>
      </div>
    </div>
  </div>
</div>
