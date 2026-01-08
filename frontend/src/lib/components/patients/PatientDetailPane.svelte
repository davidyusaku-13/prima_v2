<script>
  import { t } from "svelte-i18n";
  import { locale } from "svelte-i18n";
  import { deliveryStore } from "$lib/stores/delivery.svelte.js";
  import { toastStore } from "$lib/stores/toast.svelte.js";
  import DeliveryStatusBadge from "$lib/components/delivery/DeliveryStatusBadge.svelte";
  import ReminderHistoryView from "$lib/views/patients/ReminderHistoryView.svelte";

  /**
   * @typedef {Object} Props
   * @property {Object|null} patient - Selected patient object
   * @property {Record<string, { status: string }>} deliveryStatuses - Real-time delivery statuses
   * @property {Array<{ reminderId: string; error: string }>} failedReminders - Failed reminder objects
   * @property {string} activeTab - Current tab: 'info' | 'reminders' | 'history'
   * @property {Function} onTabChange - Callback when tab changes
   * @property {Function} [onEditPatient] - Callback for edit patient
   * @property {Function} [onAddReminder] - Callback for add reminder
   * @property {Function} [onEditReminder] - Callback for edit reminder
   * @property {Function} [onToggleReminder] - Callback to toggle reminder completion
   * @property {Function} [onDeleteReminder] - Callback to delete reminder
   * @property {Function} [onSendReminder] - Callback to send reminder
   * @property {Function} [onRetryReminder] - Callback to retry failed reminder
   * @property {string} [token] - Auth token (for ReminderHistoryView)
   */

  /** @type {Props} */
  let {
    patient = null,
    deliveryStatuses = {},
    failedReminders = [],
    activeTab = "info",
    onTabChange = () => {},
    onEditPatient = null,
    onAddReminder = null,
    onEditReminder = null,
    onToggleReminder = null,
    onDeleteReminder = null,
    onSendReminder = null,
    onRetryReminder = null,
    token = ""
  } = $props();

  // Tabs configuration
  const tabs = [
    { id: "info", label: "patients.info", icon: "info" },
    { id: "reminders", label: "patients.reminders", icon: "reminder" },
    { id: "history", label: "patients.reminderHistory", icon: "history" }
  ];

  // Track previous patient ID to reset tab when patient changes
  let previousPatientId = $state(null);

  // Reset tab to 'info' when switching to a different patient (AC26)
  $effect(() => {
    if (patient?.id !== previousPatientId) {
      if (previousPatientId !== null) {
        // Patient changed, reset to info tab
        onTabChange("info");
      }
      previousPatientId = patient?.id || null;
    }
  });

  /**
   * Format date for display
   * @param {string} dateStr
   * @returns {string}
   */
  function formatDate(dateStr) {
    if (!dateStr) return "-";
    const lang = $locale || "id";
    try {
      return new Date(dateStr).toLocaleDateString(lang, {
        day: "numeric",
        month: "long",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit"
      });
    } catch {
      return dateStr;
    }
  }

  /**
   * Format phone number for WhatsApp
   * @param {string} phone
   * @returns {string}
   */
  function formatPhoneForWhatsApp(phone) {
    if (!phone) return "";
    // Remove any non-numeric characters except +
    let cleaned = phone.replace(/[^\d+]/g, "");
    // If it starts with 0, replace with 62 (Indonesia country code)
    if (cleaned.startsWith("0")) {
      cleaned = "62" + cleaned.slice(1);
    }
    return cleaned;
  }

  /**
   * Open WhatsApp with patient phone
   */
  function openWhatsApp() {
    const phone = formatPhoneForWhatsApp(patient?.phone || "");
    if (phone) {
      window.open(`https://wa.me/${phone}`, "_blank");
    } else {
      toastStore.add($t("patients.invalidPhoneError") || "Invalid phone number", { type: "error" });
    }
  }

  /**
   * Get real-time status for a reminder
   * @param {Object} reminder
   * @returns {string}
   */
  function getReminderStatus(reminder) {
    const realtimeStatus = deliveryStatuses[reminder.id];
    return realtimeStatus?.status || reminder.delivery_status || "pending";
  }

  /**
   * Get failure reason for a reminder
   * @param {string} reminderId
   * @returns {string|null}
   */
  function getFailureReason(reminderId) {
    const failedReminder = failedReminders.find(
      (fr) => fr.reminderId === reminderId
    );
    return failedReminder?.error || null;
  }

  /**
   * Check if reminder is failed
   * @param {Object} reminder
   * @returns {boolean}
   */
  function isReminderFailed(reminder) {
    return getReminderStatus(reminder) === "failed";
  }

  /**
   * Get priority color class
   * @param {string} priority
   * @returns {string}
   */
  function getPriorityColor(priority) {
    const colors = {
      high: "bg-red-100 text-red-700",
      medium: "bg-amber-100 text-amber-700",
      low: "bg-green-100 text-green-700"
    };
    return colors[priority] || colors.medium;
  }

  /**
   * Handle edit patient click
   */
  function handleEditPatient() {
    if (onEditPatient) {
      onEditPatient();
    }
  }

  /**
   * Handle add reminder click
   */
  function handleAddReminder() {
    if (onAddReminder) {
      onAddReminder();
    }
  }

  /**
   * Handle send reminder click
   * @param {Object} reminder
   */
  function handleSendReminder(reminder) {
    if (onSendReminder) {
      onSendReminder(reminder);
    }
  }

  /**
   * Handle retry reminder click
   * @param {Object} reminder
   */
  function handleRetryReminder(reminder) {
    if (onRetryReminder) {
      onRetryReminder(reminder);
    }
  }
