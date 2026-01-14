# PRIMA Frontend Architecture

**Date:** 2026-01-13
**Type:** Web Frontend
**Framework:** Svelte 5 + Vite

## Executive Summary

The PRIMA frontend is a Single Page Application (SPA) built with Svelte 5 and Vite, providing the user interface for the Healthcare Volunteer Dashboard. It features state-based routing, real-time delivery status updates via SSE, and comprehensive patient and content management capabilities.

### Key Characteristics

- **Framework:** Svelte 5 with runes (`$state`, `$derived`, `$effect`)
- **Build Tool:** Vite 7.x
- **Styling:** Tailwind CSS 4
- **i18n:** svelte-i18n (English + Indonesian)
- **State Management:** Class-based stores with runes, legacy Svelte stores
- **Real-time:** Server-Sent Events (SSE) for delivery updates
- **Testing:** Vitest with Testing Library

## Technology Stack

| Category | Technology | Version | Purpose |
|----------|------------|---------|---------|
| Framework | Svelte | 5.43.8 | UI component framework |
| Build Tool | Vite | 7.2.4 | Development and build |
| Styling | Tailwind CSS | 4.1.18 | Utility-first CSS |
| i18n | svelte-i18n | 4.0.1 | Localization |
| Rich Text | Quill | 2.0.3 | Article editor |
| Testing | Vitest | 4.0.16 | Unit testing |
| Test Utils | Testing Library | 6.9.1 | Svelte testing utilities |

## Architecture Pattern

**Component-Based Architecture** with state-based routing:

```
┌─────────────────────────────────────────────────────────────┐
│                     App.svelte                               │
│              (Root, State-based Routing)                     │
├─────────────────────────────────────────────────────────────┤
│                      Views Layer                             │
│           (Page-level components with business logic)        │
├─────────────────────────────────────────────────────────────┤
│                    Components Layer                          │
│          (Reusable UI components with props)                 │
├─────────────────────────────────────────────────────────────┤
│                   Stores & Services                          │
│        (State management, SSE service, API utils)            │
├─────────────────────────────────────────────────────────────┤
│                  API Layer (api.js)                          │
│              (REST API communication)                        │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
frontend/
├── src/
│   ├── main.js              # Entry point (mount)
│   ├── App.svelte           # Root component
│   ├── app.css              # Tailwind imports
│   ├── i18n.js              # i18n configuration
│   ├── assets/
│   │   └── svelte.svg
│   ├── locales/
│   │   ├── en.json          # English translations
│   │   └── id.json          # Indonesian translations
│   └── lib/
│       ├── components/      # Reusable UI components
│       ├── views/           # Page-level views
│       ├── stores/          # State management
│       ├── services/        # SSE service
│       ├── utils/           # API utilities
│       └── test-utils/      # Test helpers
├── package.json
├── vite.config.js
├── svelte.config.js
└── vitest.config.js
```

### Components Directory

```
lib/components/
├── ActivityLog.svelte
├── ArticleCard.svelte
├── BottomNav.svelte
├── ConfirmModal.svelte
├── ContentListItem.svelte
├── DashboardStats.svelte
├── ImageUploader.svelte
├── PatientModal.svelte
├── PhoneEditModal.svelte
├── ProfileModal.svelte
├── QuillEditor.svelte
├── ReminderModal.svelte
├── SendReminderModal.svelte
├── Sidebar.svelte
├── UserModal.svelte
├── VideoCard.svelte
├── VideoEditModal.svelte
├── VideoModal.svelte
├── WhatsAppPreview.svelte
├── analytics/
│   ├── ContentAnalyticsWidget.svelte
│   ├── DeliveryAnalyticsWidget.svelte
│   └── FailedDeliveryCard.svelte
├── content/
│   ├── ContentChip.svelte
│   ├── ContentDisclaimer.svelte
│   ├── ContentPickerModal.svelte
│   └── ContentPreviewPanel.svelte
├── delivery/
│   ├── DeliveryStatusBadge.svelte
│   └── DeliveryStatusFilter.svelte
├── health/
│   └── SystemHealthWidget.svelte
├── indicators/
│   ├── FailedReminderBadge.svelte
│   └── QuietHoursHint.svelte
├── patients/
│   ├── PatientDetailPane.svelte
│   ├── PatientListPane.svelte
│   └── ReminderListTab.svelte
├── reminders/
│   └── CancelConfirmationModal.svelte
└── ui/
    ├── EmptyState.svelte
    └── Toast.svelte
```

