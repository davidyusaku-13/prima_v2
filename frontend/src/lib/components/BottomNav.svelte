<script>
  import { t } from 'svelte-i18n';

  export let user = null;
  export let currentView = 'dashboard';
  export let stats = { totalPatients: 0 };
  export let users = [];
  export let onNavigate = () => {};
  export let onShowProfile = () => {};
</script>

<nav class="lg:hidden fixed bottom-0 left-0 right-0 bg-white/95 backdrop-blur-md border-t border-slate-200 px-2 py-2 pb-safe z-30">
  <div class="flex items-center justify-around">
    <!-- Dashboard -->
    <button
      onclick={() => onNavigate('dashboard')}
      class="flex items-center gap-2 px-3 py-2 rounded-full transition-all duration-200 {currentView === 'dashboard' ? 'bg-teal-100 text-teal-700' : 'text-slate-500 active:bg-slate-100'}"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
      </svg>
      {#if currentView === 'dashboard'}
        <span class="text-sm font-medium">{$t('navigation.dashboard')}</span>
      {/if}
    </button>

    <!-- Patients -->
    <button
      onclick={() => onNavigate('patients')}
      class="relative flex items-center gap-2 px-3 py-2 rounded-full transition-all duration-200 {currentView === 'patients' ? 'bg-teal-100 text-teal-700' : 'text-slate-500 active:bg-slate-100'}"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
      {#if currentView === 'patients'}
        <span class="text-sm font-medium">{$t('navigation.patients')}</span>
      {/if}
      {#if stats.totalPatients > 0 && currentView !== 'patients'}
        <span class="absolute -top-1 -right-1 min-w-[18px] h-[18px] px-1 bg-teal-600 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
          {stats.totalPatients > 99 ? '99+' : stats.totalPatients}
        </span>
      {/if}
    </button>

    <!-- Users (Superadmin only) -->
    {#if user?.role === 'superadmin'}
      <button
        onclick={() => onNavigate('users')}
        class="relative flex items-center gap-2 px-3 py-2 rounded-full transition-all duration-200 {currentView === 'users' ? 'bg-purple-100 text-purple-700' : 'text-slate-500 active:bg-slate-100'}"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
        {#if currentView === 'users'}
          <span class="text-sm font-medium">{$t('navigation.users')}</span>
        {/if}
        {#if users.length > 0 && currentView !== 'users'}
          <span class="absolute -top-1 -right-1 min-w-[18px] h-[18px] px-1 bg-purple-600 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
            {users.length > 99 ? '99+' : users.length}
          </span>
        {/if}
      </button>
    {/if}

    <!-- Berita -->
    <button
      onclick={() => onNavigate('berita')}
      class="flex items-center gap-2 px-3 py-2 rounded-full transition-all duration-200 {currentView === 'berita' || currentView === 'berita-detail' ? 'bg-blue-100 text-blue-700' : 'text-slate-500 active:bg-slate-100'}"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
      </svg>
      {#if currentView === 'berita' || currentView === 'berita-detail'}
        <span class="text-sm font-medium">{$t('berita.title')}</span>
      {/if}
    </button>

    <!-- Video -->
    <button
      onclick={() => onNavigate('video')}
      class="flex items-center gap-2 px-3 py-2 rounded-full transition-all duration-200 {currentView === 'video' ? 'bg-red-100 text-red-700' : 'text-slate-500 active:bg-slate-100'}"
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
      </svg>
      {#if currentView === 'video'}
        <span class="text-sm font-medium">{$t('video.title')}</span>
      {/if}
    </button>

    <!-- Analytics (Admin/Superadmin only) -->
    {#if user?.role === 'superadmin' || user?.role === 'admin'}
      <button
        onclick={() => onNavigate('analytics')}
        class="flex items-center gap-2 px-3 py-2 rounded-full transition-all duration-200 {currentView === 'analytics' ? 'bg-amber-100 text-amber-700' : 'text-slate-500 active:bg-slate-100'}"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        {#if currentView === 'analytics'}
          <span class="text-sm font-medium">{$t('navigation.analytics')}</span>
        {/if}
      </button>
    {/if}

    <!-- Profile -->
    <button
      onclick={onShowProfile}
      class="flex items-center gap-2 px-3 py-2 rounded-full transition-all duration-200 text-slate-500 active:bg-slate-100"
    >
      <div class="w-6 h-6 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-xs">
        {user?.fullName?.charAt(0)?.toUpperCase() || user?.username?.charAt(0)?.toUpperCase() || 'U'}
      </div>
    </button>
  </div>
</nav>
