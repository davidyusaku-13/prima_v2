# Codebase Concerns

**Analysis Date:** 2026-01-17

## Tech Debt

**Password Hashing Algorithm:**
- Issue: Using `sha256.Sum256()` for password hashing in `backend/main.go:94`
- Files: `backend/main.go`, `backend/utils/hmac.go`
- Impact: SHA256 is not suitable for password hashing (fast hash, no salt, vulnerable to rainbow tables). Should use bcrypt or argon2.
- Fix approach: Migrate to `golang.org/x/crypto/bcrypt` or `argon2id` with proper salt and work factor.

**Console Logging in Production:**
- Issue: Widespread use of `console.log`, `console.error`, `console.warn` in frontend code
- Files: `frontend/src/App.svelte`, `frontend/src/lib/views/*.svelte`, `frontend/src/lib/components/*.svelte`, `frontend/src/lib/utils/api.js`, `frontend/src/lib/services/sse.js`
- Impact: Clutters browser console, exposes internal errors to users, makes debugging harder in production
- Fix approach: Replace with toast notification system or proper logging abstraction

**Skipped Frontend Test Suites:**
- Issue: Two entire test suites are skipped due to memory leaks in ContentPickerModal
- Files: `frontend/src/lib/components/content/ContentPickerModal.selection.test.js`, `frontend/src/lib/components/content/ContentPickerModal.footer.test.js`
- Impact: 12+ tests not running, component has architectural issues with Svelte 5 `$derived.by()` chains
- Fix approach: Refactor ContentPickerModal to fix reactive cascades or migrate to Playwright E2E tests

**JSON File Persistence:**
- Issue: All data stored in JSON files with `sync.RWMutex` for concurrency
- Files: `backend/main.go`, `backend/data/patients.json`, `backend/data/users.json`, `backend/data/*.json`
- Impact: No backup mechanism, potential data corruption on crashes, won't scale beyond thousands of records
- Fix approach: Migrate to SQLite or PostgreSQL with proper transactions

**Large File Complexity:**
- Issue: Multiple files exceed 500 lines with high cyclomatic complexity
- Files:
  - `backend/handlers/reminder_test.go` (2517 lines)
  - `backend/main.go` (1294 lines)
  - `backend/handlers/content.go` (1205 lines)
  - `backend/handlers/reminder.go` (1077 lines)
  - `frontend/src/lib/components/content/ContentPickerModal.svelte` (1123 lines)
  - `frontend/src/App.svelte` (856 lines)
- Impact: Hard to maintain, high bug risk, difficult to test
- Fix approach: Extract handler functions into separate files, split components

## Known Bugs

**SSE Token Exposure:**
- Issue: JWT token passed in query parameter for SSE connections (`?token=...`)
- Files: `backend/main.go:380`, `frontend/src/lib/services/sse.js`
- Symptom: Tokens logged in server logs, browser history, proxy logs
- Trigger: Any SSE connection establishes with token in URL
- Workaround: Use cookie-based authentication with SameSite=Strict

**Memory Leaks in ContentPickerModal:**
- Issue: Component has memory leaks in Svelte 5 `$derived.by()` reactive chains
- Files: `frontend/src/lib/components/content/ContentPickerModal.svelte`
- Symptom: OOM errors in Vitest happy-dom environment, tests time out
- Trigger: Rendering component with `selectedContent` prop
- Workaround: Tests use `describe.skip()` to skip entire suites

**Date Format Parsing Inconsistency:**
- Issue: Multiple date format parsing patterns throughout codebase
- Files: `backend/main.go:875-880`, `backend/services/scheduler.go:188-195`, `backend/handlers/analytics.go:138-141`
- Symptom: Reminders may not fire if due date format doesn't match expected pattern
- Trigger: Creating reminders with non-standard date formats
- Current mitigation: Multiple fallback parsing attempts

## Security Considerations

**Weak Password Storage:**
- Risk: Passwords hashed with SHA256 without salt
- Files: `backend/main.go:93-96`
- Current mitigation: None - vulnerable to rainbow tables
- Recommendations: Implement bcrypt with cost factor of 12+, add pepper salt stored separately

**JWT Token in URL Query String:**
- Risk: Tokens exposed in logs, browser history, referer headers
- Files: `backend/main.go:530-567` (sseAuthMiddleware)
- Current mitigation: None
- Recommendations: Use HTTP-only cookies for SSE authentication

**No Rate Limiting:**
- Risk: No rate limiting on public or authenticated endpoints
- Files: `backend/main.go` (all routes)
- Current mitigation: Circuit breaker only for GOWA service
- Recommendations: Add middleware-based rate limiting (e.g., golang.org/x/time/rate)

**CORS Origin Configuration:**
- Risk: CORS origin configurable but defaults to localhost
- Files: `backend/config/config.go:170-172`
- Current mitigation: Configurable in config.yaml
- Recommendations: Validate origin against allowlist in production

