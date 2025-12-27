<script>
  import { onMount } from 'svelte';
  import { locale, t } from 'svelte-i18n';

  const API_URL = 'http://localhost:8080/api';

  function setLocale(newLocale) {
    locale.set(newLocale);
    localStorage.setItem('locale', newLocale);
  }

  // Auth state
  let token = localStorage.getItem('token') || null;
  let user = null;
  let authLoading = true;

  // State
  let patients = [];
  let users = [];
  let loading = true;
  let sidebarOpen = false;
  let currentView = localStorage.getItem('currentView') || 'dashboard';

  // Auth forms
  let showAuthModal = false;
  let showProfileModal = false;
  let authMode = 'login'; // 'login' or 'register'
  let authError = '';

  // Login form
  let loginForm = {
    username: '',
    password: ''
  };

  // Register form
  let registerForm = {
    username: '',
    password: '',
    confirmPassword: '',
    fullName: ''
  };

  // Password strength
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

  // Real-time validation
  $: usernameValid = registerForm.username.length >= 3;
  $: usernameTaken = false; // Would need API check
  $: passwordValid = registerForm.password.length >= 6;
  $: formValid = usernameValid && passwordValid && passwordMatch;

  // Form states
  let showPatientModal = false;
  let showReminderModal = false;
  let editingPatient = null;
  let editingReminder = null;

  // Patient form
  let patientForm = {
    name: '',
    phone: '',
    email: '',
    notes: ''
  };

  // User management
  let showUserModal = false;
  let editingUser = null;
  let userForm = {
    role: 'volunteer',
    username: '',
    password: ''
  };
  let userFormLoading = false;

  // Reminder form
  let reminderForm = {
    patientId: '',
    title: '',
    description: '',
    dueDate: '',
    priority: 'medium',
    recurrence: {
      frequency: 'none',
      interval: 1,
      daysOfWeek: [],
      endDate: ''
    }
  };

  // Confirm modal state
  let showConfirmModal = false;
  let confirmMessage = '';
  let confirmCallback = null;
  let confirmContext = null;

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

  // Filters
  let searchQuery = '';

  // Days of week options
  const daysOfWeek = [
    { value: 0, label: 'Sun' },
    { value: 1, label: 'Mon' },
    { value: 2, label: 'Tue' },
    { value: 3, label: 'Wed' },
    { value: 4, label: 'Thu' },
    { value: 5, label: 'Fri' },
    { value: 6, label: 'Sat' }
  ];

  // Headers with auth
  function getHeaders() {
    const headers = { 'Content-Type': 'application/json' };
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    return headers;
  }

  onMount(async () => {
    if (token) {
      await fetchUser();
    } else {
      authLoading = false;
    }
  });

  async function fetchUser() {
    try {
      const res = await fetch(`${API_URL}/auth/me`, {
        headers: getHeaders()
      });
      if (res.ok) {
        const userData = await res.json();
        user = userData;
        await loadPatients();
        await loadUsers(); // Load user count on page load
      } else {
        logout();
      }
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
      const res = await fetch(`${API_URL}/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(loginForm)
      });
      const data = await res.json();
      if (res.ok) {
        token = data.token;
        user = { userId: data.userId, username: data.username, fullName: data.fullName, role: data.role };
        localStorage.setItem('token', token);
        showAuthModal = false;
        loginForm = { username: '', password: '' };
        await loadPatients();
      } else {
        authError = data.error || 'Login failed';
      }
    } catch (e) {
      authError = 'Connection error. Please check your internet connection.';
    } finally {
      authLoading = false;
    }
  }

  async function register() {
    authError = '';
    if (registerForm.password !== registerForm.confirmPassword) {
      authError = 'Passwords do not match';
      return;
    }
    authLoading = true;
    try {
      const res = await fetch(`${API_URL}/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username: registerForm.username,
          password: registerForm.password,
          fullName: registerForm.fullName
        })
      });
      const data = await res.json();
      if (res.ok) {
        // Auto-login with returned token
        token = data.token;
        user = { userId: data.userId, username: data.username, fullName: data.fullName, role: data.role };
        localStorage.setItem('token', token);
        showAuthModal = false;
        registerForm = { username: '', password: '', confirmPassword: '', fullName: '' };
        await loadPatients();
      } else {
        authError = data.error || 'Registration failed. Please try again.';
      }
    } catch (e) {
      authError = 'Connection error. Please check your internet connection.';
    } finally {
      authLoading = false;
    }
  }

  async function registerUser() {
    userFormLoading = true;
    try {
      const res = await fetch(`${API_URL}/auth/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username: userForm.username,
          password: userForm.password
        })
      });
      const data = await res.json();
      if (res.ok) {
        closeUserModal();
        loadUsers();
      } else {
        alert(data.error || 'Failed to create user');
      }
    } catch (e) {
      alert('Connection error. Please try again.');
    } finally {
      userFormLoading = false;
    }
  }

  function logout() {
    token = null;
    user = null;
    patients = [];
    localStorage.removeItem('token');
    authLoading = false;
  }

  function openAuthModal(mode = 'login') {
    authMode = mode;
    authError = '';
    showAuthModal = true;
  }

  function closeAuthModal() {
    showAuthModal = false;
    authError = '';
    loginForm = { username: '', password: '' };
    registerForm = { username: '', password: '', confirmPassword: '', fullName: '' };
  }

  function closeProfileModal() {
    showProfileModal = false;
  }

  async function loadPatients() {
    loading = true;
    try {
      const res = await fetch(`${API_URL}/patients`, {
        headers: getHeaders()
      });
      if (res.status === 401 || res.status === 403) {
        logout();
        return;
      }
      if (!res.ok) {
        throw new Error(`HTTP error: ${res.status}`);
      }
      const data = await res.json();
      patients = data.patients || [];
    } catch (e) {
      console.error('Failed to load patients:', e);
    } finally {
      loading = false;
      authLoading = false; // Ensure spinner always stops
    }
  }

  async function savePatient() {
    try {
      const method = editingPatient ? 'PUT' : 'POST';
      const url = editingPatient
        ? `${API_URL}/patients/${editingPatient.id}`
        : `${API_URL}/patients`;

      const res = await fetch(url, {
        method,
        headers: getHeaders(),
        body: JSON.stringify(patientForm)
      });

      if (res.ok) {
        await loadPatients();
        closePatientModal();
      }
    } catch (e) {
      console.error('Failed to save patient:', e);
    }
  }

  async function deletePatient(id) {
    showConfirm('Are you sure you want to delete this patient?', async () => {
      try {
        await fetch(`${API_URL}/patients/${id}`, { method: 'DELETE', headers: getHeaders() });
        await loadPatients();
      } catch (e) {
        console.error('Failed to delete patient:', e);
      }
    });
  }

  async function saveReminder() {
    try {
      const isEditing = editingReminder !== null;
      const url = isEditing
        ? `${API_URL}/patients/${reminderForm.patientId}/reminders/${editingReminder.id}`
        : `${API_URL}/patients/${reminderForm.patientId}/reminders`;

      const res = await fetch(url, {
        method: isEditing ? 'PUT' : 'POST',
        headers: getHeaders(),
        body: JSON.stringify({
          title: reminderForm.title,
          description: reminderForm.description,
          dueDate: reminderForm.dueDate,
          priority: reminderForm.priority,
          recurrence: reminderForm.recurrence
        })
      });

      if (res.ok) {
        await loadPatients();
        closeReminderModal();
      }
    } catch (e) {
      console.error('Failed to save reminder:', e);
    }
  }

  async function toggleReminderComplete(patientId, reminderId) {
    try {
      await fetch(`${API_URL}/patients/${patientId}/reminders/${reminderId}/toggle`, {
        method: 'POST',
        headers: getHeaders()
      });
      await loadPatients();
    } catch (e) {
      console.error('Failed to toggle reminder:', e);
    }
  }

  async function deleteReminder(patientId, reminderId) {
    try {
      await fetch(`${API_URL}/patients/${patientId}/reminders/${reminderId}`, {
        method: 'DELETE',
        headers: getHeaders()
      });
      await loadPatients();
    } catch (e) {
      console.error('Failed to delete reminder:', e);
    }
  }

  // User management functions
  async function loadUsers() {
    if (user?.role !== 'superadmin') return;
    try {
      const res = await fetch(`${API_URL}/users`, {
        headers: getHeaders()
      });
      if (res.ok) {
        const data = await res.json();
        users = data.users || [];
      }
    } catch (e) {
      console.error('Failed to load users:', e);
    }
  }

  async function updateUserRole() {
    if (!editingUser) return;
    try {
      const res = await fetch(`${API_URL}/users/${editingUser.id}/role`, {
        method: 'PUT',
        headers: getHeaders(),
        body: JSON.stringify({ role: userForm.role })
      });
      if (res.ok) {
        await loadUsers();
        closeUserModal();
      }
    } catch (e) {
      console.error('Failed to update user role:', e);
    }
  }

  async function deleteUser(userId) {
    showConfirm('Are you sure you want to delete this user?', async () => {
      try {
        await fetch(`${API_URL}/users/${userId}`, {
          method: 'DELETE',
          headers: getHeaders()
        });
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

  function navigateTo(view) {
    currentView = view;
    localStorage.setItem('currentView', view);
    if (view === 'users') loadUsers();
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

  function formatDate(dateStr) {
    if (!dateStr) return '';
    const lang = $locale || 'en';
    return new Date(dateStr).toLocaleDateString(lang, {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }

  function formatRecurrence(recurrence) {
    if (!recurrence || recurrence.frequency === 'none') return '';

    const freqLabels = {
      daily: 'Daily',
      weekly: 'Weekly',
      monthly: 'Monthly',
      yearly: 'Yearly'
    };

    let label = freqLabels[recurrence.frequency] || recurrence.frequency;

    if (recurrence.interval > 1) {
      label = `Every ${recurrence.interval} ${recurrence.frequency}s`;
    }

    if (recurrence.frequency === 'weekly' && recurrence.daysOfWeek?.length > 0) {
      const dayLabels = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
      const days = recurrence.daysOfWeek.map(d => dayLabels[d]).join(', ');
      label += ` on ${days}`;
    }

    return label;
  }

  function getPriorityColor(priority) {
    const colors = {
      high: 'bg-red-100 text-red-700',
      medium: 'bg-amber-100 text-amber-700',
      low: 'bg-green-100 text-green-700'
    };
    return colors[priority] || colors.medium;
  }

  $: filteredPatients = patients.filter(p => {
    const matchesSearch = p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                          p.phone.includes(searchQuery) ||
                          p.email.toLowerCase().includes(searchQuery.toLowerCase());
    return matchesSearch;
  });

  $: upcomingReminders = patients.flatMap(p =>
    (p.reminders || []).filter(r => !r.completed && r.dueDate)
      .map(r => ({ ...r, patientName: p.name, patientId: p.id }))
  ).sort((a, b) => new Date(a.dueDate) - new Date(b.dueDate)).slice(0, 5);

  $: stats = {
    totalPatients: patients.length,
    totalReminders: patients.reduce((acc, p) => acc + (p.reminders?.length || 0), 0),
    completedReminders: patients.reduce((acc, p) => acc + (p.reminders?.filter(r => r.completed).length || 0), 0),
    pendingReminders: patients.reduce((acc, p) => acc + (p.reminders?.filter(r => !r.completed).length || 0), 0)
  };
</script>

<!-- Auth Loading Screen -->
{#if authLoading}
  <div class="min-h-screen bg-slate-50 flex items-center justify-center">
    <div class="animate-spin w-10 h-10 border-4 border-teal-600 border-t-transparent rounded-full"></div>
  </div>
<!-- Login Screen -->
{:else if !token}
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-6 sm:mb-8">
        <div class="w-14 h-14 sm:w-16 sm:h-16 bg-teal-600 rounded-xl sm:rounded-2xl flex items-center justify-center mx-auto mb-3 sm:mb-4">
          <svg class="w-8 h-8 sm:w-10 sm:h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
        </div>
        <h1 class="text-xl sm:text-2xl font-bold text-slate-900">{$t('app.name')}</h1>
        <p class="text-slate-500 mt-1 text-sm sm:text-base">{$t('app.tagline')}</p>
      </div>

      <!-- Language Switcher on Login Screen -->
      <div class="mb-6">
        <div class="flex items-center gap-1 bg-slate-100 rounded-lg p-1 max-w-xs mx-auto">
          <button
            onclick={() => setLocale('en')}
            class="flex-1 py-1.5 text-xs font-medium rounded-md transition-colors duration-200 {$locale === 'en' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'}"
          >
            English
          </button>
          <button
            onclick={() => setLocale('id')}
            class="flex-1 py-1.5 text-xs font-medium rounded-md transition-colors duration-200 {$locale === 'id' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'}"
          >
            Bahasa Indonesia
          </button>
        </div>
      </div>

      <!-- Auth Card -->
      <div class="bg-white rounded-xl sm:rounded-2xl shadow-xl p-6 sm:p-8">
        <div class="flex gap-2 mb-4 sm:mb-6">
          <button
            onclick={() => authMode = 'login'}
            class="flex-1 py-2.5 rounded-xl font-medium transition-colors duration-200 {authMode === 'login' ? 'bg-teal-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
          >
            {$t('auth.login')}
          </button>
          <button
            onclick={() => authMode = 'register'}
            class="flex-1 py-2.5 rounded-xl font-medium transition-colors duration-200 {authMode === 'register' ? 'bg-teal-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
          >
            {$t('auth.register')}
          </button>
        </div>

        {#if authError}
          <div class="mb-4 p-3 bg-red-50 border border-red-200 rounded-xl text-red-700 text-sm">
            {authError}
          </div>
        {/if}

        {#if authMode === 'login'}
          <form onsubmit={(e) => { e.preventDefault(); login(); }} class="space-y-4">
            <div>
              <label for="username" class="block text-sm font-medium text-slate-700 mb-1">{$t('auth.username')}</label>
              <input
                id="username"
                type="text"
                bind:value={loginForm.username}
                required
                class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                placeholder={$t('auth.enterUsername')}
              />
            </div>
            <div>
              <label for="password" class="block text-sm font-medium text-slate-700 mb-1">{$t('auth.password')}</label>
              <input
                id="password"
                type="password"
                bind:value={loginForm.password}
                required
                class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                placeholder={$t('auth.enterPassword')}
              />
            </div>
            <button
              type="submit"
              disabled={authLoading || !loginForm.username || !loginForm.password}
              class="w-full py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 hover:shadow-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              {#if authLoading}
                <div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                {$t('auth.loggingIn')}
              {:else}
                {$t('auth.loginButton')}
              {/if}
            </button>
          </form>
        {:else}
          <form onsubmit={(e) => { e.preventDefault(); register(); }} class="space-y-4">
            <div>
              <label for="fullName" class="block text-sm font-medium text-slate-700 mb-1">{$t('auth.fullName')}</label>
              <input
                id="fullName"
                type="text"
                bind:value={registerForm.fullName}
                class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                placeholder={$t('auth.yourFullName')}
              />
            </div>
            <div>
              <label for="regUsername" class="block text-sm font-medium text-slate-700 mb-1">{$t('auth.username')} *</label>
              <input
                id="regUsername"
                type="text"
                bind:value={registerForm.username}
                required
                minlength="3"
                class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200 {registerForm.username.length > 0 && !usernameValid ? 'ring-2 ring-red-500' : ''}"
                placeholder={$t('auth.chooseUsername')}
              />
              {#if registerForm.username.length > 0 && !usernameValid}
                <p class="text-xs text-red-500 mt-1">{$t('validation.usernameMin')}</p>
              {/if}
            </div>
            <div>
              <label for="regPassword" class="block text-sm font-medium text-slate-700 mb-1">{$t('auth.password')} *</label>
              <input
                id="regPassword"
                type="password"
                bind:value={registerForm.password}
                required
                minlength="6"
                class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                placeholder={$t('auth.choosePassword')}
              />
              <!-- Password strength indicator -->
              {#if registerForm.password.length > 0}
                <div class="mt-2">
                  <div class="flex gap-1 mb-1">
                    {#each [1, 2, 3, 4, 5] as level}
                      <div class="h-1 flex-1 rounded-full transition-colors duration-200 {passwordStrength >= level ? (passwordStrength <= 2 ? 'bg-red-500' : passwordStrength <= 3 ? 'bg-amber-500' : 'bg-green-500') : 'bg-slate-200'}"></div>
                    {/each}
                  </div>
                  <p class="text-xs text-slate-500">
                    {#if passwordStrength <= 1}
                      {$t('password.strength.weak')}
                    {:else if passwordStrength <= 2}
                      {$t('password.strength.fair')}
                    {:else if passwordStrength <= 3}
                      {$t('password.strength.good')}
                    {:else}
                      {$t('password.strength.strong')}
                    {/if}
                    {#if registerForm.password.length < 6}
                      - {$t('auth.minChars', { values: { n: 6 } })}
                    {/if}
                  </p>
                </div>
              {/if}
            </div>
            <div>
              <label for="confirmPassword" class="block text-sm font-medium text-slate-700 mb-1">{$t('auth.confirmPassword')} *</label>
              <input
                id="confirmPassword"
                type="password"
                bind:value={registerForm.confirmPassword}
                required
                class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200 {registerForm.confirmPassword.length > 0 && !passwordMatch ? 'ring-2 ring-red-500' : ''}"
                placeholder={$t('auth.confirmYourPassword')}
              />
              {#if registerForm.confirmPassword.length > 0 && !passwordMatch}
                <p class="text-xs text-red-500 mt-1">{$t('auth.passwordsDoNotMatch')}</p>
              {/if}
            </div>
            <button
              type="submit"
              disabled={!formValid || authLoading}
              class="w-full py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 hover:shadow-lg transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              {#if authLoading}
                <div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                {$t('auth.creatingAccount')}
              {:else}
                {$t('auth.registerButton')}
              {/if}
            </button>
          </form>
        {/if}
      </div>
    </div>
  </div>
<!-- Main Dashboard -->
{:else}
  <!-- Sidebar (Desktop Only) -->
  <aside class="hidden lg:flex flex-col w-64 bg-white border-r border-slate-200 fixed inset-y-0 left-0 z-30">
    <!-- Logo -->
    <div class="flex items-center gap-3 px-6 py-5 border-b border-slate-100">
      <div class="w-10 h-10 bg-teal-600 rounded-xl flex items-center justify-center">
        <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>
      </div>
      <div>
        <h1 class="font-bold text-slate-900">{$t('app.name')}</h1>
        <p class="text-xs text-slate-500">{$t('navigation.volunteerDashboard')}</p>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 py-4 overflow-y-auto">
      <div class="px-4 space-y-1">
        <button onclick={() => navigateTo('dashboard')} class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-left transition-all duration-200 {currentView === 'dashboard' ? 'bg-teal-50 text-teal-700 font-medium' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          {$t('navigation.dashboard')}
        </button>
        <button onclick={() => navigateTo('patients')} class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-left transition-all duration-200 {currentView === 'patients' ? 'bg-teal-50 text-teal-700 font-medium' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          {$t('navigation.patients')}
          <span class="ml-auto bg-slate-100 text-slate-600 text-xs font-medium px-2 py-0.5 rounded-full">{stats.totalPatients}</span>
        </button>
        {#if user?.role === 'superadmin'}
          <button onclick={() => navigateTo('users')} class="w-full flex items-center gap-3 px-4 py-3 rounded-xl text-left transition-all duration-200 {currentView === 'users' ? 'bg-purple-50 text-purple-700 font-medium' : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'}">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            {$t('navigation.users')}
            <span class="ml-auto bg-purple-100 text-purple-700 text-xs font-medium px-2 py-0.5 rounded-full">{users.length}</span>
          </button>
        {/if}
      </div>
    </nav>

    <!-- User Section -->
    <div class="p-4 border-t border-slate-100">
      <!-- Language Switcher -->
      <div class="mb-3 px-4 py-2">
        <div class="flex items-center gap-1 bg-slate-100 rounded-lg p-1">
          <button onclick={() => setLocale('en')} class="flex-1 py-1.5 text-xs font-medium rounded-md transition-colors {$locale === 'en' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'}">EN</button>
          <button onclick={() => setLocale('id')} class="flex-1 py-1.5 text-xs font-medium rounded-md transition-colors {$locale === 'id' ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-600 hover:text-slate-900'}">ID</button>
        </div>
      </div>
      <div class="flex items-center gap-3 px-4 py-3 rounded-xl bg-slate-50">
        <div class="w-10 h-10 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold">{user?.fullName?.charAt(0)?.toUpperCase() || user?.username?.charAt(0)?.toUpperCase() || 'U'}</div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-slate-900 truncate">{user?.fullName || user?.username}</p>
          <div class="flex items-center gap-2">
            <p class="text-xs text-slate-500 truncate">@{user?.username}</p>
            {#if user?.role}<span class="px-1.5 py-0.5 text-xs font-medium rounded {user?.role === 'superadmin' ? 'bg-purple-100 text-purple-700' : user?.role === 'admin' ? 'bg-blue-100 text-blue-700' : 'bg-slate-200 text-slate-600'}">{user?.role}</span>{/if}
          </div>
        </div>
        <button onclick={logout} class="p-2 text-slate-400 hover:text-red-600 transition-colors" title={$t('auth.logout')}>
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
        </button>
      </div>
    </div>
  </aside>

  <!-- Main Content -->
  <main class="flex-1 lg:ml-64 p-4 lg:p-6 min-h-screen pb-20 lg:pb-6">
      {#if loading}
        <div class="flex items-center justify-center h-64">
          <div class="animate-spin w-10 h-10 border-4 border-teal-600 border-t-transparent rounded-full"></div>
        </div>
      {:else}
        <!-- Dashboard View -->
        {#if currentView === 'dashboard'}
          <!-- Header -->
          <header class="sticky top-0 z-10 bg-white/80 backdrop-blur-md border-b border-slate-200 mb-6">
            <div class="flex flex-col sm:flex-row sm:items-center justify-between py-4">
              <div class="flex items-center gap-3">
                <h1 class="text-xl font-bold text-slate-900">{$t('dashboard.title')}</h1>
                <span class="text-slate-500 text-sm hidden sm:inline">{new Date().toLocaleDateString($locale, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</span>
              </div>
              <span class="text-slate-500 text-sm sm:hidden">{new Date().toLocaleDateString($locale, { weekday: 'short', month: 'short', day: 'numeric' })}</span>
            </div>
          </header>

          <!-- Stats -->
          <div class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-6 mb-8">
            <div class="bg-white rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-slate-200 hover:shadow-lg transition-shadow duration-200">
              <div class="flex items-center justify-between mb-3 sm:mb-4">
                <div class="w-10 h-10 sm:w-12 sm:h-12 bg-teal-100 rounded-xl flex items-center justify-center">
                  <svg class="w-5 h-5 sm:w-6 sm:h-6 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                  </svg>
                </div>
              </div>
              <div class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">{stats.totalPatients}</div>
              <div class="text-slate-500 font-medium text-xs sm:text-base">{$t('dashboard.totalPatients')}</div>
            </div>

            <div class="bg-white rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-slate-200 hover:shadow-lg transition-shadow duration-200">
              <div class="flex items-center justify-between mb-3 sm:mb-4">
                <div class="w-10 h-10 sm:w-12 sm:h-12 bg-blue-100 rounded-xl flex items-center justify-center">
                  <svg class="w-5 h-5 sm:w-6 sm:h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                  </svg>
                </div>
              </div>
              <div class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">{stats.totalReminders}</div>
              <div class="text-slate-500 font-medium text-xs sm:text-base">{$t('dashboard.totalReminders')}</div>
            </div>

            <div class="bg-white rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-slate-200 hover:shadow-lg transition-shadow duration-200">
              <div class="flex items-center justify-between mb-3 sm:mb-4">
                <div class="w-10 h-10 sm:w-12 sm:h-12 bg-green-100 rounded-xl flex items-center justify-center">
                  <svg class="w-5 h-5 sm:w-6 sm:h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
              </div>
              <div class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">{stats.completedReminders}</div>
              <div class="text-slate-500 font-medium text-xs sm:text-base">{$t('dashboard.completed')}</div>
            </div>

            <div class="bg-white rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-slate-200 hover:shadow-lg transition-shadow duration-200">
              <div class="flex items-center justify-between mb-3 sm:mb-4">
                <div class="w-10 h-10 sm:w-12 sm:h-12 bg-amber-100 rounded-xl flex items-center justify-center">
                  <svg class="w-5 h-5 sm:w-6 sm:h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
              </div>
              <div class="text-2xl sm:text-3xl font-bold text-slate-900 mb-1">{stats.pendingReminders}</div>
              <div class="text-slate-500 font-medium text-xs sm:text-base">{$t('dashboard.pending')}</div>
            </div>
          </div>

          <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 lg:gap-6">
            <!-- Upcoming Reminders -->
            <div class="bg-white rounded-xl lg:rounded-2xl border border-slate-200 overflow-hidden">
              <div class="px-4 lg:px-6 py-3 lg:py-4 border-b border-slate-100">
                <h2 class="text-base lg:text-lg font-semibold text-slate-900">{$t('dashboard.upcomingReminders')}</h2>
              </div>
              <div class="p-4 lg:p-6 space-y-3 lg:space-y-4">
                {#if upcomingReminders.length === 0}
                  <p class="text-slate-500 text-center py-6 lg:py-8">{$t('dashboard.noUpcomingReminders')}</p>
                {:else}
                  {#each upcomingReminders as reminder}
                    <div class="flex items-start gap-3 p-3 lg:p-4 rounded-lg lg:rounded-xl bg-slate-50 hover:bg-slate-100 transition-colors duration-200">
                      <button
                        onclick={() => toggleReminderComplete(reminder.patientId, reminder.id)}
                        class="mt-0.5 w-5 h-5 rounded-full border-2 flex-shrink-0 transition-colors duration-200 {reminder.completed ? 'bg-teal-600 border-teal-600' : 'border-slate-300 hover:border-teal-600'}"
                      >
                        {#if reminder.completed}
                          <svg class="w-full h-full text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                          </svg>
                        {/if}
                      </button>
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-1 flex-wrap">
                          <h3 class="font-medium text-slate-900 truncate">{reminder.title}</h3>
                          <span class="px-2 py-0.5 text-xs font-medium rounded-full {getPriorityColor(reminder.priority)}">{reminder.priority}</span>
                          {#if reminder.recurrence && reminder.recurrence.frequency !== 'none'}
                            <span class="px-2 py-0.5 text-xs font-medium rounded-full bg-purple-100 text-purple-700">
                              {formatRecurrence(reminder.recurrence)}
                            </span>
                          {/if}
                        </div>
                        <p class="text-sm text-slate-600 truncate">{reminder.patientName}</p>
                        <p class="text-xs text-slate-400 mt-1">{formatDate(reminder.dueDate)}</p>
                      </div>
                    </div>
                  {/each}
                {/if}
              </div>
            </div>

            <!-- Recent Patients -->
            <div class="bg-white rounded-xl lg:rounded-2xl border border-slate-200 overflow-hidden">
              <div class="px-4 lg:px-6 py-3 lg:py-4 border-b border-slate-100 flex items-center justify-between">
                <h2 class="text-base lg:text-lg font-semibold text-slate-900">{$t('dashboard.recentPatients')}</h2>
                <button onclick={() => navigateTo('patients')} class="text-xs lg:text-sm text-teal-600 hover:text-teal-700 font-medium">{$t('common.viewAll')}</button>
              </div>
              <div class="p-4 lg:p-6 space-y-3 lg:space-y-4">
                {#if patients.length === 0}
                  <p class="text-slate-500 text-center py-6 lg:py-8">{$t('dashboard.noPatients')}</p>
                {:else}
                  {#each patients.slice(0, 5) as patient}
                    <div class="flex items-center gap-3 p-3 lg:p-4 rounded-lg lg:rounded-xl bg-slate-50 hover:bg-slate-100 transition-colors duration-200">
                      <div class="w-10 h-10 lg:w-12 lg:h-12 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-lg lg:text-xl">
                        {patient.name.charAt(0).toUpperCase()}
                      </div>
                      <div class="flex-1 min-w-0">
                        <h3 class="font-medium text-slate-900 truncate">{patient.name}</h3>
                        <p class="text-sm text-slate-500 truncate">{patient.phone || $t('patients.noPhone')}</p>
                      </div>
                      <div class="text-right">
                        <div class="text-sm font-medium text-slate-900">{patient.reminders?.filter(r => r.completed).length || 0}/{patient.reminders?.length || 0}</div>
                        <div class="text-xs text-slate-500">{$t('patients.reminders')}</div>
                      </div>
                    </div>
                  {/each}
                {/if}
              </div>
            </div>
          </div>
        {/if}

        <!-- Patients View -->
        {#if currentView === 'patients'}
          <!-- Header -->
          <header class="sticky top-0 z-10 bg-white/80 backdrop-blur-md border-b border-slate-200 mb-6">
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 py-4">
              <div class="flex flex-col sm:flex-row sm:items-center gap-4">
                <h1 class="text-xl font-bold text-slate-900">{$t('patients.title')}</h1>
                <div class="relative w-full sm:w-64">
                  <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                  <input
                    type="text"
                    bind:value={searchQuery}
                    placeholder={$t('common.searchPlaceholder')}
                    class="pl-10 pr-4 py-2.5 w-full bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                  />
                </div>
              </div>
              <button
                onclick={() => openPatientModal()}
                class="flex items-center justify-center gap-2 px-5 py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200 w-full sm:w-auto"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                {$t('patients.addPatient')}
              </button>
            </div>
          </header>

          {#if filteredPatients.length === 0}
            <div class="bg-white rounded-2xl border border-slate-200 p-12 text-center">
              <div class="w-16 h-16 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <svg class="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-slate-900 mb-2">{$t('patients.noPatients')}</h3>
              <p class="text-slate-500 mb-6">
                {#if searchQuery}
                  {$t('patients.noPatientsMatch')}
                {:else}
                  {$t('patients.getStarted')}
                {/if}
              </p>
              <button
                onclick={() => openPatientModal()}
                class="px-6 py-3 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 transition-colors duration-200"
              >
                {$t('patients.addPatient')}
              </button>
            </div>
          {:else}
            <div class="space-y-4">
              {#each filteredPatients as patient}
                <div class="bg-white rounded-2xl border border-slate-200 overflow-hidden hover:shadow-lg transition-all duration-200">
                  <div class="p-4 sm:p-6">
                    <div class="flex items-start gap-3 sm:gap-4">
                      <div class="w-12 h-12 sm:w-14 sm:h-14 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-lg sm:text-xl flex-shrink-0">
                        {patient.name.charAt(0).toUpperCase()}
                      </div>
                      <div class="flex-1 min-w-0">
                        <div class="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
                          <div class="min-w-0">
                            <h3 class="text-lg font-semibold text-slate-900 truncate">{patient.name}</h3>
                            {#if patient.phone}
                              <p class="text-slate-600 text-sm">{patient.phone}</p>
                            {/if}
                            {#if patient.email}
                              <p class="text-slate-500 text-sm truncate">{patient.email}</p>
                            {/if}
                            {#if patient.notes}
                              <p class="text-slate-500 text-sm mt-2 line-clamp-2">{patient.notes}</p>
                            {/if}
                          </div>
                          <div class="flex items-center gap-1 sm:gap-2 flex-shrink-0 self-end sm:self-start">
                            <button
                              onclick={() => openReminderModal(patient.id)}
                              class="p-2 text-teal-600 hover:bg-teal-50 rounded-lg transition-colors duration-200"
                              title={$t('patients.addReminder')}
                            >
                              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                              </svg>
                            </button>
                            <button
                              onclick={() => openPatientModal(patient)}
                              class="p-2 text-slate-600 hover:bg-slate-100 rounded-lg transition-colors duration-200"
                              title={$t('common.edit')}
                            >
                              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                              </svg>
                            </button>
                            <button
                              onclick={() => deletePatient(patient.id)}
                              class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors duration-200"
                              title={$t('common.delete')}
                            >
                              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                              </svg>
                            </button>
                          </div>
                        </div>

                        <!-- Reminders section -->
                        {#if patient.reminders && patient.reminders.length > 0}
                          <div class="mt-4 pt-4 border-t border-slate-100">
                            <h4 class="text-sm font-medium text-slate-700 mb-3">{$t('patients.reminders')}</h4>
                            <div class="space-y-2">
                              {#each patient.reminders as reminder}
                                <div class="flex flex-wrap items-center gap-2 p-3 rounded-lg {reminder.completed ? 'bg-slate-50' : 'bg-amber-50'}">
                                  <button
                                    onclick={() => toggleReminderComplete(patient.id, reminder.id)}
                                    class="w-5 h-5 rounded-full border-2 flex-shrink-0 transition-colors duration-200 {reminder.completed ? 'bg-teal-600 border-teal-600' : 'border-slate-300 hover:border-teal-600'}"
                                  >
                                    {#if reminder.completed}
                                      <svg class="w-full h-full text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                                      </svg>
                                    {/if}
                                  </button>
                                  <div class="flex-1 min-w-0">
                                    <span class="{reminder.completed ? 'text-slate-500 line-through' : 'text-slate-900'} truncate block">{reminder.title}</span>
                                    {#if reminder.dueDate}
                                      <span class="text-xs text-slate-400">{formatDate(reminder.dueDate)}</span>
                                    {/if}
                                    {#if reminder.recurrence && reminder.recurrence.frequency !== 'none'}
                                      <span class="text-xs text-purple-600 ml-1" title={formatRecurrence(reminder.recurrence)}>
                                        <svg class="w-3 h-3 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                                        </svg>
                                      </span>
                                    {/if}
                                  </div>
                                  <span class="px-2 py-0.5 text-xs font-medium rounded-full {getPriorityColor(reminder.priority)}">{reminder.priority}</span>
                                  <div class="flex items-center gap-1">
                                    <button
                                      onclick={() => openReminderModal(patient.id, reminder)}
                                      class="p-1.5 text-slate-400 hover:text-teal-600 hover:bg-white/50 rounded transition-colors duration-200"
                                      aria-label={$t('common.edit')}
                                    >
                                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                      </svg>
                                    </button>
                                    <button
                                      onclick={() => deleteReminder(patient.id, reminder.id)}
                                      class="p-1.5 text-slate-400 hover:text-red-600 hover:bg-white/50 rounded transition-colors duration-200"
                                      aria-label={$t('common.delete')}
                                    >
                                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                                      </svg>
                                    </button>
                                  </div>
                                </div>
                              {/each}
                            </div>
                          </div>
                        {/if}
                      </div>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        {/if}

        <!-- Users View (Superadmin only) -->
        {#if currentView === 'users' && user?.role === 'superadmin'}
          <!-- Header -->
          <header class="sticky top-0 z-10 bg-white/80 backdrop-blur-md border-b border-slate-200 mb-6">
            <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 py-4">
              <div class="flex items-center gap-3">
                <h1 class="text-xl font-bold text-slate-900">{$t('users.title')}</h1>
              </div>
              <div class="flex items-center gap-2">
                <button
                  onclick={loadUsers}
                  class="p-2 text-slate-400 hover:text-slate-600 rounded-lg hover:bg-slate-100 transition-colors"
                  title={$t('common.refresh')}
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                </button>
                <button
                  onclick={() => openUserModal()}
                  class="flex items-center gap-2 px-5 py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 hover:shadow-lg hover:-translate-y-0.5 transition-all duration-200 w-full sm:w-auto justify-center"
                >
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  </svg>
                  {$t('users.addUser')}
                </button>
              </div>
            </div>
          </header>

          <div class="bg-white rounded-2xl border border-slate-200 overflow-hidden">
            <div class="overflow-x-auto">
              {#if users.length === 0}
                <p class="text-slate-500 text-center py-12">{$t('users.noUsers')}</p>
              {:else}
                <div class="divide-y divide-slate-100">
                  {#each users as u}
                    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 p-4 sm:px-6 sm:py-4 hover:bg-slate-50 transition-colors">
                      <div class="flex items-center gap-3">
                        <div class="w-10 h-10 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold flex-shrink-0">
                          {u.fullName?.charAt(0)?.toUpperCase() || u.username?.charAt(0)?.toUpperCase() || 'U'}
                        </div>
                        <div class="min-w-0">
                          <p class="text-sm font-medium text-slate-900 truncate">{u.fullName || $t('users.noName')}</p>
                          <p class="text-sm text-slate-500">@{u.username}</p>
                        </div>
                      </div>
                      <div class="flex items-center justify-between sm:gap-6 ml-0">
                        <div class="flex items-center gap-4">
                          <span class="px-2 py-1 text-xs font-medium rounded-full
                            {u.role === 'superadmin' ? 'bg-purple-100 text-purple-700' :
                             u.role === 'admin' ? 'bg-blue-100 text-blue-700' :
                             'bg-slate-100 text-slate-700'}">
                            {u.role}
                          </span>
                          <span class="hidden sm:block text-sm text-slate-500">
                            {u.createdAt ? new Date(u.createdAt).toLocaleDateString($locale) : '-'}
                          </span>
                        </div>
                        <div class="flex items-center gap-1 sm:gap-2">
                          {#if u.id !== user.userId}
                            <button
                              onclick={() => openUserModal(u)}
                              class="p-2 text-slate-400 hover:text-teal-600 hover:bg-teal-50 rounded-lg transition-colors"
                              title={$t('common.edit')}
                            >
                              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                              </svg>
                            </button>
                            <button
                              onclick={() => deleteUser(u.id)}
                              class="p-2 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                              title={$t('common.delete')}
                            >
                              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                              </svg>
                            </button>
                          {:else}
                            <span class="text-xs text-slate-400 px-2">{$t('common.you')}</span>
                          {/if}
                        </div>
                      </div>
                    </div>
                  {/each}
                </div>
              {/if}
            </div>
          </div>
        {/if}
      {/if}
    </main>

    <!-- Bottom Navigation (Mobile Only) -->
    <nav class="lg:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-slate-200 px-2 py-2 pb-safe z-30">
      <div class="flex items-center justify-around">
        <button
          onclick={() => navigateTo('dashboard')}
          class="flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors {currentView === 'dashboard' ? 'text-teal-600 bg-teal-50' : 'text-slate-500 hover:text-slate-700'}"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <span class="text-xs font-medium">{$t('navigation.dashboard')}</span>
        </button>

        <button
          onclick={() => navigateTo('patients')}
          class="relative flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors {currentView === 'patients' ? 'text-teal-600 bg-teal-50' : 'text-slate-500 hover:text-slate-700'}"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          <span class="text-xs font-medium">{$t('navigation.patients')}</span>
          {#if stats.totalPatients > 0}
            <span class="absolute top-0 right-2 w-4 h-4 bg-teal-600 text-white text-[10px] font-bold rounded-full flex items-center justify-center">{stats.totalPatients}</span>
          {/if}
        </button>

        {#if user?.role === 'superadmin'}
          <button
            onclick={() => navigateTo('users')}
            class="relative flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors {currentView === 'users' ? 'text-purple-600 bg-purple-50' : 'text-slate-500 hover:text-slate-700'}"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <span class="text-xs font-medium">{$t('navigation.users')}</span>
            {#if users.length > 0}
              <span class="absolute top-0 right-2 w-4 h-4 bg-purple-600 text-white text-[10px] font-bold rounded-full flex items-center justify-center">{users.length}</span>
            {/if}
          </button>
        {/if}

        <button
          onclick={() => showProfileModal = true}
          class="flex flex-col items-center gap-1 px-4 py-2 rounded-xl transition-colors text-slate-500 hover:text-slate-700"
        >
          <div class="w-8 h-8 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-sm">
            {user?.fullName?.charAt(0)?.toUpperCase() || user?.username?.charAt(0)?.toUpperCase() || 'U'}
          </div>
        </button>
      </div>
    </nav>

  <!-- Profile Modal (Mobile) -->
  {#if showProfileModal}
    <div class="fixed inset-0 z-50 flex items-end lg:items-center justify-center p-4">
      <div
        class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
        onclick={closeProfileModal}
        onkeydown={(e) => e.key === 'Escape' && closeProfileModal()}
        role="button"
        tabindex="0"
        aria-label="Close modal"
      ></div>
      <div class="relative bg-white rounded-t-2xl lg:rounded-2xl shadow-xl w-full max-w-sm p-4 pb-8 lg:p-6" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()} role="dialog" aria-modal="true" tabindex="-1">
        <!-- Drag handle -->
        <div class="lg:hidden flex justify-center mb-4">
          <div class="w-12 h-1.5 bg-slate-200 rounded-full"></div>
        </div>

        <!-- User info -->
        <div class="flex items-center gap-3 mb-6">
          <div class="w-12 h-12 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-lg">
            {user?.fullName?.charAt(0)?.toUpperCase() || user?.username?.charAt(0)?.toUpperCase() || 'U'}
          </div>
          <div class="flex-1 min-w-0">
            <p class="font-medium text-slate-900 truncate">{user?.fullName || user?.username}</p>
            <p class="text-sm text-slate-500 truncate">@{user?.username}</p>
          </div>
          <button onclick={logout} class="p-2 text-slate-400 hover:text-red-600 transition-colors" title={$t('auth.logout')}>
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
            </svg>
          </button>
        </div>

        <!-- Language switcher -->
        <div class="border-t border-slate-100 pt-4">
          <p class="text-sm font-medium text-slate-700 mb-2">{$t('common.language')}</p>
          <div class="flex gap-2">
            <button
              onclick={() => setLocale('en')}
              class="flex-1 py-2.5 rounded-xl font-medium transition-colors duration-200 {$locale === 'en' ? 'bg-teal-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
            >
              English
            </button>
            <button
              onclick={() => setLocale('id')}
              class="flex-1 py-2.5 rounded-xl font-medium transition-colors duration-200 {$locale === 'id' ? 'bg-teal-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
            >
              Bahasa Indonesia
            </button>
          </div>
        </div>
      </div>
    </div>
  {/if}

  <!-- Patient Modal -->
  {#if showPatientModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
        onclick={closePatientModal}
        onkeydown={(e) => e.key === 'Escape' && closePatientModal()}
        role="button"
        tabindex="0"
        aria-label="Close modal"
      ></div>
      <div class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-md p-4 sm:p-6" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()} role="dialog" aria-modal="true" tabindex="-1">
        <h2 class="text-lg sm:text-xl font-semibold text-slate-900 mb-4 sm:mb-6">{editingPatient ? $t('patients.editPatient') : $t('patients.addPatient')}</h2>
        <form onsubmit={(e) => { e.preventDefault(); savePatient(); }} class="space-y-4">
          <div>
            <label for="name" class="block text-sm font-medium text-slate-700 mb-1">{$t('patients.patientName')} *</label>
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
            <label for="phone" class="block text-sm font-medium text-slate-700 mb-1">{$t('patients.whatsappNumber')} *</label>
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
            <label for="email" class="block text-sm font-medium text-slate-700 mb-1">{$t('patients.email')}</label>
            <input
              id="email"
              type="email"
              bind:value={patientForm.email}
              class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
              placeholder={$t('patients.emailPlaceholder')}
            />
          </div>
          <div>
            <label for="notes" class="block text-sm font-medium text-slate-700 mb-1">{$t('patients.notes')}</label>
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
              onclick={closePatientModal}
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

  <!-- Reminder Modal -->
  {#if showReminderModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
        onclick={closeReminderModal}
        onkeydown={(e) => e.key === 'Escape' && closeReminderModal()}
        role="button"
        tabindex="0"
        aria-label="Close modal"
      ></div>
      <div class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-md p-4 sm:p-6 max-h-[90vh] overflow-y-auto" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()} role="dialog" aria-modal="true" tabindex="-1">
        <h2 class="text-xl font-semibold text-slate-900 mb-6">{editingReminder ? $t('patients.editReminder') : $t('patients.addReminder')}</h2>
        <form onsubmit={(e) => { e.preventDefault(); saveReminder(); }} class="space-y-4">
          <div>
            <label for="title" class="block text-sm font-medium text-slate-700 mb-1">{$t('reminders.title')} *</label>
            <input
              id="title"
              type="text"
              bind:value={reminderForm.title}
              required
              class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
              placeholder={$t('reminders.titlePlaceholder')}
            />
          </div>
          <div>
            <label for="description" class="block text-sm font-medium text-slate-700 mb-1">{$t('reminders.description')}</label>
            <textarea
              id="description"
              bind:value={reminderForm.description}
              rows="3"
              class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200 resize-none"
              placeholder={$t('reminders.descriptionPlaceholder')}
            ></textarea>
          </div>
          <div>
            <label for="dueDate" class="block text-sm font-medium text-slate-700 mb-1">{$t('reminders.dueDate')}</label>
            <input
              id="dueDate"
              type="datetime-local"
              bind:value={reminderForm.dueDate}
              class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            />
          </div>
          <div>
            <label for="priority" class="block text-sm font-medium text-slate-700 mb-1">{$t('reminders.priority')}</label>
            <select
              id="priority"
              bind:value={reminderForm.priority}
              class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            >
              <option value="low">{$t('reminders.low')}</option>
              <option value="medium">{$t('reminders.medium')}</option>
              <option value="high">{$t('reminders.high')}</option>
            </select>
          </div>

          <!-- Recurrence Section -->
          <div class="pt-4 border-t border-slate-100">
            <h3 class="text-sm font-medium text-slate-700 mb-3">{$t('reminders.recurrence')}</h3>
            <div class="space-y-3">
              <div>
                <label for="frequency" class="block text-xs text-slate-500 mb-1">{$t('reminders.repeat')}</label>
                <select
                  id="frequency"
                  bind:value={reminderForm.recurrence.frequency}
                  class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                >
                  <option value="none">{$t('reminders.doesNotRepeat')}</option>
                  <option value="daily">{$t('reminders.daily')}</option>
                  <option value="weekly">{$t('reminders.weekly')}</option>
                  <option value="monthly">{$t('reminders.monthly')}</option>
                  <option value="yearly">{$t('reminders.yearly')}</option>
                </select>
              </div>

              {#if reminderForm.recurrence.frequency !== 'none'}
                <div>
                  <label for="interval" class="block text-xs text-slate-500 mb-1">{$t('reminders.repeatEvery')}</label>
                  <input
                    id="interval"
                    type="number"
                    min="1"
                    max="99"
                    bind:value={reminderForm.recurrence.interval}
                    class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                  />
                </div>

                {#if reminderForm.recurrence.frequency === 'weekly'}
                  <div>
                    <span class="block text-xs text-slate-500 mb-2">{$t('reminders.daysOfWeek')}</span>
                    <div class="flex gap-1" role="group" aria-label={$t('reminders.daysOfWeek')}>
                      {#each daysOfWeek as day}
                        <button
                          type="button"
                          onclick={() => toggleDayOfWeek(day.value)}
                          class="w-8 h-8 text-xs font-medium rounded-lg transition-colors duration-200 {reminderForm.recurrence.daysOfWeek.includes(day.value) ? 'bg-teal-600 text-white' : 'bg-slate-100 text-slate-600 hover:bg-slate-200'}"
                        >
                          {day.label}
                        </button>
                      {/each}
                    </div>
                  </div>
                {/if}

                <div>
                  <label for="endDate" class="block text-xs text-slate-500 mb-1">{$t('reminders.endDate')}</label>
                  <input
                    id="endDate"
                    type="date"
                    bind:value={reminderForm.recurrence.endDate}
                    class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                  />
                </div>
              {/if}
            </div>
          </div>

          <div class="flex flex-col sm:flex-row gap-3 pt-4">
            <button
              type="button"
              onclick={closeReminderModal}
              class="flex-1 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1"
            >
              {$t('common.cancel')}
            </button>
            <button
              type="submit"
              class="flex-1 px-4 py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 transition-colors duration-200 order-1 sm:order-2"
            >
              {editingReminder ? $t('common.save') : $t('patients.addReminder')}
            </button>
          </div>
        </form>
      </div>
    </div>
  {/if}

  <!-- User Role Modal -->
  {#if showUserModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
        onclick={closeUserModal}
        onkeydown={(e) => e.key === 'Escape' && closeUserModal()}
        role="button"
        tabindex="0"
        aria-label="Close modal"
      ></div>
      <div class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-sm p-4 sm:p-6" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()} role="dialog" aria-modal="true" tabindex="-1">
        <h2 class="text-lg sm:text-xl font-semibold text-slate-900 mb-4 sm:mb-6">{editingUser ? $t('users.editUserRole') : $t('users.addUser')}</h2>
        {#if editingUser}
          <div class="mb-6">
            <div class="flex items-center gap-3 p-3 bg-slate-50 rounded-xl mb-4">
              <div class="w-12 h-12 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-lg">
                {editingUser.fullName?.charAt(0)?.toUpperCase() || editingUser.username?.charAt(0)?.toUpperCase() || 'U'}
              </div>
              <div>
                <p class="font-medium text-slate-900">{editingUser.fullName || $t('users.noName')}</p>
                <p class="text-sm text-slate-500">@{editingUser.username}</p>
              </div>
            </div>
            <form onsubmit={(e) => { e.preventDefault(); updateUserRole(); }} class="space-y-4">
              <div>
                <label for="userRole" class="block text-sm font-medium text-slate-700 mb-1">{$t('users.role')}</label>
                <select
                  id="userRole"
                  bind:value={userForm.role}
                  class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
                >
                  <option value="admin">{$t('users.admin')}</option>
                  <option value="volunteer">{$t('users.volunteer')}</option>
                </select>
                <p class="text-xs text-slate-500 mt-2">
                  {#if userForm.role === 'admin'}
                    {$t('users.roleDescription.admin')}
                  {:else}
                    {$t('users.roleDescription.volunteer')}
                  {/if}
                </p>
              </div>
              <div class="flex flex-col sm:flex-row gap-3 pt-4">
                <button
                  type="button"
                  onclick={closeUserModal}
                  class="flex-1 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1"
                >
                  {$t('common.cancel')}
                </button>
                <button
                  type="submit"
                  class="flex-1 px-4 py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 transition-colors duration-200 order-1 sm:order-2"
                >
                  {$t('common.save')}
                </button>
              </div>
            </form>
          </div>
        {:else}
          <form onsubmit={(e) => { e.preventDefault(); registerUser(); }} class="space-y-4">
            <div>
              <label for="newUsername" class="block text-sm font-medium text-slate-700 mb-1">{$t('auth.username')}</label>
              <input
                type="text"
                id="newUsername"
                bind:value={userForm.username}
                placeholder={$t('auth.enterUsername')}
                minlength="3"
                maxlength="30"
                required
                class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
              />
            </div>
            <div>
              <label for="newPassword" class="block text-sm font-medium text-slate-700 mb-1">{$t('auth.password')}</label>
              <input
                type="password"
                id="newPassword"
                bind:value={userForm.password}
                placeholder={$t('auth.minChars', { values: { n: 6 } })}
                minlength="6"
                required
                class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
              />
            </div>
            <div class="flex flex-col sm:flex-row gap-3 pt-4">
              <button
                type="button"
                onclick={closeUserModal}
                class="flex-1 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1"
              >
                {$t('common.cancel')}
              </button>
              <button
                type="submit"
                disabled={userFormLoading}
                class="flex-1 px-4 py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 disabled:opacity-50 transition-colors duration-200 flex items-center justify-center gap-2 order-1 sm:order-2"
              >
                {#if userFormLoading}
                  <div class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                {/if}
                {$t('users.addUser')}
              </button>
            </div>
          </form>
        {/if}
      </div>
    </div>
  {/if}

  <!-- Confirm Modal -->
  {#if showConfirmModal}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm"
        onclick={closeConfirmModal}
        onkeydown={(e) => e.key === 'Escape' && closeConfirmModal()}
        role="button"
        tabindex="0"
        aria-label="Close modal"
      ></div>
      <div class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-sm p-4 sm:p-6" onclick={(e) => e.stopPropagation()} onkeydown={(e) => e.stopPropagation()} role="dialog" aria-modal="true" tabindex="-1">
        <div class="flex items-center gap-3 sm:gap-4 mb-4 sm:mb-6">
          <div class="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center flex-shrink-0">
            <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          <p class="text-slate-700">{confirmMessage}</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3">
          <button
            onclick={closeConfirmModal}
            class="flex-1 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1"
          >
            {$t('common.cancel')}
          </button>
          <button
            onclick={handleConfirm}
            class="flex-1 px-4 py-2.5 bg-red-600 text-white font-medium rounded-xl hover:bg-red-700 transition-colors duration-200 order-1 sm:order-2"
          >
            {$t('common.delete')}
          </button>
        </div>
      </div>
    </div>
  {/if}
{/if}
