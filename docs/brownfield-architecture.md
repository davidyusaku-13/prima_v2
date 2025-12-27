# CareKeeper Brownfield Architecture Document

## Introduction

This document captures the current state of the CareKeeper codebase, a healthcare volunteer dashboard with role-based access control (RBAC) and WhatsApp notification capabilities. It serves as a reference for AI agents working on enhancements, bug fixes, or feature development.

### Document Scope

Comprehensive documentation of the CareKeeper healthcare volunteer management system.

### Change Log

| Date       | Version | Description                     | Author |
| ---------- | ------- | ------------------------------- | ------ |
| 2025-12-27 | 1.0     | Initial brownfield analysis     | Claude |
| 2025-12-27 | 2.0     | Updated to CareKeeper with RBAC | Claude |

---

## Quick Reference - Key Files and Entry Points

### Critical Files for Understanding the System

| Purpose             | File Path                   |
| ------------------- | --------------------------- |
| Backend Entry       | `backend/main.go`           |
| Frontend Entry      | `frontend/src/main.js`      |
| Main UI Component   | `frontend/src/App.svelte`   |
| API URL Constant    | `frontend/src/App.svelte:4` |
| Frontend Config     | `frontend/vite.config.js`   |
| CSS/Tailwind Config | `frontend/src/app.css`      |
| Project Guidance    | `CLAUDE.md` (outdated)      |

### Project Root Structure

```
prima_v2/
├── backend/          # Go/Gin backend (port 8080)
│   ├── main.go       # Complete application (~1150 lines)
│   ├── .env.example  # Environment template
│   ├── .env          # Actual env (gitignored)
│   ├── data/         # Persistence layer
│   │   ├── patients.json    # Patient data
│   │   ├── users.json       # User accounts
│   │   └── jwt_secret.txt   # JWT signing key
│   └── go.mod        # Go module definition
├── frontend/         # Svelte 5 + Vite frontend (port 5173)
├── docs/             # Documentation
├── CLAUDE.md         # Project guidance (OUTDATED)
└── .git/             # Git repository
```

---

## High Level Architecture

### Technical Summary

CareKeeper is a healthcare volunteer management dashboard with the following capabilities:

1. **Authentication**: JWT-based auth with registration/login
2. **RBAC**: Three roles (superadmin, admin, volunteer) with different permission levels
3. **Patient Management**: CRUD operations for patient records
4. **Reminder System**: Time-based reminders with recurrence support
5. **Notifications**: Automatic WhatsApp messages via GOWA when reminders are due

### Actual Tech Stack

| Category     | Technology        | Version | Notes                                |
| ------------ | ----------------- | ------- | ------------------------------------ |
| **Backend**  |                   |         |                                      |
| Runtime      | Go                | 1.25+   | Standard Go runtime                  |
| Framework    | Gin               | Latest  | HTTP web framework                   |
| CORS         | gin-contrib/cors  | Latest  | Middleware for cross-origin requests |
| JWT          | golang-jwt/jwt/v5 | Latest  | Token-based authentication           |
| **Frontend** |                   |         |                                      |
| Runtime      | Node.js           | 16+     | For Vite development server          |
| Framework    | Svelte            | 5.43.8  | Latest Svelte 5                      |
| Build Tool   | Vite              | 7.2.4   | Fast development server and bundler  |
| Styling      | Tailwind CSS      | 4.1.18  | CSS-first utility framework          |
| Integration  | @tailwindcss/vite | 4.1.18  | Vite plugin for Tailwind 4           |
| Package Mgr  | Bun               | -       | Used for frontend dev                |

### Repository Structure

- **Type**: Polyrepo (separate `backend/` and `frontend/` directories)
- **Package Manager**: Go modules (backend), bun/bun (frontend)
- **Version Control**: Git
- **Code Guidance**: `CLAUDE.md` exists but is **OUTDATED** - describes landing page not CareKeeper

---

## Source Tree and Module Organization

### Backend Structure (`backend/`)

