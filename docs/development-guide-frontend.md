# Development Guide - Frontend (Svelte 5 + Vite)

**Generated:** January 11, 2026 (Updated)
**Project:** PRIMA Healthcare Volunteer Dashboard
**Technology:** Svelte 5.43.8 + Vite 7.2.4 (**NOT SvelteKit**)
**Scan Type:** Exhaustive Rescan

---

## CRITICAL: This is NOT SvelteKit

**PRIMA uses Vite + Svelte 5 directly. It is NOT a SvelteKit project.**

The following imports will **break the build**:

| Prohibited Import | Use Instead |
|-------------------|-------------|
| `import { goto } from '$app/navigation'` | `window.location.href = '/path'` |
| `import { page } from '$app/stores'` | Use props or context |
| `import { browser } from '$app/environment'` | `typeof window !== 'undefined'` |

**Why?** SvelteKit provides SSR, routing, and `$app/*` modules. PRIMA is a pure SPA using Vite for bundling and manual routing via localStorage.

If you see `Cannot find module '$app/...'` errors, you've used a SvelteKit pattern in a non-SvelteKit project.

---

## Quick Start

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
bun install

# Run dev server (port 5173)
bun run dev

# Run tests
bun run test

# Build for production
bun run build

# Preview production build
bun run preview
```

---

## Environment Setup

### Prerequisites

- **Bun:** 1.0+ (recommended) or **Node.js:** 18+ with npm/pnpm
- **Backend:** Running on http://localhost:8080 (see backend guide)
- **Text Editor:** VS Code with Svelte extension recommended

### Install Bun (Recommended)

**Windows:**

```powershell
# Using PowerShell
powershell -c "irm bun.sh/install.ps1 | iex"
```

**macOS/Linux:**

```bash
curl -fsSL https://bun.sh/install | bash
```

**Verify Installation:**

```bash
bun --version
# Should output: 1.x.x
```

### Alternative: Use Node.js

If you prefer Node.js over Bun:

```bash
# Using npm
npm install

# Using pnpm
pnpm install

# Run commands with npm/pnpm instead of bun
npm run dev
npm run build
```

### Clone Repository

```bash
git clone <repository-url>
cd prima_v2/frontend
```

### Install Dependencies

```bash
bun install
```

This installs all packages listed in `package.json`.

---

## Project Structure

```
frontend/
├── src/
│   ├── main.js                  # Entry point
│   ├── App.svelte               # Root component (856 lines)
│   ├── app.css                  # Global styles (Tailwind)
│   ├── i18n.js                  # Internationalization setup
│   ├── assets/                  # Images, icons
│   ├── lib/
│   │   ├── components/          # Reusable UI components
│   │   ├── views/               # Page-level components
│   │   ├── stores/              # Global state (Svelte stores)
│   │   ├── services/            # External service clients (SSE)
│   │   ├── utils/               # Utility functions (API client)
│   │   └── test/                # Test utilities
│   └── locales/                 # Translation files (en.json, id.json)
├── public/                      # Static files (favicon, etc.)
├── index.html                   # HTML entry point
├── vite.config.js               # Vite build configuration
├── svelte.config.js             # Svelte compiler options
├── jsconfig.json                # JavaScript path mappings ($lib alias)
├── vitest.config.js             # Test runner configuration
├── package.json                 # Dependencies & scripts
├── bun.lockb                    # Bun lockfile
└── README.md                    # Frontend-specific docs
```

---

## Configuration

### Vite Configuration

**`vite.config.js`:**

```javascript
import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import path from "path";

