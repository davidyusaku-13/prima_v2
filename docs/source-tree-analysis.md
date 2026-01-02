# Source Tree Analysis - PRIMA

**Generated:** January 2, 2026  
**Project Root:** `e:\Portfolio\Web\prima_v2`

---

## Project Structure Overview

```
prima_v2/
├── backend/                    # Go/Gin API Server
├── frontend/                   # Svelte 5 + Vite Web App
├── docs/                       # Generated documentation
├── _bmad/                      # BMAD Method configuration
├── _bmad-output/               # BMAD workflow outputs
├── AGENTS.md                   # AI agent guidelines
├── CLAUDE.md                   # Claude-specific instructions
├── GOWA-README.md              # WhatsApp integration docs
└── QUILL.md                    # Rich text editor docs
```

---

## Backend (Go/Gin) - Port 8080

### Entry Point

**📍 `backend/main.go`** - Application bootstrap (1240 lines)

**Responsibilities:**

- Load configuration from YAML
- Initialize stores (Patient, User, Content)
- Setup Gin router with CORS
- Register all API routes
- Initialize GOWA client with circuit breaker
- Start reminder scheduler
- Graceful shutdown handling

**Key Components Initialized:**

```go
appConfig        *config.Config
patientStore     *models.PatientStore
contentStore     *handlers.ContentStore
gowaClient       *services.GOWAClient
scheduler        *services.ReminderScheduler
reminderHandler  *handlers.ReminderHandler
webhookHandler   *handlers.WebhookHandler
sseHandler       *handlers.SSEHandler
analyticsHandler *handlers.AnalyticsHandler
healthHandler    *handlers.HealthHandler
```

---

### Directory Structure (Detailed)

```
backend/
├── main.go                     # 🚀 Application entry point (1240 LOC)
│   ├── Setup: Config, stores, GOWA, scheduler
│   ├── Routes: REST API endpoints
│   ├── Middleware: Auth, CORS
│   └── Graceful shutdown
│
├── config/                     # Configuration management
│   ├── config.go              # YAML config loader + defaults
│   └── config_test.go         # Config loading tests
│
├── config.yaml                # 📝 Runtime configuration
│   ├── Server (port, CORS)
│   ├── GOWA (endpoint, credentials, webhook secret)
│   ├── Circuit breaker (failure threshold, cooldown)
│   ├── Retry policy (attempts, delays)
│   ├── Logging (level, format)
│   ├── Disclaimer (text, enabled)
│   └── Quiet hours (start, end, timezone)
│
├── data/                      # 💾 JSON file persistence
│   ├── patients.json          # Patient records with embedded reminders
│   ├── users.json             # User accounts (hashed passwords)
│   ├── articles.json          # Health education articles (Berita)
│   ├── videos.json            # YouTube educational videos
│   ├── categories.json        # Content categories
│   ├── items.json             # Generic item storage
│   └── jwt_secret.txt         # JWT signing key (Base64)
│
├── uploads/                   # 🖼️ Image file storage
│   └── [article-uuid]-16x9.jpg  # Hero images (3 aspect ratios)
│
├── handlers/                  # 🔌 HTTP request handlers
│   ├── analytics.go           # Delivery & content analytics
│   ├── analytics_test.go
│   ├── content.go             # CMS endpoints (articles, videos, categories)
│   ├── content_test.go
│   ├── health.go              # Health check endpoints
│   ├── health_test.go
│   ├── reminder.go            # Reminder CRUD + send/retry/cancel
│   ├── reminder_test.go
│   ├── sse.go                 # Server-Sent Events for real-time updates
│   ├── sse_test.go
│   ├── webhook.go             # GOWA delivery status webhooks
│   └── webhook_test.go
│
├── models/                    # 📦 Domain entities
│   ├── patient.go             # Patient, Reminder, Recurrence, Attachment
│   │   └── PatientStore with RWMutex
│   └── content.go             # Category, Article, Video
│       └── CategoryStore, ArticleStore, VideoStore
│
├── services/                  # 🔧 Business logic layer
│   ├── gowa.go                # GOWA WhatsApp client with circuit breaker
│   ├── gowa_test.go
│   ├── scheduler.go           # Reminder auto-send scheduler (cron-like)
│   └── scheduler_test.go
│
└── utils/                     # 🛠️ Shared utilities
    ├── hmac.go                # HMAC-SHA256 webhook validation
    ├── logger.go              # Structured logging (slog)
    ├── logger_test.go
    ├── mask.go                # Phone/email masking for privacy
    ├── mask_test.go
    ├── message.go             # WhatsApp message formatting
    ├── phone.go               # Phone number normalization/validation
    ├── phone_test.go
    ├── quiethours.go          # Quiet hours enforcement (22:00-06:00 WIB)
    ├── quiethours_test.go
    └── youtube.go             # YouTube metadata fetching (noembed.com)
```

