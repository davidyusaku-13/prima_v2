<script>
  import { t } from 'svelte-i18n';
  import ReminderHistoryView from './ReminderHistoryView.svelte';

  /**
   * @typedef {Object} Props
   * @property {Object} patient - Patient object
   * @property {string} token - Auth token
   * @property {Function} onClose - Callback to close detail view
   * @property {Function} [onEditPatient] - Callback to edit patient
   * @property {Function} [onCreateReminder] - Callback to create reminder
   */

  /** @type {Props} */
  let {
    patient,
    token,
    onClose = () => {},
    onEditPatient = () => {},
    onCreateReminder = () => {}
  } = $props();

  // Active tab state
  let activeTab = $state('info'); // 'info' | 'reminders' | 'history'

  // Tabs configuration
  const tabs = [
    { id: 'info', label: 'patients.info', icon: 'info' },
    { id: 'reminders', label: 'patients.reminders', icon: 'reminder' },
    { id: 'history', label: 'patients.reminderHistory', icon: 'history' }
  ];

  // Format phone number for WhatsApp
  function formatPhoneForWhatsApp(phone) {
    if (!phone) return '';
    // Remove any non-numeric characters except +
    let cleaned = phone.replace(/[^\d+]/g, '');
    // If it starts with 0, replace with 62 (Indonesia country code)
    if (cleaned.startsWith('0')) {
      cleaned = '62' + cleaned.slice(1);
    }
    return cleaned;
  }

  function formatDate(dateStr) {
    if (!dateStr) return '-';
    try {
      return new Date(dateStr).toLocaleDateString('id-ID', {
        day: 'numeric',
        month: 'long',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch {
      return dateStr;
    }
  }

  function openWhatsApp() {
    const phone = formatPhoneForWhatsApp(patient.phone);
    if (phone) {
      window.open(`https://wa.me/${phone}`, '_blank');
    }
  }
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
  <div class="bg-white rounded-2xl shadow-2xl w-full max-w-4xl max-h-[90vh] overflow-hidden flex flex-col">
    <!-- Header -->
    <div class="flex items-center justify-between px-6 py-4 border-b border-slate-200 bg-slate-50">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-lg">
          {patient.name?.charAt(0).toUpperCase() || '?'}
        </div>
        <div>
          <h2 class="text-xl font-semibold text-slate-900">{patient.name}</h2>
          <p class="text-sm text-slate-500">{patient.phone || '-'}</p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <button
          onclick={openWhatsApp}
          class="p-2 text-green-600 hover:bg-green-50 rounded-lg transition-colors"
          title="Chat WhatsApp"
        >
          <svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
            <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413z"/>
          </svg>
        </button>
        <button
          onclick={() => onEditPatient(patient)}
          class="p-2 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
          title={$t('common.edit')}
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
          </svg>
        </button>
        <button
          onclick={onClose}
          class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-lg transition-colors"
          title={$t('common.close')}
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex border-b border-slate-200 bg-white px-4">
      {#each tabs as tab}
        <button
          onclick={() => activeTab = tab.id}
          class="px-4 py-3 text-sm font-medium border-b-2 transition-colors flex items-center gap-2 {
            activeTab === tab.id
              ? 'text-teal-600 border-teal-600'
              : 'text-slate-500 border-transparent hover:text-slate-700'
          }"
        >
          {#if tab.icon === 'info'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          {:else if tab.icon === 'reminder'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          {:else if tab.icon === 'history'}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          {/if}
          {$t(tab.label)}
        </button>
      {/each}
    </div>

    <!-- Tab Content -->
    <div class="flex-1 overflow-y-auto p-6">
      {#if activeTab === 'info'}
        <!-- Info Tab -->
        <div class="space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1">
                {$t('patients.patientName')}
              </label>
              <p class="text-slate-900">{patient.name || '-'}</p>
            </div>
            <div>
              <label class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1">
                {$t('patients.phone')}
              </label>
              <p class="text-slate-900">{patient.phone || '-'}</p>
            </div>
            <div>
              <label class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1">
                {$t('patients.email')}
              </label>
              <p class="text-slate-900">{patient.email || '-'}</p>
            </div>
            <div>
              <label class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1">
                {$t('patients.notes')}
              </label>
              <p class="text-slate-900">{patient.notes || '-'}</p>
            </div>
            <div>
              <label class="block text-xs font-medium text-slate-500 uppercase tracking-wide mb-1">
                Dibuat
              </label>
              <p class="text-slate-900">{formatDate(patient.createdAt)}</p>
            </div>
          </div>
        </div>

      {:else if activeTab === 'reminders'}
        <!-- Active Reminders Tab -->
        <div class="text-center py-8">
          <p class="text-slate-500 mb-4">{$t('patients.noReminderHistory')}</p>
          <button
            onclick={onCreateReminder}
            class="px-4 py-2 bg-teal-600 text-white rounded-lg hover:bg-teal-700 transition-colors"
          >
            {$t('patients.createFirstReminderButton')}
          </button>
        </div>

      {:else if activeTab === 'history'}
        <!-- Reminder History Tab -->
        <ReminderHistoryView
          {token}
          patientId={patient.id}
          {patient}
          {onCreateReminder}
        />
      {/if}
    </div>
  </div>
</div>
