<script>
  import { t } from 'svelte-i18n';

  export let user = null;
  export let currentView = 'dashboard';
  export let stats = { totalPatients: 0 };
  export let users = [];
  export let onNavigate = () => {};
  export let onShowProfile = () => {};
</script>

<nav class="lg:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-slate-200 px-2 py-2 pb-safe z-30">
  <div class="flex items-center justify-around">
    <button
      onclick={() => onNavigate('dashboard')}
      class="flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors {currentView === 'dashboard' ? 'text-teal-600 bg-teal-50' : 'text-slate-500 hover:text-slate-700'}"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
      </svg>
      <span class="text-xs font-medium">{$t('navigation.dashboard')}</span>
    </button>

    <button
      onclick={() => onNavigate('patients')}
      class="relative flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors {currentView === 'patients' ? 'text-teal-600 bg-teal-50' : 'text-slate-500 hover:text-slate-700'}"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
      <span class="text-xs font-medium">{$t('navigation.patients')}</span>
      {#if stats.totalPatients > 0}
        <span class="absolute top-0 right-2 w-4 h-4 bg-teal-600 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
          {stats.totalPatients}
        </span>
      {/if}
    </button>

    {#if user?.role === 'superadmin'}
      <button
        onclick={() => onNavigate('users')}
        class="relative flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors {currentView === 'users' ? 'text-purple-600 bg-purple-50' : 'text-slate-500 hover:text-slate-700'}"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
        <span class="text-xs font-medium">{$t('navigation.users')}</span>
        {#if users.length > 0}
          <span class="absolute top-0 right-2 w-4 h-4 bg-purple-600 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
            {users.length}
          </span>
        {/if}
      </button>
    {/if}

    <!-- Public Content -->
    <button
      onclick={() => onNavigate('berita')}
      class="relative flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors {currentView === 'berita' ? 'text-blue-600 bg-blue-50' : 'text-slate-500 hover:text-slate-700'}"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
      </svg>
      <span class="text-xs font-medium">{$t('berita.title')}</span>
    </button>

    <button
      onclick={() => onNavigate('video')}
      class="relative flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors {currentView === 'video' ? 'text-red-600 bg-red-50' : 'text-slate-500 hover:text-slate-700'}"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
      </svg>
      <span class="text-xs font-medium">{$t('video.title')}</span>
    </button>

    <button
      onclick={onShowProfile}
      class="flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors text-slate-500 hover:text-slate-700"
    >
      <div class="w-8 h-8 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-sm">
        {user?.fullName?.charAt(0)?.toUpperCase() || user?.username?.charAt(0)?.toUpperCase() || 'U'}
      </div>
    </button>
  </div>
</nav>
