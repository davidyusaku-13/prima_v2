# API Contracts - Backend

**Generated:** January 11, 2026 (Updated)
**Base URL:** `http://localhost:8080/api`
**Project:** PRIMA Healthcare Volunteer Dashboard
**Scan Type:** Exhaustive Rescan

---

## Authentication

**Type:** JWT Bearer Token  
**Header:** `Authorization: Bearer <token>`  
**Token Expiry:** 7 days  
**Roles:** `superadmin`, `admin`, `volunteer`

---

## Endpoints Overview

| Category          | Public | Authenticated | Admin Only | Superadmin Only |
| ----------------- | ------ | ------------- | ---------- | --------------- |
| **Auth**          | 2      | 1             | -          | -               |
| **Patients**      | -      | 5             | -          | -               |
| **Reminders**     | -      | 10            | -          | -               |
| **Content (CMS)** | 9      | 2             | 7          | -               |
| **Analytics**     | -      | -             | 5          | -               |
| **Health**        | 1      | 1             | -          | -               |
| **SSE**           | -      | 1             | -          | -               |
| **Webhooks**      | 1      | -             | -          | -               |
| **Users**         | -      | 1             | -          | 3               |
| **Config**        | 2      | -             | -          | -               |

---

## 1. Authentication

### POST /api/auth/register

Register a new user account (always creates `volunteer` role).

**Authentication:** None  
**Request Body:**

```json
{
  "username": "string (required, 3-30 chars)",
  "password": "string (required, min 6 chars)",
  "fullName": "string (optional)"
}
```

**Response (201):**

```json
{
  "message": "user registered successfully",
  "userId": "uuid",
  "username": "string",
  "fullName": "string",
  "role": "volunteer",
  "token": "jwt-token",
  "expiresAt": "2026-01-09T00:00:00Z"
}
```

**Errors:**

- `400 VALIDATION_ERROR` - Invalid username/password
- `409 CONFLICT` - Username already exists

---

### POST /api/auth/login

Login with username and password.

**Authentication:** None  
**Request Body:**

