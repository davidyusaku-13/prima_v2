# Technology Stack - PRIMA

**Generated:** January 11, 2026 (Updated)
**Project:** Healthcare Volunteer Dashboard
**Scan Type:** Exhaustive Rescan

---

## Backend (Go/Gin API)

### Core Technologies

| Category             | Technology       | Version               | Justification                                      |
| -------------------- | ---------------- | --------------------- | -------------------------------------------------- |
| **Language**         | Go               | 1.25.5                | Type-safe, compiled, excellent concurrency support |
| **Web Framework**    | Gin              | v1.11.0               | High-performance HTTP router, middleware support   |
| **Authentication**   | JWT              | golang-jwt/jwt v5.2.1 | Stateless auth, 7-day token expiry                 |
| **CORS**             | gin-contrib/cors | v1.7.6                | Cross-origin resource sharing for frontend         |
| **Configuration**    | YAML             | gopkg.in/yaml.v3      | Human-readable config files                        |
| **Image Processing** | imaging          | v1.6.2                | Image upload handling and compression              |

### Architecture Pattern

**Layered Architecture with Handlers/Services Pattern**

- **Handlers Layer:** HTTP request/response handling
- **Services Layer:** Business logic (GOWA integration, scheduler)
- **Models Layer:** Data structures and validation
- **Utils Layer:** Shared utilities (logging, masking, phone validation)

### Data Persistence

- **Storage:** JSON files in `backend/data/`
- **Concurrency:** `sync.RWMutex` for thread-safe file operations
- **Files:**
  - `patients.json` - Patient records
  - `users.json` - User accounts (hashed passwords)
  - `articles.json` - Health articles (Berita)
  - `videos.json` - Educational videos
  - `categories.json` - Content categories
  - `jwt_secret.txt` - JWT signing key

### Configuration

**File:** `config.yaml`

- Server settings (port, CORS)
- GOWA WhatsApp integration
- Circuit breaker patterns
- Retry policies
- Quiet hours enforcement
- Logging configuration

---

## Frontend (Svelte 5 + Vite)

### Core Technologies

| Category                 | Technology              | Version | Justification                                   |
| ------------------------ | ----------------------- | ------- | ----------------------------------------------- |
| **Framework**            | Svelte                  | 5.43.8  | Reactive UI without virtual DOM, Svelte 5 runes |
| **Build Tool**           | Vite                    | 7.2.4   | Fast HMR, optimized bundling                    |
| **Styling**              | Tailwind CSS            | 4.1.18  | Utility-first CSS, rapid development            |
| **Tailwind Plugin**      | @tailwindcss/vite       | 4.1.18  | Native Vite integration                         |
| **Internationalization** | svelte-i18n             | 4.0.1   | EN/ID translations                              |
| **Rich Text Editor**     | Quill                   | 2.0.3   | WYSIWYG editor for article content              |
| **Testing**              | Vitest                  | 4.0.16  | Fast unit testing with Vite integration         |
| **Testing Library**      | @testing-library/svelte | 5.3.1   | Component testing utilities                     |
| **DOM Environment**      | happy-dom               | 20.0.11 | Lightweight DOM for tests                       |

### Architecture Pattern

**Component-Based Architecture with Svelte 5 Runes**

- **Reactive State:** `$state()` rune for component-local state
- **Derived Values:** `$derived()` for computed properties
- **Side Effects:** `$effect()` for lifecycle and DOM interactions
- **Props:** `let { prop = default } = $props()`
- **Events:** Native `onclick`, `onsubmit` (no `on:` directive)

### Project Structure

```
frontend/
├── src/
│   ├── App.svelte           # Root component
│   ├── main.js              # Application bootstrap
│   ├── app.css              # Global Tailwind imports
│   ├── i18n.js              # svelte-i18n configuration
│   ├── lib/                 # Components & utilities
│   ├── locales/             # Translation files
│   │   ├── en.json
│   │   └── id.json
│   └── test/                # Test files
├── vite.config.js           # Vite configuration
├── svelte.config.js         # Svelte preprocessor
└── package.json
```

### Build Configuration

**Vite Aliases:**

- `$lib` → `./src/lib` (import components/utilities)

**Plugins:**

- `@sveltejs/vite-plugin-svelte` - Svelte compilation
- `@tailwindcss/vite` - Tailwind CSS processing

---

## External Integrations

### GOWA (WhatsApp Gateway)

| Aspect             | Details                                                     |
| ------------------ | ----------------------------------------------------------- |
| **Service**        | Golang WhatsApp Web Multi-Device                            |
| **Port**           | 3000 (local development)                                    |
| **Purpose**        | Send patient reminders via WhatsApp                         |
| **Authentication** | Basic auth + HMAC webhook validation                        |
| **Configuration**  | Environment variables (GOWA_ENDPOINT, GOWA_USER, GOWA_PASS) |
| **Features**       | Circuit breaker, retry logic, quiet hours                   |

### YouTube Metadata

| Aspect      | Details                                           |
| ----------- | ------------------------------------------------- |
| **Service** | noembed.com API                                   |
| **Purpose** | Fetch video metadata (title, thumbnail, duration) |
| **Usage**   | Video Edukasi CMS feature                         |

---

## Development Tools

### Backend

- **Testing:** `go test` with standard library
- **Linting:** `gofmt` (automatic formatting)
- **Package Manager:** Go modules (`go.mod`)

### Frontend

- **Testing:** Vitest + Testing Library
- **Package Manager:** Bun (fast npm alternative)
- **Linting:** ESLint (Svelte 5 configuration)

---

## Deployment Considerations

### Backend

- Standalone binary (no runtime dependencies)
- Configuration via `config.yaml` or environment variables
- JSON file persistence (simple, no database required)
- Port: 8080 (configurable)

### Frontend

- Static build output (`bun run build`)
- Served from `/dist` directory
- Requires backend API at runtime
- Port: 5173 (dev), production via reverse proxy

---

## Security

### Authentication

- **JWT Tokens:** 7-day expiry, signed with secret key
- **Password Storage:** Hashed (SHA256)
- **RBAC:** superadmin / admin / volunteer roles

### Data Protection

- **Phone Masking:** Partial masking in logs
- **Email Masking:** Partial masking in logs
- **HMAC Validation:** Webhook integrity verification
- **CORS:** Configured origin restrictions

---

## Architecture Decisions

### Why JSON Files Instead of Database?

- **Simplicity:** No database server setup required
- **Portability:** Easy backup and migration
- **Concurrency:** RWMutex ensures thread-safety
- **Scale:** Suitable for small-to-medium volunteer organizations

### Why Svelte 5 (Not React/Vue)?

- **Performance:** No virtual DOM overhead
- **Developer Experience:** Less boilerplate than React
- **Bundle Size:** Smaller production builds
- **Runes System:** Modern reactive programming model

### Why Vite (Not SvelteKit)?

- **Flexibility:** Direct control over routing and SSR
- **Simplicity:** No framework opinions, pure SPA
- **Speed:** Fast HMR during development

---

**Next:** See [Architecture Documentation](./architecture-backend.md) and [Architecture Documentation](./architecture-frontend.md) for detailed design.
