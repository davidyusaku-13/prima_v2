# PRIMA Brownfield Architecture Document

## Introduction

This document captures the current state of the PRIMA codebase, a healthcare volunteer dashboard with role-based access control (RBAC), WhatsApp notification capabilities, and a public-facing health education content management system. It serves as a reference for AI agents working on enhancements, bug fixes, or feature development.

### Document Scope

Comprehensive documentation of the PRIMA healthcare volunteer management and content system.

### Change Log

| Date       | Version | Description                              | Author  |
| ---------- | ------- | ---------------------------------------- | ------- |
| 2025-12-27 | 1.0     | Initial brownfield analysis              | Claude  |
| 2025-12-27 | 2.0     | Updated to CareKeeper with RBAC          | Claude  |
| 2025-12-27 | 3.0     | Added CMS features (Berita, Video Edukasi)| Claude  |
| 2025-12-28 | 3.1     | Documentation review and CLAUDE.md sync  | Claude  |
| 2025-12-28 | 3.2     | CMS dashboard redesign, draft visibility, Svelte 5 fixes | Claude  |

---

## Quick Reference - Key Files and Entry Points

### Critical Files for Understanding the System

| Purpose             | File Path                              |
| ------------------- | -------------------------------------- |
| Backend Entry       | `backend/main.go`                      |
| Backend Models      | `backend/models/content.go`            |
| Backend Handlers    | `backend/handlers/content.go`          |
| Backend Utils       | `backend/utils/youtube.go`             |
| Frontend Entry      | `frontend/src/main.js`                 |
| Main UI Component   | `frontend/src/App.svelte`              |
| API URL Constant    | `frontend/src/lib/utils/api.js:1`      |
| Frontend Config     | `frontend/vite.config.js`              |
| CSS/Tailwind Config | `frontend/src/app.css`                 |
| Project Guidance    | `CLAUDE.md`                            |

### Project Root Structure

```
prima_v2/
├── backend/                  # Go/Gin backend (port 8080)
│   ├── main.go               # Main application with auth/patient/reminder handlers
│   ├── models/
│   │   └── content.go        # Content models (Category, Article, Video)
│   ├── handlers/
│   │   └── content.go        # Content CRUD handlers, image upload
│   ├── utils/
│   │   └── youtube.go        # YouTube URL parsing and metadata fetch
│   ├── .env.example          # Environment template
│   ├── .env                  # Actual env (gitignored)
│   ├── data/                 # Persistence layer
│   │   ├── patients.json     # Patient data
│   │   ├── users.json        # User accounts
│   │   ├── categories.json   # Content categories
│   │   ├── articles.json     # Published articles
│   │   ├── videos.json       # YouTube videos
│   │   └── jwt_secret.txt    # JWT signing key
│   ├── uploads/              # Uploaded images (auto-created)
│   └── go.mod                # Go module definition
├── frontend/                 # Svelte 5 + Vite frontend (port 5173)
│   ├── src/
│   │   ├── main.js           # Svelte 5 mount entry point
│   │   ├── App.svelte        # Complete application (~610 lines core, + modals)
│   │   ├── app.css           # Tailwind CSS @import
│   │   └── lib/
│   │       ├── components/   # Reusable UI components
│   │       ├── views/        # Page views
│   │       ├── utils/
│   │       │   └── api.js    # API functions
│   │       └── i18n/         # Internationalization
│   ├── index.html            # HTML entry point
│   ├── package.json          # Dependencies and scripts
│   └── vite.config.js        # Vite + Svelte + Tailwind config
├── docs/                     # Documentation
│   ├── brownfield-architecture.md
│   ├── project-documentation.md
│   └── brainstorming-session-results.md
└── GOWA-README.md            # WhatsApp integration docs
```

---

## High Level Architecture

### Technical Summary

PRIMA is a healthcare volunteer management and health education platform with the following capabilities:

1. **Authentication**: JWT-based auth with registration/login
2. **RBAC**: Three roles (superadmin, admin, volunteer) with different permission levels
3. **Patient Management**: CRUD operations for patient records
4. **Reminder System**: Time-based reminders with recurrence support
5. **WhatsApp Notifications**: Automatic WhatsApp messages via GOWA when reminders are due
6. **Public Content (NEW)**: Berita (health news) and Video Edukasi (educational videos)
7. **CMS (NEW)**: Admin interface for managing content

