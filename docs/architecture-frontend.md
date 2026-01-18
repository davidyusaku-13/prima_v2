# Architecture: Frontend (Svelte 5 + Vite)

> Generated: 2026-01-18 | Scan Level: Exhaustive

## Executive Summary

The frontend is a Svelte 5 + Vite single-page application providing a healthcare volunteer dashboard UI. It uses Svelte 5 runes for reactivity, Tailwind CSS 4 for styling, and communicates with the Go backend via REST API and SSE for real-time updates.

**CRITICAL**: This is Vite + Svelte 5, NOT SvelteKit. No SvelteKit imports allowed.

## Technology Stack

| Category | Technology | Version | Purpose |
|----------|------------|---------|---------|
| Framework | Svelte | 5.43.8 | Reactive UI components |
| Build Tool | Vite | 7.2.4 | Development & bundling |
| CSS | Tailwind CSS | 4.1.18 | Utility-first styling |
| Rich Text | Quill | 2.0.3 | Article content editing |
| i18n | svelte-i18n | 4.0.1 | Internationalization |
| Testing | Vitest | 4.0.16 | Unit testing |
| Package Manager | Bun | - | Dependency management |

## Architecture Pattern

**Component-Based SPA** with centralized state:

```
frontend/src/
├── App.svelte           # Root component, routing, global state
├── i18n.js              # Internationalization setup
├── lib/
│   ├── components/      # Reusable UI components
│   ├── views/           # Page-level views
│   ├── stores/          # Global state management
│   ├── services/        # External service integrations
│   ├── utils/           # Utility functions
│   └── test-utils/      # Testing utilities
├── locales/             # Translation files (en, id)
└── assets/              # Static assets
```

## Svelte 5 Patterns (CRITICAL)

### Runes Syntax
```javascript
// State
let count = $state(0);
let user = $state(null);

// Derived state
let doubled = $derived(count * 2);
let isLoggedIn = $derived(!!user);

// Props
let { name, age } = $props();

// Effects
$effect(() => {
  console.log('Count changed:', count);
});
```

### Event Handling
```svelte
<!-- Svelte 5 syntax (CORRECT) -->
<button onclick={handleClick}>Click</button>
<input oninput={(e) => value = e.target.value}>

<!-- NOT Svelte 4 syntax (WRONG) -->
<button on:click={handleClick}>Click</button>
```

### Immutable Updates (Required for Reactivity)
```javascript
// Arrays - create new reference
items = [...items, newItem];           // Add
items = items.filter(i => i.id !== id); // Remove

// Objects - create new reference
obj = { ...obj, key: newValue };
```

## Component Inventory

### Views (Page Components)

| Component | Path | Description |
|-----------|------|-------------|
| `LoginScreen` | `views/LoginScreen.svelte` | Authentication UI |
| `DashboardView` | `views/DashboardView.svelte` | Main dashboard with stats |
| `PatientsView` | `views/PatientsView.svelte` | Patient list & management |
| `UsersView` | `views/UsersView.svelte` | User management (superadmin) |
| `BeritaView` | `views/BeritaView.svelte` | Article listing |
| `BeritaDetailView` | `views/BeritaDetailView.svelte` | Single article view |
| `VideoEdukasiView` | `views/VideoEdukasiView.svelte` | Video listing |
| `CMSDashboardView` | `views/CMSDashboardView.svelte` | CMS dashboard |
| `ArticleEditorView` | `views/ArticleEditorView.svelte` | Article WYSIWYG editor |
| `VideoManagerView` | `views/VideoManagerView.svelte` | Video management |
| `FailedDeliveriesView` | `views/analytics/FailedDeliveriesView.svelte` | Failed delivery analytics |
| `CmsAnalyticsView` | `views/cms/CmsAnalyticsView.svelte` | Content analytics |
| `PatientDetailView` | `views/patients/PatientDetailView.svelte` | Patient detail (desktop) |
| `ReminderHistoryView` | `views/patients/ReminderHistoryView.svelte` | Reminder history |

### UI Components

| Component | Path | Description |
|-----------|------|-------------|
| `Sidebar` | `components/Sidebar.svelte` | Desktop navigation |
| `BottomNav` | `components/BottomNav.svelte` | Mobile navigation |
| `PatientModal` | `components/PatientModal.svelte` | Patient create/edit |
| `ReminderModal` | `components/ReminderModal.svelte` | Reminder create/edit |
| `SendReminderModal` | `components/SendReminderModal.svelte` | WhatsApp send confirmation |
| `ConfirmModal` | `components/ConfirmModal.svelte` | Generic confirmation |
| `Toast` | `components/ui/Toast.svelte` | Toast notifications |
| `EmptyState` | `components/ui/EmptyState.svelte` | Empty state placeholder |
| `QuillEditor` | `components/QuillEditor.svelte` | Rich text editor wrapper |
| `ImageUploader` | `components/ImageUploader.svelte` | Hero image upload |
| `ContentPickerModal` | `components/content/ContentPickerModal.svelte` | Content attachment picker |
| `DeliveryStatusBadge` | `components/delivery/DeliveryStatusBadge.svelte` | Status indicator |
| `DeliveryStatusFilter` | `components/delivery/DeliveryStatusFilter.svelte` | Status filter |
| `FailedReminderBadge` | `components/indicators/FailedReminderBadge.svelte` | Failed count indicator |
| `QuietHoursHint` | `components/indicators/QuietHoursHint.svelte` | Quiet hours info |
| `SystemHealthWidget` | `components/health/SystemHealthWidget.svelte` | System health display |

