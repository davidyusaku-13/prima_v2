# Development Guide - Backend (Go/Gin)

**Generated:** January 11, 2026 (Updated)
**Project:** PRIMA Healthcare Volunteer Dashboard
**Technology:** Go 1.25.5 + Gin v1.11.0
**Scan Type:** Exhaustive Rescan

---

## Quick Start

```bash
# Navigate to backend directory
cd backend

# Run the server (port 8080)
go run main.go

# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Build production binary
go build -o prima-backend main.go
```

---

## Environment Setup

### Prerequisites

- **Go:** 1.21+ (tested with 1.25.5)
- **Git:** For version control
- **Text Editor:** VS Code with Go extension recommended

### Install Go

**Windows:**

```powershell
# Using Chocolatey
choco install golang

# Or download from https://go.dev/dl/
```

**macOS:**

```bash
# Using Homebrew
brew install go
```

**Linux:**

```bash
# Download and install
wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

**Verify Installation:**

```bash
go version
# Should output: go version go1.25.5 <OS>/<ARCH>
```

### Clone Repository

```bash
git clone <repository-url>
cd prima_v2/backend
```

### Install Dependencies

```bash
go mod download
```

This reads `go.mod` and downloads all dependencies listed.

---

## Project Structure

```
backend/
├── main.go                  # Entry point (1240 lines)
├── config.yaml              # Runtime configuration
├── config.example.yaml      # Template for configuration
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── config/                  # Configuration loading
│   ├── config.go
│   └── config_test.go
├── handlers/                # HTTP request handlers
│   ├── analytics.go
│   ├── content.go
│   ├── health.go
│   ├── reminder.go
│   ├── sse.go
│   ├── webhook.go
│   └── *_test.go
├── models/                  # Data entities
│   ├── patient.go
│   └── content.go
├── services/                # Business logic
│   ├── gowa.go             # WhatsApp Gateway client
│   ├── scheduler.go        # Reminder scheduler
│   └── *_test.go
├── utils/                   # Utility functions
│   ├── hmac.go
│   ├── logger.go
│   ├── mask.go
│   ├── message.go
│   ├── phone.go
│   ├── quiethours.go
│   ├── youtube.go
│   └── *_test.go
├── data/                    # JSON persistence (gitignored)
│   ├── patients.json
│   ├── users.json
│   ├── articles.json
│   ├── videos.json
│   ├── categories.json
│   └── jwt_secret.txt
└── uploads/                 # User-uploaded images (gitignored)
```

---

## Configuration

### `config.yaml`

**Location:** `backend/config.yaml`

**Template:** Copy from `config.example.yaml`

```yaml
# Server configuration
server:
  port: 8080
  cors_origin: "http://localhost:5173"

# GOWA WhatsApp Gateway
gowa:
  endpoint: "http://localhost:3000"
  username: "admin"
  password: "admin"
  webhook_secret: "your-secret-key-here"
  timeout: 30 # seconds

# Reminder scheduler
scheduler:
  interval: 60 # seconds (check every 1 minute)

# Quiet hours (local time)
quiet_hours:
  start: 20 # 8 PM
  end: 8 # 8 AM

# Logging
logging:
  level: "info" # debug, info, warn, error
  format: "text" # text or json

# Disclaimer
disclaimer:
  enabled: true
  version: "1.0"
  text: "PRIMA adalah aplikasi pengelolaan pasien dan pengingat kesehatan..."
```

**Environment Variable Override:**

Create `.env` file (optional):

```env
PORT=8080
CORS_ORIGIN=http://localhost:5173
GOWA_ENDPOINT=http://localhost:3000
GOWA_USERNAME=admin
GOWA_PASSWORD=admin
LOG_LEVEL=info
```

---

## Running Locally

### 1. Start Backend

```bash
cd backend
go run main.go
```

**Expected Output:**

```
2026/01/02 10:00:00 🚀 Server started on port 8080
```

**Access:**

- Health check: http://localhost:8080/api/health
- API base: http://localhost:8080/api

### 2. With GOWA (WhatsApp Gateway)

GOWA must be running separately for message sending functionality.

**Start GOWA** (see [GOWA-README.md](../GOWA-README.md)):

```bash
# In separate terminal
cd gowa
docker-compose up
```

**Verify GOWA:**

```bash
curl http://localhost:3000/health
```

### 3. Initialize Default Superadmin

On first run, a default superadmin user is created:

- **Username:** `superadmin`
- **Password:** `superadmin`

⚠️ **IMPORTANT:** Change this password in production!

### 4. Test Backend API

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "superadmin", "password": "superadmin"}'

# Response:
# {"token": "eyJhbGciOiJIUzI1NiIs...", "user": {...}}

# Use token for authenticated requests
TOKEN="your-token-here"

# Get patients
curl http://localhost:8080/api/patients \
  -H "Authorization: Bearer $TOKEN"
```

---

## Testing

### Run All Tests

```bash
go test ./...
```

**Output:**

