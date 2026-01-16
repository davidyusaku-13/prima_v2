# Codebase Structure

**Analysis Date:** 2026-01-17

## Directory Layout

```
/media/davidyusaku/Windows/BACKUP/Portfolio/Web/prima_v2/
├── .claude/                      # Claude Code configuration and commands
├── .planning/
│   └── codebase/                 # Architecture/quality analysis docs
├── backend/
│   ├── config/                   # Configuration loading
│   ├── data/                     # JSON data files (patients.json, users.json, etc.)
│   ├── handlers/                 # HTTP handlers
│   ├── models/                   # Data models and in-memory stores
│   ├── services/                 # Business services (GOWA, Scheduler)
│   ├── uploads/                  # Uploaded images
│   ├── utils/                    # Utility functions
│   └── main.go                   # Backend entry point
├── frontend/
│   ├── src/
│   │   ├── assets/               # Static assets (images, fonts)
│   │   ├── lib/
│   │   │   ├── components/       # Reusable UI components
│   │   │   ├── services/         # API and SSE services
│   │   │   ├── stores/           # Svelte 5 state stores
│   │   │   ├── utils/            # Client utilities
│   │   │   └── views/            # Page-level components
│   │   ├── locales/              # i18n translation files
│   │   ├── test/                 # Test utilities
│   │   └── App.svelte            # Frontend entry point
│   ├── dist/                     # Production build output
│   ├── public/                   # Public static files
│   └── node_modules/             # Dependencies
└── docs/                         # Documentation
```

## Directory Purposes

**Backend Directories:**

| Directory | Purpose | Key Files |
|-----------|---------|-----------|
| `backend/config/` | YAML configuration loading and type definitions | `config.go`, `config_test.go` |
| `backend/handlers/` | HTTP request handlers and middleware | `reminder.go`, `content.go`, `sse.go`, `webhook.go`, `analytics.go` |
| `backend/models/` | Data structures and in-memory stores | `patient.go`, `content.go` |
| `backend/services/` | Business services and external integrations | `gowa.go`, `scheduler.go` |
| `backend/utils/` | Helper functions | `phone.go`, `youtube.go`, `message.go`, `logger.go` |
| `backend/data/` | JSON persistence files | `patients.json`, `users.json`, `articles.json`, `videos.json`, `categories.json` |
| `backend/uploads/` | Uploaded image files | Generated at runtime |

**Frontend Directories:**

| Directory | Purpose | Key Files |
|-----------|---------|-----------|
| `frontend/src/lib/components/` | Reusable UI components | `Sidebar.svelte`, `PatientModal.svelte`, `ReminderModal.svelte` |
| `frontend/src/lib/components/patients/` | Patient-specific components | `PatientListPane.svelte`, `PatientDetailPane.svelte`, `ReminderListTab.svelte` |
| `frontend/src/lib/components/reminders/` | Reminder-specific components | `CancelConfirmationModal.svelte` |
| `frontend/src/lib/components/content/` | CMS content components | `ContentPickerModal.svelte`, `ArticleCard.svelte`, `VideoCard.svelte` |
| `frontend/src/lib/components/delivery/` | Delivery status components | `DeliveryStatusBadge.svelte`, `DeliveryStatusFilter.svelte` |
| `frontend/src/lib/components/analytics/` | Analytics components | `FailedDeliveryCard.svelte`, `DeliveryAnalyticsWidget.svelte` |
| `frontend/src/lib/components/ui/` | Generic UI components | `Toast.svelte`, `EmptyState.svelte`, `ConfirmModal.svelte` |
| `frontend/src/lib/components/whatsapp/` | WhatsApp preview components | `WhatsAppPreview.svelte` |
| `frontend/src/lib/components/health/` | System health components | `SystemHealthWidget.svelte` |
| `frontend/src/lib/components/indicators/` | Status indicators | `FailedReminderBadge.svelte`, `QuietHoursHint.svelte` |
| `frontend/src/lib/views/` | Page-level views | `DashboardView.svelte`, `PatientsView.svelte`, `CMSDashboardView.svelte` |
| `frontend/src/lib/views/patients/` | Patient-specific views | `PatientDetailView.svelte`, `ReminderHistoryView.svelte` |
| `frontend/src/lib/views/cms/` | CMS views | `CmsAnalyticsView.svelte` |
| `frontend/src/lib/views/analytics/` | Analytics views | `FailedDeliveriesView.svelte` |
| `frontend/src/lib/services/` | Services | `sse.js` (SSE client) |
| `frontend/src/lib/stores/` | State stores | `delivery.svelte.js`, `toast.svelte.js` |
| `frontend/src/lib/utils/` | Client utilities | `api.js` (API client) |
| `frontend/src/locales/` | i18n translations | `en.json`, `id.json` |

