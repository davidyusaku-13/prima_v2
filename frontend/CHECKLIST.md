# Pre-Implementation Checklist

**Use this checklist BEFORE starting any implementation task.**

---

## Before Writing Code

- [ ] Review story requirements and acceptance criteria
- [ ] Check for complex patterns (SSE, WebSocket, real-time)
- [ ] Create spike/prototype if needed
- [ ] Identify all integration points
- [ ] Verify this is Vite + Svelte 5, NOT SvelteKit

---

## Svelte 5 Reactivity

### Component Props Pattern

**Svelte 5 uses `$props()` instead of `export let`:**

```javascript
// ❌ OLD (Svelte 4)
export let status = 'pending';
export let onRetry = null;

// ✅ NEW (Svelte 5)
let { status = 'pending', onRetry = null } = $props();
```

**Notes:**
- Destructuring with defaults works naturally
- Props are reactive by default
- Use `$bindable()` for two-way binding if needed

### Reactivity Basics

**THE GOLDEN RULE:** Always create new object/array references when updating state.

### State Updates

| ❌ Wrong | ✅ Correct |
|----------|-----------|
| `array.push(x)` | `array = [...array, x]` |
| `array.splice(i, 1)` | `array = [...array.slice(0, i), ...array.slice(i + 1)]` |
| `object.prop = val` | `object = { ...object, prop: val }` |
| `count++` | `count = count + 1` |
| `arr[i] = x` | `arr = [...arr.slice(0, i), x, ...arr.slice(i + 1)]` |

### Derived Values

```javascript
// ✅ Use $derived for computed values
const isOverdue = $derived(dueDate < new Date());
const failedCount = $derived(
    Object.values(deliveryStatuses).filter(s => s.status === 'failed').length
);
```

### Store Pattern

```javascript
// ✅ Store should expose methods that return new references
class DeliveryStore {
    deliveryStatuses = $state({});

    updateStatus(id, status) {
        // Always return new object
        this.deliveryStatuses = {
            ...this.deliveryStatuses,
            [id]: { status, timestamp: Date.now() }
        };
    }
}
```

### Reactivity Checklist

- [ ] Using `$state()` for mutable state
- [ ] Creating NEW object/array references when updating
- [ ] Using `$derived()` for computed values
- [ ] Using `$effect()` for side effects only
- [ ] NO direct mutation: `array.push()`, `obj.prop = value`
- [ ] NO legacy `$:` reactive statements

---

## i18n

- [ ] All user-facing strings use `$t('key')` (svelte-i18n v4)
- [ ] No hardcoded text in templates
- [ ] Check hydration compatibility (use `{#if $locale}` pattern if needed)
- [ ] Translations exist in both `en.json` and `id.json`

---

## Security

- [ ] No hardcoded URLs or API endpoints
- [ ] Phone numbers masked for display: `628***890`
- [ ] HMAC validation for webhooks
- [ ] Input validation on all endpoints
- [ ] PII never logged in plain text

---

## Testing

- [ ] Unit tests for utility functions
- [ ] Component tests for UI components
- [ ] Integration tests for API handlers
- [ ] All tests pass before marking done
- [ ] Test file co-located: `Component.svelte` → `Component.test.js`

---

## Accessibility

- [ ] Interactive elements have `aria-label`
- [ ] Status updates use `role="status"` or `aria-live`
- [ ] Keyboard navigation works (Tab, Enter, Space)
- [ ] Color not the only indicator (icons + text)
- [ ] Focus ring styles for interactive elements

---

## Documentation

- [ ] File list accurate (Create vs Modify)
- [ ] Dev Notes updated after code review
- [ ] API documentation updated if needed
- [ ] No placeholder comments left behind

---

## Code Review Prep

- [ ] Self-review for hardcoded URLs
- [ ] Check for console.log/debug code
- [ ] Verify imports are valid for Vite + Svelte 5
- [ ] Run tests locally
- [ ] Check for missing error handling

---

## Common Pitfalls (Epic 3 Lessons)

1. **Direct mutation** - Always create new references
2. **Hardcoded URLs** - Use relative paths or config
3. **i18n hydration** - Use `{#if $locale}` pattern
4. **Phone masking** - Never show full numbers
5. **SvelteKit confusion** - Use native APIs, NOT `$app/*`

---

## Vite + Svelte 5 Patterns (NOT SvelteKit)

| ❌ SvelteKit | ✅ Vite + Svelte 5 |
|-------------|-------------------|
| `import { goto } from '$app/navigation'` | `window.location.href = '...'` |
| `import { page } from '$app/stores'` | Use props or context |
| `import { browser } from '$app/environment'` | `typeof window !== 'undefined'` |
| `import { onMount } from 'svelte'` | Use `$effect()` (onMount still works, $effect preferred) |
| `$: computed = ...` | `$derived(...)` |
| `export let prop` | `let { prop } = $props()` |
| SSE (Server-Sent Events) | `new EventSource('/api/sse')` |
| localStorage | `localStorage.getItem()` / `localStorage.setItem()` |
| Session storage | `sessionStorage.getItem()` / `sessionStorage.setItem()` |

---

**Remember:** When in doubt, check `project-context.md` or existing code patterns.
