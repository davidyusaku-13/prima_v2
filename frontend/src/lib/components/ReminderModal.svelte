<script>
  import { t } from "svelte-i18n";
  import WhatsAppPreview from "$lib/components/whatsapp/WhatsAppPreview.svelte";
  import ContentPickerModal from "$lib/components/content/ContentPickerModal.svelte";
  import ContentChip from "$lib/components/content/ContentChip.svelte";

  let {
    show = false,
    editingReminder = null,
    patient = null,
    userRole = "volunteer", // 'superadmin', 'admin', or 'volunteer'
    reminderForm = {
      patientId: "",
      title: "",
      description: "",
      dueDate: "",
      priority: "medium",
      recurrence: {
        frequency: "none",
        interval: 1,
        daysOfWeek: [],
        endDate: "",
      },
      attachments: [],
    },
    onClose = () => {},
    onSave = () => {},
  } = $props();

  // State for content picker modal
  let showContentPicker = $state(false);

  // Handle content selection from picker
  function handleContentSelect(selectedContent) {
    reminderForm.attachments = selectedContent;
  }

  // Handle content removal from chip
  function handleContentRemove(attachmentToRemove) {
    reminderForm.attachments = reminderForm.attachments.filter(
      (a) =>
        !(a.id === attachmentToRemove.id && a.type === attachmentToRemove.type)
    );
  }

  // Get attachment count text
  let attachmentCountText = $derived(
    reminderForm.attachments.length > 0
      ? ` (${reminderForm.attachments.length})`
      : ""
  );

  const daysOfWeek = [
    { value: 0, label: "Sun" },
    { value: 1, label: "Mon" },
    { value: 2, label: "Tue" },
    { value: 3, label: "Wed" },
    { value: 4, label: "Thu" },
    { value: 5, label: "Fri" },
    { value: 6, label: "Sat" },
  ];

  function handleSubmit(e) {
    e.preventDefault();
    onSave();
  }
</script>

