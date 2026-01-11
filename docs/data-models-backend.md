# Data Models - Backend

**Generated:** January 11, 2026 (Updated)
**Project:** PRIMA Healthcare Volunteer Dashboard
**Persistence:** JSON files with `sync.RWMutex` thread-safety
**Scan Type:** Exhaustive Rescan

---

## Overview

PRIMA uses **JSON file-based persistence** with Go structs mapped to domain entities. All data operations are thread-safe using `sync.RWMutex` for concurrent read/write access.

**Storage Location:** `backend/data/`

---

## Core Entities

### 1. Patient

Represents a patient under volunteer care.

**File:** `backend/data/patients.json`  
**Model:** `models/patient.go`

```go
type Patient struct {
    ID        string      `json:"id"`        // UUID
    Name      string      `json:"name"`      // Full name
    Phone     string      `json:"phone"`     // Normalized: 628xxxxxxxxx
    Email     string      `json:"email,omitempty"`
    Notes     string      `json:"notes,omitempty"`
    Reminders []*Reminder `json:"reminders,omitempty"` // Nested
    CreatedBy string      `json:"createdBy,omitempty"` // User ID (volunteer)
}
```

**Relationships:**

- **1:N** with `Reminder` (embedded in patient JSON)
- **N:1** with `User` (via `CreatedBy`)

**Access Control:**

- Volunteers: Own patients only
- Admin/Superadmin: All patients

**JSON Storage:**

```json
{
  "patient-uuid-1": {
    "id": "patient-uuid-1",
    "name": "John Doe",
    "phone": "628123456789",
    "email": "john@example.com",
    "notes": "Diabetes treatment",
    "reminders": [...],
    "createdBy": "volunteer-user-id"
  }
}
```

---

### 2. Reminder

Represents a scheduled reminder with WhatsApp delivery tracking.

**Model:** `models/patient.go` (nested in Patient)

```go
type Reminder struct {
    ID          string       `json:"id"`          // UUID
    Title       string       `json:"title"`       // Display title
    Description string       `json:"description"` // Detail text
    DueDate     string       `json:"dueDate,omitempty"` // ISO 8601
    Priority    string       `json:"priority"`    // low|medium|high
    Completed   bool         `json:"completed"`
    Recurrence  Recurrence   `json:"recurrence"`
    Notified    bool         `json:"notified"`
    Attachments []Attachment `json:"attachments,omitempty"` // Max 3

    // Delivery tracking
    GOWAMessageID        string `json:"gowa_message_id,omitempty"`
    DeliveryStatus       string `json:"delivery_status,omitempty"`
    DeliveryErrorMessage string `json:"delivery_error_message,omitempty"`
    MessageSentAt        string `json:"message_sent_at,omitempty"`
    DeliveredAt          string `json:"delivered_at,omitempty"`
    ReadAt               string `json:"read_at,omitempty"`
    RetryCount           int    `json:"retry_count,omitempty"`
    ScheduledDeliveryAt  string `json:"scheduled_delivery_at,omitempty"`
    CancelledAt          string `json:"cancelled_at,omitempty"`
    CancelledBy          string `json:"cancelled_by,omitempty"`
}
```

**Delivery Status State Machine:**

```
pending → scheduled (quiet hours) → sending → sent → delivered → read
pending → queued (circuit breaker) → sending → sent → delivered → read
sending → retrying → sending (transient failure)
retrying → sent (success) | failed (max retries)
any → cancelled (user action)
```

**Status Values:**

- `pending` - Created, not yet sent
- `scheduled` - Queued for quiet hours (22:00-06:00 WIB)
- `queued` - Circuit breaker open, waiting for GOWA
- `sending` - Currently being sent to GOWA
- `retrying` - Failed, waiting for retry
- `sent` - Sent to GOWA successfully
- `delivered` - Delivered to WhatsApp (webhook)
- `read` - Read by recipient (webhook)
- `failed` - All retry attempts exhausted
- `expired` - Not sent within reasonable time
- `cancelled` - User cancelled

---

### 3. Recurrence

Reminder repetition schedule.

```go
type Recurrence struct {
    Frequency  string `json:"frequency"`  // none|daily|weekly|monthly
    Interval   int    `json:"interval"`   // Repeat every N frequency units
    DaysOfWeek []int  `json:"daysOfWeek"` // For weekly: [0-6] (Sun-Sat)
    EndDate    string `json:"endDate,omitempty"` // ISO 8601
}
```

