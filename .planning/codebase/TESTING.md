# Testing Patterns

**Analysis Date:** 2025-01-17

## Test Framework

### Frontend (Svelte 5 + Vite)

**Primary Runner:**
- **Vitest** (version 4.0.16)
- Config: `/media/davidyusaku/Windows/BACKUP/Portfolio/Web/prima_v2/frontend/vitest.config.js`
- Default command: `bun run test` (runs `vitest run`)

**Alternative Runner:**
- **Jest** (version 30.2.0)
- Config: `/media/davidyusaku/Windows/BACKUP/Portfolio/Web/prima_v2/frontend/jest.config.js`
- Used for legacy compatibility

**Assertion Library:**
- Vitest built-in `expect`
- `@testing-library/jest-dom` for DOM assertions

**Run Commands:**
```bash
cd /media/davidyusaku/Windows/BACKUP/Portfolio/Web/prima_v2/frontend
bun run test              # Run all tests (vitest by default)
bun run test -- --run     # Run once (no watch mode)
bun run test api.test.js  # Run specific file
```

### Backend (Go + Gin)

**Runner:**
- Go's built-in `testing` package

**Run Commands:**
```bash
cd /media/davidyusaku/Windows/BACKUP/Portfolio/Web/prima_v2/backend
go test ./...                      # Run all tests
go test -v ./config                # Run specific package
go test -v -run TestLoad           # Run single test by name
go test -v -run "TestLoad|Default" # Run multiple tests (regex)
go test -cover ./...               # With coverage
```

---

## Test File Organization

### Frontend

**Location:**
- Co-located with source files: `src/lib/components/Component.svelte` → `src/lib/components/Component.test.js`
- Test utilities: `src/lib/test-utils/`

**Naming:**
- `*.test.js` for all test files
- Svelte components: `ContentPreviewPanel.test.js`, `PatientListPane.test.js`
- Stores: `delivery.test.js`, `toast.test.js`
- Utilities: `i18nMock.test.js`

**Structure:**
```
src/
  lib/
    components/
      delivery/
        DeliveryStatusBadge.svelte
        DeliveryStatusBadge.test.js
        DeliveryStatusFilter.test.js
      content/
        ContentPreviewPanel.svelte
        ContentPreviewPanel.test.js
    stores/
      delivery.svelte.js
      delivery.test.js
    test-utils/
      i18nMock.js
      i18nMock.test.js
  test/
    setup.js
```

### Backend

**Location:**
- Co-located with source: `config/config.go` → `config/config_test.go`
- Same directory structure as source

**Naming:**
- `*_test.go` suffix
- `config_test.go`, `logger_test.go`, `mask_test.go`

---

## Test Structure

### Frontend (Vitest + Testing Library Svelte)

**Suite Organization:**
```javascript
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, fireEvent, waitFor, cleanup } from '@testing-library/svelte';

describe('ComponentName', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    afterEach(() => {
        vi.restoreAllMocks();
        cleanup();
    });

    describe('Subfeature', () => {
        it('should do something', async () => {
            // Test implementation
        });
    });
});
```

**Store Testing Pattern:**
```javascript
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { deliveryStore } from './delivery.svelte.js';
import { sseService } from '$lib/services/sse.js';

vi.mock('$lib/services/sse.js', () => ({
    sseService: {
        on: vi.fn(),
        connect: vi.fn(),
        disconnect: vi.fn()
    }
}));

describe('DeliveryStore', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        deliveryStore.deliveryStatuses = {};
        deliveryStore.failedReminders = [];
        deliveryStore.connectionStatus = 'disconnected';
    });

    it('should initialize with empty state', () => {
        expect(deliveryStore.deliveryStatuses).toEqual({});
        expect(deliveryStore.failedReminders).toEqual([]);
    });
});
```

