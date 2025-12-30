<script>
  import { t } from 'svelte-i18n';
  import VideoCard from '$lib/components/VideoCard.svelte';
  import * as api from '$lib/utils/api.js';

  export let onWatchVideo = () => {};

  let videos = [];
  let categories = ['tutorial', 'firstAid', 'exercise', 'nutrition', 'mentalHealth', 'childHealth', 'seniorHealth'];
  let selectedCategory = null;
  let loading = true;
  let searchQuery = '';

  async function loadVideos() {
    loading = true;
    try {
      videos = await api.fetchVideos(null, selectedCategory);
    } catch (e) {
      console.error($t('common.errorLoading'), e);
      videos = [];
    } finally {
      loading = false;
    }
  }

  function selectCategory(category) {
    selectedCategory = category === selectedCategory ? null : category;
    loadVideos();
  }

  function getCategoryLabel(category) {
    const labels = {
      tutorial: $t('videoCategories.tutorial'),
      firstAid: $t('videoCategories.firstAid'),
      exercise: $t('videoCategories.exercise'),
      nutrition: $t('videoCategories.nutrition'),
      mentalHealth: $t('videoCategories.mentalHealth'),
      childHealth: $t('videoCategories.childHealth'),
      seniorHealth: $t('videoCategories.seniorHealth')
    };
    return labels[category] || category;
  }

  function getCategoryColor(category, isSelected) {
    const base = 'px-4 py-2 rounded-full text-sm font-medium transition-all duration-200';
    if (isSelected) {
      const colors = {
        tutorial: 'bg-blue-600 text-white',
        firstAid: 'bg-red-600 text-white',
        exercise: 'bg-green-600 text-white',
        nutrition: 'bg-amber-600 text-white',
        mentalHealth: 'bg-pink-600 text-white',
        childHealth: 'bg-cyan-600 text-white',
        seniorHealth: 'bg-purple-600 text-white'
      };
      return `${base} ${colors[category] || 'bg-slate-600 text-white'}`;
    }
    return `${base} bg-slate-100 text-slate-600 hover:bg-slate-200`;
  }

  $: filteredVideos = searchQuery
    ? videos.filter(v => v.title?.toLowerCase().includes(searchQuery.toLowerCase()) || v.description?.toLowerCase().includes(searchQuery.toLowerCase()))
    : videos;

  loadVideos();
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div>
      <h1 class="text-2xl font-bold text-slate-900">{$t('video.title')}</h1>
      <p class="text-slate-500 mt-1">{$t('video.titleAlt')}</p>
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
      onclick={() => { selectedCategory = null; loadVideos(); }}
      class="px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 {selectedCategory === null ? 'bg-teal-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
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
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each Array(6) as _}
        <div class="bg-white rounded-2xl border border-slate-200 overflow-hidden animate-pulse">
          <div class="aspect-video bg-slate-200"></div>
          <div class="p-4 space-y-3">
            <div class="h-4 bg-slate-200 rounded w-3/4"></div>
            <div class="h-3 bg-slate-200 rounded w-full"></div>
          </div>
        </div>
      {/each}
    </div>
  <!-- Videos Grid -->
  {:else if filteredVideos.length > 0}
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each filteredVideos as video}
        <VideoCard
          {video}
          onClick={() => onWatchVideo(video)}
        />
      {/each}
    </div>
  <!-- Empty State -->
  {:else}
    <div class="bg-white rounded-2xl border border-slate-200 p-12 text-center">
      <svg class="w-16 h-16 mx-auto text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
      </svg>
      <h3 class="mt-4 text-lg font-medium text-slate-900">{$t('video.noVideos')}</h3>
      <p class="mt-2 text-slate-500">{$t('video.noVideosMessage')}</p>
    </div>
  {/if}
</div>
