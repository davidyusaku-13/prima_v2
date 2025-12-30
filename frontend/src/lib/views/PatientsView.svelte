<script>
  import { onMount } from 'svelte';
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';
  import DeliveryStatusBadge from '$lib/components/delivery/DeliveryStatusBadge.svelte';
  import DeliveryStatusFilter from '$lib/components/delivery/DeliveryStatusFilter.svelte';
  import { deliveryStore } from '$lib/stores/delivery.svelte.js';

  let {
    patients = [],
    searchQuery = '',
    onOpenPatientModal = () => {},
    onOpenReminderModal = () => {},
    onDeletePatient = () => {},
    onToggleReminder = () => {},
    onDeleteReminder = () => {},
    onSendReminder = () => {},
    onRetryReminder = () => {}
  } = $props();

  // Reactive delivery status from store (Svelte 5 runes)
  let deliveryStatuses = $derived(deliveryStore.deliveryStatuses);
  let connectionStatus = $derived(deliveryStore.connectionStatus);
  let failedReminders = $derived(deliveryStore.failedReminders);

  // Filter state with session persistence
  let selectedFilter = $state(
    typeof window !== 'undefined'
      ? sessionStorage.getItem('reminderFilter') || 'all'
      : 'all'
  );

  // Save filter to session when changed
  $effect(() => {
    if (typeof window !== 'undefined') {
      sessionStorage.setItem('reminderFilter', selectedFilter);
    }
  });

  // Initialize SSE connection on mount
  onMount(() => {
    deliveryStore.connect();

    // Listen for navigate-to-patient event from toast action
    const handleNavigateToPatient = (event) => {
      const { patientId } = event.detail;
      // Scroll to patient card
      const patientCard = document.getElementById(`patient-${patientId}`);
      if (patientCard) {
        patientCard.scrollIntoView({ behavior: 'smooth', block: 'center' });
        // Highlight the card briefly
        patientCard.classList.add('ring-4', 'ring-red-500');
        setTimeout(() => {
          patientCard.classList.remove('ring-4', 'ring-red-500');
        }, 2000);
      }
    };

    // Listen for show-failed-reminders event from FailedReminderBadge
    const handleShowFailedReminders = () => {
      // Filter to show only failed reminders
      selectedFilter = 'failed';
    };

    window.addEventListener('navigate-to-patient', handleNavigateToPatient);
    window.addEventListener('show-failed-reminders', handleShowFailedReminders);

    // Cleanup on unmount
    return () => {
      deliveryStore.disconnect();
      window.removeEventListener('navigate-to-patient', handleNavigateToPatient);
      window.removeEventListener('show-failed-reminders', handleShowFailedReminders);
    };
  });

  // Get real-time status for a reminder
  function getReminderStatus(reminder) {
    const realtimeStatus = deliveryStatuses[reminder.id];
    return realtimeStatus?.status || reminder.delivery_status || 'pending';
  }

  // Get failure reason for a reminder
  function getFailureReason(reminderId) {
    const failedReminder = failedReminders.find(fr => fr.reminderId === reminderId);
    return failedReminder?.error || null;
  }

  // Check if reminder has failed
  function isReminderFailed(reminder) {
    const status = getReminderStatus(reminder);
    return status === 'failed';
  }

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

  // Svelte 5: Use $derived instead of $: reactive statement
  let filteredPatients = $derived(patients.filter(p => {
    const matchesSearch = p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                          p.phone.includes(searchQuery) ||
                          p.email.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesSearch;
  }));

  // Get all reminders from all patients
  let allReminders = $derived(() => {
    const reminders = [];
    filteredPatients.forEach(patient => {
      patient.reminders?.forEach(reminder => {
        reminders.push({
          ...reminder,
          patientId: patient.id,
          patientName: patient.name
        });
      });
    });
    return reminders;
  });

  // Filter reminders based on selected filter
  let filteredReminders = $derived(() => {
    const reminders = allReminders();

    if (selectedFilter === 'all') {
      return reminders;
    }

    return reminders.filter(reminder => {
      const deliveryStatus = deliveryStatuses[reminder.id];
      const status = deliveryStatus?.status || 'pending';

      switch (selectedFilter) {
        case 'pending':
          return status === 'pending' || status === 'queued';
        case 'sent':
          return status === 'sent' || status === 'delivered' || status === 'read';
        case 'failed':
          return status === 'failed';
        default:
          return true;
      }
    });
  });

  // Calculate counts for each filter
  let filterCounts = $derived(() => {
    const reminders = allReminders();
    const counts = { all: reminders.length, pending: 0, sent: 0, failed: 0 };

    reminders.forEach(reminder => {
      const deliveryStatus = deliveryStatuses[reminder.id];
      const status = deliveryStatus?.status || 'pending';

      if (status === 'pending' || status === 'queued') {
        counts.pending++;
      } else if (status === 'sent' || status === 'delivered' || status === 'read') {
        counts.sent++;
      } else if (status === 'failed') {
        counts.failed++;
      }
    });

    // Defensive: ensure all counts are non-negative
    return {
      all: Math.max(0, counts.all),
      pending: Math.max(0, counts.pending),
      sent: Math.max(0, counts.sent),
      failed: Math.max(0, counts.failed)
    };
  });

  // Handle filter change
  function handleFilterChange(filterId) {
    selectedFilter = filterId;
  }
