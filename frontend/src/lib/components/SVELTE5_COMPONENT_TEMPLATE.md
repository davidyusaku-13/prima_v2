---
title: Svelte 5 Component Template
description: Template for Svelte 5 components with proper memory leak prevention
---

# Svelte 5 + Vite Component Template

**⚠️ IMPORTANT: This project uses Svelte 5 + Vite, NOT SvelteKit**

Do NOT import from `$app/navigation` or `$app/stores`. Use `window.location` and browser APIs instead.

---

## Template

```svelte
<script>
  // ===== IMPORTS =====
  // Use browser-compatible imports only
  // import { browser } from '$app/environment'; // ❌ WRONG
  // import { onMount } from 'svelte'; // Use $effect instead
  // import { page } from '$app/stores'; // ❌ WRONG

  // ✅ CORRECT: Standard imports
  import { t } from 'svelte-i18n';
  import { someApi } from '$lib/utils/api.js';

  // ===== PROPS =====
  let {
    prop1 = default1,
    prop2 = default2,
    // Use $bindable() for two-way binding
    bindableProp = $bindable()
  } = $props();

  // ===== STATE (Svelte 5 Runes) =====
  // Use $state() for reactive state
  let stateVar = $state(initialValue);
  let computedState = $derived(expression);

  // ===== EFFECTS WITH CLEANUP =====
  // Always use AbortController for async operations
  $effect(() => {
    let aborted = false;
    const controller = new AbortController();

    // Async operation with signal
    someAsyncFunction({ signal: controller.signal })
      .then(data => {
        if (aborted) return;
        // Update state
        stateVar = data;
      })
      .catch(err => {
        if (aborted && err.name === 'AbortError') return;
        // Handle error
        console.error('Error:', err);
      });

    // Cleanup function - prevents memory leaks
    return () => {
      aborted = true;
      controller.abort();
    };
  });

  // For intervals/timeouts
  $effect(() => {
    const interval = setInterval(() => {
      // periodic work
    }, 1000);

    return () => clearInterval(interval);
  });

  // ===== DERIVED VALUES =====
  // Use $derived() for simple derived state
  let doubleValue = $derived(stateVar * 2);

  // Use $derived.by() for complex derived state
  let complexComputed = $derived.by(() => {
    // complex calculation
    return result;
  });

  // ===== EVENT HANDLERS =====
  function handleClick() {
    // Handler logic
  }

  function handleKeydown(e) {
    if (e.key === 'Enter' || e.key === ' ') {
      // Handle
    }
  }
</script>

<!-- ===== TEMPLATE ===== -->
<div
  class="container-class"
  role="region"
  aria-label="Component description"
>
  <!-- Content here -->

  <!-- Accessible click handlers -->
  <button
    onclick={handleClick}
    onkeydown={handleKeydown}
    class="btn-class"
  >
    Label
  </button>
</div>

<style>
  /* ===== STYLES ===== */
  /* Use Tailwind classes when possible */
  /* Scoped styles only when necessary */
  .container-class {
    /* styles */
  }
</style>
```

---

## Common Mistakes to Avoid

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `import { goto } from '$app/navigation'` | `window.location.href = '...'` |
| `import { page } from '$app/stores'` | Use props or context |
| `onMount(() => {...})` | `$effect(() => {...})` |
| `let count = 0` | `let count = $state(0)` |
| `$: doubled = count * 2` | `let doubled = $derived(count * 2)` |
| Missing cleanup in $effect | Always return cleanup function |

---

## Memory Leak Prevention Checklist

- [ ] All async operations use AbortController
- [ ] All intervals/timeouts are cleared in cleanup
- [ ] No SvelteKit imports used
- [ ] Props use $bindable() only when needed
- [ ] State uses $state() runes
- [ ] Derived values use $derived() / $derived.by()
- [ ] Effects have proper cleanup functions
