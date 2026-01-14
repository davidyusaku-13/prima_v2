# PRIMA Documentation Index

**Type:** Multi-part project (Client-Server)
**Primary Language:** Go (Backend), TypeScript/JavaScript (Frontend)
**Architecture:** REST API + SPA
**Last Updated:** 2026-01-13

## Project Overview

PRIMA is a Healthcare Volunteer Dashboard that enables healthcare workers to manage patients, send health education reminders via WhatsApp, and manage CMS content for health education materials.

**Tech Stack:**
- **Backend:** Go 1.25.5 + Gin Framework
- **Frontend:** Svelte 5 + Vite + Tailwind CSS 4
- **Database:** JSON file persistence (no external DB)
- **Authentication:** JWT with Role-Based Access Control
- **Messaging:** WhatsApp via GOWA server
- **i18n:** English + Indonesian

## Quick Reference

| Aspect | Details |
|--------|---------|
| **Backend Port** | 8080 |
| **Frontend Port** | 5173 (dev), 80 (prod) |
| **GOWA Port** | 3000 |
| **Default Admin** | superadmin / superadmin |
| **JWT Expiry** | 7 days |

## Generated Documentation

### Core Documentation

- [Project Overview](./architecture-overview.md) - Executive summary and high-level architecture
- [Source Tree Analysis](./source-tree-analysis.md) - Annotated directory structure

### Backend Documentation

- [Backend Architecture](./architecture-backend.md) - Detailed Go/Gin backend architecture
- [API Contracts](./api-contracts-backend.md) - API endpoints and schemas
- [Data Models](./data-models-backend.md) - Database schema and models
- [Development Guide](./development-guide.md) - Local setup and development workflow
- [Deployment Guide](./deployment-guide.md) - Deployment process and infrastructure

### Frontend Documentation

- [Frontend Architecture](./architecture-frontend.md) - Detailed Svelte 5 frontend architecture
- [Component Inventory](./component-inventory-frontend.md) - Catalog of UI components

### Integration

- [Integration Architecture](./integration-architecture.md) - How parts communicate
- [Development Guide](./development-guide.md) - Setup and dev workflow for both parts

### Existing Documentation

- [CLAUDE.md](../CLAUDE.md) - AI agent guidelines
- [AGENTS.md](../AGENTS.md) - Agent instructions
- [GOWA-README.md](../GOWA-README.md) - WhatsApp integration docs
- [QUILL.md](../QUILL.md) - Rich text editor documentation

## Getting Started

### Prerequisites

**Backend:**
- Go 1.25.5
- Git

**Frontend:**
- Node.js or Bun
- Git

### Setup

```bash
# Backend
cd backend
go mod download
go run main.go

# Frontend (new terminal)
cd frontend
bun install
bun run dev
```

### Run Tests

```bash
# Backend
cd backend && go test ./...

# Frontend
cd frontend && bun run test -- --run
```

### Build for Production

```bash
# Backend
cd backend && go build -o prima-backend main.go

# Frontend
cd frontend && bun run build
```

## For AI-Assisted Development

This documentation was generated specifically to enable AI agents to understand and extend this codebase.

### When Planning New Features:

**UI-only features:**
→ Reference: `architecture-frontend.md`, `component-inventory-frontend.md`

**API/Backend features:**
→ Reference: `architecture-backend.md`, `api-contracts-backend.md`, `data-models-backend.md`

**Full-stack features:**
→ Reference: All architecture docs + `integration-architecture.md`

**Deployment changes:**
→ Reference: `deployment-guide.md`

## Project Structure

```
prima_v2/
├── backend/                     # Go/Gin API server
│   ├── config/                  # Configuration
│   ├── handlers/                # HTTP handlers
│   ├── models/                  # Data models
│   ├── services/                # Business logic
│   ├── utils/                   # Utilities
│   ├── data/                    # JSON persistence
│   └── main.go                  # Entry point
│
├── frontend/                    # Svelte 5 SPA
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/      # UI components
│   │   │   ├── views/           # Page views
│   │   │   ├── stores/          # State management
│   │   │   ├── services/        # SSE service
│   │   │   └── utils/           # API utilities
│   │   └── locales/             # Translations
│   ├── dist/                    # Production build
│   └── package.json
│
├── docs/                        # This documentation
├── _bmad/                       # BMAD framework
└── CLAUDE.md                    # AI guidelines
```

## Role-Based Access Control

| Role | Permissions |
|------|-------------|
| `superadmin` | Full access, user management, all analytics |
| `admin` | CMS content, dashboard stats, analytics, health details |
| `volunteer` | CRUD own patients/reminders, view public content |

## API Endpoints Summary

| Category | Endpoints |
|----------|-----------|
| Auth | `/api/auth/register`, `/api/auth/login`, `/api/auth/me` |
| Patients | `/api/patients` (CRUD) |
| Reminders | `/api/patients/:id/reminders` (CRUD + send/retry/cancel) |
| Content | `/api/articles`, `/api/videos`, `/api/categories` |
| Users | `/api/users` (superadmin only) |
| Analytics | `/api/analytics/delivery`, `/api/analytics/failed-deliveries` |
| Health | `/api/health`, `/api/health/detailed` |

---

_Generated by BMAD Method `document-project` workflow_
