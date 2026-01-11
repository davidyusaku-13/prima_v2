# Component Inventory - Frontend

**Generated:** January 11, 2026 (Updated)
**Project:** PRIMA Healthcare Volunteer Dashboard
**Total Components:** 52 Svelte 5 components
**Scan Type:** Exhaustive Rescan

---

## Component Organization

Components are organized by feature/functionality in `src/lib/components/`:

```
components/
├── Core UI Components (modals, navigation)
├── Feature Components (patients, reminders, content)
├── Analytics Components
├── Indicators (badges, status)
├── UI Primitives (toast, buttons)
└── WhatsApp Integration
```

---

## Core UI Components

### Navigation

#### Sidebar.svelte

**Location:** `src/lib/components/Sidebar.svelte`

**Purpose:** Desktop sidebar navigation

**Props:**

```javascript
let { currentView, onNavigate } = $props();
```

**Events:**

- `onNavigate(view)` - Called when navigation item clicked

**Usage:**

```svelte
<Sidebar currentView={currentView} onNavigate={setView} />
```

**Features:**

- Role-based menu items (admin vs volunteer)
- Active route highlighting
- Icons for each section
- Responsive collapse

---

#### BottomNav.svelte

**Location:** `src/lib/components/BottomNav.svelte`

**Purpose:** Mobile bottom navigation

**Props:**

```javascript
let { currentView, onNavigate } = $props();
```

**Events:**

- `onNavigate(view)` - Called when tab clicked

**Usage:**

```svelte
<BottomNav currentView={currentView} onNavigate={setView} />
```

**Features:**

- 4-5 main tabs
- Icon + label
- Active state
- Fixed to bottom on mobile

---

### Modals

#### PatientModal.svelte

**Location:** `src/lib/components/PatientModal.svelte`

**Purpose:** Create/edit patient

**Props:**

```javascript
let {
  show = false,
  editingPatient = null,
  patientForm = { name: "", phone: "", email: "", notes: "" },
  onClose = () => {},
  onSave = () => {},
} = $props();
```

**Events:**

- `onClose()` - Cancel/close modal
- `onSave()` - Save patient

**Usage:**

```svelte
<PatientModal
  show={showModal}
  editingPatient={patient}
  patientForm={form}
  onClose={closeModal}
  onSave={savePatient}
/>
```

**Features:**

- Create vs edit mode
- Phone number formatting
- Email validation
- Notes textarea

---

#### ReminderModal.svelte

**Location:** `src/lib/components/ReminderModal.svelte`

**Purpose:** Create/edit reminder with content attachments

**Props:**

```javascript
let {
  show = false,
  editingReminder = null,
  reminderForm = {
    patientId: "",
    title: "",
    description: "",
    dueDate: "",
    priority: "medium",
    recurrence: { frequency: "none", interval: 1, daysOfWeek: [], endDate: "" },
    attachments: [],
  },
  onClose = () => {},
  onSave = () => {},
} = $props();
```

**Events:**

- `onClose()` - Cancel/close modal
- `onSave()` - Save reminder

**Features:**

- Title, description, due date
- Priority (low, medium, high)
- Recurrence settings (daily, weekly, monthly)
- Content attachments (articles/videos)
- Days of week picker for weekly recurrence
- End date for recurrence

---

#### UserModal.svelte

**Location:** `src/lib/components/UserModal.svelte`

**Purpose:** Create volunteer or edit user role

**Props:**

```javascript
let {
  show = false,
  editingUser = null,
  userForm = { role: "volunteer", username: "", password: "" },
  loading = false,
  onClose = () => {},
  onSaveRole = () => {},
  onRegister = () => {},
} = $props();
```

**Events:**

- `onClose()` - Cancel/close
- `onSaveRole()` - Update user role (edit mode)
- `onRegister()` - Register new user (create mode)

**Features:**

- Create mode: username, password, role
- Edit mode: change role only
- Password validation (6+ chars)
- Role selector (volunteer, admin)

---

#### ProfileModal.svelte

**Location:** `src/lib/components/ProfileModal.svelte`

**Purpose:** View user profile, change language, logout

**Props:**

