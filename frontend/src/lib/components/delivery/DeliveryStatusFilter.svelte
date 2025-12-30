<script>
  import { _ as t } from 'svelte-i18n';

  // Props
  let {
    selectedFilter = 'all',
    counts = { all: 0, pending: 0, sent: 0, failed: 0 },
    onFilterChange = null
  } = $props();

  const filters = [
    { id: 'all', label: 'delivery.filter.all' },
    { id: 'pending', label: 'delivery.filter.pending' },
    { id: 'sent', label: 'delivery.filter.sent' },
    { id: 'failed', label: 'delivery.filter.failed' }
  ];

  function handleFilterClick(filterId) {
    if (onFilterChange) {
      onFilterChange(filterId);
    }
  }
</script>

<div class="flex gap-2 mb-4" role="tablist" aria-label={$t('delivery.filter.label')}>
  {#each filters as filter}
    <button
      onclick={() => handleFilterClick(filter.id)}
      class="px-4 py-2 rounded-full text-sm font-medium transition-colors"
      class:bg-teal-600={selectedFilter === filter.id}
      class:text-white={selectedFilter === filter.id}
      class:bg-gray-100={selectedFilter !== filter.id}
      class:text-gray-700={selectedFilter !== filter.id}
      class:hover:bg-teal-700={selectedFilter === filter.id}
      class:hover:bg-gray-200={selectedFilter !== filter.id}
      role="tab"
      aria-selected={selectedFilter === filter.id}
      aria-label={$t(filter.label)}
    >
      <span>{$t(filter.label)}</span>
      {#if counts[filter.id] > 0}
        <span class="ml-2 px-2 py-0.5 rounded-full text-xs"
          class:bg-teal-700={selectedFilter === filter.id}
          class:bg-gray-200={selectedFilter !== filter.id}
        >
          {counts[filter.id]}
        </span>
      {/if}
    </button>
  {/each}
</div>
