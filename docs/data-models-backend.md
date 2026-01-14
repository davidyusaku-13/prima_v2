# PRIMA Data Models

**Date:** 2026-01-13

## Overview

This document describes the data models used in the PRIMA backend. All models are stored as JSON files in the `backend/data/` directory with thread-safe access using `sync.RWMutex`.

## Data Files

| File | Content | Key |
|------|---------|-----|
| `patients.json` | Patients and Reminders | `id` |
| `users.json` | User accounts | `id` |
| `categories.json` | Content categories | `id` |
| `articles.json` | Articles | `id` |
| `videos.json` | YouTube videos | `id` |
| `jwt_secret.txt` | JWT signing key | N/A |

---

## User Model

### JSON Structure

```json
{
  "id": "user_abc123",
  "username": "superadmin",
  "full_name": "Super Administrator",
  "password": "sha256_hash",
  "role": "superadmin",
  "created_at": "2026-01-01T00:00:00Z"
}
```

### Field Definitions

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier (auto-generated UUID) |
| `username` | string | Unique login username |
| `full_name` | string | Display name |
| `password` | string | SHA256 hash of password |
| `role` | string | `superadmin`, `admin`, or `volunteer` |
| `created_at` | string | RFC3339 UTC timestamp |

### Roles

| Role | Permissions |
|------|-------------|
| `superadmin` | Full access including user management |
| `admin` | CMS content, analytics, health checks |
| `volunteer` | CRUD own patients and reminders |

### Default Users

| Username | Password | Role |
|----------|----------|------|
| `superadmin` | `superadmin` | superadmin |

---

## Patient Model

### JSON Structure

```json
{
  "id": "patient_20260101120000_abc123",
  "name": "John Doe",
  "phone": "6281234567890",
  "email": "john@example.com",
  "notes": "Diabetes type 2 patient",
  "created_by": "user_abc123",
  "created_at": "2026-01-01T12:00:00Z",
  "updated_at": "2026-01-01T12:00:00Z",
  "reminders": []
}
```

### Field Definitions

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Format: `patient_YYYYMMDDHHMMSS_random8` |
| `name` | string | Patient's full name |
| `phone` | string | Normalized to `628...` format |
| `email` | string | Optional email address |
| `notes` | string | Optional notes about patient |
| `created_by` | string | User ID of creator |
| `created_at` | string | RFC3339 UTC timestamp |
| `updated_at` | string | RFC3339 UTC timestamp |
| `reminders` | array | Array of Reminder objects |

---

## Reminder Model

### JSON Structure

```json
{
  "id": "reminder_abc123",
  "title": "Morning Medication",
  "description": "Take metformin with breakfast",
  "due_date": "2026-01-02T08:00:00Z",
  "priority": "high",
  "completed": false,
  "recurrence": null,
  "notified": false,
  "attachments": [
    {
      "type": "article",
      "id": "article_abc123"
    }
  ],
  "gowa_message_id": "WA123",
  "delivery_status": "delivered",
  "delivery_error_message": "",
  "message_sent_at": "2026-01-02T08:00:05Z",
  "delivered_at": "2026-01-02T08:00:10Z",
  "read_at": "2026-01-02T08:15:00Z",
  "retry_count": 0,
  "scheduled_delivery_at": "2026-01-02T08:00:00Z",
  "cancelled_at": null,
  "cancelled_by": null
}
```

### Field Definitions

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier |
| `title` | string | Reminder title |
| `description` | string | Detailed description |
| `due_date` | string | RFC3339 or `2006-01-02T15:04` format |
| `priority` | string | `low`, `medium`, or `high` |
| `completed` | boolean | Whether reminder is completed |
| `recurrence` | object | Recurrence pattern (optional) |
| `notified` | boolean | Whether notification sent |
| `attachments` | array | Array of Attachment objects |
| `gowa_message_id` | string | Message ID from GOWA |
| `delivery_status` | string | See delivery status table |
| `delivery_error_message` | string | Error details if failed |
| `message_sent_at` | string | RFC3339 UTC timestamp |
| `delivered_at` | string | RFC3339 UTC timestamp (from webhook) |
| `read_at` | string | RFC3339 UTC timestamp (from webhook) |
| `retry_count` | int | Number of retry attempts |
| `scheduled_delivery_at` | string | RFC3339 UTC timestamp |
| `cancelled_at` | string | RFC3339 UTC timestamp if cancelled |
| `cancelled_by` | string | User ID who cancelled |

### Delivery Status Values

| Status | Description | Terminal? |
|--------|-------------|-----------|
| `pending` | Just created | No |
| `scheduled` | Queued for sending | No |
| `sending` | Currently sending | No |
| `sent` | Sent to GOWA | No |
| `delivered` | Delivered to WhatsApp | Yes |
| `read` | Read by recipient | Yes |
| `failed` | Send failed | Retryable |
| `retrying` | Waiting for retry | No |
| `cancelled` | User cancelled | Yes |

### Recurrence Model