```javascript
let {
  show = false,
  user = null,
  locale = "en",
  onSetLocale = () => {},
  onLogout = () => {},
  onClose = () => {},
} = $props();
```

**Events:**

- `onSetLocale(locale)` - Change language (en/id)
- `onLogout()` - Logout user
- `onClose()` - Close modal

**Features:**

- User info display (username, role)
- Language switcher (English/Indonesian)
- Logout button

---

#### ConfirmModal.svelte

**Location:** `src/lib/components/ConfirmModal.svelte`

**Purpose:** Generic confirmation dialog

**Props:**

```javascript
let {
  show = false,
  message = "",
  onClose = () => {},
  onConfirm = () => {},
} = $props();
```

**Events:**

- `onClose()` - Cancel
- `onConfirm()` - Confirm action

**Usage:**

```svelte
<ConfirmModal
  show={showConfirm}
  message="Delete this patient?"
  onClose={() => showConfirm = false}
  onConfirm={deletePatient}
/>
```

**Features:**

- Customizable message
- Cancel + Confirm buttons
- Backdrop click to close

---

#### SendReminderModal.svelte

**Location:** `src/lib/components/SendReminderModal.svelte`

**Purpose:** Preview and send reminder message

**Props:**

```javascript
let {
  show = false,
  reminder = null,
  patient = null,
  onClose = () => {},
  onSend = () => {},
} = $props();
```

**Events:**

- `onClose()` - Cancel
- `onSend()` - Send reminder now

**Features:**

- Preview formatted WhatsApp message
- Show patient info
- Show reminder details
- Confirm send button

---

#### PhoneEditModal.svelte

**Location:** `src/lib/components/PhoneEditModal.svelte`

**Purpose:** Edit patient phone number

**Props:**

```javascript
let {
  show = false,
  currentPhone = "",
  onClose = () => {},
  onSave = (newPhone) => {},
} = $props();
```

**Events:**

- `onClose()` - Cancel
- `onSave(newPhone)` - Save new phone number

**Features:**

- Phone format validation (62812...)
- Real-time validation feedback
- Save disabled if invalid

---

#### VideoModal.svelte

**Location:** `src/lib/components/VideoModal.svelte`

**Purpose:** View video details and YouTube embed

**Props:**

```javascript
let { show = false, video = null, onClose = () => {} } = $props();
```

**Events:**

- `onClose()` - Close modal

**Features:**

- YouTube video embed
- Video title, description
- Category badge
- View count
- Created date

---

#### VideoEditModal.svelte

**Location:** `src/lib/components/VideoEditModal.svelte`

**Purpose:** Edit video title, category, status

**Props:**

```javascript
let {
  video = null,
  onClose = () => {},
  onSave = () => {},
  token = null,
} = $props();
```

**Events:**

- `onClose()` - Cancel
- `onSave()` - Save changes

**Features:**

- Edit title
- Select category
- Toggle status (published/draft)
- Cannot edit YouTube URL (immutable)

---

## Feature Components

### Dashboard

#### DashboardStats.svelte

**Location:** `src/lib/components/DashboardStats.svelte`

**Purpose:** Display key metrics

**Props:**

```javascript
let {
  stats = {
    patients: 0,
    reminders: 0,
    volunteers: 0,
    articles: 0,
    videos: 0,
  },
} = $props();
```

**Features:**

- Grid layout (3 columns)
- Icon for each metric
- Number display
- Color coding

---

#### ActivityLog.svelte

**Location:** `src/lib/components/ActivityLog.svelte`

**Purpose:** Display recent activity

**Props:**

```javascript
let { activities = [] } = $props();
```

**Activity Schema:**

```javascript
{
  id: string,
  type: 'reminder_sent' | 'patient_created' | 'article_published',
  message: string,
  timestamp: string,
  user: string
}
```

**Features:**

- Reverse chronological order
- Activity type icons
- Relative timestamps
- User attribution

---

### Content (CMS)

#### ArticleCard.svelte

**Location:** `src/lib/components/ArticleCard.svelte`

**Purpose:** Display article preview

**Props:**

```javascript
let { article, onClick = () => {} } = $props();
```

**Events:**

