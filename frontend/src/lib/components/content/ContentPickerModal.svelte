<script>
  import { onMount } from 'svelte';
  import { t } from 'svelte-i18n';
  import { fetchAllContent, fetchCategories, fetchPopularContent, incrementAttachmentCount } from '$lib/utils/api.js';
  import ContentPreviewPanel from './ContentPreviewPanel.svelte';
  import EmptyState from '../ui/EmptyState.svelte';

  // Constants
  const MAX_SELECTION = 3; // Maximum number of content items that can be attached to a reminder
  const MAX_RECENT = 5;
  const MAX_POPULAR = 5;

  let {
    show = false,
    selectedContent = [],
    userRole = 'volunteer', // 'superadmin', 'admin', or 'volunteer'
    onClose = () => {},
    onSelect = () => {}
  } = $props();

  // State
  let searchQuery = $state('');
  let activeTab = $state('all'); // 'all', 'article', 'video'
  let articles = $state([]);
  let videos = $state([]);
  let categories = $state([]);
  let isLoading = $state(false);
  let errorMessage = $state('');
  let searchInput = $state(null);

  // Preview state
  let previewContent = $state(null);
  let lastFocusedCard = $state(null);

  // Session cache for content with timestamp
  let contentCache = $state(null);
  let categoriesCache = $state(null);

  // Popular content (from API)
  let popularContent = $state([]);
  let isLoadingPopular = $state(false);

  // Recent attachments (from localStorage) - isolated per user role
  let recentAttachments = $state([]);
  const isAdmin = $derived(userRole === 'superadmin' || userRole === 'admin');
  const recentAttachmentsKey = $derived(`prima_recent_attachments_${userRole}`);

  // Debounce timer for search
  let searchTimeout = $state(null);
  let debouncedQuery = $state('');

  // Filtered content based on search and tab (using debounced query)
  let filteredArticles = $derived.by(() => {
    if (!debouncedQuery.trim()) return articles;
    const query = debouncedQuery.toLowerCase();
    return articles.filter(a =>
      a.title.toLowerCase().includes(query) ||
      (a.excerpt && a.excerpt.toLowerCase().includes(query))
    );
  });

  let filteredVideos = $derived.by(() => {
    if (!debouncedQuery.trim()) return videos;
    const query = debouncedQuery.toLowerCase();
    return videos.filter(v =>
      v.title.toLowerCase().includes(query)
    );
  });

  let displayedArticles = $derived.by(() => {
    if (activeTab === 'video') return [];
    return filteredArticles;
  });

  let displayedVideos = $derived.by(() => {
    if (activeTab === 'article') return [];
    return filteredVideos;
  });

  let totalResults = $derived.by(() => displayedArticles.length + displayedVideos.length);

  // Check if content is empty (no articles and no videos)
  let hasNoContent = $derived.by(() => articles.length === 0 && videos.length === 0);

  // Check if search has no results
  let hasNoSearchResults = $derived.by(() => !hasNoContent && totalResults === 0 && debouncedQuery.trim() !== '');

  // Check if initial state (no search, has content)
  let hasInitialState = $derived.by(() => !hasNoContent && !hasNoSearchResults && debouncedQuery.trim() === '');

  // Get unique popular content that is not already selected
  let availablePopularContent = $derived.by(() => {
    return popularContent.filter(item =>
      !selectedContent.some(c => c.id === item.id && c.type === item.type)
    ).slice(0, MAX_POPULAR);
  });

  // Get recent attachments that are not already selected and exist in the content
  let availableRecentAttachments = $derived.by(() => {
    return recentAttachments
      .filter(item =>
        !selectedContent.some(c => c.id === item.id && c.type === item.type) &&
        (articles.some(a => a.id === item.id) || videos.some(v => v.id === item.id))
      )
      .slice(0, MAX_RECENT);
  });

  // Merge recent and popular content for initial state display
  let suggestedContent = $derived.by(() => {
    const result = [];
    const seenIds = new Set();

    // Add recent first
    for (const item of availableRecentAttachments) {
      if (!seenIds.has(item.id)) {
        seenIds.add(item.id);
        result.push({ ...item, source: 'recent' });
      }
    }

    // Add popular
    for (const item of availablePopularContent) {
      if (!seenIds.has(item.id)) {
        seenIds.add(item.id);
        result.push({ ...item, source: 'popular' });
      }
    }

    return result;
  });

  // Check if content is selected
  function isSelected(content) {
    return selectedContent.some(c => c.id === content.id && c.type === getContentType(content));
  }

  function getContentType(content) {
    // Determine type based on field presence
    return content.YouTubeID ? 'video' : 'article';
  }

  // Load recent attachments from localStorage
  function loadRecentAttachments() {
    try {
      const stored = localStorage.getItem(recentAttachmentsKey);
      if (stored) {
        recentAttachments = JSON.parse(stored);
      }
    } catch (e) {
      console.error('Error loading recent attachments:', e);
      recentAttachments = [];
    }
  }

  // Save recent attachments to localStorage
  function saveRecentAttachments() {
    try {
      localStorage.setItem(recentAttachmentsKey, JSON.stringify(recentAttachments));
    } catch (e) {
      console.error('Error saving recent attachments:', e);
    }
  }

  // Add content to recent attachments
  function addToRecentAttachments(content) {
    const type = getContentType(content);
    const item = {
      id: content.id,
      type,
      title: content.title,
      thumbnail: type === 'video' ? content.thumbnailURL : content.heroImages?.hero1x1,
      addedAt: Date.now()
    };

    // Remove if already exists
    recentAttachments = recentAttachments.filter(a => a.id !== content.id || a.type !== type);

    // Add to beginning
    recentAttachments = [item, ...recentAttachments];

    // Keep only MAX_RECENT items
    if (recentAttachments.length > MAX_RECENT) {
      recentAttachments = recentAttachments.slice(0, MAX_RECENT);
    }

    saveRecentAttachments();
  }

  // Load popular content from API
  async function loadPopularContent() {
    if (popularContent.length > 0) return; // Already loaded

    isLoadingPopular = true;
    try {
      const data = await fetchPopularContent(MAX_POPULAR);
      popularContent = data;
    } catch (e) {
      console.error('Error loading popular content:', e);
      popularContent = [];
    } finally {
      isLoadingPopular = false;
    }
  }

  // Load content from API with session caching
  async function loadContent() {
    if (contentCache) {
      articles = contentCache.articles;
      videos = contentCache.videos;
      return;
    }

    isLoading = true;
    errorMessage = '';

    try {
      const data = await fetchAllContent(null, 'all');
      contentCache = data;
      articles = data.articles;
      videos = data.videos;
    } catch (err) {
      errorMessage = err.message || 'Gagal memuat konten';
      console.error('Error loading content:', err);
    } finally {
      isLoading = false;
    }
  }

  // Load categories
  async function loadCategories() {
    if (categoriesCache) {
      categories = categoriesCache;
      return;
    }

    try {
      const data = await fetchCategories();
      categoriesCache = data;
      categories = data;
    } catch (err) {
      console.error('Error loading categories:', err);
    }
  }

  // Handle content selection
  function handleSelect(content) {
    const type = getContentType(content);

    // Build attachment with excerpt, slug, and youtubeId for WhatsApp preview
    let attachment = {
      id: content.id,
      type: type,
      title: content.title,
      excerpt: content.excerpt || '',
      slug: content.slug || '',
      youtubeId: content.YouTubeID || '',
      url: type === 'video'
        ? `https://www.youtube.com/watch?v=${content.YouTubeID}`
        : `/berita/${content.slug}`
    };

    if (isSelected(content)) {
      // Deselect
      onSelect(selectedContent.filter(c => !(c.id === content.id && c.type === type)));
    } else {
      // Select (max MAX_SELECTION)
      if (selectedContent.length >= MAX_SELECTION) return;
      onSelect([...selectedContent, attachment]);
      // Add to recent attachments
      addToRecentAttachments(content);
      // Increment attachment count on backend (fire and forget)
      incrementAttachmentCount(type, content.id).catch(e => {
        console.error('Failed to increment attachment count:', e);
      });
    }
  }

  // Handle search input with debounce (300ms)
  function handleSearchInput(e) {
    searchQuery = e.target.value;

    // Clear existing timeout
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }

    // Set new debounced query after 300ms
    searchTimeout = setTimeout(() => {
      debouncedQuery = searchQuery;
    }, 300);
  }

  // Clear search
  function clearSearch() {
    searchQuery = '';
    debouncedQuery = '';
    if (searchTimeout) {
      clearTimeout(searchTimeout);
      searchTimeout = null;
    }
    if (searchInput) {
      searchInput.focus();
    }
  }

  // Refresh content (invalidate cache)
  function refreshContent() {
    contentCache = null;
    loadContent();
  }

  // Close modal
  function handleClose() {
    searchQuery = '';
    debouncedQuery = '';
    if (searchTimeout) {
      clearTimeout(searchTimeout);
      searchTimeout = null;
    }
    activeTab = 'all';
    onClose();
  }

  // Open preview panel
  function openPreview(content, event) {
    lastFocusedCard = event?.target;
    previewContent = content;
  }

  // Close preview panel
  function closePreview() {
    previewContent = null;
    // Return focus to the card that opened the preview
    if (lastFocusedCard) {
      lastFocusedCard.focus();
      lastFocusedCard = null;
    }
  }

  // Handle attach from preview
  function handleAttachFromPreview() {
    if (previewContent) {
      handleSelect(previewContent);
      closePreview();
    }
  }

  // Handle keyboard navigation
  function handleKeydown(e) {
    if (e.key === 'Escape') {
      if (previewContent) {
        closePreview();
      } else {
        handleClose();
      }
    }
  }

  // Focus trap and initial focus
  function handleMount(node) {
    if (searchInput) {
      searchInput.focus();
    }

    const cleanup = node.addEventListener('keydown', handleKeydown);
    return () => {
      cleanup?.();
      node.removeEventListener('keydown', handleKeydown);
    };
  }

  // Load data when modal opens
  $effect(() => {
    if (show) {
      loadContent();
      loadCategories();
      loadRecentAttachments();
      loadPopularContent();
    }
  });
