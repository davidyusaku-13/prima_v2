<script>
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';
  import DeliveryStatusBadge from '$lib/components/delivery/DeliveryStatusBadge.svelte';

  let {
    allReminders = [],
    deliveryStatuses = {},
    failedReminders = [],
    selectedPriority = 'all',
    selectedStatus = 'all',
    onOpenReminderModal = () => {},
    onToggleReminder = () => {},
    onDeleteReminder = () => {},
    onSendReminder = () => {},
    onRetryReminder = () => {}
  } = $props();

  // Accordion state
  let filtersOpen = $state(false);

  // Confirmation modal state
  let deleteConfirmation = $state({
    isOpen: false,
    reminderId: null,
    patientId: null,
    reminderTitle: ''
  });

  function openDeleteConfirmation(patientId, reminderId, reminderTitle) {
    deleteConfirmation = {
      isOpen: true,
      patientId,
      reminderId,
      reminderTitle
    };
  }

  function closeDeleteConfirmation() {
    deleteConfirmation = {
      isOpen: false,
      reminderId: null,
      patientId: null,
      reminderTitle: ''
    };
  }

  function confirmDelete() {
    if (deleteConfirmation.patientId && deleteConfirmation.reminderId) {
      onDeleteReminder(deleteConfirmation.patientId, deleteConfirmation.reminderId);
    }
    closeDeleteConfirmation();
  }

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

  function getPriorityColor(priority) {
    const colors = {
      high: 'bg-red-100 text-red-700 border-red-200',
      medium: 'bg-amber-100 text-amber-700 border-amber-200',
      low: 'bg-green-100 text-green-700 border-green-200'
    };
    return colors[priority] || colors.medium;
  }

  function getPriorityOrder(priority) {
    const order = { high: 0, medium: 1, low: 2 };
    return order[priority] ?? 999;
  }

  function getReminderStatus(reminder) {
    const realtimeStatus = deliveryStatuses[reminder.id];
    return realtimeStatus?.status || reminder.delivery_status || 'pending';
  }

  function getFailureReason(reminderId) {
    const failedReminder = failedReminders.find(fr => fr.reminderId === reminderId);
    return failedReminder?.error || null;
  }

  function isReminderFailed(reminder) {
    const status = getReminderStatus(reminder);
    return status === 'failed';
  }

  // Filter and sort reminders
  let filteredAndSortedReminders = $derived(() => {
    let reminders = [...allReminders];

    // Filter by priority
    if (selectedPriority !== 'all') {
      reminders = reminders.filter(r => r.priority === selectedPriority);
    }

    // Filter by status
    if (selectedStatus !== 'all') {
      reminders = reminders.filter(r => {
        const status = getReminderStatus(r);
        if (selectedStatus === 'pending') {
          return status === 'pending' || status === 'queued';
        } else if (selectedStatus === 'sent') {
          return status === 'sent' || status === 'delivered' || status === 'read';
        } else if (selectedStatus === 'failed') {
          return status === 'failed';
        }
        return true;
      });
    }

    // Sort by priority, then by due date
    return reminders.sort((a, b) => {
      const priorityDiff = getPriorityOrder(a.priority) - getPriorityOrder(b.priority);
      if (priorityDiff !== 0) return priorityDiff;

      // If same priority, sort by due date (earlier first)
      if (a.dueDate && b.dueDate) {
        return new Date(a.dueDate) - new Date(b.dueDate);
      }
      return 0;
    });
  });

  // Count reminders by status
  let statusCounts = $derived(() => {
    const reminders = allReminders;
    let counts = { all: reminders.length, pending: 0, sent: 0, failed: 0 };

    reminders.forEach(reminder => {
      const status = getReminderStatus(reminder);
      if (status === 'pending' || status === 'queued') {
        counts.pending++;
      } else if (status === 'sent' || status === 'delivered' || status === 'read') {
        counts.sent++;
      } else if (status === 'failed') {
        counts.failed++;
      }
    });

    return counts;
  });

  // Count reminders by priority
  let priorityCounts = $derived(() => {
    const reminders = allReminders;
    const counts = { all: reminders.length, high: 0, medium: 0, low: 0 };

    reminders.forEach(reminder => {
      if (reminder.priority === 'high') counts.high++;
      else if (reminder.priority === 'medium') counts.medium++;
      else if (reminder.priority === 'low') counts.low++;
    });

    return counts;
  });
</script>

