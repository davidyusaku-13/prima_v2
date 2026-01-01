<script>
  import { t } from "svelte-i18n";

  let {
    show = false,
    patientName = "",
    currentPhone = "",
    onClose = () => {},
    onConfirm = () => {},
  } = $props();

  let editedPhone = $state("");
  let isValid = $derived(editedPhone.length >= 10);

  // Sync editedPhone with currentPhone when modal opens
  $effect(() => {
    if (show && currentPhone) {
      editedPhone = currentPhone;
    }
  });

  function handleConfirm() {
    if (isValid) {
      onConfirm(editedPhone);
    }
  }

  function handleKeydown(e) {
    if (e.key === "Enter" && isValid) {
      handleConfirm();
    } else if (e.key === "Escape") {
      onClose();
    }
  }
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
      onclick={onClose}
      onkeydown={(e) => {
        if (e.key === "Escape") onClose();
      }}
      role="button"
      tabindex="0"
      aria-label={$t("common.close")}
    ></div>
    <div
      class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-md p-4 sm:p-6"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => {
        if (e.key === "Escape") onClose();
      }}
      role="dialog"
      aria-modal="true"
      aria-labelledby="phone-edit-title"
      tabindex="-1"
    >
      <div class="mb-4">
        <h3
          id="phone-edit-title"
          class="text-lg font-semibold text-slate-900 mb-2"
        >
          {$t("reminder.retry.edit_phone_title")}
        </h3>
        <p class="text-sm text-slate-600 mb-4">
          {$t("reminder.retry.edit_phone_message", {
            values: { name: patientName },
          })}
        </p>

        <div class="space-y-2">
          <label
            for="phone-input"
            class="block text-sm font-medium text-slate-700"
          >
            {$t("patient.phone")}
          </label>
          <input
            id="phone-input"
            type="tel"
            bind:value={editedPhone}
            onkeydown={handleKeydown}
            placeholder="+62812..."
            class="w-full px-4 py-2.5 border border-slate-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent"
          />
          {#if !isValid && editedPhone.length > 0}
            <p class="text-xs text-red-600">
              {$t("reminder.retry.phone_too_short")}
            </p>
          {/if}
        </div>
      </div>

      <div class="flex flex-col sm:flex-row gap-3">
        <button
          onclick={onClose}
          class="flex-1 h-10 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 flex items-center justify-center"
        >
          {$t("common.cancel")}
        </button>
        <button
          onclick={handleConfirm}
          disabled={!isValid}
          class="flex-1 h-10 px-4 py-2.5 bg-emerald-600 text-white font-medium rounded-xl hover:bg-emerald-700 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
        >
          {$t("reminder.retry.update_and_retry")}
        </button>
      </div>
    </div>
  </div>
{/if}
