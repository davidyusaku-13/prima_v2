# Coding Conventions

**Analysis Date:** 2025-01-17

## Naming Patterns

### JavaScript/Svelte (Frontend)

**Files:**
- Components: `PascalCase.svelte` (e.g., `PatientListPane.svelte`, `Toast.svelte`)
- Stores: `camelCase.svelte.js` (e.g., `delivery.svelte.js`, `toast.svelte.js`)
- Utilities: `camelCase.js` (e.g., `i18nMock.js`)
- Tests: `filename.test.js` (co-located with source)

**Functions:**
- `camelCase` for all functions
- Prefix with action verb: `getStatus()`, `updateDelivery()`, `handleClick()`

**Variables:**
- `camelCase` for variables and constants
- `CONSTANT_CASE` for true constants (e.g., `API_TIMEOUT`, `MAX_RETRIES`)

**Types/Interfaces:**
- `PascalCase` for any type definitions (if using TypeScript)

**Svelte Stores:**
- `$` prefix for store subscriptions: `auth.js` exports store → use as `$auth`
- Store files named descriptively: `delivery.svelte.js`, `toast.svelte.js`

### Go (Backend)

**Files:**
- Single word or camelCase: `config.go`, `health.go`, `utils.go`

**Functions:**
- `PascalCase` for exported functions: `Load()`, `Validate()`, `NewHealthHandler()`
- `camelCase` for unexported: `getQueueCounts()`, `applyDefaults()`

**Variables:**
- `camelCase` for local variables and parameters
- `PascalCase` for package-level exported variables

**Constants:**
- `PascalCase` for exported: `DefaultPort = 8080`
- `camelCase` for unexported

**Types/Structs:**
- `PascalCase` for all type names: `Config`, `HealthHandler`, `PatientStore`

**Interfaces:**
- Single-method interfaces: `Reader`, `Writer`, `Handler`
- Or suffix for multi-method: `HealthHandler` interface

**Acronyms:**
- Use uppercase: `URL`, `ID`, `API`, `HTTPResponse`, `parseURL`
- Not: `Url`, `Id`, `Api`, `HttpResponse`

---

## Code Style

### JavaScript/Svelte (Frontend)

**Formatting:**
- ESLint config: `/media/davidyusaku/Windows/BACKUP/Portfolio/Web/prima_v2/frontend/.eslintrc.svelte5.js`
- Use `.eslintrc.js` that extends `"./.eslintrc.svelte5.js"`
- Arrays: `['item1', 'item2']` (no spacing inside brackets)
- Objects: `{ key: 'value', nested: { inner: true } }` (spacing inside braces)

**Linting Rules:**
- `no-console`: Warn, only `console.warn` and `console.error` allowed
- `no-console.log`: Forbidden in production code
- `prefer-template`: Error (use template literals over concatenation)
- `no-param-reassign`: Error with exceptions for `state`, `store`, `obj`

**Svelte 5 Specific:**
- Use runes: `$state()`, `$derived()`, `$effect()`
- Props: `let { title = 'Default' } = $props()`
- Event handlers: `<button onclick={handleClick}>` (not `on:`)
- Event forwarding: `<Child oncustom={handler}>` (not `<Child on:custom={handler}>`)

**State Management:**
- Never mutate state directly:
  - WRONG: `items.push(x)`
  - RIGHT: `items = [...items, x]`
- Always create new references for Svelte 5 reactivity

**Accessibility:**
- Required `aria-*` attributes on interactive elements
- Use `role="status"`, `role="alert"`, `aria-live="polite"`
- Screen reader-only text with `.sr-only` class
- Avoid `tabindex` without proper `role`

**SvelteKit PROHIBITED:**
- No `import { goto } from '$app/navigation'`
- No `import { page } from '$app/stores'`
- No `import { browser } from '$app/environment'`

**Use native browser APIs:**
- Navigation: `window.location.href = '...'`
- SSE: `new EventSource('/api/sse')`
- Storage: `localStorage.getItem()`, `localStorage.setItem()`

### Go (Backend)

**Formatting:**
- Use `gofmt -w` on save (no manual formatting)
- Standard Go toolchain handles style

**Imports:**
Grouped and sorted automatically:
```go
import (
    "fmt"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)
```

**Error Handling:**
- Return errors as values, never panic in production
- Wrap errors with context: `fmt.Errorf("failed to load config: %w", err)`
- Check errors immediately: `if err != nil { return err }`

**Types and Structs:**
- Use explicit types, avoid `var x` without type
- Tags for JSON/YAML: `json:"user_id,omitempty" yaml:"user_id"`

**Code Patterns:**
- Receiver on value: `func (h Handler) Method()` unless mutation needed
- Embedding for composition: `type Handler struct{ *Service }`
- `sync.RWMutex` for concurrent read/write
- Avoid global mutable state

---

## Import Organization

### Frontend Imports

**Order:**
1. Node built-ins: `import fs from 'fs'` (rare)
2. External packages: `import { render } from '@testing-library/svelte'`
3. $lib aliases: `import { toastStore } from '$lib/stores/toast.svelte.js'`
4. Relative imports: `import { createI18nMock } from '../../test-utils/i18nMock.js'`