### Actual Tech Stack

| Category     | Technology          | Version | Notes                                    |
| ------------ | ------------------- | ------- | ---------------------------------------- |
| **Backend**  |                     |         |                                          |
| Runtime      | Go                  | 1.25+   | Standard Go runtime                      |
| Framework    | Gin                 | Latest  | HTTP web framework                       |
| CORS         | gin-contrib/cors    | Latest  | Middleware for cross-origin requests     |
| JWT          | golang-jwt/jwt/v5   | Latest  | Token-based authentication               |
| Image Proc   | disintegration/imaging | Latest | Image resize for hero images             |
| **Frontend** |                     |         |                                          |
| Runtime      | Node.js             | 16+     | For Vite development server              |
| Framework    | Svelte              | 5.43.8  | Latest Svelte 5                          |
| Build Tool   | Vite                | 7.2.4   | Fast development server and bundler      |
| Styling      | Tailwind CSS        | 4.1.18  | CSS-first utility framework              |
| Integration  | @tailwindcss/vite   | 4.1.18  | Vite plugin for Tailwind 4               |
| i18n         | svelte-i18n         | Latest  | Internationalization support             |
| Package Mgr  | Bun                 | -       | Used for frontend dev                    |

### Repository Structure

- **Type**: Polyrepo (separate `backend/` and `frontend/` directories)
- **Package Manager**: Go modules (backend), bun (frontend)
- **Version Control**: Git

---

## Source Tree and Module Organization

### Backend Structure

```
backend/
├── main.go                   # Main application (~1200 lines)
│                            # - Auth handlers (register, login, me)
│                            # - Patient CRUD handlers
│                            # - Reminder CRUD handlers
│                            # - User management (superadmin)
│                            # - Role-based middleware
│                            # - WhatsApp/GOWA integration
│                            # - Reminder checker goroutine
│                            # - Content routes setup
├── models/
│   └── content.go           # Content models (Category, Article, Video)
│                            # - CategoryStore, ArticleStore, VideoStore
│                            # - Slug generation, time formatting
├── handlers/
│   └── content.go           # Content handlers
│                            # - Category CRUD
│                            # - Article CRUD (with version tracking)
│                            # - Video CRUD (YouTube URL parsing)
│                            # - Image upload with resize (16:9, 1:1, 4:3)
│                            # - Dashboard stats
├── utils/
│   └── youtube.go           # YouTube utilities
│                            # - ExtractYouTubeID()
│                            # - ValidateYouTubeURL()
│                            # - FetchYouTubeMetadata() via noembed API
│                            # - GetYouTubeThumbnailURL()
├── data/                    # JSON file persistence
│   ├── patients.json
│   ├── users.json
│   ├── categories.json
│   ├── articles.json
│   ├── videos.json
│   └── jwt_secret.txt
├── uploads/                 # Auto-created, image storage
│   ├── <id>_16x9.jpg
│   ├── <id>_1x1.jpg
│   └── <id>_4x3.jpg
├── go.mod
└── go.sum
```

### Frontend Structure

