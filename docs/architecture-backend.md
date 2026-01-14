# PRIMA Backend Architecture

**Date:** 2026-01-13
**Type:** Backend API
**Framework:** Go 1.25.5 + Gin

## Executive Summary

The PRIMA backend is a Go/Gin HTTP API server that powers the Healthcare Volunteer Dashboard. It handles patient management, reminder scheduling with WhatsApp delivery via GOWA, and CMS content management for health education materials.

### Key Characteristics

- **Language:** Go 1.25.5
- **Framework:** Gin (HTTP routing)
- **Persistence:** JSON file storage with `sync.RWMutex`
- **Authentication:** JWT (HS256, 7-day expiry)
- **Authorization:** Role-Based Access Control (superadmin/admin/volunteer)
- **Real-time Updates:** Server-Sent Events (SSE)
- **External Integrations:** GOWA (WhatsApp), noembed.com (YouTube)

## Technology Stack

| Category | Technology | Version | Purpose |
|----------|------------|---------|---------|
| Language | Go | 1.25.5 | Backend application |
| Web Framework | Gin | v1.11.0 | HTTP routing, middleware |
| JWT | golang-jwt/jwt | v5.2.1 | Token-based authentication |
| YAML | gopkg.in/yaml.v3 | v3.0.1 | Configuration files |
| Image Processing | disintegration/imaging | v1.6.2 | Hero image resizing |

## Architecture Pattern

**Layered Architecture** with clear separation:

```
┌─────────────────────────────────────────────────────────────┐
│                      Handlers Layer                          │
│           (HTTP request handling, routing)                   │
├─────────────────────────────────────────────────────────────┤
│                      Services Layer                          │
│           (Business logic, GOWA integration)                 │
├─────────────────────────────────────────────────────────────┤
│                      Models Layer                            │
│        (Data structures, persistence logic)                  │
├─────────────────────────────────────────────────────────────┤
│                     Utils Layer                              │
│         (Phone formatting, YouTube, logging)                 │
├─────────────────────────────────────────────────────────────┤
│                   Data/Persistence                           │
│           (JSON files with RWMutex locking)                  │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
backend/
├── main.go                 # Entry point, routes, auth middleware
├── config.yaml             # Runtime configuration
├── config/
│   ├── config.go           # Configuration loading
│   └── config_test.go      # Config tests
├── handlers/
│   ├── analytics.go        # Analytics endpoints
│   ├── content.go          # CMS content handlers
│   ├── health.go           # Health check endpoints
│   ├── reminder.go         # Reminder CRUD + send/retry/cancel
│   ├── sse.go              # Server-Sent Events
│   ├── webhook.go          # GOWA webhook processing
│   └── *_test.go           # Handler tests
├── models/
│   ├── content.go          # Article, Video, Category models
│   └── patient.go          # Patient, Reminder models
├── services/
│   ├── gowa.go             # GOWA client with circuit breaker
│   ├── scheduler.go        # Automatic reminder processing
│   └── *_test.go           # Service tests
├── utils/
│   ├── hmac.go             # Webhook signature validation
│   ├── logger.go           # Structured logging
│   ├── mask.go             # PII masking
│   ├── message.go          # WhatsApp message formatting
│   ├── phone.go            # Indonesian phone normalization
│   ├── quiethours.go       # Quiet hours calculation
│   └── youtube.go          # YouTube metadata (noembed)
├── data/                   # JSON persistence
│   ├── patients.json
│   ├── users.json
│   ├── categories.json
│   ├── articles.json
│   ├── videos.json
│   └── jwt_secret.txt
└── uploads/                # Uploaded images
```

## Data Models

### Patient

```go
type Patient struct {
    ID        string           // Format: YYYYMMDDHHMMSS-random8
    Name      string
    Phone     string           // Normalized to 628xxx format
    Email     string
    Notes     string
    Reminders []*Reminder
    CreatedBy string           // User ID of creator
    CreatedAt string           // RFC3339 UTC
    UpdatedAt string           // RFC3339 UTC
}
```

### Reminder

