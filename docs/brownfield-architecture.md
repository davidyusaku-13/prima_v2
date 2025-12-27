# Prima v2 Brownfield Architecture Document

## Introduction

This document captures the current state of the Prima v2 codebase, a landing page website built with Go/Gin backend and Svelte 5 + Vite frontend. It serves as a reference for AI agents working on enhancements, bug fixes, or feature development.

### Document Scope

Comprehensive documentation of entire system - no specific PRD or enhancement in focus.

### Change Log

| Date       | Version | Description                          | Author  |
| ---------- | ------- | ------------------------------------ | --------|
| 2025-12-27 | 1.0     | Initial brownfield analysis          | Claude  |
| 2025-12-27 | 1.1     | Updated to landing page architecture | Claude  |

---

## Quick Reference - Key Files and Entry Points

### Critical Files for Understanding the System

| Purpose             | File Path                 |
| ------------------- | ------------------------- |
| Backend Entry       | `backend/main.go`         |
| Frontend Entry      | `frontend/src/main.js`    |
| Main UI Component   | `frontend/src/App.svelte` |
| Frontend Config     | `frontend/vite.config.js` |
| CSS/Tailwind Config | `frontend/src/app.css`    |
| Project Guidance    | `CLAUDE.md`               |

### Project Root Structure

```
prima_v2/
├── backend/          # Go/Gin backend (port 8080)
├── frontend/         # Svelte 5 + Vite frontend (port 5173)
├── docs/             # Documentation
├── CLAUDE.md         # Claude Code guidance
├── .git/             # Git repository
└── node_modules/     # Root dependencies
```

---

## High Level Architecture

### Technical Summary

Prima v2 is a marketing landing page for a fictional product called "Prima". The application demonstrates:

1. **Backend**: A minimal Go/Gin server providing a health check endpoint
2. **Frontend**: A Svelte 5 single-page application with a complete landing page design using Tailwind CSS

The architecture is intentionally minimal - static content with simulated form interactions.

### Actual Tech Stack

| Category     | Technology        | Version    | Notes                                |
| ------------ | ----------------- | ---------- | ------------------------------------ |
| **Backend**  |                   |            |                                      |
| Runtime      | Go                | 1.25+      | Standard Go runtime                  |
| Framework    | Gin               | Latest     | HTTP web framework                   |
| CORS         | gin-contrib/cors  | Latest     | Middleware for cross-origin requests |
| **Frontend** |                   |            |                                      |
| Runtime      | Node.js           | 16+        | For Vite development server          |
| Framework    | Svelte            | 5.43.8     | Latest Svelte 5                      |
| Build Tool   | Vite              | 7.2.4      | Fast development server and bundler  |
| Styling      | Tailwind CSS      | 4.1.18     | CSS-first utility framework          |
| Integration  | @tailwindcss/vite | 4.1.18     | Vite plugin for Tailwind 4           |

### Repository Structure

- **Type**: Polyrepo (separate `backend/` and `frontend/` directories)
- **Package Manager**: Go modules (backend), npm (frontend)
- **Version Control**: Git
- **Code Guidance**: Claude Code configuration in `CLAUDE.md`

---

## Source Tree and Module Organization

### Backend Structure (`backend/`)

```
backend/
├── main.go          # Complete application (40 lines)
├── go.mod           # Go module definition
└── storage/         # [DEPRECATED] Old storage package, no longer used
    ├── go.mod
    ├── storage.go
    └── data/items.json
```

**Note**: The `storage/` directory is deprecated and can be removed in future cleanup.

### Frontend Structure (`frontend/`)

```
frontend/
├── src/
│   ├── main.js              # Svelte 5 mount entry point
│   ├── App.svelte           # Complete landing page (213 lines)
│   ├── app.css              # Tailwind CSS @import
│   ├── assets/
│   │   └── svelte.svg       # Svelte logo
│   └── lib/
│       └── Counter.svelte   # Unused demo component (can be removed)
├── index.html               # HTML entry point
├── package.json             # Dependencies and scripts
├── vite.config.js           # Vite + Svelte + Tailwind config
└── node_modules/            # Dependencies
```

