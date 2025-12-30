<script>
  import { t } from 'svelte-i18n';
  import * as api from '$lib/utils/api.js';

  export let video = null;
  export let onClose = () => {};
  export let onSave = () => {};
  export let token = null;

  let title = video?.title || '';
  let description = video?.description || '';
  let saving = false;
  let error = null;

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

  async function handleSave() {
    if (!title.trim()) {
      error = $t('videoManager.titleRequired');
      return;
    }

    saving = true;
    error = null;

    try {
      await api.updateVideo(token, video.id, {
        title: title.trim(),
        description
      });
      onSave();
      onClose();
    } catch (e) {
      error = e.message || $t('videoManager.errorSave');
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!confirm($t('videoManager.deleteConfirm'))) return;

    saving = true;
    try {
      await api.deleteVideo(token, video.id);
      onSave();
      onClose();
    } catch (e) {
      error = e.message || $t('videoManager.errorDelete');
      saving = false;
    }
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
        <h2 class="text-xl font-bold text-slate-900">{$t('videoManager.editVideo')}</h2>
        <button
          onclick={onClose}
          aria-label="Close editor"
          class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Form -->
      <div class="p-6 space-y-4">
        {#if error}
          <div class="p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">
            {error}
          </div>
        {/if}

        <!-- Thumbnail Preview -->
        {#if video?.thumbnail_url}
          <div class="aspect-video bg-slate-100 rounded-xl overflow-hidden">
            <img src={video.thumbnail_url} alt="Thumbnail" class="w-full h-full object-cover" />
          </div>
        {/if}

        <!-- YouTube URL (read-only) -->
        <div>
          <label for="youtube-url" class="block text-sm font-medium text-slate-700 mb-1">
            {$t('videoManager.youtubeUrl')}
          </label>
          <input
            id="youtube-url"
            type="text"
            value={video?.youtube_url || ''}
            disabled
            class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl text-slate-500"
          />
        </div>

        <!-- Title -->
        <div>
          <label for="video-title" class="block text-sm font-medium text-slate-700 mb-1">
            {$t('videoManager.videoTitle')}
            <span class="text-red-500">*</span>
          </label>
          <input
            id="video-title"
            type="text"
            bind:value={title}
            placeholder={$t('videoManager.titlePlaceholder')}
            class="w-full px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500"
          />
        </div>

        <!-- Description -->
        <div>
          <label for="video-description" class="block text-sm font-medium text-slate-700 mb-1">
            {$t('videoManager.description')}
          </label>
          <textarea
            id="video-description"
            bind:value={description}
            placeholder={$t('videoManager.descriptionPlaceholder')}
            rows="3"
            class="w-full px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 resize-none"
          ></textarea>
        </div>

        <!-- Category (read-only) -->
        <div>
          <label for="video-category" class="block text-sm font-medium text-slate-700 mb-1">
            {$t('videoManager.category')}
          </label>
          <input
            id="video-category"
            type="text"
            value={video?.category_id ? getCategoryLabel(video.category_id) : ''}
            disabled
            class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl text-slate-500"
          />
        </div>
      </div>

      <!-- Footer -->
      <div class="flex items-center justify-between px-6 py-4 border-t border-slate-200 bg-slate-50 rounded-b-2xl">
        <button
          onclick={handleDelete}
          disabled={saving}
          class="px-4 py-2 text-red-600 hover:bg-red-50 rounded-xl font-medium transition-colors"
        >
          {$t('common.delete')}
        </button>
        <div class="flex gap-3">
          <button
            onclick={onClose}
            disabled={saving}
            class="px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-xl font-medium transition-colors"
          >
            {$t('common.cancel')}
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