---

### Critical Backend Paths

#### Configuration

- **`config/config.go`** - YAML config loader with environment variable support
- **`config.yaml`** - Runtime settings (server, GOWA, circuit breaker, retry, quiet hours)

#### Data Persistence

- **`data/*.json`** - All application data (patients, users, content)
- Thread-safe with `sync.RWMutex`
- Pretty-printed JSON for readability

#### HTTP Handlers

| Handler                 | Responsibility        | Key Routes                         |
| ----------------------- | --------------------- | ---------------------------------- |
| `main.go`               | Auth, Patients, Users | `/api/auth/*`, `/api/patients/*`   |
| `handlers/reminder.go`  | Reminder management   | `/api/patients/:id/reminders/*`    |
| `handlers/content.go`   | CMS operations        | `/api/articles/*`, `/api/videos/*` |
| `handlers/analytics.go` | Statistics            | `/api/analytics/*`                 |
| `handlers/health.go`    | Health checks         | `/api/health`                      |
| `handlers/sse.go`       | Real-time updates     | `/api/sse/delivery-status`         |
| `handlers/webhook.go`   | GOWA callbacks        | `/api/webhook/gowa`                |

#### Services

- **`services/gowa.go`** - GOWA HTTP client
  - Circuit breaker pattern (5 failures → 5min cooldown)
  - Retry with exponential backoff (1s, 5s, 30s, 2m, 10m)
  - HMAC webhook verification
- **`services/scheduler.go`** - Background reminder processor
  - Checks every 60 seconds
  - Processes scheduled, retrying, quiet_hours reminders
  - Auto-sends when due
  - Respects quiet hours (22:00-06:00 WIB)

#### Models

- **`models/patient.go`** - Patient, Reminder, Recurrence, Attachment, PatientStore
- **`models/content.go`** - Category, Article, Video, stores for each

#### Utilities

- **`utils/phone.go`** - Normalize `08xxx` → `628xxx`, validate Indonesian format
- **`utils/quiethours.go`** - Check if current time is in quiet hours
- **`utils/mask.go`** - Mask phone (`6281234***789`) and email for logs
- **`utils/logger.go`** - Structured logging with slog
- **`utils/message.go`** - Format WhatsApp message with attachments
- **`utils/youtube.go`** - Fetch YouTube metadata via noembed.com
- **`utils/hmac.go`** - Validate GOWA webhook signatures

---

## Frontend (Svelte 5 + Vite) - Port 5173

### Entry Point

**📍 `frontend/index.html`** - HTML shell  
**📍 `frontend/src/main.js`** - JavaScript entry (imports App.svelte, i18n, CSS)

---

### Directory Structure (Detailed)

