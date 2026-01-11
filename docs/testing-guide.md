# Testing Guide - PRIMA

**Comprehensive testing documentation for backend and frontend**

---

## Overview

PRIMA uses:
- **Backend:** Go standard testing with table-driven tests
- **Frontend:** Vitest with Testing Library for Svelte

---

## Backend Testing (Go)

### Running Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test -v ./handlers
go test -v ./services
go test -v ./utils

# Run specific test by name
go test -v -run TestSendReminder ./handlers

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test File Structure

```
backend/
├── config/
│   └── config_test.go      # Configuration loading tests
├── handlers/
│   ├── analytics_test.go   # Analytics endpoint tests
│   ├── content_test.go     # CMS endpoint tests
│   ├── health_test.go      # Health check tests
│   ├── reminder_test.go    # Reminder endpoint tests
│   ├── sse_test.go         # Server-Sent Events tests
│   └── webhook_test.go     # GOWA webhook tests
├── services/
│   ├── gowa_test.go        # GOWA client tests
│   └── scheduler_test.go   # Reminder scheduler tests
└── utils/
    ├── logger_test.go      # Logging tests
    ├── mask_test.go        # Data masking tests
    ├── phone_test.go       # Phone validation tests
    └── quiethours_test.go  # Quiet hours tests
```

### Table-Driven Test Pattern

```go
func TestNormalizePhone(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid 08", "081234567890", "6281234567890", false},
        {"valid +62", "+6281234567890", "6281234567890", false},
        {"valid 62", "6281234567890", "6281234567890", false},
        {"too short", "0812345", "", true},
        {"invalid prefix", "091234567890", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := NormalizePhone(tt.input)
            if tt.wantErr {
                if err == nil {
                    t.Errorf("expected error, got nil")
                }
                return
            }
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

### Mocking HTTP Clients

For GOWA integration tests:

```go
type mockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}

func TestSendMessage_Success(t *testing.T) {
    client := &mockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            return &http.Response{
                StatusCode: 200,
                Body: io.NopCloser(strings.NewReader(`{"message_id":"test-123"}`)),
            }, nil
        },
    }

    gowaClient := NewGOWAClient(client, config)
    result, err := gowaClient.SendMessage("628123456789", "Test message")

    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result.MessageID != "test-123" {
        t.Errorf("got %s, want test-123", result.MessageID)
    }
}
```

### Test Data Management

```go
// Use t.TempDir() for temporary files
func TestSavePatients(t *testing.T) {
    tmpDir := t.TempDir()
    dataFile := filepath.Join(tmpDir, "patients.json")

    store := NewPatientStore(dataFile)
    store.Add(&Patient{ID: "test-1", Name: "Test"})
    store.Save()

    // Verify file contents
    data, _ := os.ReadFile(dataFile)
    // assertions...
}
```

### Coverage Targets

| Package | Target | Notes |
|---------|--------|-------|
| `handlers/` | 80%+ | Critical business logic |
| `services/` | 80%+ | GOWA, scheduler |
| `utils/` | 90%+ | Pure functions |
| `config/` | 70%+ | Configuration loading |

---

## Frontend Testing (Svelte 5)

### Running Tests

```bash
cd frontend

# Run all tests (watch mode)
bun run test

# Run once (CI mode)
bun run test -- --run

# Run specific file
bun run test src/lib/components/PatientModal.test.js

# Run with coverage
bun run test -- --coverage
```

### Test File Structure

```
frontend/
├── src/
│   ├── test/
│   │   ├── api.test.js         # API client tests
│   │   ├── components.test.js  # Component tests
│   │   └── utils.test.js       # Utility function tests
│   └── lib/
│       └── components/
│           └── __tests__/       # Co-located component tests
└── vitest.config.js
```

### Component Testing Pattern

```javascript
import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import PatientModal from './PatientModal.svelte';

describe('PatientModal', () => {
  it('renders when show is true', () => {
    render(PatientModal, { props: { show: true } });
    expect(screen.getByText('Add Patient')).toBeInTheDocument();
  });

  it('calls onSave when form is submitted', async () => {
    const onSave = vi.fn();
    render(PatientModal, {
      props: { show: true, onSave }
    });

    await fireEvent.input(screen.getByLabelText('Name'), {
      target: { value: 'Test Patient' }
    });
    await fireEvent.input(screen.getByLabelText('Phone'), {
      target: { value: '081234567890' }
    });
    await fireEvent.click(screen.getByText('Save'));

    expect(onSave).toHaveBeenCalledWith({
      name: 'Test Patient',
      phone: '081234567890'
    });
  });

  it('closes when cancel is clicked', async () => {
    const onClose = vi.fn();
    render(PatientModal, {
      props: { show: true, onClose }
    });

    await fireEvent.click(screen.getByText('Cancel'));
    expect(onClose).toHaveBeenCalled();
  });
});
```

### Testing Svelte 5 Stores

```javascript
import { describe, it, expect, beforeEach } from 'vitest';
import { auth } from '$lib/stores/auth.svelte.js';

