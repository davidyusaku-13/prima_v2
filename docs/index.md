# PRIMA - Project Documentation Index

> Healthcare Volunteer Dashboard | Generated: 2026-01-18

## Project Overview

**PRIMA** is a healthcare volunteer dashboard application for managing patient reminders with WhatsApp delivery. It consists of a Go/Gin backend API and a Svelte 5 + Vite frontend.

### Key Features
- Patient management with reminder scheduling
- WhatsApp message delivery via GOWA integration
- Real-time delivery status tracking (SSE)
- Content Management System (articles & videos)
- Role-based access control (superadmin, admin, volunteer)
- Quiet hours scheduling
- Analytics and failed delivery tracking

### Technology Stack

| Component | Technologies |
|-----------|--------------|
| **Backend** | Go 1.25.5, Gin 1.11.0, JWT, JSON storage |
| **Frontend** | Svelte 5.43.8, Vite 7.2.4, Tailwind CSS 4 |
| **Integration** | GOWA (WhatsApp), YouTube (noembed.com) |
| **Testing** | Go test, Vitest |

---

## Documentation Map

### Architecture

| Document | Description |
|----------|-------------|
| [Architecture: Backend](architecture-backend.md) | Go/Gin API architecture, data models, API contracts, authentication |
| [Architecture: Frontend](architecture-frontend.md) | Svelte 5 patterns, components, state management, routing |

### Development

| Document | Description |
|----------|-------------|
| [Development Guide](development-guide.md) | Setup, commands, code style, testing, debugging |
| [Project Structure](project-structure.md) | Directory tree, file organization, entry points |

### Project Configuration

| File | Location | Purpose |
|------|----------|---------|
| CLAUDE.md | `/CLAUDE.md` | Claude Code instructions and constraints |
| config.yaml | `/backend/config.yaml` | Backend configuration |
| package.json | `/frontend/package.json` | Frontend dependencies |

---

## Quick Reference

### Commands

```bash
# Backend (port 8080)
cd backend && go run main.go

# Frontend (port 5173)
cd frontend && bun run dev

# Tests
cd backend && go test ./...
cd frontend && bun run test -- --run

# Build
cd frontend && bun run build
```

### Default Credentials
- **Username**: `superadmin`
- **Password**: `superadmin`

### Key URLs
- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080/api
- **Health Check**: http://localhost:8080/api/health

---

## Critical Constraints

### Frontend (Svelte 5)
1. **NO SvelteKit imports** - Use `window.location.href`, props/context
2. **Svelte 5 runes only** - `$state()`, `$derived()`, `$effect()`, `$props()`
3. **Events**: `onclick={fn}` (not `on:click`)
4. **Immutable updates** - Always create new array/object references
5. **No console.log** - Use `console.warn`/`console.error`

### Backend (Go)
1. **Thread safety** - Use `sync.RWMutex` for stores
2. **Error wrapping** - `fmt.Errorf("context: %w", err)`
3. **gofmt** on save
4. **Grouped imports** - stdlib → external → internal

---

## Data Flow

```
┌─────────────────┐     REST API      ┌─────────────────┐
│                 │  ◄───────────────►│                 │
│  Svelte 5 SPA   │                   │   Go/Gin API    │
│  (Port 5173)    │  ◄───── SSE ─────│   (Port 8080)   │
│                 │                   │                 │
└─────────────────┘                   └────────┬────────┘
                                               │
                                    ┌──────────┴──────────┐
                                    │                     │
                              ┌─────▼─────┐        ┌──────▼──────┐
                              │   JSON    │        │    GOWA     │
                              │   Files   │        │  (WhatsApp) │
                              │ data/*.json│       │  Port 3000  │
                              └───────────┘        └─────────────┘
```

---

## Authentication & Authorization

### JWT Token
- **Expiry**: 7 days
- **Storage**: localStorage (`token`)
- **Header**: `Authorization: Bearer <token>`

### Roles

| Role | Patients | CMS | Users | Analytics |
|------|----------|-----|-------|-----------|
| superadmin | All | Full | Full | Full |
| admin | All | Full | - | Full |
| volunteer | Own only | - | - | - |

---

## External Integrations

### GOWA (WhatsApp Gateway)
- **Port**: 3000
- **Webhook**: `/api/webhook/gowa` (HMAC validated)
- **Features**: Message delivery, status tracking

### YouTube (noembed.com)
- **Purpose**: Video metadata fetching
- **Endpoint**: `https://noembed.com/embed`

---

## Testing Coverage

| Area | Framework | Command |
|------|-----------|---------|
| Backend handlers | Go test | `go test ./handlers/...` |
| Backend services | Go test | `go test ./services/...` |
| Backend utils | Go test | `go test ./utils/...` |
| Frontend components | Vitest | `bun run test -- --run` |
| Frontend stores | Vitest | `bun run test -- --run store` |

---

## File Locations

### Configuration
- Backend: `/backend/config.yaml`, `/backend/.env`
- Frontend: `/frontend/vite.config.js`, `/frontend/svelte.config.js`

### Data Persistence
- `/backend/data/patients.json` - Patient records
- `/backend/data/users.json` - User accounts
- `/backend/data/articles.json` - CMS articles
- `/backend/data/videos.json` - CMS videos

### Entry Points
- Backend: `/backend/main.go:216`
- Frontend: `/frontend/src/App.svelte`

---

## Commit Convention

```
<type>[scope]: <description>

Types: feat, fix, docs, style, refactor, test, chore
```

Example: `feat(reminder): add quiet hours scheduling`

---

## Related Documentation

- [CLAUDE.md](/CLAUDE.md) - Claude Code project instructions
- [GOWA-README.md](/GOWA-README.md) - WhatsApp integration details

---

*Documentation generated via BMAD document-project workflow*
*Scan Level: Exhaustive | Date: 2026-01-18*