**Component Testing Pattern:**
```javascript
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, cleanup } from '@testing-library/svelte';
import ContentPreviewPanel from './ContentPreviewPanel.svelte';

vi.mock('svelte-i18n', async () => {
    const { readable } = await import('svelte/store');
    const translations = {
        'content.picker.articles': 'Artikel',
        'common.close': 'Tutup',
    };
    const mockT = (key) => translations[key] || key;
    const tStore = readable(mockT);
    const t = Object.assign(mockT, { subscribe: tStore.subscribe });
    return { t, _: t, locale: readable('id'), ... };
});

describe('ContentPreviewPanel', () => {
    it('should render article preview', async () => {
        const onClose = vi.fn();
        const onAttach = vi.fn();

        render(ContentPreviewPanel, {
            props: { content: mockArticle, isSelected: false, onClose, onAttach }
        });

        await waitFor(() => {
            expect(screen.getByText('Test Article Title')).toBeInTheDocument();
        });
    });
});
```

### Backend (Go)

**Suite Organization:**
```go
import (
    "testing"
)

func TestLoad(t *testing.T) {
    // Create a temporary config file
    content := `server:
  port: 9090`
    // ... test setup

    // Test loading the config
    cfg, err := Load(tmpFile.Name())
    if err != nil {
        t.Fatalf("Failed to load config: %v", err)
    }

    // Verify
    if cfg.Server.Port != 9090 {
        t.Errorf("Expected port 9090, got %d", cfg.Server.Port)
    }
}
```

**Table-Driven Tests Pattern:**
```go
func TestMaskPhone(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"Indonesian format with 62", "628123456789", "628***789"},
        {"Short number", "12345", "***"},
        {"Empty string", "", "***"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := MaskPhone(tt.input)
            if result != tt.expected {
                t.Errorf("MaskPhone(%q) = %q, want %q", tt.input, result, tt.expected)
            }
        })
    }
}
```

---

## Mocking

### Frontend Mocking

**Framework:** Vitest's `vi.mock()` and `vi.fn()`

**Mock Pattern - svelte-i18n:**
```javascript
vi.mock('svelte-i18n', async () => {
    const { readable } = await import('svelte/store');

    const translations = {
        'common.close': 'Tutup',
        'common.loading': 'Memuat...',
    };

    const mockT = (key, options) => translations[key] || key;
    const tStore = readable(mockT);
    const t = Object.assign(mockT, { subscribe: tStore.subscribe });

    return {
        t,
        _: mockT,
        locale: readable('id'),
        locales: readable(['id', 'en']),
        loading: readable(false),
        init: vi.fn(),
        getLocaleFromNavigator: vi.fn(() => 'id'),
        addMessages: vi.fn(),
    };
});
```

**Mock Pattern - Reusable i18n Mock:**
```javascript
import { createI18nMock } from '../../test-utils/i18nMock.js';

vi.mock('svelte-i18n', () => createI18nMock({
    'custom.key': 'Custom Translation',
    'common.close': 'Override Close',
}));
```

**Mock Pattern - Services:**
```javascript
vi.mock('$lib/services/sse.js', () => ({
    sseService: {
        on: vi.fn(),
        connect: vi.fn(),
        disconnect: vi.fn()
    }
}));

vi.mock('./toast.svelte.js', () => ({
    toastStore: {
        add: vi.fn()
    }
}));
```

**Mock Pattern - Components:**
```javascript
import Sidebar from './Sidebar.svelte';

vi.mock("svelte-i18n", async () => {
    // ... mock implementation
});

describe("Sidebar", () => {
    it("should show analytics link for superadmin", async () => {
        const onNavigate = vi.fn();
        const { container } = render(Sidebar, {
            props: {
                user: { role: "superadmin", username: "admin" },
                currentView: "dashboard",
                stats: { totalPatients: 5 },
                users: [],
                onNavigate,
            },
        });
        // ... assertions
    });
});
```

### Backend Mocking

**Approach:** Minimal mocking, test with real implementations where possible

**Temp Files for Config:**
```go
tmpFile, err := os.CreateTemp("", "config-*.yaml")
if err != nil {
    t.Fatalf("Failed to create temp file: %v", err)
}
defer os.Remove(tmpFile.Name())

if _, err := tmpFile.WriteString(content); err != nil {
    t.Fatalf("Failed to write temp file: %v", err)
}
tmpFile.Close()

cfg, err := Load(tmpFile.Name())
if err != nil {
    t.Fatalf("Failed to load config: %v", err)
}
```

