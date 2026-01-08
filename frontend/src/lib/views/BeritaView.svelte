<script>
  import { t } from 'svelte-i18n';
  import ArticleCard from '$lib/components/ArticleCard.svelte';
  import * as api from '$lib/utils/api.js';

  export let onNavigateToArticle = () => {};

  let articles = [];
  let categories = ['latest', 'policy', 'research', 'outbreak', 'lifestyle', 'local'];
  let selectedCategory = null;
  let loading = true;
  let searchQuery = '';

  async function loadArticles() {
    loading = true;
    try {
      articles = await api.fetchArticles(null, selectedCategory);
    } catch (e) {
      console.error($t('common.errorLoading'), e);
      articles = [];
    } finally {
      loading = false;
    }
  }

  function selectCategory(category) {
    selectedCategory = category === selectedCategory ? null : category;
    loadArticles();
  }

  function getCategoryLabel(category) {
    return $t(`articleCategories.${category}`);
  }

  function getCategoryColor(category, isSelected) {
    const base = 'px-4 py-2 rounded-full text-sm font-medium transition-all duration-200';
    if (isSelected) {
      const colors = {
        all: 'bg-teal-600 text-white',
        latest: 'bg-red-600 text-white',
        policy: 'bg-blue-600 text-white',
        research: 'bg-purple-600 text-white',
        outbreak: 'bg-orange-600 text-white',
        lifestyle: 'bg-green-600 text-white',
        local: 'bg-amber-600 text-white'
      };
      return `${base} ${colors[category] || 'bg-slate-600 text-white'}`;
    }
    return `${base} bg-slate-100 text-slate-600 hover:bg-slate-200`;
  }

  $: filteredArticles = searchQuery
    ? articles.filter(a => a.title?.toLowerCase().includes(searchQuery.toLowerCase()) || a.excerpt?.toLowerCase().includes(searchQuery.toLowerCase()))
    : articles;

  loadArticles();
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div>
      <h1 class="text-2xl font-bold text-slate-900">{$t('berita.title')}</h1>
      <p class="text-slate-500 mt-1">{$t('berita.titleAlt')}</p>
    </div>
  </div>

  <!-- Search -->
  <div class="relative">
    <svg class="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
    </svg>
    <input
      type="text"
      placeholder={$t('common.searchPlaceholder')}
      bind:value={searchQuery}
      class="w-full pl-12 pr-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent"
    />
  </div>

  <!-- Categories -->
  <div class="flex flex-wrap gap-2">
    <button
      onclick={() => { selectedCategory = null; loadArticles(); }}
      class={getCategoryColor('all', selectedCategory === null)}
    >
      {$t('berita.allCategories')}
    </button>
    {#each categories as category}
      <button
        onclick={() => selectCategory(category)}
        class={getCategoryColor(category, selectedCategory === category)}
      >
        {getCategoryLabel(category)}
      </button>
    {/each}
  </div>

  <!-- Loading -->
  {#if loading}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each Array(6) as _}
        <div class="bg-white rounded-2xl border border-slate-200 overflow-hidden animate-pulse">
          <div class="aspect-[16/9] bg-slate-200"></div>
          <div class="p-4 space-y-3">
            <div class="h-4 bg-slate-200 rounded w-3/4"></div>
            <div class="h-3 bg-slate-200 rounded w-full"></div>
            <div class="h-3 bg-slate-200 rounded w-2/3"></div>
          </div>
        </div>
      {/each}
    </div>
  <!-- Articles Grid -->
  {:else if filteredArticles.length > 0}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each filteredArticles as article}
        <ArticleCard
          {article}
          onClick={() => onNavigateToArticle(article)}
        />
      {/each}
    </div>
  <!-- Empty State -->
  {:else}
    <div class="bg-white rounded-2xl border border-slate-200 p-12 text-center">
      <svg class="w-16 h-16 mx-auto text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
      </svg>
      <h3 class="mt-4 text-lg font-medium text-slate-900">{$t('berita.noArticles')}</h3>
      <p class="mt-2 text-slate-500">{$t('berita.noArticlesMessage')}</p>
    </div>
  {/if}
</div>
