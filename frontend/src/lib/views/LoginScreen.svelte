<script>
  import { t } from 'svelte-i18n';
  import { locale } from 'svelte-i18n';

  export let authMode = 'login';
  export let authError = '';
  export let authLoading = false;
  export let loginForm = { username: '', password: '' };
  export let registerForm = { username: '', password: '', confirmPassword: '', fullName: '' };
  export let passwordStrength = 0;
  export let passwordMatch = false;
  export let usernameValid = false;
  export let formValid = false;
  export let setLocale = () => {};
  export let onLogin = () => {};
  export let onRegister = () => {};

  function getPasswordStrengthLabel(strength) {
    if (strength <= 1) return 'weak';
    if (strength <= 2) return 'fair';
    if (strength <= 3) return 'good';
    return 'strong';
  }

  function handleLoginSubmit(e) {
    e.preventDefault();
    onLogin();
  }

  function handleRegisterSubmit(e) {
    e.preventDefault();
    onRegister();
  }
</script>

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
        <form onsubmit={handleLoginSubmit} class="space-y-4">
          <div>
            <label for="username" class="block text-sm font-medium text-slate-700 mb-1">
              {$t('auth.username')}
            </label>
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
            <label for="password" class="block text-sm font-medium text-slate-700 mb-1">
              {$t('auth.password')}
            </label>
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
        <form onsubmit={handleRegisterSubmit} class="space-y-4">
          <div>
            <label for="fullName" class="block text-sm font-medium text-slate-700 mb-1">
              {$t('auth.fullName')}
            </label>
            <input
              id="fullName"
              type="text"
              bind:value={registerForm.fullName}
              class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
              placeholder={$t('auth.yourFullName')}
            />
          </div>
          <div>
            <label for="regUsername" class="block text-sm font-medium text-slate-700 mb-1">
              {$t('auth.username')} *
            </label>
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
            <label for="regPassword" class="block text-sm font-medium text-slate-700 mb-1">
              {$t('auth.password')} *
            </label>
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
                    <div
                      class="h-1 flex-1 rounded-full transition-colors duration-200 {passwordStrength >= level
                        ? (passwordStrength <= 2 ? 'bg-red-500' : passwordStrength <= 3 ? 'bg-amber-500' : 'bg-green-500')
                        : 'bg-slate-200'}"
                    ></div>
                  {/each}
                </div>
                <p class="text-xs text-slate-500">
                  {$t(`password.strength.${getPasswordStrengthLabel(passwordStrength)}`)}
                  {#if registerForm.password.length < 6}
                    - {$t('auth.minChars', { values: { n: 6 } })}
                  {/if}
                </p>
              </div>
            {/if}
          </div>
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-slate-700 mb-1">
              {$t('auth.confirmPassword')} *
            </label>
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
