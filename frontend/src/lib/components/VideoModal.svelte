<script>
  import { t } from 'svelte-i18n';

  export let show = false;
  export let video = null;
  export let onClose = () => {};

  function getYouTubeVideoId(url) {
    if (!url) return null;
    const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|&v=)([^#&?]*).*/;
    const match = url.match(regExp);
    return (match && match[2].length === 11) ? match[2] : null;
  }

  function formatDuration(minutes) {
    if (!minutes) return '';
    if (minutes < 60) return $t('video.min', { n: minutes });
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return mins > 0 ? `${hours}h ${mins}m` : `${hours}h`;
  }

  $: videoId = video ? getYouTubeVideoId(video.youtube_url || video.url || video.youtubeUrl) : null;
  $: embedUrl = videoId ? `https://www.youtube.com/embed/${videoId}` : '';
</script>

{#if show && video && embedUrl}
  <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
  <div
    class="fixed inset-0 bg-black/80 z-50 flex items-center justify-center p-4"
    onclick={onClose}
    onkeydown={(e) => e.key === 'Escape' && onClose()}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
  >
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div class="w-full max-w-4xl" onclick={(e) => e.stopPropagation()}>
      <!-- Close Button -->
      <button
        onclick={onClose}
        aria-label="Close video"
        class="absolute top-4 right-4 p-2 text-white/80 hover:text-white hover:bg-white/10 rounded-full transition-colors"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <!-- Video Container (16:9) -->
      <div class="relative aspect-video bg-slate-900 rounded-xl overflow-hidden shadow-2xl">
        <iframe
          src={embedUrl}
          title={video.title}
          class="absolute inset-0 w-full h-full"
          frameborder="0"
          allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
          allowfullscreen
        ></iframe>
      </div>

      <!-- Video Info -->
      <div class="mt-4 text-white">
        <h3 class="text-xl font-bold">{video.title}</h3>
        {#if video.description}
          <p class="mt-2 text-white/70 text-sm leading-relaxed">{video.description}</p>
        {/if}
        {#if video.duration}
          <p class="mt-2 text-white/50 text-sm">
            {$t('video.duration')}: {formatDuration(video.duration)}
          </p>
        {/if}
        <a
          href={video.url || `https://youtube.com/watch?v=${videoId}`}
          target="_blank"
          rel="noopener noreferrer"
          class="inline-flex items-center gap-2 mt-3 text-teal-400 hover:text-teal-300 text-sm font-medium"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
            <path d="M23.498 6.186a3.016 3.016 0 0 0-2.122-2.136C19.505 3.545 12 3.545 12 3.545s-7.505 0-9.377.505A3.017 3.017 0 0 0 .502 6.186C0 8.07 0 12 0 12s0 3.93.502 5.814a3.016 3.016 0 0 0 2.122 2.136c1.871.505 9.376.505 9.376.505s7.505 0 9.377-.505a3.015 3.015 0 0 0 2.122-2.136C24 15.93 24 12 24 12s0-3.93-.502-5.814zM9.545 15.568V8.432L15.818 12l-6.273 3.568z"/>
          </svg>
          {$t('video.watchOnYoutube')}
        </a>
      </div>
    </div>
  </div>
{/if}