```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Response (200):**

```json
{
  "token": "jwt-token",
  "userId": "uuid",
  "username": "string",
  "fullName": "string",
  "role": "volunteer|admin|superadmin",
  "expiresAt": "2026-01-09T00:00:00Z"
}
```

**Errors:**

- `401 UNAUTHORIZED` - Invalid credentials

---

### GET /api/auth/me

Get current authenticated user details.

**Authentication:** Required  
**Response (200):**

```json
{
  "userId": "uuid",
  "username": "string",
  "fullName": "string",
  "role": "volunteer|admin|superadmin"
}
```

---

## 2. Patients

### GET /api/patients

Get list of patients (filtered by role: volunteers see only their own, admins see all).

**Authentication:** Required  
**Roles:** All authenticated users

**Response (200):**

```json
{
  "patients": [
    {
      "id": "uuid",
      "name": "string",
      "phone": "628123456789 (normalized)",
      "email": "string",
      "notes": "string",
      "reminders": [],
      "createdBy": "userId"
    }
  ]
}
```

---

### POST /api/patients

Create a new patient.

**Authentication:** Required  
**Roles:** All authenticated users

**Request Body:**

```json
{
  "name": "string (required)",
  "phone": "string (required, Indonesian format)",
  "email": "string (optional)",
  "notes": "string (optional)"
}
```

**Response (201):** Patient object

**Errors:**

- `400 INVALID_PHONE` - Invalid Indonesian phone number

---

### GET /api/patients/:id

Get a single patient by ID.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Response (200):** Patient object

**Errors:**

- `403 FORBIDDEN` - Cannot access other volunteer's patients
- `404 NOT_FOUND` - Patient not found

---

### PUT /api/patients/:id

Update patient information.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Request Body:** Partial patient object

**Response (200):** Updated patient object

---

### DELETE /api/patients/:id

Delete a patient.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Response (200):**

```json
{
  "message": "patient deleted"
}
```

---

## 3. Reminders

### GET /api/patients/:id/reminders

Get reminders for a patient with pagination.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Query Parameters:**

- `history` - boolean (default: false) - include cancelled reminders
- `page` - integer (default: 1)
- `limit` - integer (default: 20, max: 100)

**Response (200):**

```json
{
  "data": [
    {
      "id": "uuid",
      "title": "string",
      "message": "string",
      "message_preview": "string (first 100 chars)",
      "scheduled_at": "2026-01-03T10:00:00Z",
      "delivery_status": "pending|scheduled|queued|sending|retrying|sent|delivered|read|failed|expired|cancelled",
      "delivery_error": "string",
      "sent_at": "2026-01-03T10:00:05Z",
      "delivered_at": "2026-01-03T10:00:10Z",
      "read_at": "2026-01-03T10:05:00Z",
      "cancelled_at": null,
      "attachments": [
        {
          "type": "article|video",
          "id": "content-id",
          "title": "string",
          "url": "string"
        }
      ],
      "attachment_count": 1
    }
  ],
  "message": "Success",
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45,
    "has_more": true
  }
}
```

---

### POST /api/patients/:id/reminders

Create a new reminder for a patient.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Request Body:**

```json
{
  "title": "string (required, max 200 chars)",
  "description": "string (optional)",
  "dueDate": "2026-01-10T14:00:00Z",
  "priority": "low|medium|high (default: medium)",
  "recurrence": {
    "frequency": "none|daily|weekly|monthly (default: none)",
    "interval": 1
  },
  "attachments": [
    {
      "type": "article|video",
      "id": "content-id (must exist in CMS)",
      "title": "string (max 200 chars)",
      "url": "string (optional)"
    }
  ]
}
```

**Response (201):** Reminder object

**Errors:**

- `400 MAX_ATTACHMENTS_EXCEEDED` - Maximum 3 content attachments allowed
- `400 INVALID_ATTACHMENT` - Content ID not found or invalid type

---

### PUT /api/patients/:id/reminders/:reminderId

Update a reminder.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Request Body:** Partial reminder object

**Response (200):** Updated reminder object

---

### POST /api/patients/:id/reminders/:reminderId/toggle

Toggle reminder completed status.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Response (200):** Updated reminder object

---

### DELETE /api/patients/:id/reminders/:reminderId

Delete a reminder.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Response (200):**

```json
{
  "message": "reminder deleted"
}
```

---

### POST /api/patients/:id/reminders/:reminderId/send

Manually send a reminder via WhatsApp (respects quiet hours 22:00-06:00 WIB).

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Response (200):**

```json
{
  "data": {
    "id": "reminder-id",
    "delivery_status": "sent|scheduled",
    "gowa_message_id": "gowa-id",
    "sent_at": "2026-01-02T10:00:00Z"
  },
  "message": "Reminder berhasil dikirim",
  "scheduled": false,
  "scheduled_at": null
}
```

**Response (200 - Quiet Hours):**

```json
{
  "data": {
    "id": "reminder-id",
    "delivery_status": "scheduled"
  },
  "message": "Reminder dijadwalkan (quiet hours)",
  "scheduled": true,
  "scheduled_at": "2026-01-03T06:00:00Z"
}
```

**Errors:**

- `400 INVALID_PHONE` - Invalid WhatsApp number format
- `409 ALREADY_SENDING` - Reminder is already being sent
- `503 GOWA_UNAVAILABLE` - GOWA service down (circuit breaker open)
- `503 RETRY_SCHEDULED` - GOWA error, scheduled for retry

---

### GET /api/reminders/:id/status

Get delivery status of a reminder.

**Authentication:** Required

**Response (200):**

```json
{
  "data": {
    "id": "uuid",
    "title": "string",
    "delivery_status": "pending|scheduled|queued|sending|retrying|sent|delivered|read|failed|expired|cancelled",
    "retry_count": 2,
    "max_attempts": 5,
    "error_message": "string",
    "scheduled_at": "2026-01-03T06:00:00Z",
    "sent_at": "2026-01-02T10:00:00Z",
    "gowa_message_id": "gowa-msg-id",
    "patient_id": "patient-uuid"
  }
}
```

---

### POST /api/reminders/:id/retry

Manually retry sending a failed reminder.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Response (200):**

```json
{
  "data": {
    "reminder_id": "uuid",
    "status": "sent|retrying",
    "message_id": "gowa-msg-id"
  },
  "message": "Reminder berhasil dikirim ulang"
}
```

**Errors:**

- `400 INVALID_STATUS` - Only failed reminders can be retried
- `404 REMINDER_NOT_FOUND` - Reminder doesn't exist
- `503 GOWA_ERROR` - GOWA service unavailable

---

### POST /api/reminders/:id/cancel

Cancel a pending or scheduled reminder.

**Authentication:** Required  
**Roles:** Owner (volunteer) or admin/superadmin

**Response (200):**

```json
{
  "data": {
    "id": "uuid",
    "title": "string",
    "status": "cancelled",
    "cancelled_at": "2026-01-02T11:00:00Z"
  },
  "message": "Reminder berhasil dibatalkan"
}
```

**Errors:**

- `400 CANNOT_CANCEL` - Reminder status does not allow cancellation (already sent/delivered)

---

## 4. Content (CMS)

### GET /api/categories

Get all content categories (public).

**Authentication:** None

**Response (200):**

```json
{
  "categories": [
    {
      "id": "uuid",
      "name": "Gizi",
      "type": "article|video",
      "createdAt": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

### GET /api/categories/:type

Get categories filtered by type (public).

**Authentication:** None  
**Path Parameters:** `type` - "article" or "video"

**Response (200):** Array of category objects

---

### POST /api/categories

Create a new category.

**Authentication:** Required  
**Roles:** admin, superadmin

**Request Body:**

```json
{
  "name": "string (required)",
  "type": "article|video (required)"
}
```

**Response (201):** Category object

---

### GET /api/articles

Get published articles (or all if user is admin with `all=true`).

**Authentication:** None (public for published)

**Query Parameters:**

- `category` - string (category ID, optional)
- `all` - boolean (default: false, requires admin role)

**Response (200):**

```json
{
  "articles": [
    {
      "id": "uuid",
      "title": "string",
      "slug": "kebersihan-tangan",
      "excerpt": "string",
      "content": "string (Quill Delta JSON)",
      "authorId": "user-id",
      "categoryId": "category-id",
      "heroImages": {
        "hero_16x9": "/uploads/article-uuid-16x9.jpg",
        "hero_1x1": "/uploads/article-uuid-1x1.jpg",
        "hero_4x3": "/uploads/article-uuid-4x3.jpg"
      },
      "status": "draft|published",
      "version": 3,
      "viewCount": 125,
      "attachmentCount": 8,
      "createdAt": "2026-01-01T00:00:00Z",
      "updatedAt": "2026-01-02T10:00:00Z",
      "publishedAt": "2026-01-01T12:00:00Z"
    }
  ]
}
```

---

### GET /api/articles/:slug

Get a single article by slug (increments view count).

**Authentication:** None  
**Path Parameters:** `slug` - Article slug (e.g., "kebersihan-tangan")

**Response (200):** Article object

**Errors:**

- `404 NOT_FOUND` - Article not found

---

### POST /api/articles

Create a new article.

**Authentication:** Required  
**Roles:** admin, superadmin

**Request Body:**

```json
{
  "title": "string (required)",
  "excerpt": "string (optional)",
  "content": "string (Quill Delta JSON, optional)",
  "category_id": "uuid (optional)",
  "status": "draft|published (default: draft)",
  "hero_images": {
    "hero_16x9": "/uploads/image-16x9.jpg",
    "hero_1x1": "/uploads/image-1x1.jpg",
    "hero_4x3": "/uploads/image-4x3.jpg"
  }
}
```

**Response (201):** Article object with auto-generated slug

---

### PUT /api/articles/:id

Update an article.

**Authentication:** Required  
**Roles:** admin, superadmin

**Path Parameters:** `id` - Article ID

**Request Body:** Partial article object

**Response (200):** Updated article object (version incremented)

---

### DELETE /api/articles/:id

Delete an article.

**Authentication:** Required  
**Roles:** admin, superadmin

**Response (200):**

```json
{
  "message": "article deleted"
}
```

---

### GET /api/videos

Get published videos.

**Authentication:** None

**Query Parameters:**

- `category` - string (category ID, optional)

**Response (200):**

```json
{
  "videos": [
    {
      "id": "uuid",
      "youtubeUrl": "https://www.youtube.com/watch?v=VIDEO_ID",
      "youtubeId": "VIDEO_ID",
      "title": "string (from YouTube metadata)",
      "description": "string",
      "channelName": "string",
      "thumbnailUrl": "https://i.ytimg.com/vi/VIDEO_ID/maxresdefault.jpg",
      "duration": "10:25",
      "categoryId": "uuid",
      "status": "published",
      "viewCount": 0,
      "attachmentCount": 5,
      "createdAt": "2026-01-01T00:00:00Z",
      "updatedAt": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

### POST /api/videos

Add a YouTube video (fetches metadata from noembed.com).

**Authentication:** Required  
**Roles:** admin, superadmin

**Request Body:**

```json
{
  "youtube_url": "https://www.youtube.com/watch?v=VIDEO_ID (required)",
  "category_id": "uuid (optional)"
}
```

**Response (201):** Video object

**Errors:**

- `400 INVALID_URL` - Invalid YouTube URL
- `409 CONFLICT` - Video already exists

---

### DELETE /api/videos/:id

Delete a video.

**Authentication:** Required  
**Roles:** admin, superadmin

**Response (200):**

```json
{
  "message": "video deleted"
}
```

---

### GET /api/content

Get all published content (articles and videos combined).

**Authentication:** None

**Query Parameters:**

- `type` - string (all|article|video, default: all)
- `category` - string (category ID, optional)

**Response (200):**

```json
{
  "articles": [],
  "videos": []
}
```

---

### GET /api/content/popular

Get most frequently attached content.

**Authentication:** None

**Query Parameters:**

- `limit` - integer (default: 5, max: 20)

**Response (200):**

```json
{
  "content": [
    {
      "id": "uuid",
      "type": "article|video",
      "title": "string",
      "thumbnail": "/uploads/image.jpg"
    }
  ]
}
```

---

### POST /api/content/:type/:id/increment-attachment

Increment attachment count for content (internal use when reminder is created with attachment).

**Authentication:** Required

**Path Parameters:**

- `type` - "article" or "video"
- `id` - Content ID

**Response (200):**

```json
{
  "message": "attachment count incremented"
}
```

---

### POST /api/upload/image

Upload and resize image for article hero images (3 aspect ratios: 16:9, 1:1, 4:3).

**Authentication:** Required  
**Roles:** admin, superadmin

**Content-Type:** `multipart/form-data`

**Form Data:**

- `image` - file (JPEG, PNG, or WebP, max 10MB)

**Response (200):**

```json
{
  "images": {
    "hero_16x9": "/uploads/uuid-16x9.jpg",
    "hero_1x1": "/uploads/uuid-1x1.jpg",
    "hero_4x3": "/uploads/uuid-4x3.jpg"
  }
}
```

**Errors:**

- `400 INVALID_FILE` - Invalid file format or size

---

## 5. Analytics

### GET /api/analytics/delivery

Get delivery statistics for admin dashboard.

**Authentication:** Required  
**Roles:** admin, superadmin

**Query Parameters:**

- `period` - string (today|7d|30d|all, default: all)

**Response (200):**

```json
{
  "data": {
    "totalSent": 1250,
    "successRate": 94.5,
    "failedLast7Days": 15,
    "avgDeliveryTime": "2.5 seconds",
    "breakdown": {
      "pending": 5,
      "scheduled": 12,
      "queued": 3,
      "sending": 2,
      "retrying": 4,
      "sent": 800,
      "delivered": 650,
      "read": 420,
      "failed": 18,
      "expired": 2
    },
    "period": "all",
    "periodStartDate": "2025-12-01T00:00:00Z",
    "periodEndDate": "2026-01-02T23:59:59Z"
  },
  "message": "success"
}
```

---

### GET /api/analytics/failed-deliveries

Get paginated list of failed deliveries.

**Authentication:** Required  
**Roles:** admin, superadmin

**Query Parameters:**

- `page` - integer (default: 1)
- `limit` - integer (default: 20, max: 100)
- `reason` - string (invalid_phone|gowa_timeout|message_rejected|other, optional)

**Response (200):**

```json
{
  "data": {
    "items": [
      {
        "reminder_id": "uuid",
        "patient_name_masked": "John D***",
        "phone_masked": "6281234***789",
        "volunteer_name": "Jane Smith",
        "reminder_title": "Minum Obat Pagi",
        "failure_reason": "GOWA Timeout",
        "failure_reason_code": "gowa_timeout",
        "failure_timestamp": "2026-01-02T10:30:00Z",
        "retry_count": 5,
        "delivery_error_message": "Request timeout after 30 seconds"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 45,
      "total_pages": 3
    },
    "filter_counts": {
      "invalid_phone": 8,
      "gowa_timeout": 25,
      "message_rejected": 5,
      "other": 7
    }
  },
  "message": "Failed deliveries retrieved successfully"
}
```

---

### GET /api/analytics/failed-deliveries/export

Export failed deliveries as CSV.

**Authentication:** Required  
**Roles:** admin, superadmin

**Query Parameters:**

- `reason` - string (filter by reason, optional)

**Response (200):**

- **Content-Type:** `text/csv`
- **Header:** `Content-Disposition: attachment; filename=failed-deliveries-2026-01-02.csv`

**CSV Columns:**

```
Reminder ID,Patient Name,Phone,Volunteer,Title,Reason,Reason Code,Timestamp,Retry Count,Error
```

---

### GET /api/analytics/failed-deliveries/:id

Get detailed information about a single failed delivery.

**Authentication:** Required  
**Roles:** admin, superadmin

**Path Parameters:** `id` - Reminder ID

**Response (200):**

```json
{
  "data": {
    "reminder_id": "uuid",
    "patient_name_masked": "John D***",
    "phone_masked": "6281234***789",
    "volunteer_name": "Jane Smith",
    "reminder_title": "Minum Obat Pagi",
    "delivery_error_message": "Request timeout after 30 seconds",
    "failure_reason": "GOWA Timeout",
    "failure_reason_code": "gowa_timeout",
    "failure_timestamp": "2026-01-02T10:30:00Z",
    "retry_count": 5,
    "retry_attempts": [
      {
        "attempt": 1,
        "timestamp": "2026-01-02T10:00:00Z",
        "error": "Connection refused"
      },
      {
        "attempt": 2,
        "timestamp": "2026-01-02T10:05:00Z",
        "error": "Timeout"
      }
    ]
  },
  "message": "Failed delivery details retrieved"
}
```

---

### GET /api/analytics/content

Get content attachment statistics.

**Authentication:** Required  
**Roles:** admin, superadmin

**Response (200):**

```json
{
  "data": {
    "articles": [
      {
        "id": "uuid",
        "title": "Kebersihan Tangan",
        "attachmentCount": 45,
        "type": "article"
      }
    ],
    "videos": [
      {
        "id": "uuid",
        "title": "Cara Cuci Tangan",
        "attachmentCount": 32,
        "type": "video"
      }
    ],
    "topContent": []
  }
}
```

---

### GET /api/dashboard/stats

Get CMS dashboard statistics.

**Authentication:** Required  
**Roles:** admin, superadmin

**Response (200):**

```json
{
  "categories": {
    "articles": 8,
    "videos": 5
  },
  "articles": {
    "published": 25,
    "drafts": 5,
    "total": 30
  },
  "videos": {
    "total": 18
  },
  "total_views": {
    "articles": 1250
  }
}
```

---

## 6. Health

### GET /api/health

Basic health check endpoint (public).

**Authentication:** None

**Response (200):**

```json
{
  "data": {
    "status": "ok",
    "timestamp": "2026-01-02T10:00:00Z"
  },
  "message": "Health check successful"
}
```

---

### GET /api/health/detailed

Detailed health status with GOWA connectivity and queue info.

**Authentication:** Required  
**Roles:** admin, superadmin

**Response (200):**

```json
{
  "data": {
    "status": "ok",
    "timestamp": "2026-01-02T10:00:00Z",
    "gowa": {
      "connected": true,
      "last_ping": "2026-01-02T09:59:50Z",
      "endpoint": "http://localhost:3000"
    },
    "circuit_breaker": {
      "state": "closed",
      "failure_count": 0,
      "cooldown_remaining_seconds": 0
    },
    "queue": {
      "total": 15,
      "scheduled": 8,
      "retrying": 3,
      "quiet_hours": 4
    }
  },
  "message": "Detailed health status retrieved"
}
```

---

## 7. Server-Sent Events (SSE)

### GET /api/sse/delivery-status

Real-time delivery status updates via Server-Sent Events.

**Authentication:** Query parameter (JWT token)  
**Query Parameters:**

- `token` - string (JWT token, required)

**Content-Type:** `text/event-stream`

**Events:**

1. **connection.established**

```json
{
  "message": "Connected to delivery status updates",
  "timestamp": "2026-01-02T10:00:00Z"
}
```

2. **delivery.status.updated**

```json
{
  "reminder_id": "uuid",
  "status": "sent|delivered|read",
  "timestamp": "2026-01-02T10:00:05Z"
}
```

3. **delivery.failed**

```json
{
  "reminder_id": "uuid",
  "patient_id": "uuid",
  "patient_name": "John Doe",
  "error": "GOWA timeout",
  "timestamp": "2026-01-02T10:00:10Z"
}
```

**Usage Example (JavaScript):**

```javascript
const token = localStorage.getItem("token");
const eventSource = new EventSource(`/api/sse/delivery-status?token=${token}`);

eventSource.addEventListener("delivery.status.updated", (e) => {
  const data = JSON.parse(e.data);
  console.log("Status updated:", data);
});
```

---

## 8. Webhooks

### POST /api/webhook/gowa

GOWA webhook callback for delivery status updates (uses HMAC-SHA256 validation).

**Authentication:** HMAC signature in header  
**Headers:**

- `X-Webhook-Signature` - string (HMAC-SHA256 signature, required)

**Request Body:**

```json
{
  "event": "message.ack",
  "message": {
    "id": "gowa-message-id",
    "status": "delivered|read|failed"
  }
}
```

**Response (200):**

```json
{
  "data": {
    "message_id": "gowa-msg-id",
    "reminder_id": "uuid",
    "delivery_status": "delivered"
  },
  "message": "Reminder status updated to 'delivered'"
}
```

**Errors:**

- `401 MISSING_SIGNATURE` - Missing X-Webhook-Signature header
- `401 INVALID_SIGNATURE` - HMAC signature verification failed
- `400 INVALID_PAYLOAD` - Malformed webhook payload

---

## 9. Users

### GET /api/users

Get all users.

**Authentication:** Required  
**Roles:** superadmin

**Response (200):**

```json
{
  "users": [
    {
      "id": "uuid",
      "username": "string",
      "fullName": "string",
      "role": "superadmin|admin|volunteer",
      "createdAt": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

### GET /api/users/:id

Get user details by ID (for author lookup).

**Authentication:** Required

**Response (200):** User object (without password)

---

### PUT /api/users/:id/role

Update user role (cannot change own role).

**Authentication:** Required  
**Roles:** superadmin

**Request Body:**

```json
{
  "role": "admin|volunteer (required)"
}
```

**Response (200):**

```json
{
  "message": "role updated successfully",
  "user": {
    "id": "uuid",
    "username": "string",
    "role": "admin"
  }
}
```

**Errors:**

- `403 FORBIDDEN` - Cannot change own role

---

### DELETE /api/users/:id

Delete a user (cannot delete yourself).

**Authentication:** Required  
**Roles:** superadmin

**Response (200):**

```json
{
  "message": "user deleted successfully"
}
```

**Errors:**

- `403 FORBIDDEN` - Cannot delete yourself

---

## 10. Config

### GET /api/config/disclaimer

Get disclaimer configuration.

**Authentication:** None

**Response (200):**

```json
{
  "data": {
    "text": "Informasi ini untuk tujuan edukasi. Konsultasikan dengan tenaga kesehatan untuk kondisi spesifik Anda.",
    "enabled": true
  }
}
```

---

### GET /api/config/quiet-hours

Get quiet hours configuration.

**Authentication:** None

**Response (200):**

```json
{
  "data": {
    "start_hour": 21,
    "end_hour": 6,
    "timezone": "WIB"
  }
}
```

---

## Error Response Format

All API errors follow this structure:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE (optional)",
  "details": "Additional context (optional)"
}
```

**Common HTTP Status Codes:**

- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., duplicate)
- `500 Internal Server Error` - Server error
- `503 Service Unavailable` - External service down (e.g., GOWA)

---

## Special Features

### Circuit Breaker Pattern

GOWA integration uses circuit breaker to prevent cascading failures:

- **Closed:** Normal operation
- **Open:** After 5 consecutive failures, stops sending for 5 minutes
- **Half-Open:** Tests with single request after cooldown

### Retry Strategy

Failed reminder deliveries are automatically retried with exponential backoff:

1. **Attempt 1:** Immediate
2. **Attempt 2:** +1 second
3. **Attempt 3:** +5 seconds
4. **Attempt 4:** +30 seconds
5. **Attempt 5:** +2 minutes
6. **Attempt 6:** +10 minutes

After max attempts, status becomes `failed`.

### Quiet Hours Enforcement

WhatsApp reminders are NOT sent between **22:00-06:00 WIB**. They are automatically scheduled for 06:00 WIB the next day.

### Real-Time Updates

Use SSE endpoint `/api/sse/delivery-status` for real-time delivery notifications without polling.

---

**Generated by:** BMad Document Project Workflow  
**Last Updated:** January 2, 2026
