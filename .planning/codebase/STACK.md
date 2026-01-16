# Technology Stack

**Analysis Date:** 2025-01-17

## Languages

**Primary:**
- Go 1.25.5 - Backend API server and services
- JavaScript/TypeScript - Frontend Svelte 5 components

**Secondary:**
- CSS (Tailwind CSS 4) - Styling
- JSON - Data persistence format
- YAML - Configuration files

## Runtime

**Environment:**
- Go runtime (built binary or `go run main.go`)
- Node.js/Bun for frontend development

**Package Manager:**
- Go modules (`go.mod`/`go.sum`)
- Bun (`bun install`/`bun run`) for frontend
- Lockfiles present: `backend/go.sum`, `frontend/bun.lockb`

## Frameworks

**Backend:**
- Gin v1.11.0 - HTTP web framework
  - Location: `backend/handlers/`, `backend/main.go`
  - Router, middleware, request binding

**Frontend:**
- Svelte 5.43.8 - UI framework
  - Location: `frontend/src/`
  - Uses runes: `$state()`, `$derived()`, `$effect()`
- Vite 7.2.4 - Build tool and dev server
  - Config: `frontend/vite.config.js`
- Tailwind CSS 4.1.18 - CSS framework
  - Config: `@tailwindcss/vite` plugin in vite.config.js

**Testing:**
- Vitest 4.0.16 - Frontend unit testing
  - Config: `frontend/vitest.config.js`
  - Test files: `**/*.test.js` co-located with components
- Jest 30.2.0 - Alternative test runner
  - Config: `frontend/jest.config.js`
- Go testing - Backend unit testing
  - `go test ./...` in backend directory

## Key Dependencies

**Backend Critical:**
- `github.com/gin-gonic/gin v1.11.0` - Web framework
- `github.com/golang-jwt/jwt/v5 v5.2.1` - JWT authentication
- `github.com/disintegration/imaging v1.6.2` - Image processing
- `gopkg.in/yaml.v3 v3.0.1` - YAML configuration parsing

**Frontend Critical:**
- `svelte-i18n ^4.0.1` - Internationalization (EN/ID)
- `quill ^2.0.3` - Rich text editor for CMS content
- `@tailwindcss/vite ^4.1.18` - Tailwind CSS integration

**Development:**
- `@sveltejs/vite-plugin-svelte ^6.2.1` - Svelte integration
- `@testing-library/svelte ^5.3.1` - Component testing
- `happy-dom ^20.0.11` - DOM environment for tests

## Configuration

**Environment:**
- Backend: `backend/.env` and `backend/config.yaml`
- Frontend: `frontend/.env` with `VITE_*` prefix
- Environment variable expansion in YAML: `${VAR_NAME}` syntax

**Key configs required:**
- `GOWA_ENDPOINT` - WhatsApp gateway URL
- `GOWA_USER` / `GOWA_PASS` - GOWA authentication
- `VITE_API_URL` - Backend API endpoint
- `VITE_CLERK_PUBLISHABLE_KEY` - Clerk auth (configured but not actively used)

**Build config files:**
- `backend/config/config.go` - Configuration loading with defaults
- `frontend/vite.config.js` - Vite with Svelte and Tailwind plugins
- `frontend/svelte.config.js` - Svelte preprocessing

## Platform Requirements

**Development:**
- Go 1.25+
- Node.js/Bun
- GOWA server (WhatsApp gateway) running on port 3000
- Backend on port 8080, Frontend on port 5173

**Production:**
- Go binary deployment
- Static frontend build (`vite build`)
- GOWA WhatsApp gateway endpoint
- JSON file storage for persistence

---

*Stack analysis: 2025-01-17*
