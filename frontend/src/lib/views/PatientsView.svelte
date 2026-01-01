<script>
  import { onMount } from 'svelte';
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';
  import DeliveryStatusBadge from '$lib/components/delivery/DeliveryStatusBadge.svelte';
  import PatientViewTabs from '$lib/components/patients/PatientViewTabs.svelte';
  import PatientListTab from '$lib/components/patients/PatientListTab.svelte';
  import ReminderListTab from '$lib/components/patients/ReminderListTab.svelte';
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
  let activeTab = $state(
    typeof window !== 'undefined'
      ? sessionStorage.getItem('patientsTab') || 'patients'
      : 'patients'
  );

  let selectedPriority = $state('all');
  let selectedStatus = $state('all');

  let selectedFilter = $state(
    typeof window !== 'undefined'
      ? sessionStorage.getItem('reminderFilter') || 'all'
      : 'all'
  );

  // Save active tab to session when changed
  $effect(() => {
    if (typeof window !== 'undefined') {
      sessionStorage.setItem('patientsTab', activeTab);
    }
  });

  function handleTabChange(tabId) {
    activeTab = tabId;
  }

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

<!-- Tab Navigation -->
<PatientViewTabs {activeTab} onTabChange={handleTabChange} />

<!-- Patient Tab View -->
{#if activeTab === 'patients'}
  <PatientListTab
    filteredPatients={filteredPatients}
    {searchQuery}
    {deliveryStatuses}
    {failedReminders}
    {onOpenPatientModal}
    {onOpenReminderModal}
    {onDeletePatient}
    {onToggleReminder}
    {onDeleteReminder}
    {onSendReminder}
    {onRetryReminder}
  />
{:else if activeTab === 'reminders'}
  <ReminderListTab
    allReminders={allReminders()}
    {deliveryStatuses}
    {failedReminders}
    {selectedPriority}
    {selectedStatus}
    {onOpenReminderModal}
    {onToggleReminder}
    {onDeleteReminder}
    {onSendReminder}
    {onRetryReminder}
  />
{/if}
