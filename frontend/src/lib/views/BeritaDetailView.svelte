<script>
  import { t, locale } from 'svelte-i18n';
  import * as api from '$lib/utils/api.js';

  export let articleId;
  export let onBack = () => {};

  let article = null;
  let loading = true;
  let error = null;

  function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString($locale || 'en-US', {
      weekday: 'long',
      day: 'numeric',
      month: 'long',
      year: 'numeric'
    });
  }

  function getCategoryLabel(category) {
    const categories = {
      latest: $t('articleCategories.latest'),
      policy: $t('articleCategories.policy'),
      research: $t('articleCategories.research'),
      outbreak: $t('articleCategories.outbreak'),
      lifestyle: $t('articleCategories.lifestyle'),
      local: $t('articleCategories.local')
    };
    return categories[category] || category || $t('articleCategories.latest');
  }

  function getCategoryColor(category) {
    const colors = {
      latest: 'bg-red-100 text-red-700',
      policy: 'bg-blue-100 text-blue-700',
      research: 'bg-purple-100 text-purple-700',
      outbreak: 'bg-orange-100 text-orange-700',
      lifestyle: 'bg-green-100 text-green-700',
      local: 'bg-amber-100 text-amber-700'
    };
    return colors[category] || 'bg-slate-100 text-slate-700';
  }

  async function loadArticle() {
    loading = true;
    error = null;
    try {
      article = await api.fetchArticle(null, articleId);
    } catch (e) {
      console.error($t('berita.errorLoading'), e);
      error = e.message || $t('berita.errorLoading');
    } finally {
      loading = false;
    }
  }

  loadArticle();
</script>

<div class="max-w-4xl mx-auto">
  <!-- Back Button -->
  <button
    onclick={onBack}
    class="mb-6 flex items-center gap-2 text-slate-600 hover:text-teal-600 transition-colors"
  >
    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
    </svg>
    <span class="text-sm font-medium">{$t('common.close')}</span>
  </button>

  <!-- Loading -->
  {#if loading}
    <div class="bg-white rounded-2xl border border-slate-200 overflow-hidden animate-pulse">
      <div class="aspect-[21/9] bg-slate-200"></div>
      <div class="p-8 space-y-4">
        <div class="h-8 bg-slate-200 rounded w-3/4"></div>
        <div class="h-4 bg-slate-200 rounded w-1/4"></div>
        <div class="h-4 bg-slate-200 rounded w-full"></div>
        <div class="h-4 bg-slate-200 rounded w-full"></div>
        <div class="h-4 bg-slate-200 rounded w-2/3"></div>
      </div>
    </div>
  <!-- Error -->
  {:else if error}
    <div class="bg-white rounded-2xl border border-slate-200 p-12 text-center">
      <svg class="w-16 h-16 mx-auto text-red-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="mt-4 text-lg font-medium text-slate-900">{$t('berita.errorLoading')}</h3>
      <p class="mt-2 text-slate-500">{error}</p>
      <button
        onclick={loadArticle}
        class="mt-4 px-4 py-2 bg-teal-600 text-white rounded-lg hover:bg-teal-700 transition-colors"
      >
        {$t('common.refresh')}
      </button>
    </div>
  <!-- Article Content -->
  {:else}
    <article class="bg-white rounded-2xl border border-slate-200 overflow-hidden">
      <!-- Hero Image -->
      {#if article.heroImage}
        <div class="aspect-[21/9] bg-slate-100 overflow-hidden">
          <img
            src={article.heroImage}
            alt={article.title}
            class="w-full h-full object-cover"
          />
        </div>
      {/if}

      <div class="p-6 md:p-8 lg:p-10">
        <!-- Category & Meta -->
        <div class="flex flex-wrap items-center gap-3 mb-4">
          {#if article.category}
            <span class="px-3 py-1 text-sm font-medium rounded-full {getCategoryColor(article.category)}">
              {getCategoryLabel(article.category)}
            </span>
          {/if}
          <span class="text-sm text-slate-400">
            {$t('berita.publishedOn')} {formatDate(article.createdAt || article.publishedAt)}
          </span>
        </div>

        <!-- Title -->
        <h1 class="text-2xl md:text-3xl lg:text-4xl font-bold text-slate-900 leading-tight">
          {article.title}
        </h1>

        <!-- Excerpt -->
        {#if article.excerpt}
          <p class="mt-4 text-lg text-slate-600 leading-relaxed">
            {article.excerpt}
          </p>
        {/if}

        <!-- Divider -->
        <hr class="my-8 border-slate-200" />

        <!-- Content -->
        <div class="prose prose-slate max-w-none">
          {#if article.content}
            {@html article.content}
          {:else}
            <p class="text-slate-500 italic">{$t('berita.noContent')}</p>
          {/if}
        </div>

        <!-- Author/Source -->
        {#if article.author || article.source}
          <div class="mt-8 pt-6 border-t border-slate-100">
            <p class="text-sm text-slate-500">
              {#if article.author}
                {$t('common.byAuthor', { author: article.author })}
              {/if}
              {#if article.source}
                <span class="mx-2">|</span>
                {$t('common.source', { source: article.source })}
              {/if}
            </p>
          </div>
        {/if}
      </div>
    </article>
  {/if}
</div>
