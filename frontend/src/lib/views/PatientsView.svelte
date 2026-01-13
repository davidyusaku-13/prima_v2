<script>
  import { onMount } from "svelte";
  import { t } from "svelte-i18n";
  import { locale } from "svelte-i18n";
  import { deliveryStore } from "$lib/stores/delivery.svelte.js";
  import PatientListPane from "$lib/components/patients/PatientListPane.svelte";
  import PatientDetailPane from "$lib/components/patients/PatientDetailPane.svelte";
  import ConfirmModal from "$lib/components/ConfirmModal.svelte";

  let {
    patients = [],
    token = "",
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

  // Search query state
  let searchQuery = $state("");

  // Selected patient state with session storage persistence
  let selectedPatientId = $state(null);
  let showMobileDetailPane = $state(false);

  // Detail pane tab state
  let activeTab = $state("info");

  // Delete confirmation modal state
  let showDeleteModal = $state(false);
  let pendingDeleteId = $state(null);

  // Initialize selectedPatientId from session storage on mount
  onMount(() => {
    try {
      const stored = sessionStorage.getItem("prima-selectedPatientId");
      if (stored) {
        // Validate that the patient exists in the current list
        const patientExists = patients.some((p) => p.id === stored);
        if (patientExists) {
          selectedPatientId = stored;
        } else {
          // Clear stale patientId
          sessionStorage.removeItem("prima-selectedPatientId");
        }
      }
    } catch (e) {
      // Handle private mode or quota exceeded errors gracefully
      console.warn("Could not read session storage:", e.message);
    }
  });

  // Save selectedPatientId to session storage when it changes
  $effect(() => {
    try {
      if (selectedPatientId) {
        sessionStorage.setItem("prima-selectedPatientId", selectedPatientId);
      } else {
        sessionStorage.removeItem("prima-selectedPatientId");
      }
    } catch (e) {
      // Handle private mode or quota exceeded errors gracefully
      console.warn("Could not write to session storage:", e.message);
    }
  });

  // Compute selected patient from patients array
  let selectedPatient = $derived(
    selectedPatientId
      ? patients.find((p) => p.id === selectedPatientId) || null
      : null
  );

  // Filter patients based on search query
  let filteredPatients = $derived(
    patients.filter(
      (p) =>
        p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        p.phone.includes(searchQuery) ||
        p.email.toLowerCase().includes(searchQuery.toLowerCase())
    )
  );

  // Initialize SSE connection on mount
  onMount(() => {
    deliveryStore.connect();

    // Hydrate delivery store with existing delivery statuses from backend data
    // This ensures filter counts are accurate on page load
    patients.forEach(p => {
      p.reminders?.forEach(r => {
        if (r.delivery_status) {
          deliveryStore.updateStatus(r.id, r.delivery_status, r.message_sent_at || new Date().toISOString());
        }
      });
    });

    // Listen for navigate-to-patient event from toast action
    const handleNavigateToPatient = (event) => {
      const { patientId } = event.detail;
      handleSelectPatient(patientId);
    };

    // Listen for show-failed-reminders event from FailedReminderBadge
    const handleShowFailedReminders = () => {
      // Could filter to show only failed reminders in the future
    };

    window.addEventListener("navigate-to-patient", handleNavigateToPatient);
    window.addEventListener("show-failed-reminders", handleShowFailedReminders);

    // Cleanup on unmount
    return () => {
      deliveryStore.disconnect();
      window.removeEventListener("navigate-to-patient", handleNavigateToPatient);
      window.removeEventListener(
        "show-failed-reminders",
        handleShowFailedReminders
      );
    };
  });

  /**
   * Handle patient selection
   * @param {string} patientId
   */
  function handleSelectPatient(patientId) {
    selectedPatientId = patientId;
    // On mobile, show the detail pane as slide-over
    if (typeof window !== "undefined" && window.innerWidth < 768) {
      showMobileDetailPane = true;
    }
  }

  /**
   * Handle back button on mobile
   */
  function handleBack() {
    showMobileDetailPane = false;
    // Small delay to clear selection after animation
    setTimeout(() => {
      selectedPatientId = null;
    }, 300);
  }

  /**
   * Handle tab change in detail pane
   * @param {string} tabId
   */
  function handleTabChange(tabId) {
    activeTab = tabId;
  }

  /**
   * Handle add patient button click
   */
  function handleOpenPatientModal() {
    onOpenPatientModal();
  }

  /**
   * Handle edit patient action
   */
  function handleEditPatient() {
    if (selectedPatient) {
      onOpenPatientModal(selectedPatient);
    }
  }

  /**
   * Handle delete patient action
   * @param {string} patientId
   */
  function handleDeletePatient(patientId) {
    pendingDeleteId = patientId;
    showDeleteModal = true;
  }

  /**
   * Confirm patient deletion from modal
   */
  function confirmDeletePatient() {
    if (pendingDeleteId) {
      onDeletePatient(pendingDeleteId);
      // Clear selection if deleted patient was selected
      if (selectedPatientId === pendingDeleteId) {
        selectedPatientId = null;
        showMobileDetailPane = false;
      }
    }
    showDeleteModal = false;
    pendingDeleteId = null;
  }

  /**
   * Handle add reminder action
   */
  function handleAddReminder() {
    if (selectedPatient) {
      onOpenReminderModal(selectedPatient.id);
    }
  }

  /**
   * Handle edit reminder action
   * @param {Object} reminder
   */
  function handleEditReminder(reminder) {
    if (selectedPatient) {
      onOpenReminderModal(selectedPatient.id, reminder);
    }
  }
</script>

<div class="h-full flex flex-col">
  <!-- Header -->
  <header
    class="sticky top-0 z-40 bg-white/80 backdrop-blur-md border-b border-slate-200"
  >
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3 py-3 px-4">
      <div class="flex flex-col sm:flex-row sm:items-center gap-3">
        <div class="flex items-center gap-2">
          <h1 class="text-lg sm:text-xl font-bold text-slate-900">
            {$t("patients.title")}
          </h1>
          <!-- SSE Connection Status Indicator -->
          {#if connectionStatus === "connected"}
            <span
              class="flex items-center gap-1 px-2 py-0.5 text-xs text-green-700 bg-green-50 rounded-full"
              title={$t("delivery.status.connected")}
            >
              <span
                class="w-1.5 h-1.5 bg-green-500 rounded-full animate-pulse"
              ></span>
              {$t("delivery.status.live")}
            </span>
          {:else if connectionStatus === "connecting"}
            <span
              class="flex items-center gap-1 px-2 py-0.5 text-xs text-amber-700 bg-amber-50 rounded-full"
              title={$t("delivery.status.connectingStatus")}
            >
              <span
                class="w-1.5 h-1.5 bg-amber-500 rounded-full animate-pulse"
              ></span>
              {$t("delivery.status.connecting")}
            </span>
          {:else if connectionStatus === "disconnected"}
            <span
              class="flex items-center gap-1 px-2 py-0.5 text-xs text-slate-500 bg-slate-100 rounded-full"
              title={$t("delivery.status.disconnected")}
            >
              <span class="w-1.5 h-1.5 bg-slate-400 rounded-full"></span>
              {$t("delivery.status.offline")}
            </span>
          {/if}
        </div>
        <div class="relative w-full sm:w-56 md:w-64 lg:w-72 xl:w-80">
          <svg
            class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 sm:w-5 sm:h-5 text-slate-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
          <input
            type="text"
            bind:value={searchQuery}
            placeholder={$t("common.searchPlaceholder")}
            aria-label={$t("common.searchPlaceholder")}
            class="pl-9 sm:pl-10 pr-4 py-2 sm:py-2.5 w-full bg-slate-100 border-0 rounded-lg sm:rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
          />
        </div>
      </div>
      <button
        onclick={handleOpenPatientModal}
        class="flex items-center justify-center gap-2 px-4 sm:px-5 py-2 sm:py-2.5 bg-teal-600 text-white font-medium rounded-lg sm:rounded-xl hover:bg-teal-700 hover:shadow-lg active:scale-[0.98] transition-all duration-200 w-full sm:w-auto text-sm sm:text-base"
      >
        <svg
          class="w-4 h-4 sm:w-5 sm:h-5"
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
        {$t("patients.addPatient")}
      </button>
    </div>
  </header>

  <!-- Two-Pane Layout -->
  <div class="flex-1 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-12 overflow-hidden">
    <!-- Left Pane: Patient List -->
    <div
      class="lg:col-span-5 bg-white border-r border-slate-200 overflow-hidden {showMobileDetailPane
        ? 'hidden lg:block'
        : 'block'}"
    >
      <PatientListPane
        patients={filteredPatients}
        {selectedPatientId}
        {searchQuery}
        onSelect={handleSelectPatient}
        onAddPatient={handleOpenPatientModal}
        onEditPatient={handleEditPatient}
        onDeletePatient={handleDeletePatient}
      />
    </div>

    <!-- Right Pane: Patient Detail -->
    <div
      class="lg:col-span-7 bg-slate-50 overflow-hidden transition-transform duration-300 ease-out {showMobileDetailPane
        ? 'translate-x-0'
        : 'translate-x-full lg:translate-x-0 lg:block hidden'}"
    >
      <!-- Mobile Back Button -->
      <div
        class="lg:hidden flex items-center gap-2 px-4 py-3 bg-white border-b border-slate-200 {showMobileDetailPane
          ? 'opacity-100'
          : 'opacity-0 pointer-events-none'}"
      >
        <button
          onclick={handleBack}
          class="flex items-center gap-2 text-slate-600 hover:text-slate-900 transition-colors"
          aria-label={$t("patients.backToPatientList")}
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
              d="M15 19l-7-7 7-7"
            />
          </svg>
          <span class="text-sm font-medium">{$t("patients.back")}</span>
        </button>
      </div>

      <!-- Mobile Slide-Over Backdrop -->
      <div
        class="lg:hidden fixed inset-0 bg-black/20 transition-opacity duration-300 {showMobileDetailPane
          ? 'opacity-100'
          : 'opacity-0 pointer-events-none'}"
        style="z-index: -1"
        onclick={handleBack}
        aria-hidden="true"
      ></div>

      <!-- Detail Pane -->
      <div class="h-full overflow-hidden bg-white lg:bg-slate-50">
        <PatientDetailPane
          patient={selectedPatient}
          {token}
          {deliveryStatuses}
          {failedReminders}
          {activeTab}
          onTabChange={handleTabChange}
          onEditPatient={handleEditPatient}
          onAddReminder={handleAddReminder}
          onEditReminder={handleEditReminder}
          onToggleReminder={onToggleReminder}
          onDeleteReminder={onDeleteReminder}
          onSendReminder={onSendReminder}
          onRetryReminder={onRetryReminder}
        />
      </div>
    </div>
  </div>

  <!-- Delete Confirmation Modal -->
  <ConfirmModal
    show={showDeleteModal}
    message={$t("patients.deleteConfirmation")}
    onClose={() => {
      showDeleteModal = false;
      pendingDeleteId = null;
    }}
    onConfirm={confirmDeletePatient}
  />
</div>