```
frontend/
├── index.html                 # HTML shell
├── src/
│   ├── main.js                # 🚀 Application bootstrap
│   ├── App.svelte             # Root Svelte component
│   ├── app.css                # Global Tailwind imports
│   ├── i18n.js                # svelte-i18n configuration (EN/ID)
│   │
│   ├── assets/                # Static assets (icons, images)
│   │
│   ├── lib/                   # 📦 Components & utilities
│   │   ├── components/        # Reusable Svelte 5 components
│   │   │   ├── auth/          # Login, Register components
│   │   │   ├── patients/      # Patient list, form, detail
│   │   │   ├── reminders/     # Reminder CRUD, status display
│   │   │   ├── content/       # Article/video list, detail, editor
│   │   │   ├── dashboard/     # Admin dashboard, stats
│   │   │   ├── analytics/     # Delivery analytics, failed deliveries
│   │   │   ├── shared/        # Buttons, modals, tables, alerts
│   │   │   └── SVELTE5_COMPONENT_TEMPLATE.md  # Component template
│   │   │
│   │   ├── stores/            # Svelte stores (*.svelte.js)
│   │   │   ├── auth.svelte.js        # User auth state
│   │   │   ├── theme.svelte.js       # Theme preferences
│   │   │   └── notifications.svelte.js  # Toast notifications
│   │   │
│   │   ├── api/               # API client functions
│   │   │   ├── auth.js        # Login, register, me
│   │   │   ├── patients.js    # Patient CRUD
│   │   │   ├── reminders.js   # Reminder operations
│   │   │   ├── content.js     # CMS API calls
│   │   │   └── analytics.js   # Analytics endpoints
│   │   │
│   │   └── utils/             # Helper functions
│   │       ├── api.js         # Fetch wrapper with auth
│   │       ├── date.js        # Date formatting
│   │       ├── validation.js  # Form validation
│   │       └── constants.js   # App constants
│   │
│   ├── locales/               # 🌐 Translations
│   │   ├── en.json            # English
│   │   └── id.json            # Indonesian
│   │
│   └── test/                  # 🧪 Vitest tests
│       ├── api.test.js
│       ├── components.test.js
│       └── utils.test.js
│
├── package.json               # Dependencies & scripts
├── vite.config.js             # Vite build config (aliases, plugins)
├── svelte.config.js           # Svelte preprocessor
├── vitest.config.js           # Test configuration
├── jsconfig.json              # VS Code IntelliSense
├── README.md                  # Frontend documentation
└── CHECKLIST.md               # Development checklist
```

---

### Critical Frontend Paths

#### Application Entry

- **`src/main.js`** - Bootstrap app, mount to DOM
- **`src/App.svelte`** - Root component, router logic
- **`src/i18n.js`** - i18n initialization (EN/ID)

#### Components

**Feature-Based Organization:**

```
lib/components/
├── auth/          # Authentication flows
├── patients/      # Patient management
├── reminders/     # Reminder CRUD + delivery tracking
├── content/       # CMS (articles, videos, categories)
├── dashboard/     # Admin overview
├── analytics/     # Delivery & content stats
└── shared/        # Reusable UI components
```

#### Stores (Svelte 5 Runes)

- **`lib/stores/auth.svelte.js`** - User session, token, role
- **`lib/stores/theme.svelte.js`** - Dark/light mode
- **`lib/stores/notifications.svelte.js`** - Toast notifications

**Access pattern:**

```javascript
import { auth } from "$lib/stores/auth.svelte.js";

// Read
$auth.user.role;

// Update
auth.login(token, user);
```

#### API Client

- **`lib/api/*.js`** - Fetch wrappers for backend endpoints
- Automatic JWT token injection
- Error handling with notifications

**Example:**

```javascript
// lib/api/patients.js
export async function getPatients() {
  return api.get("/api/patients");
}
```

#### Localization

- **`locales/en.json`** - English translations
- **`locales/id.json`** - Indonesian translations

**Usage:**

```svelte
<script>
  import { _ } from 'svelte-i18n';
</script>

<h1>{$_('patients.title')}</h1>
```

---

## Integration Points

### Frontend → Backend (REST API)

**Base URL:** `http://localhost:8080/api`

**Authentication:**

```javascript
headers: {
  'Authorization': `Bearer ${token}`,
  'Content-Type': 'application/json'
}
```

**Key Endpoints:**

- `/api/auth/login` - Authentication
- `/api/patients` - Patient CRUD
- `/api/patients/:id/reminders` - Reminder management
- `/api/articles` - Article CMS
- `/api/videos` - Video CMS
- `/api/analytics/*` - Statistics
- `/api/sse/delivery-status?token=<jwt>` - Real-time updates

---

### Backend → GOWA (WhatsApp Gateway)

**Endpoint:** `http://localhost:3000` (configurable via `config.yaml`)

**Send Message:**

```
POST /send/message
Authorization: Basic <base64(user:password)>

{
  "phone": "628123456789",
  "message": "Reminder: Minum obat..."
}
```

**Response:**

```json
{
  "message_id": "gowa-msg-id-123"
}
```

**Circuit Breaker:**