```go
type Reminder struct {
    ID                   string
    Title                string
    Description          string
    DueDate              string           // RFC3339 or 2006-01-02T15:04 format
    Priority             string
    Completed            bool
    Recurrence           Recurrence
    Notified             bool
    Attachments          []Attachment     // Article/Video content

    // Delivery Tracking
    GOWAMessageID        string
    DeliveryStatus       string           // pending, scheduled, sending, sent, delivered, read, failed, retrying, cancelled
    DeliveryErrorMessage string
    MessageSentAt        string
    DeliveredAt          string
    ReadAt               string
    RetryCount           int
    ScheduledDeliveryAt  string
    CancelledAt          string
    CancelledBy          string
}
```

### User

```go
type User struct {
    ID        string
    Username  string           // Unique
    FullName  string
    Password  string           // SHA256 hash
    Role      Role             // "superadmin", "admin", "volunteer"
    CreatedAt string
}
```

### Content

```go
type Article struct {
    ID              string
    Title           string
    Slug            string
    Excerpt         string
    Content         string           // HTML from Quill editor
    AuthorID        string
    CategoryID      string
    HeroImages      HeroImages       // 16x9, 1x1, 4x3 URLs
    Status          ArticleStatus    // "draft", "published"
    Version         int
    ViewCount       int
    AttachmentCount int              // How many reminders this is attached to
    CreatedAt       string
    PublishedAt     string
    UpdatedAt       string
}

type Video struct {
    ID              string
    YouTubeURL      string
    YouTubeID       string
    Title           string
    Description     string
    ChannelName     string
    ThumbnailURL    string
    Duration        string
    CategoryID      string
    Status          VideoStatus
    ViewCount       int
    AttachmentCount int
    CreatedAt       string
    UpdatedAt       string
}
```

## API Endpoints

### Authentication

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/register` | Register new user (default: volunteer) |
| POST | `/api/auth/login` | Login, returns JWT token |
| GET | `/api/auth/me` | Get current user |

### Patient Management

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/patients` | JWT | List patients (volunteers see only own) |
| POST | `/api/patients` | JWT | Create patient |
| GET | `/api/patients/:id` | JWT | Get patient by ID |
| PUT | `/api/patients/:id` | JWT | Update patient |
| DELETE | `/api/patients/:id` | JWT | Delete patient |

### Reminder Management

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/patients/:id/reminders` | JWT | List reminders |
| POST | `/api/patients/:id/reminders` | JWT | Create reminder |
| PUT | `/api/patients/:id/reminders/:reminderId` | JWT | Update reminder |
| POST | `/api/patients/:id/reminders/:reminderId/toggle` | JWT | Toggle completed |
| DELETE | `/api/patients/:id/reminders/:reminderId` | JWT | Delete reminder |
| POST | `/api/patients/:id/reminders/:reminderId/send` | JWT | Send via WhatsApp |
| GET | `/api/reminders/:id/status` | JWT | Get delivery status |
| POST | `/api/reminders/:id/retry` | JWT | Retry failed reminder |
| POST | `/api/reminders/:id/cancel` | JWT | Cancel scheduled |

### Content Management (CMS)

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/categories` | Public | List categories |
| GET | `/api/categories/:type` | Public | Get by type |
| POST | `/api/categories` | Admin+ | Create category |
| GET | `/api/articles` | Public | List published articles |
| GET | `/api/articles/:slug` | Public | Get by slug |
| POST | `/api/articles` | Admin+ | Create article |
| PUT | `/api/articles/:id` | Admin+ | Update article |
| DELETE | `/api/articles/:id` | Admin+ | Delete article |
| GET | `/api/videos` | Public | List published videos |
| POST | `/api/videos` | Admin+ | Create video (YouTube URL) |
| DELETE | `/api/videos/:id` | Admin+ | Delete video |
| GET | `/api/content` | Public | List all content |
| GET | `/api/content/popular` | Public | Popular by attachment count |
| POST | `/api/upload/image` | Admin+ | Upload hero image |

