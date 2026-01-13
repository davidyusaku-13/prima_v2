<script>
  import { onMount } from 'svelte';
  import { locale, t } from 'svelte-i18n';

  // Components
  import LoginScreen from '$lib/views/LoginScreen.svelte';
  import DashboardView from '$lib/views/DashboardView.svelte';
  import PatientsView from '$lib/views/PatientsView.svelte';
  import UsersView from '$lib/views/UsersView.svelte';
  import BeritaView from '$lib/views/BeritaView.svelte';
  import BeritaDetailView from '$lib/views/BeritaDetailView.svelte';
  import VideoEdukasiView from '$lib/views/VideoEdukasiView.svelte';
  import CMSDashboardView from '$lib/views/CMSDashboardView.svelte';
  import ArticleEditorView from '$lib/views/ArticleEditorView.svelte';
  import VideoManagerView from '$lib/views/VideoManagerView.svelte';
  import FailedDeliveriesView from '$lib/views/analytics/FailedDeliveriesView.svelte';
  import CmsAnalyticsView from '$lib/views/cms/CmsAnalyticsView.svelte';
  import Sidebar from '$lib/components/Sidebar.svelte';
  import BottomNav from '$lib/components/BottomNav.svelte';
  import PatientModal from '$lib/components/PatientModal.svelte';
  import ReminderModal from '$lib/components/ReminderModal.svelte';
  import UserModal from '$lib/components/UserModal.svelte';
  import ConfirmModal from '$lib/components/ConfirmModal.svelte';
  import ProfileModal from '$lib/components/ProfileModal.svelte';
  import VideoModal from '$lib/components/VideoModal.svelte';
  import VideoEditModal from '$lib/components/VideoEditModal.svelte';
  import SendReminderModal from '$lib/components/SendReminderModal.svelte';
  import PhoneEditModal from '$lib/components/PhoneEditModal.svelte';
  import FailedReminderBadge from '$lib/components/indicators/FailedReminderBadge.svelte';
  import Toast from '$lib/components/ui/Toast.svelte';

  // API
  import * as api from '$lib/utils/api.js';

  // Stores
  import { toastStore } from '$lib/stores/toast.svelte.js';
  import { deliveryStore } from '$lib/stores/delivery.svelte.js';

  // Auth state
  let token = $state(localStorage.getItem('token') || null);
  let user = $state(null);
  let authLoading = $state(true);

  // State
  let patients = $state([]);
  let users = $state([]);
  let loading = $state(true);
  let currentView = $state(localStorage.getItem('currentView') || 'dashboard');

  // Auth forms
  let showAuthModal = $state(false);
  let showProfileModal = $state(false);
  let authMode = $state('login');
  let authError = $state('');
  let loginForm = $state({ username: '', password: '' });
  let registerForm = $state({ username: '', password: '', confirmPassword: '', fullName: '' });

  // Password validation
  function getPasswordStrength(password) {
    let score = 0;
    if (password.length >= 6) score++;
    if (password.length >= 10) score++;
    if (/[a-z]/.test(password) && /[A-Z]/.test(password)) score++;
    if (/\d/.test(password)) score++;
    if (/[^a-zA-Z0-9]/.test(password)) score++;
    return score;
  }

  let passwordStrength = $derived(getPasswordStrength(registerForm.password));
  let passwordMatch = $derived(registerForm.password && registerForm.confirmPassword && registerForm.password === registerForm.confirmPassword);
  let usernameValid = $derived(registerForm.username.length >= 3);
  let formValid = $derived(usernameValid && registerForm.password.length >= 6 && passwordMatch);

  // Form states
  let showPatientModal = $state(false);
  let showReminderModal = $state(false);
  let editingPatient = $state(null);
  let editingReminder = $state(null);
  let patientForm = $state({ name: '', phone: '', email: '', notes: '' });

  let showUserModal = $state(false);
  let editingUser = $state(null);
  let userForm = $state({ role: 'volunteer', username: '', password: '' });
  let userFormLoading = $state(false);

  let reminderForm = $state({
    patientId: '',
    title: '',
    description: '',
    dueDate: '',
    priority: 'medium',
    recurrence: { frequency: 'none', interval: 1, daysOfWeek: [], endDate: '' },
    attachments: []
  });

  // Confirm modal state
  let showConfirmModal = $state(false);
  let confirmMessage = $state('');
  let confirmCallback = $state(null);
  let confirmContext = $state(null);

  // Filters
  let searchQuery = $state('');

  // CMS State
  let showArticleEditor = $state(false);
  let editingArticle = $state(null);
  let showVideoManager = $state(false);
  let showVideoModal = $state(false);
  let showVideoEditModal = $state(false);
  let editingVideo = $state(null);
  let currentVideo = $state(null);
  let currentArticleId = $state(localStorage.getItem('currentArticleId') || null);

  // Send Reminder Modal State
  let showSendReminderModal = $state(false);
  let sendReminderPatient = $state(null);
  let sendReminderReminder = $state(null);
  let sendReminderStatus = $state('idle'); // 'idle' | 'sending' | 'success' | 'error' | 'scheduled'
  let sendReminderError = $state('');
  let quietHoursConfig = $state(null);
  let isCurrentlyQuietHours = $state(false);
  let scheduledDeliveryTime = $state(null);

  // Phone Edit Modal State
  let showPhoneEditModal = $state(false);
  let phoneEditPatient = $state(null);
  let phoneEditReminder = $state(null);

  // CMS Functions
  function openArticleEditor(article = null) {
    editingArticle = article;
    showArticleEditor = true;
  }

  function closeArticleEditor() {
    showArticleEditor = false;
    editingArticle = null;
  }

  function handleArticleSave(status) {
    console.log('Article saved:', status);
  }

  function openVideoManager() {
    showVideoManager = true;
  }

  function closeVideoManager() {
    showVideoManager = false;
  }

  function handleVideoSave() {
    console.log('Video saved');
  }

  function openVideoEditModal(video) {
    editingVideo = video;
    showVideoEditModal = true;
  }

  function closeVideoEditModal() {
    showVideoEditModal = false;
    editingVideo = null;
  }

  function watchVideo(video) {
    currentVideo = video;
    showVideoModal = true;
  }

  function closeVideoModal() {
    showVideoModal = false;
    currentVideo = null;
  }

  function viewArticle(article) {
    currentArticleId = article.slug;
    localStorage.setItem('currentArticleId', article.slug);
    navigateTo('berita-detail');
  }

  // Locale
  function setLocale(newLocale) {
    locale.set(newLocale);
    localStorage.setItem('locale', newLocale);
  }

  onMount(async () => {
    const savedLocale = localStorage.getItem('locale');
    if (savedLocale) {
      locale.set(savedLocale);
    }

    if (token) {
      await fetchUser();
    } else {
      authLoading = false;
    }
  });

  // Stats
  let stats = $derived({
    totalPatients: patients.length,
    totalReminders: patients.reduce((acc, p) => acc + (p.reminders?.length || 0), 0),
    completedReminders: patients.reduce((acc, p) => acc + (p.reminders?.filter(r => r.completed).length || 0), 0),
    pendingReminders: patients.reduce((acc, p) => acc + (p.reminders?.filter(r => !r.completed).length || 0), 0)
  });

  // Auth functions
  async function fetchUser() {
    try {
      const userData = await api.fetchUser(token);
      user = userData;
      await loadPatients();
      await loadUsers();
    } catch (e) {
      logout();
    } finally {
      authLoading = false;
    }
  }

  async function login() {
    authError = '';
    authLoading = true;
    try {
      const data = await api.login(loginForm.username, loginForm.password);
      token = data.token;
      user = { userId: data.userId, username: data.username, fullName: data.fullName, role: data.role };
      localStorage.setItem('token', token);
      loginForm = { username: '', password: '' };
      await loadPatients();
    } catch (e) {
      authError = e.message || 'Login failed';
    } finally {
      authLoading = false;
    }
  }

  async function register() {
    authError = '';
    authLoading = true;
    try {
      const data = await api.register(registerForm.username, registerForm.password, registerForm.fullName);
      token = data.token;
      user = { userId: data.userId, username: data.username, fullName: data.fullName, role: data.role };
      localStorage.setItem('token', token);
      registerForm = { username: '', password: '', confirmPassword: '', fullName: '' };
      await loadPatients();
    } catch (e) {
      authError = e.message || 'Registration failed';
    } finally {
      authLoading = false;
    }
  }

  function logout() {
    token = null;
    user = null;
    patients = [];
    localStorage.removeItem('token');
    authLoading = false;
  }

  // Data functions
  async function loadPatients() {
    loading = true;
    try {
      patients = await api.fetchPatients(token);
    } catch (e) {
      console.error('Failed to load patients:', e);
      if (e.message === 'Unauthorized') {
        logout();
      }
    } finally {
      loading = false;
    }
  }

  async function loadUsers() {
    if (user?.role !== 'superadmin') return;
    try {
      users = await api.fetchUsers(token);
    } catch (e) {
      console.error('Failed to load users:', e);
    }
  }

  // Patient functions
  async function savePatient() {
    try {
      await api.savePatient(token, patientForm, editingPatient?.id);
      await loadPatients();
      closePatientModal();
    } catch (e) {
      console.error('Failed to save patient:', e);
    }
  }

  function deletePatient(id) {
    showConfirm('Are you sure you want to delete this patient?', async () => {
      try {
        await api.deletePatient(token, id);
        await loadPatients();
      } catch (e) {
        console.error('Failed to delete patient:', e);
      }
    });
  }

  function openPatientModal(patient = null) {
    editingPatient = patient;
    if (patient) {
      patientForm = { name: patient.name, phone: patient.phone, email: patient.email, notes: patient.notes };
    } else {
      patientForm = { name: '', phone: '', email: '', notes: '' };
    }
    showPatientModal = true;
  }

  function closePatientModal() {
    showPatientModal = false;
    editingPatient = null;
    patientForm = { name: '', phone: '', email: '', notes: '' };
  }

  // Reminder functions
  async function saveReminder() {
    try {
      await api.saveReminder(token, reminderForm.patientId, reminderForm, editingReminder?.id);
      await loadPatients();
      closeReminderModal();
    } catch (e) {
      console.error('Failed to save reminder:', e);
    }
  }

  async function toggleReminder(patientId, reminderId) {
    try {
      await api.toggleReminder(token, patientId, reminderId);
      await loadPatients();
    } catch (e) {
      console.error('Failed to toggle reminder:', e);
    }
  }

  async function deleteReminder(patientId, reminderId) {
    try {
      await api.deleteReminder(token, patientId, reminderId);
      await loadPatients();
    } catch (e) {
      console.error('Failed to delete reminder:', e);
    }
  }

  function openReminderModal(patientId, reminder = null) {
    editingReminder = reminder;
    if (reminder) {
      reminderForm = {
        patientId,
        title: reminder.title,
        description: reminder.description || '',
        dueDate: reminder.dueDate || '',
        priority: reminder.priority || 'medium',
        recurrence: reminder.recurrence || { frequency: 'none', interval: 1, daysOfWeek: [], endDate: '' },
        attachments: reminder.attachments || []
      };
    } else {
      reminderForm = {
        patientId,
        title: '',
        description: '',
        dueDate: '',
        priority: 'medium',
        recurrence: { frequency: 'none', interval: 1, daysOfWeek: [], endDate: '' },
        attachments: []
      };
    }
    showReminderModal = true;
  }

  function closeReminderModal() {
    showReminderModal = false;
    editingReminder = null;
    reminderForm = {
      patientId: '',
      title: '',
      description: '',
      dueDate: '',
      priority: 'medium',
      recurrence: { frequency: 'none', interval: 1, daysOfWeek: [], endDate: '' },
      attachments: []
    };
  }

  function toggleDayOfWeek(day) {
    const idx = reminderForm.recurrence.daysOfWeek.indexOf(day);
    if (idx === -1) {
      reminderForm.recurrence.daysOfWeek = [...reminderForm.recurrence.daysOfWeek, day].sort();
    } else {
      reminderForm.recurrence.daysOfWeek = reminderForm.recurrence.daysOfWeek.filter(d => d !== day);
    }
  }

  // Send Reminder functions
  async function openSendReminderModal(patient, reminder) {
    sendReminderPatient = patient;
    sendReminderReminder = reminder;
    sendReminderStatus = 'idle';
    sendReminderError = '';
    scheduledDeliveryTime = null;

    // Fetch quiet hours config and check if currently in quiet hours
    try {
      quietHoursConfig = await api.fetchQuietHoursConfig();
      isCurrentlyQuietHours = api.isQuietHours(quietHoursConfig);
      if (isCurrentlyQuietHours) {
        scheduledDeliveryTime = api.getNextActiveTime(quietHoursConfig);
      }
    } catch (e) {
      console.error('Failed to fetch quiet hours config:', e);
      isCurrentlyQuietHours = false;
    }

    showSendReminderModal = true;
  }

  function closeSendReminderModal() {
    showSendReminderModal = false;
    sendReminderPatient = null;
    sendReminderReminder = null;
    sendReminderStatus = 'idle';
    sendReminderError = '';
    isCurrentlyQuietHours = false;
    scheduledDeliveryTime = null;
  }

  async function handleSendReminder() {
    if (!sendReminderPatient || !sendReminderReminder) return;

    sendReminderStatus = 'sending';
    sendReminderError = '';

    // Optimistic UI update - set reminder status to 'sending' locally
    const patientIndex = patients.findIndex(p => p.id === sendReminderPatient.id);
    if (patientIndex !== -1) {
      const reminderIndex = patients[patientIndex].reminders?.findIndex(r => r.id === sendReminderReminder.id);
      if (reminderIndex !== -1 && reminderIndex !== undefined) {
        patients[patientIndex].reminders[reminderIndex].delivery_status = 'sending';
        patients = [...patients]; // Trigger reactivity
      }
    }

    try {
      await api.sendReminder(token, sendReminderPatient.id, sendReminderReminder.id);
      sendReminderStatus = 'success';
      await loadPatients(); // Refresh to get updated delivery status
      closeSendReminderModal();
    } catch (e) {
      sendReminderStatus = 'error';
      sendReminderError = e.message || 'Failed to send reminder';
      // Revert optimistic update on error
      await loadPatients();
    }
  }

  // Retry failed reminder
  async function handleRetryReminder(patient, reminder) {
    // Optimistic UI update - set reminder status to 'sending' locally
    deliveryStore.updateStatus(reminder.id, 'sending');

    try {
      const result = await api.retryReminder(token, reminder.id);

      // Update delivery store with new status
      deliveryStore.updateStatus(reminder.id, result.status);

      // Show appropriate toast based on status
      if (result.status === 'queued') {
        toastStore.show({
          type: 'warning',
          message: $t('reminder.retry.queued_success')
        });
      } else {
        toastStore.show({
          type: 'success',
          message: $t('reminder.retry.success')
        });
      }

      // Refresh patients to get updated data
      await loadPatients();
    } catch (error) {
      // Revert to failed state
      deliveryStore.updateStatus(reminder.id, 'failed', error.message);

      // Show error toast with specific message based on error code
      console.error('Retry failed:', error);

      // Handle INVALID_PHONE error with phone edit modal
      if (error.code === 'INVALID_PHONE') {
        phoneEditPatient = patient;
        phoneEditReminder = reminder;
        showPhoneEditModal = true;
        return; // Don't show toast, modal will handle it
      }

      let toastMessage = $t('reminder.retry.error');

      if (error.code === 'CIRCUIT_BREAKER_OPEN') {
        toastMessage = $t('reminder.retry.circuit_breaker_open');
      } else if (error.code === 'INVALID_STATUS') {
        toastMessage = $t('reminder.retry.invalid_status');
      }

      toastStore.show({
        type: 'error',
        message: toastMessage
      });

      // Refresh patients to sync state
      await loadPatients();
    }
  }

  // Handle phone update and retry
  async function handlePhoneUpdateAndRetry(newPhone) {
    showPhoneEditModal = false;

    try {
      // Update patient phone number
      await api.updatePatient(token, phoneEditPatient.id, {
        ...phoneEditPatient,
        phone: newPhone
      });

      // Refresh patients to get updated data
      await loadPatients();

      // Retry the reminder with updated phone
      await handleRetryReminder(phoneEditPatient, phoneEditReminder);
    } catch (error) {
      console.error('Failed to update phone:', error);
      toastStore.show({
        type: 'error',
        message: $t('patient.updateFailed')
      });
    }
  }

  // User functions
  async function registerUser() {
    userFormLoading = true;
    try {
      await api.registerUser(token, userForm.username, userForm.password);
      closeUserModal();
      loadUsers();
    } catch (e) {
      alert(e.message || 'Failed to create user');
    } finally {
      userFormLoading = false;
    }
  }

  async function updateUserRole() {
    if (!editingUser) return;
    try {
      await api.updateUserRole(token, editingUser.id, userForm.role);
      await loadUsers();
      closeUserModal();
    } catch (e) {
      console.error('Failed to update user role:', e);
    }
  }

  function deleteUser(userId) {
    showConfirm('Are you sure you want to delete this user?', async () => {
      try {
        await api.deleteUser(token, userId);
        await loadUsers();
      } catch (e) {
        console.error('Failed to delete user:', e);
      }
    });
  }

  function openUserModal(userToEdit = null) {
    editingUser = userToEdit;
    if (userToEdit) {
      userForm = { role: userToEdit.role, username: '', password: '' };
    } else {
      userForm = { role: 'volunteer', username: '', password: '' };
    }
    showUserModal = true;
  }

  function closeUserModal() {
    showUserModal = false;
    editingUser = null;
    userForm = { role: 'volunteer', username: '', password: '' };
  }

  // Navigation
  function navigateTo(view) {
    currentView = view;
    localStorage.setItem('currentView', view);

    // Clear toasts on navigation to prevent stale notifications
    toastStore.clear();

    if (view === 'users') loadUsers();
    if (view === 'berita-detail' && !currentArticleId) {
      currentView = 'berita';
      localStorage.setItem('currentView', 'berita');
    }
  }

  // Confirm modal
  function showConfirm(message, callback, context = null) {
    confirmMessage = message;
    confirmCallback = callback;
    confirmContext = context;
    showConfirmModal = true;
  }

  function handleConfirm() {
    if (confirmCallback) {
      confirmCallback(confirmContext);
    }
    closeConfirmModal();
  }

  function closeConfirmModal() {
    showConfirmModal = false;
    confirmMessage = '';
    confirmCallback = null;
    confirmContext = null;
  }