{#if show}
  <div
    class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4"
  >
    <div
      class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
      onclick={onClose}
      onkeydown={(e) => e.key === "Escape" && onClose()}
      role="button"
      tabindex="0"
      aria-label="Close modal"
    ></div>
    <div
      class="relative bg-white rounded-t-2xl sm:rounded-2xl shadow-xl w-full sm:max-w-md p-4 sm:p-6 max-h-[90vh] overflow-y-auto"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <!-- Drag handle for mobile -->
      <div class="sm:hidden flex justify-center mb-3">
        <div class="w-10 h-1 bg-slate-200 rounded-full"></div>
      </div>
      <h2
        class="text-base sm:text-lg md:text-xl font-semibold text-slate-900 mb-4 sm:mb-6"
      >
        {editingReminder
          ? $t("patients.editReminder")
          : $t("patients.addReminder")}
      </h2>
      <form onsubmit={handleSubmit} class="space-y-3 sm:space-y-4">
        <div>
          <label
            for="title"
            class="block text-xs sm:text-sm font-medium text-slate-700 mb-1"
          >
            {$t("reminders.title")} *
          </label>
          <input
            id="title"
            type="text"
            bind:value={reminderForm.title}
            required
            class="w-full px-3 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            placeholder={$t("reminders.titlePlaceholder")}
          />
        </div>
        <div>
          <label
            for="description"
            class="block text-xs sm:text-sm font-medium text-slate-700 mb-1"
          >
            {$t("reminders.description")}
          </label>
          <textarea
            id="description"
            bind:value={reminderForm.description}
            rows="2"
            class="w-full px-3 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200 resize-none"
            placeholder={$t("reminders.descriptionPlaceholder")}
          ></textarea>
        </div>

        <!-- Content Attachment Section -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <span class="block text-xs sm:text-sm font-medium text-slate-700">
              {$t("content.picker.title")}
            </span>
            <button
              type="button"
              onclick={() => (showContentPicker = true)}
              class="text-xs text-teal-600 hover:text-teal-700 font-medium underline"
              aria-label="{$t('content.picker.attach')}{attachmentCountText}"
            >
              {$t("content.picker.attach")}{attachmentCountText}
            </button>
          </div>

          <!-- Attached content chips -->
          {#if reminderForm.attachments.length > 0}
            <div class="flex flex-wrap gap-2">
              {#each reminderForm.attachments as attachment (attachment.id + "-" + attachment.type)}
                <ContentChip {attachment} onRemove={handleContentRemove} />
              {/each}
            </div>
          {:else}
            <p class="text-xs text-slate-500">
              {$t("content.picker.empty")}
            </p>
          {/if}
        </div>

        <div class="grid grid-cols-2 gap-3">
          <div>
            <label
              for="dueDate"
              class="block text-xs sm:text-sm font-medium text-slate-700 mb-1"
            >
              {$t("reminders.dueDate")}
            </label>
            <input
              id="dueDate"
              type="datetime-local"
              bind:value={reminderForm.dueDate}
              class="w-full px-2 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-xs sm:text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            />
          </div>
          <div>
            <label
              for="priority"
              class="block text-xs sm:text-sm font-medium text-slate-700 mb-1"
            >
              {$t("reminders.priority")}
            </label>
            <select
              id="priority"
              bind:value={reminderForm.priority}
              class="w-full px-2 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-xs sm:text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            >
              <option value="low">{$t("reminders.low")}</option>
              <option value="medium">{$t("reminders.medium")}</option>
              <option value="high">{$t("reminders.high")}</option>
            </select>
          </div>
        </div>

        <!-- Recurrence Section -->
        <div class="pt-3 sm:pt-4 border-t border-slate-100">
          <h3
            class="text-xs sm:text-sm font-medium text-slate-700 mb-2 sm:mb-3"
          >
            {$t("reminders.recurrence")}
          </h3>
          <div class="space-y-2 sm:space-y-3">
            <div>
              <label
                for="frequency"
                class="block text-[10px] sm:text-xs text-slate-500 mb-1"
              >
                {$t("reminders.repeat")}
              </label>
              <select
                id="frequency"
                bind:value={reminderForm.recurrence.frequency}
                class="w-full px-3 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-xs sm:text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
              >
                <option value="none">{$t("reminders.doesNotRepeat")}</option>
                <option value="daily">{$t("reminders.daily")}</option>
                <option value="weekly">{$t("reminders.weekly")}</option>
                <option value="monthly">{$t("reminders.monthly")}</option>
                <option value="yearly">{$t("reminders.yearly")}</option>
              </select>
            </div>

            {#if reminderForm.recurrence.frequency !== "none"}
              <div>
                <label
                  for="interval"
                  class="block text-[10px] sm:text-xs text-slate-500 mb-1"
                >
                  {$t("reminders.repeatEvery")}
                </label>
                <input
                  id="interval"
                  type="number"
                  min="1"
                  max="99"
                  bind:value={reminderForm.recurrence.interval}
                  class="w-full px-3 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-xs sm:text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                />
              </div>

              {#if reminderForm.recurrence.frequency === "weekly"}
                <div>
                  <span
                    class="block text-[10px] sm:text-xs text-slate-500 mb-1.5 sm:mb-2"
                    >{$t("reminders.daysOfWeek")}</span
                  >
                  <div
                    class="flex gap-1"
                    role="group"
                    aria-label={$t("reminders.daysOfWeek")}
                  >
                    {#each daysOfWeek as day}
                      <button
                        type="button"
                        onclick={() => onToggleDay(day.value)}
                        class="w-7 h-7 sm:w-8 sm:h-8 text-[10px] sm:text-xs font-medium rounded-md sm:rounded-lg transition-colors duration-200 {reminderForm.recurrence.daysOfWeek.includes(
                          day.value
                        )
                          ? 'bg-teal-600 text-white'
                          : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
                      >
                        {day.label}
                      </button>
                    {/each}
                  </div>
                </div>
              {/if}

              <div>
                <label
                  for="endDate"
                  class="block text-[10px] sm:text-xs text-slate-500 mb-1"
                >
                  {$t("reminders.endDate")}
                </label>
                <input
                  id="endDate"
                  type="date"
                  bind:value={reminderForm.recurrence.endDate}
                  class="w-full px-3 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-xs sm:text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                />
              </div>
            {/if}
          </div>
        </div>

        <!-- WhatsApp Preview Section -->
        {#if reminderForm.title || reminderForm.description || reminderForm.attachments.length > 0}
          <div class="pt-3 sm:pt-4 border-t border-slate-100">
            <WhatsAppPreview
              message={reminderForm.description}
              patientName={patient?.name || ""}
              reminderTitle={reminderForm.title}
              attachments={reminderForm.attachments}
            />
          </div>
        {/if}

        <div
          class="flex flex-col-reverse sm:flex-row gap-2 sm:gap-3 pt-3 sm:pt-4"
        >
          <button
            type="button"
            onclick={onClose}
            class="flex-1 h-10 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-lg sm:rounded-xl hover:bg-slate-50 active:bg-slate-100 transition-colors duration-200 text-sm flex items-center justify-center"
          >
            {$t("common.cancel")}
          </button>
          <button
            type="submit"
            class="flex-1 h-10 px-4 py-2.5 bg-teal-600 text-white font-medium rounded-lg sm:rounded-xl hover:bg-teal-700 active:bg-teal-800 transition-colors duration-200 text-sm flex items-center justify-center"
          >
            {editingReminder ? $t("common.save") : $t("patients.addReminder")}
          </button>
        </div>
      </form>
    </div>
  </div>

  <!-- Content Picker Modal -->
  <ContentPickerModal
    show={showContentPicker}
    selectedContent={reminderForm.attachments}
    {userRole}
    onClose={() => (showContentPicker = false)}
    onSelect={handleContentSelect}
  />
{/if}