```
ok  	github.com/davidyusaku-13/prima_v2/config	0.123s
ok  	github.com/davidyusaku-13/prima_v2/handlers	0.456s
ok  	github.com/davidyusaku-13/prima_v2/services	0.789s
ok  	github.com/davidyusaku-13/prima_v2/utils	0.234s
```

### Run Specific Package

```bash
go test ./config
go test ./handlers
go test ./utils
```

### Run Specific Test

```bash
go test -v -run TestLoadConfig ./config
go test -v -run TestMaskPhone ./utils
```

### Run with Coverage

```bash
# All packages
go test -cover ./...

# Specific package with verbose coverage
go test -v -cover ./config

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

**Open `coverage.html` in browser to see detailed coverage.**

### Test Pattern

**File:** `config/config_test.go`

```go
package config

import (
    "os"
    "testing"
)

func TestLoadConfig(t *testing.T) {
    // Create temporary config file
    tmpFile, err := os.CreateTemp("", "config-*.yaml")
    if err != nil {
        t.Fatalf("Failed to create temp file: %v", err)
    }
    defer os.Remove(tmpFile.Name())

    // Write test config
    configData := `
server:
  port: 9090
  cors_origin: "http://example.com"
`
    if _, err := tmpFile.WriteString(configData); err != nil {
        t.Fatalf("Failed to write config: %v", err)
    }
    tmpFile.Close()

    // Load config
    cfg := Load(tmpFile.Name())

    // Assert
    if cfg.Server.Port != 9090 {
        t.Errorf("Expected port 9090, got %d", cfg.Server.Port)
    }
    if cfg.Server.CORSOrigin != "http://example.com" {
        t.Errorf("Expected cors_origin http://example.com, got %s", cfg.Server.CORSOrigin)
    }
}

func TestDefault(t *testing.T) {
    cfg := Default()

    if cfg.Server.Port != 8080 {
        t.Errorf("Expected default port 8080, got %d", cfg.Server.Port)
    }
}
```

---

## Code Style

### Formatting

**Always use `gofmt`:**

```bash
# Format all files
gofmt -w .

# Check formatting without modifying
gofmt -d .
```

**VS Code:** Install Go extension and enable "Format on Save"

### Naming Conventions

**Exported (public):**

- Types: `PascalCase` → `type Patient struct{}`
- Functions: `PascalCase` → `func NewPatientStore()`
- Constants: `PascalCase` → `const RoleSuperadmin`

**Unexported (private):**

- Variables: `camelCase` → `var userID string`
- Functions: `camelCase` → `func verifyToken()`
- Parameters: `camelCase` → `func GetPatient(patientID string)`

**Acronyms:**

- `HTTPResponse` not `HttpResponse`
- `parseURL` not `parseUrl`
- `userID` not `userId`

**Interfaces:**

- Single method: `Reader`, `Writer`, `Closer`
- Multiple methods: `PatientStore`, `RemindHandler`

### Error Handling

**Always check errors:**

```go
data, err := os.ReadFile("config.yaml")
if err != nil {
    return fmt.Errorf("failed to read config: %w", err)
}
```

**Wrap errors with context:**

```go
err := savePatient(patient)
if err != nil {
    return fmt.Errorf("failed to save patient %s: %w", patient.ID, err)
}
```

**Use `%w` for error wrapping (Go 1.13+):**

```go
// Good (preserves error chain)
return fmt.Errorf("failed to load config: %w", err)

// Bad (breaks error chain)
return fmt.Errorf("failed to load config: %v", err)
```

### Struct Tags

**JSON tags:**

```go
type Patient struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Phone string `json:"phone,omitempty"`  // Omit if empty
}
```

**YAML tags:**

```go
type Config struct {
    Port int    `yaml:"port"`
    Host string `yaml:"host"`
}
```

### Comments

**Package comment:**

```go
// Package handlers implements HTTP request handlers for the PRIMA API.
package handlers
```

**Exported function:**

```go
// NewPatientStore creates a new PatientStore with the given onChange callback.
// The callback is triggered whenever data is modified.
func NewPatientStore(onChange func()) *PatientStore {
    return &PatientStore{
        Patients: make(map[string]*Patient),
        onChange: onChange,
    }
}
```

---

## Adding New Features

### Add New Endpoint

**1. Define Handler Function**

`handlers/example.go`:

```go
package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type ExampleHandler struct {
    store *models.PatientStore
}

func NewExampleHandler(store *models.PatientStore) *ExampleHandler {
    return &ExampleHandler{store: store}
}

func (h *ExampleHandler) GetExample(c *gin.Context) {
    id := c.Param("id")

    // Business logic
    data, ok := h.store.GetPatient(id)
    if !ok {
        c.JSON(404, gin.H{"error": "not found"})
        return
    }

    c.JSON(200, gin.H{"data": data})
}
```

**2. Register Route**

`main.go`:

```go
// Initialize handler
exampleHandler := handlers.NewExampleHandler(patientStore)

// Register route (protected)
api := router.Group("/api")
api.Use(authMiddleware())
{
    api.GET("/example/:id", exampleHandler.GetExample)
}
```

**3. Add Test**

`handlers/example_test.go`:

```go
package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/gin-gonic/gin"
)

