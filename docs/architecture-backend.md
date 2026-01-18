# Architecture: Backend (Go/Gin API)

> Generated: 2026-01-18 | Scan Level: Exhaustive

## Executive Summary

The backend is a Go/Gin REST API server providing healthcare volunteer dashboard functionality. It handles patient management, reminder scheduling with WhatsApp delivery via GOWA, content management (articles/videos), and real-time updates via Server-Sent Events (SSE).

## Technology Stack

| Category | Technology | Version | Purpose |
|----------|------------|---------|---------|
| Language | Go | 1.25.5 | Server-side API |
| Framework | Gin | 1.11.0 | HTTP router, middleware |
| Auth | golang-jwt/jwt | v5.2.1 | JWT token management |
| Security | gin-contrib/cors | 1.7.6 | CORS handling |
| Image Processing | disintegration/imaging | 1.6.2 | Hero image generation |
| Config | gopkg.in/yaml.v3 | 3.0.1 | YAML config parsing |
| Data Storage | JSON files | - | File-based persistence |

## Architecture Pattern

**Service-Oriented Monolith** with clean separation:

```
backend/
├── main.go              # Application entry point, routes, auth middleware
├── config/              # Configuration loading (YAML)
├── handlers/            # HTTP request handlers (controllers)
├── models/              # Domain models and data stores
├── services/            # Business logic (GOWA client, scheduler)
├── utils/               # Shared utilities (logging, phone validation)
└── data/                # JSON data files (persistence)
```

## Data Models

### Core Entities

#### Patient (`models/patient.go`)
```go
type Patient struct {
    ID        string      `json:"id"`
    Name      string      `json:"name"`
    Phone     string      `json:"phone"`
    Email     string      `json:"email,omitempty"`
    Notes     string      `json:"notes,omitempty"`
    Reminders []*Reminder `json:"reminders,omitempty"`
    CreatedBy string      `json:"createdBy,omitempty"`
    CreatedAt string      `json:"created_at"`
    UpdatedAt string      `json:"updated_at"`
}
```

#### Reminder (`models/patient.go`)
```go
type Reminder struct {
    ID                   string       `json:"id"`
    Title                string       `json:"title"`
    Description          string       `json:"description"`
    DueDate              string       `json:"dueDate,omitempty"`
    Priority             string       `json:"priority"`
    Completed            bool         `json:"completed"`
    Recurrence           Recurrence   `json:"recurrence"`
    Attachments          []Attachment `json:"attachments,omitempty"`
    DeliveryStatus       string       `json:"delivery_status,omitempty"`
    GOWAMessageID        string       `json:"gowa_message_id,omitempty"`
    DeliveryErrorMessage string       `json:"delivery_error_message,omitempty"`
    RetryCount           int          `json:"retry_count,omitempty"`
    // ... timestamp fields
}
```

**Delivery Status State Machine:**
```
pending → scheduled (quiet hours) → sending → sent → delivered → read
pending → queued → sending → sent → delivered → read
sending → failed
sending → retrying → sending (on transient failure)
any → cancelled (user cancelled)
```

#### Content Models (`models/content.go`)
- **Category**: Content categorization (article/video)
- **Article**: News/educational articles with hero images, slug, status
- **Video**: YouTube video references with metadata

### Data Stores

All stores use `sync.RWMutex` for thread-safe operations:
- `PatientStore`: Patients with nested reminders
- `UserStore`: User accounts with username index
- `CategoryStore`, `ArticleStore`, `VideoStore`: Content management

## API Contracts

### Authentication

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/auth/register` | Register new user | Public |
| POST | `/api/auth/login` | Login, get JWT | Public |
| GET | `/api/auth/me` | Get current user | JWT |

### Patients

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/patients` | List patients | JWT |
| POST | `/api/patients` | Create patient | JWT |
| GET | `/api/patients/:id` | Get patient | JWT |
| PUT | `/api/patients/:id` | Update patient | JWT |
| DELETE | `/api/patients/:id` | Delete patient | JWT |