- Closed: Normal operation
- Open: After 5 failures, no requests for 5 minutes
- Half-Open: Test with single request after cooldown

---

### GOWA → Backend (Webhook)

**Endpoint:** `POST /api/webhook/gowa`

**Authentication:** HMAC-SHA256 signature in `X-Webhook-Signature` header

**Payload:**

```json
{
  "event": "message.ack",
  "message": {
    "id": "gowa-msg-id-123",
    "status": "delivered"
  }
}
```

**Backend Action:**

1. Validate HMAC signature
2. Lookup reminder by `gowa_message_id`
3. Update `delivery_status` and timestamps
4. Broadcast via SSE to connected clients
5. Save to `patients.json`

---

### Backend → Frontend (Server-Sent Events)

**Endpoint:** `GET /api/sse/delivery-status?token=<jwt>`

**Content-Type:** `text/event-stream`

**Events:**

```javascript
// connection.established
{
  "message": "Connected to delivery status updates",
  "timestamp": "2026-01-02T10:00:00Z"
}

// delivery.status.updated
{
  "reminder_id": "uuid",
  "status": "delivered",
  "timestamp": "2026-01-02T10:00:05Z"
}

// delivery.failed
{
  "reminder_id": "uuid",
  "patient_id": "uuid",
  "patient_name": "John Doe",
  "error": "GOWA timeout",
  "timestamp": "2026-01-02T10:00:10Z"
}
```

**Frontend Usage:**

```javascript
const eventSource = new EventSource(`/api/sse/delivery-status?token=${token}`);

eventSource.addEventListener("delivery.status.updated", (e) => {
  const data = JSON.parse(e.data);
  updateReminderStatus(data.reminder_id, data.status);
});
```

---

### Backend → YouTube (Metadata Fetch)

**Endpoint:** `https://noembed.com/embed?url=<youtube_url>`

**Usage:** When admin adds a video via `/api/videos`

**Response:**

```json
{
  "title": "Video Title",
  "author_name": "Channel Name",
  "thumbnail_url": "https://i.ytimg.com/vi/VIDEO_ID/maxresdefault.jpg",
  "duration": "5:32"
}
```

---

## Development Workflow

### Backend Development

```bash
cd backend
go run main.go           # Start server (port 8080)
go test ./...            # Run all tests
go test -v ./handlers    # Test specific package
```

### Frontend Development

```bash
cd frontend
bun run dev              # Start dev server (port 5173)
bun run build            # Production build
bun run test             # Run Vitest tests
```

### Full Stack

**Terminal 1:**

```bash
cd backend && go run main.go
```

**Terminal 2:**

```bash
cd frontend && bun run dev
```

**Browser:** `http://localhost:5173`

---

## Build Artifacts

### Backend

```
backend/
├── prima_v2              # Compiled binary (Linux)
├── prima_v2.exe          # Compiled binary (Windows)
└── data/                 # Runtime data (must exist)
```

**Build:**

```bash
go build -o prima_v2
```

### Frontend

```
frontend/
└── dist/                 # Production build output
    ├── index.html
    ├── assets/
    │   ├── index-[hash].js
    │   └── index-[hash].css
    └── uploads/          # Copy from backend/uploads
```

**Build:**

```bash
bun run build
```

**Serve:**

- Static file server (Nginx, Caddy, etc.)
- Point API requests to backend via reverse proxy

---

## Key Design Patterns

### Backend

1. **Layered Architecture**

   - Handlers (HTTP layer)
   - Services (Business logic)
   - Models (Data structures)
   - Utils (Shared functions)

2. **Repository Pattern**

   - PatientStore, ContentStore
   - Thread-safe with RWMutex
   - In-memory with JSON persistence

3. **Circuit Breaker**

   - GOWA integration resilience
   - Prevents cascading failures

4. **Event Broadcasting**
   - SSE for real-time updates
   - Webhook triggers broadcasts

### Frontend

1. **Component-Based**

   - Svelte 5 components with runes
   - Feature-based organization

2. **Reactive State**

   - `$state()` for local state
   - Stores for global state
   - `$derived()` for computed values

3. **API Client Layer**
   - Centralized fetch logic
   - Automatic auth injection
   - Error handling

---

**Next:** See [Integration Architecture](./integration-architecture.md) for detailed data flow diagrams.
