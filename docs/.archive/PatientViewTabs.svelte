<script>
  import { t } from 'svelte-i18n';

  let { activeTab = 'patients', onTabChange = () => {} } = $props();

  const tabs = [
    { id: 'patients', label: 'patients.tabs.patients', icon: 'users' },
    { id: 'reminders', label: 'patients.tabs.reminders', icon: 'bell' }
  ];

  function handleTabClick(tabId) {
    activeTab = tabId;
    onTabChange(tabId);
  }
</script>

<div class="flex gap-1 border-b border-slate-200 mb-4 sm:mb-6">
  {#each tabs as tab}
    <button
      onclick={() => handleTabClick(tab.id)}
      class="relative px-4 sm:px-6 py-3 font-medium text-sm sm:text-base transition-colors duration-200 {
        activeTab === tab.id
          ? 'text-teal-600'
          : 'text-slate-600 hover:text-slate-900'
      }"
    >
      <div class="flex items-center gap-2">
        {#if tab.icon === 'users'}
          <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
        {:else if tab.icon === 'bell'}
          <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
        {/if}
        <span>{$t(tab.label)}</span>
      </div>
      {#if activeTab === tab.id}
        <div class="absolute bottom-0 left-0 right-0 h-1 bg-teal-600 rounded-t-full"></div>
      {/if}
    </button>
  {/each}
</div>
