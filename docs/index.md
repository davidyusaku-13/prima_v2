# PRIMA Documentation Index

**Healthcare Volunteer Dashboard - Technical Documentation**  
**Generated:** January 2, 2026  
**Version:** Exhaustive Scan  
**Project Type:** Backend (Go/Gin) + Frontend (Svelte 5 + Vite)

---

## 🎯 Quick Navigation

- [Technology Stack](#-technology-stack)
- [Architecture](#-architecture)
- [API Reference](#-api-reference)
- [Data Models](#-data-models)
- [Development Guides](#-development-guides)
- [Existing Documentation](#-existing-documentation)

---

## 📦 Technology Stack

**Document:** [technology-stack.md](./technology-stack.md)

Complete overview of all technologies, frameworks, libraries, and external services used in PRIMA.

**Key Technologies:**

- **Backend:** Go 1.25.5 + Gin v1.11.0
- **Frontend:** Svelte 5.43.8 + Vite 7.2.4 (NOT SvelteKit)
- **Styling:** Tailwind CSS 4.1.18
- **Data:** JSON file persistence with `sync.RWMutex`
- **External:** GOWA WhatsApp Gateway, noembed.com for YouTube

**Covers:**

- Backend dependencies (21 packages)
- Frontend dependencies (31 packages)
- External API integrations
- Architecture decisions & rationale

---

## 🏗️ Architecture

### Integration Architecture

**Document:** [integration-architecture.md](./integration-architecture.md)

Complete system integration architecture showing how all parts communicate.

**Covers:**

- System overview diagram
- 5 integration points (Frontend↔Backend, SSE, GOWA, Webhooks, YouTube API)
- Data flow diagrams (reminder send flow, CMS article creation)
- Security measures (JWT, RBAC, CORS, HMAC)
- Performance considerations (circuit breaker, SSE connection management)
- Deployment architecture (dev + production)

### Backend Architecture

**Document:** [architecture-backend.md](./architecture-backend.md)

Detailed architecture of the Go/Gin backend service.

**Covers:**

- Layered architecture (handlers/services/models/utils)
- Routing structure (54 endpoints)
- Middleware chain (CORS, JWT, rate limiting, logging)
- Data persistence strategy (JSON + RWMutex)
- Scheduler implementation (cron-based reminders)
- Circuit breaker pattern (for GOWA integration)
- Error handling strategy
- Concurrency & thread safety (RWMutex patterns)
- Security architecture (JWT, RBAC, CORS)
- Performance patterns (SSE broadcasting, JSON optimization)

### Frontend Architecture

**Document:** [architecture-frontend.md](./architecture-frontend.md)

Detailed architecture of the Svelte 5 + Vite frontend application.

**Covers:**

- Component hierarchy (App.svelte → Views → Components)
- State management (Svelte stores with runes: $state, $derived, $effect)
- Manual routing implementation (NOT svelte-spa-router)
- API client layer (async/await pattern)
- SSE real-time updates (EventSource)
- i18n implementation (svelte-i18n for EN/ID)
- Styling system (Tailwind CSS 4 + custom dark mode)
- Build & bundle strategy (Vite 7 with code splitting)

### Source Tree Analysis

**Document:** [source-tree-analysis.md](./source-tree-analysis.md)

Complete annotated directory structure for both backend and frontend.

**Covers:**

- Backend structure (main.go entry, handlers/, models/, services/, utils/, data/)
- Frontend structure (src/lib/components/, stores/, api/, locales/)
- Integration points (REST API, SSE, webhooks)
- Development workflow (running, testing, building)
- Build artifacts & outputs
- Design patterns used

---

## 📡 API Reference

### API Contracts (Backend)

**Document:** [api-contracts-backend.md](./api-contracts-backend.md)

Complete reference for all 54 backend REST API endpoints.

**Categories:**

1. **Authentication** (3 endpoints) - Login, register, profile
2. **Patients** (5 endpoints) - CRUD operations
3. **Reminders** (10 endpoints) - Scheduling, sending, delivery tracking
4. **Content** (13 endpoints) - Articles, videos, categories (CMS)
5. **Analytics** (5 endpoints) - Dashboard metrics, logs
6. **Health** (2 endpoints) - Service health checks
7. **SSE** (1 endpoint) - Real-time delivery status updates
8. **Webhooks** (1 endpoint) - GOWA callbacks
9. **Users** (4 endpoints) - User management (superadmin)
10. **Config** (2 endpoints) - System configuration

**Each Endpoint Includes:**

- HTTP method & path
- Authentication requirements
- Role-based access control
- Request body schema
- Response schema with examples
- Error codes & messages
- Rate limiting info (where applicable)

---

## 💾 Data Models

### Data Models (Backend)

**Document:** [data-models-backend.md](./data-models-backend.md)

Complete data entity definitions and relationships.

**Entities:**

1. **Patient** - Patient records with embedded reminders
2. **Reminder** - Medication/appointment reminders
3. **Recurrence** - Reminder repeat patterns
4. **Attachment** - Files attached to reminders
5. **User** - System users (volunteers, admins)
6. **Category** - Content categorization
7. **Article** - Educational health articles
8. **Video** - YouTube educational videos

**Covers:**

- Go struct definitions with JSON tags
- Field validation rules
- Relationships & foreign keys
- Delivery status state machine (11 states)
- Thread-safe store operations (RWMutex)
- Persistence implementation (JSON marshaling)
- Entity-Relationship Diagram
- Security considerations (masking, access control)

---

## 👨‍💻 Development Guides

### Getting Started

**Quick Start Commands:**

```bash
# Backend (port 8080)
cd backend && go run main.go

# Frontend (port 5173)
cd frontend && bun run dev

# Backend Tests
cd backend && go test ./...

# Frontend Tests
cd frontend && bun run test
```

### Development Guide (Backend)

### Development Guide (Backend)

**Document:** [development-guide-backend.md](./development-guide-backend.md)

**Covers:**

- Environment setup (Go installation, dependencies)
- Configuration (`config.yaml`)
- Running locally (`go run main.go`)
- Testing strategy (unit tests with table-driven tests)
- Code style (gofmt, naming conventions)
- Adding new endpoints (handler → route → test)
- Adding new models (JSON persistence)
- Debugging techniques (Delve debugger)
- Common issues & solutions
- Production deployment checklist

### Development Guide (Frontend)

### Development Guide (Frontend)

**Document:** [development-guide-frontend.md](./development-guide-frontend.md)

**Covers:**

- Environment setup (Bun installation, dependencies)
- Project structure (src/lib/components/, stores/, views/)
- Svelte 5 runes (`$state`, `$derived`, `$effect`, `$props`)
- Component creation patterns (reactive state, props, events)
- API integration (api.js client)
- State management with stores (auth, delivery, toast)
- i18n (adding translations with svelte-i18n)
- Styling with Tailwind CSS 4 (dark mode, custom utilities)
- Testing components (Vitest + Testing Library)
- Building for production (`bun run build`)
- Common issues & solutions

### Component Inventory (Frontend)

### Component Inventory (Frontend)

**Document:** [component-inventory-frontend.md](./component-inventory-frontend.md)

**Covers:**

- Complete list of 40+ Svelte components
- Component hierarchy (navigation, modals, features, analytics, UI primitives)
- Props & events for each component
- Reusable UI components (Toast, ImageUploader, QuillEditor)
- Feature-specific components (DashboardStats, ActivityLog, ArticleCard)
- Layout components (Sidebar, BottomNav)
- Modal components (PatientModal, ReminderModal, UserModal, etc.)
- Usage examples & best practices
- Component template for creating new components

---

## 📚 Existing Documentation

These documents existed before this exhaustive scan and remain useful:

### AGENTS.md

**Path:** [../AGENTS.md](../AGENTS.md)

Project overview and AI agent guidance.

**Contents:**

- Project overview (Vite + Svelte 5, **NOT** SvelteKit)
- Commands (backend, frontend, tests, build)
- Code style guidelines (Go, JavaScript/Svelte)
- Architecture summary
- Commit message conventions

### CLAUDE.md

**Path:** [../CLAUDE.md](../CLAUDE.md)

Specific guidance for Claude AI assistant when working on this project.

### GOWA-README.md

**Path:** [../GOWA-README.md](../GOWA-README.md)

Documentation for the GOWA WhatsApp Gateway integration.

**Contents:**

- GOWA API reference
- Authentication (Basic Auth)
- Send message endpoint
- Webhook callback specification
- HMAC signature validation
- Status codes

### QUILL.md

**Path:** [../QUILL.md](../QUILL.md)

_(Content unknown - exists in workspace)_

### Frontend Checklist

**Path:** [../frontend/CHECKLIST.md](../frontend/CHECKLIST.md)

Development checklist for frontend features and tasks.

### Bun Testing Guide

**Path:** [../frontend/bun-test.md](../frontend/bun-test.md)

Guide for running tests with Bun runtime.

---

## 🔧 Configuration Files

### Backend Configuration

**File:** `backend/config.yaml`

Runtime configuration for the backend service.

**Key Settings:**

- Server port (default: 8080)
- CORS origin (default: http://localhost:5173)
- GOWA endpoint, credentials, webhook secret
- Scheduler interval
- Log level

**Example:**

```yaml
server:
  port: 8080
  cors_origin: "http://localhost:5173"

gowa:
  endpoint: "http://localhost:3000"
  username: "admin"
  password: "admin"
  webhook_secret: "your-secret-key"

scheduler:
  interval: 60 # seconds

log_level: "info"
```

### Frontend Configuration

**Files:**

- `frontend/vite.config.js` - Vite build configuration
- `frontend/svelte.config.js` - Svelte compiler options
- `frontend/jsconfig.json` - JavaScript path mappings (`$lib` alias)
- `frontend/vitest.config.js` - Test runner configuration

---

## 🔐 Security

### Authentication & Authorization

**JWT Authentication:**

- Algorithm: HS256
- Token expiry: 7 days
- Secret stored in: `backend/data/jwt_secret.txt`
- Claims: `userId`, `username`, `role`

**Password Security:**

- Hash: SHA256
- Encoding: Base64
- Stored in: `backend/data/users.json`

**Role-Based Access Control:**

- `superadmin` - All operations + user management
- `admin` - CMS management + analytics
- `volunteer` - Patient management (own patients only)

**Default Credentials:**

- Username: `superadmin`
- Password: `superadmin`
- ⚠️ **Change in production!**

### Data Security

**Masking:**

- Phone numbers: `628123456789` → `6281234***789`
- Emails: `user@example.com` → `u***@example.com`
- Used in: Logs, admin analytics views

**HMAC Webhook Validation:**

- Algorithm: HMAC-SHA256
- Header: `X-Webhook-Signature`
- Validates GOWA webhook authenticity

---

## 📊 Data Persistence

### JSON File Storage

**Location:** `backend/data/`

**Files:**

- `patients.json` - Patient records with embedded reminders
- `users.json` - System users (volunteers, admins)
- `articles.json` - Educational articles (CMS)
- `videos.json` - YouTube videos (CMS)
- `categories.json` - Content categories
- `items.json` - _(Purpose unknown)_
- `jwt_secret.txt` - JWT signing secret

**Thread Safety:**

- `sync.RWMutex` for concurrent read/write
- Multiple readers allowed
- Single writer at a time

**Backup Strategy:**

- Manual: Copy `data/` folder
- Automated: _(To be implemented)_

---

## 🚀 Deployment

### Development

**Backend:**

```bash
cd backend
go run main.go
# Listens on http://localhost:8080
```

**Frontend:**

```bash
cd frontend
bun run dev
# Listens on http://localhost:5173
```

**GOWA:**

```bash
# Run in separate terminal or Docker
# Listens on http://localhost:3000
```

### Production

**Build Frontend:**

```bash
cd frontend
bun run build
# Output: frontend/dist/
```

**Build Backend:**

```bash
cd backend
go build -o prima-backend main.go
# Output: backend/prima-backend (or prima-backend.exe on Windows)
```

**Deployment Guide:**

**Document:** [deployment-guide.md](./deployment-guide.md)

**Covers:**

- Prerequisites (server requirements, domain setup)
- Backend deployment (systemd service, configuration)
- Frontend deployment (Vite build, static hosting)
- Nginx reverse proxy configuration (with SSL)
- SSL/TLS setup (Let's Encrypt/Certbot)
- Process management (systemd, alternative PM2)
- Database backup (automated cron jobs)
- Monitoring (health checks, logs, uptime)
- Troubleshooting common issues
- Security checklist & firewall configuration
- Updating application (rollback strategy)
- Docker deployment (alternative)

---

## 🧪 Testing

### Backend Tests

**Run All Tests:**

```bash
cd backend
go test ./...
```

**Run Specific Package:**

```bash
go test -v ./config
```

**Run Specific Test:**

```bash
go test -v -run TestLoadConfig
```

**With Coverage:**

```bash
go test -cover ./...
```

**Test Structure:**

- Unit tests: `*_test.go` files alongside source
- Test data: Inline or in `testdata/` folders
- Mocking: Manual or with interfaces

### Frontend Tests

**Run All Tests:**

```bash
cd frontend
bun run test
```

**Run Once (No Watch):**

```bash
bun run test -- --run
```

**Run Specific File:**

```bash
bun run test api.test.js
```

**Test Structure:**

- Unit tests: `src/test/*.test.js`
- Component tests: `src/lib/components/**/*.test.js`
- Testing library: Vitest
- DOM testing: jsdom environment

---

## 🌍 Internationalization (i18n)

**Supported Languages:**

- English (en)
- Indonesian (id)

**Implementation:**

- Library: `svelte-i18n`
- Translation files: `frontend/src/locales/en.json`, `id.json`
- Setup: `frontend/src/i18n.js`

**Usage in Components:**

```javascript
import { _ } from 'svelte-i18n';

<h1>{$_('dashboard.title')}</h1>
<button>{$_('common.save')}</button>
```

**Adding Translations:**

1. Add key to `en.json` and `id.json`
2. Use `$_('key')` in component
3. Translations load asynchronously on init

---

## 📈 Analytics & Monitoring

### Built-in Analytics

**Endpoints:**

- `GET /api/analytics/dashboard` - Dashboard metrics
- `GET /api/analytics/delivery-log` - Delivery history
- `GET /api/analytics/summary` - Summary statistics

**Metrics:**

- Total patients, reminders, volunteers
- Delivery success/failure rates
- Content views (articles, videos)
- Most active volunteers
- Top viewed content

### Health Monitoring

**Public Health Check:**

```bash
curl http://localhost:8080/api/health
```

**Detailed Health (Admin):**

```bash
curl -H "Authorization: Bearer TOKEN" \
  http://localhost:8080/api/health/detailed
```

**Includes:**

- Backend status
- GOWA connectivity
- Circuit breaker state
- Queue statistics (scheduled, retrying)

---

## 🐛 Troubleshooting

### Common Issues

#### Backend won't start

**Error:** `Address already in use`

**Solution:**

```bash
# Find process using port 8080
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # macOS/Linux

# Kill process
taskkill /PID <PID> /F        # Windows
kill -9 <PID>                 # macOS/Linux
```

#### Frontend can't connect to backend

**Error:** `Network error` or `CORS error`

**Solution:**

1. Check backend is running (`http://localhost:8080/api/health`)
2. Verify CORS origin in `backend/config.yaml`:
   ```yaml
   server:
     cors_origin: "http://localhost:5173"
   ```
3. Restart backend after config change

#### GOWA messages not sending

**Error:** Circuit breaker open / Connection refused

**Solution:**

1. Check GOWA is running (`http://localhost:3000`)
2. Verify credentials in `backend/config.yaml`
3. Check circuit breaker state:
   ```bash
   curl -H "Authorization: Bearer TOKEN" \
     http://localhost:8080/api/health/detailed
   ```
4. If open, wait 5 minutes for half-open state

#### SSE not updating UI

**Issue:** Delivery status not updating in real-time

**Solution:**

1. Check browser console for SSE errors
2. Verify SSE connection:
   ```javascript
   // Should see in Network tab:
   // /api/sse/delivery-status?token=...
   // Type: eventsource
   ```
3. Check token is valid (not expired)
4. Restart SSE connection:
   ```javascript
   eventSource.close();
   // Reconnect
   ```

---

## 📝 Code Style & Conventions

### Go (Backend)

**Formatting:**

- Use `gofmt -w .` on save
- Imports grouped: stdlib, external, internal
- Naming: `PascalCase` (exported), `camelCase` (unexported)
- Acronyms: `HTTPServer`, `parseURL` (uppercase)

**Error Handling:**

```go
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}
```

**Struct Tags:**

```go
type Patient struct {
    ID   string `json:"id" yaml:"id"`
    Name string `json:"name" yaml:"name"`
}
```

### JavaScript/Svelte (Frontend)

**Formatting:**

- ESLint config: `.eslintrc.svelte5.js`
- Naming: `camelCase` (variables), `PascalCase` (components)
- Constants: `CONSTANT_CASE`

**Svelte 5 Runes:**

```javascript
// State
let count = $state(0);

// Derived
let doubled = $derived(count * 2);

// Effect
$effect(() => {
  console.log("Count changed:", count);
});

// Props
let { title = "Default" } = $props();
```

**Event Handlers:**

```svelte
<!-- Svelte 5 -->
<button onclick={handleClick}>Click</button>

<!-- NOT SvelteKit -->
<!-- ✗ import { goto } from '$app/navigation' -->
<!-- ✓ window.location.href = '/path' -->
```

---

## 🤝 Contributing

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer]
```

**Types:**

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation changes
- `style` - Code style (formatting, no logic change)
- `refactor` - Code refactoring
- `test` - Adding/updating tests
- `chore` - Maintenance tasks

**Examples:**

```
feat(cms): add video thumbnail preview
fix(reminder): respect quiet hours on send
docs(api): update webhook signature validation
```

### Development Workflow

1. **Create branch:**

   ```bash
   git checkout -b feat/your-feature
   ```

2. **Make changes** (follow code style)

3. **Test changes:**

   ```bash
   # Backend
   cd backend && go test ./...

   # Frontend
   cd frontend && bun run test
   ```

4. **Commit with conventional message:**

   ```bash
   git commit -m "feat(patients): add export to CSV"
   ```

5. **Push & create PR:**
   ```bash
   git push origin feat/your-feature
   ```

---

## 📖 Additional Resources

### External Documentation

- **Go:** https://go.dev/doc/
- **Gin:** https://gin-gonic.com/docs/
- **Svelte 5:** https://svelte-5-preview.vercel.app/docs
- **Vite:** https://vitejs.dev/guide/
- **Tailwind CSS 4:** https://tailwindcss.com/docs

### Project-Specific Guides

- **GOWA Integration:** [../GOWA-README.md](../GOWA-README.md)
- **Agent Guidelines:** [../AGENTS.md](../AGENTS.md)
- **Claude Assistant Guide:** [../CLAUDE.md](../CLAUDE.md)

---

## 📋 Document Status

### ✅ Complete (11 Documents)

- [x] Technology Stack
- [x] API Contracts (Backend)
- [x] Data Models (Backend)
- [x] Integration Architecture
- [x] Source Tree Analysis
- [x] Architecture (Backend)
- [x] Architecture (Frontend)
- [x] Development Guide (Backend)
- [x] Development Guide (Frontend)
- [x] Component Inventory (Frontend)
- [x] Deployment Guide

### 📝 To Be Generated (0 Documents)

_All planned documentation has been generated._

### 📌 Notes

- This index was generated via exhaustive scan (reads ALL source files)
- Scan date: January 2, 2026
- Documentation completed: January 2, 2026
- Project version: Current (no formal versioning)
- Documentation will be updated as project evolves

---

**Last Updated:** January 2, 2026  
**Scan Type:** Exhaustive  
**Total Documents:** 5 complete, 6 pending  
**For Questions:** Refer to existing docs or create GitHub issue
