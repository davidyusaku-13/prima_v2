<script>
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';
  import EmptyState from '$lib/components/ui/EmptyState.svelte';
  import DeliveryStatusBadge from '$lib/components/delivery/DeliveryStatusBadge.svelte';
  import CancelConfirmationModal from '$lib/components/reminders/CancelConfirmationModal.svelte';
  import * as api from '$lib/utils/api.js';

  /**
   * @typedef {Object} Props
   * @property {string} token - Auth token
   * @property {string} patientId - Patient ID
   * @property {Object} [patient] - Patient object (optional, for displaying patient info)
   * @property {Function} [onCreateReminder] - Callback to open reminder creation modal
   */

  /** @type {Props} */
  let {
    token,
    patientId,
    patient = null,
    onCreateReminder = () => {}
  } = $props();

  // State
  let reminders = $state([]);
  let loading = $state(true);
  let loadingMore = $state(false);
  let error = $state(null);
  let expandedId = $state(null);
  let page = $state(1);
  let hasMore = $state(true);
  let limit = 20;

  // Cancel modal state
  let showCancelModal = $state(false);
  let reminderToCancel = $state(null);
  let cancelling = $state(false);
  let cancelError = $state(null);

  // Format date helper
  function formatDate(dateStr) {
    if (!dateStr) return '-';
    const lang = $locale || 'id';
    try {
      return new Date(dateStr).toLocaleDateString(lang, {
        month: 'short',
        day: 'numeric',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return dateStr;
    }
  }

  // Check if reminder is cancelled
  function isCancelled(status) {
    return status === 'cancelled';
  }

  // Toggle expanded row
  function toggleExpand(id) {
    expandedId = expandedId === id ? null : id;
  }

  // Load reminders
  async function loadReminders(reset = false) {
    if (reset) {
      page = 1;
      reminders = [];
      hasMore = true;
    }

    if (!hasMore && !reset) return;

    const isFirstLoad = reminders.length === 0;
    if (isFirstLoad) {
      loading = true;
    } else {
      loadingMore = true;
    }

    error = null;

    try {
      const response = await api.fetchPatientReminders(token, patientId, {
        history: true,
        page,
        limit
      });

      const newReminders = response.data || [];
      reminders = reset ? newReminders : [...reminders, ...newReminders];
      hasMore = response.pagination?.has_more || false;
      page++;
    } catch (e) {
      error = e.message || $t('common.errorLoading');
      console.warn('Failed to load reminder history:', e.message);
    } finally {
      loading = false;
      loadingMore = false;
    }
  }

  // Load more reminders (pagination)
  async function loadMore() {
    if (loadingMore || !hasMore) return;
    await loadReminders();
  }

  let initialized = $state(false);

  // Initial load - run once when token and patientId are available
  $effect(() => {
    if (token && patientId && !initialized) {
      initialized = true;
      loadReminders(true);
    }
  });

  // Get attachment icon
  function getAttachmentIcon(type) {
    return type === 'article'
      ? '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"></path></svg>'
      : '<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>';
  }

  // Check if reminder can be cancelled
  function canCancelReminder(status) {
    return status === 'pending' || status === 'scheduled';
  }

  // Open cancel modal
  function openCancelModal(reminder) {
    reminderToCancel = reminder;
    showCancelModal = true;
    cancelError = null;
  }

  // Close cancel modal
  function closeCancelModal() {
    showCancelModal = false;
    reminderToCancel = null;
    cancelError = null;
  }

  // Handle cancel confirmation
  async function handleCancelConfirm() {
    if (!reminderToCancel) return;

    cancelling = true;
    cancelError = null;

    try {
      await api.cancelReminder(token, reminderToCancel.id);

      // Update the reminder in the list
      const index = reminders.findIndex(r => r.id === reminderToCancel.id);
      if (index !== -1) {
        reminders[index] = {
          ...reminders[index],
          delivery_status: 'cancelled',
          cancelled_at: new Date().toISOString()
        };
      }

      closeCancelModal();
      
      // Show success toast (AC requirement)
      showToast($t('reminder.cancel.success') || 'Reminder dibatalkan', 'success');
    } catch (e) {
      cancelError = e.message || $t('reminder.cancel.error') || 'Failed to cancel reminder';
      console.warn('Failed to cancel reminder:', e.message);
    } finally {
      cancelling = false;
    }
  }
  
  // Toast notification helper
  function showToast(message, type = 'info') {
    // Dispatch custom event for toast notification
    window.dispatchEvent(new CustomEvent('show-toast', {
      detail: { message, type }
    }));
  }
</script>

<div class="space-y-4 px-4 pt-2">
  <!-- Reminder List -->

  <!-- Loading state -->
  {#if loading}
    <div class="flex items-center justify-center py-12" role="status">
      <div class="animate-spin w-8 h-8 border-3 border-teal-600 border-t-transparent rounded-full"></div>
    </div>
  {:else if error}
    <!-- Error state -->
    <div class="text-center py-8">
      <p class="text-red-600 mb-4">{error}</p>
      <button
        onclick={() => loadReminders(true)}
        class="px-4 py-2 bg-teal-600 text-white rounded-lg hover:bg-teal-700"
      >
        {$t('common.refresh')}
      </button>
    </div>
  {:else if reminders.length === 0}
    <!-- Empty state -->
    <EmptyState
      icon={`<svg class="w-12 h-12 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
      </svg>`}
      title={$t('patients.noReminderHistory') || 'Belum ada riwayat reminder untuk pasien ini'}
      description={$t('patients.createFirstReminder') || 'Buat reminder pertama untuk memulai pelacakan'}
      actionLabel={$t('patients.createFirstReminderButton') || 'Buat Reminder Pertama'}
      onAction={onCreateReminder}
    />
  {:else}
    <!-- Reminder list -->
    <div class="space-y-3">
      {#each reminders as reminder (reminder.id)}
        {@const cancelled = isCancelled(reminder.delivery_status)}
        {@const isExpanded = expandedId === reminder.id}

        <div
          class="bg-white rounded-xl border transition-all duration-200 {
            cancelled ? 'border-slate-200 bg-slate-50' :
            isExpanded ? 'border-teal-300 shadow-md' :
            'border-slate-200 hover:shadow-sm'
          } {cancelled ? 'opacity-75' : ''}"
        >
          <!-- Summary row (clickable) -->
          <button
            onclick={() => toggleExpand(reminder.id)}
            class="w-full p-4 text-left flex items-start gap-3"
            aria-expanded={isExpanded}
          >
            <!-- Delivery Status Badge -->
            <div class="flex-shrink-0 mt-0.5">
              <DeliveryStatusBadge status={reminder.delivery_status} />
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0 flex-1">
                  <h4 class="font-medium text-slate-900 truncate {cancelled ? 'line-through text-slate-500' : ''}">
                    {reminder.title}
                  </h4>
                  {#if reminder.message_preview}
                    <p class="text-sm text-slate-500 truncate mt-0.5">
                      {reminder.message_preview}
                    </p>
                  {/if}
                </div>

                <!-- Expand/collapse icon -->
                <svg
                  class="w-5 h-5 text-slate-400 flex-shrink-0 transition-transform duration-200 {isExpanded ? 'rotate-180' : ''}"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                </svg>
              </div>

              <!-- Meta info row -->
              <div class="flex items-center gap-3 mt-2 text-xs text-slate-500">
                <span class="flex items-center gap-1">
                  <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  {formatDate(reminder.scheduled_at)}
                </span>

                {#if reminder.attachment_count > 0}
                  <span class="flex items-center gap-1">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"></path>
                    </svg>
                    {reminder.attachment_count} {$t('patients.attachments') || 'konten'}
                  </span>
                {/if}

                {#if cancelled && reminder.cancelled_at}
                  <span class="flex items-center gap-1 text-amber-600">
                    <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                    </svg>
                    {$t('patients.cancelled') || 'Dibatalkan'} {formatDate(reminder.cancelled_at)}
                  </span>
                {/if}
              </div>
            </div>
          </button>

          <!-- Expanded details -->
          {#if isExpanded}
            <div class="px-4 pb-4 pt-0 border-t border-slate-100">
              <div class="pt-4 space-y-4">
                <!-- Full message -->
                <div>
                  <h5 class="text-xs font-medium text-slate-500 uppercase tracking-wide mb-1">
                    {$t('patients.message') || 'Pesan'}
                  </h5>
                  <p class="text-sm text-slate-700 {cancelled ? 'line-through' : ''}">
                    {reminder.message || reminder.description || '-'}
                  </p>
                </div>

                <!-- Attachments -->
                {#if reminder.attachments && reminder.attachments.length > 0}
                  <div>
                    <h5 class="text-xs font-medium text-slate-500 uppercase tracking-wide mb-2">
                      {$t('patients.attachmentsTitle') || 'Lampiran'}
                    </h5>
                    <div class="space-y-2">
                      {#each reminder.attachments as attachment}
                        <a
                          href={attachment.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          class="flex items-center gap-2 p-2 bg-slate-50 rounded-lg hover:bg-slate-100 transition-colors"
                        >
                          <span class="text-slate-500">
                            {@html getAttachmentIcon(attachment.type)}
                          </span>
                          <span class="text-sm text-slate-700 truncate flex-1">
                            {attachment.title}
                          </span>
                          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"></path>
                          </svg>
                        </a>
                      {/each}
                    </div>
                  </div>
                {/if}

                <!-- Delivery timeline -->
                <div>
                  <h5 class="text-xs font-medium text-slate-500 uppercase tracking-wide mb-2">
                    {$t('patients.deliveryTimeline') || 'Timeline Pengiriman'}
                  </h5>
                  <div class="space-y-2">
                    {#if reminder.sent_at}
                      <div class="flex items-center gap-2 text-sm">
                        <span class="w-2 h-2 rounded-full bg-slate-400"></span>
                        <span class="text-slate-600">
                          {$t('patients.sent') || 'Dikirim'}: {formatDate(reminder.sent_at)}
                        </span>
                      </div>
                    {/if}
                    {#if reminder.delivered_at}
                      <div class="flex items-center gap-2 text-sm">
                        <span class="w-2 h-2 rounded-full bg-emerald-500"></span>
                        <span class="text-slate-600">
                          {$t('patients.delivered') || 'Diterima'}: {formatDate(reminder.delivered_at)}
                        </span>
                      </div>
                    {/if}
                    {#if reminder.read_at}
                      <div class="flex items-center gap-2 text-sm">
                        <span class="w-2 h-2 rounded-full bg-blue-500"></span>
                        <span class="text-slate-600">
                          {$t('patients.read') || 'Dibaca'}: {formatDate(reminder.read_at)}
                        </span>
                      </div>
                    {/if}
                    {#if !reminder.sent_at && !reminder.delivered_at && !reminder.read_at}
                      <p class="text-sm text-slate-400 italic">
                        {$t('patients.noDeliveryInfo') || 'Belum ada informasi pengiriman'}
                      </p>
                    {/if}
                  </div>
                </div>

                <!-- Error message (if failed) -->
                {#if reminder.delivery_status === 'failed' && reminder.delivery_error}
                  <div class="p-3 bg-red-50 rounded-lg border border-red-200">
                    <div class="flex items-start gap-2">
                      <svg class="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
                      </svg>
                      <div>
                        <p class="text-sm font-medium text-red-800">
                          {$t('patients.deliveryError') || 'Kesalahan Pengiriman'}
                        </p>
                        <p class="text-sm text-red-600 mt-1">
                          {reminder.delivery_error}
                        </p>
                      </div>
                    </div>
                  </div>
                {/if}

                <!-- Cancel button (only for cancellable reminders) -->
                {#if canCancelReminder(reminder.delivery_status)}
                  <div class="pt-2">
                    <button
                      onclick={() => openCancelModal(reminder)}
                      class="flex items-center gap-2 px-4 py-2 text-amber-600 bg-amber-50 border border-amber-200 rounded-lg hover:bg-amber-100 transition-colors duration-200"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                      </svg>
                      <span class="font-medium">{$t('reminder.cancel.button') || 'Batalkan Reminder'}</span>
                    </button>
                  </div>
                {/if}
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>

    <!-- Load More button -->
    {#if hasMore}
      <div class="text-center pt-4">
        <button
          onclick={loadMore}
          disabled={loadingMore}
          class="px-6 py-2.5 bg-white border border-slate-300 text-slate-700 font-medium rounded-xl hover:bg-slate-50 hover:border-slate-400 active:scale-[0.98] transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {#if loadingMore}
            <span class="flex items-center gap-2">
              <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {$t('common.loading')}
            </span>
          {:else}
            {$t('patients.loadMore') || 'Muat Lebih Banyak'}
          {/if}
        </button>
      </div>
    {:else if reminders.length > 0}
      <p class="text-center text-sm text-slate-400 py-4">
        {$t('patients.allHistoryLoaded') || 'Semua riwayat telah dimuat'}
      </p>
    {/if}
  {/if}
</div>

<!-- Cancel Confirmation Modal -->
<CancelConfirmationModal
  show={showCancelModal}
  reminder={reminderToCancel}
  onClose={closeCancelModal}
  onConfirm={handleCancelConfirm}
/>
