# Project Structure & Source Tree

> Generated: 2026-01-18 | Scan Level: Exhaustive

## Repository Overview

```
prima_v2/
├── backend/                 # Go/Gin REST API
├── frontend/                # Svelte 5 + Vite SPA
├── docs/                    # Generated documentation
├── _bmad/                   # BMAD workflow configuration
├── CLAUDE.md                # Claude Code instructions
├── GOWA-README.md           # WhatsApp integration docs
└── .git/                    # Git repository
```

## Backend Structure

```
backend/
├── main.go                  # Entry point, routes, auth handlers
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── config.yaml              # Application configuration
├── .env                     # Environment variables (local)
│
├── config/
│   ├── config.go            # Configuration struct & loading
│   └── config_test.go       # Config tests
│
├── handlers/
│   ├── analytics.go         # Analytics API handlers
│   ├── analytics_test.go    # Analytics tests
│   ├── content.go           # CMS content handlers
│   ├── content_test.go      # Content tests
│   ├── health.go            # Health check handlers
│   ├── health_test.go       # Health tests
│   ├── reminder.go          # Reminder CRUD & delivery
│   ├── reminder_test.go     # Reminder tests
│   ├── sse.go               # Server-Sent Events handler
│   ├── sse_test.go          # SSE tests
│   ├── webhook.go           # GOWA webhook handler
│   └── webhook_test.go      # Webhook tests
│
├── models/
│   ├── patient.go           # Patient, Reminder, stores
│   └── content.go           # Article, Video, Category
│
├── services/
│   ├── gowa.go              # GOWA WhatsApp client
│   ├── gowa_test.go         # GOWA tests
│   ├── scheduler.go         # Reminder scheduler
│   └── scheduler_test.go    # Scheduler tests
│
├── utils/
│   ├── hmac.go              # HMAC signature validation
│   ├── logger.go            # Structured logging
│   ├── logger_test.go       # Logger tests
│   ├── mask.go              # Data masking utilities
│   ├── mask_test.go         # Mask tests
│   ├── message.go           # WhatsApp message formatting
│   ├── phone.go             # Phone number validation
│   ├── phone_test.go        # Phone tests
│   ├── quiethours.go        # Quiet hours logic
│   ├── quiethours_test.go   # Quiet hours tests
│   └── youtube.go           # YouTube metadata fetching
│
├── data/                    # Persistent data (JSON files)
│   ├── patients.json        # Patient records
│   ├── users.json           # User accounts
│   ├── categories.json      # Content categories
│   ├── articles.json        # Article content
│   ├── videos.json          # Video references
│   └── jwt_secret.txt       # JWT signing secret
│
├── uploads/                 # User-uploaded files
│   └── images/              # Hero images (16x9, 1x1, 4x3)
│
├── cmd/
│   └── check_db/            # Database check utility
│
└── bin/                     # Compiled binaries
```

## Frontend Structure