```
backend/
├── main.go              # Complete application (~1150 lines, single file)
├── .env.example         # Environment variables template
├── .env                 # Actual environment (gitignored)
├── data/                # Persistence layer
│   ├── patients.json    # Patient data (auto-created)
│   ├── users.json       # User accounts (auto-created)
│   └── jwt_secret.txt   # JWT signing key (auto-created)
└── go.mod               # Go module definition
```

**Note**: The entire backend logic resides in a single `main.go` file (~1150 lines). This is a design choice for simplicity but may need refactoring for larger deployments.

### Frontend Structure (`frontend/`)

```
frontend/
├── src/
│   ├── main.js              # Svelte 5 mount entry point
│   ├── App.svelte           # Complete application (~1470 lines)
│   ├── app.css              # Tailwind CSS @import
│   └── assets/
│       └── svelte.svg       # Svelte logo (unused)
├── index.html               # HTML entry point
├── package.json             # Dependencies and scripts
├── vite.config.js           # Vite + Svelte + Tailwind config
└── node_modules/            # Dependencies
```

---

## Backend Implementation Details

### Entry Point: `backend/main.go`

All backend logic resides in a single file. Key sections:

**Imports and Constants (lines 1-28)**

```go
import (
    "bytes"
    "crypto/sha256"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "github.com/golang-jwt/jwt/v5"
)

const (
    dataFile       = "data/patients.json"
    usersDataFile  = "data/users.json"
    jwtSecretFile  = "data/jwt_secret.txt"
)
const tokenExpiry = 24 * 7 * time.Hour // 1 week
```

### Data Models

**User Model (with RBAC)**

```go
type Role string

const (
    RoleSuperadmin Role = "superadmin"
    RoleAdmin      Role = "admin"
    RoleVolunteer  Role = "volunteer"
)

type User struct {
    ID        string `json:"id"`
    Username  string `json:"username"`
    FullName  string `json:"fullName,omitempty"`
    Password  string `json:"password"` // SHA256 hashed
    Role      Role   `json:"role"`
    CreatedAt string `json:"createdAt"`
}
```

**Patient Model**

```go
type Patient struct {
    ID        string     `json:"id"`
    Name      string     `json:"name"`
    Phone     string     `json:"phone"` // Indonesian format (08xx or 628xx)
    Email     string     `json:"email,omitempty"`
    Notes     string     `json:"notes,omitempty"`
    Reminders []*Reminder `json:"reminders,omitempty"`
    CreatedBy string     `json:"createdBy,omitempty"` // Volunteer who created
}
```

**Reminder Model with Recurrence**

```go
type Reminder struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    DueDate     string     `json:"dueDate,omitempty"` // Format: "2006-01-02T15:04"
    Priority    string     `json:"priority"` // low, medium, high
    Completed   bool       `json:"completed"`
    Notified    bool       `json:"notified"` // Prevents duplicate WhatsApp
    Recurrence  Recurrence `json:"recurrence"`
}

type Recurrence struct {
    Frequency  string `json:"frequency"` // none, daily, weekly, monthly, yearly
    Interval   int    `json:"interval"` // Repeat every N periods
    DaysOfWeek []int  `json:"daysOfWeek"` // 0=Sun, 6=Sat for weekly
    EndDate    string `json:"endDate,omitempty"`
}
```

### In-Memory Storage with File Persistence

```go
type PatientStore struct {
    mu       sync.RWMutex
    patients map[string]*Patient
}

type UserStore struct {
    mu     sync.RWMutex
    users  map[string]*User
    byName map[string]string // username -> userID for fast lookup
}

var (
    store     = PatientStore{patients: make(map[string]*Patient)}
    userStore = UserStore{users: make(map[string]*User), byName: make(map[string]string)}
)
```

**Persistence Behavior**:

- Data loaded from `data/patients.json` and `data/users.json` on startup
- Data saved after every write operation (create/update/delete)
- File writes are asynchronous (`saveData()` spawns goroutine)
- JWT secret stored in `data/jwt_secret.txt`, auto-generated if missing