### Analytics Components

| Component | Path | Description |
|-----------|------|-------------|
| `ContentAnalyticsWidget` | `components/analytics/ContentAnalyticsWidget.svelte` | Content stats |
| `DeliveryAnalyticsWidget` | `components/analytics/DeliveryAnalyticsWidget.svelte` | Delivery stats |
| `FailedDeliveryCard` | `components/analytics/FailedDeliveryCard.svelte` | Failed delivery item |

## State Management

### Stores

#### auth.js (Svelte Store)
```javascript
// Traditional Svelte store for auth state
export const auth = createAuthStore();
// Methods: setUser, setToken, logout, reset
```

#### delivery.svelte.js (Svelte 5 Runes)
```javascript
// Svelte 5 runes-based store
class DeliveryStore {
    deliveryStatuses = $state({});
    connectionStatus = $state('disconnected');
    failedReminders = $state([]);
    failedCount = $derived(this.failedReminders.length);
}
export const deliveryStore = new DeliveryStore();
```

#### toast.svelte.js (Svelte 5 Runes)
```javascript
// Toast notification store
class ToastStore {
    toasts = $state([]);
    add(message, options) { /* ... */ }
    remove(id) { /* ... */ }
    clear() { /* ... */ }
}
export const toastStore = new ToastStore();
```

## Services

### SSE Service (`services/sse.js`)
Real-time delivery status updates via Server-Sent Events:
```javascript
sseService.connect();           // Connect to SSE endpoint
sseService.disconnect();        // Disconnect
sseService.on('event', cb);     // Subscribe to events

// Events:
// - delivery.status.updated
// - delivery.failed
// - connection.status
```

### API Utilities (`utils/api.js`)
Centralized API call functions with error handling:
- Authentication: `login()`, `register()`, `fetchUser()`
- Patients: `fetchPatients()`, `savePatient()`, `deletePatient()`
- Reminders: `saveReminder()`, `sendReminder()`, `retryReminder()`
- Content: `fetchArticles()`, `fetchVideos()`, etc.

## Navigation

### Routing (localStorage-based)
```javascript
function navigateTo(view) {
    currentView = view;
    localStorage.setItem('currentView', view);
}

// Available views:
// dashboard, patients, users, cms, berita, berita-detail, video
// failed-deliveries, analytics
```

### View Visibility Rules

| View | Required Role |
|------|---------------|
| `dashboard` | Any authenticated |
| `patients` | Any authenticated |
| `users` | superadmin only |
| `cms` | admin, superadmin |
| `berita`, `video` | Any authenticated |
| `failed-deliveries` | admin, superadmin |
| `analytics` | admin, superadmin |

## Internationalization (i18n)

### Setup
```javascript
// src/i18n.js
import { init, locale } from 'svelte-i18n';

init({
    fallbackLocale: 'id',
    initialLocale: 'id'
});
```

### Usage
```svelte
<script>
import { t } from 'svelte-i18n';
</script>

<h1>{$t('dashboard.title')}</h1>
<p>{$t('reminder.status.sent')}</p>
```

### Supported Locales
- `id` - Indonesian (default)
- `en` - English

## Development Commands

```bash
# Start dev server (port 5173)
cd frontend && bun run dev

# Run tests
cd frontend && bun run test -- --run

# Build for production
cd frontend && bun run build
```

## Key Patterns

### Props Pattern
```svelte
<script>
let {
    patient,
    onSave,
    onClose,
    loading = false
} = $props();
</script>
```

### Event Callbacks
```svelte
<!-- Parent -->
<PatientModal
    {patient}
    onSave={handleSave}
    onClose={handleClose}
/>

<!-- Child -->
<button onclick={onSave}>Save</button>
```

### Conditional Rendering
```svelte
{#if loading}
    <Spinner />
{:else if error}
    <ErrorMessage {error} />
{:else}
    <Content {data} />
{/if}
```

### List Rendering
```svelte
{#each patients as patient (patient.id)}
    <PatientCard {patient} />
{/each}
```

## Testing

### Test Files
- Co-located: `ComponentName.test.js` next to component
- Run: `bun run test -- --run`
- Framework: Vitest

### Test Utilities
Located in `src/lib/test-utils/`:
- Component mounting helpers
- Mock stores
- API mocking utilities

## Constraints

1. **NO SvelteKit imports** - Use:
   - `window.location.href` instead of `goto`
   - Props/context instead of `$app/stores`
   - `typeof window !== 'undefined'` instead of `$app/environment`

2. **Svelte 5 runes required** - No legacy reactive statements (`$:`)

3. **Immutable updates** - Always create new array/object references

4. **ARIA required** - All interactive elements need accessibility attributes

5. **No console.log** - Use `console.warn` or `console.error` only

## Entry Points

- **Root Component**: `src/App.svelte`
- **Entry HTML**: `index.html`
- **Vite Config**: `vite.config.js`
- **Svelte Config**: `svelte.config.js`