export default defineConfig({
  plugins: [svelte(), tailwindcss()],
  resolve: {
    alias: {
      $lib: path.resolve(__dirname, "./src/lib"),
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
});
```

**Features:**

- **`$lib` alias:** Import from `src/lib/` using `$lib/...`
- **Proxy:** `/api` requests forwarded to backend (http://localhost:8080)
- **Tailwind:** Integrated via `@tailwindcss/vite` plugin

### API Base URL

**`src/lib/utils/api.js`:**

```javascript
const API_URL = "http://localhost:8080/api";
```

**For Production:** Update to your backend URL:

```javascript
const API_URL =
  import.meta.env.VITE_API_URL || "https://api.prima.example.com/api";
```

**Environment Variables (`.env`):**

```env
VITE_API_URL=http://localhost:8080/api
```

**Access in code:**

```javascript
const apiUrl = import.meta.env.VITE_API_URL;
```

---

## Running Locally

### 1. Start Backend First

```bash
cd backend
go run main.go
```

**Verify backend:**

```bash
curl http://localhost:8080/api/health
# {"status": "ok", "timestamp": "..."}
```

### 2. Start Frontend

```bash
cd frontend
bun run dev
```

**Expected Output:**

```
  VITE v7.2.4  ready in 234 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
  ➜  press h + enter to show help
```

**Access:**

- Frontend: http://localhost:5173
- Login with: `superadmin` / `superadmin`

### 3. Hot Module Replacement (HMR)

Vite automatically reloads changes:

- **Svelte components:** Instant update (state preserved when possible)
- **CSS:** Instant injection (no page reload)
- **JavaScript:** Fast reload

**Example:** Edit `src/App.svelte` and save → changes appear immediately in browser

---

## Testing

### Run All Tests

```bash
bun run test
```

**Runs tests in watch mode (re-runs on file changes).**

### Run Once (No Watch)

```bash
bun run test -- --run
```

### Run Specific Test File

```bash
bun run test api.test.js
bun run test stores/delivery.test.js
```

### Run with Coverage

```bash
bun run test -- --coverage
```

**Output:** Coverage report in terminal + `coverage/` folder

### Test Pattern

**File:** `src/lib/utils/api.test.js`

```javascript
import { describe, it, expect, vi, beforeEach } from "vitest";
import * as api from "./api.js";

// Mock fetch globally
global.fetch = vi.fn();

describe("API Client", () => {
  beforeEach(() => {
    fetch.mockClear();
  });

  it("should login successfully", async () => {
    // Arrange
    const mockResponse = {
      token: "test-token",
      user: { id: "1", username: "test" },
    };
    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => mockResponse,
    });

    // Act
    const result = await api.login("test", "password");

    // Assert
    expect(fetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/auth/login",
      expect.objectContaining({
        method: "POST",
        body: JSON.stringify({ username: "test", password: "password" }),
      })
    );
    expect(result).toEqual(mockResponse);
  });

  it("should throw error on login failure", async () => {
    // Arrange
    fetch.mockResolvedValueOnce({
      ok: false,
      json: async () => ({ error: "Invalid credentials" }),
    });

    // Act & Assert
    await expect(api.login("test", "wrong")).rejects.toThrow(
      "Invalid credentials"
    );
  });
});
```

**Component Test Example:**

`src/lib/components/ui/Toast.test.js`:

```javascript
import { render, fireEvent } from "@testing-library/svelte";
import { describe, it, expect } from "vitest";
import Toast from "./Toast.svelte";

describe("Toast Component", () => {
  it("should render toast message", () => {
    const { getByText } = render(Toast, {
      props: {
        toast: {
          id: 1,
          message: "Test message",
          type: "info",
        },
      },
    });

    expect(getByText("Test message")).toBeInTheDocument();
  });

  it("should call onClose when close button clicked", async () => {
    let closed = false;
    const { getByRole } = render(Toast, {
      props: {
        toast: { id: 1, message: "Test", type: "info" },
        onClose: () => {
          closed = true;
        },
      },
    });

    const closeButton = getByRole("button", { name: /close/i });
    await fireEvent.click(closeButton);

    expect(closed).toBe(true);
  });
});
```

---

## Code Style

### Svelte 5 Runes (Critical!)

**PRIMA uses Svelte 5 runes, NOT legacy reactivity.**

#### ❌ Legacy Svelte 3/4 (Don't Use)

```svelte
<script>
  export let title = 'Default';  // Old props syntax
  let count = 0;
  $: doubled = count * 2;  // Old reactive statement
  $: {  // Old reactive block
    console.log(count);
  }
</script>

<button on:click={() => count++}>Click</button>
```

#### ✅ Svelte 5 (Use This)

```svelte
<script>
  let { title = 'Default' } = $props();  // New props syntax
  let count = $state(0);  // Reactive state
  let doubled = $derived(count * 2);  // Derived state
  $effect(() => {  // Effect (side effects)
    console.log(count);
  });
</script>