</script>

<!-- Header -->
<header class="sticky top-0 z-10 bg-white/80 backdrop-blur-md border-b border-slate-200 mb-4 sm:mb-6">
  <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 sm:gap-4 py-3 sm:py-4">
    <div class="flex flex-col sm:flex-row sm:items-center gap-3 sm:gap-4">
      <div class="flex items-center gap-2">
        <h1 class="text-lg sm:text-xl font-bold text-slate-900">{$t('patients.title')}</h1>
        <!-- SSE Connection Status Indicator -->
        {#if connectionStatus === 'connected'}
          <span class="flex items-center gap-1 px-2 py-0.5 text-xs text-green-700 bg-green-50 rounded-full" title="Real-time updates active">
            <span class="w-1.5 h-1.5 bg-green-500 rounded-full animate-pulse"></span>
            Live
          </span>
        {:else if connectionStatus === 'connecting'}
          <span class="flex items-center gap-1 px-2 py-0.5 text-xs text-amber-700 bg-amber-50 rounded-full" title="Connecting...">
            <span class="w-1.5 h-1.5 bg-amber-500 rounded-full animate-pulse"></span>
            Connecting
          </span>
        {:else if connectionStatus === 'disconnected'}
          <span class="flex items-center gap-1 px-2 py-0.5 text-xs text-slate-500 bg-slate-100 rounded-full" title="Real-time updates unavailable">
            <span class="w-1.5 h-1.5 bg-slate-400 rounded-full"></span>
            Offline
          </span>
        {/if}
      </div>
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

<!-- Delivery Status Filter -->
{#if allReminders().length > 0}
  <div class="mb-4">
    <DeliveryStatusFilter
      selectedFilter={selectedFilter}
      counts={filterCounts()}
      onFilterChange={handleFilterChange}
    />
  </div>
{/if}

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
{:else if filteredReminders().length === 0 && allReminders().length > 0}
  <div class="bg-white rounded-xl sm:rounded-2xl border border-slate-200 p-8 sm:p-12 text-center">
    <div class="w-12 h-12 sm:w-16 sm:h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4">
      <svg class="w-6 h-6 sm:w-8 sm:h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
      </svg>
    </div>
    <h3 class="text-base sm:text-lg font-semibold text-slate-900 mb-2">{$t(`delivery.filter.empty.${selectedFilter}`)}</h3>
  </div>
{:else}
  <div class="space-y-3 sm:space-y-4">
    {#each filteredPatients as patient}
      <div id="patient-{patient.id}" class="bg-white rounded-xl sm:rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg transition-all duration-200">
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
                {@const visibleReminders = filteredReminders().filter(r => r.patientId === patient.id)}
                {#if visibleReminders.length > 0}
                  <div class="mt-3 sm:mt-4 pt-3 sm:pt-4 border-t border-slate-100">
                    <h4 class="text-xs sm:text-sm font-medium text-slate-700 mb-2 sm:mb-3">{$t('patients.reminders')}</h4>
                    <div class="space-y-1.5 sm:space-y-2">
                      {#each visibleReminders as reminder}
                      {@const isFailed = isReminderFailed(reminder)}
                      {@const failureReason = getFailureReason(reminder.id)}
                      <div class="flex flex-wrap items-center gap-1.5 sm:gap-2 p-2 sm:p-3 rounded-lg border-2 {
                        isFailed ? 'bg-red-50 border-red-500' :
                        reminder.completed ? 'bg-slate-50 border-transparent' :
                        'bg-amber-50 border-transparent'
                      }">
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
                          <span class="text-xs sm:text-sm {reminder.completed ? 'text-slate-500 line-through' : isFailed ? 'text-red-900 font-medium' : 'text-slate-900'} truncate block">{reminder.title}</span>
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
                          {#if isFailed && failureReason}
                            <div class="text-[10px] sm:text-xs text-red-600 mt-0.5 flex items-center gap-1">
                              <svg class="w-3 h-3 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
                              </svg>
                              <span class="truncate">{failureReason}</span>
                            </div>
                          {/if}
                        </div>
                        <span class="px-1.5 sm:px-2 py-0.5 text-[10px] sm:text-xs font-medium rounded-full {getPriorityColor(reminder.priority)}">{reminder.priority}</span>
                        <!-- Delivery Status Badge with real-time updates -->
                        {#if reminder.delivery_status || deliveryStatuses[reminder.id]}
                          <DeliveryStatusBadge
                            status={getReminderStatus(reminder)}
                            onRetry={() => onRetryReminder(patient, reminder)}
                            isRetrying={getReminderStatus(reminder) === 'sending'}
                          />
                        {/if}
                        <!-- Attachment count badge -->
                        {#if reminder.attachments && reminder.attachments.length > 0}
                          <span class="px-1.5 sm:px-2 py-0.5 text-[10px] sm:text-xs font-medium rounded-full bg-teal-100 text-teal-700" title={reminder.attachments.length + ' konten terlampir'}>
                            {reminder.attachments.length} konten
                          </span>
                        {/if}
                        <div class="flex items-center gap-0.5 sm:gap-1">
                          <!-- Send via WhatsApp button - prominent for failed reminders -->
                          <button
                            onclick={() => onSendReminder(patient, reminder)}
                            class="p-1 sm:p-1.5 rounded transition-colors duration-200 {
                              isFailed ? 'text-white bg-red-600 hover:bg-red-700 font-medium px-2 sm:px-3' :
                              'text-green-600 hover:text-green-700 hover:bg-green-50'
                            } {reminder.delivery_status === 'sending' ? 'opacity-50 cursor-not-allowed' : ''}"
                            aria-label={isFailed ? $t('reminder.retry') : $t('reminder.send.button')}
                            disabled={reminder.delivery_status === 'sending'}
                            title={isFailed ? $t('reminder.retry') : $t('reminder.send.title')}
                          >
                            {#if isFailed}
                              <span class="flex items-center gap-1 text-xs sm:text-sm">
                                <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                                </svg>
                                Coba Lagi
                              </span>
                            {:else}
                              <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/>
                              </svg>
                            {/if}
                          </button>
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
              {/if}
            </div>
          </div>
        </div>
      </div>
    {/each}
  </div>
{/if}
