<script>
  import { t } from 'svelte-i18n';
  import { onMount } from 'svelte';
  import ArticleCard from '$lib/components/ArticleCard.svelte';
  import VideoCard from '$lib/components/VideoCard.svelte';
  import ContentListItem from '$lib/components/ContentListItem.svelte';
  import * as api from '$lib/utils/api.js';

  export let onNavigateToArticleEditor = () => {};
  export let onNavigateToVideoManager = () => {};
  export let onEditArticle = () => {};
  export let onEditVideo = () => {};
  export let token = null;

  // Data
  let stats = { totalArticles: 0, totalVideos: 0, totalViews: 0, drafts: 0 };
  let allArticles = [];
  let allVideos = [];
  let loading = true;

  // UI State
  let viewMode = 'list'; // 'list' or 'grid'
  let activeFilter = 'all'; // 'all', 'articles', 'videos', 'drafts'
  let sortBy = 'newest'; // 'newest', 'oldest', 'views', 'titleAZ'
  let searchQuery = '';
  let selectedItems = new Set();
  let currentPage = 1;
  const itemsPerPage = 10;

  // Load view preference from localStorage
  onMount(() => {
    const savedView = localStorage.getItem('cmsViewMode');
    if (savedView) viewMode = savedView;
    loadDashboardData();
  });

  async function loadDashboardData() {
    loading = true;
    try {
      const [dashboardStats, articles, videos] = await Promise.all([
        api.fetchDashboardStats(token),
        api.fetchArticles(token, null, true), // all=true to include drafts
        api.fetchVideos(token)
      ]);

      // Transform backend stats
      stats = {
        totalArticles: dashboardStats.articles?.total || 0,
        totalVideos: dashboardStats.videos?.total || 0,
        totalViews: (dashboardStats.total_views?.articles || 0) + (dashboardStats.total_views?.videos || 0),
        drafts: dashboardStats.articles?.drafts || 0
      };

      allArticles = articles || [];
      allVideos = videos || [];
    } catch (e) {
      console.error($t('common.errorLoading'), e);
    } finally {
      loading = false;
    }
  }

  // Combine and filter content
  $: combinedContent = (() => {
    let content = [];

    if (activeFilter === 'all' || activeFilter === 'articles') {
      content = [...content, ...allArticles.map(a => ({ ...a, type: 'article' }))];
    }
    if (activeFilter === 'all' || activeFilter === 'videos') {
      content = [...content, ...allVideos.map(v => ({ ...v, type: 'video' }))];
    }
    if (activeFilter === 'drafts') {
      content = allArticles.filter(a => a.status === 'draft').map(a => ({ ...a, type: 'article' }));
    }

    // Apply search filter
    if (searchQuery.trim()) {
      const query = searchQuery.toLowerCase();
      content = content.filter(item =>
        item.title?.toLowerCase().includes(query)
      );
    }

    // Apply sorting
    content.sort((a, b) => {
      if (sortBy === 'newest') {
        const dateA = new Date(a.published_at || a.publishedAt || a.created_at || a.createdAt || 0);
        const dateB = new Date(b.published_at || b.publishedAt || b.created_at || b.createdAt || 0);
        return dateB - dateA;
      } else if (sortBy === 'oldest') {
        const dateA = new Date(a.published_at || a.publishedAt || a.created_at || a.createdAt || 0);
        const dateB = new Date(b.published_at || b.publishedAt || b.created_at || b.createdAt || 0);
        return dateA - dateB;
      } else if (sortBy === 'views') {
        return (b.view_count || b.viewCount || 0) - (a.view_count || a.viewCount || 0);
      } else if (sortBy === 'titleAZ') {
        return (a.title || '').localeCompare(b.title || '');
      }
      return 0;
    });

    return content;
  })();

  // Pagination
  $: totalPages = Math.ceil(combinedContent.length / itemsPerPage);
  $: paginatedContent = combinedContent.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  // Reset to page 1 when filters change
  $: if (activeFilter || searchQuery || sortBy) {
    currentPage = 1;
    selectedItems = new Set();
  }

  // View mode toggle
  function toggleViewMode(mode) {
    viewMode = mode;
    localStorage.setItem('cmsViewMode', mode);
  }

  // Selection handlers
  function toggleSelectAll() {
    if (selectedItems.size === paginatedContent.length) {
      selectedItems = new Set();
    } else {
      selectedItems = new Set(paginatedContent.map(item => item.id));
    }
  }

  function toggleItemSelection(itemId) {
    const newSet = new Set(selectedItems);
    if (newSet.has(itemId)) {
      newSet.delete(itemId);
    } else {
      newSet.add(itemId);
    }
    selectedItems = newSet;
  }

  // Content actions
  async function handleToggleStatus(item) {
    if (item.type !== 'article') return;

    try {
      const newStatus = item.status === 'draft' ? 'published' : 'draft';
      await api.updateArticle(token, item.id, { status: newStatus });
      await loadDashboardData();
    } catch (e) {
      console.error('Failed to toggle status:', e);
      alert($t('common.error') + ': ' + e.message);
    }
  }

  async function handleDeleteItem(item) {
    const confirmMsg = item.type === 'article'
      ? $t('articleEditor.deleteConfirm')
      : $t('videoManager.deleteConfirm');

    if (!confirm(confirmMsg)) return;

    try {
      if (item.type === 'article') {
        await api.deleteArticle(token, item.id);
      } else {
        await api.deleteVideo(token, item.id);
      }
      await loadDashboardData();
      const newSet = new Set(selectedItems);
      newSet.delete(item.id);
      selectedItems = newSet;
    } catch (e) {
      console.error('Failed to delete:', e);
      alert($t('common.error') + ': ' + e.message);
    }
  }

  async function handleBulkDelete() {
    if (selectedItems.size === 0) return;

    const confirmMsg = $t('cms.bulkActions.deleteSelected') + ` (${selectedItems.size})?`;
    if (!confirm(confirmMsg)) return;

    try {
      const deletePromises = [];
      paginatedContent.forEach(item => {
        if (selectedItems.has(item.id)) {
          if (item.type === 'article') {
            deletePromises.push(api.deleteArticle(token, item.id));
          } else {
            deletePromises.push(api.deleteVideo(token, item.id));
          }
        }
      });

      await Promise.all(deletePromises);
      await loadDashboardData();
      selectedItems = new Set();
    } catch (e) {
      console.error('Failed to bulk delete:', e);
      alert($t('common.error') + ': ' + e.message);
    }
  }

  function handleEditItem(item) {
    if (item.type === 'article') {
      onEditArticle(item);
    } else {
      onEditVideo(item);
    }
  }

  // Empty state message
  $: emptyStateMessage = (() => {
    if (searchQuery.trim()) {
      return $t('cms.noSearchResults', { values: { query: searchQuery } });
    }
    if (activeFilter === 'articles') return $t('cms.noArticlesFound');
    if (activeFilter === 'videos') return $t('cms.noVideosFound');
    if (activeFilter === 'drafts') return $t('cms.noDraftsFound');
    return $t('cms.noContent');
  })();
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
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {#each Array(4) as _}
          <div class="h-28 bg-slate-200 rounded-2xl"></div>
        {/each}
      </div>
      <div class="h-64 bg-slate-200 rounded-2xl"></div>
    </div>
  {:else}
    <!-- Stats Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Total Articles -->
      <div class="bg-white rounded-2xl border border-slate-200 p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-500">{$t('cms.stats.totalArticles')}</p>
            <p class="text-2xl font-bold text-slate-900 mt-1">{stats.totalArticles}</p>
          </div>
          <div class="w-12 h-12 bg-teal-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Total Videos -->
      <div class="bg-white rounded-2xl border border-slate-200 p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-500">{$t('cms.stats.totalVideos')}</p>
            <p class="text-2xl font-bold text-slate-900 mt-1">{stats.totalVideos}</p>
          </div>
          <div class="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Drafts -->
      <div class="bg-white rounded-2xl border border-slate-200 p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-500">{$t('cms.stats.drafts')}</p>
            <p class="text-2xl font-bold text-slate-900 mt-1">{stats.drafts}</p>
          </div>
          <div class="w-12 h-12 bg-yellow-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Total Views -->
      <div class="bg-white rounded-2xl border border-slate-200 p-5">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-sm text-slate-500">{$t('cms.stats.totalViews')}</p>
            <p class="text-2xl font-bold text-slate-900 mt-1">{stats.totalViews}</p>
          </div>
          <div class="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Content Management Section -->
    <div class="bg-white rounded-2xl border border-slate-200">
      <!-- Search Bar -->
      <div class="p-5 border-b border-slate-100">
        <div class="relative">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            type="text"
            bind:value={searchQuery}
            placeholder={$t('cms.searchPlaceholder')}
            class="w-full pl-10 pr-4 py-2.5 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent"
          />
        </div>
      </div>

      <!-- Filter Tabs & Controls -->
      <div class="px-5 py-4 border-b border-slate-100 flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <!-- Filter Tabs -->
        <div class="flex gap-2 overflow-x-auto">
          <button
            onclick={() => activeFilter = 'all'}
            class="px-4 py-2 rounded-lg font-medium transition-colors whitespace-nowrap {activeFilter === 'all' ? 'bg-teal-100 text-teal-700' : 'text-slate-600 hover:bg-slate-100'}"
          >
            {$t('cms.filters.all')}
          </button>
          <button
            onclick={() => activeFilter = 'articles'}
            class="px-4 py-2 rounded-lg font-medium transition-colors whitespace-nowrap {activeFilter === 'articles' ? 'bg-teal-100 text-teal-700' : 'text-slate-600 hover:bg-slate-100'}"
          >
            {$t('cms.filters.articles')}
          </button>
          <button
            onclick={() => activeFilter = 'videos'}
            class="px-4 py-2 rounded-lg font-medium transition-colors whitespace-nowrap {activeFilter === 'videos' ? 'bg-teal-100 text-teal-700' : 'text-slate-600 hover:bg-slate-100'}"
          >
            {$t('cms.filters.videos')}
          </button>
          <button
            onclick={() => activeFilter = 'drafts'}
            class="px-4 py-2 rounded-lg font-medium transition-colors whitespace-nowrap {activeFilter === 'drafts' ? 'bg-teal-100 text-teal-700' : 'text-slate-600 hover:bg-slate-100'}"
          >
            {$t('cms.filters.drafts')}
          </button>
        </div>

        <!-- Sort & View Controls -->
        <div class="flex items-center gap-3">
          <!-- Sort Dropdown -->
          <select
            bind:value={sortBy}
            class="px-3 py-2 border border-slate-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-teal-500"
          >
            <option value="newest">{$t('cms.sort.newest')}</option>
            <option value="oldest">{$t('cms.sort.oldest')}</option>
            <option value="views">{$t('cms.sort.views')}</option>
            <option value="titleAZ">{$t('cms.sort.titleAZ')}</option>
          </select>

          <!-- View Toggle -->
          <div class="flex border border-slate-200 rounded-lg overflow-hidden">
            <button
              onclick={() => toggleViewMode('list')}
              class="p-2 transition-colors {viewMode === 'list' ? 'bg-teal-100 text-teal-700' : 'text-slate-600 hover:bg-slate-100'}"
              title={$t('cms.view.list')}
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
            <button
              onclick={() => toggleViewMode('grid')}
              class="p-2 transition-colors {viewMode === 'grid' ? 'bg-teal-100 text-teal-700' : 'text-slate-600 hover:bg-slate-100'}"
              title={$t('cms.view.grid')}
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Bulk Actions Bar -->
      {#if selectedItems.size > 0}
        <div class="px-5 py-3 bg-teal-50 border-b border-teal-100 flex items-center justify-between">
          <span class="text-sm font-medium text-teal-700">
            {$t('cms.bulkActions.selected', { values: { n: selectedItems.size } })}
          </span>
          <button
            onclick={handleBulkDelete}
            class="px-4 py-2 bg-red-600 text-white rounded-lg font-medium hover:bg-red-700 transition-colors flex items-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            {$t('cms.bulkActions.deleteSelected')}
          </button>
        </div>
      {/if}

      <!-- Content List/Grid -->
      {#if paginatedContent.length === 0}
        <!-- Empty State -->
        <div class="p-12 text-center">
          <div class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <p class="text-slate-600 font-medium">{emptyStateMessage}</p>
        </div>
      {:else if viewMode === 'list'}
        <!-- List View -->
        <div>
          <!-- List Header -->
          <div class="flex items-center gap-4 px-5 py-3 bg-slate-50 border-b border-slate-100 text-xs font-medium text-slate-500 uppercase tracking-wider">
            <input
              type="checkbox"
              checked={selectedItems.size === paginatedContent.length && paginatedContent.length > 0}
              onchange={toggleSelectAll}
              class="w-4 h-4 text-teal-600 border-slate-300 rounded focus:ring-teal-500 focus:ring-2"
            />
            <div class="w-8"></div>
            <div class="w-20"></div>
            <div class="flex-1">{$t('articleEditor.articleTitle')}</div>
            <div class="w-24 text-center">{$t('articleEditor.status')}</div>
            <div class="w-32">{$t('articleEditor.category')}</div>
            <div class="w-28">{$t('berita.publishedOn')}</div>
            <div class="w-20">{$t('cms.stats.totalViews')}</div>
            <div class="w-24 text-center">{$t('users.actions')}</div>
          </div>

          <!-- List Items -->
          {#each paginatedContent as item (item.id)}
            <ContentListItem
              {item}
              type={item.type}
              selected={selectedItems.has(item.id)}
              onSelect={() => toggleItemSelection(item.id)}
              onEdit={() => handleEditItem(item)}
              onDelete={() => handleDeleteItem(item)}
              onToggleStatus={() => handleToggleStatus(item)}
            />
          {/each}
        </div>
      {:else}
        <!-- Grid View -->
        <div class="p-5 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {#each paginatedContent as item (item.id)}
            {#if item.type === 'article'}
              <ArticleCard
                article={item}
                onClick={() => handleEditItem(item)}
              />
            {:else}
              <VideoCard
                video={item}
                onClick={() => handleEditItem(item)}
                showActions={true}
                onEdit={() => handleEditItem(item)}
                onDelete={() => handleDeleteItem(item)}
              />
            {/if}
          {/each}
        </div>
      {/if}

      <!-- Pagination -->
      {#if totalPages > 1}
        <div class="px-5 py-4 border-t border-slate-100 flex items-center justify-between">
          <div class="text-sm text-slate-500">
            {$t('common.search')} {(currentPage - 1) * itemsPerPage + 1}-{Math.min(currentPage * itemsPerPage, combinedContent.length)} of {combinedContent.length}
          </div>
          <div class="flex items-center gap-2">
            <button
              onclick={() => currentPage = Math.max(1, currentPage - 1)}
              disabled={currentPage === 1}
              class="px-3 py-2 border border-slate-200 rounded-lg text-sm font-medium text-slate-600 hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {$t('common.no')}
            </button>

            {#each Array(totalPages) as _, i}
              {#if i + 1 === 1 || i + 1 === totalPages || (i + 1 >= currentPage - 1 && i + 1 <= currentPage + 1)}
                <button
                  onclick={() => currentPage = i + 1}
                  class="px-3 py-2 rounded-lg text-sm font-medium transition-colors {currentPage === i + 1 ? 'bg-teal-600 text-white' : 'text-slate-600 hover:bg-slate-100'}"
                >
                  {i + 1}
                </button>
              {:else if i + 1 === currentPage - 2 || i + 1 === currentPage + 2}
                <span class="px-2 text-slate-400">...</span>
              {/if}
            {/each}

            <button
              onclick={() => currentPage = Math.min(totalPages, currentPage + 1)}
              disabled={currentPage === totalPages}
              class="px-3 py-2 border border-slate-200 rounded-lg text-sm font-medium text-slate-600 hover:bg-slate-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {$t('common.yes')}
            </button>
          </div>
        </div>
      {/if}
    </div>
  {/if}
</div>