### Reminders

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/patients/:id/reminders` | List reminders | JWT |
| POST | `/api/patients/:id/reminders` | Create reminder | JWT |
| PUT | `/api/patients/:id/reminders/:rid` | Update reminder | JWT |
| DELETE | `/api/patients/:id/reminders/:rid` | Delete reminder | JWT |
| POST | `/api/patients/:id/reminders/:rid/send` | Send via WhatsApp | JWT |
| POST | `/api/patients/:id/reminders/:rid/toggle` | Toggle completion | JWT |
| GET | `/api/reminders/:id/status` | Get delivery status | JWT |
| POST | `/api/reminders/:id/retry` | Retry failed send | JWT |
| POST | `/api/reminders/:id/cancel` | Cancel pending | JWT |

### Content (CMS)

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/categories` | List all categories | Public |
| GET | `/api/articles` | List articles | Public |
| GET | `/api/articles/:slug` | Get article by slug | Public |
| POST | `/api/articles` | Create article | Admin+ |
| PUT | `/api/articles/:id` | Update article | Admin+ |
| DELETE | `/api/articles/:id` | Delete article | Admin+ |
| GET | `/api/videos` | List videos | Public |
| POST | `/api/videos` | Add video | Admin+ |
| DELETE | `/api/videos/:id` | Delete video | Admin+ |
| POST | `/api/upload/image` | Upload hero image | Admin+ |

### Analytics & Health

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/health` | Basic health check | Public |
| GET | `/api/health/detailed` | Detailed health | Admin+ |
| GET | `/api/analytics/content` | Content statistics | Admin+ |
| GET | `/api/analytics/delivery` | Delivery statistics | Admin+ |
| GET | `/api/analytics/failed-deliveries` | Failed deliveries | Admin+ |

### Real-time Updates

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/sse/delivery-status` | SSE stream | Query token |
| POST | `/api/webhook/gowa` | GOWA webhook | HMAC |

## Authentication & Authorization

### JWT Authentication
- **Token Expiry**: 7 days
- **Algorithm**: HS256
- **Secret**: Stored in `data/jwt_secret.txt`

### Role-Based Access Control (RBAC)

| Role | Permissions |
|------|-------------|
| `superadmin` | Full access, user management |
| `admin` | CMS access, all patient access |
| `volunteer` | Own patients only |

### Middleware Chain
1. `cors.New()` - CORS handling
2. `authMiddleware()` - JWT validation (protected routes)
3. `requireRole()` - Role-based authorization
4. `sseAuthMiddleware()` - Query param token for SSE

## External Integrations

### GOWA (WhatsApp Gateway)
- **Port**: 3000
- **Client**: `services/gowa.go`
- **Features**:
  - Circuit breaker pattern for resilience
  - Retry with exponential backoff
  - Message delivery tracking

### YouTube (noembed.com)
- **Purpose**: Video metadata fetching
- **Fields**: Title, thumbnail, channel name

### Webhook Integration
- **Endpoint**: `/api/webhook/gowa`
- **Auth**: HMAC signature validation
- **Events**: Delivery status updates (sent, delivered, read, failed)

## Configuration

### config.yaml Structure
```yaml
server:
  port: 8080
  cors_origin: "http://localhost:5173"

gowa:
  base_url: "http://localhost:3000"
  device_id: "default"

quiet_hours:
  start_hour: 22
  end_hour: 6
  timezone: "Asia/Jakarta"

retry:
  max_attempts: 3
  delays: [30s, 60s, 120s]

logging:
  level: "info"
  format: "json"

disclaimer:
  enabled: true
  text: "..."
```

## Concurrency & Thread Safety

- All data stores use `sync.RWMutex`
- Lock acquisition order: Store lock → Operation → Unlock
- Async save operations via goroutines
- Graceful shutdown with signal handling

## Entry Points

- **Main**: `backend/main.go:216` - Server initialization
- **Routes**: Defined in `main()` after middleware setup
- **Handlers**: `backend/handlers/*.go`
- **Services**: `backend/services/*.go`

## Testing

- Test files: `*_test.go` alongside source files
- Run: `cd backend && go test ./...`
- Coverage: Unit tests for handlers, services, utils