**Examples:**

- Daily: `{frequency: "daily", interval: 1}`
- Every 2 days: `{frequency: "daily", interval: 2}`
- Weekly on Mon/Wed/Fri: `{frequency: "weekly", interval: 1, daysOfWeek: [1,3,5]}`
- Monthly: `{frequency: "monthly", interval: 1}`

---

### 4. Attachment

Content attached to reminder (Article or Video).

```go
type Attachment struct {
    Type  string `json:"type"`  // "article" | "video"
    ID    string `json:"id"`    // Content ID (validated against CMS)
    Title string `json:"title"` // Cached title for WhatsApp message
    URL   string `json:"url"`   // Article slug or video watch URL
}
```

**Validation:**

- Maximum 3 attachments per reminder
- Content ID must exist in CMS (articles.json or videos.json)
- Increments `attachment_count` in content entity

---

### 5. User

User account (volunteer, admin, superadmin).

**File:** `backend/data/users.json`  
**Model:** `main.go` (User struct)

```go
type User struct {
    ID        string `json:"id"`         // UUID
    Username  string `json:"username"`   // Unique, 3-30 chars
    FullName  string `json:"fullName,omitempty"`
    Password  string `json:"password"`   // SHA256 hash (Base64)
    Role      Role   `json:"role"`       // superadmin|admin|volunteer
    CreatedAt string `json:"createdAt"`  // ISO 8601
}
```

**Roles:**

- `superadmin` - Full system access, user management
- `admin` - CMS management, analytics access
- `volunteer` - Patient management (own patients only)

**Password Hashing:**

```go
hash := sha256.Sum256([]byte(password))
encoded := base64.StdEncoding.EncodeToString(hash[:])
```

**JSON Storage:**

```json
{
  "user-uuid": {
    "id": "user-uuid",
    "username": "volunteer1",
    "fullName": "Jane Smith",
    "password": "hashed-password-base64",
    "role": "volunteer",
    "createdAt": "2026-01-01T00:00:00Z"
  }
}
```

---

### 6. Category

Content category for articles and videos.

**File:** `backend/data/categories.json`  
**Model:** `models/content.go`

```go
type Category struct {
    ID        string       `json:"id"`         // UUID
    Name      string       `json:"name"`       // Display name
    Type      CategoryType `json:"type"`       // "article" | "video"
    CreatedAt string       `json:"created_at"` // ISO 8601
}
```

**JSON Storage:**

```json
{
  "category-uuid": {
    "id": "category-uuid",
    "name": "Gizi",
    "type": "article",
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

---

### 7. Article

Health education article (Berita).

**File:** `backend/data/articles.json`  
**Model:** `models/content.go`

```go
type Article struct {
    ID              string        `json:"id"`         // UUID
    Title           string        `json:"title"`
    Slug            string        `json:"slug"`       // URL-friendly, unique
    Excerpt         string        `json:"excerpt"`    // Short description
    Content         string        `json:"content"`    // Quill Delta JSON
    AuthorID        string        `json:"author_id"`  // User ID
    CategoryID      string        `json:"category_id"`
    HeroImages      HeroImages    `json:"hero_images"` // 3 aspect ratios
    Status          ArticleStatus `json:"status"`     // draft|published
    Version         int           `json:"version"`    // Incremented on update
    ViewCount       int           `json:"view_count"` // Incremented on read
    AttachmentCount int           `json:"attachment_count"` // Reminder attachments
    CreatedAt       string        `json:"created_at"`
    PublishedAt     string        `json:"published_at"`
    UpdatedAt       string        `json:"updated_at"`
}
```

**Hero Images:**

```go
type HeroImages struct {
    Hero16x9 string `json:"hero_16x9"` // 1920x1080
    Hero1x1  string `json:"hero_1x1"`  // 1080x1080
    Hero4x3  string `json:"hero_4x3"`  // 1600x1200
}
```

**Slug Generation:**

- Lowercase
- Replace spaces with hyphens
- Remove special characters
- Ensure uniqueness

**JSON Storage:**

```json
{
  "article-uuid": {
    "id": "article-uuid",
    "title": "Kebersihan Tangan",
    "slug": "kebersihan-tangan",
    "excerpt": "Pentingnya mencuci tangan...",
    "content": "{\"ops\":[...]}",
    "author_id": "admin-user-id",
    "category_id": "category-uuid",
    "hero_images": {
      "hero_16x9": "/uploads/article-uuid-16x9.jpg",
      "hero_1x1": "/uploads/article-uuid-1x1.jpg",
      "hero_4x3": "/uploads/article-uuid-4x3.jpg"
    },
    "status": "published",
    "version": 2,
    "view_count": 125,
    "attachment_count": 8,
    "created_at": "2026-01-01T00:00:00Z",
    "published_at": "2026-01-01T12:00:00Z",
    "updated_at": "2026-01-02T10:00:00Z"
  }
}
```

---

### 8. Video

Educational YouTube video (Video Edukasi).

**File:** `backend/data/videos.json`  
**Model:** `models/content.go`

```go
type Video struct {
    ID              string      `json:"id"`         // UUID
    YouTubeURL      string      `json:"youtube_url"`
    YouTubeID       string      `json:"youtube_id"` // Extracted from URL
    Title           string      `json:"title"`      // From YouTube metadata
    Description     string      `json:"description"`
    ChannelName     string      `json:"channel_name"`
    ThumbnailURL    string      `json:"thumbnail_url"`
    Duration        string      `json:"duration"`   // e.g., "10:25"
    CategoryID      string      `json:"category_id"`
    Status          VideoStatus `json:"status"`     // "published" only
    ViewCount       int         `json:"view_count"`
    AttachmentCount int         `json:"attachment_count"`
    CreatedAt       string      `json:"created_at"`
    UpdatedAt       string      `json:"updated_at"`
}
```

**Metadata Fetching:**

- Source: noembed.com API
- URL: `https://noembed.com/embed?url=<youtube_url>`
- Automatic on video creation

