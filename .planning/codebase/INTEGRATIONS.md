# External Integrations

**Analysis Date:** 2025-01-17

## APIs & External Services

**Messaging:**
- GOWA (Go WhatsApp) - WhatsApp gateway service
  - Used for: Sending reminder messages to patients via WhatsApp
  - SDK/Client: Custom `GOWAClient` in `backend/services/gowa.go`
  - Endpoint: `GOWA_ENDPOINT` (default: `http://localhost:3000`)
  - Auth: Basic auth with `GOWA_USER` and `GOWA_PASSWORD`
  - Webhook: `POST /api/webhook/gowa` for delivery status updates
  - Circuit breaker: Prevents cascade failures (5 failures = open circuit, 5min cooldown)

**Video/Media:**
- noembed.com - YouTube video metadata API
  - Used for: Fetching video title, description, author, thumbnail
  - Endpoint: `https://noembed.com/embed?url=...`
  - Implementation: `backend/utils/youtube.go`
  - Fetches: Title, description, author_name, thumbnail_url, duration, HTML embed

**Authentication:**
- Clerk - Identity provider (configured but not actively used)
  - Publishable key: `VITE_CLERK_PUBLISHABLE_KEY` in frontend .env
  - Active auth: Custom JWT implementation in `backend/main.go`
  - JWT expiry: 7 days
  - Roles: superadmin, admin, volunteer

## Data Storage

**Databases:**
- None (No SQL or NoSQL database)
  - Persistence: JSON files in `backend/data/`
  - Files: `patients.json`, `users.json`, `articles.json`, `videos.json`, `categories.json`
  - Concurrency: `sync.RWMutex` for thread-safe access

**File Storage:**
- Local filesystem
  - Uploads: `backend/uploads/` directory
  - Served at: `/uploads` route
  - Image processing: `github.com/disintegration/imaging` library

**Caching:**
- None (in-memory only)
  - Content store in `handlers/content.go` loads from JSON on startup

## Authentication & Identity

**Auth Provider:**
- Custom JWT implementation
  - Secret: `backend/data/jwt_secret.txt` (auto-generated)
  - Token expiry: 7 days (168 hours)
  - Middleware: `authMiddleware()` in `backend/main.go`
  - SSE auth: Query parameter `?token=` for EventSource API

**RBAC (Role-Based Access Control):**
- superadmin - Full access, user management
- admin - CMS features, analytics, all patients
- volunteer - Create/manage own patients only

## Monitoring & Observability

**Error Tracking:**
- None (no external error tracking service)
- Logging: Structured JSON logging via `slog`
  - Config: `backend/config.yaml` logging.level/format

**Logs:**
- Output: stdout (JSON or text format)
- Level: debug, info, warn, error (configurable)

## CI/CD & Deployment

**Hosting:**
- Not specified (self-hosted deployment)
- Backend: Go binary or `go run main.go`
- Frontend: Vite static build (`npm run build`)

**CI Pipeline:**
- Not detected (no GitHub Actions, CircleCI, etc.)

## Environment Configuration

**Required env vars:**
- Backend (`backend/.env`):
  - `GOWA_ENDPOINT` - WhatsApp gateway URL
  - `GOWA_USER` - GOWA username
  - `GOWA_PASS` - GOWA password
- Frontend (`frontend/.env`):
  - `VITE_API_URL` - Backend API base URL
  - `VITE_CLERK_PUBLISHABLE_KEY` - Clerk publishable key

**Secrets location:**
- `backend/.env` - Environment variables (not committed)
- `backend/data/jwt_secret.txt` - Auto-generated JWT secret
- `backend/.env.example` - Template (committed)

## Webhooks & Callbacks

**Incoming:**
- GOWA delivery status webhook
  - Endpoint: `POST /api/webhook/gowa`
  - Auth: HMAC validation via `webhook_secret`
  - Handler: `handlers/webhook.go`
  - Purpose: Update reminder delivery status in real-time

**Outgoing:**
- GOWA API calls
  - Send message: `POST {endpoint}/send/message`
  - Auth: Basic auth (user:password base64 encoded)

---

*Integration audit: 2025-01-17*