```json
{
  "frequency": "daily|weekly|monthly",
  "interval": 1,
  "days_of_week": [1, 2, 3, 4, 5]
}
```

### Attachment Model

```json
{
  "type": "article|video",
  "id": "content_id"
}
```

---

## Category Model

### JSON Structure

```json
{
  "id": "cat_abc123",
  "name": "Diabetes Care",
  "type": "article",
  "created_at": "2026-01-01T00:00:00Z"
}
```

### Field Definitions

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier |
| `name` | string | Category name |
| `type` | string | `article` or `video` |
| `created_at` | string | RFC3339 UTC timestamp |

---

## Article Model

### JSON Structure

```json
{
  "id": "article_abc123",
  "title": "Understanding Diabetes",
  "slug": "understanding-diabetes",
  "excerpt": "A comprehensive guide...",
  "content": "<p>Full HTML content...</p>",
  "author_id": "user_abc123",
  "category_id": "cat_abc123",
  "hero_images": {
    "16x9": "/uploads/articles/abc123_16x9.jpg",
    "1x1": "/uploads/articles/abc123_1x1.jpg",
    "4x3": "/uploads/articles/abc123_4x3.jpg"
  },
  "status": "published",
  "version": 1,
  "view_count": 150,
  "attachment_count": 5,
  "created_at": "2026-01-01T00:00:00Z",
  "published_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

### Field Definitions

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier |
| `title` | string | Article title |
| `slug` | string | URL-friendly identifier |
| `excerpt` | string | Brief summary (max ~200 chars) |
| `content` | string | Full HTML content (Quill editor output) |
| `author_id` | string | User ID of author |
| `category_id` | string | Category ID |
| `hero_images` | object | Uploaded image URLs by aspect ratio |
| `status` | string | `draft` or `published` |
| `version` | int | Content version (for future use) |
| `view_count` | int | Number of views |
| `attachment_count` | int | How many reminders include this |
| `created_at` | string | RFC3339 UTC timestamp |
| `published_at` | string | RFC3339 UTC timestamp |
| `updated_at` | string | RFC3339 UTC timestamp |

---

## Video Model

### JSON Structure

```json
{
  "id": "video_abc123",
  "you_tube_url": "https://www.youtube.com/watch?v=abc123",
  "you_tube_id": "abc123",
  "title": "Healthy Eating Tips",
  "description": "Learn about balanced diets...",
  "channel_name": "Health Education",
  "thumbnail_url": "https://img.youtube.com/vi/abc123/maxresdefault.jpg",
  "duration": "10:30",
  "category_id": "cat_abc123",
  "status": "published",
  "view_count": 500,
  "attachment_count": 3,
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-01-01T00:00:00Z"
}
```

### Field Definitions

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier |
| `you_tube_url` | string | Original YouTube URL |
| `you_tube_id` | string | YouTube video ID |
| `title` | string | Video title (from noembed) |
| `description` | string | Video description (from noembed) |
| `channel_name` | string | YouTube channel name |
| `thumbnail_url` | string | Thumbnail image URL |
| `duration` | string | Duration (format: MM:SS or HH:MM:SS) |
| `category_id` | string | Category ID |
| `status` | string | `published` |
| `view_count` | int | View count (from noembed) |
| `attachment_count` | int | How many reminders include this |
| `created_at` | string | RFC3339 UTC timestamp |
| `updated_at` | string | RFC3339 UTC timestamp |

---

## Store Interfaces

### UserStore

```go
type UserStore struct {
    mu      sync.RWMutex
        users    map[string]*User
        usernames map[string]string  // username -> id
}

func (s *UserStore) GetByID(id string) (*User, error)
func (s *UserStore) GetByUsername(username string) (*User, error)
func (s *UserStore) Create(user *User) error
func (s *UserStore) Update(user *User) error
func (s *UserStore) Delete(id string) error
func (s *UserStore) List() ([]*User, error)
```

### PatientStore

```go
type PatientStore struct {
    mu      sync.RWMutex
    patients map[string]*Patient
}

func (s *PatientStore) GetByID(id string) (*Patient, error)
func (s *PatientStore) GetByCreator(userID string) ([]*Patient, error)
func (s *PatientStore) Create(patient *Patient) error
func (s *PatientStore) Update(patient *Patient) error
func (s *PatientStore) Delete(id string) error
func (s *PatientStore) List() ([]*Patient, error)
```

### ReminderStore

```go
type ReminderStore struct {
    mu        sync.RWMutex
    reminders map[string]*Reminder
}

func (s *ReminderStore) GetByID(id string) (*Reminder, error)
func (s *ReminderStore) GetByPatient(patientID string) ([]*Reminder, error)
func (s *ReminderStore) Create(reminder *Reminder) error
func (s *ReminderStore) Update(reminder *Reminder) error
func (s *ReminderStore) Delete(id string) error
func (s *ReminderStore) GetScheduled() ([]*Reminder, error)
func (s *ReminderStore) GetFailed() ([]*Reminder, error)
```

---

_Generated using BMAD Method `document-project` workflow_