## Performance Bottlenecks

**JSON File Writes:**
- Problem: Every data change triggers file write via `saveData()` goroutine
- Files: `backend/main.go:755-776`
- Cause: Atomic write pattern with tmp file rename, but no batching
- Improvement path: Add write coalescing or batch multiple changes

**Patient List Loading:**
- Problem: All patients loaded into memory for every GET /api/patients request
- Files: `backend/main.go:953-988`
- Cause: In-memory map iterated completely for each request
- Improvement path: Implement pagination, caching with TTL

**Image Processing Synchronous:**
- Problem: Image uploads block request handler during resize operations
- Files: `backend/handlers/content.go:737-753`
- Cause: `imaging.Save()` called synchronously for 3 image sizes
- Improvement path: Move to background goroutine pool

**Svelte Component Initial Render:**
- Problem: ContentPickerModal (1123 lines) loads all content on mount
- Files: `frontend/src/lib/components/content/ContentPickerModal.svelte`
- Cause: Single component handles filtering, selection, pagination
- Improvement path: Lazy load content, virtualize list rendering

## Fragile Areas

**ContentStore Locking:**
- Files: `backend/handlers/content.go`
- Why fragile: Multiple mutex locks (`Mu.RLock()`, `Mu.RUnlock()`) on Articles and Videos, inconsistent lock ordering
- Safe modification: Always acquire Articles lock before Videos lock, use defer for unlock
- Test coverage: Moderate - content_test.go covers main paths

**Reminder Scheduler Timing:**
- Files: `backend/services/scheduler.go`
- Why fragile: Time-based scheduling with quiet hours, retry logic, multiple date format support
- Safe modification: Add integration tests with time Travel pattern
- Test coverage: Good - scheduler_test.go exists but limited edge cases

**GOWA Circuit Breaker:**
- Files: `backend/services/gowa.go:18-110`
- Why fragile: State machine with open/half-open/closed transitions: Add state transition validation
- Safe modification
- Test coverage: Good - gowa_test.go covers circuit breaker

**SSE Connection Management:**
- Files: `backend/handlers/sse.go`
- Why fragile: Concurrent client map access, connection lifecycle management
- Safe modification: Use channel-based shutdown coordination
- Test coverage: Limited - sse_test.go basic cases only

## Scaling Limits

**In-Memory User Store:**
- Current capacity: ~10,000 users (memory-dependent)
- Limit: Single server instance, no clustering
- Scaling path: Redis for session storage, database for users

**Patient Data File:**
- Current capacity: ~5,000-10,000 patients (file I/O bound)
- Limit: Single JSON file, lock contention under load
- Scaling path: SQLite with proper indexing, eventually PostgreSQL

**GOWA Connection Pool:**
- Current capacity: 1 connection per server instance
- Limit: Sequential message sending, no parallelization
- Scaling path: Connection pool with goroutine workers

## Dependencies at Risk

**github.com/gin-gonic/gin:**
- Risk: Active maintenance, no immediate concern
- Impact: Core web framework
- Migration plan: Migration to Fiber or Echo straightforward

**github.com/golang-jwt/jwt/v5:**
- Risk: Active maintenance
- Impact: Authentication
- Migration plan: Minimal changes needed for v4->v5

**github.com/disintegration/imaging:**
- Risk: Low maintenance activity
- Impact: Image processing for content uploads
- Migration plan: Migrate to `github.com/imgproxy/imgproxy-go` or standard library

## Missing Critical Features

**User Password Reset:**
- Problem: No password reset flow
- Blocks: Self-service account recovery
- Impact: Admin must manually reset, security risk

**Data Export/Backup:**
- Problem: No automated backup or export functionality
- Blocks: Disaster recovery, data portability
- Impact: Data loss risk on server failure

**Multi-Language Content Management:**
- Problem: Content (articles, videos) stored in single language
- Blocks: Localization expansion beyond EN/ID
- Impact: Cannot support additional locales

## Test Coverage Gaps

**Frontend Components - No Unit Tests:**
- What's not tested: Most Svelte components have no unit tests
- Files: `frontend/src/lib/views/*.svelte`, `frontend/src/lib/components/*.svelte` (39+ components)
- Risk: UI regressions undetected, refactoring dangerous
- Priority: High

**Backend Integration Tests:**
- What's not tested: HTTP handler integration tests incomplete
- Files: Missing integration tests for webhook, SSE handlers under load
- Risk: Production issues with concurrent access patterns
- Priority: Medium

**Error Path Testing:**
- What's not tested: Error responses in many handlers
- Files: `backend/handlers/*.go`
- Risk: Error messages leak internal details, improper HTTP status codes
- Priority: Medium

**Frontend - No E2E Tests:**
- What's not tested: Full user flows, critical paths
- Files: N/A
- Risk: Integration issues between components
- Priority: Medium (consider Playwright)

---

*Concerns audit: 2026-01-17*