<div class="space-y-4">
  <!-- Filters Accordion -->
  <details bind:open={filtersOpen} class="bg-white rounded-lg sm:rounded-xl border border-slate-200 overflow-hidden">
    <summary class="cursor-pointer px-3 sm:px-4 py-3 sm:py-3.5 font-medium text-slate-900 hover:bg-slate-50 flex items-center justify-between select-none">
      <span class="flex items-center gap-2">
        <svg class="w-4 h-4 sm:w-5 sm:h-5 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
        </svg>
        {$t('patients.filters')} 
        {#if selectedStatus !== 'all' || selectedPriority !== 'all'}
          <span class="px-2 py-0.5 text-xs font-semibold bg-teal-100 text-teal-700 rounded-full">
            {(selectedStatus !== 'all' ? 1 : 0) + (selectedPriority !== 'all' ? 1 : 0)}
          </span>
        {/if}
      </span>
      <svg class="w-4 h-4 sm:w-5 sm:h-5 text-slate-600 transition-transform duration-200 {filtersOpen ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
      </svg>
    </summary>

    <div class="border-t border-slate-200 p-3 sm:p-4 space-y-3 sm:space-y-4 bg-slate-50">
      <!-- Status Filter -->
      <div>
        <label for="status-filter" class="text-xs sm:text-sm font-medium text-slate-700 mb-2 block">{$t('patients.filterByStatus')}</label>
        <div id="status-filter" class="flex flex-wrap gap-2">
          {#each [
            { id: 'all', label: `All (${statusCounts().all})`, color: 'slate' },
            { id: 'pending', label: `Pending (${statusCounts().pending})`, color: 'amber' },
            { id: 'sent', label: `Sent (${statusCounts().sent})`, color: 'green' },
            { id: 'failed', label: `Failed (${statusCounts().failed})`, color: 'red' }
          ] as filter}
            <button
              onclick={() => selectedStatus = filter.id}
              class="px-3 sm:px-4 py-1.5 sm:py-2 text-xs sm:text-sm font-medium rounded-full transition-colors duration-200 {
                selectedStatus === filter.id
                  ? `bg-${filter.color}-600 text-white`
                  : `bg-${filter.color}-50 text-${filter.color}-700 hover:bg-${filter.color}-100 border border-${filter.color}-200`
              }"
            >
              {filter.label}
            </button>
          {/each}
        </div>
      </div>

      <!-- Priority Filter -->
      <div>
        <label for="priority-filter" class="text-xs sm:text-sm font-medium text-slate-700 mb-2 block">{$t('patients.filterByPriority')}</label>
        <div id="priority-filter" class="flex flex-wrap gap-2">
          {#each [
            { id: 'all', label: `All (${priorityCounts().all})`, color: 'slate' },
            { id: 'high', label: `High (${priorityCounts().high})`, color: 'red' },
            { id: 'medium', label: `Medium (${priorityCounts().medium})`, color: 'amber' },
            { id: 'low', label: `Low (${priorityCounts().low})`, color: 'green' }
          ] as filter}
            <button
              onclick={() => selectedPriority = filter.id}
              class="px-3 sm:px-4 py-1.5 sm:py-2 text-xs sm:text-sm font-medium rounded-full transition-colors duration-200 {
                selectedPriority === filter.id
                  ? `bg-${filter.color}-600 text-white`
                  : `bg-${filter.color}-50 text-${filter.color}-700 hover:bg-${filter.color}-100 border border-${filter.color}-200`
              }"
            >
              {filter.label}
            </button>
          {/each}
        </div>
      </div>
    </div>
  </details>

  <!-- Reminders List -->
  {#if filteredAndSortedReminders().length === 0}
    <div class="bg-white rounded-xl sm:rounded-2xl border border-slate-200 p-8 sm:p-12 text-center">
      <div class="w-12 h-12 sm:w-16 sm:h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4">
        <svg class="w-6 h-6 sm:w-8 sm:h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
      </div>
      <h3 class="text-base sm:text-lg font-semibold text-slate-900 mb-2">{$t('patients.noReminders')}</h3>
      <p class="text-slate-500 text-sm">{$t('patients.noRemindersMatch')}</p>
    </div>
  {:else}
    <div class="space-y-2 sm:space-y-3">
      {#each filteredAndSortedReminders() as reminder}
        {@const isFailed = isReminderFailed(reminder)}
        {@const failureReason = getFailureReason(reminder.id)}
        <div class="bg-white rounded-lg sm:rounded-xl border-2 p-3 sm:p-4 transition-all duration-200 {
          isFailed ? 'border-red-500 bg-red-50' :
          reminder.completed ? 'border-slate-200 bg-slate-50' :
          'border-slate-200 hover:border-teal-300 hover:shadow-md'
        }">
          <div class="flex items-start gap-3 sm:gap-4">
            <!-- Checkbox -->
            <button
              onclick={() => onToggleReminder(reminder.patientId, reminder.id)}
              class="w-5 h-5 sm:w-6 sm:h-6 rounded-full border-2 flex-shrink-0 flex items-center justify-center transition-colors duration-200 {reminder.completed ? 'bg-teal-600 border-teal-600' : 'border-slate-300 hover:border-teal-600'}"
            >
              {#if reminder.completed}
                <svg class="w-full h-full text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                </svg>
              {/if}
            </button>

            <!-- Reminder Content -->
            <div class="flex-1 min-w-0">
              <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-2 sm:gap-4">
                <div class="min-w-0 flex-1">
                  <!-- Title + Patient Name -->
                  <div>
                    <h4 class="text-sm sm:text-base font-semibold {reminder.completed ? 'text-slate-500 line-through' : isFailed ? 'text-red-900' : 'text-slate-900'}">
                      {reminder.title}
                    </h4>
                    <p class="text-xs sm:text-sm text-slate-500 mt-0.5">
                      <a href="#patient-{reminder.patientId}" class="hover:text-teal-600 underline">{reminder.patientName}</a>
                    </p>
                  </div>

                  <!-- Meta info -->
                  <div class="flex flex-wrap items-center gap-2 mt-2 sm:mt-3">
                    {#if reminder.dueDate}
                      <span class="text-xs sm:text-sm text-slate-600 bg-slate-100 px-2 sm:px-3 py-1 rounded-full">
                        {formatDate(reminder.dueDate)}
                      </span>
                    {/if}
                    <span class="text-xs sm:text-sm font-medium rounded-full px-2 sm:px-3 py-1 border {getPriorityColor(reminder.priority)}">
                      {reminder.priority}
                    </span>
                    {#if reminder.delivery_status || deliveryStatuses[reminder.id]}
                      <DeliveryStatusBadge
                        status={getReminderStatus(reminder)}
                        onRetry={() => onRetryReminder({ id: reminder.patientId, name: reminder.patientName }, reminder)}
                        isRetrying={getReminderStatus(reminder) === 'sending'}
                      />
                    {/if}
                  </div>

                  <!-- Error message if failed -->
                  {#if isFailed && failureReason}
                    <div class="text-xs sm:text-sm text-red-600 mt-2 flex items-center gap-2 bg-red-50 px-2 sm:px-3 py-1.5 rounded-md">
                      <svg class="w-4 h-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
                      </svg>
                      {failureReason}
                    </div>
                  {/if}
                </div>

                <!-- Action buttons -->
                <div class="flex items-center gap-1 sm:gap-2 flex-shrink-0">
                  <button
                    onclick={() => onSendReminder({ id: reminder.patientId, name: reminder.patientName }, reminder)}
                    class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium transition-colors duration-200 {
                      isFailed ? 'text-white bg-red-600 hover:bg-red-700' :
                      'text-green-600 bg-green-50 hover:bg-green-100 border border-green-200'
                    } {reminder.delivery_status === 'sending' ? 'opacity-50 cursor-not-allowed' : ''}"
                    disabled={reminder.delivery_status === 'sending'}
                  >
                    {#if isFailed}
                      <span class="flex items-center gap-1">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                        </svg>
                        Retry
                      </span>
                    {:else}
                      <span class="flex items-center gap-1">
                        <svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
                          <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/>
                        </svg>
                        Send
                      </span>
                    {/if}
                  </button>
                  <button
                    onclick={() => onOpenReminderModal(reminder.patientId, reminder)}
                    class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium text-slate-600 bg-slate-50 hover:bg-slate-100 border border-slate-200 transition-colors duration-200 flex items-center gap-1"
                    title={$t('common.edit')}
                    aria-label={$t('common.edit')}
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                    {$t('common.edit')}
                  </button>
                  <button
                    onclick={() => openDeleteConfirmation(reminder.patientId, reminder.id, reminder.title)}
                    class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium text-red-600 bg-red-50 hover:bg-red-100 border border-red-200 transition-colors duration-200 flex items-center gap-1"
                    title={$t('common.delete')}
                    aria-label={$t('common.delete')}
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    {$t('common.delete')}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
<!-- Delete Confirmation Modal -->
{#if deleteConfirmation.isOpen}
  <div class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
    <div class="bg-white rounded-xl sm:rounded-2xl max-w-sm w-full shadow-xl">
      <div class="p-4 sm:p-6">
        <!-- Icon -->
        <div class="w-12 h-12 sm:w-14 sm:h-14 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4">
          <svg class="w-6 h-6 sm:w-8 sm:h-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </div>

        <!-- Content -->
        <h3 class="text-lg sm:text-xl font-semibold text-slate-900 text-center mb-2">
          {$t('reminder.delete.title')}
        </h3>
        <p class="text-slate-600 text-sm sm:text-base text-center mb-4 sm:mb-6">
          {$t('reminder.delete.confirm', { values: { title: deleteConfirmation.reminderTitle } })}
        </p>

        <!-- Buttons -->
        <div class="flex gap-3 sm:gap-4">
          <button
            onclick={closeDeleteConfirmation}
            class="flex-1 px-4 py-2 sm:py-2.5 border border-slate-300 text-slate-700 font-medium rounded-lg sm:rounded-xl hover:bg-slate-50 transition-colors duration-200 text-sm sm:text-base"
          >
            {$t('common.cancel')}
          </button>
          <button
            onclick={confirmDelete}
            class="flex-1 px-4 py-2 sm:py-2.5 bg-red-600 text-white font-medium rounded-lg sm:rounded-xl hover:bg-red-700 transition-colors duration-200 text-sm sm:text-base"
          >
            {$t('common.delete')}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}