**JSON Storage:**

```json
{
  "video-uuid": {
    "id": "video-uuid",
    "youtube_url": "https://www.youtube.com/watch?v=VIDEO_ID",
    "youtube_id": "VIDEO_ID",
    "title": "Cara Cuci Tangan yang Benar",
    "description": "Tutorial lengkap...",
    "channel_name": "Kementerian Kesehatan",
    "thumbnail_url": "https://i.ytimg.com/vi/VIDEO_ID/maxresdefault.jpg",
    "duration": "5:32",
    "category_id": "category-uuid",
    "status": "published",
    "view_count": 0,
    "attachment_count": 12,
    "created_at": "2026-01-01T00:00:00Z",
    "updated_at": "2026-01-01T00:00:00Z"
  }
}
```

---

## Data Stores (Thread-Safety)

### PatientStore

```go
type PatientStore struct {
    Mu       sync.RWMutex
    Patients map[string]*Patient
    SaveFunc func()
}
```

**Methods:**

- `GetPatient(id) (*Patient, bool)` - Read lock
- `SaveData()` - Triggers JSON file write
- `Lock() / Unlock()` - Write operations
- `RLock() / RUnlock()` - Read operations

**Usage Pattern:**

```go
store.RLock()
patient, exists := store.GetPatient(id)
store.RUnlock()

store.Lock()
store.Patients[id] = newPatient
store.SaveData()
store.Unlock()
```

---

### ContentStore

Manages categories, articles, and videos in memory.

**Structure:**

```go
type CategoryStore struct {
    Mu         sync.RWMutex
    Categories map[string]*Category
    ByType     map[CategoryType][]string
}

type ArticleStore struct {
    Mu         sync.RWMutex
    Articles   map[string]*Article
    BySlug     map[string]string       // slug -> article ID
    ByCategory map[string][]string     // category ID -> article IDs
}

type VideoStore struct {
    Mu         sync.RWMutex
    Videos     map[string]*Video
    ByCategory map[string][]string
}
```

**Indexes:**

- `BySlug` - Fast article lookup by slug
- `ByCategory` - Filter content by category
- `ByType` - Filter categories by type

---

## Relationships

### Entity Relationship Diagram

```
User (1) ─────── (N) Patient
                    │
                    │ (1)
                    │
                    │ (N)
                 Reminder ─── (N) Attachment
                    │
                    │
                 (references)
                    │
            ┌───────┴────────┐
            │                │
        Article          Video
            │                │
            │                │
        Category         Category
```

### Foreign Keys (Logical)