### Views Directory

```
lib/views/
├── ArticleEditorView.svelte
├── BeritaDetailView.svelte
├── BeritaView.svelte
├── CMSDashboardView.svelte
├── DashboardView.svelte
├── LoginScreen.svelte
├── PatientsView.svelte
├── UsersView.svelte
├── VideoEdukasiView.svelte
├── VideoManagerView.svelte
├── analytics/
│   └── FailedDeliveriesView.svelte
├── cms/
│   └── CmsAnalyticsView.svelte
└── patients/
    ├── PatientDetailView.svelte
    └── ReminderHistoryView.svelte
```

## State Management

### Svelte 5 Runes Pattern

The codebase uses modern Svelte 5 runes for reactivity:

```svelte
// Reactive state
let searchQuery = $state("");
let patients = $state([]);
let showModal = $state(false);

// Derived values
let filteredPatients = $derived(
  patients.filter(p => p.name.includes(searchQuery()))
);

// Side effects
$effect(() => {
  if (token) loadPatients();
});

// Props
let { patients = [], onSave = () => {} } = $props();
```

### Class-Based Stores

**`delivery.svelte.js`** - Real-time delivery status

```javascript
export class DeliveryStore {
  deliveryStatuses = $state({});
  connectionStatus = $state('disconnected');
  failedReminders = $state([]);

  updateStatus(reminderId, status) { /* ... */ }
  connect() { /* SSE connection */ }
  disconnect() { /* Cleanup */ }
}
```

**`toast.svelte.js`** - Toast notifications

```javascript
export class ToastStore {
  toasts = $state([]);

  add(message, type = 'info', duration = 3000) { /* ... */ }
  remove(id) { /* ... */ }
}
```

### Legacy Store Pattern

**`auth.js`** - Authentication state

```javascript
export const auth = createAuthStore();
// Provides: token, user, loading, setUser, setToken, logout
```

## API Integration

**File:** `lib/utils/api.js`

All API calls use fetch with JWT Bearer tokens:

```javascript
const API_URL = 'http://localhost:8080/api';

async function fetchWithAuth(endpoint, options = {}) {
  const token = localStorage.getItem('token');
  const response = await fetch(`${API_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers
    }
  });
  return response.json();
}
```

### API Categories

| Category | Endpoints |
|----------|-----------|
| Auth | login, register, fetchUser |
| Patients | fetchPatients, savePatient, deletePatient |
| Reminders | saveReminder, toggleReminder, sendReminder, retryReminder |
| CMS | fetchArticles, createArticle, fetchVideos, createVideo |
| Content | fetchAllContent, fetchPopularContent |
| Analytics | getContentAnalytics, getDeliveryAnalytics, getFailedDeliveries |
| Health | getHealth, getHealthDetailed |

## Component Hierarchy

```
App.svelte
├── Sidebar.svelte (Desktop nav)
├── BottomNav.svelte (Mobile nav)
├── Toast.svelte (Global notifications)
└── Views (conditional rendering):
    ├── LoginScreen.svelte
    ├── DashboardView.svelte
    ├── PatientsView.svelte
    │   ├── PatientListPane.svelte
    │   └── PatientDetailPane.svelte
    ├── UsersView.svelte (superadmin only)
    ├── CMSDashboardView.svelte (admin only)
    ├── ArticleEditorView.svelte
    ├── VideoManagerView.svelte
    ├── BeritaView.svelte
    ├── VideoEdukasiView.svelte
    ├── analytics/FailedDeliveriesView.svelte
    └── cms/CmsAnalyticsView.svelte