```
frontend/src/
├── main.js                  # Svelte 5 mount entry point
├── App.svelte               # Main application orchestrator
├── app.css                  # Tailwind CSS @import
└── lib/
    ├── components/          # Reusable UI components
    │   ├── Sidebar.svelte       # Navigation sidebar
    │   ├── BottomNav.svelte     # Mobile bottom navigation
    │   ├── PatientModal.svelte  # Add/edit patient
    │   ├── ReminderModal.svelte # Add/edit reminder
    │   ├── UserModal.svelte     # User management
    │   ├── ConfirmModal.svelte  # Delete confirmation
    │   ├── ProfileModal.svelte  # User profile/logout
    │   ├── VideoModal.svelte    # YouTube video embed modal
    │   ├── ImageUploader.svelte # Drag & drop image upload
    │   ├── ArticleCard.svelte   # Article preview card
    │   ├── VideoCard.svelte     # Video preview card
    │   ├── ContentListItem.svelte # CMS list view row item
    │   ├── DashboardStats.svelte # CMS stats display
    │   └── ActivityLog.svelte   # Recent activity feed
    ├── views/               # Page views
    │   ├── LoginScreen.svelte   # Login/Register screen
    │   ├── DashboardView.svelte  # Dashboard stats
    │   ├── PatientsView.svelte   # Patient list
    │   ├── UsersView.svelte      # User management (superadmin)
    │   ├── BeritaView.svelte     # Health news list
    │   ├── BeritaDetailView.svelte # Article reader
    │   ├── VideoEdukasiView.svelte # YouTube video gallery
    │   ├── CMSDashboardView.svelte # CMS admin dashboard
    │   ├── ArticleEditorView.svelte # Create/edit articles
    │   └── VideoManagerView.svelte  # Add YouTube videos
    ├── utils/
    │   └── api.js           # API functions (all endpoints)
    └── i18n/                # Internationalization
        ├── en.json          # English translations
        └── id.json          # Indonesian translations
```

---

## Backend Implementation Details

### Entry Point: `backend/main.go`

All backend logic resides in a single file with content handlers in separate files.

**Imports and Constants (lines 1-32)**

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
    "github.com/davidyusaku-13/prima_v2/handlers"
)

