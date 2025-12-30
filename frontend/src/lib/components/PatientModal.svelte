<script>
  import { t } from 'svelte-i18n';

  export let show = false;
  export let editingPatient = null;
  export let patientForm = { name: '', phone: '', email: '', notes: '' };
  export let onClose = () => {};
  export let onSave = () => {};

  function handleSubmit(e) {
    e.preventDefault();
    onSave();
  }
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-end sm:items-center justify-center p-0 sm:p-4">
    <div
      class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
      onclick={onClose}
      onkeydown={(e) => e.key === 'Escape' && onClose()}
      role="button"
      tabindex="0"
      aria-label="Close modal"
    ></div>
    <div
      class="relative bg-white rounded-t-2xl sm:rounded-2xl shadow-xl w-full sm:max-w-md max-h-[90vh] overflow-y-auto p-4 sm:p-6"
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
      <h2 class="text-base sm:text-lg md:text-xl font-semibold text-slate-900 mb-4 sm:mb-6">
        {editingPatient ? $t('patients.editPatient') : $t('patients.addPatient')}
      </h2>
      <form onsubmit={handleSubmit} class="space-y-3 sm:space-y-4">
        <div>
          <label for="name" class="block text-xs sm:text-sm font-medium text-slate-700 mb-1">
            {$t('patients.patientName')} *
          </label>
          <input
            id="name"
            type="text"
            bind:value={patientForm.name}
            required
            class="w-full px-3 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            placeholder={$t('patients.patientName')}
          />
        </div>
        <div>
          <label for="phone" class="block text-xs sm:text-sm font-medium text-slate-700 mb-1">
            {$t('patients.whatsappNumber')} *
          </label>
          <input
            id="phone"
            type="tel"
            bind:value={patientForm.phone}
            required
            class="w-full px-3 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            placeholder="6281234567890"
          />
          <p class="text-[10px] sm:text-xs text-slate-500 mt-1">{$t('patients.whatsappNote')}</p>
        </div>
        <div>
          <label for="email" class="block text-xs sm:text-sm font-medium text-slate-700 mb-1">
            {$t('patients.email')}
          </label>
          <input
            id="email"
            type="email"
            bind:value={patientForm.email}
            class="w-full px-3 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            placeholder={$t('patients.emailPlaceholder')}
          />
        </div>
        <div>
          <label for="notes" class="block text-xs sm:text-sm font-medium text-slate-700 mb-1">
            {$t('patients.notes')}
          </label>
          <textarea
            id="notes"
            bind:value={patientForm.notes}
            rows="3"
            class="w-full px-3 sm:px-4 py-2 sm:py-2.5 bg-slate-100 border-0 rounded-lg sm:rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200 resize-none"
            placeholder={$t('patients.notesPlaceholder')}
          ></textarea>
        </div>
        <div class="flex flex-col-reverse sm:flex-row gap-2 sm:gap-3 pt-3 sm:pt-4">
          <button
            type="button"
            onclick={onClose}
            class="flex-1 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-lg sm:rounded-xl hover:bg-slate-50 active:bg-slate-100 transition-colors duration-200 text-sm"
          >
            {$t('common.cancel')}
          </button>
          <button
            type="submit"
            class="flex-1 px-4 py-2.5 bg-teal-600 text-white font-medium rounded-lg sm:rounded-xl hover:bg-teal-700 active:bg-teal-800 transition-colors duration-200 text-sm"
          >
            {editingPatient ? $t('common.save') : $t('patients.addPatient')}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
