<script>
  import { t } from "svelte-i18n";
  import { locale } from "svelte-i18n";
  import DeliveryStatusBadge from "$lib/components/delivery/DeliveryStatusBadge.svelte";
  import { deliveryStore } from "$lib/stores/delivery.svelte.js";

  let {
    filteredPatients = [],
    searchQuery = "",
    deliveryStatuses = {},
    failedReminders = [],
    onOpenPatientModal = () => {},
    onOpenReminderModal = () => {},
    onDeletePatient = () => {},
    onToggleReminder = () => {},
    onDeleteReminder = () => {},
    onSendReminder = () => {},
    onRetryReminder = () => {},
  } = $props();

  function formatDate(dateStr) {
    if (!dateStr) return "";
    const lang = $locale || "en";
    return new Date(dateStr).toLocaleDateString(lang, {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  function formatRecurrence(recurrence) {
    if (!recurrence || recurrence.frequency === "none") return "";
    const freqLabels = {
      daily: "Daily",
      weekly: "Weekly",
      monthly: "Monthly",
      yearly: "Yearly",
    };
    let label = freqLabels[recurrence.frequency] || recurrence.frequency;
    if (recurrence.interval > 1) {
      label = `Every ${recurrence.interval} ${recurrence.frequency}s`;
    }
    if (
      recurrence.frequency === "weekly" &&
      recurrence.daysOfWeek?.length > 0
    ) {
      const dayLabels = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];
      const days = recurrence.daysOfWeek.map((d) => dayLabels[d]).join(", ");
      label += ` on ${days}`;
    }
    return label;
  }

  function getPriorityColor(priority) {
    const colors = {
      high: "bg-red-100 text-red-700",
      medium: "bg-amber-100 text-amber-700",
      low: "bg-green-100 text-green-700",
    };
    return colors[priority] || colors.medium;
  }

  function getReminderStatus(reminder) {
    const realtimeStatus = deliveryStatuses[reminder.id];
    return realtimeStatus?.status || reminder.delivery_status || "pending";
  }

  function getFailureReason(reminderId) {
    const failedReminder = failedReminders.find(
      (fr) => fr.reminderId === reminderId
    );
    return failedReminder?.error || null;
  }

  function isReminderFailed(reminder) {
    const status = getReminderStatus(reminder);
    return status === "failed";
  }

  // Count reminders per patient
  function countPatientReminders(patient, visibleReminders = null) {
    if (visibleReminders) {
      return visibleReminders.filter((r) => r.patientId === patient.id).length;
    }
    return patient.reminders?.length || 0;
  }
</script>

{#if filteredPatients.length === 0}
  <div
    class="bg-white rounded-xl sm:rounded-2xl border border-slate-200 p-8 sm:p-12 text-center"
  >
    <div
      class="w-12 h-12 sm:w-16 sm:h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4"
    >
      <svg
        class="w-6 h-6 sm:w-8 sm:h-8 text-slate-400"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
        />
      </svg>
    </div>
    <h3 class="text-base sm:text-lg font-semibold text-slate-900 mb-2">
      {$t("patients.noPatients")}
    </h3>
    <p class="text-slate-500 mb-4 sm:mb-6 text-sm">
      {#if searchQuery}
        {$t("patients.noPatientsMatch")}
      {:else}
        {$t("patients.getStarted")}
      {/if}
    </p>
    <button
      onclick={() => onOpenPatientModal()}
      class="px-5 sm:px-6 py-2.5 sm:py-3 bg-teal-600 text-white font-medium rounded-lg sm:rounded-xl hover:bg-teal-700 transition-colors duration-200 text-sm sm:text-base"
    >
      {$t("patients.addPatient")}
    </button>
  </div>
{:else}
  <div class="space-y-3 sm:space-y-4">
    {#each filteredPatients as patient}
      <div
        id="patient-{patient.id}"
        class="bg-white rounded-xl sm:rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg transition-all duration-200"
      >
        <div class="p-3 sm:p-4 md:p-5 lg:p-6">
          <div class="flex items-start gap-2.5 sm:gap-3 md:gap-4">
            <!-- Avatar -->
            <div
              class="w-10 h-10 sm:w-12 sm:h-12 md:w-14 md:h-14 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-base sm:text-lg md:text-xl flex-shrink-0"
            >
              {patient.name.charAt(0).toUpperCase()}
            </div>

            <!-- Patient Info + Actions -->
            <div class="flex-1 min-w-0">
              <div
                class="flex flex-col sm:flex-row sm:items-start justify-between gap-2 sm:gap-4"
              >
                <div class="min-w-0 flex-1">
                  <h3
                    class="text-base sm:text-lg font-semibold text-slate-900 truncate"
                  >
                    {patient.name}
                  </h3>
                  {#if patient.phone}
                    <p class="text-slate-600 text-xs sm:text-sm">
                      {patient.phone}
                    </p>
                  {/if}
                  {#if patient.email}
                    <p class="text-slate-500 text-xs sm:text-sm truncate">
                      {patient.email}
                    </p>
                  {/if}
                  {#if patient.notes}
                    <p
                      class="text-slate-500 text-xs sm:text-sm mt-1.5 sm:mt-2 line-clamp-2"
                    >
                      {patient.notes}
                    </p>
                  {/if}

                  <!-- Reminder count badge -->
                  {#if patient.reminders && patient.reminders.length > 0}
                    <div class="flex items-center gap-2 mt-2 sm:mt-3">
                      <span
                        class="inline-flex items-center gap-1 px-2.5 sm:px-3 py-1 text-xs sm:text-sm font-medium rounded-full bg-teal-50 text-teal-700 border border-teal-200"
                      >
                        <svg
                          class="w-3.5 h-3.5 sm:w-4 sm:h-4"
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
                        {patient.reminders.length}
                        {patient.reminders.length === 1
                          ? "reminder"
                          : "reminders"}
                      </span>
                    </div>
                  {/if}
                </div>

                <!-- Action buttons -->
                <div class="flex items-center gap-1 sm:gap-2 flex-shrink-0">
                  <button
                    onclick={() => onOpenReminderModal(patient.id)}
                    class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium text-teal-600 bg-teal-50 hover:bg-teal-100 border border-teal-200 transition-colors duration-200 flex items-center gap-1"
                    title={$t("patients.addReminder")}
                  >
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
                        d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                      />
                    </svg>
                    {$t("patients.addReminder")}
                  </button>
                  <button
                    onclick={() => onOpenPatientModal(patient)}
                    class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium text-slate-600 bg-slate-50 hover:bg-slate-100 border border-slate-200 transition-colors duration-200 flex items-center gap-1"
                    title={$t("common.edit")}
                  >
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
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                      />
                    </svg>
                    {$t("common.edit")}
                  </button>
                  <button
                    onclick={() => onDeletePatient(patient.id)}
                    class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium text-red-600 bg-red-50 hover:bg-red-100 border border-red-200 transition-colors duration-200 flex items-center gap-1"
                    title={$t("common.delete")}
                  >
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
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                    {$t("common.delete")}
                  </button>
                </div>
              </div>

              <!-- Collapsed reminders summary -->
              {#if patient.reminders && patient.reminders.length > 0}
                <div
                  class="mt-3 sm:mt-4 pt-3 sm:pt-4 border-t border-slate-100"
                >
                  <details class="cursor-pointer">
                    <summary
                      class="text-xs sm:text-sm font-medium text-slate-700 hover:text-slate-900 select-none flex items-center gap-2"
                    >
                      <svg
                        class="w-4 h-4 transition-transform duration-200"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M9 5l7 7-7 7"
                        />
                      </svg>
                      {$t("patients.reminders")} ({patient.reminders.length})
                    </summary>

                    <div class="mt-3 space-y-1.5 sm:space-y-2">
                      {#each patient.reminders as reminder}
                        {@const isFailed = isReminderFailed(reminder)}
                        {@const failureReason = getFailureReason(reminder.id)}
                        <div
                          class="flex flex-wrap items-center gap-1.5 sm:gap-2 p-2 sm:p-3 rounded-lg border-2 {isFailed
                            ? 'bg-red-50 border-red-500'
                            : reminder.completed
                              ? 'bg-slate-50 border-transparent'
                              : 'bg-amber-50 border-transparent'}"
                        >
                          <button
                            onclick={() =>
                              onToggleReminder(patient.id, reminder.id)}
                            class="w-4 h-4 sm:w-5 sm:h-5 rounded-full border-2 flex-shrink-0 transition-colors duration-200 {reminder.completed
                              ? 'bg-teal-600 border-teal-600'
                              : 'border-slate-300 hover:border-teal-600'}"
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
                          <div class="flex-1 min-w-0">
                            <span
                              class="text-xs sm:text-sm {reminder.completed
                                ? 'text-slate-500 line-through'
                                : isFailed
                                  ? 'text-red-900 font-medium'
                                  : 'text-slate-900'} truncate block"
                              >{reminder.title}</span
                            >
                            {#if reminder.dueDate}
                              <span
                                class="text-[10px] sm:text-xs text-slate-400"
                                >{formatDate(reminder.dueDate)}</span
                              >
                            {/if}
                          </div>
                          <span
                            class="px-1.5 sm:px-2 py-0.5 text-[10px] sm:text-xs font-medium rounded-full {getPriorityColor(
                              reminder.priority
                            )}">{reminder.priority}</span
                          >
                          {#if reminder.delivery_status || deliveryStatuses[reminder.id]}
                            <DeliveryStatusBadge
                              status={getReminderStatus(reminder)}
                              onRetry={() => onRetryReminder(patient, reminder)}
                              isRetrying={getReminderStatus(reminder) ===
                                "sending"}
                            />
                          {/if}
                          <div class="flex items-center gap-1 sm:gap-2">
                            <button
                              onclick={() => onSendReminder(patient, reminder)}
                              class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium transition-colors duration-200 {isFailed
                                ? 'text-white bg-red-600 hover:bg-red-700'
                                : 'text-green-600 bg-green-50 hover:bg-green-100 border border-green-200'} {reminder.delivery_status ===
                              'sending'
                                ? 'opacity-50 cursor-not-allowed'
                                : ''} flex items-center gap-1"
                              disabled={reminder.delivery_status === "sending"}
                              title={isFailed
                                ? $t("reminder.retry")
                                : $t("reminder.send.title")}
                            >
                              {#if isFailed}
                                <span class="flex items-center gap-1">
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
                                      d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                                    />
                                  </svg>
                                  Retry
                                </span>
                              {:else}
                                <span class="flex items-center gap-1">
                                  <svg
                                    class="w-4 h-4"
                                    viewBox="0 0 24 24"
                                    fill="currentColor"
                                  >
                                    <path
                                      d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"
                                    />
                                  </svg>
                                  Send
                                </span>
                              {/if}
                            </button>
                            <button
                              onclick={() =>
                                onOpenReminderModal(patient.id, reminder)}
                              class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium text-slate-600 bg-slate-50 hover:bg-slate-100 border border-slate-200 transition-colors duration-200 flex items-center gap-1"
                              title={$t("common.edit")}
                            >
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
                                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                                />
                              </svg>
                              {$t("common.edit")}
                            </button>
                            <button
                              onclick={() =>
                                onDeleteReminder(patient.id, reminder.id)}
                              class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg text-xs sm:text-sm font-medium text-red-600 bg-red-50 hover:bg-red-100 border border-red-200 transition-colors duration-200 flex items-center gap-1"
                              title={$t("common.delete")}
                            >
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
                                  d="M6 18L18 6M6 6l12 12"
                                />
                              </svg>
                              {$t("common.delete")}
                            </button>
                          </div>
                        </div>
                      {/each}
                    </div>
                  </details>
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>
    {/each}
  </div>
{/if}
