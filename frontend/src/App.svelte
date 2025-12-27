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
  import Sidebar from '$lib/components/Sidebar.svelte';
  import BottomNav from '$lib/components/BottomNav.svelte';
  import PatientModal from '$lib/components/PatientModal.svelte';
  import ReminderModal from '$lib/components/ReminderModal.svelte';
  import UserModal from '$lib/components/UserModal.svelte';
  import ConfirmModal from '$lib/components/ConfirmModal.svelte';
  import ProfileModal from '$lib/components/ProfileModal.svelte';
  import VideoModal from '$lib/components/VideoModal.svelte';

  // API
  import * as api from '$lib/utils/api.js';

  // Auth state
  let token = localStorage.getItem('token') || null;
  let user = null;
  let authLoading = true;

  // State
  let patients = [];
  let users = [];
  let loading = true;
  let currentView = localStorage.getItem('currentView') || 'dashboard';

  // Auth forms
  let showAuthModal = false;
  let showProfileModal = false;
  let authMode = 'login';
  let authError = '';
  let loginForm = { username: '', password: '' };
  let registerForm = { username: '', password: '', confirmPassword: '', fullName: '' };

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

  $: passwordStrength = getPasswordStrength(registerForm.password);
  $: passwordMatch = registerForm.password && registerForm.confirmPassword && registerForm.password === registerForm.confirmPassword;
  $: usernameValid = registerForm.username.length >= 3;
  $: formValid = usernameValid && registerForm.password.length >= 6 && passwordMatch;

  // Form states
  let showPatientModal = false;
  let showReminderModal = false;
  let editingPatient = null;
  let editingReminder = null;
  let patientForm = { name: '', phone: '', email: '', notes: '' };

  let showUserModal = false;
  let editingUser = null;
  let userForm = { role: 'volunteer', username: '', password: '' };
  let userFormLoading = false;

  let reminderForm = {
    patientId: '',
    title: '',
    description: '',
    dueDate: '',
    priority: 'medium',
    recurrence: { frequency: 'none', interval: 1, daysOfWeek: [], endDate: '' }
  };

  // Confirm modal state
  let showConfirmModal = false;
  let confirmMessage = '';
  let confirmCallback = null;
  let confirmContext = null;

  // Filters
  let searchQuery = '';

  // CMS State
  let showArticleEditor = false;
  let editingArticle = null;
  let showVideoManager = false;
  let showVideoModal = false;
  let currentVideo = null;
  let currentArticleId = null;

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

  function watchVideo(video) {
    currentVideo = video;
    showVideoModal = true;
  }

  function closeVideoModal() {
    showVideoModal = false;
    currentVideo = null;
  }

  function viewArticle(article) {
    currentArticleId = article.id;
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
  $: stats = {
    totalPatients: patients.length,
    totalReminders: patients.reduce((acc, p) => acc + (p.reminders?.length || 0), 0),
    completedReminders: patients.reduce((acc, p) => acc + (p.reminders?.filter(r => r.completed).length || 0), 0),
    pendingReminders: patients.reduce((acc, p) => acc + (p.reminders?.filter(r => !r.completed).length || 0), 0)
  };

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
        recurrence: reminder.recurrence || { frequency: 'none', interval: 1, daysOfWeek: [], endDate: '' }
      };
    } else {
      reminderForm = {
        patientId,
        title: '',
        description: '',
        dueDate: '',
        priority: 'medium',
        recurrence: { frequency: 'none', interval: 1, daysOfWeek: [], endDate: '' }
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
      recurrence: { frequency: 'none', interval: 1, daysOfWeek: [], endDate: '' }
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
    if (view === 'users') loadUsers();
    if (view === 'berita-detail' && !currentArticleId) {
      currentView = 'berita';
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

  <main class="flex-1 lg:ml-64 p-4 lg:p-6 min-h-screen pb-20 lg:pb-6">
    {#if loading}
      <div class="flex items-center justify-center h-64">
        <div class="animate-spin w-10 h-10 border-4 border-teal-600 border-t-transparent rounded-full"></div>
      </div>
    {:else}
      {#if currentView === 'dashboard'}
        <DashboardView {patients} onToggleReminder={toggleReminder} />
      {:else if currentView === 'patients'}
        <PatientsView
          {patients}
          bind:searchQuery
          onOpenPatientModal={openPatientModal}
          onOpenReminderModal={openReminderModal}
          onDeletePatient={deletePatient}
          onToggleReminder={toggleReminder}
          onDeleteReminder={deleteReminder}
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
          onNavigateToArticleEditor={() => openArticleEditor()}
          onNavigateToVideoManager={openVideoManager}
          onNavigateToArticle={viewArticle}
        />
      {:else if currentView === 'berita'}
        <BeritaView onNavigateToArticle={viewArticle} />
      {:else if currentView === 'berita-detail'}
        <BeritaDetailView
          articleId={currentArticleId}
          onBack={() => { currentArticleId = null; currentView = 'berita'; }}
        />
      {:else if currentView === 'video'}
        <VideoEdukasiView onWatchVideo={watchVideo} />
      {/if}
    {/if}
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
    {patientForm}
    onClose={closePatientModal}
    onSave={savePatient}
  />

  <!-- Reminder Modal -->
  <ReminderModal
    show={showReminderModal}
    {editingReminder}
    {reminderForm}
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
{/if}