### API Endpoints

| Method              | Endpoint                                         | Auth             | Description                         |
| ------------------- | ------------------------------------------------ | ---------------- | ----------------------------------- |
| **Health**          |                                                  |                  |                                     |
| GET                 | `/api/health`                                    | None             | Health check `{"status": "ok"}`     |
| **Auth**            |                                                  |                  |                                     |
| POST                | `/api/auth/register`                             | None             | Register new user (role: volunteer) |
| POST                | `/api/auth/login`                                | None             | Login, returns JWT token            |
| GET                 | `/api/auth/me`                                   | JWT              | Get current user info               |
| **Patients**        |                                                  |                  |                                     |
| GET                 | `/api/patients`                                  | JWT              | List patients (filtered by role)    |
| POST                | `/api/patients`                                  | JWT              | Create patient                      |
| GET                 | `/api/patients/:id`                              | JWT              | Get patient by ID                   |
| PUT                 | `/api/patients/:id`                              | JWT              | Update patient                      |
| DELETE              | `/api/patients/:id`                              | JWT              | Delete patient                      |
| **Reminders**       |                                                  |                  |                                     |
| POST                | `/api/patients/:id/reminders`                    | JWT              | Create reminder                     |
| PUT                 | `/api/patients/:id/reminders/:reminderId`        | JWT              | Update reminder                     |
| POST                | `/api/patients/:id/reminders/:reminderId/toggle` | JWT              | Toggle completion                   |
| DELETE              | `/api/patients/:id/reminders/:reminderId`        | JWT              | Delete reminder                     |
| **User Management** |                                                  |                  |                                     |
| GET                 | `/api/users`                                     | JWT + Superadmin | List all users                      |
| PUT                 | `/api/users/:id/role`                            | JWT + Superadmin | Update user role                    |
| DELETE              | `/api/users/:id`                                 | JWT + Superadmin | Delete user                         |

### Role-Based Access Control (RBAC)

**Role Hierarchy**:

1. **superadmin** - Full access, can manage users
2. **admin** - Can view all patients, cannot manage users
3. **volunteer** - Can only see patients they created

**Access Control Implementation**:

```go
// In getPatients, getPatient, updatePatient, deletePatient
if role == string(RoleVolunteer) {
    if p.CreatedBy != userID {
        // Deny access
    }
}
```

**User Management Protection**:

```go
// Only superadmin can access user management routes
api.GET("/users", requireRole(RoleSuperadmin), getUsers)
api.PUT("/users/:id/role", requireRole(RoleSuperadmin), updateUserRole)
api.DELETE("/users/:id", requireRole(RoleSuperadmin), deleteUser)
```

### Default Superadmin Account

```go
// Created automatically on first run
Username: superadmin
Password: superadmin
Role: superadmin
ID: superadmin
```

**Constraint**: The default superadmin password is hardcoded to "superadmin" - should be changed in production.

### CORS Configuration

CORS is explicitly configured for frontend origin only:

```go
router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    AllowCredentials: true,
}))
```

**Constraint**: Frontend URL is hardcoded, not configurable via environment.

### JWT Authentication

**Token Structure**:

```go
type Claims struct {
    UserID   string `json:"userId"`
    Username string `json:"username"`
    Role     Role   `json:"role"`
    jwt.RegisteredClaims
}
```

**Token Features**:

- Algorithm: HS256
- Expiry: 7 days (configurable via `tokenExpiry`)
- Secret: Stored in file, auto-generated on first run

**Auth Middleware**:

```go
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        // Expects "Bearer <token>"
        // Validates token and sets user context
    }
}
```

### WhatsApp Integration (GOWA)

**Configuration**:

```go
func getGowaEndpoint() string { return getEnv("GOWA_ENDPOINT", "http://localhost:3000") }
func getGowaUser() string     { return getEnv("GOWA_USER", "admin") }
func getGowaPass() string     { return getEnv("GOWA_PASS", "password123") }
```

**Phone Number Formatting**:

```go
func formatWhatsAppNumber(phone string) string {
    // 08xx -> 628xx (Indonesian format)
    // Returns: 6281234567890@s.whatsapp.net
}
```

**Sending Messages**:

```go
func sendWhatsAppMessage(phone, message string) error {
    endpoint := getGowaEndpoint() + "/send/message"
    // HTTP POST with Basic auth
}
```

### Automatic Reminder Checker

Background goroutine runs every minute:

```go
func checkReminders() {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        // Check all reminders
        // If due (within 5 minute window) and not notified
        // Send WhatsApp message
        // Mark as notified
    }
}
```

**Reminder Due Window**: Reminders are checked if `now >= dueTime` AND `now < dueTime + 5 minutes`

---

## Frontend Implementation Details

### Entry Point: `frontend/src/main.js`

Uses Svelte 5's new mounting API:

```javascript
import { mount } from "svelte";
import "./app.css";
import App from "./App.svelte";

const app = mount(App, {
  target: document.getElementById("app"),
});

export default app;
```

### Main Component: `frontend/src/App.svelte`

**Single component handles entire application UI (~1470 lines).**

#### State Management

```javascript
// Auth state
let token = localStorage.getItem("token") || null;
let user = null;
let authLoading = true;

// Application state
let patients = [];
let users = [];
let loading = true;
let sidebarOpen = true;
let currentView = "dashboard"; // 'dashboard', 'patients', 'users'
```

#### Authentication Flow

1. **onMount**: Check for existing token, fetch user profile
2. **Login**: POST to `/api/auth/login`, store token in localStorage
3. **Register**: POST to `/api/auth/register`, auto-login with returned token
4. **Logout**: Clear token and state

#### Role-Based UI Rendering

```svelte
{#if user?.role === 'superadmin'}
  <button onclick={() => { currentView = 'users'; loadUsers(); }}>
    Users
  </button>
{/if}
```

#### Views

| View      | Access          | Description                                |
| --------- | --------------- | ------------------------------------------ |
| dashboard | All roles       | Stats, upcoming reminders, recent patients |
| patients  | All roles       | Patient list with search, CRUD operations  |
| users     | Superadmin only | User management table                      |

#### Modals

1. **Auth Modal**: Login/Register forms with validation
2. **Patient Modal**: Add/Edit patient (name, phone, email, notes)
3. **Reminder Modal**: Add/Edit reminder (title, description, due date, priority, recurrence)
4. **User Modal**: Edit user role (superadmin only)

#### Key Functions

| Function                   | Purpose                    |
| -------------------------- | -------------------------- |
| `login()`                  | Authenticate user          |
| `register()`               | Create new account         |
| `logout()`                 | Clear auth state           |
| `loadPatients()`           | Fetch patients from API    |
| `savePatient()`            | Create/update patient      |
| `deletePatient()`          | Remove patient             |
| `saveReminder()`           | Create/update reminder     |
| `toggleReminderComplete()` | Toggle reminder completion |
| `deleteReminder()`         | Remove reminder            |
| `loadUsers()`              | Fetch users (superadmin)   |
| `updateUserRole()`         | Change user role           |
| `deleteUser()`             | Remove user                |

#### Computed Values (Svelte Reactivity)

```javascript
$: filteredPatients = patients.filter(p => {
    // Search filter by name, phone, email
});

$: upcomingReminders = patients.flatMap(p => ...)
    .filter(r => !r.completed && r.dueDate)
    .sort((a, b) => new Date(a.dueDate) - new Date(b.dueDate))
    .slice(0, 5);

$: stats = {
    totalPatients: patients.length,
    totalReminders: patients.reduce(...),
    completedReminders: patients.reduce(...),
    pendingReminders: patients.reduce(...)
};
```

### Tailwind CSS Integration

- **Version**: Tailwind CSS 4 (CSS-first configuration)
- **Config**: `@tailwindcss/vite` plugin in `vite.config.js`
- **Import**: `@import "tailwindcss"` in `app.css`
- **Usage**: Utility classes directly in Svelte template

---

## Technical Debt and Known Issues

