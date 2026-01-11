# Architecture Decision Records (ADRs) - PRIMA

**Documenting key technical decisions and their rationale**

---

## What are ADRs?

Architecture Decision Records capture important technical decisions, their context, and consequences. They help future developers understand *why* the codebase is the way it is.

---

## ADR-001: Svelte 5 + Vite over SvelteKit

**Date:** 2025 (Original Development)
**Status:** Accepted
**Context:**

PRIMA needed a modern, reactive frontend framework. The options considered were:
1. SvelteKit (full-stack framework)
2. Svelte 5 + Vite (SPA approach)
3. React + Vite

**Decision:**

Use **Svelte 5 with Vite directly**, without SvelteKit.

**Rationale:**

- **No SSR needed:** PRIMA is a dashboard application, not a public-facing website. SEO is irrelevant.
- **Simpler deployment:** No Node.js server required for the frontend; static files only.
- **Lighter bundle:** No SvelteKit runtime overhead.
- **Cleaner separation:** Backend is Go; mixing with SvelteKit's server functions would create confusion.
- **Control over routing:** Custom localStorage-based routing fits the dashboard UX better.

**Consequences:**

- Cannot use `$app/*` imports (goto, page, browser).
- Must implement manual routing via localStorage and conditional rendering.
- No server-side rendering; initial load may show blank state briefly.
- Developers familiar with SvelteKit may initially use prohibited patterns.

**Migration Path:**

If SEO or SSR becomes necessary, evaluate SvelteKit migration. The component structure is compatible.

---

## ADR-002: JSON File Persistence over Database

**Date:** 2025 (Original Development)
**Status:** Accepted
**Context:**

PRIMA needed to store patient data, reminders, articles, and user accounts. Options considered:
1. SQLite (embedded database)
2. PostgreSQL (production database)
3. JSON files (file-based persistence)

**Decision:**

Use **JSON file persistence** with `sync.RWMutex` for thread safety.

**Rationale:**

- **Simple deployment:** No database server to install, configure, or maintain.
- **Human-readable:** Data can be inspected and edited directly if needed.
- **Easy backup:** Copy files to backup.
- **Sufficient for scale:** Target is <1000 patients per deployment.
- **Fast prototyping:** Changes to schema don't require migrations.

**Consequences:**

- All data loaded into memory at startup.
- No complex queries (must filter in Go code).
- Horizontal scaling not possible (single server only).
- File corruption possible if write fails mid-operation.
- No ACID transactions.

**Mitigations:**

- Atomic writes via temp file + rename pattern.
- RWMutex prevents concurrent write corruption.
- Regular backups recommended.

**Migration Path:**

When patient count exceeds 5000 or multi-server deployment is needed:
1. Add PostgreSQL/SQLite support with repository pattern.
2. Create migration script from JSON to database.
3. Keep RWMutex pattern for in-memory caching layer.

---

## ADR-003: Go/Gin over Node.js for Backend

**Date:** 2025 (Original Development)
**Status:** Accepted
**Context:**

Backend framework selection. Options considered:
1. Node.js with Express
2. Python with FastAPI
3. Go with Gin

**Decision:**

Use **Go with Gin framework**.

**Rationale:**

- **Performance:** Go compiles to native binary; lower memory, faster execution.
- **Type safety:** Compile-time type checking reduces runtime errors.
- **Concurrency:** Goroutines handle SSE connections and scheduler efficiently.
- **Single binary:** No runtime dependencies to install on server.
- **Standard library:** HTTP, JSON, crypto all in stdlib.

**Consequences:**

- Smaller talent pool than Node.js.
- More verbose than dynamic languages.
- Must handle JSON marshaling explicitly.

**Migration Path:**

No migration anticipated. Go backend is stable and performant.

---

## ADR-004: Custom Routing over svelte-spa-router

**Date:** 2025 (Original Development)
**Status:** Accepted
**Context:**

SPA routing solution needed. Options considered:
1. svelte-spa-router (npm package)
2. tinro (npm package)
3. Custom localStorage-based routing

**Decision:**

Implement **custom routing** using localStorage and conditional rendering.

**Rationale:**

- **Minimal dependencies:** No external router library needed.
- **Dashboard UX:** Users rarely share URLs; deep linking less important.
- **Simple logic:** Switch statement in App.svelte handles all views.
- **State preservation:** localStorage remembers view across refreshes.

**Consequences:**

