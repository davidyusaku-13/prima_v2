<script>
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';

  export let patients = [];
  export let searchQuery = '';
  export let onOpenPatientModal = () => {};
  export let onOpenReminderModal = () => {};
  export let onDeletePatient = () => {};
  export let onToggleReminder = () => {};
  export let onDeleteReminder = () => {};

  const daysOfWeek = [
    { value: 0, label: 'Sun' },
    { value: 1, label: 'Mon' },
    { value: 2, label: 'Tue' },
    { value: 3, label: 'Wed' },
    { value: 4, label: 'Thu' },
    { value: 5, label: 'Fri' },
    { value: 6, label: 'Sat' }
  ];

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

  $: filteredPatients = patients.filter(p => {
    const matchesSearch = p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                          p.phone.includes(searchQuery) ||
                          p.email.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesSearch;
  });
</script>

<!-- Header -->
<header class="sticky top-0 z-10 bg-white/80 backdrop-blur-md border-b border-slate-200 mb-4 sm:mb-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 sm:gap-4 py-3 sm:py-4">
    <div class="flex flex-col sm:flex-row sm:items-center gap-3 sm:gap-4">
      <h1 class="text-lg sm:text-xl font-bold text-slate-900">{$t('patients.title')}</h1>
      <div class="relative w-full sm:w-56 md:w-64 lg:w-72 xl:w-80">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 sm:w-5 sm:h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          type="text"
          bind:value={searchQuery}
          placeholder={$t('common.searchPlaceholder')}
          class="pl-9 sm:pl-10 pr-4 py-2 sm:py-2.5 w-full bg-slate-100 border-0 rounded-lg sm:rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
        />
      </div>
    </div>
    <button
      onclick={onOpenPatientModal}
      class="flex items-center justify-center gap-2 px-4 sm:px-5 py-2 sm:py-2.5 bg-teal-600 text-white font-medium rounded-lg sm:rounded-xl hover:bg-teal-700 hover:shadow-lg active:scale-[0.98] transition-all duration-200 w-full sm:w-auto text-sm sm:text-base"
    >
      <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
      </svg>
      {$t('patients.addPatient')}
    </button>
  </div>
</header>