</script>

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm"
    onclick={handleClose}
    onkeydown={(e) => e.key === 'Escape' && handleClose()}
    role="button"
    tabindex="0"
    aria-label="Tutup modal"
  >
    <div
      class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-2xl max-h-[80vh] flex flex-col overflow-hidden"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
      role="dialog"
      aria-modal="true"
      aria-labelledby="content-picker-title"
      tabindex="-1"
      use:handleMount
    >
      <!-- Header -->
      <div class="flex-shrink-0 p-4 sm:p-6 border-b border-slate-100">
        <div class="flex items-center justify-between mb-4">
          <h2 id="content-picker-title" class="text-lg sm:text-xl font-semibold text-slate-900">
            {$t('content.picker.title') || 'Pilih Konten Edukasi'}
          </h2>
          <button
            onclick={handleClose}
            class="p-2 rounded-full hover:bg-slate-100 transition-colors"
            aria-label="Tutup"
          >
            <svg class="w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Search Input -->
        <div class="relative">
          <svg
            class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            bind:this={searchInput}
            value={searchQuery}
            oninput={handleSearchInput}
            type="text"
            placeholder={$t('content.picker.search') || 'Cari artikel atau video...'}
            class="w-full pl-10 pr-10 py-3 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent"
            aria-label={$t('content.picker.search') || 'Cari konten'}
            tabindex="0"
          />
          {#if searchQuery}
            <button
              onclick={clearSearch}
              class="absolute right-3 top-1/2 -translate-y-1/2 p-1 rounded-full hover:bg-slate-100"
              aria-label="Hapus pencarian"
            >
              <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          {/if}
        </div>

        <!-- Category Tabs -->
        <div class="flex gap-2 mt-4" role="tablist" aria-label="Jenis konten">
          {#each [
            { id: 'all', label: $t('content.picker.tab.all') || 'Semua' },
            { id: 'article', label: $t('content.picker.tab.articles') || 'Artikel' },
            { id: 'video', label: $t('content.picker.tab.videos') || 'Video' }
          ] as tab}
            <button
              onclick={() => activeTab = tab.id}
              role="tab"
              aria-selected={activeTab === tab.id}
              aria-controls={`${tab.id}-panel`}
              class="px-4 py-2 rounded-lg font-medium text-sm transition-colors duration-200
                {activeTab === tab.id
                  ? 'bg-teal-50 text-teal-700 border border-teal-200'
                  : 'text-slate-600 hover:bg-slate-50 border border-transparent'}"
            >
              {tab.label}
            </button>
          {/each}
        </div>

        <!-- Selection count -->
        {#if selectedContent.length > 0}
          <div class="mt-3 px-3 py-2 bg-teal-50 rounded-lg flex items-center justify-between">
            <span class="text-sm text-teal-700">
              {selectedContent.length} dipilih
            </span>
            <button
              onclick={() => onSelect([])}
              class="text-xs text-teal-600 hover:text-teal-800"
            >
              Hapus semua
            </button>
          </div>
        {/if}
      </div>

      <!-- Content List -->
      <div
        class="flex-1 overflow-y-auto p-4 sm:p-6"
        role="tabpanel"
        id={`${activeTab}-panel`}
        aria-label="{activeTab} konten"
      >
        {#if isLoading}
          <div class="flex items-center justify-center py-12">
            <svg class="animate-spin h-8 w-8 text-teal-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>
        {:else if errorMessage}
          <div class="text-center py-8">
            <p class="text-red-500 mb-2">{errorMessage}</p>
            <div class="flex gap-2 justify-center">
              <button
                onclick={refreshContent}
                class="px-4 py-2 bg-teal-500 text-white rounded-lg hover:bg-teal-600"
              >
                Refresh
              </button>
              <button
                onclick={loadContent}
                class="px-4 py-2 border border-slate-200 rounded-lg hover:bg-slate-50"
              >
                Coba lagi
              </button>
            </div>
          </div>
        {:else if hasNoContent}
          <!-- No content in system - Admin only CTA -->
          <EmptyState
            icon={`<svg class="w-12 h-12 mx-auto text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>`}
            title={$t('content.picker.empty') || 'Belum ada konten edukasi'}
            description="Tambahkan artikel dan video edukasi untuk pasien Anda."
            actionLabel={isAdmin ? $t('cms.addArticle') || 'Tambah Konten' : undefined}
            onAction={isAdmin ? () => window.location.href = '/cms' : undefined}
          />
        {:else if hasNoSearchResults}
          <!-- No search results with popular content suggestions -->
          <EmptyState
            icon={`<svg class="w-12 h-12 mx-auto text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>`}
            title={$t('content.picker.noResults').replace('{query}', `"${searchQuery}"`) || `Tidak ada hasil untuk "${searchQuery}"`}
            description={$t('content.picker.noResultsGuidance') || 'Coba kata kunci lain atau browse kategori'}
          />
          <!-- Popular Content suggestions -->
          {#if suggestedContent.length > 0}
            <div class="mt-4">
              <h4 class="text-sm font-medium text-slate-500 mb-3 uppercase tracking-wide">
                {$t('content.picker.popularContent') || 'Konten Populer'}
              </h4>
              <div class="space-y-3">
                {#each suggestedContent as item (item.id + '-' + item.type)}
                  <button
                    onclick={() => {
                      const content = item.type === 'video'
                        ? videos.find(v => v.id === item.id)
                        : articles.find(a => a.id === item.id);
                      if (content) openPreview(content);
                    }}
                    class="w-full text-left p-3 rounded-xl border border-slate-200 hover:border-slate-300 hover:bg-slate-50 transition-all duration-150"
                  >
                    <div class="flex items-center gap-3">
                      <div class="flex-shrink-0 w-12 h-12 bg-slate-100 rounded-lg overflow-hidden">
                        {#if item.thumbnail}
                          <img src={item.thumbnail} alt="" class="w-full h-full object-cover" />
                        {:else}
                          <div class="w-full h-full flex items-center justify-center text-slate-400">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                            </svg>
                          </div>
                        {/if}
                      </div>
                      <div class="flex-1 min-w-0">
                        <p class="font-medium text-slate-900 truncate">{item.title}</p>
                        <p class="text-xs text-slate-500 capitalize">{item.type === 'video' ? 'Video' : 'Artikel'} · {item.source === 'recent' ? ($t('content.picker.recentLabel') || 'Baru saja') : ($t('content.picker.popularLabel') || 'Populer')}</p>
                      </div>
                    </div>
                  </button>
                {/each}
              </div>
            </div>
          {/if}
        {:else if hasInitialState}
          <!-- Initial state with Recent Attachments and Popular Content -->
          {#if availableRecentAttachments.length > 0}
            <div class="mb-6">
              <h4 class="text-sm font-medium text-slate-500 mb-3 uppercase tracking-wide">
                {$t('content.picker.recentAttachments') || 'Lampiran Terbaru'}
              </h4>
              <div class="space-y-3">
                {#each availableRecentAttachments as item (item.id + '-' + item.type)}
                  <button
                    onclick={() => {
                      const content = item.type === 'video'
                        ? videos.find(v => v.id === item.id)
                        : articles.find(a => a.id === item.id);
                      if (content) openPreview(content);
                    }}
                    class="w-full text-left p-3 rounded-xl border border-slate-200 hover:border-slate-300 hover:bg-slate-50 transition-all duration-150"
                  >
                    <div class="flex items-center gap-3">
                      <div class="flex-shrink-0 w-12 h-12 bg-slate-100 rounded-lg overflow-hidden">
                        {#if item.thumbnail}
                          <img src={item.thumbnail} alt="" class="w-full h-full object-cover" />
                        {:else}
                          <div class="w-full h-full flex items-center justify-center text-slate-400">
                            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                            </svg>
                          </div>
                        {/if}
                      </div>
                      <div class="flex-1 min-w-0">
                        <p class="font-medium text-slate-900 truncate">{item.title}</p>
                        <p class="text-xs text-slate-500 capitalize">{item.type === 'video' ? 'Video' : 'Artikel'}</p>
                      </div>
                    </div>
                  </button>
                {/each}
              </div>
            </div>
          {/if}
          {#if availablePopularContent.length > 0}
            <div>
              <h4 class="text-sm font-medium text-slate-500 mb-3 uppercase tracking-wide">
                {$t('content.picker.popularContent') || 'Konten Populer'}
              </h4>
              {#if isLoadingPopular}
                <div class="flex items-center justify-center py-4">
                  <svg class="animate-spin h-5 w-5 text-teal-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </div>
              {:else}
                <div class="space-y-3">
                  {#each availablePopularContent as item (item.id + '-' + item.type)}
                    <button
                      onclick={() => {
                        const content = item.type === 'video'
                          ? videos.find(v => v.id === item.id)
                          : articles.find(a => a.id === item.id);
                        if (content) openPreview(content);
                      }}
                      class="w-full text-left p-3 rounded-xl border border-slate-200 hover:border-slate-300 hover:bg-slate-50 transition-all duration-150"
                    >
                      <div class="flex items-center gap-3">
                        <div class="flex-shrink-0 w-12 h-12 bg-slate-100 rounded-lg overflow-hidden">
                          {#if item.thumbnail}
                            <img src={item.thumbnail} alt="" class="w-full h-full object-cover" />
                          {:else}
                            <div class="w-full h-full flex items-center justify-center text-slate-400">
                              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                              </svg>
                            </div>
                          {/if}
                        </div>
                        <div class="flex-1 min-w-0">
                          <p class="font-medium text-slate-900 truncate">{item.title}</p>
                          <p class="text-xs text-slate-500 capitalize">{item.type === 'video' ? 'Video' : 'Artikel'}</p>
                        </div>
                      </div>
                    </button>
                  {/each}
                </div>
              {/if}
            </div>
          {/if}
        {:else}
          <!-- Articles -->
          {#if displayedArticles.length > 0}
            <div class="mb-6">
              <h3 class="text-sm font-medium text-slate-500 mb-3 uppercase tracking-wide">
                {$t('content.picker.articles') || 'Artikel'}
              </h3>
              <div class="space-y-3">
                {#each displayedArticles as article (article.id)}
                  <button
                    onclick={(e) => {
                      if (e.ctrlKey || e.metaKey) {
                        handleSelect(article);
                      } else {
                        openPreview(article, e);
                      }
                    }}
                    onkeydown={(e) => {
                      if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        openPreview(article, e);
                      }
                    }}
                    class="w-full text-left p-3 rounded-xl border transition-all duration-150 cursor-pointer
                      {isSelected(article)
                        ? 'border-teal-500 bg-teal-50'
                        : 'border-slate-200 hover:border-slate-300 hover:bg-slate-50 focus:ring-2 focus:ring-teal-500 focus:border-transparent active:bg-slate-100'}"
                    aria-pressed={isSelected(article)}
                    aria-label="Pratinjau {article.title}"
                    title="Klik untuk pratinjau, Ctrl+klik untuk langsung memilih"
                  >
                    <div class="flex items-start gap-3">
                      <div class="flex-shrink-0 w-16 h-16 bg-slate-100 rounded-lg overflow-hidden">
                        {#if article.heroImages?.hero1x1}
                          <img
                            src={article.heroImages.hero1x1}
                            alt=""
                            class="w-full h-full object-cover"
                          />
                        {:else}
                          <div class="w-full h-full flex items-center justify-center text-slate-400">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                            </svg>
                          </div>
                        {/if}
                      </div>
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-1">
                          <span class="inline-block px-2 py-0.5 bg-teal-50 text-teal-700 text-xs font-medium rounded">
                            {$t('content.picker.articles') || 'Artikel'}
                          </span>
                          {#if article.categoryName}
                            <span class="text-xs text-slate-500">{article.categoryName}</span>
                          {/if}
                        </div>
                        <h4 class="font-medium text-slate-900 truncate">{article.title}</h4>
                        {#if article.excerpt}
                          <p class="text-sm text-slate-500 line-clamp-2 mt-1">{article.excerpt}</p>
                        {/if}
                      </div>
                      <div class="flex-shrink-0">
                        {#if isSelected(article)}
                          <svg class="w-6 h-6 text-teal-500" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                          </svg>
                        {:else}
                          <div class="w-6 h-6 rounded-full border-2 border-slate-300"></div>
                        {/if}
                      </div>
                    </div>
                  </button>
                {/each}
              </div>
            </div>
          {/if}

          <!-- Videos -->
          {#if displayedVideos.length > 0}
            <div>
              <h3 class="text-sm font-medium text-slate-500 mb-3 uppercase tracking-wide">
                {$t('content.picker.videos') || 'Video'}
              </h3>
              <div class="space-y-3">
                {#each displayedVideos as video (video.id)}
                  <button
                    onclick={(e) => {
                      if (e.ctrlKey || e.metaKey) {
                        handleSelect(video);
                      } else {
                        openPreview(video, e);
                      }
                    }}
                    onkeydown={(e) => {
                      if (e.key === 'Enter' || e.key === ' ') {
                        e.preventDefault();
                        openPreview(video, e);
                      }
                    }}
                    class="w-full text-left p-3 rounded-xl border transition-all duration-150 cursor-pointer
                      {isSelected(video)
                        ? 'border-teal-500 bg-teal-50'
                        : 'border-slate-200 hover:border-slate-300 hover:bg-slate-50 focus:ring-2 focus:ring-teal-500 focus:border-transparent active:bg-slate-100'}"
                    aria-pressed={isSelected(video)}
                    aria-label="Pratinjau {video.title}"
                    title="Klik untuk pratinjau, Ctrl+klik untuk langsung memilih"
                  >
                    <div class="flex items-start gap-3">
                      <div class="flex-shrink-0 w-24 h-16 bg-slate-100 rounded-lg overflow-hidden relative">
                        {#if video.thumbnailURL}
                          <img
                            src={video.thumbnailURL}
                            alt=""
                            class="w-full h-full object-cover"
                          />
                          {#if video.duration}
                            <div class="absolute bottom-1 right-1 px-1 bg-black/70 text-white text-xs rounded">
                              {video.duration}
                            </div>
                          {/if}
                        {:else}
                          <div class="w-full h-full flex items-center justify-center text-slate-400">
                            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                          </div>
                        {/if}
                      </div>
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-1">
                          <span class="inline-block px-2 py-0.5 bg-purple-50 text-purple-700 text-xs font-medium rounded">
                            {$t('content.picker.videos') || 'Video'}
                          </span>
                        </div>
                        <h4 class="font-medium text-slate-900 truncate">{video.title}</h4>
                        {#if video.channelName}
                          <p class="text-sm text-slate-500 mt-1">{video.channelName}</p>
                        {/if}
                      </div>
                      <div class="flex-shrink-0">
                        {#if isSelected(video)}
                          <svg class="w-6 h-6 text-teal-500" fill="currentColor" viewBox="0 0 20 20">
                            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                          </svg>
                        {:else}
                          <div class="w-6 h-6 rounded-full border-2 border-slate-300"></div>
                        {/if}
                      </div>
                    </div>
                  </button>
                {/each}
              </div>
            </div>
          {/if}
        {/if}
      </div>

      <!-- Footer -->
      <div class="flex-shrink-0 p-4 sm:p-6 border-t border-slate-100 bg-slate-50">
        <div class="flex flex-col sm:flex-row gap-3">
          <button
            onclick={handleClose}
            class="flex-1 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-100 transition-colors"
          >
            {$t('common.cancel') || 'Batal'}
          </button>
          <button
            onclick={() => {
              onSelect(selectedContent);
              handleClose();
            }}
            disabled={selectedContent.length === 0}
            class="flex-1 px-4 py-2.5 bg-teal-500 text-white font-medium rounded-xl hover:bg-teal-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <span>{$t('content.picker.attach') || 'Lampirkan'} ({selectedContent.length})</span>
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Content Preview Panel -->
  <ContentPreviewPanel
    content={previewContent}
    isSelected={previewContent ? isSelected(previewContent) : false}
    onClose={closePreview}
    onAttach={handleAttachFromPreview}
  />
{/if}