**Buffer for Logger Output:**
```go
var buf bytes.Buffer
logger := NewLogger(LoggerConfig{
    Level:  "info",
    Format: "json",
    Output: &buf,
})

logger.Info("test message")
output := buf.String()
```

---

## Fixtures and Factories

### Frontend

**Test Data - Inline:**
```javascript
const mockArticle = {
    id: '1',
    title: 'Test Article Title',
    excerpt: 'This is a test excerpt...',
    publishedAt: '2024-01-15T10:00:00Z',
    heroImages: {
        hero16x9: '/uploads/articles/hero-16x9-1.jpg',
        hero1x1: '/uploads/articles/hero-1x1-1.jpg'
    }
};

const mockVideo = {
    id: '2',
    title: 'Test Video Title',
    YouTubeID: 'dQw4w9WgXcQ',
    thumbnailURL: 'https://img.youtube.com/vi/dQw4w9WgXcQ/maxresdefault.jpg',
    duration: '10:30',
    channelName: 'Test Channel'
};
```

**Default Props Factory:**
```javascript
const defaultProps = {
    patients: mockPatients,
    selectedPatientId: null,
    searchQuery: "",
    onSelect: vi.fn(),
    onAddPatient: vi.fn(),
    onEditPatient: null,
    onDeletePatient: null
};
```

### Backend

**Table-Driven Test Data:**
```go
tests := []struct {
    name     string
    input    string
    expected string
}{
    {"Indonesian format with 62", "628123456789", "628***789"},
    {"Short number", "12345", "***"},
    {"Empty string", "", "***"},
}
```

---

## Setup and Teardown

### Frontend Setup

**Test Environment Setup:**
File: `/media/davidyusaku/Windows/BACKUP/Portfolio/Web/prima_v2/frontend/src/test/setup.js`

```javascript
// Setup happy-dom environment BEFORE any imports
import { Window } from 'happy-dom';

const window = new Window({
    url: 'http://localhost',
    pretendToBeVisual: true
});

// Set up global objects BEFORE any test code runs
globalThis.window = window;
globalThis.document = window.document;
globalThis.navigator = window.navigator;
globalThis.Element = window.Element;
// ... more globals

import { beforeAll, afterAll, vi } from 'vitest';
import '@testing-library/jest-dom';

// Mock ResizeObserver, MutationObserver, matchMedia, etc.
globalThis.ResizeObserver = class ResizeObserver {
    observe() {}
    unobserve() {}
    disconnect() {}
};

// Mock localStorage
const localStorageMock = {
    store: {},
    getItem: vi.fn((key) => this.store[key] || null),
    setItem: vi.fn((key, value) => { this.store[key] = value; }),
    // ...
};

// Suppress console errors during tests
const originalError = console.error;
beforeAll(() => {
    console.error = (...args) => {
        if (
            args[0]?.includes?.('Hydration') ||
            args[0]?.includes?.('was passed to') ||
            // ... more suppressions
        ) {
            return;
        }
        originalError.call(console, ...args);
    };
});
```

**Vitest Config:**
```javascript
export default defineConfig({
    test: {
        globals: true,
        include: ["src/**/*.test.js"],
        setupFiles: ["src/test/setup.js"],
        environment: "happy-dom",
        singleThread: true,
        server: {
            deps: {
                inline: [/svelte/, "@testing-library/svelte"],
            },
        },
    },
});
```

### Backend Setup

**Minimal Setup:**
- Go's `testing` package handles setup automatically
- Use `t.Fatalf()` for fatal failures during setup
- Use `defer` for cleanup (e.g., `defer os.Remove(tmpFile.Name())`)

---

## Coverage

### Frontend

**Coverage Command:**
```bash
# Not configured in package.json, but available via vitest
npx vitest run --coverage
```

**No coverage target enforced currently**

### Backend

**Coverage Command:**
```bash
go test -cover ./...
```

**Output Example:**
```
PASS: coverage: 85.4% of statements
```

---

## Test Types

### Frontend Tests

**Unit Tests:**
- Store logic: `delivery.test.js` - tests state mutations, reactive updates
- Utilities: `i18nMock.test.js` - tests mock creation and behavior

