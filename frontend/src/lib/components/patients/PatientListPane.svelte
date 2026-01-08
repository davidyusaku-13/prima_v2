<script>
  import { t } from "svelte-i18n";

  /**
   * @typedef {Object} Props
   * @property {Array<{
   *   id: string;
   *   name: string;
   *   phone: string;
   *   email: string;
   *   notes: string;
   *   reminders?: Array<{ id: string }>;
   * }>} patients - Array of patient objects
   * @property {string|null} selectedPatientId - Currently selected patient ID
   * @property {string} searchQuery - Current search text
   * @property {Function} onSelect - Callback when patient is selected
   * @property {Function} onAddPatient - Callback for add patient button
   * @property {Function} [onEditPatient] - Callback for edit patient action
   * @property {Function} [onDeletePatient] - Callback for delete patient action
   */

  /** @type {Props} */
  let {
    patients = [],
    selectedPatientId = null,
    searchQuery = "",
    onSelect = () => {},
    onAddPatient = () => {},
    onEditPatient = null,
    onDeletePatient = null
  } = $props();

  /**
   * Handle patient selection
   * @param {string} patientId
   */
  function handleSelect(patientId) {
    onSelect(patientId);
  }

  /**
   * Handle add patient click
   */
  function handleAddPatient() {
    onAddPatient();
  }

  /**
   * Handle edit patient click
   * @param {Event} e
   * @param {string} patientId
   */
  function handleEditPatient(e, patientId) {
    e.stopPropagation();
    if (onEditPatient) {
      onEditPatient(patientId);
    }
  }

  /**
   * Handle delete patient click
   * @param {Event} e
   * @param {string} patientId
   */
  function handleDeletePatient(e, patientId) {
    e.stopPropagation();
    if (onDeletePatient) {
      onDeletePatient(patientId);
    }
  }

  /**
   * Get reminder count for a patient
   * @param {Object} patient
   * @returns {number}
   */
  function getReminderCount(patient) {
    return patient.reminders?.length || 0;
  }
</script>

<div class="flex flex-col h-full">
  <!-- Patient List -->
  {#if patients.length === 0}
    <!-- Empty State -->
    <div class="flex-1 flex items-center justify-center p-4">
      <div class="text-center">
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
              d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
            />
          </svg>
        </div>
        <h3 class="text-sm font-semibold text-slate-900 mb-1">
          {$t("patients.noPatients")}
        </h3>
        <p class="text-xs text-slate-500">
          {#if searchQuery}
            {$t("patients.noPatientsMatch")}
          {:else}
            {$t("patients.getStarted")}
          {/if}
        </p>
      </div>
    </div>
  {:else}
    <div
      class="flex-1 overflow-y-auto"
      role="listbox"
      aria-label={$t("patients.title")}
    >
      {#each patients as patient (patient.id)}
        {@const isSelected = selectedPatientId === patient.id}
        {@const reminderCount = getReminderCount(patient)}

        <button
          onclick={() => handleSelect(patient.id)}
          role="option"
          aria-selected={isSelected}
          aria-label={patient.name}
          class="w-full p-3 text-left border-b border-slate-100 hover:bg-slate-50 transition-colors duration-150 {
            isSelected ? 'bg-teal-50 border-l-4 border-l-teal-600' : 'border-l-4 border-l-transparent'
          }"
        >
          <div class="flex items-start gap-3">
            <!-- Avatar -->
            <div
              class="w-10 h-10 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-sm flex-shrink-0"
              aria-hidden="true"
            >
              {patient.name.charAt(0).toUpperCase()}
            </div>

            <!-- Patient Info -->
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0 flex-1">
                  <h4
                    class="text-sm font-medium text-slate-900 truncate {isSelected ? 'text-teal-900' : ''}"
                  >
                    {patient.name}
                  </h4>
                  {#if patient.phone}
                    <p class="text-xs text-slate-500 truncate">
                      {patient.phone}
                    </p>
                  {/if}
                </div>

                <!-- Reminder Count Badge -->
                {#if reminderCount > 0}
                  <span
                    class="flex-shrink-0 inline-flex items-center gap-1 px-2 py-0.5 text-xs font-medium rounded-full bg-teal-100 text-teal-700"
                    aria-label={$t("patients.reminderCount", { count: reminderCount })}
                  >
                    <svg
                      class="w-3 h-3"
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
                    {reminderCount}
                  </span>
                {/if}
              </div>

              <!-- Action Buttons (Edit/Delete) -->
              {#if isSelected}
                <div class="flex items-center gap-1 mt-2" role="group" aria-label={$t("patients.patientActions")}>
                  {#if onEditPatient}
                    <span
                      onclick={(e) => handleEditPatient(e, patient.id)}
                      onkeydown={(e) => e.key === 'Enter' && handleEditPatient(e, patient.id)}
                      role="button"
                      tabindex="0"
                      aria-label={$t("common.edit")}
                      class="flex-1 h-8 px-3 py-1.5 text-xs font-medium text-slate-600 bg-white border border-slate-200 rounded-lg hover:bg-slate-50 hover:border-slate-300 transition-colors flex items-center justify-center gap-1 cursor-pointer"
                    >
                      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                        />
                      </svg>
                      {$t("common.edit")}
                    </span>
                  {/if}
                  {#if onDeletePatient}
                    <span
                      onclick={(e) => handleDeletePatient(e, patient.id)}
                      onkeydown={(e) => e.key === 'Enter' && handleDeletePatient(e, patient.id)}
                      role="button"
                      tabindex="0"
                      aria-label={$t("common.delete")}
                      class="flex-1 h-8 px-3 py-1.5 text-xs font-medium text-red-600 bg-white border border-slate-200 rounded-lg hover:bg-red-50 hover:border-red-300 transition-colors flex items-center justify-center gap-1 cursor-pointer"
                    >
                      <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path
                          stroke-linecap="round"
                          stroke-linejoin="round"
                          stroke-width="2"
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                      {$t("common.delete")}
                    </span>
                  {/if}
                </div>
              {/if}
            </div>
          </div>
        </button>
      {/each}
    </div>
  {/if}
</div>
