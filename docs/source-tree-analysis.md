# PRIMA Source Tree Analysis

**Date:** 2026-01-13

## Overview

PRIMA is a multi-part project with separate frontend (Svelte 5 + Vite) and backend (Go + Gin) directories. The project structure follows a clean separation of concerns with clear boundaries between the client and server applications.

## Multi-Part Structure

This project consists of 2 distinct parts:

- **Backend** (`backend/`): Go/Gin API server
- **Frontend** (`frontend/`): Svelte 5 + Vite SPA

## Complete Directory Structure

```
prima_v2/
├── backend/                     # Go/Gin Backend API
│   ├── main.go                  # Entry point, routes, auth middleware
│   ├── config.yaml              # Runtime configuration
│   ├── go.mod                   # Go module definition
│   ├── go.sum                   # Go dependencies
│   ├── config/                  # Configuration handling
│   │   ├── config.go
│   │   └── config_test.go
│   ├── data/                    # JSON persistence (excluded from scan)
│   │   ├── patients.json
│   │   ├── users.json
│   │   ├── categories.json
│   │   ├── articles.json
│   │   ├── videos.json
│   │   └── jwt_secret.txt
│   ├── handlers/                # HTTP request handlers
│   │   ├── analytics.go
│   │   ├── content.go
│   │   ├── health.go
│   │   ├── reminder.go
│   │   ├── sse.go
│   │   ├── webhook.go
│   │   └── *_test.go
│   ├── models/                  # Data models
│   │   ├── content.go
│   │   └── patient.go
│   ├── services/                # Business logic
│   │   ├── gowa.go              # WhatsApp integration
│   │   ├── scheduler.go         # Automatic reminder processing
│   │   └── *_test.go
│   ├── utils/                   # Utility functions
│   │   ├── hmac.go              # Webhook validation
│   │   ├── logger.go
│   │   ├── mask.go              # PII masking
│   │   ├── message.go           # WhatsApp formatting
│   │   ├── phone.go             # Phone normalization
│   │   ├── quiethours.go        # Quiet hours logic
│   │   └── youtube.go           # YouTube metadata
│   └── uploads/                 # Uploaded images
│
├── frontend/                    # Svelte 5 + Vite Frontend
│   ├── src/
│   │   ├── main.js              # Entry point (mount)
│   │   ├── App.svelte           # Root component, routing
│   │   ├── app.css              # Tailwind imports
│   │   ├── i18n.js              # i18n configuration
│   │   ├── assets/
│   │   │   └── svelte.svg
│   │   ├── locales/             # Translations
│   │   │   ├── en.json
│   │   │   └── id.json
│   │   └── lib/
│   │       ├── components/      # Reusable UI components
│   │       ├── views/           # Page-level views
│   │       ├── stores/          # State management
│   │       ├── services/        # SSE service
│   │       ├── utils/           # API utilities
│   │       └── test-utils/      # Test helpers
│   ├── package.json
│   ├── vite.config.js
│   ├── svelte.config.js
│   ├── vitest.config.js
│   ├── dist/                    # Production build
│   └── node_modules/
│
├── docs/                        # Project documentation
│   ├── .archive/
│   │   └── (archived files)
│   └── svelte/
│       └── llms-full.txt
│
├── _bmad/                       # BMAD framework
│   ├── core/
│   │   ├── tasks/
│   │   └── workflows/
│   └── bmm/
│       ├── config.yaml
│       └── workflows/
│
├── _bmad-output/                # BMAD workflow outputs
│   └── project-scan-report.json
│
├── CLAUDE.md                    # AI agent guidelines
├── AGENTS.md                    # Agent instructions
├── GOWA-README.md               # WhatsApp integration docs
└── QUILL.md                     # Rich text editor docs
```

## Critical Directories

### `backend/`

| Directory | Purpose | Contains |
|-----------|---------|----------|
| `handlers/` | HTTP request handlers | Route implementations for all API endpoints |
| `models/` | Data structures | Patient, Reminder, Article, Video, User models |
| `services/` | Business logic | GOWA client, ReminderScheduler |
| `config/` | Configuration | Config loading and validation |
| `utils/` | Helpers | Phone formatting, YouTube, logging |
| `data/` | Persistence | JSON files with thread-safe access |

### `frontend/src/`

| Directory | Purpose | Contains |
|-----------|---------|----------|
| `lib/components/` | UI components | Buttons, modals, forms, cards |
| `lib/views/` | Page views | Dashboard, Patients, CMS, Analytics |
| `lib/stores/` | State management | Auth, Delivery, Toast stores |
| `lib/services/` | External services | SSE connection |
| `lib/utils/` | Utilities | API client, helpers |
| `locales/` | Translations | English and Indonesian |
| `assets/` | Static assets | Images, icons |

## Entry Points

### Backend Entry Point

- **File:** `backend/main.go`
- **Bootstrap:** Initializes Gin router, loads config, sets up JWT middleware, registers all routes, starts scheduler goroutine

### Frontend Entry Point

- **File:** `frontend/src/main.js`
- **Bootstrap:** Mounts App.svelte to DOM, initializes i18n

### Root Component

- **File:** `frontend/src/App.svelte`
- **Responsibilities:** State-based routing, view orchestration, navigation, toast display

## File Organization Patterns

| Pattern | Location | Examples |
|---------|----------|----------|
| `*_test.go` | Same dir as source | `config_test.go`, `handlers/*_test.go` |
| `*.svelte.js` | `lib/stores/` | `delivery.svelte.js`, `toast.svelte.js` |
| `*View.svelte` | `lib/views/` | `DashboardView.svelte`, `PatientsView.svelte` |
| `*Modal.svelte` | `lib/components/` | `PatientModal.svelte`, `ReminderModal.svelte` |
| `*.test.js` | Same dir as source | `delivery.test.js`, `Toast.test.js` |

## Configuration Files

| File | Purpose |
|------|---------|
| `backend/config.yaml` | Backend runtime configuration |
| `backend/go.mod` | Go module definition |
| `frontend/package.json` | Frontend dependencies |
| `frontend/vite.config.js` | Vite build configuration |
| `frontend/svelte.config.js` | Svelte compiler options |
| `frontend/vitest.config.js` | Vitest configuration |
| `CLAUDE.md` | AI agent guidelines |

## Integration Points

### Frontend → Backend

- **Type:** REST API over HTTP
- **Base URL:** `http://localhost:8080/api`
- **Authentication:** JWT Bearer token
- **Real-time:** Server-Sent Events at `/api/sse/delivery-status`

### Backend → GOWA (WhatsApp)

- **Type:** HTTP POST
- **Endpoint:** GOWA server on port 3000
- **Auth:** Basic Auth
- **Webhook:** `/api/webhook/gowa` for delivery status

## Asset Locations

| Type | Location | Description |
|------|----------|-------------|
| Uploaded images | `backend/uploads/` | Hero images for articles |
| Static assets | `frontend/src/assets/` | Icons, images |
| Built assets | `frontend/dist/` | Production build output |

---

_Generated using BMAD Method `document-project` workflow_
