<script>
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';

  export let users = [];
  export let currentUser = null;
  export let onLoadUsers = () => {};
  export let onOpenUserModal = () => {};
  export let onDeleteUser = () => {};
</script>

<!-- Header -->
<header class="sticky top-0 z-10 bg-white/80 backdrop-blur-md border-b border-slate-200 mb-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 py-4">
    <div class="flex items-center gap-3">
      <h1 class="text-xl font-bold text-slate-900">{$t('users.title')}</h1>
    </div>
    <div class="flex items-center gap-2">
      <button
        onclick={onLoadUsers}
        class="p-2 text-slate-400 hover:text-slate-600 rounded-lg hover:bg-slate-100 transition-colors"
        title={$t('common.refresh')}
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
      <button
        onclick={() => onOpenUserModal()}
        class="flex items-center gap-2 px-5 py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200 w-full sm:w-auto justify-center"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
        </svg>
        {$t('users.addUser')}
      </button>
    </div>
  </div>
</header>

<div class="bg-white rounded-2xl border border-slate-200 overflow-hidden">
  <div class="overflow-x-auto">
    {#if users.length === 0}
      <p class="text-slate-500 text-center py-12">{$t('users.noUsers')}</p>
    {:else}
      <div class="divide-y divide-slate-100">
        {#each users as u}
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 sm:px-6 sm:py-4 hover:bg-slate-50 transition-colors">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold flex-shrink-0">
                {u.fullName?.charAt(0)?.toUpperCase() || u.username?.charAt(0)?.toUpperCase() || 'U'}
              </div>
              <div class="min-w-0">
                <p class="text-sm font-medium text-slate-900 truncate">{u.fullName || $t('users.noName')}</p>
                <p class="text-sm text-slate-500">@{u.username}</p>
              </div>
            </div>
            <div class="flex items-center justify-between sm:gap-6 ml-0">
              <div class="flex items-center gap-4">
                <span class="px-2 py-1 text-xs font-medium rounded-full
                  {u.role === 'superadmin' ? 'bg-purple-100 text-purple-700' :
                   u.role === 'admin' ? 'bg-blue-100 text-blue-700' :
                   'bg-slate-100 text-slate-700'}">
                  {u.role}
                </span>
                <span class="hidden sm:block text-sm text-slate-500">
                  {u.createdAt ? new Date(u.createdAt).toLocaleDateString($locale) : '-'}
                </span>
              </div>
              <div class="flex items-center gap-1 sm:gap-2">
                {#if u.id !== currentUser?.userId}
                  <button
                    onclick={() => onOpenUserModal(u)}
                    class="p-2 text-slate-400 hover:text-teal-600 hover:bg-teal-50 rounded-lg transition-colors"
                    title={$t('common.edit')}
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button
                    onclick={() => onDeleteUser(u.id)}
                    class="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                    title={$t('common.delete')}
                  >
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                {:else}
                  <span class="text-xs text-slate-400 px-2">{$t('common.you')}</span>
                {/if}
              </div>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>