<button onclick={() => count++}>Click</button>
```

### Naming Conventions

**Components:**

- `PascalCase` → `PatientModal.svelte`, `DashboardView.svelte`

**Variables:**

- `camelCase` → `let patientList = $state([]);`

**Constants:**

- `CONSTANT_CASE` → `const API_URL = '...';`

**Stores:**

- `camelCase` with `.svelte.js` extension → `auth.svelte.js`, `toast.svelte.js`

### Component Structure

**Order:**

```svelte
<script>
  // 1. Imports
  import { onMount } from 'svelte';
  import * as api from '$lib/utils/api.js';

  // 2. Props
  let { patient, onSave, onClose } = $props();

  // 3. State
  let loading = $state(false);
  let form = $state({ name: patient?.name || '', phone: patient?.phone || '' });

  // 4. Derived state
  let isValid = $derived(form.name && form.phone);

  // 5. Functions
  async function save() {
    loading = true;
    try {
      await onSave(form);
    } finally {
      loading = false;
    }
  }

  // 6. Effects
  $effect(() => {
    console.log('Patient changed:', patient);
  });
</script>

<!-- 7. Template -->
<div class="modal">
  <h2>{patient ? 'Edit' : 'New'} Patient</h2>
  <input bind:value={form.name} placeholder="Name" />
  <input bind:value={form.phone} placeholder="Phone" />
  <button onclick={save} disabled={!isValid || loading}>
    {loading ? 'Saving...' : 'Save'}
  </button>
</div>

<!-- 8. Styles (scoped) -->
<style>
  .modal {
    padding: 2rem;
    background: white;
    border-radius: 8px;
  }
</style>
```

### Reactivity Rules (Svelte 5)

**Always create new references for objects/arrays:**

```javascript
// ❌ Won't trigger reactivity
items.push(newItem);
items[0].name = "Updated";

// ✅ Creates new reference, triggers reactivity
items = [...items, newItem];
items = items.map((item, i) => (i === 0 ? { ...item, name: "Updated" } : item));

// ✅ For objects
state = { ...state, field: "new value" };
```

### Imports

**Use `$lib` alias:**

```javascript
// ✅ Good
import * as api from "$lib/utils/api.js";
import { toastStore } from "$lib/stores/toast.svelte.js";

// ❌ Bad (relative paths)
import * as api from "../../utils/api.js";
```

**Configured in `jsconfig.json`:**

```json
{
  "compilerOptions": {
    "paths": {
      "$lib": ["./src/lib"],
      "$lib/*": ["./src/lib/*"]
    }
  }
}
```

---

## Adding New Features

### Add New Component

**1. Create Component File**

`src/lib/components/ExampleCard.svelte`:

```svelte
<script>
  let { item, onClick } = $props();
</script>

<div class="card" onclick={onClick}>
  <h3>{item.title}</h3>
  <p>{item.description}</p>
</div>

<style>
  .card {
    padding: 1rem;
    border: 1px solid #e5e7eb;
    border-radius: 0.5rem;
    cursor: pointer;
  }

  .card:hover {
    background-color: #f9fafb;
  }
</style>
```

**2. Use Component**

`src/lib/views/ExampleView.svelte`:

```svelte
<script>
  import ExampleCard from '$lib/components/ExampleCard.svelte';

  let items = $state([
    { id: 1, title: 'Item 1', description: 'Description 1' },
    { id: 2, title: 'Item 2', description: 'Description 2' },
  ]);

  function handleClick(item) {
    console.log('Clicked:', item);
  }
</script>