---

## Backend Implementation Details

### Entry Point: `backend/main.go`

**All backend logic resides in a single, minimal file.**

```go
package main

import (
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

func main() {
    router := gin.Default()
    // CORS configuration...
    router.GET("/api/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
    // Graceful shutdown...
    router.Run(":8080")
}
```

### API Endpoints

| Method | Endpoint       | Response                    | Purpose       |
| ------ | -------------- | --------------------------- | ------------- |
| GET    | `/api/health`  | `{"status": "ok"}`          | Health check  |

### CORS Configuration

CORS is explicitly configured for frontend origin only:

```go
AllowOrigins: []string{"http://localhost:5173"}
```

**Constraint**: The frontend URL is hardcoded.

### Graceful Shutdown

The backend handles SIGINT and SIGTERM for graceful shutdown:

```go
go func() {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down server...")
}()
```

---

## Frontend Implementation Details

### Entry Point: `frontend/src/main.js`

Uses Svelte 5's new mounting API:

```javascript
import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'

const app = mount(App, {
  target: document.getElementById('app'),
})

export default app
```

### Main Component: `frontend/src/App.svelte`

**Single component handles entire landing page UI.**

#### State Management (Standard Svelte)

```javascript
let email = '';      // Newsletter email input
let subscribed = false;   // Subscription success state
let loading = false;      // Form submission state
```

#### Functions

| Function          | Purpose                                      |
| ----------------- | -------------------------------------------- |
| `handleSubscribe()` | Simulates email subscription (1s delay)    |

#### Landing Page Sections

1. **Navigation**
   - Logo ("Prima")
   - Menu links (Features, About, Contact)
   - "Get Started" CTA button

2. **Hero Section**
   - Headline: "Build better products, faster."
   - Subtitle describing the product
   - Two CTA buttons: "Start Free Trial", "Watch Demo"

3. **Stats Section**
   - 4 stat cards: 10k+ Users, 99.9% Uptime, 50M+ API Requests, 24/7 Support

4. **Features Section**
   - 3 feature cards: Lightning Fast, Secure by Default, Team Collaboration

5. **About Section**
   - Company description
   - Stats grid (2024 Founded, 150+ Team Members, 500+ Clients, Global Team)

6. **Newsletter Section**
   - Email input form
   - Subscribe button with loading state
   - Success message after "subscription"

7. **Footer**
   - Logo
   - Link columns (Features, About, Contact, Privacy, Terms)
   - Copyright notice

### Tailwind CSS Integration

- **Version**: Tailwind CSS 4 (CSS-first configuration)
- **Config**: `@tailwindcss/vite` plugin in `vite.config.js`
- **Import**: `@import "tailwindcss"` in `app.css`
- **Usage**: Utility classes directly in Svelte template

### Unused Components

- `frontend/src/lib/Counter.svelte` - Demo component, not imported
- `backend/storage/` directory - Deprecated, no longer used

---

## Technical Debt and Known Issues

### Critical Technical Debt

1. **Hardcoded Values**
   - Newsletter email subscription is simulated (no backend endpoint)
   - No actual email collection/storage
   - All statistics are hardcoded

2. **No Form Backend Integration**
   - Newsletter form uses `setTimeout` to simulate submission
   - No actual email service integration (SendGrid, Mailchimp, etc.)

3. **Single-File Frontend Component**
   - All landing page sections in one file
   - No reusable UI components extracted

4. **No Responsive Menu Mobile**
   - Navigation menu hidden on mobile (uses `hidden md:flex`)
   - No mobile hamburger menu

### Workarounds and Gotchas

- **CORS Strict**: Backend only accepts `http://localhost:5173`
- **No Persistence**: Landing page is static, no user data stored
- **No Analytics**: No tracking or analytics integration
- **Mock Statistics**: All company stats are placeholder values

---

## Integration Points and External Dependencies

### External Services

| Service   | Purpose        | Integration Type | Key Files                 |
| --------- | -------------- | ---------------- | ------------------------- |
| Go stdlib | Core runtime   | Built-in         | `backend/main.go`         |
| Node.js   | Frontend build | Runtime          | `frontend/package.json`   |
| Vite      | Dev server     | npm dependency   | `frontend/vite.config.js` |