- `onclick()` - Navigate to article detail

**Features:**

- Hero image (16:9 aspect ratio)
- Title, excerpt
- Category badge
- View count, created date
- Published/draft status

---

#### VideoCard.svelte

**Location:** `src/lib/components/VideoCard.svelte`

**Purpose:** Display video preview

**Props:**

```javascript
let {
  video,
  onClick = () => {},
  showActions = false,
  onDelete = () => {},
  onEdit = () => {},
} = $props();
```

**Events:**

- `onclick()` - Open video modal
- `onDelete()` - Delete video (admin)
- `onEdit()` - Edit video (admin)

**Features:**

- YouTube thumbnail
- Play icon overlay
- Title, description snippet
- Category badge
- View count
- Action buttons (admin only)

---

#### ContentListItem.svelte

**Location:** `src/lib/components/ContentListItem.svelte`

**Purpose:** Generic content list item (articles or videos)

**Props:**

```javascript
let {
  item,
  type, // 'article' or 'video'
  selected = false,
  onSelect = () => {},
  onEdit = () => {},
  onDelete = () => {},
  onToggleStatus = () => {},
} = $props();
```

**Events:**

- `onSelect()` - Select item for attachment
- `onEdit()` - Edit item
- `onDelete()` - Delete item
- `onToggleStatus()` - Toggle published/draft

**Features:**

- Thumbnail preview
- Title, excerpt
- Status badge
- Selection checkbox (for attachments)
- Action buttons

---

#### ImageUploader.svelte

**Location:** `src/lib/components/ImageUploader.svelte`

**Purpose:** Upload images for articles

**Props:**

```javascript
let { imageUrl = "", required = false, label = "", token = null } = $props();
```

**Events:**

- Updates `imageUrl` prop on upload success

**Features:**

- Drag & drop
- File picker
- Image preview
- Upload progress
- Generates 3 aspect ratios (16:9, 4:3, 1:1)

---

#### QuillEditor.svelte

**Location:** `src/lib/components/QuillEditor.svelte`

**Purpose:** Rich text editor for articles

**Props:**

```javascript
let { content = "" } = $props();
```

**Features:**

- Bold, italic, underline
- Headings, lists
- Links, images
- Code blocks
- Undo/redo
- HTML output

---

### Indicators

#### FailedReminderBadge.svelte

**Location:** `src/lib/components/indicators/FailedReminderBadge.svelte`

**Purpose:** Show count of failed reminders

**Props:**

```javascript
let { count = 0 } = $props();
```

**Features:**

- Red badge with count
- Pulsing animation when count > 0
- Click to view failed deliveries
- Hidden if count === 0

---

## UI Primitives

### Toast.svelte

**Location:** `src/lib/components/ui/Toast.svelte`

**Purpose:** Notification toast

**Props:**

```javascript
let {
  toast = {
    id: number,
    message: string,
    type: "info" | "success" | "error" | "warning",
    action: { label: string, onClick: () => {} } | null,
  },
  onClose = () => {},
} = $props();
```

**Events:**

- `onClose()` - Dismiss toast

**Features:**

- Auto-dismiss after duration
- Manual dismiss button
- Optional action button
- Color coded by type
- Slide-in animation
- Stack multiple toasts

---

## Specialized Components

### Analytics Components

Located in `src/lib/components/analytics/`:

- **DeliveryChart.svelte** - Charts for delivery statistics
- **FailedDeliveryTable.svelte** - Table of failed deliveries
- **AnalyticsFilters.svelte** - Date range, status filters

### Delivery Components

Located in `src/lib/components/delivery/`:

- **DeliveryStatusBadge.svelte** - Color-coded status badge
- **DeliveryTimeline.svelte** - Delivery status timeline

### Health Components

Located in `src/lib/components/health/`:

- **HealthIndicator.svelte** - System health status
- **CircuitBreakerStatus.svelte** - Circuit breaker state

### Patient Components

Located in `src/lib/components/patients/`:

- **PatientCard.svelte** - Patient list item
- **PatientDetailHeader.svelte** - Patient info header

### Reminder Components

Located in `src/lib/components/reminders/`:

