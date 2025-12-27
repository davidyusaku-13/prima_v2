<script>
  import { t } from 'svelte-i18n';

  export let user = null;
  export let currentView = 'dashboard';
  export let stats = { totalPatients: 0 };
  export let users = [];
  export let locale = 'en';
  export let onNavigate = () => {};
  export let onSetLocale = () => {};
  export let onLogout = () => {};
</script>

<aside class="hidden lg:flex flex-col w-64 bg-white border-r border-slate-200 fixed inset-y-0 left-0 z-30">
  <!-- Logo -->
  <div class="flex items-center gap-3 px-6 py-5 border-b border-slate-100">
    <div class="w-10 h-10 bg-teal-600 rounded-xl flex items-center justify-center">
      <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
      </svg>
    </div>
    <div>
      <h1 class="font-bold text-slate-900">{$t('app.name')}</h1>
      <p class="text-xs text-slate-500">{$t('navigation.volunteerDashboard')}</p>
    </div>
  </div>

  <!-- Navigation -->
  <nav class="flex-1 py-4 overflow-y-auto">
    <div class="px-4 space-y-1">
      <button
        onclick={() => onNavigate('dashboard')}
        class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-left transition-all duration-200 {currentView === 'dashboard' ? 'bg-teal-50 text-teal-700 font-medium' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
        </svg>
        {$t('navigation.dashboard')}
      </button>
      <button
        onclick={() => onNavigate('patients')}
        class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-left transition-all duration-200 {currentView === 'patients' ? 'bg-teal-50 text-teal-700 font-medium' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        {$t('navigation.patients')}
        <span class="ml-auto bg-slate-100 text-slate-600 text-xs font-medium px-2 py-0.5 rounded-full">{stats.totalPatients}</span>
      </button>
      {#if user?.role === 'superadmin'}
        <button
          onclick={() => onNavigate('users')}
          class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-left transition-all duration-200 {currentView === 'users' ? 'bg-purple-50 text-purple-700 font-medium' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          {$t('navigation.users')}
          <span class="ml-auto bg-purple-100 text-purple-700 text-xs font-medium px-2 py-0.5 rounded-full">{users.length}</span>
        </button>
      {/if}
    </div>
  </nav>

  <!-- User Section -->
  <div class="p-4 border-t border-slate-100">
    <!-- Language Switcher -->
    <div class="mb-3 px-4 py-2">
      <div class="flex items-center gap-1 bg-slate-100 rounded-lg p-1">
        <button
          onclick={() => onSetLocale('en')}
          class="flex-1 py-1.5 text-xs font-medium rounded-md transition-colors {locale === 'en' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'}"
        >
          EN
        </button>
        <button
          onclick={() => onSetLocale('id')}
          class="flex-1 py-1.5 text-xs font-medium rounded-md transition-colors {locale === 'id' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'}"
        >
          ID
        </button>
      </div>
    </div>
    <div class="flex items-center gap-3 px-4 py-3 rounded-xl bg-slate-50">
      <div class="w-10 h-10 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold">
        {user?.fullName?.charAt(0)?.toUpperCase() || user?.username?.charAt(0)?.toUpperCase() || 'U'}
      </div>
      <div class="flex-1 min-w-0">
        <p class="text-sm font-medium text-slate-900 truncate">{user?.fullName || user?.username}</p>
        <div class="flex items-center gap-2">
          <p class="text-xs text-slate-500 truncate">@{user?.username}</p>
          {#if user?.role}
            <span class="px-1.5 py-0.5 text-xs font-medium rounded {user?.role === 'superadmin' ? 'bg-purple-100 text-purple-700' : user?.role === 'admin' ? 'bg-blue-100 text-blue-700' : 'bg-slate-200 text-slate-600'}">
              {user?.role}
            </span>
          {/if}
        </div>
      </div>
      <button onclick={onLogout} class="p-2 text-slate-400 hover:text-red-600 transition-colors" title={$t('auth.logout')}>
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
      </button>
    </div>
  </div>
</aside>
