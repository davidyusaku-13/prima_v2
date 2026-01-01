<script>
  import { t } from "svelte-i18n";

  export let show = false;
  export let editingUser = null;
  export let userForm = { role: "volunteer", username: "", password: "" };
  export let loading = false;
  export let onClose = () => {};
  export let onSaveRole = () => {};
  export let onRegister = () => {};

  function handleSubmit(e) {
    e.preventDefault();
    if (editingUser) {
      onSaveRole();
    } else {
      onRegister();
    }
  }
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
      class="relative bg-white rounded-xl sm:rounded-2xl shadow-xl w-full max-w-sm p-4 sm:p-6"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <h2 class="text-lg sm:text-xl font-semibold text-slate-900 mb-4 sm:mb-6">
        {editingUser ? $t("users.editUserRole") : $t("users.addUser")}
      </h2>

      {#if editingUser}
        <div class="mb-6">
          <div class="flex items-center gap-3 p-3 bg-slate-50 rounded-xl mb-4">
            <div
              class="w-12 h-12 bg-teal-600 rounded-full flex items-center justify-center text-white font-semibold text-lg"
            >
              {editingUser.fullName?.charAt(0)?.toUpperCase() ||
                editingUser.username?.charAt(0)?.toUpperCase() ||
                "U"}
            </div>
            <div>
              <p class="font-medium text-slate-900">
                {editingUser.fullName || $t("users.noName")}
              </p>
              <p class="text-sm text-slate-500">@{editingUser.username}</p>
            </div>
          </div>
          <form onsubmit={handleSubmit} class="space-y-4">
            <div>
              <label
                for="userRole"
                class="block text-sm font-medium text-slate-700 mb-1"
              >
                {$t("users.role")}
              </label>
              <select
                id="userRole"
                bind:value={userForm.role}
                class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
              >
                <option value="admin">{$t("users.admin")}</option>
                <option value="volunteer">{$t("users.volunteer")}</option>
              </select>
              <p class="text-xs text-slate-500 mt-2">
                {#if userForm.role === "admin"}
                  {$t("users.roleDescription.admin")}
                {:else}
                  {$t("users.roleDescription.volunteer")}
                {/if}
              </p>
            </div>
            <div class="flex flex-col sm:flex-row gap-3 pt-4">
              <button
                type="button"
                onclick={onClose}
                class="flex-1 h-10 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1 flex items-center justify-center"
              >
                {$t("common.cancel")}
              </button>
              <button
                type="submit"
                class="flex-1 h-10 px-4 py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 transition-colors duration-200 order-1 sm:order-2 flex items-center justify-center"
              >
                {$t("common.save")}
              </button>
            </div>
          </form>
        </div>
      {:else}
        <form onsubmit={handleSubmit} class="space-y-4">
          <div>
            <label
              for="newUsername"
              class="block text-sm font-medium text-slate-700 mb-1"
            >
              {$t("auth.username")}
            </label>
            <input
              type="text"
              id="newUsername"
              bind:value={userForm.username}
              placeholder={$t("auth.enterUsername")}
              minlength="3"
              maxlength="30"
              required
              class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            />
          </div>
          <div>
            <label
              for="newPassword"
              class="block text-sm font-medium text-slate-700 mb-1"
            >
              {$t("auth.password")}
            </label>
            <input
              type="password"
              id="newPassword"
              bind:value={userForm.password}
              placeholder={$t("auth.minChars", { values: { n: 6 } })}
              minlength="6"
              required
              class="w-full px-4 py-2.5 bg-slate-100 border-0 rounded-xl focus:outline-none focus:ring-2 focus:ring-teal-500 focus:bg-white transition-all duration-200"
            />
          </div>
          <div class="flex flex-col sm:flex-row gap-3 pt-4">
            <button
              type="button"
              onclick={onClose}
              class="flex-1 h-10 px-4 py-2.5 border border-slate-200 text-slate-700 font-medium rounded-xl hover:bg-slate-50 transition-colors duration-200 order-2 sm:order-1 flex items-center justify-center"
            >
              {$t("common.cancel")}
            </button>
            <button
              type="submit"
              disabled={loading}
              class="flex-1 h-10 px-4 py-2.5 bg-teal-600 text-white font-medium rounded-xl hover:bg-teal-700 disabled:opacity-50 transition-colors duration-200 flex items-center justify-center gap-2 order-1 sm:order-2"
            >
              {#if loading}
                <div
                  class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"
                ></div>
              {/if}
              {$t("users.addUser")}
            </button>
          </div>
        </form>
      {/if}
    </div>
  </div>
{/if}