describe('auth store', () => {
  beforeEach(() => {
    auth.logout();
  });

  it('starts logged out', () => {
    expect(auth.isAuthenticated).toBe(false);
    expect(auth.user).toBe(null);
  });

  it('stores user on login', () => {
    auth.login('test-token', { id: '1', username: 'test', role: 'volunteer' });

    expect(auth.isAuthenticated).toBe(true);
    expect(auth.user.username).toBe('test');
    expect(auth.token).toBe('test-token');
  });

  it('clears state on logout', () => {
    auth.login('test-token', { id: '1', username: 'test', role: 'volunteer' });
    auth.logout();

    expect(auth.isAuthenticated).toBe(false);
    expect(auth.user).toBe(null);
  });
});
```

### Mocking API Calls

```javascript
import { vi } from 'vitest';
import * as api from '$lib/utils/api.js';

vi.mock('$lib/utils/api.js', () => ({
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn()
}));

describe('PatientList', () => {
  it('fetches patients on mount', async () => {
    api.get.mockResolvedValue({
      patients: [{ id: '1', name: 'Test Patient' }]
    });

    render(PatientList);

    await waitFor(() => {
      expect(screen.getByText('Test Patient')).toBeInTheDocument();
    });
    expect(api.get).toHaveBeenCalledWith('/api/patients');
  });
});
```

### Testing i18n

```javascript
import { render, screen } from '@testing-library/svelte';
import { init, locale } from 'svelte-i18n';
import Dashboard from './Dashboard.svelte';

beforeAll(async () => {
  init({
    fallbackLocale: 'en',
    initialLocale: 'en'
  });
  await locale.set('en');
});

it('displays translated text', () => {
  render(Dashboard);
  expect(screen.getByText('Dashboard')).toBeInTheDocument();
});
```

---

## Integration Testing

### End-to-End Flow Testing

For critical user flows, test the full stack:

```bash
# Terminal 1: Start backend
cd backend && go run main.go

# Terminal 2: Start frontend
cd frontend && bun run dev

# Terminal 3: Run integration tests
cd frontend && bun run test:e2e
```

### SSE Testing

```javascript
describe('SSE delivery updates', () => {
  it('receives status updates', async () => {
    const events = [];
    const eventSource = new EventSource('/api/sse/delivery-status?token=test');

    eventSource.addEventListener('delivery.status.updated', (e) => {
      events.push(JSON.parse(e.data));
    });

    // Trigger a delivery status change via API
    await api.post('/api/patients/1/reminders/1/send');

    // Wait for SSE event
    await waitFor(() => {
      expect(events.length).toBeGreaterThan(0);
    });

    eventSource.close();
  });
});
```

### Webhook Testing

```go
func TestGOWAWebhook(t *testing.T) {
    // Create test payload
    payload := `{"event":"message.ack","message":{"id":"test-123","status":"delivered"}}`

    // Calculate HMAC signature
    mac := hmac.New(sha256.New, []byte(webhookSecret))
    mac.Write([]byte(payload))
    signature := hex.EncodeToString(mac.Sum(nil))

    // Create request
    req := httptest.NewRequest("POST", "/api/webhook/gowa", strings.NewReader(payload))
    req.Header.Set("X-Webhook-Signature", signature)
    req.Header.Set("Content-Type", "application/json")

    // Execute
    w := httptest.NewRecorder()
    handler.HandleWebhook(w, req)

    // Assert
    if w.Code != http.StatusOK {
        t.Errorf("got status %d, want 200", w.Code)
    }
}
```

---

## Continuous Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - run: cd backend && go test -v -cover ./...

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: oven-sh/setup-bun@v1
      - run: cd frontend && bun install
      - run: cd frontend && bun run test -- --run
```

---

## Best Practices

### Do

- Write tests alongside new features
- Use table-driven tests for multiple inputs
- Mock external dependencies (GOWA, YouTube API)
- Test error paths, not just happy paths
- Use meaningful test names that describe behavior

### Don't

- Don't test implementation details
- Don't rely on test order
- Don't use production data in tests
- Don't skip tests without explanation
- Don't write flaky tests (use proper async handling)

---

**Last Updated:** January 11, 2026
