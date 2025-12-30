<script>
  import { createEventDispatcher } from 'svelte';
  import { t } from 'svelte-i18n';
  import * as api from '$lib/utils/api.js';

  export let imageUrl = '';
  export let required = false;
  export let label = '';
  export let token = null;

  const dispatch = createEventDispatcher();
  const BACKEND_URL = 'http://localhost:8080';

  // Helper to get full image URL
  function getFullImageUrl(path) {
    if (!path) return '';
    if (path.startsWith('http')) return path;
    return BACKEND_URL + path;
  }

  let dragging = false;
  let fileInput;
  let uploading = false;

  function handleDragOver(e) {
    e.preventDefault();
    dragging = true;
  }

  function handleDragLeave() {
    dragging = false;
  }

  function handleDrop(e) {
    e.preventDefault();
    dragging = false;
    const files = e.dataTransfer?.files;
    if (files?.length) {
      handleFile(files[0]);
    }
  }

  function handleFileSelect(e) {
    const files = e.target?.files;
    if (files?.length) {
      handleFile(files[0]);
    }
  }

  async function handleFile(file) {
    if (!file.type.startsWith('image/')) {
      alert($t('imageUploader.selectImage'));
      return;
    }
    if (file.size > 5 * 1024 * 1024) {
      alert($t('imageUploader.sizeLimit'));
      return;
    }

    uploading = true;

    try {
      // Upload to backend and get hero images
      const heroImages = await api.uploadImage(token, file);
      imageUrl = heroImages;
      dispatch('change', { file, imageUrl: heroImages });
    } catch (e) {
      alert($t('imageUploader.uploadFailed', { error: e.message }));
    } finally {
      uploading = false;
    }
  }

  function removeImage() {
    imageUrl = '';
    dispatch('change', { file: null, imageUrl: '' });
  }

  function triggerFileSelect() {
    fileInput?.click();
  }
</script>

<div class="space-y-2">
  {#if label}
    <span class="block text-sm font-medium text-slate-700">
      {label}
      {#if required}
        <span class="text-red-500">*</span>
      {/if}
    </span>
  {/if}

  {#if imageUrl && imageUrl.hero_16x9}
    <div class="relative group">
      <div class="relative aspect-[16/9] rounded-xl overflow-hidden bg-slate-100">
        <img src={getFullImageUrl(imageUrl.hero_16x9)} alt="Preview" class="w-full h-full object-cover" />
        <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
          <button
            type="button"
            onclick={triggerFileSelect}
            class="px-4 py-2 bg-white text-slate-900 text-sm font-medium rounded-lg hover:bg-slate-100 transition-colors"
          >
            {$t('imageUploader.change')}
          </button>
          <button
            type="button"
            onclick={removeImage}
            class="px-4 py-2 bg-red-500 text-white text-sm font-medium rounded-lg hover:bg-red-600 transition-colors"
          >
            {$t('imageUploader.remove')}
          </button>
        </div>
      </div>
    </div>
  {:else}
    <button
      type="button"
      onclick={triggerFileSelect}
      ondragover={handleDragOver}
      ondragleave={handleDragLeave}
      ondrop={handleDrop}
      class="w-full aspect-[16/9] border-2 border-dashed rounded-xl transition-all duration-200
        {dragging
          ? 'border-teal-500 bg-teal-50'
          : 'border-slate-300 hover:border-teal-400 hover:bg-slate-50'}"
      disabled={uploading}
    >
      <div class="h-full flex flex-col items-center justify-center gap-3 p-6 text-slate-500">
        {#if uploading}
          <div class="animate-spin w-10 h-10 border-4 border-teal-600 border-t-transparent rounded-full"></div>
          <p class="text-sm font-medium">{$t('imageUploader.uploading')}</p>
        {:else}
          <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <div class="text-center">
            <p class="text-sm font-medium">{$t('imageUploader.clickToUpload')}</p>
            <p class="text-xs text-slate-400 mt-1">{$t('imageUploader.fileTypes')}</p>
          </div>
        {/if}
      </div>
    </button>
    <input
      bind:this={fileInput}
      type="file"
      accept="image/*"
      onchange={handleFileSelect}
      class="hidden"
    />
  {/if}
</div>
