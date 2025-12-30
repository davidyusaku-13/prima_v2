<script>
  /**
   * EmptyState - A reusable empty state component for displaying friendly messages
   * when there's no content to show.
   *
   * Props:
   * - icon: SVG string or component for visual indicator
   * - title: Main message (required)
   * - description: Detailed guidance (optional)
   * - actionLabel: CTA button text (optional)
   * - onAction: Callback function for CTA click (optional)
   * - showAsCard: Boolean for card vs inline styling (optional)
   */
  let {
    icon = null,
    title = '',
    description = '',
    actionLabel = '',
    onAction = null,
    showAsCard = true
  } = $props();

  function handleActionClick() {
    if (onAction && typeof onAction === 'function') {
      onAction();
    }
  }
</script>

{#if showAsCard}
  <div
    class="text-center py-8 px-4"
    role="status"
    aria-live="polite"
  >
    {#if icon}
      <div class="mb-3 flex justify-center">
        {@html icon}
      </div>
    {/if}
    <p class="text-lg font-medium text-slate-700 mb-1">{title}</p>
    {#if description}
      <p class="text-sm text-slate-500 mb-4">{description}</p>
    {/if}
    {#if actionLabel && onAction}
      <button
        onclick={handleActionClick}
        class="px-4 py-2 bg-teal-500 text-white font-medium rounded-lg hover:bg-teal-600 transition-colors focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-offset-2"
        aria-label={actionLabel}
      >
        {actionLabel}
      </button>
    {/if}
  </div>
{:else}
  <div
    class="flex items-center gap-3 py-4"
    role="status"
    aria-live="polite"
  >
    {#if icon}
      <div class="flex-shrink-0">
        {@html icon}
      </div>
    {/if}
    <div class="flex-1">
      <p class="font-medium text-slate-700">{title}</p>
      {#if description}
        <p class="text-sm text-slate-500">{description}</p>
      {/if}
    </div>
    {#if actionLabel && onAction}
      <button
        onclick={handleActionClick}
        class="px-3 py-1.5 text-sm bg-teal-500 text-white font-medium rounded-lg hover:bg-teal-600 transition-colors focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-offset-2"
        aria-label={actionLabel}
      >
        {actionLabel}
      </button>
    {/if}
  </div>
{/if}