### User Management

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/users` | Superadmin | List all users |
| GET | `/api/users/:id` | JWT | Get user by ID |
| PUT | `/api/users/:id/role` | Superadmin | Update role |
| DELETE | `/api/users/:id` | Superadmin | Delete user |

### Analytics

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/dashboard/stats` | Admin | Dashboard statistics |
| GET | `/api/analytics/content` | Admin | Content attachment stats |
| POST | `/api/analytics/content/sync` | Admin | Sync attachment counts |
| GET | `/api/analytics/delivery` | Admin | Delivery statistics |
| GET | `/api/analytics/failed-deliveries` | Admin | Failed deliveries list |
| GET | `/api/analytics/failed-deliveries/export` | Admin | Export CSV |

### Health & Webhooks

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/api/health` | Public | Basic health check |
| GET | `/api/health/detailed` | Admin | Detailed health |
| POST | `/api/webhook/gowa` | HMAC | GOWA webhook |
| GET | `/api/sse/delivery-status` | JWT | SSE stream |

## Authentication & Authorization

### JWT Authentication

- **Algorithm:** HS256
- **Expiry:** 7 days
- **Secret:** Stored in `data/jwt_secret.txt` (auto-generated)
- **Header:** `Authorization: Bearer <token>`

### Role-Based Access Control

| Role | Permissions |
|------|-------------|
| `superadmin` | Full access, user management, all analytics |
| `admin` | CMS content, dashboard stats, analytics, health details |
| `volunteer` | CRUD own patients/reminders, view public content |

### Middleware

1. `authMiddleware()` - Validates JWT token
2. `sseAuthMiddleware()` - Accepts token from query param (for EventSource)
3. `requireRole(roles...)` - Checks user role

## Key Business Logic

### Reminder Sending Flow

```
1. User triggers send (POST /send)
2. Validate phone (normalize to 628xxx)
3. Check quiet hours (21:00-06:00 WIB)
4. Check circuit breaker status
5. Send to GOWA (/send/message)
6. Update status (sending → sent/failed)
7. Broadcast SSE to connected clients
```

### Automatic Scheduler

Runs every 1 minute, processes:
- **Scheduled reminders:** Quiet hours delivery at 06:00
- **Retrying reminders:** Exponential backoff retry
- **Auto-send:** Past due date (within 24h window)

### Circuit Breaker Pattern

- **States:** `closed` (normal), `open` (blocked)
- **Threshold:** 5 consecutive failures → open
- **Cooldown:** 5 minutes before retry
- **Behavior:** Returns 503 with "queued" flag when open

### GOWA Webhook Processing

1. Validates HMAC-SHA256 signature (`X-Webhook-Signature`)
2. Idempotent processing (tracks `messageId:status` pairs)
3. Updates reminder status: `sent` → `delivered`/`read`/`failed`
4. Broadcasts status via SSE

## Configuration

**File:** `config.yaml`

| Section | Fields |
|---------|--------|
| Server | port, cors_origin |
| GOWA | endpoint, username, password, webhook_secret |
| CircuitBreaker | failure_threshold, cooldown_minutes |
| Retry | max_attempts, delays |
| Logging | level, format |
| Disclaimer | text (Indonesian health disclaimer) |
| QuietHours | start_hour, end_hour, timezone |

## Data Persistence

All data stored as JSON files in `backend/data/`:

| File | Content |
|------|---------|
| `patients.json` | Patients and Reminders |
| `users.json` | User accounts |
| `categories.json` | Content categories |
| `articles.json` | Articles |
| `videos.json` | YouTube videos |
| `jwt_secret.txt` | JWT signing key |

All stores use `sync.RWMutex` for thread-safe access with async saves.

## External Integrations

### GOWA (WhatsApp Gateway)

- **Port:** 3000
- **Auth:** Basic Auth
- **Endpoint:** `/send/message`
- **Webhook:** `/api/webhook/gowa` with HMAC validation

### YouTube (via noembed.com)

- **Usage:** Fetch video metadata from YouTube URL
- **Endpoint:** `https://noembed.com/embed?url=<youtube_url>`
- **Returns:** title, description, thumbnail_url, duration

## Development

### Commands

```bash
# Run development server
cd backend && go run main.go

# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

### Default Credentials

- **Superadmin:** `superadmin` / `superadmin`

---

_Generated using BMAD Method `document-project` workflow_