### Critical Technical Debt

1. **Single-File Backend**

   - Entire backend in `main.go` (~1150 lines)
   - No separation of concerns (handlers, services, models)
   - Difficult to test individual components

2. **Single-File Frontend Component**

   - All UI logic in `App.svelte` (~1470 lines)
   - No reusable UI components extracted
   - Modal logic intermixed with view logic

3. **Hardcoded Values**

   - Default superadmin password: "superadmin"
   - Frontend API URL hardcoded in `App.svelte:4`
   - CORS origin hardcoded in `main.go`
   - GOWA default credentials hardcoded

4. **No Tests**

   - No unit tests
   - No integration tests
   - No E2E tests

5. **Outdated Documentation**
   - `CLAUDE.md` describes wrong project (landing page)
   - `docs/brownfield-architecture.md` partially outdated
   - `docs/project-documentation.md` describes older version

### Workarounds and Gotchas

- **CORS Strict**: Backend only accepts `http://localhost:5173`
- **Password Hashing**: Uses SHA256 (not recommended for production, use bcrypt)
- **JWT Secret**: Auto-generated and stored in file (ok for single instance)
- **File Persistence**: Async writes could lead to data loss on crash
- **Indonesian Phone Format**: Phone numbers must be Indonesian format for WhatsApp
- **Volunteer Isolation**: Volunteers can only see their own patients (security feature)
- **No Pagination**: Large patient lists will impact performance
- **No Soft Delete**: Deletion is permanent

---

## Integration Points and External Dependencies

### External Services

| Service   | Purpose            | Integration Type | Key Files               |
| --------- | ------------------ | ---------------- | ----------------------- |
| GOWA      | WhatsApp messaging | REST API         | `main.go:570-597`       |
| Go stdlib | Core runtime       | Built-in         | `main.go`               |
| Node.js   | Frontend build     | Runtime          | `frontend/package.json` |

### Internal Integration Points

- **Frontend → Backend**: REST API on port 8080 with JWT auth
- **Svelte → Tailwind**: CSS utility classes in templates
- **Vite → Svelte**: Build plugin in `vite.config.js`
- **GOWA → WhatsApp**: External service for messaging

### Environment Variables

| Variable      | Default                 | Purpose                  |
| ------------- | ----------------------- | ------------------------ |
| GOWA_ENDPOINT | `http://localhost:3000` | WhatsApp API URL         |
| GOWA_USER     | `admin`                 | GOWA basic auth username |
| GOWA_PASS     | `password123`           | GOWA basic auth password |

**Note**: Environment variables loaded from `.env` file, not system environment.

### Ports and Hosts

| Service  | Host:Port      | Configuration Location |
| -------- | -------------- | ---------------------- |
| Backend  | localhost:8080 | `main.go:307`          |
| Frontend | localhost:5173 | Vite default           |
| GOWA     | localhost:3000 | `main.go:143`          |

---

## Development and Deployment

### Local Development Setup

#### Prerequisites

- Go 1.25+
- Node.js 16+
- Bun (recommended) or bun
- GOWA server running (for WhatsApp notifications)

#### Starting the Backend

```bash
cd backend
# Create .env if needed
go run main.go
```

Backend runs on `http://localhost:8080`

#### Starting the Frontend

```bash
cd frontend
bun run dev
# or: bun run dev
```

Frontend runs on `http://localhost:5173`

#### Starting GOWA (Optional)

```bash
# Ensure GOWA is running on port 3000 for WhatsApp features
```

#### Accessing the App

1. Open `http://localhost:5173` in browser
2. Login with default superadmin: `superadmin` / `superadmin`
3. Or register a new account (becomes volunteer)

### Build and Deployment Process

- **Frontend Build**: `bun run build` (outputs to `frontend/dist/`)
- **Backend Build**: `go build -o carekeeper main.go` (standalone binary)
- **No CI/CD**: Manual deployment required
- **No Docker**: No containerization configured

### Environment Configuration

