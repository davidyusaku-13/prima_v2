# Architecture - Frontend (Svelte 5 + Vite)

**Generated:** January 2, 2026  
**Project:** PRIMA Healthcare Volunteer Dashboard  
**Technology:** Svelte 5.43.8 + Vite 7.2.4 (**NOT SvelteKit**)

---

## Table of Contents

1. [Overview](#overview)
2. [Component Architecture](#component-architecture)
3. [Project Structure](#project-structure)
4. [State Management](#state-management)
5. [Routing Strategy](#routing-strategy)
6. [API Client Layer](#api-client-layer)
7. [Real-Time Updates (SSE)](#real-time-updates-sse)
8. [Internationalization](#internationalization)
9. [Styling System](#styling-system)
10. [Build & Bundle Strategy](#build--bundle-strategy)

---

## Overview

### Architecture Style

**Single-Page Application (SPA)** using component-based architecture:

```
┌─────────────────────────────────────────────────────┐
│              Browser (Index.html)                   │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│               App.svelte (Root)                     │
│         • Router (view switching)                   │
│         • Auth state management                     │
│         • Global modals                             │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│              Views (Pages)                          │
│  • DashboardView, PatientsView, etc.                │
│  • Route-level components                           │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│           Reusable Components                       │
│  • PatientModal, ReminderCard, etc.                 │
│  • Leaf components                                  │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│         Services & Stores                           │
│  • API client ($lib/utils/api.js)                   │
│  • SSE service ($lib/services/sse.js)               │
│  • Svelte stores ($lib/stores/*.svelte.js)          │
└─────────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Component Composition** - Build complex UIs from small, reusable components
2. **Unidirectional Data Flow** - Props down, events up
3. **Svelte 5 Runes** - Modern reactive primitives (`$state`, `$derived`, `$effect`)
4. **No SvelteKit** - Pure Vite + Svelte (no SSR, no file-based routing)
5. **Colocation** - Components, tests, and related files live together

---

## Component Architecture

### Component Hierarchy

```
App.svelte (856 lines) ◄──── Root application component
├── Sidebar.svelte
├── BottomNav.svelte
├── Views/ (Route-level components)
│   ├── LoginScreen.svelte
│   ├── DashboardView.svelte
│   ├── PatientsView.svelte
│   ├── UsersView.svelte
│   ├── BeritaView.svelte (Articles)
│   ├── BeritaDetailView.svelte
│   ├── VideoEdukasiView.svelte (Videos)
│   ├── CMSDashboardView.svelte
│   ├── ArticleEditorView.svelte
│   ├── VideoManagerView.svelte
│   ├── analytics/
│   │   └── FailedDeliveriesView.svelte
│   └── cms/
│       └── CmsAnalyticsView.svelte
├── Modals/ (Global modal components)
│   ├── PatientModal.svelte
│   ├── ReminderModal.svelte
│   ├── UserModal.svelte
│   ├── ProfileModal.svelte
│   ├── VideoModal.svelte
│   ├── VideoEditModal.svelte
│   ├── SendReminderModal.svelte
│   ├── PhoneEditModal.svelte
│   └── ConfirmModal.svelte
├── Feature Components/
│   ├── DashboardStats.svelte
│   ├── ActivityLog.svelte
│   ├── ArticleCard.svelte
│   ├── VideoCard.svelte
│   ├── ContentListItem.svelte
│   ├── ImageUploader.svelte
│   ├── QuillEditor.svelte
│   ├── patients/ (Patient-specific components)
│   ├── reminders/ (Reminder-specific components)
│   ├── delivery/ (Delivery status components)
│   ├── analytics/ (Analytics components)
│   ├── content/ (CMS content components)
│   ├── health/ (Health check components)
│   ├── indicators/ (Badge/indicator components)
│   │   └── FailedReminderBadge.svelte
│   ├── whatsapp/ (WhatsApp-related components)
│   └── ui/ (Reusable UI components)
│       └── Toast.svelte
└── Stores/ (Global state management)
    ├── auth.js (Legacy, using localStorage)
    ├── toast.svelte.js (Svelte 5 runes)
    └── delivery.svelte.js (Svelte 5 runes + SSE)
```

### Component Patterns

#### 1. Container/Presentational Pattern

**Container (Smart Component):**

- Manages state
- Fetches data from API
- Handles business logic
- Passes data via props

```svelte
<!-- DashboardView.svelte (container) -->
<script>
  import DashboardStats from '$lib/components/DashboardStats.svelte';
  import * as api from '$lib/utils/api.js';

  let stats = $state({ patients: 0, reminders: 0, volunteers: 0 });
  let loading = $state(true);

  async function loadStats() {
    loading = true;
    try {
      const data = await api.fetchDashboardStats(token);
      stats = data;
    } catch (err) {
      console.error('Failed to load stats:', err);
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    loadStats();
  });
</script>

<DashboardStats {stats} {loading} />
```

**Presentational (Dumb Component):**

- Receives data via props
- No API calls
- Pure UI rendering
- Emits events for actions

```svelte
<!-- DashboardStats.svelte (presentational) -->
<script>
  let { stats = { patients: 0, reminders: 0, volunteers: 0 } } = $props();
</script>

<div class="grid grid-cols-3 gap-4">
  <div class="stat-card">
    <h3>Patients</h3>
    <p>{stats.patients}</p>
  </div>
  <div class="stat-card">
    <h3>Reminders</h3>
    <p>{stats.reminders}</p>
  </div>
  <div class="stat-card">
    <h3>Volunteers</h3>
    <p>{stats.volunteers}</p>
  </div>
</div>
```

#### 2. Modal Pattern

**Global Modal Management (in App.svelte):**

```svelte
<script>
  let showPatientModal = $state(false);
  let editingPatient = $state(null);
  let patientForm = $state({ name: '', phone: '', email: '', notes: '' });

  function openPatientModal(patient = null) {
    editingPatient = patient;
    patientForm = patient ? { ...patient } : { name: '', phone: '', email: '', notes: '' };
    showPatientModal = true;
  }

  function closePatientModal() {
    showPatientModal = false;
    editingPatient = null;
  }

  async function savePatient() {
    try {
      await api.savePatient(token, patientForm, editingPatient?.id);
      await loadPatients(); // Refresh list
      closePatientModal();
    } catch (err) {
      console.error('Failed to save patient:', err);
    }
  }
</script>

<PatientModal
  show={showPatientModal}
  editingPatient={editingPatient}
  patientForm={patientForm}
  onClose={closePatientModal}
  onSave={savePatient}
/>
```

**Modal Component:**

```svelte
<!-- PatientModal.svelte -->
<script>
  let { show = false, editingPatient = null, patientForm, onClose, onSave } = $props();
</script>

{#if show}
  <div class="modal-backdrop" onclick={onClose}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <h2>{editingPatient ? 'Edit Patient' : 'New Patient'}</h2>

      <input type="text" bind:value={patientForm.name} placeholder="Name" />
      <input type="tel" bind:value={patientForm.phone} placeholder="Phone" />
      <input type="email" bind:value={patientForm.email} placeholder="Email" />

      <div class="modal-actions">
        <button onclick={onClose}>Cancel</button>
        <button onclick={onSave}>Save</button>
      </div>
    </div>
  </div>
{/if}
```

#### 3. Svelte 5 Runes Pattern

**State Management:**

```svelte
<script>
  // Reactive state (replaces let with $: reactive statements)
  let count = $state(0);

  // Derived state (replaces $: count * 2)
  let doubled = $derived(count * 2);

  // Effects (replaces $: { console.log() })
  $effect(() => {
    console.log('Count changed:', count);
  });

  // Props (new syntax)
  let { title = 'Default Title', onUpdate } = $props();
</script>

<h1>{title}</h1>
<p>Count: {count}, Doubled: {doubled}</p>
<button onclick={() => count++}>Increment</button>
```

**Comparison with Legacy Svelte:**

| Legacy Svelte 3/4               | Svelte 5                                  |
| ------------------------------- | ----------------------------------------- |
| `export let title = 'Default';` | `let { title = 'Default' } = $props();`   |
| `let count = 0;`                | `let count = $state(0);`                  |
| `$: doubled = count * 2;`       | `let doubled = $derived(count * 2);`      |
| `$: console.log(count);`        | `$effect(() => { console.log(count); });` |
| `on:click={handler}`            | `onclick={handler}`                       |

---

## Project Structure

```
frontend/
├── src/
│   ├── main.js                  # Entry point (mounts App.svelte)
│   ├── app.css                  # Global styles (Tailwind)
│   ├── i18n.js                  # i18n setup (svelte-i18n)
│   ├── App.svelte               # Root component (856 lines)
│   ├── assets/                  # Static assets (images, icons)
│   ├── lib/
│   │   ├── components/          # Reusable Svelte components
│   │   │   ├── ActivityLog.svelte
│   │   │   ├── ArticleCard.svelte
│   │   │   ├── BottomNav.svelte
│   │   │   ├── ConfirmModal.svelte
│   │   │   ├── ContentListItem.svelte
│   │   │   ├── DashboardStats.svelte
│   │   │   ├── ImageUploader.svelte
│   │   │   ├── PatientModal.svelte
│   │   │   ├── ProfileModal.svelte
│   │   │   ├── QuillEditor.svelte
│   │   │   ├── ReminderModal.svelte
│   │   │   ├── SendReminderModal.svelte
│   │   │   ├── Sidebar.svelte
│   │   │   ├── UserModal.svelte
│   │   │   ├── VideoCard.svelte
│   │   │   ├── VideoEditModal.svelte
│   │   │   ├── VideoModal.svelte
│   │   │   ├── PhoneEditModal.svelte
│   │   │   ├── analytics/       # Analytics components
│   │   │   ├── content/         # CMS content components
│   │   │   ├── delivery/        # Delivery status components
│   │   │   ├── health/          # Health check components
│   │   │   ├── indicators/      # Badges, status indicators
│   │   │   │   └── FailedReminderBadge.svelte
│   │   │   ├── patients/        # Patient-specific components
│   │   │   ├── reminders/       # Reminder-specific components
│   │   │   ├── ui/              # Generic UI components
│   │   │   │   └── Toast.svelte
│   │   │   └── whatsapp/        # WhatsApp integration components
│   │   ├── views/               # Page-level components
│   │   │   ├── LoginScreen.svelte
│   │   │   ├── DashboardView.svelte
│   │   │   ├── PatientsView.svelte
│   │   │   ├── UsersView.svelte
│   │   │   ├── BeritaView.svelte
│   │   │   ├── BeritaDetailView.svelte
│   │   │   ├── VideoEdukasiView.svelte
│   │   │   ├── CMSDashboardView.svelte
│   │   │   ├── ArticleEditorView.svelte
│   │   │   ├── VideoManagerView.svelte
│   │   │   ├── analytics/
│   │   │   │   └── FailedDeliveriesView.svelte
│   │   │   ├── cms/
│   │   │   │   └── CmsAnalyticsView.svelte
│   │   │   └── patients/        # Patient detail views
│   │   ├── stores/              # Svelte stores (global state)
│   │   │   ├── auth.js          # Auth store (legacy)
│   │   │   ├── toast.svelte.js  # Toast notifications (Svelte 5)
│   │   │   └── delivery.svelte.js # Delivery status (Svelte 5 + SSE)
│   │   ├── services/            # External service clients
│   │   │   └── sse.js           # Server-Sent Events client
│   │   ├── utils/               # Utility functions
│   │   │   ├── api.js           # Backend API client (539 lines)
│   │   │   └── ...
│   │   └── test/                # Test files
│   │       └── ...
│   └── locales/                 # i18n translation files
│       ├── en.json              # English translations
│       └── id.json              # Indonesian translations
├── public/                      # Static files (served as-is)
├── index.html                   # HTML entry point
├── vite.config.js               # Vite build configuration
├── svelte.config.js             # Svelte compiler options
├── jsconfig.json                # JavaScript path mappings ($lib alias)
├── vitest.config.js             # Test runner configuration
├── package.json                 # Dependencies
└── README.md                    # Project documentation
```

---

## State Management

### 1. Local Component State (Svelte 5 Runes)

**Usage:** Component-specific state that doesn't need to be shared

```svelte
<script>
  let count = $state(0);
  let name = $state('');
  let items = $state([]);

  function addItem(item) {
    items = [...items, item]; // Create new array for reactivity
  }
</script>
```

**CRITICAL for Svelte 5 Reactivity:**

- Always create new object/array references when mutating
- ❌ `items.push(newItem)` - Won't trigger reactivity
- ✅ `items = [...items, newItem]` - Creates new reference, triggers reactivity

### 2. Svelte Stores (Global State)

**Location:** `src/lib/stores/*.svelte.js`

#### Toast Store (Svelte 5 Runes)

```javascript
// stores/toast.svelte.js
class ToastStore {
  toasts = $state([]);

  add(message, options = {}) {
    const id = Date.now() + Math.random();
    const toast = {
      id,
      message,
      type: options.type || "info", // info, success, error, warning
      duration: options.duration || 3000,
      action: options.action || null, // { label: 'Undo', onClick: () => {} }
    };

    this.toasts = [...this.toasts, toast];

    if (toast.duration > 0) {
      setTimeout(() => this.remove(id), toast.duration);
    }
  }

  remove(id) {
    this.toasts = this.toasts.filter((t) => t.id !== id);
  }
}

export const toastStore = new ToastStore();
```

**Usage:**

```svelte
<script>
  import { toastStore } from '$lib/stores/toast.svelte.js';

  function saveData() {
    try {
      // ... save logic
      toastStore.add('Data saved successfully', { type: 'success' });
    } catch (err) {
      toastStore.add('Failed to save data', { type: 'error' });
    }
  }
</script>
```

#### Delivery Store (SSE Integration)

```javascript
// stores/delivery.svelte.js
import { sseService } from "$lib/services/sse.js";

class DeliveryStore {
  deliveryStatuses = $state({});
  connectionStatus = $state("disconnected");
  failedReminders = $state([]);

  failedCount = $derived(this.failedReminders.length);

  constructor() {
    // Subscribe to SSE events
    sseService.on("delivery.status.updated", (data) => {
      this.updateStatus(data.reminder_id, data.status, data.timestamp);
    });

    sseService.on("connection.status", (status) => {
      this.connectionStatus = status;
    });

    sseService.on("delivery.failed", (data) => {
      this.addFailedReminder(data);
      toastStore.add(`Delivery failed: ${data.patient_name}`, {
        type: "error",
        action: {
          label: "View Details",
          onClick: () => {
            window.dispatchEvent(
              new CustomEvent("navigate-to-patient", {
                detail: { patientId: data.patient_id },
              })
            );
          },
        },
      });
    });
  }

  updateStatus(reminderId, status, timestamp) {
    this.deliveryStatuses = {
      ...this.deliveryStatuses,
      [reminderId]: { status, timestamp, updatedAt: new Date().toISOString() },
    };
  }

  getStatus(reminderId) {
    return this.deliveryStatuses[reminderId]?.status || null;
  }

  connect() {
    sseService.connect();
  }

  disconnect() {
    sseService.disconnect();
  }
}

export const deliveryStore = new DeliveryStore();
```

### 3. LocalStorage (Auth Token)

**Simple persistence without reactivity:**

```javascript
// Set token
localStorage.setItem("token", data.token);

// Get token
const token = localStorage.getItem("token");

// Remove token
localStorage.removeItem("token");
```

**In App.svelte:**

```svelte
<script>
  let token = $state(localStorage.getItem('token') || null);

  function login(username, password) {
    const data = await api.login(username, password);
    token = data.token;
    localStorage.setItem('token', token);
  }

  function logout() {
    token = null;
    localStorage.removeItem('token');
  }
</script>
```

---

## Routing Strategy

### Manual Routing (No Router Library)

**Why No Router?**

- Small SPA with limited routes
- Avoid SvelteKit dependency (this is Vite + Svelte only)
- Simple state-based routing sufficient

**Implementation in App.svelte:**

```svelte
<script>
  let currentView = $state(localStorage.getItem('currentView') || 'dashboard');

  function setView(view) {
    currentView = view;
    localStorage.setItem('currentView', view); // Persist across page reloads
  }
</script>

<!-- Sidebar navigation -->
<Sidebar currentView={currentView} onNavigate={setView} />

<!-- View switching -->
{#if !token}
  <LoginScreen />
{:else if currentView === 'dashboard'}
  <DashboardView {token} {user} />
{:else if currentView === 'patients'}
  <PatientsView {token} {user} />
{:else if currentView === 'users'}
  <UsersView {token} />
{:else if currentView === 'berita'}
  <BeritaView {token} {user} />
{:else if currentView === 'berita-detail'}
  <BeritaDetailView {token} {articleSlug} onBack={() => setView('berita')} />
{:else if currentView === 'video-edukasi'}
  <VideoEdukasiView {token} {user} />
{:else if currentView === 'cms-dashboard'}
  <CMSDashboardView {token} {user} />
{:else if currentView === 'article-editor'}
  <ArticleEditorView {token} onBack={() => setView('cms-dashboard')} />
{:else if currentView === 'video-manager'}
  <VideoManagerView {token} onBack={() => setView('cms-dashboard')} />
{:else if currentView === 'failed-deliveries'}
  <FailedDeliveriesView {token} />
{:else if currentView === 'cms-analytics'}
  <CmsAnalyticsView {token} />
{/if}
```

**Navigation with Context:**

```javascript
// Navigate to article detail
function openArticle(article) {
  articleSlug = article.slug;
  setView("berita-detail");
}

// Navigate back
function goBack() {
  setView("berita");
}
```

---

## API Client Layer

### Centralized API Module

**Location:** `src/lib/utils/api.js` (539 lines)

**Pattern:**

```javascript
const API_URL = "http://localhost:8080/api";

function getHeaders(token) {
  const headers = { "Content-Type": "application/json" };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  return headers;
}

// Generic fetch wrapper
async function apiCall(endpoint, options = {}, token = null) {
  const res = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: getHeaders(token),
  });

  const data = await res.json();
  if (!res.ok) {
    throw new Error(data.error || "API call failed");
  }

  return data;
}

// Auth
export async function login(username, password) {
  const res = await fetch(`${API_URL}/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || "Login failed");
  return data;
}

export async function fetchUser(token) {
  const res = await fetch(`${API_URL}/auth/me`, {
    headers: getHeaders(token),
  });
  if (!res.ok) throw new Error("Unauthorized");
  return res.json();
}

// Patients
export async function fetchPatients(token) {
  const res = await fetch(`${API_URL}/patients`, {
    headers: getHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to fetch patients");
  const data = await res.json();
  return data.patients || [];
}

export async function savePatient(token, patient, editingId = null) {
  const method = editingId ? "PUT" : "POST";
  const url = editingId
    ? `${API_URL}/patients/${editingId}`
    : `${API_URL}/patients`;
  const res = await fetch(url, {
    method,
    headers: getHeaders(token),
    body: JSON.stringify(patient),
  });
  if (!res.ok) throw new Error("Failed to save patient");
  return res.json();
}

// Reminders
export async function sendReminder(token, patientId, reminderId) {
  const res = await fetch(
    `${API_URL}/patients/${patientId}/reminders/${reminderId}/send`,
    {
      method: "POST",
      headers: getHeaders(token),
    }
  );
  if (!res.ok) throw new Error("Failed to send reminder");
  return res.json();
}

// CMS - Articles
export async function fetchArticles(token) {
  const res = await fetch(`${API_URL}/articles`, {
    headers: getHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to fetch articles");
  const data = await res.json();
  return data.articles || [];
}

export async function createArticle(token, article) {
  const res = await fetch(`${API_URL}/articles`, {
    method: "POST",
    headers: getHeaders(token),
    body: JSON.stringify(article),
  });
  if (!res.ok) throw new Error("Failed to create article");
  return res.json();
}

// CMS - Videos
export async function createVideo(token, video) {
  const res = await fetch(`${API_URL}/videos`, {
    method: "POST",
    headers: getHeaders(token),
    body: JSON.stringify(video),
  });
  if (!res.ok) throw new Error("Failed to create video");
  return res.json();
}

// Analytics
export async function fetchDeliveryAnalytics(token) {
  const res = await fetch(`${API_URL}/analytics/delivery`, {
    headers: getHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to fetch analytics");
  return res.json();
}

export async function fetchFailedDeliveries(token, filters = {}) {
  const params = new URLSearchParams(filters);
  const res = await fetch(`${API_URL}/analytics/failed-deliveries?${params}`, {
    headers: getHeaders(token),
  });
  if (!res.ok) throw new Error("Failed to fetch failed deliveries");
  return res.json();
}
```

**Usage in Components:**

```svelte
<script>
  import * as api from '$lib/utils/api.js';

  let patients = $state([]);
  let loading = $state(true);

  async function loadPatients() {
    loading = true;
    try {
      patients = await api.fetchPatients(token);
    } catch (err) {
      console.error('Failed to load patients:', err);
      toastStore.add('Failed to load patients', { type: 'error' });
    } finally {
      loading = false;
    }
  }

  $effect(() => {
    loadPatients();
  });
</script>
```

---

## Real-Time Updates (SSE)

### SSE Service

**Location:** `src/lib/services/sse.js`

**Implementation:**

```javascript
class SSEService {
  eventSource = null;
  listeners = {};
  token = null;

  connect() {
    if (this.eventSource) {
      this.disconnect();
    }

    this.token = localStorage.getItem("token");
    if (!this.token) {
      console.error("No token available for SSE connection");
      return;
    }

    // EventSource doesn't support custom headers, use query parameter
    this.eventSource = new EventSource(
      `http://localhost:8080/api/sse/delivery-status?token=${this.token}`
    );

    this.eventSource.addEventListener("connection.established", (e) => {
      console.log("SSE connected:", JSON.parse(e.data));
      this.emit("connection.status", "connected");
    });

    this.eventSource.addEventListener("delivery.status.updated", (e) => {
      const data = JSON.parse(e.data);
      this.emit("delivery.status.updated", data);
    });

    this.eventSource.addEventListener("delivery.failed", (e) => {
      const data = JSON.parse(e.data);
      this.emit("delivery.failed", data);
    });

    this.eventSource.onerror = () => {
      console.error("SSE connection error");
      this.emit("connection.status", "error");
      // Auto-reconnect after 5 seconds
      setTimeout(() => this.connect(), 5000);
    };
  }

  disconnect() {
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
      this.emit("connection.status", "disconnected");
    }
  }

  on(eventType, callback) {
    if (!this.listeners[eventType]) {
      this.listeners[eventType] = [];
    }
    this.listeners[eventType].push(callback);
  }

  emit(eventType, data) {
    if (this.listeners[eventType]) {
      this.listeners[eventType].forEach((callback) => callback(data));
    }
  }
}

export const sseService = new SSEService();
```

**Integration with Delivery Store:**

```javascript
// stores/delivery.svelte.js
constructor() {
    sseService.on('delivery.status.updated', (data) => {
        this.updateStatus(data.reminder_id, data.status, data.timestamp);
    });
}
```

**Lifecycle Management in App.svelte:**

```svelte
<script>
  import { deliveryStore } from '$lib/stores/delivery.svelte.js';

  $effect(() => {
    if (token) {
      deliveryStore.connect(); // Connect when logged in
    }

    return () => {
      deliveryStore.disconnect(); // Cleanup on logout
    };
  });
</script>
```

---

## Internationalization

### svelte-i18n Setup

**Installation:**

```bash
bun add svelte-i18n
```

**Configuration (`src/i18n.js`):**

```javascript
import { register, init, locale } from "svelte-i18n";

register("en", () => import("./locales/en.json"));
register("id", () => import("./locales/id.json"));

init({
  fallbackLocale: "en",
  initialLocale: localStorage.getItem("locale") || "en",
});

// Export locale for setting
export { locale };
```

**Translation Files:**

`src/locales/en.json`:

```json
{
  "common": {
    "save": "Save",
    "cancel": "Cancel",
    "delete": "Delete",
    "edit": "Edit",
    "loading": "Loading..."
  },
  "dashboard": {
    "title": "Dashboard",
    "stats": {
      "patients": "Patients",
      "reminders": "Reminders",
      "volunteers": "Volunteers"
    }
  },
  "reminder": {
    "send": "Send Reminder",
    "sendSuccess": "Reminder sent successfully",
    "sendFailed": "Failed to send reminder",
    "failedNotification": "Delivery failed for {patientName}",
    "viewDetails": "View Details"
  }
}
```

`src/locales/id.json`:

```json
{
  "common": {
    "save": "Simpan",
    "cancel": "Batal",
    "delete": "Hapus",
    "edit": "Edit",
    "loading": "Memuat..."
  },
  "dashboard": {
    "title": "Dasbor",
    "stats": {
      "patients": "Pasien",
      "reminders": "Pengingat",
      "volunteers": "Relawan"
    }
  },
  "reminder": {
    "send": "Kirim Pengingat",
    "sendSuccess": "Pengingat berhasil dikirim",
    "sendFailed": "Gagal mengirim pengingat",
    "failedNotification": "Pengiriman gagal untuk {patientName}",
    "viewDetails": "Lihat Detail"
  }
}
```

**Usage in Components:**

```svelte
<script>
  import { _, locale } from 'svelte-i18n';

  function changeLanguage(newLocale) {
    locale.set(newLocale);
    localStorage.setItem('locale', newLocale);
  }
</script>

<h1>{$_('dashboard.title')}</h1>
<button onclick={() => changeLanguage('en')}>English</button>
<button onclick={() => changeLanguage('id')}>Bahasa Indonesia</button>

<button onclick={sendReminder}>
  {$_('reminder.send')}
</button>

<p>{$_('reminder.failedNotification', { values: { patientName: patient.name } })}</p>
```

---

## Styling System

### Tailwind CSS 4

**Configuration (`app.css`):**

```css
@import "tailwindcss";

/* Custom theme variables */
@theme {
  --color-primary: #3b82f6;
  --color-secondary: #10b981;
  --color-danger: #ef4444;
}

/* Global styles */
body {
  @apply bg-gray-50 text-gray-900;
}

.btn {
  @apply px-4 py-2 rounded-lg font-medium transition-colors;
}

.btn-primary {
  @apply bg-blue-600 text-white hover:bg-blue-700;
}

.btn-secondary {
  @apply bg-gray-200 text-gray-800 hover:bg-gray-300;
}

.btn-danger {
  @apply bg-red-600 text-white hover:bg-red-700;
}
```

**Usage in Components:**

```svelte
<div class="max-w-7xl mx-auto px-4 py-8">
  <h1 class="text-3xl font-bold mb-6">Dashboard</h1>

  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <div class="bg-white rounded-lg shadow p-6">
      <h3 class="text-lg font-semibold">Patients</h3>
      <p class="text-3xl font-bold text-blue-600">{stats.patients}</p>
    </div>
  </div>

  <button class="btn btn-primary mt-4">
    Add Patient
  </button>
</div>
```

---

## Build & Bundle Strategy

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
  build: {
    outDir: "dist",
    sourcemap: true,
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

### Build Commands

**Development:**

```bash
bun run dev
# Starts dev server on http://localhost:5173
# Hot module replacement enabled
# Proxy /api to backend (http://localhost:8080)
```

**Production Build:**

```bash
bun run build
# Output: frontend/dist/
# Optimized bundle with minification
# Generates source maps
```

**Preview Production Build:**

```bash
bun run preview
# Serves production build locally for testing
```

### Bundle Optimization

**Code Splitting:**

- Vite automatically splits vendor code (node_modules) into separate chunk
- Dynamic imports create separate chunks (not used extensively in PRIMA)

**Tree Shaking:**

- Unused exports from modules are eliminated
- Tailwind CSS purges unused classes

**Minification:**

- JavaScript: Terser (enabled by default in production)
- CSS: cssnano (built into Vite)

---

**Next:** See [Component Inventory](./component-inventory-frontend.md) for detailed component documentation.
