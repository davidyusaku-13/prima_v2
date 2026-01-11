# AI Context Summary - PRIMA

**Optimized context for AI coding assistants (<4K tokens)**

---

## Project Identity

**PRIMA** = Healthcare Volunteer Dashboard
- **Backend:** Go 1.25 + Gin (port 8080)
- **Frontend:** Svelte 5 + Vite (port 5173) - **NOT SvelteKit**
- **Database:** JSON files with RWMutex
- **Styling:** Tailwind CSS 4
- **i18n:** svelte-i18n (EN/ID)

---

## Critical Constraint

**NEVER use SvelteKit imports. They will break the build.**

| Prohibited | Use Instead |
|------------|-------------|
| `import { goto } from '$app/navigation'` | `window.location.href = '...'` |
| `import { page } from '$app/stores'` | Props or context |
| `import { browser } from '$app/environment'` | `typeof window !== 'undefined'` |

---

## Frontend Patterns

### Svelte 5 Runes

```svelte
<script>
  // State
  let count = $state(0);

  // Derived
  let doubled = $derived(count * 2);

  // Effect
  $effect(() => {
    console.log('Count:', count);
  });

  // Props
  let { title = 'Default', onClose } = $props();
</script>
```

### Event Handlers

```svelte
<!-- Svelte 5 syntax -->
<button onclick={handleClick}>Click</button>
<input oninput={(e) => value = e.target.value} />
```

### Store Pattern

```javascript
// $lib/stores/auth.svelte.js
class AuthStore {
  user = $state(null);
  token = $state(null);

  get isAuthenticated() {
    return this.token !== null;
  }

  login(token, user) {
    this.token = token;
    this.user = user;
    localStorage.setItem('token', token);
  }

  logout() {
    this.token = null;
    this.user = null;
    localStorage.removeItem('token');
  }
}

export const auth = new AuthStore();
```

### Import Paths

```javascript
// Use $lib alias
import { auth } from '$lib/stores/auth.svelte.js';
import PatientModal from '$lib/components/PatientModal.svelte';
import { getPatients } from '$lib/api/patients.js';
```

---

## Backend Patterns

### Handler Structure

```go
func (h *Handler) CreatePatient(c *gin.Context) {
    var req CreatePatientRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": "VALIDATION_ERROR"})
        return
    }

    patient := &models.Patient{
        ID:        uuid.New().String(),
        Name:      req.Name,
        Phone:     utils.NormalizePhone(req.Phone),
        CreatedBy: c.GetString("userId"),
    }

    h.store.Lock()
    h.store.Patients[patient.ID] = patient
    h.store.SaveData()
    h.store.Unlock()

    c.JSON(201, patient)
}
```

### Thread-Safe Store Access

```go
// Read
store.RLock()
patient, exists := store.Patients[id]
store.RUnlock()

// Write
store.Lock()
store.Patients[id] = patient
store.SaveData()
store.Unlock()
```

---

## API Quick Reference

### Auth

```
POST /api/auth/login    { username, password } → { token, user }
POST /api/auth/register { username, password, fullName } → { token, user }
GET  /api/auth/me       → { user }
```

### Patients

```
GET    /api/patients           → { patients[] }
POST   /api/patients           { name, phone } → patient
GET    /api/patients/:id       → patient
PUT    /api/patients/:id       { name?, phone? } → patient
DELETE /api/patients/:id       → { message }
```

### Reminders

```
GET  /api/patients/:id/reminders     → { data[], pagination }
POST /api/patients/:id/reminders     { title, dueDate } → reminder
POST /api/patients/:id/reminders/:rid/send  → { status }
POST /api/reminders/:id/cancel       → { status }
POST /api/reminders/:id/retry        → { status }
```

### Content

```
GET  /api/articles        → { articles[] }
POST /api/articles        { title, content } → article
GET  /api/videos          → { videos[] }
POST /api/videos          { youtube_url } → video
```

---

## Key Files

| Purpose | Path |
|---------|------|
| App entry | `frontend/src/App.svelte` |
| Main backend | `backend/main.go` |
| Auth store | `frontend/src/lib/stores/auth.svelte.js` |
| API client | `frontend/src/lib/utils/api.js` |
| Patient model | `backend/models/patient.go` |
| GOWA client | `backend/services/gowa.go` |
| Config | `backend/config.yaml` |

---

## Common Components

| Component | Props | Purpose |
|-----------|-------|---------|
| `PatientModal` | show, editingPatient, onClose, onSave | Add/edit patient |
| `ReminderModal` | show, patient, reminder, onClose, onSave | Add/edit reminder |
| `ConfirmModal` | show, title, message, onConfirm, onCancel | Confirmation dialog |
| `Toast` | message, type (success/error/info) | Notifications |
| `DeliveryStatusBadge` | status | Show delivery state |

---

## Roles

| Role | Permissions |
|------|-------------|
| `volunteer` | Own patients only |
| `admin` | All patients + CMS |
| `superadmin` | All + user management |

---

## Error Codes

| Code | HTTP | Meaning |
|------|------|---------|
| `VALIDATION_ERROR` | 400 | Bad input |
| `INVALID_PHONE` | 400 | Wrong phone format |
| `UNAUTHORIZED` | 401 | Bad/missing token |
| `FORBIDDEN` | 403 | Wrong role |
| `NOT_FOUND` | 404 | Resource missing |
| `GOWA_UNAVAILABLE` | 503 | WhatsApp down |

---

## Quick Commands

```bash
# Backend
cd backend && go run main.go
cd backend && go test ./...

# Frontend
cd frontend && bun run dev
cd frontend && bun run test
cd frontend && bun run build
```

---

**Default login:** superadmin / superadmin
