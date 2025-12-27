<script>
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';

  export let patients = [];
  export let onToggleReminder = () => {};

  function formatDate(dateStr) {
    if (!dateStr) return '';
    const lang = $locale || 'en';
    return new Date(dateStr).toLocaleDateString(lang, {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function formatRecurrence(recurrence) {
    if (!recurrence || recurrence.frequency === 'none') return '';
    const freqLabels = {
      daily: 'Daily',
      weekly: 'Weekly',
      monthly: 'Monthly',
      yearly: 'Yearly'
    };
    let label = freqLabels[recurrence.frequency] || recurrence.frequency;
    if (recurrence.interval > 1) {
      label = `Every ${recurrence.interval} ${recurrence.frequency}s`;
    }
    if (recurrence.frequency === 'weekly' && recurrence.daysOfWeek?.length > 0) {
      const dayLabels = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
      const days = recurrence.daysOfWeek.map(d => dayLabels[d]).join(', ');
      label += ` on ${days}`;
    }
    return label;
  }

  function getPriorityColor(priority) {
    const colors = {
      high: 'bg-red-100 text-red-700',
      medium: 'bg-amber-100 text-amber-700',
      low: 'bg-green-100 text-green-700'
    };
    return colors[priority] || colors.medium;
  }

  $: stats = {
    totalPatients: patients.length,
    totalReminders: patients.reduce((acc, p) => acc + (p.reminders?.length || 0), 0),
    completedReminders: patients.reduce((acc, p) => acc + (p.reminders?.filter(r => r.completed).length || 0), 0),
    pendingReminders: patients.reduce((acc, p) => acc + (p.reminders?.filter(r => !r.completed).length || 0), 0)
  };

  $: upcomingReminders = patients.flatMap(p =>
    (p.reminders || []).filter(r => !r.completed && r.dueDate)
      .map(r => ({ ...r, patientName: p.name, patientId: p.id }))
  ).sort((a, b) => new Date(a.dueDate) - new Date(b.dueDate)).slice(0, 5);
</script>

<!-- Header -->
<header class="sticky top-0 z-10 bg-white/80 backdrop-blur-md border-b border-slate-200 mb-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between py-4">
    <div class="flex items-center gap-3">
      <h1 class="text-xl font-bold text-slate-900">{$t('dashboard.title')}</h1>
      <span class="text-slate-500 text-sm hidden sm:inline">
        {new Date().toLocaleDateString($locale, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
      </span>
    </div>
    <span class="text-slate-500 text-sm sm:hidden">
      {new Date().toLocaleDateString($locale, { weekday: 'short', month: 'short', day: 'numeric' })}
    </span>
  </div>
</header>

<!-- Stats -->
<div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-6 mb-8">
  <div class="bg-white rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-slate-200 hover:shadow-lg transition-shadow duration-200">
    <div class="flex items-center justify-between mb-3 sm:mb-4">
      <div class="w-10 h-10 sm:w-12 sm:h-12 bg-teal-100 rounded-xl flex items-center justify-center">
        <svg class="w-5 h-5 sm:w-6 sm:h-6 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </div>
    </div>
    <div class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">{stats.totalPatients}</div>
    <div class="text-slate-500 font-medium text-xs sm:text-base">{$t('dashboard.totalPatients')}</div>
  </div>

  <div class="bg-white rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-slate-200 hover:shadow-lg transition-shadow duration-200">
    <div class="flex items-center justify-between mb-3 sm:mb-4">
      <div class="w-10 h-10 sm:w-12 sm:h-12 bg-blue-100 rounded-xl flex items-center justify-center">
        <svg class="w-5 h-5 sm:w-6 sm:h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
      </div>
    </div>
    <div class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">{stats.totalReminders}</div>
    <div class="text-slate-500 font-medium text-xs sm:text-base">{$t('dashboard.totalReminders')}</div>
  </div>

  <div class="bg-white rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-slate-200 hover:shadow-lg transition-shadow duration-200">
    <div class="flex items-center justify-between mb-3 sm:mb-4">
      <div class="w-10 h-10 sm:w-12 sm:h-12 bg-green-100 rounded-xl flex items-center justify-center">
        <svg class="w-5 h-5 sm:w-6 sm:h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
    </div>
    <div class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">{stats.completedReminders}</div>
    <div class="text-slate-500 font-medium text-xs sm:text-base">{$t('dashboard.completed')}</div>
  </div>

  <div class="bg-white rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-slate-200 hover:shadow-lg transition-shadow duration-200">
    <div class="flex items-center justify-between mb-3 sm:mb-4">
      <div class="w-10 h-10 sm:w-12 sm:h-12 bg-amber-100 rounded-xl flex items-center justify-center">
        <svg class="w-5 h-5 sm:w-6 sm:h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
    </div>
    <div class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">{stats.pendingReminders}</div>
    <div class="text-slate-500 font-medium text-xs sm:text-base">{$t('dashboard.pending')}</div>
  </div>
</div>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-4 lg:gap-6">
  <!-- Upcoming Reminders -->
  <div class="bg-white rounded-xl lg:rounded-2xl border border-slate-200 overflow-hidden">
    <div class="px-4 lg:px-6 py-3 lg:py-4 border-b border-slate-100">
      <h2 class="text-base lg:text-lg font-semibold text-slate-900">{$t('dashboard.upcomingReminders')}</h2>
    </div>
    <div class="p-4 lg:p-6 space-y-3 lg:space-y-4">
      {#if upcomingReminders.length === 0}
        <p class="text-slate-500 text-center py-6 lg:py-8">{$t('dashboard.noUpcomingReminders')}</p>
      {:else}
        {#each upcomingReminders as reminder}
          <div class="flex items-start gap-3 p-3 lg:p-4 rounded-lg lg:rounded-xl bg-slate-50 hover:bg-slate-100 transition-colors duration-200">
            <button
              onclick={() => onToggleReminder(reminder.patientId, reminder.id)}
              class="mt-0.5 w-5 h-5 rounded-full border-2 flex-shrink-0 transition-colors duration-200 {reminder.completed ? 'bg-teal-600 border-teal-600' : 'border-slate-300 hover:border-teal-600'}"
            >
              {#if reminder.completed}
                <svg class="w-full h-full text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                </svg>
              {/if}
            </button>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1 flex-wrap">
                <h3 class="font-medium text-slate-900 truncate">{reminder.title}</h3>
                <span class="px-2 py-0.5 text-xs font-medium rounded-full {getPriorityColor(reminder.priority)}">{reminder.priority}</span>
                {#if reminder.recurrence && reminder.recurrence.frequency !== 'none'}
                  <span class="px-2 py-0.5 text-xs font-medium rounded-full bg-purple-100 text-purple-700">
                    {formatRecurrence(reminder.recurrence)}
                  </span>
                {/if}
              </div>
              <p class="text-sm text-slate-600 truncate">{reminder.patientName}</p>
              <p class="text-xs text-slate-400 mt-1">{formatDate(reminder.dueDate)}</p>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>

  <!-- Recent Patients -->
  <div class="bg-white rounded-xl lg:rounded-2xl border border-slate-200 overflow-hidden">
    <div class="px-4 lg:px-6 py-3 lg:py-4 border-b border-slate-100 flex items-center justify-between">
      <h2 class="text-base lg:text-lg font-semibold text-slate-900">{$t('dashboard.recentPatients')}</h2>
      <slot name="viewAllPatients">
        <button class="text-xs lg:text-sm text-teal-600 hover:text-teal-700 font-medium">
          {$t('common.viewAll')}
        </button>
      </slot>
    </div>
    <div class="p-4 lg:p-6 space-y-3 lg:space-y-4">
      {#if patients.length === 0}
        <p class="text-slate-500 text-center py-6 lg:py-8">{$t('dashboard.noPatients')}</p>
      {:else}
        {#each patients.slice(0, 5) as patient}
          <div class="flex items-center gap-3 p-3 lg:p-4 rounded-lg lg:rounded-xl bg-slate-50 hover:bg-slate-100 transition-colors duration-200">
            <div class="w-10 h-10 lg:w-12 lg:h-12 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-lg lg:text-xl">
              {patient.name.charAt(0).toUpperCase()}
            </div>
            <div class="flex-1 min-w-0">
              <h3 class="font-medium text-slate-900 truncate">{patient.name}</h3>
              <p class="text-sm text-slate-500 truncate">{patient.phone || $t('patients.noPhone')}</p>
            </div>
            <div class="text-right">
              <div class="text-sm font-medium text-slate-900">
                {patient.reminders?.filter(r => r.completed).length || 0}/{patient.reminders?.length || 0}
              </div>
              <div class="text-xs text-slate-500">{$t('patients.reminders')}</div>
            </div>
          </div>
        {/each}
      {/if}
    </div>
  </div>
</div>