## Key File Locations

**Entry Points:**

| File | Purpose |
|------|---------|
| `backend/main.go` | Backend application entry point, router configuration, initialization |
| `frontend/src/App.svelte` | Frontend SPA entry point, view routing, state management |

**Configuration:**

| File | Purpose |
|------|---------|
| `backend/config/config.go` | Configuration loading and type definitions |
| `backend/config.yaml` | Application configuration (server, GOWA, logging, etc.) |

**Core Logic:**

| File | Purpose |
|------|---------|
| `backend/handlers/reminder.go` | Reminder CRUD and send operations |
| `backend/handlers/content.go` | Articles, videos, categories CMS operations |
| `backend/services/gowa.go` | GOWA WhatsApp integration with circuit breaker |
| `backend/services/scheduler.go` | Background reminder scheduler |
| `frontend/src/lib/utils/api.js` | API client functions |
| `frontend/src/lib/stores/delivery.svelte.js` | Real-time delivery status store |

**Persistence:**

| File | Purpose |
|------|---------|
| `backend/models/patient.go` | Patient and Reminder models, PatientStore |
| `backend/data/patients.json` | Patient data persistence |
| `backend/data/users.json` | User data persistence |
| `backend/data/articles.json` | Article data persistence |
| `backend/data/videos.json` | Video data persistence |

## Naming Conventions

**Backend (Go):**

| Pattern | Example |
|---------|---------|
| Files | `camelCase.go` for most files, `snake_case_test.go` for tests |
| Packages | Single word or short phrase: `handlers`, `services`, `models` |
| Types/Exports | `PascalCase`: `ReminderHandler`, `PatientStore` |
| Variables/Functions | `camelCase`: `generateToken`, `store.patients` |
| Constants | `PascalCase` or `CONSTANT_CASE`: `DeliveryStatusPending`, `tokenExpiry` |
| Struct Tags | JSON: `json:"user_id,omitempty"`, YAML: `yaml:"user_id"` |

**Frontend (Svelte/JavaScript):**

| Pattern | Example |
|---------|---------|
| Component Files | `PascalCase.svelte`: `PatientModal.svelte`, `DashboardView.svelte` |
| JS/TS Files | `camelCase.js` / `camelCase.svelte.js`: `api.js`, `delivery.svelte.js` |
| Components | `PascalCase`: `<ReminderModal />`, `<DeliveryStatusBadge />` |
| Variables/Functions | `camelCase`: `savePatient()`, `reminderForm` |
| Constants | `CONSTANT_CASE`: `API_URL` |
| Stores | `$` prefix for Svelte store imports (if using legacy stores) |
| Props | `$props()` in Svelte 5: `let { title = 'Default' } = $props()` |
| State | `$state()` in Svelte 5: `let patients = $state([])` |
| Derived | `$derived()` in Svelte 5: `let completedReminders = $derived(...)` |

## Where to Add New Code

**New Backend Handler:**
1. Create file in `backend/handlers/` following existing handler patterns
2. Register routes in `backend/main.go` under appropriate middleware group
3. Add tests in same directory with `_test.go` suffix

**New Backend Model/Store:**
1. Create file in `backend/models/` with `type XStore struct{...}` pattern
2. Add thread-safe methods with `sync.RWMutex`
3. Register store initialization in `backend/main.go`

**New Backend Service:**
1. Create file in `backend/services/` for business logic
2. Use dependency injection for config and logger
3. Add tests in same directory

**New Frontend Component:**
1. Generic UI component → `frontend/src/lib/components/ui/`
2. Domain-specific component → appropriate subdirectory (`patients/`, `content/`, etc.)
3. Page-level view → `frontend/src/lib/views/`

**New Frontend Store:**
1. Create `*.svelte.js` file in `frontend/src/lib/stores/`
2. Use Svelte 5 runes: `$state()`, `$derived()`
3. Export singleton instance: `export const storeName = new StoreClass()`

**New API Endpoint:**
1. Backend → Add handler function in `backend/handlers/` and route in `backend/main.go`
2. Frontend → Add function in `frontend/src/lib/utils/api.js`

**New Utility Function:**
1. Backend → `backend/utils/` following file organization (phone, youtube, message, etc.)
2. Frontend → `frontend/src/lib/utils/` or appropriate subdirectory

## Special Directories

| Directory | Purpose | Generated | Committed |
|-----------|---------|-----------|-----------|
| `backend/data/` | JSON persistence | Yes | Yes |
| `backend/uploads/` | Uploaded images | Yes | No (.gitignore) |
| `frontend/dist/` | Production build | Yes | Yes (committed) |
| `frontend/node_modules/` | Dependencies | Yes | No (.gitignore) |

---

*Structure analysis: 2026-01-17*