1. Copy `backend/.env.example` to `backend/.env`
2. Configure GOWA credentials
3. Consider changing default superadmin password

---

## Testing Reality

### Current Test Coverage

- **Unit Tests**: None
- **Integration Tests**: None
- **E2E Tests**: None
- **Manual Testing**: Primary QA method

### Running Tests

```bash
# No test scripts configured
bun test   # Not available
go test    # No test files
```

**Gap**: The project has no automated tests. This is a significant risk for a healthcare application.

---

## Enhancement Considerations

### Files That May Need Modification

| Enhancement Type          | Files to Modify                            |
| ------------------------- | ------------------------------------------ |
| Add email notifications   | `main.go` (new handler), `App.svelte`      |
| Add SMS support           | `main.go` (new provider), `Reminder` model |
| Add patient categories    | `main.go` (model + handlers), `App.svelte` |
| Add patient notes history | `main.go` (model + handlers), `App.svelte` |
| Add reporting/analytics   | `main.go` (new endpoints), `App.svelte`    |
| Add data export           | `main.go` (export handlers)                |
| Add two-factor auth       | `main.go` (auth flow), `App.svelte`        |
| Improve password hashing  | `main.go:80-88` (bcrypt instead of SHA256) |

### New Files/Modules That May Be Needed

- `backend/api/auth.go` - Authentication handlers (extract from main.go)
- `backend/api/patients.go` - Patient handlers
- `backend/api/reminders.go` - Reminder handlers
- `backend/services/whatsapp.go` - WhatsApp integration
- `frontend/src/components/Navigation.svelte` - Sidebar component
- `frontend/src/components/PatientCard.svelte` - Reusable card
- `frontend/src/components/ReminderItem.svelte` - Reminder list item
- `backend/middleware/auth.go` - Auth middleware (extract)
- `backend/models/models.go` - Shared models

### Integration Considerations

- **Email service**: Add SendGrid, Mailgun, or similar for email notifications
- **SMS provider**: Add Twilio or local provider for SMS
- **Database**: Consider PostgreSQL for production (current: JSON file)
- **Caching**: Add Redis for session/jwt caching in multi-instance deployment
- **Rate limiting**: Add Gin middleware for abuse prevention

---

## Appendix - Useful Commands and Scripts

### Backend Commands

```bash
# Run backend
cd backend && go run main.go

# Build backend binary
cd backend && go build -o carekeeper main.go

# Check Go syntax
cd backend && go vet

# List dependencies
cd backend && go list -m
```

### Frontend Commands

```bash
# Install dependencies
cd frontend && bun install

# Start dev server
cd frontend && bun run dev

# Build for production
cd frontend && bun run build

# Preview production build
cd frontend && bun run preview
```

### Debugging and Troubleshooting

- **Backend logs**: Console output from `go run main.go`
- **Frontend logs**: Browser DevTools console
- **Network**: Browser DevTools Network tab
- **CORS errors**: Check browser console, verify backend CORS config
- **Auth issues**: Check token expiration, verify JWT secret
- **WhatsApp issues**: Verify GOWA server running, check credentials

### Common Issues

| Issue                  | Solution                                      |
| ---------------------- | --------------------------------------------- |
| CORS error in browser  | Ensure backend is running before frontend     |
| Connection refused     | Check backend is on port 8080                 |
| Build fails            | Run `bun install` in frontend directory       |
| Login fails            | Check credentials, verify backend is running  |
| WhatsApp not sending   | Verify GOWA running, check phone format       |
| Can't see patients     | Check role - volunteers only see own patients |
| Can't access Users tab | Must be superadmin role                       |

---

## References

- **CLAUDE.md**: Project guidance (OUTDATED - describes wrong project)
- **Gin Framework**: https://gin-gonic.com/
- **Svelte 5**: https://svelte.dev/blog/runes
- **Tailwind CSS 4**: https://tailwindcss.com/
- **Vite**: https://vitejs.dev/
- **JWT**: https://github.com/golang-jwt/jwt/
- **GOWA**: Go WhatsApp Web Multi-Device library
