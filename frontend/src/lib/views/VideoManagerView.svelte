<script>
  import { t } from 'svelte-i18n';
  import * as api from '$lib/utils/api.js';

  export let onClose = () => {};
  export let onSave = () => {};
  export let token = null;

  let youtubeUrl = '';
  let videoTitle = '';
  let videoDescription = '';
  let videoCategory = '';
  let videoThumbnail = '';

  let saving = false;
  let error = null;
  let fetchingPreview = false;

  const categories = ['tutorial', 'firstAid', 'exercise', 'nutrition', 'mentalHealth', 'childHealth', 'seniorHealth'];

  function getYouTubeVideoId(url) {
    if (!url) return null;
    const regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|&v=)([^#&?]*).*/;
    const match = url.match(regExp);
    return (match && match[2].length === 11) ? match[2] : null;
  }

  async function fetchYouTubePreview() {
    const videoId = getYouTubeVideoId(youtubeUrl);
    if (!videoId) {
      error = $t('videoManager.errorInvalidUrl');
      return;
    }

    fetchingPreview = true;
    error = null;

    try {
      const response = await fetch(`https://noembed.com/embed?url=https://www.youtube.com/watch?v=${videoId}`);
      const data = await response.json();

      if (data.title) {
        videoTitle = data.title;
      }
      videoThumbnail = data.thumbnail_url || `https://img.youtube.com/vi/${videoId}/hqdefault.jpg`;
    } catch (e) {
      videoThumbnail = `https://img.youtube.com/vi/${videoId}/hqdefault.jpg`;
      error = $t('videoManager.errorPreview');
    } finally {
      fetchingPreview = false;
    }
  }

  function clearForm() {
    youtubeUrl = '';
    videoTitle = '';
    videoDescription = '';
    videoCategory = '';
    videoThumbnail = '';
    error = null;
  }

  async function handleSave() {
    if (!videoTitle.trim()) {
      error = $t('videoManager.titleRequired');
      return;
    }

    const videoId = getYouTubeVideoId(youtubeUrl);
    if (!videoId) {
      error = $t('videoManager.errorInvalidUrl');
      return;
    }

    saving = true;
    error = null;

    try {
      const videoData = {
        youtube_url: youtubeUrl,
        category_id: videoCategory
      };

      await api.createVideo(token, videoData);

      clearForm();
      onSave();
    } catch (e) {
      error = e.message || $t('videoManager.errorSave');
    } finally {
      saving = false;
    }
  }

  function getCategoryLabel(cat) {
    return $t(`videoCategories.${cat}`);
  }
</script>

<div
  class="fixed inset-0 bg-black/50 z-50 overflow-y-auto"
  onclick={onClose}
  role="button"
  tabindex="0"
  aria-label="Close modal"
  onkeydown={(e) => e.key === 'Escape' && onClose()}
>
  <div
    class="min-h-screen flex items-start justify-center p-4 pt-10"
    onclick={(e) => e.stopPropagation()}
    onkeydown={(e) => e.key === 'Escape' && onClose()}
    role="presentation"
  >
    <div
      class="w-full max-w-lg bg-white rounded-2xl shadow-2xl"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-200">
        <h2 class="text-xl font-bold text-slate-900">{$t('videoManager.addVideo')}</h2>
        <button
          onclick={onClose}
          aria-label="Close video manager"
          class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <div class="p-6 space-y-4">
        {#if error}
          <div class="p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">
            {error}
          </div>
        {/if}

        <!-- YouTube URL -->
        <label class="block text-sm font-medium text-slate-700 mb-1">
          {$t('videoManager.youtubeUrl')}
          <span class="text-red-500">*</span>
          <div class="flex gap-2 mt-1">
            <input
              type="text"
              bind:value={youtubeUrl}
              placeholder={$t('videoManager.urlPlaceholder')}
              class="flex-1 px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500"
            />
            <button
              onclick={fetchYouTubePreview}
              disabled={fetchingPreview}
              class="px-4 py-2 bg-slate-100 text-slate-700 rounded-xl font-medium hover:bg-slate-200 transition-colors disabled:opacity-50"
            >
              {fetchingPreview ? $t('common.loading') : $t('videoManager.autoFetch')}
            </button>
          </div>
        </label>

        <!-- Thumbnail Preview -->
        {#if videoThumbnail}
          <div class="aspect-video bg-slate-100 rounded-xl overflow-hidden">
            <img src={videoThumbnail} alt="Thumbnail" class="w-full h-full object-cover" />
          </div>
        {/if}

        <!-- Title -->
        <label class="block text-sm font-medium text-slate-700 mb-1">
          {$t('videoManager.videoTitle')}
          <span class="text-red-500">*</span>
          <input
            type="text"
            bind:value={videoTitle}
            placeholder={$t('videoManager.titlePlaceholder')}
            class="w-full px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500"
          />
        </label>

        <!-- Description -->
        <label class="block text-sm font-medium text-slate-700 mb-1">
          {$t('videoManager.description')}
          <textarea
            bind:value={videoDescription}
            placeholder={$t('videoManager.descriptionPlaceholder')}
            rows="3"
            class="w-full px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 resize-none"
          ></textarea>
        </label>

        <!-- Category -->
        <label class="block text-sm font-medium text-slate-700 mb-1">
          {$t('videoManager.category')}
          <select
            bind:value={videoCategory}
            class="w-full px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500"
          >
            <option value="">{$t('videoManager.selectCategory')}</option>
            {#each categories as cat}
              <option value={cat}>{getCategoryLabel(cat)}</option>
            {/each}
          </select>
        </label>

        <!-- Actions -->
        <div class="flex gap-3 pt-4">
          <button
            onclick={clearForm}
            class="flex-1 px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-xl font-medium transition-colors"
          >
            {$t('common.clear')}
          </button>
          <button
            onclick={handleSave}
            disabled={saving}
            class="px-6 py-2 bg-teal-600 text-white rounded-xl font-medium hover:bg-teal-700 transition-colors disabled:opacity-50"
          >
            {saving ? $t('common.loading') : $t('common.save')}
          </button>
        </div>
      </div>
    </div>
  </div>
</div>