Modals:
├── PatientModal.svelte
├── ReminderModal.svelte
│   ├── ContentPickerModal.svelte
│   └── WhatsAppPreview.svelte
├── SendReminderModal.svelte
├── PhoneEditModal.svelte
└── UserModal.svelte
```

## Real-Time Updates (SSE)

**File:** `lib/services/sse.js`

```javascript
class SSEService {
  connect() {
    const token = localStorage.getItem('token');
    const url = `http://localhost:8080/api/sse/delivery-status?token=${token}`;
    this.eventSource = new EventSource(url);

    this.eventSource.addEventListener('delivery.status.updated', (e) => {
      const data = JSON.parse(e.data);
      deliveryStore.updateStatus(data.reminder_id, data.status);
    });
  }
}
```

### SSE Events

- `connection.established` - Connected
- `connection.status` - Status changes
- `delivery.status.updated` - Real-time delivery update
- `delivery.failed` - Failed delivery notification

## Navigation Pattern

**State-based routing** (NOT SvelteKit):

```svelte
// In App.svelte
let currentView = $state(localStorage.getItem('currentView') || 'dashboard');

function navigateTo(view) {
  currentView = view;
  localStorage.setItem('currentView', view);
}

// Conditional rendering
{#if currentView === 'dashboard'}
  <DashboardView />
{:else if currentView === 'patients'}
  <PatientsView />
{/if}
```

**Important:** Per CLAUDE.md, SvelteKit imports are PROHIBITED:
- NO `goto` from '$app/navigation'
- NO `page` from '$app/stores'
- NO `browser` from '$app/environment'

## Key Components

### Authentication

| Component | Purpose |
|-----------|---------|
| `LoginScreen.svelte` | Login/register with password strength validation |

### Patient Management

| Component | Purpose |
|-----------|---------|
| `PatientsView.svelte` | Two-pane list + detail view |
| `PatientListPane.svelte` | Searchable patient list |
| `PatientDetailPane.svelte` | Patient details + reminders |
| `PatientModal.svelte` | Add/edit patient form |

### Reminders

| Component | Purpose |
|-----------|---------|
| `ReminderModal.svelte` | Create/edit with content attachment |
| `SendReminderModal.svelte` | Confirm and send with preview |
| `ContentPickerModal.svelte` | Select articles/videos to attach |
| `ContentPreviewPanel.svelte` | Preview selected content |

### CMS (Admin)

| Component | Purpose |
|-----------|---------|
| `CMSDashboardView.svelte` | Content management dashboard |
| `ArticleEditorView.svelte` | Quill rich text editor |
| `VideoManagerView.svelte` | Add YouTube videos |

### Analytics

| Component | Purpose |
|-----------|---------|
| `FailedDeliveriesView.svelte` | Failed deliveries with CSV export |
| `DeliveryAnalyticsWidget.svelte` | Delivery statistics |
| `ContentAnalyticsWidget.svelte` | Content attachment stats |

## Internationalization

**Files:** `locales/en.json`, `locales/id.json`

The app supports English and Indonesian with:
- Full UI translation
- Date/time formatting
- Number formatting
- Locale switching in ProfileModal

## Testing

**Framework:** Vitest with Testing Library

```bash
# Run tests
bun run test

# Run once (no watch)
bun run test -- --run

# Run specific file
bun run test api.test.js
```

### Test Structure

- Unit tests for stores (`.svelte.js`)
- Component tests using Testing Library
- Rendering, interaction, and edge case tests

## Build & Deployment

```bash
# Development
cd frontend && bun run dev

# Production build
cd frontend && bun run build

# Preview production build
bun run preview
```

### Build Output

- Output directory: `dist/`
- Assets optimized by Vite
- Ready for static hosting

---

_Generated using BMAD Method `document-project` workflow_
