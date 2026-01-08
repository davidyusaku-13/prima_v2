<script>
  import { t, locale } from 'svelte-i18n';

  export let article;
  export let onClick = () => {};

  const BACKEND_URL = 'http://localhost:8080';

  // Helper to get full image URL
  function getFullImageUrl(path) {
    if (!path) return '';
    if (path.startsWith('http')) return path;
    return BACKEND_URL + path;
  }

  function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString($locale || 'en-US', { day: 'numeric', month: 'short', year: 'numeric' });
  }

  function getCategoryLabel(category) {
    return $t(`articleCategories.${category}`);
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

  // Handle both old format (heroImage) and new format (hero_images.hero_16x9)
  $: heroImageUrl = article.hero_images?.hero_16x9 || article.heroImage || '';
  $: categoryId = article.category_id || article.category || '';
</script>

<button
  onclick={onClick}
  class="group w-full text-left bg-white rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg hover:border-teal-200 transition-all duration-300"
>
  <!-- Hero Image -->
  <div class="relative aspect-[16/9] bg-slate-100 overflow-hidden">
    {#if heroImageUrl}
      <img
        src={getFullImageUrl(heroImageUrl)}
        alt={article.title}
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
      />
    {:else}
      <div class="w-full h-full flex items-center justify-center text-slate-300">
        <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      </div>
    {/if}
    <!-- Category Badge -->
    <span class="absolute top-3 left-3 px-2.5 py-1 text-xs font-medium rounded-full {getCategoryColor(categoryId)}">
      {getCategoryLabel(categoryId)}
    </span>
  </div>

  <!-- Content -->
  <div class="p-4">
    <h3 class="font-semibold text-slate-900 group-hover:text-teal-700 line-clamp-2 leading-snug transition-colors">
      {article.title}
    </h3>
    {#if article.excerpt}
      <p class="mt-2 text-sm text-slate-500 line-clamp-2 leading-relaxed">
        {article.excerpt}
      </p>
    {/if}
    <div class="mt-3 flex items-center gap-2 text-xs text-slate-400">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
      <span>{formatDate(article.createdAt || article.publishedAt)}</span>
    </div>
  </div>
</button>