- No URL-based routing (URLs don't reflect current view).
- Back button doesn't work as expected.
- Cannot bookmark specific views.
- All routing logic in App.svelte.

**Trade-offs:**

For a dashboard application used by a small team of volunteers, the simplicity outweighs the limitations. Public-facing apps would need proper URL routing.

**Migration Path:**

If URL-based routing becomes necessary:
1. Install svelte-spa-router or tinro.
2. Replace localStorage checks with route definitions.
3. Update navigation to use router's goto/push functions.

---

## ADR-005: GOWA over Twilio for WhatsApp

**Date:** 2025 (Original Development)
**Status:** Accepted
**Context:**

WhatsApp message delivery needed. Options considered:
1. Twilio WhatsApp API
2. WhatsApp Business API (direct)
3. GOWA (self-hosted gateway)

**Decision:**

Use **GOWA** (Go WhatsApp Gateway).

**Rationale:**

- **Cost:** No per-message fees; GOWA is self-hosted and free.
- **Indonesia focus:** GOWA works well with Indonesian phone numbers.
- **Simplicity:** HTTP API with Basic Auth.
- **Control:** Self-hosted; no third-party dependency.

**Consequences:**

- Must maintain GOWA instance.
- Less reliable than Twilio (unofficial API).
- May break if WhatsApp changes their protocol.
- No message templates or business features.

**Mitigations:**

- Circuit breaker pattern handles GOWA unavailability.
- Retry with exponential backoff.
- Quiet hours prevent message flooding.

**Migration Path:**

If GOWA becomes unreliable:
1. Implement Twilio adapter with same interface.
2. Configure API credentials in config.yaml.
3. Switch adapter via configuration.

---

## ADR-006: SHA256 Password Hashing

**Date:** 2025 (Original Development)
**Status:** Accepted (with caveat)
**Context:**

Password storage mechanism needed. Options considered:
1. bcrypt
2. Argon2
3. SHA256

**Decision:**

Use **SHA256** hashing with Base64 encoding.

**Rationale:**

- **Simplicity:** SHA256 is in Go stdlib; no external dependencies.
- **Low-risk context:** Internal healthcare volunteer system, not public-facing.
- **Fast implementation:** Prototype speed was prioritized.

**Consequences:**

- Less secure than bcrypt/Argon2 (no salt, fast to brute-force).
- Not suitable for high-security applications.
- Security auditors may flag this.

**Caveat:**

This decision prioritized development speed. For production deployments handling sensitive data:

```go
// Replace SHA256 with bcrypt:
import "golang.org/x/crypto/bcrypt"

hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
err := bcrypt.CompareHashAndPassword(hash, []byte(password))
```

**Migration Path:**

1. Add bcrypt dependency.
2. Update User struct and password verification.
3. Force password reset for all users.
4. Remove SHA256 code.

---

## ADR-007: Circuit Breaker for External Services

**Date:** 2025 (Original Development)
**Status:** Accepted
**Context:**

GOWA integration could fail; needed resilience pattern.

**Decision:**

Implement **circuit breaker pattern** with exponential backoff retry.

**Configuration:**
- 5 failures → circuit opens
- 5-minute cooldown
- Half-open: single test request
- Retry delays: 1s, 5s, 30s, 2m, 10m

**Rationale:**

- Prevents cascade failures when GOWA is down.
- Automatic recovery when service returns.
- Graceful degradation for users.

**Consequences:**

- Messages queue during outage.
- May delay message delivery up to 10 minutes on retries.
- State lost on server restart (in-memory only).

**Implementation:**

```go
type CircuitBreaker struct {
    failureCount int
    state        string // closed, open, half-open
    lastFailure  time.Time
    cooldown     time.Duration
}
```

---

## ADR-008: Server-Sent Events over WebSockets

**Date:** 2025 (Original Development)
**Status:** Accepted
**Context:**

Real-time delivery status updates needed. Options:
1. WebSockets (bidirectional)
2. Server-Sent Events (unidirectional)
3. Long polling

**Decision:**

Use **Server-Sent Events (SSE)**.

**Rationale:**

- **Simpler:** SSE is HTTP; no upgrade handshake complexity.
- **Unidirectional fits:** Clients only need to receive updates, not send.
- **Automatic reconnection:** Browser handles reconnection automatically.
- **Firewall friendly:** Uses standard HTTP/HTTPS ports.

**Consequences:**

- Unidirectional only (client cannot send via SSE).
- Connection limits per domain (6 in HTTP/1.1).
- No binary data support.

**Implementation:**

```go
// Backend
c.Writer.Header().Set("Content-Type", "text/event-stream")
fmt.Fprintf(c.Writer, "event: delivery.status.updated\ndata: %s\n\n", json)

// Frontend
const es = new EventSource('/api/sse/delivery-status?token=' + token);
es.addEventListener('delivery.status.updated', handler);
```

---

## Template for New ADRs

```markdown
## ADR-XXX: Title

**Date:** YYYY-MM-DD
**Status:** Proposed | Accepted | Deprecated | Superseded by ADR-YYY

**Context:**

What is the issue or decision that we're facing?

**Decision:**

What is the decision we've made?

**Rationale:**

Why did we make this decision over alternatives?

**Consequences:**

What are the positive and negative results of this decision?

**Migration Path:**

If this decision needs to change, what's the path forward?
```

---

**Last Updated:** January 11, 2026
