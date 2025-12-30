<script>
  import { t } from 'svelte-i18n';

  let {
    attachment = null,
    contentTitle = '',
    contentType = 'article',
    onRemove = () => {},
    showRemove = true,
    onClick = () => {} // Added onClick prop for preview
  } = $props();

  // Get display title from attachment or prop
  let displayTitle = $derived(attachment?.title || contentTitle || 'Konten tidak tersedia');

  // Get content type from attachment or prop
  let type = $derived(attachment?.type || contentType);

  // Get icon based on content type
  let typeIcon = $derived(type === 'video' ? '🎬' : '📄');

  // Handle remove action
  function handleRemove(e) {
    e.stopPropagation();
    onRemove(attachment);
  }

  // Handle click - trigger onClick for preview
  function handleClick(e) {
    onClick(attachment);
  }

  // Handle keyboard interaction
  function handleKeydown(e) {
    if (e.key === 'Enter' || e.key === ' ') {
      e.preventDefault();
      handleClick(e); // Trigger click action
    } else if (e.key === 'Backspace' || e.key === 'Delete') {
      if (showRemove) {
        handleRemove(e);
      }
    }
  }
</script>

<div
  class="inline-flex items-center gap-2 px-3 py-1.5 bg-teal-50 text-teal-700 rounded-full text-sm font-medium cursor-pointer hover:bg-teal-100 transition-colors duration-150"
  role="button"
  tabindex="0"
  aria-label="{displayTitle}, tekan hapus untuk melepas"
  onclick={handleClick}
  onkeydown={handleKeydown}
>
  <!-- Content type icon -->
  <span class="flex-shrink-0 w-4 h-4 flex items-center justify-center text-xs" aria-hidden="true">
    {typeIcon}
  </span>

  <!-- Truncated title -->
  <span class="max-w-[150px] truncate">
    {displayTitle}
  </span>

  <!-- Remove button -->
  {#if showRemove}
    <button
      onclick={handleRemove}
      class="flex-shrink-0 w-5 h-5 flex items-center justify-center rounded-full hover:bg-teal-200 transition-colors duration-150"
      aria-label="Hapus {displayTitle}"
    >
      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
      </svg>
    </button>
  {/if}
</div>