### Internal Integration Points

- **Frontend → Backend**: Only `/api/health` endpoint used (optional)
- **Svelte → Tailwind**: CSS utility classes in Svelte templates
- **Vite → Svelte**: Build plugin in `vite.config.js`

### Ports and Hosts

| Service  | Host:Port      | Configuration Location            |
| -------- | -------------- | --------------------------------- |
| Backend  | localhost:8080 | `backend/main.go:38`              |
| Frontend | localhost:5173 | Vite default                      |

---

## Development and Deployment

### Local Development Setup

#### Prerequisites

- Go 1.25+
- Node.js 16+
- npm (comes with Node)

#### Starting the Backend

```bash
cd backend
go run main.go
```

Backend runs on `http://localhost:8080`

#### Starting the Frontend

```bash
cd frontend
npm run dev
```

Frontend runs on `http://localhost:5173`

#### Accessing the App

Open `http://localhost:5173` in browser. Backend health check will work automatically.

### Build and Deployment Process

- **Frontend Build**: `npm run build` (outputs to `frontend/dist/`)
- **Backend Build**: `go build -o prima_v2 main.go` (standalone binary)
- **No CI/CD**: Manual deployment required
- **No Docker**: No containerization configured

### Environment Variables

**No environment variables currently used** - all configuration is hardcoded.

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
npm test   # Not available
```

**Gap**: The project has no automated tests.

---

## Future Enhancement Considerations

### Files That May Need Modification

| Enhancement Type         | Files to Modify                              |
| ------------------------ | -------------------------------------------- |
| Add email service        | `frontend/src/App.svelte` (add API call)     |
| Add mobile menu          | `frontend/src/App.svelte` (nav section)      |
| Add analytics            | `frontend/src/App.svelte` + `index.html`     |
| Add form validation      | `frontend/src/App.svelte` (handleSubscribe)  |
| Extract UI components    | Create `frontend/src/components/` directory  |

### New Files/Modules That May Be Needed

- `frontend/src/components/Navigation.svelte` - Reusable nav component
- `frontend/src/components/Footer.svelte` - Reusable footer component
- `frontend/src/components/NewsletterForm.svelte` - Newsletter with real backend
- `backend/api/newsletter.go` - Email subscription endpoint
- `.env` - Environment variables for services

### Integration Considerations for Enhancements

- **Email service**: Connect SendGrid, Mailchimp, or similar for newsletter
- **Analytics**: Add Google Analytics, Plausible, or similar
- **Backend endpoints**: Create `/api/newsletter` for form submissions
- **Form validation**: Add proper email validation and error handling

---

## Appendix - Useful Commands and Scripts

### Backend Commands

```bash
# Run backend
cd backend && go run main.go

# Build backend binary
cd backend && go build -o prima_v2 main.go

# Check Go syntax
cd backend && go vet

# Clean up deprecated storage
rm -rf backend/storage backend/data
```

### Frontend Commands

```bash
# Install dependencies
cd frontend && npm install

# Start dev server
cd frontend && npm run dev

# Build for production
cd frontend && npm run build

# Preview production build
cd frontend && npm run preview
```

### Debugging and Troubleshooting

- **Backend logs**: Console output from `go run main.go`
- **Frontend logs**: Browser DevTools console
- **Network**: Browser DevTools Network tab
- **CORS errors**: Check browser console, verify backend CORS config

### Common Issues

| Issue                    | Solution                                       |
| ------------------------ | ---------------------------------------------- |
| CORS error in browser    | Ensure backend is running before frontend      |
| Connection refused       | Check backend is on port 8080                  |
| Build fails              | Run `npm install` in frontend directory        |
| Storage errors           | Backend no longer uses storage, can be removed |

---

## References

- **CLAUDE.md**: Project-specific guidance for Claude Code
- **Gin Framework**: https://gin-gonic.com/
- **Svelte 5**: https://svelte.dev/blog/runes
- **Tailwind CSS 4**: https://tailwindcss.com/
- **Vite**: https://vitejs.dev/