<div class="grid grid-cols-2 gap-4">
  {#each items as item (item.id)}
    <ExampleCard {item} onClick={() => handleClick(item)} />
  {/each}
</div>
```

### Add New Store

**1. Create Store File**

`src/lib/stores/example.svelte.js`:

```javascript
class ExampleStore {
  items = $state([]);
  loading = $state(false);

  // Derived state
  count = $derived(this.items.length);

  async fetchItems() {
    this.loading = true;
    try {
      const data = await api.fetchItems();
      this.items = data;
    } catch (err) {
      console.error("Failed to fetch items:", err);
    } finally {
      this.loading = false;
    }
  }

  addItem(item) {
    this.items = [...this.items, item];
  }

  removeItem(id) {
    this.items = this.items.filter((item) => item.id !== id);
  }
}

export const exampleStore = new ExampleStore();
```

**2. Use Store**

```svelte
<script>
  import { exampleStore } from '$lib/stores/example.svelte.js';

  $effect(() => {
    exampleStore.fetchItems();
  });
</script>

<p>Total items: {exampleStore.count}</p>

{#if exampleStore.loading}
  <p>Loading...</p>
{:else}
  <ul>
    {#each exampleStore.items as item (item.id)}
      <li>{item.name}</li>
    {/each}
  </ul>
{/if}
```

### Add New View (Page)

**1. Create View Component**

`src/lib/views/ExampleView.svelte`:

```svelte
<script>
  import * as api from '$lib/utils/api.js';

  let { token } = $props();
  let data = $state([]);
  let loading = $state(true);

  async function loadData() {
    loading = true;
    try {
      data = await api.fetchExample(token);
    } catch (err) {
      console.error('Failed to load data:', err);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    loadData();
  });
</script>

<div class="container mx-auto px-4 py-8">
  <h1 class="text-3xl font-bold mb-6">Example View</h1>

  {#if loading}
    <p>Loading...</p>
  {:else}
    <div class="grid gap-4">
      {#each data as item (item.id)}
        <div class="card">
          <h2>{item.name}</h2>
        </div>
      {/each}
    </div>
  {/if}
</div>
```

**2. Add Route in App.svelte**

```svelte
<script>
  import ExampleView from '$lib/views/ExampleView.svelte';

  let currentView = $state('dashboard');

  function setView(view) {
    currentView = view;
  }
</script>

<!-- Add to navigation -->
<Sidebar {currentView} onNavigate={setView} />

<!-- Add to view switching -->
{#if currentView === 'example'}
  <ExampleView {token} />
{:else if currentView === 'dashboard'}
  <DashboardView {token} />
{/if}
```

**3. Update Sidebar Navigation**

`src/lib/components/Sidebar.svelte`:

```svelte
<button
  class="nav-item"
  class:active={currentView === 'example'}
  onclick={() => onNavigate('example')}
>
  Example
</button>
```

---

## Internationalization (i18n)

### Add New Translation Key

**1. Update Translation Files**

`src/locales/en.json`:

```json
{
  "example": {
    "title": "Example Page",
    "action": "Click Me",
    "message": "Hello, {name}!"
  }
}
```

`src/locales/id.json`:

```json
{
  "example": {
    "title": "Halaman Contoh",
    "action": "Klik Saya",
    "message": "Halo, {name}!"
  }
}
```

**2. Use in Component**

```svelte
<script>
  import { _ } from 'svelte-i18n';

  let name = $state('John');
</script>

<h1>{$_('example.title')}</h1>
<button>{$_('example.action')}</button>
<p>{$_('example.message', { values: { name } })}</p>
```

### Change Language

```svelte
<script>
  import { locale } from 'svelte-i18n';

  function setLanguage(lang) {
    locale.set(lang);
    localStorage.setItem('locale', lang);
  }
</script>

<button onclick={() => setLanguage('en')}>English</button>
<button onclick={() => setLanguage('id')}>Bahasa</button>
```

---

## Styling with Tailwind CSS

### Tailwind Configuration

**`app.css`:**

```css
@import "tailwindcss";

@theme {
  --color-primary: #3b82f6;
  --color-secondary: #10b981;
  --color-danger: #ef4444;
}

.btn {
  @apply px-4 py-2 rounded-lg font-medium transition-colors;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700;
}
```

### Use Tailwind Classes

```svelte
<div class="container mx-auto px-4 py-8">
  <h1 class="text-3xl font-bold text-gray-900 mb-6">
    Dashboard
  </h1>

  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <div class="bg-white rounded-lg shadow-md p-6">
      <h3 class="text-lg font-semibold text-gray-700">Card Title</h3>
      <p class="text-3xl font-bold text-blue-600 mt-2">123</p>
    </div>
  </div>

  <button class="btn btn-primary mt-4">
    Save Changes
  </button>
</div>
```

### Conditional Classes

```svelte
<script>
  let active = $state(false);
  let type = $state('success');
</script>

<!-- Using class: directive -->
<div class="card" class:active>Content</div>

<!-- Using ternary -->
<div class="badge {type === 'success' ? 'bg-green-500' : 'bg-red-500'}">
  Status
</div>

<!-- Multiple conditions -->
<button
  class="btn"
  class:btn-primary={type === 'primary'}
  class:btn-danger={type === 'danger'}
  class:disabled={loading}
>
  Button
</button>
```

---

## Debugging

### Vue/React DevTools Alternative

**Svelte DevTools:** Install browser extension

- Chrome: https://chrome.google.com/webstore (search "Svelte DevTools")
- Firefox: https://addons.mozilla.org/firefox/ (search "Svelte DevTools")

**Features:**

- Inspect component hierarchy
- View component state/props
- Monitor store values
- Track events

### Console Debugging

```svelte
<script>
  let data = $state({ name: 'John', age: 30 });

  $effect(() => {
    console.log('Data changed:', data);
    console.table(data);  // Table format
  });

  function handleClick() {
    console.group('Click Handler');
    console.log('Before:', data);
    data = { ...data, age: data.age + 1 };
    console.log('After:', data);
    console.groupEnd();
  }
</script>
```

### Network Debugging

**Check API calls:**

1. Open DevTools (F12)
2. Go to Network tab
3. Filter by "Fetch/XHR"
4. Click request to see details

**Check SSE Connection:**

1. Network tab → Filter by "EventSource"
2. Should see `/api/sse/delivery-status?token=...`
3. Click to see events stream

---

## Common Issues

### Port 5173 Already in Use

**Solution:**

```bash
# Kill process using port 5173
lsof -i :5173  # macOS/Linux
netstat -ano | findstr :5173  # Windows

# Or change port in vite.config.js
server: {
  port: 5174,
}
```

### Backend API Errors (CORS)

**Error:** `CORS policy: No 'Access-Control-Allow-Origin' header`

**Check:**

1. Backend `config.yaml` has correct CORS origin:
   ```yaml
   server:
     cors_origin: "http://localhost:5173"
   ```
2. Restart backend after config change

### Reactivity Not Working

**Problem:** UI doesn't update after data change

**Solutions:**

```javascript
// ❌ Mutating directly (won't work)
items.push(newItem);
state.field = "new value";

// ✅ Create new reference
items = [...items, newItem];
state = { ...state, field: "new value" };

// ✅ For arrays
items = items.filter((item) => item.id !== id);
items = items.map((item) =>
  item.id === id ? { ...item, field: "updated" } : item
);
```

### SSE Connection Fails

**Error:** EventSource error / disconnected

**Check:**

1. Backend running: `curl http://localhost:8080/api/health`
2. Valid token in localStorage
3. Token not expired (7-day expiry)
4. Check browser console for errors

**Force Reconnect:**

```javascript
// In browser console
deliveryStore.disconnect();
deliveryStore.connect();
```

---

## Performance Tips

### Lazy Loading Components

**Not implemented yet, but can add:**

```svelte
<script>
  import { onMount } from 'svelte';

  let HeavyComponent;

  onMount(async () => {
    const module = await import('$lib/components/HeavyComponent.svelte');
    HeavyComponent = module.default;
  });
</script>

{#if HeavyComponent}
  <svelte:component this={HeavyComponent} />
{/if}
```

### Optimize Large Lists

**Use `{#each}` with key:**

```svelte
<!-- ✅ With key (efficient updates) -->
{#each items as item (item.id)}
  <ItemCard {item} />
{/each}

<!-- ❌ Without key (full re-render) -->
{#each items as item}
  <ItemCard {item} />
{/each}
```

### Avoid Expensive Computations in Template

```svelte
<!-- ❌ Bad (recalculates on every render) -->
<p>{expensiveComputation(data)}</p>

<!-- ✅ Good (cached with $derived) -->
<script>
  let result = $derived(expensiveComputation(data));
</script>
<p>{result}</p>
```

---

## Production Build

### Build Command

```bash
bun run build
```

**Output:** `dist/` folder with optimized assets

**Build Size:**

- Typical: 200-500 KB (gzipped)
- Includes Svelte runtime, components, Tailwind CSS

### Preview Build

```bash
bun run preview
```

**Serves production build locally for testing.**

### Deployment

See [Deployment Guide](./deployment-guide.md) for full instructions.

**Quick:**

1. Build frontend: `bun run build`
2. Copy `dist/` folder to web server
3. Configure Nginx/Apache to serve `index.html` for all routes (SPA)

---

## Production Checklist

- [ ] Update `API_URL` to production backend URL
- [ ] Test production build locally (`bun run preview`)
- [ ] Verify all API calls work with production backend
- [ ] Check SSE connection works (may need wss:// for HTTPS)
- [ ] Enable gzip/brotli compression on web server
- [ ] Set up CDN for static assets (optional)
- [ ] Configure caching headers (1 year for assets, no-cache for index.html)
- [ ] Test all routes/views
- [ ] Test in multiple browsers (Chrome, Firefox, Safari, Edge)
- [ ] Test mobile responsiveness

---

**Next:** See [Component Inventory](./component-inventory-frontend.md) for detailed component documentation and [Deployment Guide](./deployment-guide.md) for production setup.
