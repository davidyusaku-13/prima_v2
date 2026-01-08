<script>
  import { t } from 'svelte-i18n';
  import ImageUploader from '$lib/components/ImageUploader.svelte';
  import QuillEditor from '$lib/components/QuillEditor.svelte';
  import ConfirmModal from '$lib/components/ConfirmModal.svelte';
  import * as api from '$lib/utils/api.js';

  export let article = null;
  export let onClose = () => {};
  export let onSave = () => {};
  export let token = null;

  let title = article?.title || '';
  let content = article?.content || '';
  let excerpt = article?.excerpt || '';
  let category = article?.category_id || '';
  let status = article?.status || 'draft';
  let heroImage = article?.hero_images || '';
  let slug = article?.slug || '';

  let saving = false;
  let error = null;
  let isPreview = false;
  let showDeleteModal = false;

  const categories = ['latest', 'policy', 'research', 'outbreak', 'lifestyle', 'local'];

  function generateSlug(text) {
    return text
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/(^-|-$)/g, '');
  }

  function handleTitleChange() {
    if (!article && !slug) {
      slug = generateSlug(title);
    }
  }

  function getCategoryLabel(cat) {
    const labels = {
      latest: $t('articleCategories.latest'),
      policy: $t('articleCategories.policy'),
      research: $t('articleCategories.research'),
      outbreak: $t('articleCategories.outbreak'),
      lifestyle: $t('articleCategories.lifestyle'),
      local: $t('articleCategories.local')
    };
    return labels[cat] || cat;
  }

  async function handleImageUpload(file) {
    try {
      const heroImages = await api.uploadImage(token, file);
      return heroImages.hero_16x9;
    } catch (e) {
      console.error('Failed to upload image:', e);
      return null;
    }
  }

  async function handleSave(isPublish = false) {
    if (!title.trim()) {
      error = $t('articleEditor.titleRequired');
      return;
    }
    if (!heroImage || !heroImage.hero_16x9) {
      error = $t('articleEditor.imageRequired');
      return;
    }

    saving = true;
    error = null;

    try {
      const articleData = {
        title: title.trim(),
        content,
        excerpt,
        category_id: category,
        status: isPublish ? 'published' : status,
        hero_images: {
          hero_16x9: heroImage.hero_16x9 || '',
          hero_1x1: heroImage.hero_1x1 || '',
          hero_4x3: heroImage.hero_4x3 || ''
        },
        slug: slug || generateSlug(title)
      };

      if (article?.id) {
        await api.updateArticle(token, article.id, articleData);
      } else {
        await api.createArticle(token, articleData);
      }

      onSave(isPublish ? 'published' : 'draft');
      onClose();
    } catch (e) {
      error = e.message || $t('articleEditor.errorSave');
    } finally {
      saving = false;
    }
  }

  async function handleDelete() {
    if (!article?.id) return;
    showDeleteModal = true;
  }

  async function confirmDelete() {
    if (!article?.id) return;
    showDeleteModal = false;
    saving = true;
    try {
      await api.deleteArticle(token, article.id);
      onSave('deleted');
      onClose();
    } catch (e) {
      error = e.message || $t('articleEditor.errorDelete');
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
    role="presentation"
  >
    <div
      class="w-full max-w-4xl bg-white rounded-2xl shadow-2xl"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => {}}
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-slate-200">
        <div class="flex items-center gap-3">
          <h2 class="text-xl font-bold text-slate-900">
            {article?.id ? $t('articleEditor.editArticle') : $t('articleEditor.newArticle')}
          </h2>
        </div>
        <div class="flex items-center gap-2">
          {#if !article?.id}
            <button
              onclick={() => isPreview = !isPreview}
              class="px-3 py-1.5 text-sm font-medium rounded-lg transition-colors"
              class:bg-teal-100={isPreview}
              class:text-teal-700={isPreview}
              class:bg-slate-100={!isPreview}
              class:text-slate-600={!isPreview}
            >
              {isPreview ? $t('articleEditor.editMode') : $t('articleEditor.previewMode')}
            </button>
          {/if}
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
      </div>

      <!-- Form -->
      <div class="p-6 space-y-6">
        {#if error}
          <div class="p-4 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">
            {error}
          </div>
        {/if}

        {#if isPreview}
          <!-- Preview Mode -->
          <div class="space-y-6">
            <div>
              <h1 class="text-2xl font-bold text-slate-900 mb-2">{title || 'Untitled'}</h1>
              <div class="flex items-center gap-4 text-sm text-slate-500">
                {#if category}
                  <span class="px-2 py-0.5 bg-slate-100 rounded-full">{getCategoryLabel(category)}</span>
                {/if}
                {#if excerpt}
                  <span>{excerpt}</span>
                {/if}
              </div>
            </div>
            {#if heroImage?.hero_16x9}
              <img
                src="http://localhost:8080{heroImage.hero_16x9}"
                alt={title}
                class="w-full aspect-[16/9] object-cover rounded-xl"
              />
            {/if}
            <div class="prose max-w-none">
              {@html content || '<p class="text-slate-400 italic">No content yet...</p>'}
            </div>
          </div>
        {:else}
          <!-- Edit Mode -->
          <!-- Title -->
          <label class="block text-sm font-medium text-slate-700 mb-1">
            {$t('articleEditor.articleTitle')}
            <span class="text-red-500">*</span>
            <input
              type="text"
              bind:value={title}
              oninput={handleTitleChange}
              placeholder={$t('articleEditor.titlePlaceholder')}
              class="w-full px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:border-transparent"
            />
          </label>

          <!-- Hero Image -->
          <ImageUploader
            bind:imageUrl={heroImage}
            label={$t('articleEditor.heroImage')}
            required
            {token}
            on:change={(e) => heroImage = e.detail.imageUrl}
          />

          <!-- Slug & Category -->
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <label class="block text-sm font-medium text-slate-700 mb-1">
              {$t('articleEditor.slug')}
              <input
                type="text"
                bind:value={slug}
                placeholder={$t('articleEditor.slugPlaceholder')}
                class="w-full px-4 py-3 bg-slate-50 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500"
              />
            </label>
            <label class="block text-sm font-medium text-slate-700 mb-1">
              {$t('articleEditor.category')}
              <select
                bind:value={category}
                class="w-full px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500"
              >
                <option value="">{$t('articleEditor.selectCategory')}</option>
                {#each categories as cat}
                  <option value={cat}>{getCategoryLabel(cat)}</option>
                {/each}
              </select>
            </label>
          </div>

          <!-- Excerpt -->
          <label class="block text-sm font-medium text-slate-700 mb-1">
            {$t('articleEditor.excerpt')}
            <textarea
              bind:value={excerpt}
              placeholder={$t('articleEditor.excerptPlaceholder')}
              rows="2"
              class="w-full px-4 py-3 bg-white border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 resize-none"
            ></textarea>
          </label>

          <!-- Content with Quill -->
          <label class="block text-sm font-medium text-slate-700 mb-1">
            {$t('articleEditor.content')}
            <QuillEditor bind:value={content} onUploadImage={handleImageUpload} />
          </label>

          <!-- Status (for existing articles) -->
          {#if article?.id}
            <div class="flex items-center gap-4 p-4 bg-slate-50 rounded-xl">
              <span class="text-sm font-medium text-slate-700">{$t('articleEditor.status')}:</span>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" bind:group={status} value="draft" class="w-4 h-4 text-teal-600" />
                <span class="text-sm text-slate-600">{$t('articleEditor.draft')}</span>
              </label>
              <label class="flex items-center gap-2 cursor-pointer">
                <input type="radio" bind:group={status} value="published" class="w-4 h-4 text-teal-600" />
                <span class="text-sm text-slate-600">{$t('articleEditor.published')}</span>
              </label>
            </div>
          {/if}
        {/if}
      </div>

      <!-- Footer -->
      {#if !isPreview}
        <div class="flex items-center justify-between px-6 py-4 border-t border-slate-200 bg-slate-50 rounded-b-2xl">
          {#if article?.id}
            <button
              onclick={handleDelete}
              disabled={saving}
              class="px-4 py-2 text-red-600 hover:bg-red-50 rounded-xl font-medium transition-colors"
            >
              {$t('articleEditor.deleteArticle')}
            </button>
          {:else}
            <div></div>
          {/if}
          <div class="flex gap-3">
            <button
              onclick={onClose}
              disabled={saving}
              class="px-4 py-2 text-slate-600 hover:bg-slate-100 rounded-xl font-medium transition-colors"
            >
              {$t('common.cancel')}
            </button>
            <button
              onclick={() => handleSave(false)}
              disabled={saving}
              class="px-4 py-2 bg-slate-600 text-white rounded-xl font-medium hover:bg-slate-700 transition-colors disabled:opacity-50"
            >
              {saving ? $t('common.loading') : $t('articleEditor.saveDraft')}
            </button>
            <button
              onclick={() => handleSave(true)}
              disabled={saving}
              class="px-4 py-2 bg-teal-600 text-white rounded-xl font-medium hover:bg-teal-700 transition-colors disabled:opacity-50"
            >
              {saving ? $t('common.loading') : $t('articleEditor.publish')}
            </button>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- Delete Confirmation Modal -->
<ConfirmModal
  show={showDeleteModal}
  message={$t('articleEditor.deleteConfirm')}
  onClose={() => showDeleteModal = false}
  onConfirm={confirmDelete}
/>
