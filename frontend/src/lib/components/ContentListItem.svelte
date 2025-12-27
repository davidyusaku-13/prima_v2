<script>
  import { t, locale } from 'svelte-i18n';

  export let item;
  export let type; // 'article' or 'video'
  export let selected = false;
  export let onSelect = () => {};
  export let onEdit = () => {};
  export let onDelete = () => {};
  export let onToggleStatus = () => {};

  const BACKEND_URL = 'http://localhost:8080';

  function getFullImageUrl(path) {
    if (!path) return '';
    if (path.startsWith('http')) return path;
    return BACKEND_URL + path;
  }

  function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString($locale || 'en-US', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    });
  }

  function getCategoryLabel(categoryId) {
    if (type === 'article') {
      const categories = {
        latest: $t('articleCategories.latest'),
        policy: $t('articleCategories.policy'),
        research: $t('articleCategories.research'),
        outbreak: $t('articleCategories.outbreak'),
        lifestyle: $t('articleCategories.lifestyle'),
        local: $t('articleCategories.local')
      };
      return categories[categoryId] || categoryId || '-';
    } else {
      const categories = {
        tutorial: $t('videoCategories.tutorial'),
        firstAid: $t('videoCategories.firstAid'),
        exercise: $t('videoCategories.exercise'),
        nutrition: $t('videoCategories.nutrition'),
        mentalHealth: $t('videoCategories.mentalHealth'),
        childHealth: $t('videoCategories.childHealth'),
        seniorHealth: $t('videoCategories.seniorHealth')
      };
      return categories[categoryId] || categoryId || '-';
    }
  }

  $: thumbnail = type === 'article'
    ? (item.hero_images?.hero_16x9 || item.heroImage || '')
    : (item.thumbnail_url || item.thumbnail || '');

  $: categoryId = item.category_id || item.category || '';
  $: displayDate = formatDate(item.published_at || item.publishedAt || item.created_at || item.createdAt);
  $: viewCount = item.view_count || item.viewCount || 0;
  $: status = item.status || 'published';
  $: isDraft = status === 'draft';
</script>

<div class="group flex items-center gap-4 p-4 bg-white border-b border-slate-100 hover:bg-slate-50 transition-colors">
  <!-- Checkbox -->
  <input
    type="checkbox"
    checked={selected}
    onchange={onSelect}
    class="w-4 h-4 text-teal-600 border-slate-300 rounded focus:ring-teal-500 focus:ring-2"
  />

  <!-- Type Icon -->
  <div class="flex-shrink-0 w-8 h-8 flex items-center justify-center text-slate-400">
    {#if type === 'article'}
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
    {:else}
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
      </svg>
    {/if}
  </div>

  <!-- Thumbnail -->
  <div class="flex-shrink-0 w-20 h-14 bg-slate-100 rounded-lg overflow-hidden">
    {#if thumbnail}
      <img
        src={getFullImageUrl(thumbnail)}
        alt={item.title}
        class="w-full h-full object-cover"
      />
    {:else}
      <div class="w-full h-full flex items-center justify-center text-slate-300">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      </div>
    {/if}
  </div>

  <!-- Title -->
  <div class="flex-1 min-w-0">
    <h3 class="font-medium text-slate-900 truncate group-hover:text-teal-700 transition-colors">
      {item.title}
    </h3>
  </div>

  <!-- Status Badge -->
  <div class="flex-shrink-0">
    {#if isDraft}
      <span class="px-2.5 py-1 text-xs font-medium bg-yellow-100 text-yellow-700 rounded-full">
        {$t('articleEditor.draft')}
      </span>
    {:else}
      <span class="px-2.5 py-1 text-xs font-medium bg-green-100 text-green-700 rounded-full">
        {$t('articleEditor.published')}
      </span>
    {/if}
  </div>

  <!-- Category -->
  <div class="flex-shrink-0 w-32 text-sm text-slate-600 truncate">
    {getCategoryLabel(categoryId)}
  </div>

  <!-- Date -->
  <div class="flex-shrink-0 w-28 text-sm text-slate-500">
    {displayDate}
  </div>

  <!-- Views -->
  <div class="flex-shrink-0 w-20 flex items-center gap-1.5 text-sm text-slate-500">
    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
    </svg>
    <span>{viewCount}</span>
  </div>

  <!-- Actions -->
  <div class="flex-shrink-0 flex items-center gap-1">
    <!-- Toggle Status -->
    {#if type === 'article'}
      <button
        onclick={onToggleStatus}
        class="p-2 text-slate-400 hover:text-teal-600 hover:bg-teal-50 rounded-lg transition-colors"
        title={isDraft ? $t('articleEditor.publish') : 'Unpublish'}
      >
        {#if isDraft}
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        {:else}
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        {/if}
      </button>
    {/if}

    <!-- Edit -->
    <button
      onclick={onEdit}
      class="p-2 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
      title={$t('common.edit')}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
      </svg>
    </button>

    <!-- Delete -->
    <button
      onclick={onDelete}
      class="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
      title={$t('common.delete')}
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
      </svg>
    </button>
  </div>
</div>