</script>

<div class="flex flex-col h-full bg-white">
  {#if !patient}
    <!-- Empty State: No patient selected -->
    <div class="flex-1 flex items-center justify-center p-4">
      <div class="text-center max-w-sm">
        <div
          class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-4"
        >
          <svg
            class="w-8 h-8 text-slate-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
            />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-slate-900 mb-2">
          {$t("patients.selectPatient")}
        </h3>
        <p class="text-slate-500">
          {$t("patients.selectPatientDescription")}
        </p>
      </div>
    </div>
  {:else}
    <!-- Patient Header -->
    <div class="flex-shrink-0 px-4 py-4 border-b border-slate-200 bg-slate-50">
      <div class="flex items-start justify-between gap-3">
        <div class="flex items-center gap-3 min-w-0">
          <div
            class="w-12 h-12 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-lg flex-shrink-0"
          >
            {patient.name?.charAt(0).toUpperCase() || "?"}
          </div>
          <div class="min-w-0 flex-1">
            <h2 class="text-lg font-semibold text-slate-900 truncate">
              {patient.name}
            </h2>
            <p class="text-sm text-slate-500 truncate">
              {patient.phone || "-"}
            </p>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center gap-1 flex-shrink-0">
          <button
            onclick={openWhatsApp}
            class="p-2 text-green-600 hover:bg-green-50 rounded-lg transition-colors {patient.phone ? '' : 'opacity-50 cursor-not-allowed'}"
            title={patient.phone ? $t("patients.whatsappChat") : $t("patients.whatsappNoPhone")}
            aria-label={patient.phone ? $t("patients.whatsappOpen") : $t("patients.whatsappNoPhoneAvailable")}
          >
            <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
              <path
                d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"
              />
            </svg>
          </button>
          {#if onEditPatient}
            <button
              onclick={handleEditPatient}
              class="p-2 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
              title={$t("common.edit")}
              aria-label={$t("common.edit")}
            >
              <svg
                class="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                />
              </svg>
            </button>
          {/if}
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div
      class="flex-shrink-0 border-b border-slate-200 bg-white px-4"
      role="tablist"
    >
      {#each tabs as tab}
        <button
          onclick={() => onTabChange(tab.id)}
          role="tab"
          aria-selected={activeTab === tab.id}
          aria-controls="panel-{tab.id}"
          class="px-4 py-3 text-sm font-medium border-b-2 transition-colors flex items-center gap-2 {activeTab ===
          tab.id
            ? 'text-teal-600 border-teal-600'
            : 'text-slate-500 border-transparent hover:text-slate-700'}"
        >
          {#if tab.icon === "info"}
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          {:else if tab.icon === "reminder"}
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          {:else if tab.icon === "history"}
            <svg
              class="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
              />
            </svg>
          {/if}
          {$t(tab.label)}
        </button>
      {/each}
    </div>

    <!-- Tab Content -->
    <div
      id="panel-info"
      role="tabpanel"
      aria-labelledby="tab-info"
      class="flex-1 overflow-y-auto p-4 {activeTab !== 'info' ? 'hidden' : ''}"
    >
      {#if activeTab === "info"}
        <div class="space-y-6">
          <div class="grid grid-cols-1 gap-4">
            <div>
              <label
                class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1"
              >
                {$t("patients.patientName")}
              </label>
              <p class="text-slate-900">{patient.name || "-"}</p>
            </div>
            <div>
              <label
                class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1"
              >
                {$t("patients.phone")}
              </label>
              <p class="text-slate-900">{patient.phone || "-"}</p>
            </div>
            <div>
              <label
                class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1"
              >
                {$t("patients.email")}
              </label>
              <p class="text-slate-900">{patient.email || "-"}</p>
            </div>
            {#if patient.notes}
              <div>
                <label
                  class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1"
                >
                  {$t("patients.notes")}
                </label>
                <p class="text-slate-900">{patient.notes}</p>
              </div>
            {/if}
            <div>
              <label
                class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1"
              >
                {$t("patients.createdAt")}
              </label>
              <p class="text-slate-900">{formatDate(patient.createdAt)}</p>
            </div>
          </div>
        </div>
      {/if}
    </div>

    <div
      id="panel-reminders"
      role="tabpanel"
      aria-labelledby="tab-reminders"
      class="flex-1 overflow-y-auto p-4 {activeTab !== 'reminders' ? 'hidden' : ''}"
    >
      {#if activeTab === "reminders"}
        {#if !patient.reminders || patient.reminders.length === 0}
          <!-- Empty State: No reminders -->
          <div class="text-center py-8">
            <div
              class="w-12 h-12 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-3"
            >
              <svg
                class="w-6 h-6 text-slate-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
                />
              </svg>
            </div>
            <h3 class="text-sm font-semibold text-slate-900 mb-1">
              {$t("patients.noReminders")}
            </h3>
            <p class="text-xs text-slate-500 mb-4">
              {$t("patients.noRemindersDescription")}
            </p>
            {#if onAddReminder}
              <button
                onclick={handleAddReminder}
                class="px-4 py-2 bg-teal-600 text-white text-sm font-medium rounded-lg hover:bg-teal-700 transition-colors"
              >
                {$t("patients.addReminder")}
              </button>
            {/if}
          </div>
        {:else}
          <!-- Reminders List -->
          <div class="space-y-3">
            {#each patient.reminders as reminder (reminder.id)}
              {@const status = getReminderStatus(reminder)}
              {@const isFailed = status === "failed"}
              {@const failureReason = getFailureReason(reminder.id)}
              {@const isSending = status === "sending"}

              <div
                class="p-3 rounded-xl border-2 transition-colors {isFailed
                  ? 'bg-red-50 border-red-500'
                  : reminder.completed
                    ? 'bg-slate-50 border-transparent'
                    : 'bg-amber-50 border-transparent'}"
              >
                <div class="flex items-start gap-2">
                  <!-- Toggle Button -->
                  {#if onToggleReminder}
                    <button
                      onclick={() =>
                        onToggleReminder(patient.id, reminder.id)}
                      class="w-5 h-5 rounded-full border-2 flex-shrink-0 mt-0.5 transition-colors {reminder.completed
                        ? 'bg-teal-600 border-teal-600'
                        : 'border-slate-300 hover:border-teal-600'}"
                      aria-label={reminder.completed
                        ? $t("patients.markAsIncomplete")
                        : $t("patients.markAsComplete")}
                    >
                      {#if reminder.completed}
                        <svg
                          class="w-full h-full text-white"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="3"
                            d="M5 13l4 4L19 7"
                          />
                        </svg>
                      {/if}
                    </button>
                  {/if}

                  <!-- Reminder Content -->
                  <div class="flex-1 min-w-0">
                    <div class="flex items-start justify-between gap-2">
                      <div class="min-w-0 flex-1">
                        <span
                          class="text-sm block {reminder.completed
                            ? 'text-slate-500 line-through'
                            : isFailed
                              ? 'text-red-900 font-medium'
                              : 'text-slate-900'}"
                        >
                          {reminder.title}
                        </span>
                        {#if reminder.dueDate}
                          <span class="text-xs text-slate-400">
                            {formatDate(reminder.dueDate)}
                          </span>
                        {/if}
                      </div>

                      <!-- Priority Badge -->
                      <span
                        class="flex-shrink-0 px-2 py-0.5 text-xs font-medium rounded-full {getPriorityColor(
                          reminder.priority
                        )}"
                      >
                        {reminder.priority}
                      </span>
                    </div>

                    <!-- Actions Row -->
                    <div class="flex items-center gap-1 mt-2">
                      <!-- Delivery Status -->
                      <DeliveryStatusBadge
                        {status}
                        onRetry={() => handleRetryReminder(reminder)}
                        isRetrying={isSending}
                      />

                      <!-- Send Button -->
                      {#if onSendReminder}
                        <button
                          onclick={() => handleSendReminder(reminder)}
                          disabled={isSending}
                          class="px-3 py-1.5 text-xs font-medium rounded-lg transition-colors flex items-center gap-1 {isFailed
                            ? 'text-white bg-red-600 hover:bg-red-700'
                            : 'text-green-600 bg-green-50 hover:bg-green-100 border border-green-200'} {isSending
                            ? 'opacity-50 cursor-not-allowed'
                            : ''}"
                          title={isFailed
                            ? $t("reminder.retry")
                            : $t("reminder.send.title")}
                        >
                          {#if isFailed}
                            <svg
                              class="w-3.5 h-3.5"
                              fill="none"
                              stroke="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                              />
                            </svg>
                            Retry
                          {:else}
                            <svg
                              class="w-3.5 h-3.5"
                              viewBox="0 0 24 24"
                              fill="currentColor"
                            >
                              <path
                                d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"
                              />
                            </svg>
                            Send
                          {/if}
                        </button>
                      {/if}

                      <!-- Edit Button -->
                      {#if onEditReminder}
                        <button
                          onclick={() => onEditReminder(reminder)}
                          class="px-3 py-1.5 text-xs font-medium text-slate-600 bg-slate-50 hover:bg-slate-100 rounded-lg transition-colors flex items-center gap-1"
                          title={$t("common.edit")}
                        >
                          <svg
                            class="w-3.5 h-3.5"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                            />
                          </svg>
                          {$t("common.edit")}
                        </button>
                      {/if}

                      <!-- Delete Button -->
                      {#if onDeleteReminder}
                        <button
                          onclick={() =>
                            onDeleteReminder(patient.id, reminder.id)}
                          class="px-3 py-1.5 text-xs font-medium text-red-600 bg-red-50 hover:bg-red-100 rounded-lg transition-colors flex items-center gap-1"
                          title={$t("common.delete")}
                        >
                          <svg
                            class="w-3.5 h-3.5"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M6 18L18 6M6 6l12 12"
                            />
                          </svg>
                          {$t("common.delete")}
                        </button>
                      {/if}
                    </div>

                    <!-- Error Message -->
                    {#if failureReason}
                      <p class="mt-2 text-xs text-red-600">{failureReason}</p>
                    {/if}
                  </div>
                </div>
              </div>
            {/each}

            <!-- Add Reminder Button -->
            {#if onAddReminder}
              <button
                onclick={handleAddReminder}
                class="w-full py-3 border-2 border-dashed border-slate-200 rounded-xl text-slate-500 hover:text-slate-700 hover:border-slate-300 transition-colors flex items-center justify-center gap-2"
              >
                <svg
                  class="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                  />
                </svg>
                {$t("patients.addReminder")}
              </button>
            {/if}
          </div>
        {/if}
      {/if}
    </div>

    <div
      id="panel-history"
      role="tabpanel"
      aria-labelledby="tab-history"
      class="flex-1 overflow-y-auto {activeTab !== 'history' ? 'hidden' : ''}"
    >
      {#if activeTab === "history" && token}
        <ReminderHistoryView
          {token}
          patientId={patient.id}
          {patient}
          onCreateReminder={onAddReminder}
        />
      {:else if activeTab === "history"}
        <div class="p-4 text-center text-slate-500">
          <p class="text-sm">{$t("patients.tokenRequired")}</p>
        </div>
      {/if}
    </div>
  {/if}
</div>
