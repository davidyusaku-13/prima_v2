<script>
  import { t } from "svelte-i18n";
  import WhatsAppPreview from "$lib/components/whatsapp/WhatsAppPreview.svelte";
  import QuietHoursHint from "$lib/components/indicators/QuietHoursHint.svelte";

  let {
    show = false,
    patient = null,
    reminder = null,
    status = "idle", // 'idle' | 'sending' | 'success' | 'error' | 'scheduled'
    errorMessage = "",
    scheduledTime = null,
    isQuietHours = false,
    onClose = () => {},
    onConfirm = () => {},
  } = $props();
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
      onclick={onClose}
      onkeydown={(e) => e.key === "Escape" && onClose()}
      role="button"
      tabindex="0"
      aria-label="Close modal"
    ></div>
    <div
      class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-md p-4 sm:p-6"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <!-- Header with WhatsApp icon -->
      <div class="flex items-center gap-3 sm:gap-4 mb-4 sm:mb-6">
        <div
          class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center flex-shrink-0"
        >
          <svg
            class="w-6 h-6 text-green-600"
            viewBox="0 0 24 24"
            fill="currentColor"
          >
            <path
              d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"
            />
          </svg>
        </div>
        <div>
          <h2 class="text-base sm:text-lg font-semibold text-slate-900">
            {$t("reminder.send.title")}
          </h2>
          {#if patient}
            <p class="text-sm text-slate-500">{patient.phone}</p>
          {/if}
        </div>
      </div>

      <!-- Patient name prominently displayed -->
      {#if patient}
        <div class="mb-4 p-3 bg-slate-50 rounded-xl">
          <p class="text-xs text-slate-500 mb-1">
            {$t("patients.patient") || "Patient"}
          </p>
          <p class="text-lg font-semibold text-slate-900">{patient.name}</p>
        </div>
      {/if}

      <!-- WhatsApp Message Preview -->
      {#if reminder && patient}
        <div class="mb-4">
          <WhatsAppPreview
            message={reminder.description || ""}
            patientName={patient.name}
            reminderTitle={reminder.title}
            isScheduled={isQuietHours}
            {scheduledTime}
          />
        </div>
      {/if}

      <!-- Quiet Hours Hint -->
      {#if isQuietHours}
        <div class="mb-4">
          <QuietHoursHint {scheduledTime} />
        </div>
      {/if}

      <!-- Confirmation message -->
      <p class="text-sm text-slate-700 mb-4">
        {$t("reminder.send.confirm", {
          values: { patientName: patient?.name || "" },
        })}
      </p>

      <!-- Error message -->
      {#if status === "error" && errorMessage}
        <div class="mb-4 p-3 bg-red-50 border border-red-100 rounded-xl">
          <p class="text-sm text-red-600">{errorMessage}</p>
        </div>
      {/if}

      <!-- Buttons -->
      <div class="flex flex-col sm:flex-row gap-3">
        <button
          onclick={onClose}
          disabled={status === "sending"}
          class="flex-1 h-10 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
        >
          {$t("reminder.send.cancel")}
        </button>
        <button
          onclick={onConfirm}
          disabled={status === "sending"}
          class="flex-1 h-10 px-4 py-2.5 bg-green-600 text-white font-medium rounded-xl hover:bg-green-700 transition-colors duration-200 order-1 sm:order-2 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          {#if status === "sending"}
            <svg
              class="animate-spin h-4 w-4"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            {$t("reminder.send.sending")}
          {:else}
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
              <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
            </svg>
            {$t("reminder.send.button")}
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
