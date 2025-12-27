import { writable } from 'svelte/store';

function createAuthStore() {
  const storedToken = typeof localStorage !== 'undefined' ? localStorage.getItem('token') : null;

  const { subscribe, set, update } = writable({
    token: storedToken,
    user: null,
    loading: true
  });

  return {
    subscribe,
    setUser: (user) => update(state => ({ ...state, user, loading: false })),
    setLoading: (loading) => update(state => ({ ...state, loading })),
    setToken: (token) => {
      if (typeof localStorage !== 'undefined') {
        if (token) {
          localStorage.setItem('token', token);
        } else {
          localStorage.removeItem('token');
        }
      }
      update(state => ({ ...state, token }));
    },
    logout: () => {
      if (typeof localStorage !== 'undefined') {
        localStorage.removeItem('token');
      }
      set({ token: null, user: null, loading: false });
    },
    reset: () => set({ token: null, user: null, loading: false })
  };
}

export const auth = createAuthStore();
