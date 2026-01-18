# Development Guide

> Generated: 2026-01-18 | Scan Level: Exhaustive

## Quick Start

### Prerequisites
- **Go** 1.25.5+ for backend
- **Bun** (or Node.js 18+) for frontend
- **Git** for version control

### Setup

```bash
# Clone repository
git clone <repository-url>
cd prima_v2

# Backend setup
cd backend
go mod download
cp .env.example .env  # If exists, configure GOWA

# Frontend setup
cd ../frontend
bun install
```

### Running the Application

```bash
# Terminal 1: Backend (port 8080)
cd backend && go run main.go

# Terminal 2: Frontend (port 5173)
cd frontend && bun run dev
```

**Access:** http://localhost:5173
**Default credentials:** `superadmin` / `superadmin`

## Development Commands

### Backend

```bash
# Run development server
cd backend && go run main.go

# Run tests
cd backend && go test ./...

# Run tests with coverage
cd backend && go test -cover ./...

# Format code
gofmt -w .

# Build binary
cd backend && go build -o bin/server main.go
```

### Frontend

```bash
# Development server with hot reload
cd frontend && bun run dev

# Run specific test
cd frontend && bun run test -- --run "test name"

# Run all tests
cd frontend && bun run test -- --run

# Build for production
cd frontend && bun run build

# Preview production build
cd frontend && bun run preview
```

## Code Style Guidelines

### Go (Backend)

1. **Formatting**: Run `gofmt -w` on save
2. **Imports**: Group stdlib → external → internal
3. **Naming**:
   - `PascalCase` for exported identifiers
   - `camelCase` for unexported
   - `ID`, `URL`, `API` uppercase in names
4. **Error handling**:
   - Check errors immediately
   - Wrap with context: `fmt.Errorf("context: %w", err)`
5. **Receivers**: Value receivers unless mutation needed
6. **Concurrency**: Use `sync.RWMutex` for shared state

**Example:**
```go
func (h *Handler) GetPatient(c *gin.Context) {
    id := c.Param("id")

    h.store.RLock()
    patient, exists := h.store.Patients[id]
    h.store.RUnlock()

    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
        return
    }

    c.JSON(http.StatusOK, patient)
}
```

### Svelte 5 (Frontend)

1. **Runes only**: Use `$state()`, `$derived()`, `$effect()`, `$props()`
2. **Events**: `onclick={fn}` (not `on:click`)
3. **Store files**: Name as `*.svelte.js`
4. **Immutable updates**: Always create new references
5. **No console.log**: Use `console.warn` or `console.error`
6. **Accessibility**: Require `aria-*` on interactive elements

**Example:**
```svelte
<script>
let { patient, onSave } = $props();
let name = $state(patient?.name || '');
let isValid = $derived(name.length >= 2);

function handleSubmit() {
    if (isValid) {
        onSave({ ...patient, name });
    }
}
</script>

<form onsubmit|preventDefault={handleSubmit}>
    <input bind:value={name} aria-label="Patient name">
    <button type="submit" disabled={!isValid}>Save</button>
</form>
```

### SvelteKit Constraints (CRITICAL)

**DO NOT USE** SvelteKit imports. This is Vite + Svelte 5 only.

| Instead of | Use |
|------------|-----|
| `import { goto } from '$app/navigation'` | `window.location.href = '/path'` |
| `import { page } from '$app/stores'` | Props or context |
| `$app/environment` | `typeof window !== 'undefined'` |

## Git Workflow

### Branch Strategy
- `main` - Production-ready code
- Feature branches from `main`

### Commit Convention

[Conventional Commits](https://conventionalcommits.org):

```
<type>[scope]: <description>

Types: feat, fix, docs, style, refactor, test, chore
```

**Examples:**
```
feat(reminder): add quiet hours scheduling
fix(auth): correct token expiry validation
docs: update API documentation
refactor(frontend): migrate to Svelte 5 runes
test(backend): add reminder handler tests
chore: update dependencies
```

## Environment Variables

### Backend (.env)

```bash
# Server
PORT=8080
GIN_MODE=release  # or debug

# CORS
CORS_ORIGIN=http://localhost:5173

# GOWA WhatsApp Gateway
GOWA_BASE_URL=http://localhost:3000
GOWA_DEVICE_ID=default
GOWA_HMAC_SECRET=your-secret

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### Frontend (.env)

```bash
# API endpoint
VITE_API_URL=http://localhost:8080
```

## Testing

### Backend Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./handlers/...

# Run with coverage
go test -cover ./...
```

### Frontend Tests

```bash
# Run all tests once
bun run test -- --run

# Run specific test file
bun run test -- --run PatientModal

# Watch mode
bun run test

# Coverage
bun run test -- --coverage
```

## API Development

### Adding a New Endpoint

1. **Define handler** in `backend/handlers/`:
```go
func (h *Handler) NewEndpoint(c *gin.Context) {
    // Implementation
}
```

2. **Register route** in `backend/main.go`:
```go
api.GET("/new-endpoint", newHandler.NewEndpoint)
```

3. **Add API function** in `frontend/src/lib/utils/api.js`:
```javascript
export async function fetchNewData(token) {
    return request('/api/new-endpoint', { token });
}
```

### Adding a New Component

1. **Create component** in `frontend/src/lib/components/`:
```svelte
<!-- NewComponent.svelte -->
<script>
let { data, onAction } = $props();
</script>

<div>
    <!-- Template -->
</div>
```

2. **Add tests** in same directory:
```javascript
// NewComponent.test.js
import { describe, it, expect } from 'vitest';
// Tests
```

3. **Import in parent**:
```svelte
import NewComponent from '$lib/components/NewComponent.svelte';
```

## Debugging

### Backend

```go
// Structured logging
appLogger.Info("Message",
    "key", value,
    "user_id", userID,
)

appLogger.Error("Failed operation",
    "error", err.Error(),
)
```

### Frontend

```javascript
// Use console.warn/error, not console.log
console.warn('Debug info:', data);
console.error('Error occurred:', error);
```

### Browser DevTools
- Network tab for API calls
- Application tab for localStorage/auth token
- Console for errors

## Common Issues

### CORS Errors
Verify `CORS_ORIGIN` in backend config matches frontend URL.

### Authentication Failures
1. Check token in localStorage
2. Verify token not expired (7 days)
3. Check role permissions

### SSE Connection Issues
1. Verify token passed in query parameter
2. Check GOWA service running
3. Monitor backend logs for errors

### Build Failures
```bash
# Frontend - clear cache
rm -rf node_modules/.vite
bun install

# Backend - clean build
go clean -cache
go mod download
```

## IDE Setup

### VS Code Extensions
- **Go** - Go language support
- **Svelte for VS Code** - Svelte syntax
- **Tailwind CSS IntelliSense** - CSS classes
- **Error Lens** - Inline errors

### Settings
```json
{
    "editor.formatOnSave": true,
    "go.formatTool": "gofmt",
    "svelte.enable-ts-plugin": true
}
```

## Deployment Notes

### Production Build

```bash
# Backend
cd backend
go build -ldflags="-s -w" -o server main.go

# Frontend
cd frontend
bun run build  # Output in dist/
```

### Environment
- Set `GIN_MODE=release` for production
- Configure proper CORS origin
- Use HTTPS for API endpoints
- Set strong JWT secret