func TestGetExample(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    handler := NewExampleHandler(mockStore)
    router.GET("/example/:id", handler.GetExample)

    // Test
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/example/123", nil)
    router.ServeHTTP(w, req)

    // Assert
    if w.Code != 200 {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

### Add New Model

**1. Define Entity**

`models/example.go`:

```go
package models

import "sync"

type Example struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    CreatedAt string `json:"createdAt"`
}

type ExampleStore struct {
    Examples map[string]*Example
    mu       sync.RWMutex
    onChange func()
}

func NewExampleStore(onChange func()) *ExampleStore {
    return &ExampleStore{
        Examples: make(map[string]*Example),
        onChange: onChange,
    }
}

func (s *ExampleStore) Get(id string) (*Example, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    example, ok := s.Examples[id]
    return example, ok
}

func (s *ExampleStore) Save(example *Example) {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.Examples[example.ID] = example
    if s.onChange != nil {
        s.onChange()
    }
}
```

**2. Add Persistence**

`main.go`:

```go
const examplesDataFile = "data/examples.json"

var exampleStore *models.ExampleStore

func loadExamples() {
    data, err := os.ReadFile(examplesDataFile)
    if err != nil {
        return
    }

    var examples map[string]*models.Example
    json.Unmarshal(data, &examples)
    exampleStore.Examples = examples
}

func saveExamples() {
    exampleStore.mu.RLock()
    defer exampleStore.mu.RUnlock()

    data, _ := json.MarshalIndent(exampleStore.Examples, "", "  ")
    os.WriteFile(examplesDataFile, data, 0644)
}

// In main()
exampleStore = models.NewExampleStore(saveExamples)
loadExamples()
```

---

## Debugging

### Enable Debug Logging

**config.yaml:**

```yaml
logging:
  level: "debug"
  format: "json"
```

**Code:**

```go
logger.Debug("Processing patient",
    "patient_id", patientID,
    "reminder_count", len(patient.Reminders),
)
```

### Use Delve Debugger

**Install:**

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

**Debug with breakpoints:**

```bash
dlv debug main.go
```

**In Delve:**

```
(dlv) break main.main
(dlv) continue
(dlv) print patientID
(dlv) next
```

**VS Code:** Use built-in debugger with `launch.json`:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Backend",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend",
      "args": []
    }
  ]
}
```

### Check GOWA Connectivity

```bash
# Health check
curl http://localhost:8080/api/health/detailed \
  -H "Authorization: Bearer $TOKEN"

# Response includes GOWA status:
# {
#   "gowa": {
#     "connected": true,
#     "last_ping": "2026-01-02T10:00:00Z"
#   },
#   "circuit_breaker": {
#     "state": "closed",
#     "failure_count": 0
#   }
# }
```

---

## Common Issues

### Port Already in Use

**Error:** `bind: address already in use`

**Solution:**

```bash
# Find process using port 8080
lsof -i :8080  # macOS/Linux
netstat -ano | findstr :8080  # Windows

# Kill process
kill -9 <PID>  # macOS/Linux
taskkill /PID <PID> /F  # Windows
```

### GOWA Connection Failed

**Error:** Circuit breaker open / Connection refused

**Check:**

1. GOWA running: `curl http://localhost:3000/health`
2. Credentials correct in `config.yaml`
3. Circuit breaker state: `/api/health/detailed`

**Reset:**

- Wait 5 minutes for circuit breaker to auto-reset
- Or restart backend server

### JWT Secret Missing

**Error:** `Failed to load JWT secret`

**Solution:**

- Delete `data/jwt_secret.txt` and restart server
- Server will generate a new secret automatically
- All existing tokens will be invalidated (users must re-login)

### Data Corruption

**Symptoms:** Server crashes on startup, JSON parse errors

**Recovery:**

1. Backup corrupted file: `cp data/patients.json data/patients.json.backup`
2. Restore from backup or delete file
3. Restart server (will start with empty data)

---

## Performance Tips

### Monitor Goroutines

```go
import "runtime"

func logGoroutineCount() {
    count := runtime.NumGoroutine()
    logger.Info("Goroutine count", "count", count)
}
```

### Profile CPU Usage

```bash
# Enable profiling
go run main.go -cpuprofile=cpu.prof

# Analyze
go tool pprof cpu.prof
```

### Profile Memory Usage

```bash
# Enable profiling
go run main.go -memprofile=mem.prof

# Analyze
go tool pprof mem.prof
```

---

## Production Checklist

- [ ] Change default superadmin password
- [ ] Set strong JWT secret in `data/jwt_secret.txt`
- [ ] Use production GOWA endpoint (not localhost)
- [ ] Set `logging.level` to `info` or `warn`
- [ ] Enable HTTPS (use reverse proxy like Nginx)
- [ ] Set up log rotation
- [ ] Configure backup strategy for `data/` folder
- [ ] Monitor disk space (uploads can grow)
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure firewall rules (only expose reverse proxy)

---

**Next:** See [Development Guide - Frontend](./development-guide-frontend.md) for frontend setup.
