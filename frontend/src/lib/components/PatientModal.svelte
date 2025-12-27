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
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div
      class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
      onclick={onClose}
      onkeydown={(e) => e.key === 'Escape' && onClose()}
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
      <h2 class="text-lg sm:text-xl font-semibold text-slate-900 mb-4 sm:mb-6">
        {editingPatient ? $t('patients.editPatient') : $t('patients.addPatient')}
      </h2>
      <form onsubmit={handleSubmit} class="space-y-4">
        <div>
          <label for="name" class="block text-sm font-medium text-slate-700 mb-1">
            {$t('patients.patientName')} *
          </label>
          <input
            id="name"
            type="text"
            bind:value={patientForm.name}
            required
            class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            placeholder={$t('patients.patientName')}
          />
        </div>
        <div>
          <label for="phone" class="block text-sm font-medium text-slate-700 mb-1">
            {$t('patients.whatsappNumber')} *
          </label>
          <input
            id="phone"
            type="tel"
            bind:value={patientForm.phone}
            required
            class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            placeholder="6281234567890"
          />
          <p class="text-xs text-slate-500 mt-1">{$t('patients.whatsappNote')}</p>
        </div>
        <div>
          <label for="email" class="block text-sm font-medium text-slate-700 mb-1">
            {$t('patients.email')}
          </label>
          <input
            id="email"
            type="email"
            bind:value={patientForm.email}
            class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            placeholder={$t('patients.emailPlaceholder')}
          />
        </div>
        <div>
          <label for="notes" class="block text-sm font-medium text-slate-700 mb-1">
            {$t('patients.notes')}
          </label>
          <textarea
            id="notes"
            bind:value={patientForm.notes}
            rows="3"
            class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200 resize-none"
            placeholder={$t('patients.notesPlaceholder')}
          ></textarea>
        </div>
        <div class="flex flex-col sm:flex-row gap-3 pt-4">
          <button
            type="button"
            onclick={onClose}
            class="flex-1 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1"
          >
            {$t('common.cancel')}
          </button>
          <button
            type="submit"
            class="flex-1 px-4 py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 transition-colors duration-200 order-1 sm:order-2"
          >
            {editingPatient ? $t('common.save') : $t('patients.addPatient')}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