- **ReminderCard.svelte** - Reminder list item
- **ReminderList.svelte** - List of reminders
- **RecurrenceEditor.svelte** - Recurrence settings UI

### WhatsApp Components

Located in `src/lib/components/whatsapp/`:

- **MessagePreview.svelte** - Preview WhatsApp message format
- **WhatsAppLink.svelte** - Clickable WhatsApp link

---

## Component Usage Patterns

### Modal Pattern

**All modals follow this structure:**

```svelte
<script>
  let { show = false, onClose, data } = $props();
</script>

{#if show}
  <div class="modal-backdrop" onclick={onClose}>
    <div class="modal-content" onclick={(e) => e.stopPropagation()}>
      <!-- Modal content -->
      <button onclick={onClose}>Close</button>
    </div>
  </div>
{/if}

<style>
  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 50;
  }

  .modal-content {
    background: white;
    border-radius: 0.5rem;
    padding: 2rem;
    max-width: 90%;
    max-height: 90vh;
    overflow-y: auto;
  }
</style>
```

### Card Pattern

**Reusable card components:**

```svelte
<script>
  let { item, onClick } = $props();
</script>

<div class="card" onclick={onClick}>
  {#if item.image}
    <img src={item.image} alt={item.title} class="card-image" />
  {/if}

  <div class="card-body">
    <h3 class="card-title">{item.title}</h3>
    <p class="card-description">{item.description}</p>
  </div>

  <div class="card-footer">
    <span class="badge">{item.category}</span>
    <span class="text-sm text-gray-500">{item.views} views</span>
  </div>
</div>

<style>
  .card {
    background: white;
    border-radius: 0.5rem;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s;
  }

  .card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
</style>
```

### List Pattern

**Rendering lists with keys:**

```svelte
<script>
  let { items } = $props();
</script>

<div class="list">
  {#each items as item (item.id)}
    <ItemCard {item} />
  {/each}

  {#if items.length === 0}
    <p class="empty-state">No items found</p>
  {/if}
</div>
```

---

## Component Template

**When creating new components, use this template:**

```svelte
<script>
  // Imports
  import { onMount } from 'svelte';
  import * as api from '$lib/utils/api.js';

  // Props
  let {
    data,
    onAction = () => {}
  } = $props();

  // State
  let loading = $state(false);
  let error = $state(null);

  // Derived state
  let isEmpty = $derived(data.length === 0);

  // Functions
  async function handleAction() {
    loading = true;
    error = null;
    try {
      await onAction();
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  // Effects
  $effect(() => {
    console.log('Data changed:', data);
  });
</script>

<!-- Template -->
<div class="component-name">
  {#if loading}
    <p>Loading...</p>
  {:else if error}
    <p class="error">{error}</p>
  {:else if isEmpty}
    <p class="empty">No data</p>
  {:else}
    <!-- Main content -->
    <button onclick={handleAction}>Action</button>
  {/if}
</div>

<!-- Scoped styles -->
<style>
  .component-name {
    /* Component styles */
  }

  .error {
    color: red;
  }

  .empty {
    color: gray;
  }
</style>
```

---

## Component Guidelines

### Props

- Use `$props()` rune for all props
- Provide default values where sensible
- Document expected types in comments

### Events

- Use callback props (e.g., `onClose`, `onSave`)
- Don't use `createEventDispatcher` (Svelte 4 pattern)
- Prefix callbacks with `on`

### State

- Use `$state()` for reactive local state
- Use `$derived()` for computed values
- Use `$effect()` for side effects

### Styling

- Use Tailwind utility classes primarily
- Use scoped `<style>` for component-specific styles
- Avoid global styles except in `app.css`

### Accessibility

- Add `aria-label` to icon-only buttons
- Use semantic HTML (`<button>`, `<nav>`, `<main>`)
- Test keyboard navigation
- Ensure sufficient color contrast

### Testing

- Create `.test.js` file alongside component
- Test props, events, and state changes
- Use `@testing-library/svelte`

---

**Next:** See [Architecture - Frontend](./architecture-frontend.md) for architectural patterns and [Development Guide - Frontend](./development-guide-frontend.md) for development workflow.