const (
    dataFile           = "data/patients.json"
    usersDataFile      = "data/users.json"
    jwtSecretFile      = "data/jwt_secret.txt"
    categoriesDataFile = "data/categories.json"
    articlesDataFile   = "data/articles.json"
    videosDataFile     = "data/videos.json"
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

**Content Models (NEW - `backend/models/content.go`)**

```go
type Category struct {
    ID        string       `json:"id"`
    Name      string       `json:"name"`
    Type      CategoryType `json:"type"` // "article" or "video"
    CreatedAt string       `json:"created_at"`
}

type Article struct {
    ID          string       `json:"id"`
    Title       string       `json:"title"`
    Slug        string       `json:"slug"` // URL-friendly identifier
    Excerpt     string       `json:"excerpt"`
    Content     string       `json:"content"`
    AuthorID    string       `json:"author_id"`
    CategoryID  string       `json:"category_id"`
    HeroImages  HeroImages   `json:"hero_images"` // 16:9, 1:1, 4:3 versions
    Status      ArticleStatus `json:"status"` // "draft" or "published"
    Version     int          `json:"version"` // For version history
    ViewCount   int          `json:"view_count"`
    CreatedAt   string       `json:"created_at"`
    PublishedAt string       `json:"published_at"`
    UpdatedAt   string       `json:"updated_at"`
}

type Video struct {
    ID           string     `json:"id"`
    YouTubeURL   string     `json:"youtube_url"`
    YouTubeID    string     `json:"youtube_id"`
    Title        string     `json:"title"`
    Description  string     `json:"description"`
    ChannelName  string     `json:"channel_name"`
    ThumbnailURL string     `json:"thumbnail_url"`
    Duration     string     `json:"duration"`
    CategoryID   string     `json:"category_id"`
    Status       string     `json:"status"` // "published"
    ViewCount    int        `json:"view_count"`
    CreatedAt    string     `json:"created_at"`
    UpdatedAt    string     `json:"updated_at"`
}

type HeroImages struct {
    Hero16x9 string `json:"hero_16x9"` // 1920x1080
    Hero1x1  string `json:"hero_1x1"`  // 1080x1080
    Hero4x3  string `json:"hero_4x3"`  // 1600x1200
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

// Content stores with mutex protection
type ContentStore struct {
    Categories *CategoryStore
    Articles   *ArticleStore
    Videos     *VideoStore
}

var (
    store        = PatientStore{patients: make(map[string]*Patient)}
    userStore    = UserStore{users: make(map[string]*User), byName: make(map[string]string)}
    contentStore *handlers.ContentStore
)
```

**Persistence Behavior**:

- Data loaded from JSON files on startup
- Data saved after every write operation (async goroutine)
- File writes use tmp file + rename for atomicity
- JWT secret stored in `data/jwt_secret.txt`, auto-generated if missing

---

## API Endpoints

### Core Endpoints

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| **Health** |||
| GET | `/api/health` | None | Health check |
| **Auth** |||
| POST | `/api/auth/register` | None | Register user (role: volunteer) |
| POST | `/api/auth/login` | None | Login, returns JWT token |
| GET | `/api/auth/me` | JWT | Get current user |
| **Patients** |||
| GET | `/api/patients` | JWT | List patients (role-filtered) |
| POST | `/api/patients` | JWT | Create patient |
| GET | `/api/patients/:id` | JWT | Get patient by ID |
| PUT | `/api/patients/:id` | JWT | Update patient |
| DELETE | `/api/patients/:id` | JWT | Delete patient |
| **Reminders** |||
| POST | `/api/patients/:id/reminders` | JWT | Create reminder |
| PUT | `/api/patients/:id/reminders/:reminderId` | JWT | Update reminder |
| POST | `/api/patients/:id/reminders/:reminderId/toggle` | JWT | Toggle completion |
| DELETE | `/api/patients/:id/reminders/:reminderId` | JWT | Delete reminder |
| **User Management** |||
| GET | `/api/users` | JWT+Superadmin | List all users |
| PUT | `/api/users/:id/role` | JWT+Superadmin | Update user role |
| DELETE | `/api/users/:id` | JWT+Superadmin | Delete user |

### Content Endpoints (NEW)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| **Categories** |||
| GET | `/api/categories` | Public | List all categories |
| GET | `/api/categories/:type` | Public | Get categories by type (article/video) |
| POST | `/api/categories` | Admin+ | Create category |
| **Articles** |||
| GET | `/api/articles` | Public | List published articles (use `?all=true` for drafts) |
| GET | `/api/articles/:slug` | Public | Get article by slug |
| POST | `/api/articles` | Admin+ | Create article |
| PUT | `/api/articles/:id` | Admin+ | Update article |
| DELETE | `/api/articles/:id` | Admin+ | Delete article |
| **Videos** |||
| GET | `/api/videos` | Public | List published videos |
| POST | `/api/videos` | Admin+ | Add YouTube video (auto-fetches metadata) |
| DELETE | `/api/videos/:id` | Admin+ | Delete video |
| **Upload** |||
| POST | `/api/upload/image` | Admin+ | Upload & resize image (16:9, 1:1, 4:3) |
| **Dashboard** |||
| GET | `/api/dashboard/stats` | Admin+ | Get CMS statistics |

### Role-Based Access Control (RBAC)

**Role Hierarchy**:
1. **superadmin** - Full access, can manage users, all content
2. **admin** - Can view all patients, manage content, cannot manage users
3. **volunteer** - Can only see patients they created, read-only content access

**Content Access**:
- **Public**: GET articles, GET videos, GET categories
- **Admin/Superadmin**: Full content CRUD
- **Volunteer**: Read-only content (no CMS access)

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

Single component orchestrates entire application.

**State Management**:

```javascript
// Auth state
let token = localStorage.getItem('token') || null;
let user = null;
let authLoading = true;

// Application state
let patients = [];
let users = [];
let loading = true;
let currentView = localStorage.getItem('currentView') || 'dashboard';

// CMS State
let showArticleEditor = false;
let editingArticle = null;
let showVideoManager = false;
let showVideoModal = false;
let currentVideo = null;
let currentArticleId = null;
```

**Views (currentView values)**:
- `dashboard` - Stats, upcoming reminders, recent patients
- `patients` - Patient list with search, CRUD
- `users` - User management (superadmin only)
- `berita` - Health news list (public)
- `berita-detail` - Article reader (public)
- `video` - Video gallery (public)
- `cms` - CMS dashboard (admin+)

### API Functions (`frontend/src/lib/utils/api.js`)

All API calls centralized here:

```javascript
// Auth
login(username, password), register(username, password, fullName), fetchUser(token)

// Patients
fetchPatients(token), savePatient(token, patient, editingId), deletePatient(token, id)

// Reminders
saveReminder(token, patientId, reminder, editingId), toggleReminder(token, patientId, reminderId)

// Users (superadmin)
fetchUsers(token), updateUserRole(token, userId, role), deleteUser(token, userId)

// CMS - Articles
fetchArticles(token, category, all), fetchArticle(token, id), createArticle(token, article),
updateArticle(token, id, article), deleteArticle(token, id)

// CMS - Videos
fetchVideos(token, category), createVideo(token, video), updateVideo(token, id, video), deleteVideo(token, id)

// CMS - Categories
fetchCategories(token)

// CMS - Dashboard
fetchDashboardStats(token), uploadImage(token, file)
```

### Navigation Components

**Sidebar** (`src/lib/components/Sidebar.svelte`):
- Dashboard link (all roles)
- Patients link (all roles)
- Users link (superadmin only)
- CMS link (admin+)
- Berita link (all roles)
- Video Edukasi link (all roles)

**BottomNav** (`src/lib/components/BottomNav.svelte`):
- Mobile navigation with same links as sidebar

### Content Views (NEW)

**BeritaView.svelte** - Health news list:
- Category filtering
- Search functionality
- Article cards with hero images
- Navigate to detail on click

**BeritaDetailView.svelte** - Article reader:
- Full article content
- Hero image display
- Back navigation

**VideoEdukasiView.svelte** - Video gallery:
- YouTube embed grid
- Category filtering
- Click to play in modal

**CMSDashboardView.svelte** - Admin dashboard:
- Stats overview (articles, videos, views, drafts)
- Content list with list/grid view toggle
- Search, filter (all/articles/videos/drafts), and sort functionality
- Bulk selection with select-all checkbox
- Bulk delete action for selected items
- Pagination for large content lists
- Quick action buttons (add article, add video)

**ArticleEditorView.svelte** - Article management:
- Title, excerpt, content, category
- Hero image uploader (16:9, 1:1, 4:3)
- Draft/publish toggle
- Version tracking ready

**VideoManagerView.svelte** - Video management:
- YouTube URL input
- Auto-fetch metadata preview
- Category selection
- Delete videos

---

## Technical Debt and Known Issues

### Critical Technical Debt

1. **Single-File Backend Core**
   - `main.go` contains ~1200 lines of mixed concerns
   - Handlers separated to `handlers/content.go` (good)
   - No separation of concerns for patient/reminder handlers
   - Difficult to test individual components

2. **Password Hashing**
   - Uses SHA256 (not recommended for production)
   - Should use bcrypt or argon2

3. **Hardcoded Values**
   - Default superadmin password: "superadmin"
   - Frontend API URL hardcoded in `api.js:1`
   - CORS origin hardcoded in `main.go`
   - GOWA default credentials hardcoded

4. **No Tests**
   - No unit tests
   - No integration tests
   - No E2E tests

5. **Documentation**
   - `CLAUDE.md` is now up-to-date with CMS features
   - This document is current

### Workarounds and Gotchas

- **CORS Strict**: Backend only accepts `http://localhost:5173`
- **Indonesian Phone Format**: Phone numbers must be Indonesian format for WhatsApp
- **Volunteer Isolation**: Volunteers can only see their own patients
- **No Pagination**: Large lists will impact performance (except CMS dashboard)
- **No Soft Delete**: Deletion is permanent
- **File Persistence**: Async writes could lead to data loss on crash
- **Activity Log Endpoint**: `fetchActivityLog()` returns empty - backend doesn't have this endpoint yet
- **Slug Uniqueness**: Articles get numbered suffixes if slug collides
- **Svelte 5 Set Reactivity**: Must create new `Set()` instances instead of mutating in place (`.add()`, `.delete()`, `.clear()` don't trigger reactivity)
- **Draft Articles**: Use `?all=true` query param on `/api/articles` to include drafts (CMS dashboard uses this)

---

## Integration Points and External Dependencies

### External Services

| Service | Purpose | Integration Type | Key Files |
|---------|---------|-----------------|-----------|
| GOWA | WhatsApp messaging | REST API | `main.go:617-644` |
| YouTube | Video hosting | oEmbed/noembed API | `backend/utils/youtube.go` |
| Node.js | Frontend build | Runtime | `frontend/package.json` |

### Internal Integration Points

- **Frontend → Backend**: REST API on port 8080 with JWT auth
- **Svelte → Tailwind**: CSS utility classes in templates
- **Vite → Svelte**: Build plugin in `vite.config.js`
- **GOWA → WhatsApp**: External service for messaging

### Environment Variables

| Variable | Default | Purpose |
|----------|---------|---------|
| GOWA_ENDPOINT | `http://localhost:3000` | WhatsApp API URL |
| GOWA_USER | `admin` | GOWA basic auth username |
| GOWA_PASS | `password123` | GOWA basic auth password |

### Ports and Hosts

| Service | Host:Port | Configuration Location |
|---------|-----------|------------------------|
| Backend | localhost:8080 | `main.go:354` |
| Frontend | localhost:5173 | Vite default |
| GOWA | localhost:3000 | `main.go:148` |

---

## Development and Deployment

### Local Development Setup

#### Prerequisites

- Go 1.25+
- Node.js 16+
- Bun
- GOWA server (optional, for WhatsApp features)

#### Starting the Backend

```bash
cd backend
go run main.go
```

Backend runs on `http://localhost:8080`

#### Starting the Frontend

```bash
cd frontend
bun run dev
```

Frontend runs on `http://localhost:5173`

#### Accessing the App

1. Open `http://localhost:5173` in browser
2. Login with default superadmin: `superadmin` / `superadmin`
3. Or register a new account (becomes volunteer)

### Build and Deployment

- **Frontend Build**: `bun run build` (outputs to `frontend/dist/`)
- **Backend Build**: `go build -o prima main.go`
- **No CI/CD**: Manual deployment required
- **No Docker**: No containerization configured

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

---

## Enhancement Considerations

### Files That May Need Modification

| Enhancement Type | Files to Modify |
|-----------------|-----------------|
| Add email notifications | `main.go` (new handler), `App.svelte` |
| Add SMS support | `main.go` (new provider), `Reminder` model |
| Add patient categories | `main.go` (handlers), `App.svelte` |
| Add version history UI | `ArticleEditorView.svelte`, `handlers/content.go` |
| Add activity logging | `handlers/content.go`, new endpoint |
| Add search | `handlers/content.go`, frontend search UI |
| Improve password hashing | `main.go:85-88` (bcrypt instead of SHA256) |

### New Files/Modules That May Be Needed

- `backend/services/notification.go` - Email/SMS notifications
- `backend/handlers/activity.go` - Activity logging
- `frontend/src/components/SearchBar.svelte` - Search component
- `frontend/src/components/Pagination.svelte` - List pagination

### Integration Considerations

- **Email service**: Add SendGrid, Mailgun, or similar
- **SMS provider**: Add Twilio or local provider
- **Database**: Consider PostgreSQL for production
- **Caching**: Add Redis for sessions in multi-instance deployment
- **Rate limiting**: Add Gin middleware for abuse prevention

---

## Appendix - Useful Commands and Scripts

### Backend Commands

```bash
# Run backend
cd backend && go run main.go

# Build backend binary
cd backend && go build -o prima main.go

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

| Issue | Solution |
|-------|----------|
| CORS error in browser | Ensure backend is running before frontend |
| Connection refused | Check backend is on port 8080 |
| Build fails | Run `bun install` in frontend directory |
| Login fails | Check credentials, verify backend is running |
| WhatsApp not sending | Verify GOWA running, check phone format |
| Can't see patients | Check role - volunteers only see own patients |
| Can't access Users tab | Must be superadmin role |
| Can't access CMS | Must be admin or superadmin role |
| Articles not showing | Check status is "published" not "draft" |
| Video fetch fails | Check YouTube URL format, noembed API |
| Drafts not showing in CMS | Ensure `?all=true` param is passed to articles API |
| Bulk selection not working | Svelte 5 requires new Set() instances for reactivity |
| CMS columns misaligned | Check header and row widths match in CMSDashboardView |

---

## References

- **CLAUDE.md**: Project guidance (up-to-date)
- **Gin Framework**: https://gin-gonic.com/
- **Svelte 5**: https://svelte.dev/blog/runes
- **Tailwind CSS 4**: https://tailwindcss.com/
- **Vite**: https://vitejs.dev/
- **JWT**: https://github.com/golang-jwt/jwt/
- **GOWA**: Go WhatsApp Web Multi-Device - see `GOWA-README.md`
- **Noembed API**: https://noembed.com/
