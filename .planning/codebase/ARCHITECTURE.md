# Architecture

**Analysis Date:** 2026-01-17

## Pattern Overview

**Overall:** Layered Monolithic Architecture with In-Memory Store and File Persistence

**Key Characteristics:**
- Clear separation between frontend (Svelte 5 SPA) and backend (Go/Gin REST API)
- In-memory data stores with `sync.RWMutex` for thread-safe operations
- Asynchronous file persistence to JSON for data durability
- JWT-based authentication with Role-Based Access Control (RBAC)
- Real-time updates via Server-Sent Events (SSE)
- Circuit breaker pattern for external service resilience
- Background scheduler for reminder auto-sending

## Layers

**1. Backend - Presentation Layer (Gin HTTP Handlers):**
- Location: `backend/handlers/`
- Handles HTTP requests/responses, authentication middleware, request validation
- Routes defined in `backend/main.go`
- Key handlers: `ReminderHandler`, `ContentHandler`, `AnalyticsHandler`, `SSEHandler`, `WebhookHandler`

**2. Backend - Service Layer:**
- Location: `backend/services/`
- `GOWAClient` - WhatsApp gateway integration with circuit breaker
- `ReminderScheduler` - Background job for scheduled reminder delivery
- Business logic coordination, external API calls, retry logic

**3. Backend - Model/Store Layer:**
- Location: `backend/models/`
- `PatientStore`, `UserStore`, `ArticleStore`, `VideoStore`, `CategoryStore`
- Thread-safe in-memory storage with `sync.RWMutex`
- File persistence via JSON (asynchronous writes)

**4. Backend - Utility Layer:**
- Location: `backend/utils/`
- Phone validation, YouTube metadata fetching, message formatting, logging

**5. Frontend - View Layer:**
- Location: `frontend/src/lib/views/`
- Svelte 5 components for major views: `DashboardView`, `PatientsView`, `CMSDashboardView`, etc.

**6. Frontend - Component Layer:**
- Location: `frontend/src/lib/components/`
- Reusable UI components organized by domain: `patients/`, `reminders/`, `content/`, `delivery/`, `ui/`

**7. Frontend - State/Service Layer:**
- Location: `frontend/src/lib/stores/` and `frontend/src/lib/services/`
- `delivery.svelte.js` - Svelte 5 rune-based store for real-time delivery status
- `sse.js` - SSE client service with reconnection logic
- `api.js` - API client functions

## Data Flow

**Patient and Reminder Creation Flow:**

1. Frontend (`App.svelte`) calls `api.saveReminder()` in `frontend/src/lib/utils/api.js`
2. API client sends POST to `/api/patients/:id/reminders`
3. `ReminderHandler.Create()` in `backend/handlers/reminder.go` validates request
4. Handler acquires write lock on `PatientStore`, appends reminder to patient
5. Handler releases lock, triggers async `store.SaveData()`
6. Response returns created reminder to frontend

**Real-Time Delivery Status Update Flow:**

1. User triggers reminder send via `SendReminderModal`
2. `ReminderHandler.Send()` sends message via `GOWAClient`
3. On success, handler calls `sseHandler.BroadcastDeliveryStatusUpdate()`
4. SSE handler broadcasts to all connected clients via `SSEHandler.BroadcastDeliveryStatusUpdate()`
5. Frontend `sseService` receives event, updates `deliveryStore`
6. Svelte 5 reactivity updates UI automatically

**Scheduled Reminder Auto-Send Flow:**

1. `ReminderScheduler.run()` in `backend/services/scheduler.go` checks every minute
2. Scans for reminders with `DeliveryStatusScheduled` or `DeliveryStatusRetrying`
3. Calls `sendScheduledReminder()` or `processRetryReminder()`
4. Formats message with content attachments
5. Sends via `GOWAClient.SendMessage()`
6. Updates delivery status and broadcasts via SSE

**Authentication Flow:**

1. User submits credentials via `LoginScreen.svelte`
2. Backend `login()` handler in `backend/main.go` verifies password hash
3. Generates JWT token with 7-day expiry
4. Frontend stores token in `localStorage`
5. Subsequent requests include `Authorization: Bearer <token>` header
6. `authMiddleware()` in `backend/main.go` validates JWT claims
7. RBAC middleware (`requireRole()`) checks user role for protected routes

## Key Abstractions

**PatientStore (In-Memory Thread-Safe Storage):**
- Location: `backend/models/patient.go`
- Wraps map with `sync.RWMutex`
- Exposes `Lock()`, `Unlock()`, `RLock()`, `RUnlock()` methods
- `SaveFunc()` callback for persistence
- Pattern used for all domain stores

**ContentStore (Centralized CMS Data):**
- Location: `backend/handlers/content.go`
- Aggregates `CategoryStore`, `ArticleStore`, `VideoStore`
- Provides unified access for content CRUD operations
- Handles author name resolution via user store

**Circuit Breaker (Resilience Pattern):**
- Location: `backend/services/gowa.go`
- Protects GOWA WhatsApp service from cascade failures
- States: `closed`, `open`, `cooldown`
- Configurable failure threshold and cooldown duration
- Prevents request flooding when service is down

**SSE Handler (Real-Time Communication):**
- Location: `backend/handlers/sse.go`
- Manages persistent connections to clients
- Broadcasts to all connected clients via channel
- Automatic cleanup on client disconnect or server shutdown

**Svelte 5 Rune-Based State:**
- Location: `frontend/src/lib/stores/delivery.svelte.js`
- Uses `$state()` for reactive state
- Uses `$derived()` for computed values
- Subscribes to SSE events via `sseService`

## Entry Points

**Backend Entry Point:**
- Location: `backend/main.go`
- Initializes all stores, handlers, services
- Configures Gin router with middleware and routes
- Starts scheduler goroutine
- Sets up graceful shutdown signal handling

**Frontend Entry Point:**
- Location: `frontend/src/App.svelte`
- Single-page application with client-side routing via conditional views
- Initializes auth state, SSE connection, and stores
- All views are conditional components, not separate pages

**SSE Endpoint:**
- Route: `GET /api/sse/delivery-status`
- Location: `backend/handlers/sse.go` - `HandleDeliveryStatusSSE()`
- Auth: Query parameter token (EventSource API limitation)
- Streams delivery status updates to connected clients

## Error Handling

**Strategy:** Consistent JSON error responses with HTTP status codes

**Patterns:**
- Validation errors: `400 Bad Request` with error message
- Authentication errors: `401 Unauthorized`
- Authorization errors: `403 Forbidden` (RBAC)
- Not found: `404 Not Found`
- Service unavailable: `503 Service Unavailable` (circuit breaker open)

**Circuit Breaker Pattern:**
- `GOWAClient` blocks requests when circuit is open
- Failed deliveries queued for retry when service recovers
- Exponential backoff for retry scheduling

## Cross-Cutting Concerns

**Authentication:** JWT tokens with 7-day expiry, stored in localStorage
- Middleware: `authMiddleware()`, `sseAuthMiddleware()`, `requireRole()`

**Logging:** Structured JSON logging via `slog`
- Configurable log level via `config.yaml`
- All handlers and services use injected logger

**Validation:**
- Go struct binding with `binding:"required"` tags
- Phone number validation via `utils.ValidatePhoneNumber()`
- Attachment content validation in `ReminderHandler`

**Concurrency:**
- `sync.RWMutex` for read-heavy workloads
- Goroutines for background jobs (scheduler, health checks)
- Async file writes to avoid blocking requests

---

*Architecture analysis: 2026-01-17*
