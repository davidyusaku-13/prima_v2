<script>
  import { onMount } from 'svelte';
  import { fade, slide } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import { locale, t } from 'svelte-i18n';
  import ContentDisclaimer from './ContentDisclaimer.svelte';

  let {
    content = null,
    isSelected = false,
    onClose = () => {},
    onAttach = () => {}
  } = $props();

  let panelElement = $state(null);
  let imageError = $state(false);

  // Reset image error when content changes
  $effect(() => {
    if (content) {
      imageError = false;
    }
  });

  // Determine content type
  let contentType = $derived(content?.YouTubeID ? 'video' : 'article');

  // Truncate excerpt to 200 characters
  let truncatedExcerpt = $derived.by(() => {
    if (!content?.excerpt) return '';
    const excerpt = content.excerpt;
    if (excerpt.length <= 200) return excerpt;
    return excerpt.substring(0, 197) + '...';
  });

  // Format publish date using i18n locale
  let formattedDate = $derived.by(() => {
    if (!content?.publishedAt && !content?.createdAt) return '';
    const date = content.publishedAt || content.createdAt;
    try {
      // Use locale from svelte-i18n or fall back to Indonesian
      const localeValue = $locale || 'id-ID';
      return new Date(date).toLocaleDateString(localeValue, {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      });
    } catch {
      return '';
    }
  });

  // Focus panel on mount
  onMount(() => {
    if (panelElement) {
      panelElement.focus();
    }
  });

  // Handle keyboard navigation with focus trap
  function handleKeydown(e) {
    if (e.key === 'Escape') {
      e.preventDefault();
      onClose();
      return;
    }

    // Focus trap: keep focus within panel
    if (e.key === 'Tab' && panelElement) {
      const focusableElements = panelElement.querySelectorAll(
        'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
      );
      const firstElement = focusableElements[0];
      const lastElement = focusableElements[focusableElements.length - 1];

      if (e.shiftKey && document.activeElement === firstElement) {
        e.preventDefault();
        lastElement.focus();
      } else if (!e.shiftKey && document.activeElement === lastElement) {
        e.preventDefault();
        firstElement.focus();
      }
    }
  }

  // Handle click outside - FIXED: Changed role from button to presentation
  function handleBackdropClick(e) {
    if (e.target === e.currentTarget) {
      onClose();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if content}
  <!-- Backdrop -->
  <div
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-900/50 backdrop-blur-sm"
    onclick={handleBackdropClick}
    role="presentation"
    aria-hidden="true"
  >
    <!-- Preview Panel -->
    <div
      bind:this={panelElement}
      class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-lg max-h-[80vh] overflow-hidden outline-none"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.key === 'Escape' && onClose()}
      role="dialog"
      aria-modal="true"
      aria-labelledby="preview-title"
      aria-describedby="preview-description"
      tabindex="-1"
      transition:slide={{ axis: 'y', duration: 200, easing: cubicOut }}
    >
      <span id="preview-description" class="sr-only">
        {#if contentType === 'article'}
          Pratinjau artikel dengan gambar, judul, kutipan, dan informasi penulis.
        {:else}
          Pratinjau video dengan thumbnail, judul, durasi, dan informasi channel.
        {/if}
        Tekan tombol lampirkan untuk menambahkan ke pengingat.
      </span>
      <!-- Close button -->
      <button
        onclick={onClose}
        class="absolute top-3 right-3 p-2 rounded-full hover:bg-slate-100 transition-colors z-10"
        aria-label="Tutup pratinjau"
      >
        <svg class="w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <!-- Article Preview -->
      {#if contentType === 'article'}
        <div class="flex flex-col">
          <!-- Hero Image -->
          {#if content.heroImages?.hero16x9 && !imageError}
            <div class="w-full h-48 bg-slate-100 overflow-hidden">
              <img
                src={content.heroImages.hero16x9}
                alt={content.title}
                class="w-full h-full object-cover"
                onerror={() => imageError = true}
              />
            </div>
          {:else if imageError}
            <div class="w-full h-48 bg-slate-100 flex items-center justify-center">
              <svg class="w-12 h-12 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            </div>
          {/if}

          <!-- Content -->
          <div class="p-5">
            <span class="inline-block px-2 py-1 bg-teal-50 text-teal-700 text-xs font-medium rounded-lg mb-3">
              {t('content.picker.articles') || 'Artikel'}
            </span>

            <h2 id="preview-title" class="text-lg font-semibold text-slate-900 mb-2 pr-8">
              {content.title}
            </h2>

            {#if content.excerpt}
              <p class="text-sm text-slate-600 mb-3 leading-relaxed">
                {truncatedExcerpt}
              </p>
            {/if}

            <!-- Attribution Section -->
            <div class="flex flex-wrap gap-3 text-xs text-slate-500 mb-3">
              {#if content.authorName}
                <span class="flex items-center gap-1" aria-label={$t('content.attribution.author')}>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  {content.authorName}
                </span>
              {/if}
              {#if formattedDate}
                <span class="flex items-center gap-1" aria-label={$t('content.attribution.publishedOn')}>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  {formattedDate}
                </span>
              {/if}
            </div>

            <!-- Health Disclaimer -->
            <ContentDisclaimer />

            <!-- Attach Button -->
            <button
              onclick={onAttach}
              class="w-full py-2.5 px-4 rounded-xl font-medium transition-colors flex items-center justify-center gap-2 mt-4
                {isSelected
                  ? 'bg-teal-100 text-teal-700 hover:bg-teal-200'
                  : 'bg-teal-500 text-white hover:bg-teal-600'}"
              aria-pressed={isSelected}
            >
              {#if isSelected}
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <span>{$t('content.preview.selected') || 'Dipilih'}</span>
              {:else}
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                <span>{$t('content.preview.attach') || 'Lampirkan'}</span>
              {/if}
            </button>
          </div>
        </div>

      <!-- Video Preview -->
      {:else if contentType === 'video'}
        <div class="flex flex-col">
          <!-- Thumbnail -->
          <div class="w-full h-48 bg-slate-100 relative overflow-hidden">
            {#if content.thumbnailURL && !imageError}
              <img
                src={content.thumbnailURL}
                alt={content.title}
                class="w-full h-full object-cover"
                onerror={() => imageError = true}
              />
            {:else if imageError}
              <div class="w-full h-full flex items-center justify-center text-slate-400">
                <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            {:else}
              <div class="w-full h-full flex items-center justify-center text-slate-400">
                <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            {/if}

            <!-- Duration overlay -->
            {#if content.duration}
              <div class="absolute bottom-2 right-2 z-10 px-2 py-0.5 bg-black/80 text-white text-xs rounded">
                {content.duration}
              </div>
            {/if}

            <!-- Play icon overlay -->
            <div class="absolute inset-0 flex items-center justify-center bg-black/20">
              <div class="w-12 h-12 rounded-full bg-white/90 flex items-center justify-center shadow-lg">
                <svg class="w-6 h-6 text-teal-600 ml-1" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z" />
                </svg>
              </div>
            </div>
          </div>

          <!-- Content -->
          <div class="p-5">
            <span class="inline-block px-2 py-1 bg-purple-50 text-purple-700 text-xs font-medium rounded-lg mb-3">
              {$t('content.picker.videos') || 'Video'}
            </span>

            <h2 id="preview-title" class="text-lg font-semibold text-slate-900 mb-2 pr-8">
              {content.title}
            </h2>

            <!-- Attribution Section -->
            <div class="flex flex-wrap gap-3 text-xs text-slate-500 mb-3">
              {#if content.channelName}
                <span class="flex items-center gap-1" aria-label={$t('content.attribution.channel')}>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                  {content.channelName}
                </span>
              {/if}
              {#if formattedDate}
                <span class="flex items-center gap-1" aria-label={$t('content.attribution.uploadedOn')}>
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  {formattedDate}
                </span>
              {/if}
            </div>

            <!-- Health Disclaimer -->
            <ContentDisclaimer />

            <!-- Attach Button -->
            <button
              onclick={onAttach}
              class="w-full py-2.5 px-4 rounded-xl font-medium transition-colors flex items-center justify-center gap-2 mt-4
                {isSelected
                  ? 'bg-teal-100 text-teal-700 hover:bg-teal-200'
                  : 'bg-teal-500 text-white hover:bg-teal-600'}"
              aria-pressed={isSelected}
            >
              {#if isSelected}
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <span>{$t('content.preview.selected') || 'Dipilih'}</span>
              {:else}
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
                <span>{$t('content.preview.attach') || 'Lampirkan'}</span>
              {/if}
            </button>
          </div>
        </div>
      {/if}
    </div>
  </div>
{/if}
