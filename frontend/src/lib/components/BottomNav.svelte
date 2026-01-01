<script>
  import { t } from 'svelte-i18n';

  let {
    user = null,
    currentView = 'dashboard',
    stats = { totalPatients: 0 },
    users = [],
    onNavigate = () => {},
    onShowProfile = () => {}
  } = $props();

  let adminMenuOpen = $state(false);
</script>

<nav class="lg:hidden fixed bottom-0 left-0 right-0 bg-white/95 backdrop-blur-md border-t border-slate-200 px-0.5 py-2 pb-safe z-30">
  <div class="flex items-center justify-around gap-0.5">
    <!-- Dashboard -->
    <button
      onclick={() => onNavigate('dashboard')}
      class="relative flex flex-col items-center gap-0.5 px-1 py-1 rounded-lg transition-colors duration-200 flex-1 {currentView === 'dashboard' ? 'text-teal-700' : 'text-slate-500'}"
      class:border-b-2={currentView === 'dashboard'}
      class:border-teal-700={currentView === 'dashboard'}
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
      </svg>
      <span class="text-xs font-medium">{$t('navigation.dashboard')}</span>
    </button>

    <!-- Patients -->
    <button
      onclick={() => onNavigate('patients')}
      class="relative flex flex-col items-center gap-0.5 px-1 py-1 rounded-lg transition-colors duration-200 flex-1 {currentView === 'patients' ? 'text-teal-700' : 'text-slate-500'}"
      class:border-b-2={currentView === 'patients'}
      class:border-teal-700={currentView === 'patients'}
    >
      <div class="relative">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        {#if stats.totalPatients > 0}
          <span class="absolute -top-2 -right-2 min-w-[18px] h-[18px] px-1 bg-teal-600 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
            {stats.totalPatients > 99 ? '99+' : stats.totalPatients}
          </span>
        {/if}
      </div>
      <span class="text-xs font-medium">{$t('navigation.patients')}</span>
    </button>

    <!-- Berita -->
    <button
      onclick={() => onNavigate('berita')}
      class="relative flex flex-col items-center gap-0.5 px-1 py-1 rounded-lg transition-colors duration-200 flex-1 {currentView === 'berita' || currentView === 'berita-detail' ? 'text-blue-700' : 'text-slate-500'}"
      class:border-b-2={currentView === 'berita' || currentView === 'berita-detail'}
      class:border-blue-700={currentView === 'berita' || currentView === 'berita-detail'}
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z" />
      </svg>
      <span class="text-xs font-medium">{$t('berita.title')}</span>
    </button>

    <!-- Video -->
    <button
      onclick={() => onNavigate('video')}
      class="relative flex flex-col items-center gap-0.5 px-1 py-1 rounded-lg transition-colors duration-200 flex-1 {currentView === 'video' ? 'text-red-700' : 'text-slate-500'}"
      class:border-b-2={currentView === 'video'}
      class:border-red-700={currentView === 'video'}
    >
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
      </svg>
      <span class="text-xs font-medium">{$t('video.title')}</span>
    </button>

    <!-- Admin Menu (Superadmin/Admin only) -->
    {#if user?.role === 'superadmin' || user?.role === 'admin'}
      <div class="relative flex flex-1">
        <button
          onclick={() => adminMenuOpen = !adminMenuOpen}
          class="relative flex flex-col items-center gap-0.5 px-1 py-1 rounded-lg transition-colors duration-200 w-full {adminMenuOpen ? 'text-amber-700' : 'text-slate-500'}"
          class:border-b-2={adminMenuOpen}
          class:border-amber-700={adminMenuOpen}
          title="Admin Menu"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2m0 7a1 1 0 110-2 1 1 0 010 2m0 7a1 1 0 110-2 1 1 0 010 2" />
          </svg>
          <span class="text-xs font-medium">Admin</span>
        </button>

        <!-- Admin dropdown menu -->
        {#if adminMenuOpen}
          <div class="absolute bottom-full right-0 mb-2 bg-white rounded-lg shadow-lg border border-slate-200 z-50 min-w-[140px]">
            <button
              onclick={() => {
                onNavigate('users');
                adminMenuOpen = false;
              }}
              class="w-full px-4 py-2 text-left text-sm hover:bg-slate-50 transition-colors {currentView === 'users' ? 'text-purple-700 font-medium' : 'text-slate-700'}"
            >
              👥 {$t('navigation.users')}
              {#if users.length > 0}
                <span class="ml-2 inline-block min-w-[18px] h-[18px] px-1 bg-purple-600 text-white text-[10px] font-bold rounded-full">
                  {users.length > 99 ? '99+' : users.length}
                </span>
              {/if}
            </button>
            <button
              onclick={() => {
                onNavigate('analytics');
                adminMenuOpen = false;
              }}
              class="w-full px-4 py-2 text-left text-sm hover:bg-slate-50 transition-colors border-t border-slate-200 {currentView === 'analytics' ? 'text-amber-700 font-medium' : 'text-slate-700'}"
            >
              📊 {$t('navigation.analytics')}
            </button>
          </div>
        {/if}
      </div>
    {/if}

    <!-- Profile -->
    <button
      onclick={onShowProfile}
      class="relative flex flex-col items-center justify-center px-2 py-2 rounded-lg transition-colors duration-200 flex-1 hover:bg-slate-50"
      title={$t('profile.profile')}
    >
      <div class="w-8 h-8 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-sm hover:ring-2 hover:ring-teal-400 transition-all">
        {user?.fullName?.charAt(0)?.toUpperCase() || user?.username?.charAt(0)?.toUpperCase() || 'U'}
      </div>
    </button>
  </div>
</nav>