```
frontend/
├── index.html               # Entry HTML
├── package.json             # Dependencies & scripts
├── vite.config.js           # Vite configuration
├── svelte.config.js         # Svelte configuration
├── vitest.config.js         # Test configuration
├── bun.lockb                # Bun lock file
│
├── src/
│   ├── App.svelte           # Root component
│   ├── main.js              # Application entry
│   ├── app.css              # Global styles (Tailwind)
│   ├── i18n.js              # i18n configuration
│   │
│   ├── lib/
│   │   ├── components/      # Reusable UI components
│   │   │   ├── Sidebar.svelte
│   │   │   ├── BottomNav.svelte
│   │   │   ├── PatientModal.svelte
│   │   │   ├── ReminderModal.svelte
│   │   │   ├── SendReminderModal.svelte
│   │   │   ├── PhoneEditModal.svelte
│   │   │   ├── ConfirmModal.svelte
│   │   │   ├── ProfileModal.svelte
│   │   │   ├── UserModal.svelte
│   │   │   ├── VideoModal.svelte
│   │   │   ├── VideoEditModal.svelte
│   │   │   ├── QuillEditor.svelte
│   │   │   ├── ImageUploader.svelte
│   │   │   ├── ArticleCard.svelte
│   │   │   ├── VideoCard.svelte
│   │   │   ├── ContentListItem.svelte
│   │   │   ├── DashboardStats.svelte
│   │   │   ├── ActivityLog.svelte
│   │   │   │
│   │   │   ├── ui/
│   │   │   │   ├── Toast.svelte
│   │   │   │   ├── Toast.test.js
│   │   │   │   └── EmptyState.svelte
│   │   │   │
│   │   │   ├── content/
│   │   │   │   ├── ContentPickerModal.svelte
│   │   │   │   ├── ContentPickerModal.*.test.js
│   │   │   │   ├── ContentPreviewPanel.svelte
│   │   │   │   ├── ContentChip.svelte
│   │   │   │   └── ContentDisclaimer.svelte
│   │   │   │
│   │   │   ├── delivery/
│   │   │   │   ├── DeliveryStatusBadge.svelte
│   │   │   │   ├── DeliveryStatusBadge.test.js
│   │   │   │   ├── DeliveryStatusFilter.svelte
│   │   │   │   └── DeliveryStatusFilter.test.js
│   │   │   │
│   │   │   ├── analytics/
│   │   │   │   ├── ContentAnalyticsWidget.svelte
│   │   │   │   ├── DeliveryAnalyticsWidget.svelte
│   │   │   │   └── FailedDeliveryCard.svelte
│   │   │   │
│   │   │   ├── health/
│   │   │   │   └── SystemHealthWidget.svelte
│   │   │   │
│   │   │   ├── indicators/
│   │   │   │   ├── FailedReminderBadge.svelte
│   │   │   │   ├── FailedReminderBadge.test.js
│   │   │   │   └── QuietHoursHint.svelte
│   │   │   │
│   │   │   ├── patients/
│   │   │   │   ├── PatientListPane.svelte
│   │   │   │   ├── PatientListPane.test.js
│   │   │   │   ├── PatientDetailPane.svelte
│   │   │   │   ├── PatientDetailPane.test.js
│   │   │   │   └── ReminderListTab.svelte
│   │   │   │
│   │   │   ├── reminders/
│   │   │   │   ├── CancelConfirmationModal.svelte
│   │   │   │   └── CancelConfirmationModal.test.js
│   │   │   │
│   │   │   └── whatsapp/
│   │   │       └── (WhatsApp integration components)
│   │   │
│   │   ├── views/           # Page-level components
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
│   │   │   │
│   │   │   ├── analytics/
│   │   │   │   └── FailedDeliveriesView.svelte
│   │   │   │
│   │   │   ├── cms/
│   │   │   │   └── CmsAnalyticsView.svelte
│   │   │   │
│   │   │   └── patients/
│   │   │       ├── PatientDetailView.svelte
│   │   │       └── ReminderHistoryView.svelte
│   │   │
│   │   ├── stores/          # State management
│   │   │   ├── auth.js
│   │   │   ├── toast.svelte.js
│   │   │   ├── toast.test.js
│   │   │   ├── delivery.svelte.js
│   │   │   └── delivery.test.js
│   │   │
│   │   ├── services/        # External integrations
│   │   │   └── sse.js
│   │   │
│   │   ├── utils/           # Utility functions
│   │   │   └── api.js
│   │   │
│   │   └── test-utils/      # Test utilities
│   │       └── (test helpers)
│   │
│   ├── locales/             # Translation files
│   │   ├── en.json
│   │   └── id.json
│   │
│   ├── assets/              # Static assets
│   │   └── (images, icons)
│   │
│   └── test/
│       └── __mocks__/       # Test mocks
│
├── public/                  # Public static files
│
└── node_modules/            # Dependencies (git-ignored)
```

## File Count Summary

| Directory | Files | Description |
|-----------|-------|-------------|
| `backend/` | ~32 | Go source files |
| `backend/handlers/` | 12 | API handlers + tests |
| `backend/models/` | 2 | Data models |
| `backend/services/` | 4 | Business services |
| `backend/utils/` | 12 | Utilities |
| `frontend/src/` | ~70 | Svelte/JS files |
| `frontend/src/lib/components/` | ~40 | UI components |
| `frontend/src/lib/views/` | ~14 | Page views |
| `frontend/src/lib/stores/` | 5 | State stores |

## Key Entry Points

| File | Purpose |
|------|---------|
| `backend/main.go` | Backend entry point |
| `frontend/src/main.js` | Frontend entry point |
| `frontend/src/App.svelte` | Root Svelte component |
| `frontend/index.html` | HTML shell |

## Configuration Files

| File | Purpose |
|------|---------|
| `backend/config.yaml` | Backend configuration |
| `backend/.env` | Environment variables |
| `frontend/vite.config.js` | Vite bundler config |
| `frontend/svelte.config.js` | Svelte compiler config |
| `frontend/vitest.config.js` | Test runner config |
| `CLAUDE.md` | Claude Code instructions |

## Data Files

| File | Purpose |
|------|---------|
| `backend/data/patients.json` | Patient records |
| `backend/data/users.json` | User accounts |
| `backend/data/articles.json` | Article content |
| `backend/data/videos.json` | Video references |
| `backend/data/categories.json` | Content categories |
| `backend/data/jwt_secret.txt` | JWT signing key |

## Build Artifacts

| Directory | Purpose |
|-----------|---------|
| `backend/bin/` | Compiled Go binaries |
| `frontend/dist/` | Production build output |
| `frontend/node_modules/` | npm dependencies |
