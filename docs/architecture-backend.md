# Architecture - Backend (Go/Gin)

**Generated:** January 11, 2026 (Updated)
**Project:** PRIMA Healthcare Volunteer Dashboard
**Technology:** Go 1.25.5 + Gin v1.11.0
**Scan Type:** Exhaustive Rescan

---

## Table of Contents

1. [Overview](#overview)
2. [Layered Architecture](#layered-architecture)
3. [Project Structure](#project-structure)
4. [Routing & Middleware](#routing--middleware)
5. [Data Persistence](#data-persistence)
6. [Services Layer](#services-layer)
7. [Error Handling](#error-handling)
8. [Concurrency & Thread Safety](#concurrency--thread-safety)
9. [Security Architecture](#security-architecture)
10. [Performance Patterns](#performance-patterns)

---

## Overview

### Architecture Style

**Layered Monolith** with clear separation of concerns:

```
┌─────────────────────────────────────────────────────┐
│                   HTTP Layer                        │
│        (Gin Router + Middleware Chain)              │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│                 Handlers Layer                      │
│    (Request/Response, Validation, Auth Check)       │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│                 Services Layer                      │
│  (Business Logic, Circuit Breaker, Scheduler)       │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│                 Models Layer                        │
│     (Entities, Stores, Thread-Safe Operations)      │
└──────────────────┬──────────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────────┐
│               Data Layer (JSON)                     │
│          (File-based Persistence)                   │
└─────────────────────────────────────────────────────┘
```

### Key Design Principles

1. **Single Responsibility** - Each package has one clear purpose
2. **Dependency Injection** - Handlers receive dependencies via constructors
3. **Thread-Safe Operations** - All stores use `sync.RWMutex`
4. **Fail-Fast** - Input validation at handler level
5. **Graceful Degradation** - Circuit breaker for external services

---

## Layered Architecture

### Layer 1: HTTP Layer (Gin Router)

**Location:** `main.go` (router setup)

**Responsibilities:**

- Request routing
- CORS handling
- Static file serving
- Middleware chain execution

**Key Code:**

```go
router := gin.Default()

// CORS middleware
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{appConfig.Server.CORSOrigin},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))

// Static file serving
router.Static("/uploads", "./uploads")
```

### Layer 2: Middleware Chain

**Location:** `main.go` (middleware functions)

**Pipeline:**

```
Request → CORS → Auth → Role Check → Handler → Response
```

**Middleware Types:**

1. **CORS Middleware** (Gin built-in)

   - Applied globally via `cors.New()`
   - Allows `http://localhost:5173` (configurable)

2. **Auth Middleware** (`authMiddleware()`)

   - Extracts JWT from `Authorization: Bearer <token>`
   - Validates token signature and expiry
   - Sets user context (`userID`, `username`, `role`)
   - Returns 401 on failure

3. **SSE Auth Middleware** (`sseAuthMiddleware()`)

   - Accepts token from query parameter (EventSource limitation)
   - Falls back to Authorization header
   - Same validation as `authMiddleware()`

4. **Role Check Middleware** (`requireRole()`)
   - Validates user role against allowed roles
   - Returns 403 if insufficient permissions
   - Used for admin/superadmin-only endpoints

**Example:**

```go
// Public endpoint (no middleware)
router.POST("/api/auth/login", login)

// Protected endpoint (auth only)
api.GET("/patients", getPatients)

// Admin-only endpoint (auth + role check)
api.POST("/articles", requireRole(RoleAdmin, RoleSuperadmin), contentStore.CreateArticle)
```

### Layer 3: Handlers Layer

**Location:** `handlers/` package

**Responsibilities:**

- HTTP request/response handling
- Input validation (Gin binding)
- Authorization checks (role-based)
- Calling services/models
- Error response formatting

**Handlers:**

| Handler                           | File           | Responsibilities                                                 |
| --------------------------------- | -------------- | ---------------------------------------------------------------- |
| **ReminderHandler**               | `reminder.go`  | Create, update, delete, send reminders; delivery status tracking |
| **ContentHandler** (ContentStore) | `content.go`   | CMS for articles/videos/categories; image uploads                |
| **AnalyticsHandler**              | `analytics.go` | Dashboard metrics, delivery logs, failed delivery reports        |
| **HealthHandler**                 | `health.go`    | System health checks, GOWA status, queue statistics              |
| **SSEHandler**                    | `sse.go`       | Server-Sent Events for real-time delivery updates                |
| **WebhookHandler**                | `webhook.go`   | GOWA webhook callbacks with HMAC validation                      |

**Pattern Example (ReminderHandler):**

```go
// handlers/reminder.go
type ReminderHandler struct {
    store        *models.PatientStore // Data access
    config       *config.Config       // Configuration
    gowaClient   *services.GOWAClient // External service
    logger       *slog.Logger         // Logging
    generateID   func() string        // ID generation (injected)
    contentStore *ContentStore        // Dependency on content
}

func NewReminderHandler(
    store *models.PatientStore,
    config *config.Config,
    gowaClient *services.GOWAClient,
    logger *slog.Logger,
    generateID func() string,
    contentStore *ContentStore,
) *ReminderHandler {
    return &ReminderHandler{
        store:        store,
        config:       config,
        gowaClient:   gowaClient,
        logger:       logger,
        generateID:   generateID,
        contentStore: contentStore,
    }
}

func (h *ReminderHandler) Create(c *gin.Context) {
    // 1. Extract path parameters
    patientID := c.Param("id")

    // 2. Bind and validate request body
    var req CreateReminderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "invalid input"})
        return
    }

    // 3. Authorization check (volunteers can only access own patients)
    userRole := c.GetString("role")
    if userRole == "volunteer" {
        // Check patient ownership
    }

    // 4. Business logic via service/model
    reminder := h.store.CreateReminder(patientID, req)

    // 5. Success response
    c.JSON(201, gin.H{"reminder": reminder})
}
```

### Layer 4: Services Layer

**Location:** `services/` package

**Services:**

1. **GOWAClient** (`gowa.go`) - WhatsApp Gateway integration

   - HTTP client with Basic Auth
   - Circuit breaker pattern
   - Retry logic with exponential backoff
   - Connection health tracking

2. **ReminderScheduler** (`scheduler.go`) - Automatic reminder sending
   - Cron-based scheduler (1-minute interval)
   - Quiet hours enforcement
   - Auto-send for due reminders
   - Retry failed deliveries
   - SSE broadcasting integration

**Service Pattern Example:**

```go
// services/gowa.go
type GOWAClient struct {
    endpoint       string
    username       string
    password       string
    httpClient     *http.Client
    circuitBreaker *CircuitBreaker
    logger         *slog.Logger
}

func (c *GOWAClient) SendMessage(phone, message string) (string, error) {
    // 1. Check circuit breaker state
    if c.circuitBreaker.IsOpen() {
        return "", ErrCircuitBreakerOpen
    }

    // 2. Prepare HTTP request
    payload := map[string]string{
        "phone":   phone,
        "message": message,
    }

    // 3. Send with Basic Auth
    resp, err := c.httpClient.Post(c.endpoint + "/send/message", payload)

    // 4. Update circuit breaker state
    if err != nil {
        c.circuitBreaker.RecordFailure()
        return "", err
    }

    c.circuitBreaker.RecordSuccess()
    return resp.MessageID, nil
}
```

### Layer 5: Models Layer

**Location:** `models/` package

**Models:**

1. **Patient** (`patient.go`)

   - Patient entity with embedded Reminders array
   - Delivery status state machine
   - PatientStore with RWMutex for thread-safe operations

2. **Content** (`content.go`)
   - Article, Video, Category entities
   - Stores for each entity type
   - Slug generation, view tracking

**Store Pattern:**

```go
// models/patient.go
type PatientStore struct {
    Patients map[string]*Patient
    mu       sync.RWMutex
    onChange func() // Persistence callback
}

func (s *PatientStore) GetPatient(id string) (*Patient, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    patient, ok := s.Patients[id]
    return patient, ok
}

func (s *PatientStore) UpdatePatient(patient *Patient) {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.Patients[patient.ID] = patient
    if s.onChange != nil {
        s.onChange() // Trigger persistence
    }
}
```

### Layer 6: Data Layer

**Location:** `data/` directory (JSON files)

**Files:**

- `patients.json` - Patients with embedded reminders
- `users.json` - System users (volunteers, admins)
- `articles.json` - Educational articles
- `videos.json` - YouTube videos
- `categories.json` - Content categories
- `jwt_secret.txt` - JWT signing secret

**Persistence Strategy:**

```go
// main.go
func saveData() {
    store.mu.RLock()
    defer store.mu.RUnlock()

    data, err := json.MarshalIndent(store.patients, "", "  ")
    if err != nil {
        log.Printf("Failed to marshal data: %v", err)
        return
    }

    if err := os.WriteFile(dataFile, data, 0644); err != nil {
        log.Printf("Failed to save data: %v", err)
    }
}
```

---

## Project Structure

```
backend/
├── main.go                  # Entry point, routing, bootstrapping (1240 lines)
├── config/
│   ├── config.go            # Configuration loading (YAML)
│   └── config_test.go       # Config tests
├── handlers/
│   ├── analytics.go         # Analytics endpoints
│   ├── analytics_test.go
│   ├── content.go           # CMS endpoints (572 lines)
│   ├── content_test.go
│   ├── health.go            # Health check endpoints
│   ├── health_test.go
│   ├── reminder.go          # Reminder management (745 lines)
│   ├── reminder_test.go
│   ├── sse.go               # Server-Sent Events
│   ├── sse_test.go
│   ├── webhook.go           # GOWA webhooks
│   └── webhook_test.go
├── models/
│   ├── patient.go           # Patient, Reminder entities (200 lines)
│   └── content.go           # Article, Video, Category (150 lines)
├── services/
│   ├── gowa.go              # GOWA client with circuit breaker
│   ├── gowa_test.go
│   ├── scheduler.go         # Reminder scheduler (515 lines)
│   └── scheduler_test.go
├── utils/
│   ├── hmac.go              # HMAC signature validation
│   ├── logger.go            # Structured logging (slog)
│   ├── logger_test.go
│   ├── mask.go              # Data masking (phone, email)
│   ├── mask_test.go
│   ├── message.go           # Message template formatting
│   ├── phone.go             # Phone number validation
│   ├── phone_test.go
│   ├── quiethours.go        # Quiet hours enforcement
│   ├── quiethours_test.go
│   └── youtube.go           # YouTube metadata fetcher
├── data/                    # JSON persistence
│   ├── patients.json
│   ├── users.json
│   ├── articles.json
│   ├── videos.json
│   ├── categories.json
│   └── jwt_secret.txt
├── uploads/                 # User-uploaded images
├── config.yaml              # Runtime configuration
├── config.example.yaml      # Template for configuration
├── go.mod                   # Go module dependencies
└── go.sum                   # Dependency checksums
```

---

## Routing & Middleware

### Route Organization

**Route Groups:**

1. **Public Routes** (no auth)

   ```go
   router.POST("/api/auth/register", register)
   router.POST("/api/auth/login", login)
   router.GET("/api/health", healthHandler.GetHealth)
   router.POST("/api/webhook/gowa", webhookHandler.HandleGOWAWebhook)
   router.GET("/api/config/disclaimer", getDisclaimerConfig)
   router.GET("/api/config/quiet-hours", getQuietHoursConfig)
   ```

2. **Public Content Routes** (read-only CMS)

   ```go
   contentPublic := router.Group("/api")
   contentPublic.GET("/categories", contentStore.ListCategories)
   contentPublic.GET("/articles", contentStore.ListArticles)
   contentPublic.GET("/articles/:slug", contentStore.GetArticle)
   contentPublic.GET("/videos", contentStore.ListVideos)
   contentPublic.GET("/content", contentStore.ListAllContent)
   contentPublic.GET("/content/popular", contentStore.GetPopularContent)
   ```

3. **SSE Route** (query param auth)

   ```go
   router.GET("/api/sse/delivery-status", sseAuthMiddleware(), sseHandler.HandleDeliveryStatusSSE)
   ```

4. **Protected Routes** (auth required)

   ```go
   api := router.Group("/api")
   api.Use(authMiddleware())
   {
       // Current user
       api.GET("/auth/me", getCurrentUser)

       // Patients (volunteers see own, admins see all)
       api.GET("/patients", getPatients)
       api.POST("/patients", createPatient)
       api.GET("/patients/:id", getPatient)
       api.PUT("/patients/:id", updatePatient)
       api.DELETE("/patients/:id", deletePatient)

       // Reminders
       api.POST("/patients/:id/reminders", reminderHandler.Create)
       api.PUT("/patients/:id/reminders/:reminderId", reminderHandler.Update)
       api.POST("/patients/:id/reminders/:reminderId/send", reminderHandler.Send)
       // ... more reminder endpoints
   }
   ```

5. **Admin-Only Routes**

   ```go
   // User management (superadmin only)
   api.GET("/users", requireRole(RoleSuperadmin), getUsers)
   api.PUT("/users/:id/role", requireRole(RoleSuperadmin), updateUserRole)
   api.DELETE("/users/:id", requireRole(RoleSuperadmin), deleteUser)

   // CMS management (admin + superadmin)
   api.POST("/articles", requireRole(RoleAdmin, RoleSuperadmin), contentStore.CreateArticle)
   api.PUT("/articles/:id", requireRole(RoleAdmin, RoleSuperadmin), contentStore.UpdateArticle)
   api.DELETE("/articles/:id", requireRole(RoleAdmin, RoleSuperadmin), contentStore.DeleteArticle)
   api.POST("/videos", requireRole(RoleAdmin, RoleSuperadmin), contentStore.CreateVideo)
   api.POST("/upload/image", requireRole(RoleAdmin, RoleSuperadmin), contentStore.UploadImage)

   // Analytics (admin + superadmin)
   api.GET("/analytics/delivery", requireRole(RoleAdmin, RoleSuperadmin), analyticsHandler.GetDeliveryAnalytics)
   api.GET("/analytics/failed-deliveries", requireRole(RoleAdmin, RoleSuperadmin), analyticsHandler.GetFailedDeliveries)
   ```

### Complete Route List

**Total:** 54 endpoints across 10 categories

See [API Contracts](./api-contracts-backend.md) for detailed documentation of all endpoints.

---

## Data Persistence

### JSON File Strategy

**Why JSON Files?**

1. Simple deployment (no database server required)
2. Human-readable (easy debugging)
3. Version control friendly
4. Atomic writes (replace entire file)

**Trade-offs:**

- ❌ Limited concurrency (in-memory + RWMutex)
- ❌ No complex queries (load all → filter in memory)
- ❌ No transactions (single-file writes)
- ✅ Simple backups (copy `data/` folder)
- ✅ Easy migration (JSON → SQL is straightforward)

### Persistence Flow

```
┌──────────────┐
│  HTTP Request│
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Handler    │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Store.Update │ ◄── Acquires write lock (RWMutex)
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  onChange()  │ ◄── Callback triggered
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  saveData()  │ ◄── Marshal to JSON + write file
└──────────────┘
```

**Code Example:**

```go
// main.go - Initialization
patientStore = models.NewPatientStore(saveData) // Inject persistence callback
patientStore.Patients = store.patients          // Share in-memory map

// models/patient.go - Store operation
func (s *PatientStore) UpdateReminder(patientID, reminderID string, updates func(*Reminder)) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    patient, ok := s.Patients[patientID]
    if !ok {
        return fmt.Errorf("patient not found")
    }

    for i, reminder := range patient.Reminders {
        if reminder.ID == reminderID {
            updates(&patient.Reminders[i]) // Modify in-place

            if s.onChange != nil {
                s.onChange() // Trigger file write
            }
            return nil
        }
    }

    return fmt.Errorf("reminder not found")
}

// main.go - Persistence callback
func saveData() {
    store.mu.RLock()
    defer store.mu.RUnlock()

    data, err := json.MarshalIndent(store.patients, "", "  ")
    if err != nil {
        log.Printf("Failed to marshal data: %v", err)
        return
    }

    // Atomic write: write to temp file, then rename
    tmpFile := dataFile + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0644); err != nil {
        log.Printf("Failed to write temp file: %v", err)
        return
    }

    if err := os.Rename(tmpFile, dataFile); err != nil {
        log.Printf("Failed to rename temp file: %v", err)
        return
    }
}
```

### Data Loading (Startup)

```go
// main.go
func loadData() {
    data, err := os.ReadFile(dataFile)
    if err != nil {
        log.Printf("No existing data file, starting fresh")
        return
    }

    store.mu.Lock()
    defer store.mu.Unlock()

    if err := json.Unmarshal(data, &store.patients); err != nil {
        log.Printf("Failed to load data: %v", err)
    }
}

func loadUsers() {
    data, err := os.ReadFile(usersDataFile)
    if err != nil {
        log.Printf("No existing users file, starting fresh")
        return
    }

    userStore.mu.Lock()
    defer userStore.mu.Unlock()

    var users map[string]*User
    if err := json.Unmarshal(data, &users); err != nil {
        log.Printf("Failed to load users: %v", err)
        return
    }

    userStore.users = users
    for id, user := range users {
        userStore.byName[user.Username] = id
    }
}
```

---

## Services Layer

### GOWAClient (WhatsApp Gateway)

**Location:** `services/gowa.go`

**Features:**

- HTTP client with configurable timeout
- Basic authentication
- Circuit breaker pattern (5 failures → 5min open)
- Exponential backoff retry (1s, 5s, 30s, 2m, 10m)
- Connection health tracking

**Circuit Breaker State Machine:**

```
┌─────────────┐
│   CLOSED    │ ◄─── Normal: All requests allowed
│  (Healthy)  │      Success → stay CLOSED
└──────┬──────┘      5 consecutive failures → OPEN
       │
       ▼
┌─────────────┐
│    OPEN     │ ◄─── Failing: No requests allowed
│  (Circuit   │      Fast-fail with ErrCircuitBreakerOpen
│   Tripped)  │      Wait 5 minutes → HALF_OPEN
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  HALF_OPEN  │ ◄─── Testing: Allow 1 test request
│  (Testing)  │      Success → CLOSED
└──────┬──────┘      Failure → OPEN
       │
       └─────────────► (back to OPEN or CLOSED)
```

**Usage:**

```go
// Initialize
gowaClient := services.NewGOWAClientFromConfig(appConfig, appLogger)

// Send message
messageID, err := gowaClient.SendMessage("628123456789", "Test message")
if err == services.ErrCircuitBreakerOpen {
    // Handle circuit breaker open (schedule retry later)
    return
}
if err != nil {
    // Handle other errors (network, timeout, etc.)
    return
}

// Check availability
if gowaClient.IsAvailable() {
    // Circuit breaker is CLOSED, can send messages
}
```

### ReminderScheduler (Automatic Sending)

**Location:** `services/scheduler.go`

**Features:**

- Runs every 1 minute (configurable)
- Finds all due reminders (`scheduled_time <= now`)
- Enforces quiet hours (8 PM - 8 AM by default)
- Sends via GOWAClient (respects circuit breaker)
- Updates delivery status
- Broadcasts SSE events
- Handles retries for failed deliveries

**Scheduler Flow:**

```
┌──────────────┐
│ Ticker (1min)│
└──────┬───────┘
       │
       ▼
┌────────────────────────────────────┐
│ processScheduledReminders()        │
│ • Scan all patients' reminders     │
│ • Find where scheduled_time <= now │
│ • Filter by enabled=true           │
└──────┬─────────────────────────────┘
       │
       ▼
┌────────────────────────────────────┐
│ For each due reminder:             │
│ 1. Check quiet hours               │
│ 2. Format message (with content)   │
│ 3. Send via GOWAClient             │
│ 4. Update status (sent/failed)     │
│ 5. Broadcast SSE event             │
│ 6. Handle recurrence               │
└────────────────────────────────────┘
```

**Start/Stop:**

```go
// main.go
scheduler := services.NewReminderScheduler(patientStore, gowaClient, appConfig, appLogger)
scheduler.SetSSEHandler(sseHandler) // Link for broadcasting
scheduler.Start() // Runs in background goroutine

// Graceful shutdown
defer scheduler.Stop() // Wait for current cycle to finish
```

---

## Error Handling

### Error Response Format

**Standard Error:**

```json
{
  "error": "Human-readable error message"
}
```

**With Error Code:**

```json
{
  "error": "Invalid phone number format",
  "code": "INVALID_PHONE"
}
```

### HTTP Status Codes

| Code                          | Usage                    | Example                         |
| ----------------------------- | ------------------------ | ------------------------------- |
| **200 OK**                    | Successful GET/PUT       | Patient retrieved               |
| **201 Created**               | Successful POST          | Reminder created                |
| **400 Bad Request**           | Invalid input            | Missing required field          |
| **401 Unauthorized**          | Missing/invalid token    | Token expired                   |
| **403 Forbidden**             | Insufficient permissions | Volunteer accessing admin route |
| **404 Not Found**             | Resource not found       | Patient ID doesn't exist        |
| **409 Conflict**              | Duplicate resource       | Username already taken          |
| **500 Internal Server Error** | Unexpected error         | JSON marshal failed             |
| **503 Service Unavailable**   | External service down    | Circuit breaker open            |

### Error Handling Pattern

```go
// handlers/reminder.go
func (h *ReminderHandler) Send(c *gin.Context) {
    patientID := c.Param("id")
    reminderID := c.Param("reminderId")

    // 1. Validate input
    patient, ok := h.store.GetPatient(patientID)
    if !ok {
        c.JSON(404, gin.H{"error": "patient not found"})
        return
    }

    reminder := findReminder(patient, reminderID)
    if reminder == nil {
        c.JSON(404, gin.H{"error": "reminder not found"})
        return
    }

    // 2. Business logic validation
    if !reminder.Enabled {
        c.JSON(400, gin.H{
            "error": "cannot send disabled reminder",
            "code":  "REMINDER_DISABLED",
        })
        return
    }

    // 3. External service call
    messageID, err := h.gowaClient.SendMessage(patient.Phone, reminder.Message)
    if err == services.ErrCircuitBreakerOpen {
        c.JSON(503, gin.H{
            "error": "WhatsApp service temporarily unavailable",
            "code":  "CIRCUIT_BREAKER_OPEN",
        })
        return
    }
    if err != nil {
        h.logger.Error("Failed to send message",
            "patient_id", patientID,
            "reminder_id", reminderID,
            "error", err,
        )
        c.JSON(500, gin.H{"error": "failed to send message"})
        return
    }

    // 4. Success
    c.JSON(200, gin.H{
        "message_id": messageID,
        "status":     "sent",
    })
}
```

### Logging

**Structured Logging (slog):**

```go
// utils/logger.go
func InitDefaultLogger(level, format string) {
    var handler slog.Handler

    if format == "json" {
        handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: parseLevel(level),
        })
    } else {
        handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
            Level: parseLevel(level),
        })
    }

    DefaultLogger = slog.New(handler)
}

// Usage in handlers
h.logger.Info("Reminder sent",
    "reminder_id", reminderID,
    "patient_id", patientID,
    "phone", utils.MaskPhone(phone), // Masked for privacy
    "status", "sent",
    "message_id", messageID,
)

h.logger.Error("GOWA send failed",
    "reminder_id", reminderID,
    "error", err.Error(),
    "retry_count", retryCount,
)
```

**Log Levels:**

- `DEBUG` - Development details (verbose)
- `INFO` - Normal operations (default)
- `WARN` - Recoverable errors (circuit breaker open, retry scheduled)
- `ERROR` - Critical failures (GOWA unreachable, data corruption)

---

## Concurrency & Thread Safety

### RWMutex Pattern

**Why RWMutex?**

- Multiple concurrent readers (GET requests)
- Single writer at a time (POST/PUT/DELETE)
- Better performance than regular `sync.Mutex` for read-heavy workloads

**Usage:**

```go
// models/patient.go
type PatientStore struct {
    Patients map[string]*Patient
    mu       sync.RWMutex
    onChange func()
}

// Read operation (shared lock)
func (s *PatientStore) GetPatient(id string) (*Patient, bool) {
    s.mu.RLock() // Multiple goroutines can hold RLock simultaneously
    defer s.mu.RUnlock()

    patient, ok := s.Patients[id]
    return patient, ok
}

// Write operation (exclusive lock)
func (s *PatientStore) UpdatePatient(patient *Patient) {
    s.mu.Lock() // Only one goroutine can hold Lock at a time
    defer s.mu.Unlock()

    s.Patients[patient.ID] = patient
    if s.onChange != nil {
        s.onChange() // Trigger persistence
    }
}
```

### Goroutine Safety Checklist

✅ **Thread-Safe:**

- All store operations (protected by RWMutex)
- JWT token generation/verification (stateless)
- Handler functions (Gin handles concurrency)
- Logger (slog is thread-safe)

⚠️ **Requires Lock:**

- Modifying shared maps (`store.patients`, `userStore.users`)
- Reading/writing shared state (circuit breaker state)

❌ **Not Thread-Safe:**

- Direct map access without lock
- SSE client map modifications (uses separate mutex)

### Goroutine Usage

**Background Tasks:**

1. **Reminder Scheduler** (`scheduler.Start()`)

   - Runs every 1 minute
   - Scans for due reminders
   - Sends via GOWAClient

2. **Health Check Goroutine**

   - Runs every 60 seconds
   - Pings GOWA for connectivity
   - Updates health handler state

3. **SSE Connections**
   - One goroutine per SSE client
   - Listens for broadcast events
   - Writes to client stream

**Graceful Shutdown:**

```go
// main.go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

appLogger.Info("Shutting down server...")

// 1. Stop scheduler first (no new sends)
scheduler.Stop()

// 2. Close all SSE connections
sseHandler.Shutdown()

// 3. Shutdown HTTP server with 5s timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
srv.Shutdown(ctx)

appLogger.Info("Server exited gracefully")
```

---

## Security Architecture

### Authentication (JWT)

**Token Generation:**

```go
func generateToken(userID, username string, role Role) (string, error) {
    secret := getJWTSecret() // From data/jwt_secret.txt

    claims := &Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

**Token Verification:**

```go
func verifyToken(tokenString string) (*Claims, error) {
    secret := getJWTSecret()

    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, fmt.Errorf("invalid token")
}
```

**Security Properties:**

- Algorithm: HS256 (HMAC with SHA-256)
- Secret: Random Base64 string, persisted in file
- Expiry: 7 days (configurable via `tokenExpiry`)
- Stateless: No server-side session storage

### Authorization (RBAC)

**Roles:**

- `superadmin` - Full access (user management, CMS, analytics)
- `admin` - CMS management, analytics (no user management)
- `volunteer` - Patient/reminder management (own patients only)

**Enforcement:**

```go
// Middleware-level (route protection)
api.GET("/users", requireRole(RoleSuperadmin), getUsers)

// Handler-level (data filtering)
func getPatients(c *gin.Context) {
    userRole := Role(c.GetString("role"))
    userID := c.GetString("userID")

    if userRole == RoleVolunteer {
        // Filter: only return patients assigned to this volunteer
        return getVolunteerPatients(userID)
    }

    // Admin/superadmin: return all patients
    return getAllPatients()
}
```

### Password Security

**Hashing:**

```go
func hashPassword(password string) string {
    hash := sha256.Sum256([]byte(password))
    return base64.StdEncoding.EncodeToString(hash[:])
}
```

**Verification:**

```go
func verifyPassword(password, hash string) bool {
    return hashPassword(password) == hash
}
```

**Properties:**

- Algorithm: SHA-256 (one-way hash)
- Encoding: Base64
- No salt (acceptable for low-risk internal system)
- **Production Recommendation:** Use bcrypt or Argon2 for better security

### HMAC Webhook Validation

**Purpose:** Verify GOWA webhooks are authentic

**Implementation:**

```go
// utils/hmac.go
func ValidateHMAC(secret, payload, signature string) bool {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(payload))
    expectedMAC := hex.EncodeToString(mac.Sum(nil))

    return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// handlers/webhook.go
func (h *WebhookHandler) HandleGOWAWebhook(c *gin.Context) {
    signature := c.GetHeader("X-Webhook-Signature")
    if signature == "" {
        c.JSON(401, gin.H{"error": "missing signature"})
        return
    }

    body, _ := io.ReadAll(c.Request.Body)

    if !utils.ValidateHMAC(h.config.GOWA.WebhookSecret, string(body), signature) {
        c.JSON(401, gin.H{"error": "invalid signature"})
        return
    }

    // Process webhook
}
```

### Data Masking

**Phone Numbers:**

```go
// utils/mask.go
func MaskPhone(phone string) string {
    if len(phone) <= 6 {
        return phone
    }
    return phone[:6] + "***" + phone[len(phone)-3:]
}
// 628123456789 → 6281234***789
```

**Emails:**

```go
func MaskEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return email
    }

    local := parts[0]
    if len(local) > 1 {
        local = string(local[0]) + "***"
    }

    return local + "@" + parts[1]
}
// john.doe@example.com → j***@example.com
```

**Usage:** Logs, admin analytics views (privacy protection)

---

## Performance Patterns

### Circuit Breaker Benefits

**Problem:** GOWA service becomes unavailable

- Without circuit breaker: Every request waits for timeout (30s)
- 100 pending reminders × 30s = 50 minutes of blocked operations

**Solution:** Circuit breaker pattern

- After 5 consecutive failures, open circuit
- Fast-fail for next 5 minutes (no network calls)
- Reminders marked as `failed` and scheduled for retry
- Backend remains responsive

**Implementation:**

```go
type CircuitBreaker struct {
    failureThreshold int
    timeout          time.Duration
    failureCount     int
    lastFailureTime  time.Time
    state            CircuitBreakerState // CLOSED, OPEN, HALF_OPEN
    mu               sync.RWMutex
}

func (cb *CircuitBreaker) IsOpen() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()

    if cb.state == CircuitBreakerStateOpen {
        // Check if timeout expired → transition to HALF_OPEN
        if time.Since(cb.lastFailureTime) > cb.timeout {
            return false // Allow one test request
        }
        return true
    }

    return false
}
```

### Connection Pooling (HTTP Client)

**Default HTTP Client Issues:**

- Creates new TCP connection for each request
- No connection reuse
- Higher latency

**Solution:** Reusable HTTP client with connection pool

```go
func NewGOWAClientFromConfig(cfg *config.Config, logger *slog.Logger) *GOWAClient {
    return &GOWAClient{
        endpoint: cfg.GOWA.Endpoint,
        username: cfg.GOWA.Username,
        password: cfg.GOWA.Password,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        100,              // Total idle connections
                MaxIdleConnsPerHost: 10,               // Per-host idle connections
                IdleConnTimeout:     90 * time.Second, // Keep alive duration
            },
        },
        circuitBreaker: NewCircuitBreaker(5, 5*time.Minute),
        logger:         logger,
    }
}
```

### SSE Broadcasting Optimization

**Naive Approach:** Store sends to all clients in loop

- O(n) iteration for each update
- Blocked if any client is slow

**Optimized Approach:** Channel-based broadcasting

```go
// handlers/sse.go
type SSEHandler struct {
    clients   map[string]chan Event // clientID → event channel
    broadcast chan Event            // Broadcast channel
    mu        sync.RWMutex
}

// Client goroutine listens to its channel
func (h *SSEHandler) HandleDeliveryStatusSSE(c *gin.Context) {
    clientID := generateID()
    eventChan := make(chan Event, 10) // Buffered channel

    h.addClient(clientID, eventChan)
    defer h.removeClient(clientID)

    for {
        select {
        case event := <-eventChan:
            c.SSEvent(event.Type, event.Data)
            c.Writer.Flush()
        case <-c.Request.Context().Done():
            return // Client disconnected
        }
    }
}

// Broadcast goroutine distributes events to all clients
func (h *SSEHandler) broadcastLoop() {
    for event := range h.broadcast {
        h.mu.RLock()
        for _, eventChan := range h.clients {
            select {
            case eventChan <- event: // Non-blocking send
            default: // Channel full, skip (client too slow)
            }
        }
        h.mu.RUnlock()
    }
}

// External code broadcasts events
func (h *SSEHandler) BroadcastDeliveryStatusUpdate(reminderID, status, timestamp string) {
    h.broadcast <- Event{
        Type: "delivery.status.updated",
        Data: fmt.Sprintf(`{"reminder_id": "%s", "status": "%s", "timestamp": "%s"}`, reminderID, status, timestamp),
    }
}
```

### Quiet Hours Enforcement

**Problem:** Don't send reminders late at night (8 PM - 8 AM)

**Solution:** Time zone-aware scheduling

```go
// utils/quiethours.go
func IsQuietHours(t time.Time, quietStart, quietEnd int) bool {
    // Convert to local time (server's time zone)
    hour := t.Local().Hour()

    if quietEnd > quietStart {
        // Normal range: e.g., 22:00 - 08:00 (across midnight)
        return hour >= quietStart || hour < quietEnd
    }

    // Inverted range (shouldn't happen, but handle it)
    return hour >= quietStart && hour < quietEnd
}

// services/scheduler.go
func (s *ReminderScheduler) processScheduledReminders() {
    now := time.Now().UTC()

    for _, reminder := range dueReminders {
        if utils.IsQuietHours(now, s.config.QuietHours.Start, s.config.QuietHours.End) {
            s.logger.Info("Skipping reminder send (quiet hours)",
                "reminder_id", reminder.ID,
                "scheduled_time", reminder.ScheduledTime,
            )
            continue // Skip, will be checked again next cycle
        }

        // Send reminder
    }
}
```

---

**Next:** See [Architecture - Frontend](./architecture-frontend.md) for frontend architecture details.
