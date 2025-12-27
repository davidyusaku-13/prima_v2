<script>
  import { t } from 'svelte-i18n';

  export let video;
  export let onClick = () => {};
  export let showActions = false;
  export let onDelete = () => {};

  function getYouTubeVideoId(url) {
    if (!url) return null;
    const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|&v=)([^#&?]*).*/;
    const match = url.match(regExp);
    return (match && match[2].length === 11) ? match[2] : null;
  }

  function getYouTubeThumbnail(videoId) {
    return `https://img.youtube.com/vi/${videoId}/hqdefault.jpg`;
  }

  function formatDuration(minutes) {
    if (!minutes) return '';
    if (minutes < 60) return `${minutes} ${$t('video.minutes', { n: minutes })}`;
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
  }

  function getCategoryLabel(category) {
    const categories = {
      tutorial: $t('videoCategories.tutorial'),
      firstAid: $t('videoCategories.firstAid'),
      exercise: $t('videoCategories.exercise'),
      nutrition: $t('videoCategories.nutrition'),
      mentalHealth: $t('videoCategories.mentalHealth'),
      childHealth: $t('videoCategories.childHealth'),
      seniorHealth: $t('videoCategories.seniorHealth')
    };
    return categories[category] || category || $t('videoCategories.all');
  }

  function getCategoryColor(category) {
    const colors = {
      tutorial: 'bg-blue-100 text-blue-700',
      firstAid: 'bg-red-100 text-red-700',
      exercise: 'bg-green-100 text-green-700',
      nutrition: 'bg-amber-100 text-amber-700',
      mentalHealth: 'bg-pink-100 text-pink-700',
      childHealth: 'bg-cyan-100 text-cyan-700',
      seniorHealth: 'bg-purple-100 text-purple-700'
    };
    return colors[category] || 'bg-slate-100 text-slate-700';
  }

  $: videoId = getYouTubeVideoId(video.youtube_url || video.url || video.youtubeUrl);
  $: thumbnail = video.thumbnail_url || video.thumbnail || (videoId ? getYouTubeThumbnail(videoId) : '');
  $: categoryId = video.category_id || video.category || '';
</script>

<div class="group bg-white rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg transition-all duration-300">
  <!-- Video Thumbnail -->
  <button onclick={onClick} class="w-full relative aspect-video bg-slate-100 overflow-hidden">
    {#if thumbnail}
      <img
        src={thumbnail}
        alt={video.title}
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
      />
    {:else}
      <div class="w-full h-full flex items-center justify-center text-slate-300">
        <svg class="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
    {/if}
    <!-- Play Button Overlay -->
    <div class="absolute inset-0 flex items-center justify-center bg-black/20 group-hover:bg-black/30 transition-colors">
      <div class="w-14 h-14 bg-white/90 rounded-full flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform">
        <svg class="w-6 h-6 text-teal-600 ml-1" fill="currentColor" viewBox="0 0 24 24">
          <path d="M8 5v14l11-7z" />
        </svg>
      </div>
    </div>
    <!-- Duration Badge -->
    {#if video.duration}
      <span class="absolute bottom-2 right-2 px-2 py-0.5 bg-black/80 text-white text-xs font-medium rounded">
        {formatDuration(video.duration)}
      </span>
    {/if}
  </button>

  <!-- Content -->
  <div class="p-4">
    <div class="flex items-start gap-3">
      <div class="flex-1 min-w-0">
        <h3 class="font-semibold text-slate-900 group-hover:text-teal-700 line-clamp-2 leading-snug transition-colors">
          {video.title}
        </h3>
        {#if video.description}
          <p class="mt-1 text-sm text-slate-500 line-clamp-2 leading-relaxed">
            {video.description}
          </p>
        {/if}
        <div class="mt-2 flex items-center gap-2">
          {#if categoryId}
            <span class="px-2 py-0.5 text-xs font-medium rounded-full {getCategoryColor(categoryId)}">
              {getCategoryLabel(categoryId)}
            </span>
          {/if}
        </div>
      </div>
    </div>

    <!-- Actions -->
    {#if showActions}
      <div class="mt-3 pt-3 border-t border-slate-100 flex items-center justify-end gap-2">
        <button
          onclick={onClick}
          class="px-3 py-1.5 text-xs font-medium text-slate-600 hover:text-teal-600 hover:bg-teal-50 rounded-lg transition-colors"
        >
          {$t('common.edit')}
        </button>
        <button
          onclick={onDelete}
          class="px-3 py-1.5 text-xs font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors"
        >
          {$t('common.delete')}
        </button>
      </div>
    {/if}
  </div>
</div>
