<script>
  import { t } from 'svelte-i18n';
  import VideoCard from '$lib/components/VideoCard.svelte';
  import * as api from '$lib/utils/api.js';

  export let onClose = () => {};
  export let onSave = () => {};
  export const token = null;

  let videos = [];
  let loading = true;

  let editingVideo = null;
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
      videoThumbnail = `https://img.youtube.com/vi/${videoId}/hqdefault.jpg`;
      if (!videoTitle) {
        videoTitle = $t('videoManager.youtubeVideo', { id: videoId });
      }
    } catch (e) {
      error = $t('videoManager.errorPreview');
    } finally {
      fetchingPreview = false;
    }
  }

  function clearForm() {
    editingVideo = null;
    youtubeUrl = '';
    videoTitle = '';
    videoDescription = '';
    videoCategory = '';
    videoThumbnail = '';
    error = null;
  }

  function editVideo(video) {
    editingVideo = video;
    youtubeUrl = video.youtube_url || '';
    videoTitle = video.title;
    videoDescription = video.description || '';
    videoCategory = video.category_id || '';
    videoThumbnail = video.thumbnail_url || '';
  }

  async function handleSave() {
    if (!videoTitle.trim()) {
      error = $t('videoManager.titleRequired');
      return;
    }

    const videoId = getYouTubeVideoId(youtubeUrl);
    if (!videoId && !editingVideo) {
      error = $t('videoManager.errorInvalidUrl');
      return;
    }

    saving = true;
    error = null;

    try {
      // Backend expects youtube_url and category_id
      const videoData = {
        youtube_url: youtubeUrl,
        category_id: videoCategory
      };

      if (editingVideo?.id) {
        // For update, we send what changed (title, description)
        await api.updateVideo(null, editingVideo.id, {
          title: videoTitle.trim(),
          description: videoDescription
        });
      } else {
        // For create, we need youtube_url and category_id
        await api.createVideo(null, videoData);
      }

      clearForm();
      loadVideos();
      onSave();
    } catch (e) {
      error = e.message || $t('videoManager.errorSave');
    } finally {
      saving = false;
    }
  }

  async function handleDelete(video) {
    if (!confirm($t('videoManager.deleteConfirm'))) return;

    try {
      await api.deleteVideo(null, video.id);
      loadVideos();
      onSave();
    } catch (e) {
      console.error($t('common.error'), e);
    }
  }

  async function loadVideos() {
    loading = true;
    try {
      videos = await api.fetchVideos();
    } catch (e) {
      console.error($t('common.errorLoading'), e);
      videos = [];
    } finally {
      loading = false;
    }
  }

  function getCategoryLabel(cat) {
    const labels = {
      tutorial: $t('videoCategories.tutorial'),
      firstAid: $t('videoCategories.firstAid'),
      exercise: $t('videoCategories.exercise'),
      nutrition: $t('videoCategories.nutrition'),
      mentalHealth: $t('videoCategories.mentalHealth'),
      childHealth: $t('videoCategories.childHealth'),
      seniorHealth: $t('videoCategories.seniorHealth')
    };
    return labels[cat] || cat;
  }

  loadVideos();
</script>

<div class="fixed inset-0 bg-black/50 z-50 overflow-y-auto" onclick={onClose}>
  <div class="min-h-screen flex items-start justify-center p-4 pt-10">
    <div class="w-full max-w-5xl bg-white rounded-2xl shadow-2xl" onclick={(e) => e.stopPropagation()}>
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-200">
        <h2 class="text-xl font-bold text-slate-900">{$t('videoManager.title')}</h2>
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

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 p-6">
        <!-- Add/Edit Form -->
        <div class="space-y-4">
          <h3 class="font-semibold text-slate-900">
            {editingVideo ? $t('videoManager.editVideo') : $t('videoManager.addVideo')}
          </h3>

          {#if error}
            <div class="p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">
              {error}
            </div>
          {/if}

          <!-- YouTube URL -->
          <label class="block text-sm font-medium text-slate-700 mb-1">
            {$t('videoManager.youtubeUrl')}
            {#if !editingVideo}<span class="text-red-500">*</span>{/if}
            <div class="flex gap-2 mt-1">
              <input
                type="text"
                bind:value={youtubeUrl}
                placeholder={$t('videoManager.urlPlaceholder')}
                disabled={!!editingVideo}
                class="flex-1 px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500"
              />
              {#if !editingVideo}
                <button
                  onclick={fetchYouTubePreview}
                  disabled={fetchingPreview}
                  class="px-4 py-2 bg-slate-100 text-slate-700 rounded-xl font-medium hover:bg-slate-200 transition-colors disabled:opacity-50"
                >
                  {fetchingPreview ? $t('common.loading') : $t('videoManager.autoFetch')}
                </button>
              {/if}
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
            {#if editingVideo}
              <button
                onclick={clearForm}
                class="flex-1 px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-xl font-medium transition-colors"
              >
                {$t('common.cancel')}
              </button>
            {:else}
              <div class="flex-1"></div>
            {/if}
            <button
              onclick={handleSave}
              disabled={saving}
              class="px-6 py-2 bg-teal-600 text-white rounded-xl font-medium hover:bg-teal-700 transition-colors disabled:opacity-50"
            >
              {saving ? $t('common.loading') : $t('common.save')}
            </button>
          </div>
        </div>

        <!-- Video List -->
        <div class="space-y-4">
          <h3 class="font-semibold text-slate-900">{$t('cms.existingVideos', { n: videos.length })}</h3>

          {#if loading}
            <div class="space-y-4">
              {#each Array(3) as _}
                <div class="bg-slate-100 rounded-xl p-4 animate-pulse">
                  <div class="flex gap-3">
                    <div class="w-32 h-20 bg-slate-200 rounded"></div>
                    <div class="flex-1 space-y-2">
                      <div class="h-4 bg-slate-200 rounded w-3/4"></div>
                      <div class="h-3 bg-slate-200 rounded w-full"></div>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {:else if videos.length > 0}
            <div class="space-y-3 max-h-[600px] overflow-y-auto">
              {#each videos as video}
                <VideoCard
                  {video}
                  showActions={true}
                  onClick={() => editVideo(video)}
                  onDelete={() => handleDelete(video)}
                />
              {/each}
            </div>
          {:else}
            <div class="text-center py-12 text-slate-500">
              <svg class="w-12 h-12 mx-auto text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              <p class="mt-4">{$t('video.noVideos')}</p>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