</script>

<!-- Auth Loading Screen -->
{#if authLoading}
  <div class="min-h-screen bg-slate-50 flex items-center justify-center">
    <div class="animate-spin w-10 h-10 border-4 border-teal-600 border-t-transparent rounded-full"></div>
  </div>
<!-- Login Screen -->
{:else if !token}
  <LoginScreen
    bind:authMode
    bind:authError
    bind:authLoading
    bind:loginForm
    bind:registerForm
    {passwordStrength}
    {passwordMatch}
    {usernameValid}
    {formValid}
    {setLocale}
    onLogin={login}
    onRegister={register}
  />
<!-- Main Dashboard -->
{:else}
  <Sidebar
    {user}
    {currentView}
    {stats}
    {users}
    locale={$locale}
    onNavigate={navigateTo}
    onSetLocale={setLocale}
    onLogout={logout}
  />

  <!-- Failed Reminder Badge - Persistent indicator in navigation -->
  <div class="fixed top-4 left-4 lg:left-72 z-40">
    <FailedReminderBadge />
  </div>

  <main class="flex-1 lg:ml-64 p-3 sm:p-4 md:p-5 lg:p-6 xl:p-8 min-h-screen pb-24 lg:pb-6">
    <div class="max-w-7xl mx-auto">
    {#if loading}
      <div class="flex items-center justify-center h-64">
        <div class="animate-spin w-10 h-10 border-4 border-teal-600 border-t-transparent rounded-full"></div>
      </div>
    {:else}
      {#if currentView === 'dashboard'}
        <DashboardView {patients} onToggleReminder={toggleReminder} onViewAllPatients={() => navigateTo('patients')} />
      {:else if currentView === 'patients'}
        <PatientsView
          {patients}
          {token}
          bind:searchQuery
          onOpenPatientModal={openPatientModal}
          onOpenReminderModal={openReminderModal}
          onDeletePatient={deletePatient}
          onToggleReminder={toggleReminder}
          onDeleteReminder={deleteReminder}
          onSendReminder={openSendReminderModal}
          onRetryReminder={handleRetryReminder}
        />
      {:else if currentView === 'users' && user?.role === 'superadmin'}
        <UsersView
          {users}
          {user}
          onLoadUsers={loadUsers}
          onOpenUserModal={openUserModal}
          onDeleteUser={deleteUser}
        />
      {:else if currentView === 'cms' && (user?.role === 'superadmin' || user?.role === 'admin')}
        <CMSDashboardView
          {token}
          onNavigateToArticleEditor={() => openArticleEditor()}
          onNavigateToVideoManager={openVideoManager}
          onEditArticle={(article) => openArticleEditor(article)}
          onEditVideo={(video) => openVideoEditModal(video)}
        />
      {:else if currentView === 'berita'}
        <BeritaView onNavigateToArticle={viewArticle} />
      {:else if currentView === 'berita-detail'}
        <BeritaDetailView
          articleId={currentArticleId}
          onBack={() => { currentArticleId = null; localStorage.removeItem('currentArticleId'); navigateTo('berita'); }}
        />
      {:else if currentView === 'video'}
        <VideoEdukasiView onWatchVideo={watchVideo} />
      {:else if currentView === 'failed-deliveries' && (user?.role === 'admin' || user?.role === 'superadmin')}
        <FailedDeliveriesView token={token} />
      {:else if currentView === 'analytics' && (user?.role === 'admin' || user?.role === 'superadmin')}
        <CmsAnalyticsView token={token} />
      {/if}
    {/if}
    </div>
  </main>

  <BottomNav
    {user}
    {currentView}
    {stats}
    {users}
    onNavigate={navigateTo}
    onShowProfile={() => showProfileModal = true}
  />

  <!-- Profile Modal -->
  <ProfileModal
    show={showProfileModal}
    {user}
    locale={$locale}
    onSetLocale={setLocale}
    onLogout={logout}
    onClose={() => showProfileModal = false}
  />

  <!-- Patient Modal -->
  <PatientModal
    show={showPatientModal}
    {editingPatient}
    bind:patientForm
    onClose={closePatientModal}
    onSave={savePatient}
  />

  <!-- Reminder Modal -->
  <ReminderModal
    show={showReminderModal}
    {editingReminder}
    bind:reminderForm
    patient={patients.find(p => p.id === reminderForm.patientId) || null}
    userRole={user?.role || 'volunteer'}
    onClose={closeReminderModal}
    onSave={saveReminder}
    onToggleDay={toggleDayOfWeek}
  />

  <!-- User Modal -->
  <UserModal
    show={showUserModal}
    {editingUser}
    {userForm}
    loading={userFormLoading}
    onClose={closeUserModal}
    onSaveRole={updateUserRole}
    onRegister={registerUser}
  />

  <!-- Confirm Modal -->
  <ConfirmModal
    show={showConfirmModal}
    message={confirmMessage}
    onClose={closeConfirmModal}
    onConfirm={handleConfirm}
  />

  <!-- Article Editor Modal -->
  {#if showArticleEditor}
    <ArticleEditorView
      article={editingArticle}
      {token}
      onClose={closeArticleEditor}
      onSave={handleArticleSave}
    />
  {/if}

  <!-- Video Manager Modal -->
  {#if showVideoManager}
    <VideoManagerView
      {token}
      onClose={closeVideoManager}
      onSave={handleVideoSave}
    />
  {/if}

  <!-- Video Modal -->
  <VideoModal
    show={showVideoModal}
    video={currentVideo}
    onClose={closeVideoModal}
  />

  <!-- Video Edit Modal -->
  {#if showVideoEditModal}
    <VideoEditModal
      video={editingVideo}
      {token}
      onClose={closeVideoEditModal}
      onSave={handleVideoSave}
    />
  {/if}

  <!-- Send Reminder Modal -->
  <SendReminderModal
    show={showSendReminderModal}
    patient={sendReminderPatient}
    reminder={sendReminderReminder}
    status={sendReminderStatus}
    errorMessage={sendReminderError}
    isQuietHours={isCurrentlyQuietHours}
    scheduledTime={scheduledDeliveryTime}
    onClose={closeSendReminderModal}
    onConfirm={handleSendReminder}
  />

  <!-- Phone Edit Modal -->
  <PhoneEditModal
    show={showPhoneEditModal}
    patientName={phoneEditPatient?.name || ''}
    currentPhone={phoneEditPatient?.phone || ''}
    onClose={() => showPhoneEditModal = false}
    onConfirm={handlePhoneUpdateAndRetry}
  />

  <!-- Toast Notifications - Global -->
  <Toast />
{/if}