**Integration Tests:**
- Component rendering: `ContentPreviewPanel.test.js`
- Component interactions: `PatientListPane.test.js`
- Full component suites: `ContentPickerModal.*.test.js`

**Test Coverage Areas:**
- Props and defaults
- State updates and reactivity
- Event handlers and callbacks
- Accessibility attributes (ARIA)
- i18n integration
- Conditional rendering
- Error handling

### Backend Tests

**Unit Tests:**
- Config loading and defaults: `config_test.go`
- Utility functions: `logger_test.go`, `mask_test.go`
- Validation logic: `config_test.go` - `TestQuietHoursValidation_*`

**Integration Tests:**
- Handler tests with HTTP: `handlers/health_test.go`
- Service tests with external dependencies mocked

**Test Coverage Areas:**
- Config file parsing
- Environment variable expansion
- Default value application
- Validation rules
- Error handling paths
- Boundary conditions

---

## Common Patterns

### Async Testing

**Frontend (Svelte + Vitest):**
```javascript
it('should update delivery status correctly', () => {
    const reminderId = 'reminder-123';
    const status = 'sent';
    const timestamp = '2025-12-30T10:00:00Z';

    deliveryStore.updateStatus(reminderId, status, timestamp);

    expect(deliveryStore.deliveryStatuses[reminderId]).toMatchObject({
        status,
        timestamp
    });
});
```

**Frontend (Component with async):**
```javascript
it('should render article preview with title, excerpt, and date', async () => {
    const onClose = vi.fn();
    const onAttach = vi.fn();

    render(ContentPreviewPanel, {
        props: { content: mockArticle, isSelected: false, onClose, onAttach }
    });

    await waitFor(() => {
        expect(screen.getByText('Test Article Title')).toBeInTheDocument();
    });
});
```

### Error Testing

**Frontend:**
```javascript
it('should return null for non-existent reminder status', () => {
    const status = deliveryStore.getStatus('non-existent');
    expect(status).toBeNull();
});
```

**Backend:**
```go
func TestQuietHoursValidation_InvalidStartHour(t *testing.T) {
    startHour := 25 // Invalid: > 23
    endHour := 6
    cfg := &QuietHoursConfig{
        StartHour: &startHour,
        EndHour:   &endHour,
        Timezone:  "WIB",
    }

    err := cfg.Validate()
    if err == nil {
        t.Error("Expected error for invalid start_hour, got nil")
    }
}
```

### Svelte 5 Reactivity Testing

**Store Reactivity:**
```javascript
it('should create new object reference when updating status (Svelte 5 reactivity)', () => {
    const originalRef = deliveryStore.deliveryStatuses;

    deliveryStore.updateStatus('reminder-123', 'sent', '2025-12-30T10:00:00Z');

    // Reference should change
    expect(deliveryStore.deliveryStatuses).not.toBe(originalRef);
});

it('should create new array reference when adding failed reminder (Svelte 5 reactivity)', () => {
    const originalRef = deliveryStore.failedReminders;

    deliveryStore.addFailedReminder({ /* ... */ });

    // Reference should change
    expect(deliveryStore.failedReminders).not.toBe(originalRef);
});
```

### Accessibility Testing

```javascript
describe('Accessibility', () => {
    it('should have role="status" for screen readers', async () => {
        render(DeliveryStatusBadge, { status: 'delivered' });
        expect(screen.getByRole('status')).toBeInTheDocument();
    });

    it('should have aria-label with status text', async () => {
        render(DeliveryStatusBadge, { status: 'delivered' });
        const status = screen.getByRole('status');
        expect(status).toHaveAttribute('aria-label', 'Diterima');
    });

    it('should have aria-live="polite" for real-time updates', async () => {
        render(DeliveryStatusBadge, { status: 'delivered' });
        const status = screen.getByRole('status');
        expect(status).toHaveAttribute('aria-live', 'polite');
    });
});
```

### Mock Cleanup

```javascript
afterEach(() => {
    vi.restoreAllMocks();
    cleanup();
});
```

---

*Testing analysis: 2025-01-17*