**Path Aliases:**
- `$lib/*` maps to `/media/davidyusaku/Windows/BACKUP/Portfolio/Web/prima_v2/frontend/src/lib/*`
- Defined in `vite.config.js`

**Example:**
```javascript
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/svelte';
import { toastStore } from '$lib/stores/toast.svelte.js';
import { createI18nMock } from '../../test-utils/i18nMock.js';
```

### Backend Imports

**Order:**
1. Standard library
2. Third-party packages
3. Local application imports

**Example:**
```go
import (
    "fmt"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"

    "github.com/davidyusaku-13/prima_v2/models"
    "github.com/davidyusaku-13/prima_v2/services"
)
```

---

## Error Handling

### Frontend

**Console Usage:**
- `console.warn()` for warnings
- `console.error()` for errors
- `console.log` is FORBIDDEN

**Error Boundaries:**
- Use try/catch for async operations
- Show user-friendly error messages via toast store

**Example:**
```javascript
async function fetchPatient(id) {
    try {
        const response = await fetch(`/api/patients/${id}`);
        if (!response.ok) throw new Error('Failed to fetch patient');
        return await response.json();
    } catch (error) {
        console.error('Failed to fetch patient:', error);
        toastStore.add('Gagal memuat data pasien', { type: 'error' });
        throw error;
    }
}
```

### Backend

**Error Returns:**
```go
func (h *HealthHandler) GetHealthDetailed(c *gin.Context) {
    role := c.GetString("role")
    if role != "admin" && role != "superadmin" {
        c.JSON(http.StatusForbidden, gin.H{"error": "admin access required", "code": "ADMIN_REQUIRED"})
        return
    }
    // ...
}
```

**Error Wrapping:**
```go
data, err := os.ReadFile(path)
if err != nil {
    return nil, fmt.Errorf("failed to read config file: %w", err)
}
```

---

## Logging

### Frontend

**Framework:** Console-based (`console.warn`, `console.error`)

**Pattern:** Minimal logging, prefer user feedback via UI

### Backend

**Framework:** `log/slog` (structured logging)

**Example:**
```go
logger.Info("patient reminder sent",
    "reminder_id", reminder.ID,
    "patient_id", patient.ID,
    "status", "sent",
)
```

**Masking for NFR-S6:**
```go
// MaskPhone masks phone number for logging
// Input: "628123456789" → Output: "628***789"
func MaskPhone(phone string) string {
    // ... implementation
}
```

---

## Comments

### When to Comment

**Frontend:**
- JSDoc for component props and complex functions
- Comments for Svelte 5 reactivity notes (CRITICAL comments)
- File header comment explaining component purpose

**Example (Svelte 5):**
```javascript
/**
 * Delivery status store using Svelte 5 runes
 * Manages real-time delivery status updates via SSE
 *
 * CRITICAL: This is Vite + Svelte 5, NOT SvelteKit!
 * - Uses Svelte 5 runes ($state, $derived)
 * - No legacy reactive statements ($:)
 * - No SvelteKit imports
 */
```

### JSDoc/TSDoc

**Props Documentation:**
```javascript
/**
 * Add a new toast notification
 * @param {string} message - The message to display
 * @param {Object} options - Toast options
 * @param {string} options.type - Toast type: 'success' | 'error' | 'warning' | 'info'
 * @param {Object} options.action - Action button: { label: string, onClick: function }
 * @param {number} options.duration - Auto-dismiss duration in ms (0 = no auto-dismiss)
 * @returns {number} Toast ID
 */
```

---

## Function Design

### Frontend

**Size:** Keep functions focused, single responsibility

**Parameters:**
- Destructuring for objects: `function handleClick({ id, name })`
- Default values for optional params

**Return Values:**
- Always return for clarity
- Use explicit returns

**Example:**
```javascript
/**
 * Get delivery status for a reminder
 * @param {string} reminderId - The reminder ID
 * @returns {string|null} Status or null if not found
 */
function getStatus(reminderId) {
    return deliveryStatuses[reminderId]?.status || null;
}
```

### Backend

**Size:** Keep functions focused (typically < 50 lines)

**Parameters:**
- Explicit types required
- Pointer receivers for mutation, value for read-only

**Return Values:**
- Return error as last value: `(*Config, error)`
- Multiple returns with clear purpose

**Example:**
```go
// Load reads configuration from a YAML file and returns a Config struct
func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    // ...
}
```

---

## Module Design

### Frontend Exports

**Named exports preferred:**
```javascript
// delivery.svelte.js
export const deliveryStore = new DeliveryStore();

// Toast.svelte
export default Toast;
```

**Default exports for components:**
```javascript
import ContentPreviewPanel from './ContentPreviewPanel.svelte';
```

**Barrel Files:** Not used; import directly from feature directories

### Backend Package Structure

**Exported types/functions:**
```go
package config

type Config struct { ... }
func Load(path string) (*Config, error) { ... }
func LoadOrDefault(path string) *Config { ... }
```

**Unexported (private):**
```go
func (c *Config) applyDefaults() { ... }
```

---

*Convention analysis: 2025-01-17*