{#if filteredPatients.length === 0}
  <div class="bg-white rounded-xl sm:rounded-2xl border border-slate-200 p-8 sm:p-12 text-center">
    <div class="w-12 h-12 sm:w-16 sm:h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4">
      <svg class="w-6 h-6 sm:w-8 sm:h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
    </div>
    <h3 class="text-base sm:text-lg font-semibold text-slate-900 mb-2">{$t('patients.noPatients')}</h3>
    <p class="text-slate-500 mb-4 sm:mb-6 text-sm">
      {#if searchQuery}
        {$t('patients.noPatientsMatch')}
      {:else}
        {$t('patients.getStarted')}
      {/if}
    </p>
    <button
      onclick={onOpenPatientModal}
      class="px-5 sm:px-6 py-2.5 sm:py-3 bg-teal-600 text-white font-medium rounded-lg sm:rounded-xl hover:bg-teal-700 transition-colors duration-200 text-sm sm:text-base"
    >
      {$t('patients.addPatient')}
    </button>
  </div>
{:else}
  <div class="space-y-3 sm:space-y-4">
    {#each filteredPatients as patient}
      <div class="bg-white rounded-xl sm:rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg transition-all duration-200">
        <div class="p-3 sm:p-4 md:p-5 lg:p-6">
          <div class="flex items-start gap-2.5 sm:gap-3 md:gap-4">
            <div class="w-10 h-10 sm:w-12 sm:h-12 md:w-14 md:h-14 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-base sm:text-lg md:text-xl flex-shrink-0">
              {patient.name.charAt(0).toUpperCase()}
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-2 sm:gap-4">
                <div class="min-w-0 flex-1">
                  <h3 class="text-base sm:text-lg font-semibold text-slate-900 truncate">{patient.name}</h3>
                  {#if patient.phone}
                    <p class="text-slate-600 text-xs sm:text-sm">{patient.phone}</p>
                  {/if}
                  {#if patient.email}
                    <p class="text-slate-500 text-xs sm:text-sm truncate">{patient.email}</p>
                  {/if}
                  {#if patient.notes}
                    <p class="text-slate-500 text-xs sm:text-sm mt-1.5 sm:mt-2 line-clamp-2">{patient.notes}</p>
                  {/if}
                </div>
                <!-- Action buttons - horizontal on all screens -->
                <div class="flex items-center gap-0.5 sm:gap-1 md:gap-2 flex-shrink-0">
                  <button
                    onclick={() => onOpenReminderModal(patient.id)}
                    class="p-1.5 sm:p-2 text-teal-600 hover:bg-teal-50 active:bg-teal-100 rounded-lg transition-colors duration-200"
                    title={$t('patients.addReminder')}
                  >
                    <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                  </button>
                  <button
                    onclick={() => onOpenPatientModal(patient)}
                    class="p-1.5 sm:p-2 text-slate-600 hover:bg-slate-100 active:bg-slate-200 rounded-lg transition-colors duration-200"
                    title={$t('common.edit')}
                  >
                    <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button
                    onclick={() => onDeletePatient(patient.id)}
                    class="p-1.5 sm:p-2 text-red-600 hover:bg-red-50 active:bg-red-100 rounded-lg transition-colors duration-200"
                    title={$t('common.delete')}
                  >
                    <svg class="w-4 h-4 sm:w-5 sm:h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </div>

              <!-- Reminders section -->
              {#if patient.reminders && patient.reminders.length > 0}
                <div class="mt-3 sm:mt-4 pt-3 sm:pt-4 border-t border-slate-100">
                  <h4 class="text-xs sm:text-sm font-medium text-slate-700 mb-2 sm:mb-3">{$t('patients.reminders')}</h4>
                  <div class="space-y-1.5 sm:space-y-2">
                    {#each patient.reminders as reminder}
                      <div class="flex flex-wrap items-center gap-1.5 sm:gap-2 p-2 sm:p-3 rounded-lg {reminder.completed ? 'bg-slate-50' : 'bg-amber-50'}">
                        <button
                          onclick={() => onToggleReminder(patient.id, reminder.id)}
                          class="w-4 h-4 sm:w-5 sm:h-5 rounded-full border-2 flex-shrink-0 transition-colors duration-200 {reminder.completed ? 'bg-teal-600 border-teal-600' : 'border-slate-300 hover:border-teal-600'}"
                        >
                          {#if reminder.completed}
                            <svg class="w-full h-full text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                            </svg>
                          {/if}
                        </button>
                        <div class="flex-1 min-w-0">
                          <span class="text-xs sm:text-sm {reminder.completed ? 'text-slate-500 line-through' : 'text-slate-900'} truncate block">{reminder.title}</span>
                          {#if reminder.dueDate}
                            <span class="text-[10px] sm:text-xs text-slate-400">{formatDate(reminder.dueDate)}</span>
                          {/if}
                          {#if reminder.recurrence && reminder.recurrence.frequency !== 'none'}
                            <span class="text-[10px] sm:text-xs text-purple-600 ml-1" title={formatRecurrence(reminder.recurrence)}>
                              <svg class="w-2.5 h-2.5 sm:w-3 sm:h-3 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                              </svg>
                            </span>
                          {/if}
                        </div>
                        <span class="px-1.5 sm:px-2 py-0.5 text-[10px] sm:text-xs font-medium rounded-full {getPriorityColor(reminder.priority)}">{reminder.priority}</span>
                        <div class="flex items-center gap-0.5 sm:gap-1">
                          <button
                            onclick={() => onOpenReminderModal(patient.id, reminder)}
                            class="p-1 sm:p-1.5 text-slate-400 hover:text-teal-600 hover:bg-white/50 rounded transition-colors duration-200"
                            aria-label={$t('common.edit')}
                          >
                            <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                            </svg>
                          </button>
                          <button
                            onclick={() => onDeleteReminder(patient.id, reminder.id)}
                            class="p-1 sm:p-1.5 text-slate-400 hover:text-red-600 hover:bg-white/50 rounded transition-colors duration-200"
                            aria-label={$t('common.delete')}
                          >
                            <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    {/each}
                  </div>
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>
    {/each}
  </div>
{/if}
