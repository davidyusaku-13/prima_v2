# PRIMA Development Guide

**Date:** 2026-01-13

## Prerequisites

### Backend

| Requirement | Version | Notes |
|-------------|---------|-------|
| Go | 1.25.5 | Install from golang.org |
| Git | Any recent version | For dependency management |

### Frontend

| Requirement | Version | Notes |
|-------------|---------|-------|
| Node.js/Bun | Latest | Bun recommended for speed |
| Git | Any recent version | For dependency management |

## Getting Started

### 1. Clone and Setup

```bash
# Navigate to project directory
cd prima_v2

# Backend dependencies
cd backend && go mod download

# Frontend dependencies
cd ../frontend && bun install
```

### 2. Environment Configuration

#### Backend

```bash
cd backend

# Copy example config if needed
cp config.example.yaml config.yaml

# Edit configuration (optional - defaults work out of box)
nano config.yaml

# Default superadmin credentials: superadmin / superadmin
```

**Backend Configuration Options:**

| Setting | Default | Description |
|---------|---------|-------------|
| `server.port` | 8080 | API server port |
| `server.cors_origin` | `*` | CORS allowed origin |
| `gowa.endpoint` | `http://localhost:3000` | GOWA server URL |
| `gowa.username` | `admin` | GOWA basic auth |
| `gowa.password` | `password123` | GOWA basic auth |
| `quiet_hours.start_hour` | 21 | Quiet hours start (21:00) |
| `quiet_hours.end_hour` | 6 | Quiet hours end (06:00) |

#### Frontend

```bash
cd frontend

# No .env file needed - API URL is hardcoded to localhost:8080
# For production, update API_URL in lib/utils/api.js
```

### 3. Start Development Servers

**Terminal 1 - Backend:**

```bash
cd backend
go run main.go
```

Expected output:
```
[GIN-debug] [WARNING] Running in "debug" mode...
[GIN-debug] POST   /api/auth/register
[GIN-debug] POST   /api/auth/login
...
[GIN-debug] Listening and serving HTTP on :8080
```

**Terminal 2 - Frontend:**

```bash
cd frontend
bun run dev
```

Expected output:
```
VITE v7.2.4  ready in 340 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

### 4. Access Application

| Service | URL | Credentials |
|---------|-----|-------------|
| Frontend | http://localhost:5173 | N/A (login required) |
| Backend API | http://localhost:8080/api | N/A |
| Health Check | http://localhost:8080/api/health | N/A |

**Default Admin Login:**
- Username: `superadmin`
- Password: `superadmin`

## Development Workflow

### Running Tests

#### Backend Tests

```bash
cd backend

# Run all tests
go test ./...

# Run specific package
go test -v ./config
go test -v -run TestLoad

# Run with coverage
go test -cover ./...
```

#### Frontend Tests

```bash
cd frontend

# Run all tests (watch mode)
bun run test

# Run once (no watch)
bun run test -- --run

# Run specific file
bun run test api.test.js
```

### Building for Production

#### Backend

```bash
cd backend

# Build binary
go build -o prima-backend.exe main.go

# Or use provided build script
```

#### Frontend

```bash
cd frontend

# Production build
bun run build

# Preview production build
bun run preview
```

Output: `frontend/dist/` directory ready for static hosting.

### Code Style

#### Go (Backend)

- Use `gofmt -w` on save
- Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- Group imports (stdlib first, then third-party)
- Error handling: Return errors as values, wrap with context

#### Svelte 5 (Frontend)

- Use Svelte 5 runes: `$state()`, `$derived()`, `$effect()`
- Props: `let { title = 'Default' } = $props()`
- NO SvelteKit imports (`goto`, `page`, `browser`)
- Use native browser APIs: `window.location.href`

### File Organization

```
prima_v2/
├── backend/
│   ├── config/         # Configuration
│   ├── handlers/       # HTTP handlers
│   ├── models/         # Data models
│   ├── services/       # Business logic
│   ├── utils/          # Utilities
│   └── data/           # JSON persistence (auto-generated)
├── frontend/
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/  # Reusable components
│   │   │   ├── views/       # Page views
│   │   │   ├── stores/      # State management
│   │   │   ├── services/    # External services
│   │   │   └── utils/       # Utilities
│   │   └── locales/         # i18n translations
│   └── dist/           # Production build
└── docs/               # Documentation
```

## Common Development Tasks

### Adding a New API Endpoint

1. **Backend:** Create handler in `backend/handlers/`
2. **Backend:** Register route in `backend/main.go`
3. **Frontend:** Add API function in `frontend/src/lib/utils/api.js`
4. **Frontend:** Create/update component to use the API

### Adding a New Component

1. Create `.svelte` file in appropriate directory
2. Use Svelte 5 runes for reactivity
3. Export props with `$props()`
4. Add to parent component

Example:
```svelte
<script>
  let { title = 'Default', onClick = () => {} } = $props();
</script>

<button onclick={onClick}>{title}</button>
```

### Adding Internationalization

1. Add key to `frontend/src/locales/en.json`
2. Add translation to `frontend/src/locales/id.json`
3. Use in component: `{$t('key')}` (with i18n setup)

### Modifying Data Models

1. Update struct in `backend/models/`
2. Update JSON tags as needed
3. Run tests to verify persistence
4. Update frontend TypeScript/interfaces if applicable

## External Services

### GOWA (WhatsApp Integration)

| Setting | Value |
|---------|-------|
| Port | 3000 |
| Endpoint | `/send/message` |
| Webhook | `/api/webhook/gowa` |

Start GOWA server before testing WhatsApp features:
```bash
cd backend && ./whatsapp.exe
# or
cd backend && go run -tags "gowa" main.go
```

### YouTube Integration

- Uses noembed.com API for metadata
- No API key required
- Endpoint: `https://noembed.com/embed?url=<youtube_url>`

## Troubleshooting

### Backend Issues

| Issue | Solution |
|-------|----------|
| "jwt_secret.txt not found" | File auto-generated on first run |
| CORS errors | Check `cors_origin` in config.yaml |
| GOWA connection failed | Ensure GOWA server is running on port 3000 |

### Frontend Issues

| Issue | Solution |
|-------|----------|
| 401 errors | Check login/token storage |
| SSE not connecting | Verify token in SSE URL |
| Styles not loading | Check Tailwind CSS setup |

### Common Commands

```bash
# Backend - reload config without restart
# (Not supported - restart required)

# Frontend - clear cache and rebuild
cd frontend
rm -rf node_modules/.vite
bun run dev

# Frontend - check for updates
npm outdated
```

## Testing Checklist

Before committing:

- [ ] Backend: `go test ./...` passes
- [ ] Frontend: `bun run test -- --run` passes
- [ ] Build succeeds: `go build` and `bun run build`
- [ ] No console errors in browser dev tools

---

_Generated using BMAD Method `document-project` workflow_