| Child Entity | Parent Entity | Foreign Key Field | Relationship   |
| ------------ | ------------- | ----------------- | -------------- |
| Patient      | User          | `createdBy`       | N:1            |
| Reminder     | Patient       | (embedded)        | N:1            |
| Attachment   | Article/Video | `id` (validated)  | N:1            |
| Reminder     | User          | `cancelled_by`    | N:1 (optional) |
| Article      | User          | `author_id`       | N:1            |
| Article      | Category      | `category_id`     | N:1            |
| Video        | Category      | `category_id`     | N:1            |

---

## Validation Rules

### Phone Numbers

- Format: Indonesian mobile (+62 or 08)
- Normalization: `628xxxxxxxxx` (remove leading 0, add 62)
- Validation: 10-13 digits after 62

**Example:**

- Input: `08123456789` or `+628123456789`
- Stored: `628123456789`

### Email

- Optional field
- Standard email format validation

### Reminder Attachments

- Maximum: 3 attachments
- Types: `article` or `video`
- Content ID must exist in CMS
- Increments `attachment_count` on content entity

### Article Slug

- Auto-generated from title
- Lowercase, hyphens for spaces
- Unique constraint
- Regenerated on title change

---

## Concurrency Strategy

**Read Operations:**

```go
store.RLock()
defer store.RUnlock()
// Safe concurrent reads
```

**Write Operations:**

```go
store.Lock()
defer store.Unlock()
// Exclusive write access
store.SaveData()
```

**Why RWMutex?**

- Multiple goroutines can read simultaneously
- Only one goroutine can write
- Prevents race conditions on map access
- No database connection overhead

---

## Persistence Implementation

### File Write Pattern

```go
func saveData() {
    store.Lock()
    defer store.Unlock()

    data, err := json.MarshalIndent(store.Patients, "", "  ")
    if err != nil {
        log.Printf("Marshal error: %v", err)
        return
    }

    err = os.WriteFile("data/patients.json", data, 0644)
    if err != nil {
        log.Printf("Write error: %v", err)
    }
}
```

### Load on Startup

```go
func loadData() {
    data, err := os.ReadFile("data/patients.json")
    if err != nil {
        if os.IsNotExist(err) {
            return // Empty store
        }
        log.Fatalf("Load error: %v", err)
    }

    store.Lock()
    defer store.Unlock()

    err = json.Unmarshal(data, &store.Patients)
    if err != nil {
        log.Fatalf("Unmarshal error: %v", err)
    }
}
```

---

## Backup & Migration

### Backup Strategy

1. **Automatic:** Copy JSON files before writes (optional)
2. **Manual:** Copy entire `backend/data/` directory
3. **Version Control:** Exclude from git (add to `.gitignore`)

### Data Migration

**To Database (Future):**

1. Read JSON files
2. Map structs to database schema
3. Bulk insert with transactions
4. Verify data integrity
5. Switch persistence layer

**Data remains compatible** - Go structs use JSON tags compatible with database ORMs (GORM, sqlx).

---

## Security Considerations

### Password Storage

- **Algorithm:** SHA256
- **Encoding:** Base64
- **Salt:** None (consider adding for production)

### Phone Number Masking

In logs and admin views:

```go
// 628123456789 → 6281234***789
masked := phone[:7] + "***" + phone[len(phone)-3:]
```

### Email Masking

```go
// john@example.com → j***@example.com
parts := strings.Split(email, "@")
masked := string(parts[0][0]) + "***@" + parts[1]
```

---

## Performance Considerations

### In-Memory Storage

**Advantages:**

- ✅ Fast read/write (no network/disk I/O per operation)
- ✅ Simple deployment (no database server)
- ✅ Easy backup (copy files)

**Limitations:**

- ⚠️ All data in RAM
- ⚠️ No query optimization (full scans)
- ⚠️ Single-server only (no horizontal scaling)

**Suitable for:**

- Small to medium datasets (<10,000 patients)
- Single-server deployments
- Rapid prototyping

### Memory Footprint Estimate

- Patient with 10 reminders: ~2KB
- 1,000 patients: ~2MB
- 10,000 patients: ~20MB
- Articles (with content): ~50KB each
- 100 articles: ~5MB

**Total for 1,000 patients + 100 articles:** ~10MB RAM

---

**Next:** See [Integration Architecture](./integration-architecture.md) for how backend communicates with GOWA and frontend.
