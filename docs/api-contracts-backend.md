# PRIMA API Contracts

**Date:** 2026-01-13

## Overview

This document describes the REST API contracts for the PRIMA Backend. All endpoints return JSON responses.

## Base URL

```
http://localhost:8080/api
```

## Authentication

All protected endpoints require:
- Header: `Authorization: Bearer <jwt_token>`

## Response Formats

### Success Response

```json
{
  "success": true,
  "data": { ... }
}
```

### Error Response

```json
{
  "success": false,
  "error": "Error message here"
}
```

---

## Authentication Endpoints

### Register User

**POST** `/auth/register`

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required, min 6 chars)",
  "full_name": "string (optional)"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "string",
      "username": "string",
      "full_name": "string",
      "role": "volunteer"
    },
    "token": "jwt_token_string"
  }
}
```

### Login

**POST** `/auth/login`

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "string",
      "username": "string",
      "full_name": "string",
      "role": "superadmin|admin|volunteer"
    },
    "token": "jwt_token_string"
  }
}
```

### Get Current User

**GET** `/auth/me`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "string",
    "username": "string",
    "full_name": "string",
    "role": "string",
    "created_at": "string"
  }
}
```

---

## Patient Endpoints

### List Patients

**GET** `/patients`

**Query Parameters:**
- `limit` (optional, default 50)
- `offset` (optional, default 0)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "patients": [
      {
        "id": "string",
        "name": "string",
        "phone": "628xxx",
        "email": "string",
        "notes": "string",
        "reminders": [],
        "created_by": "string",
        "created_at": "string",
        "updated_at": "string"
      }
    ],
    "total": 0
  }
}
```

**Note:** Volunteers only see their own patients.

### Get Patient

**GET** `/patients/:id`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "patient": { ... }
  }
}
```

### Create Patient

**POST** `/patients`

**Request Body:**
```json
{
  "name": "string (required)",
  "phone": "string (required, Indonesian format)",
  "email": "string (optional)",
  "notes": "string (optional)"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "patient": { ... }
  }
}
```

### Update Patient

**PUT** `/patients/:id`

**Request Body:**
```json
{
  "name": "string",
  "phone": "string",
  "email": "string",
  "notes": "string"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "patient": { ... }
  }
}
```

### Delete Patient

**DELETE** `/patients/:id`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Patient deleted successfully"
}
```

---

## Reminder Endpoints

### List Reminders

**GET** `/patients/:patient_id/reminders`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "reminders": [
      {
        "id": "string",
        "title": "string",
        "description": "string",
        "due_date": "string",
        "priority": "low|medium|high",
        "completed": false,
        "attachments": [],
        "delivery_status": "pending",
        "created_at": "string"
      }
    ]
  }
}
```

### Create Reminder

**POST** `/patients/:patient_id/reminders`

**Request Body:**
```json
{
  "title": "string (required)",
  "description": "string (optional)",
  "due_date": "string (required, RFC3339 or 2006-01-02T15:04)",
  "priority": "string (optional, default 'medium')",
  "attachments": [
    {
      "type": "article|video",
      "id": "string"
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "reminder": { ... }
  }
}
```

### Update Reminder

**PUT** `/patients/:patient_id/reminders/:reminder_id`

**Request Body:**
```json
{
  "title": "string",
  "description": "string",
  "due_date": "string",
  "priority": "string"
}
```

### Toggle Completed

**POST** `/patients/:patient_id/reminders/:reminder_id/toggle`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "reminder": { ... }
  }
}
```

### Send Reminder via WhatsApp

**POST** `/patients/:patient_id/reminders/:reminder_id/send`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "message": "Reminder queued for sending",
    "delivery_status": "scheduled"
  }
}
```

### Get Delivery Status

**GET** `/reminders/:id/status`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "string",
    "delivery_status": "pending|scheduled|sending|sent|delivered|read|failed|retrying|cancelled",
    "message_sent_at": "string",
    "delivered_at": "string",
    "read_at": "string",
    "retry_count": 0,
    "delivery_error_message": "string"
  }
}
```

### Retry Failed Reminder

**POST** `/reminders/:id/retry`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Reminder queued for retry"
}
```

### Cancel Reminder

**POST** `/reminders/:id/cancel`

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Reminder cancelled"
}
```

---

## Content Management Endpoints

### List Categories

**GET** `/categories`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "categories": [
      {
        "id": "string",
        "name": "string",
        "type": "article|video",
        "created_at": "string"
      }
    ]
  }
}
```

### Get Categories by Type

**GET** `/categories/:type`

- `type`: `article` or `video`

### Create Category

**POST** `/categories` (Admin+)

**Request Body:**
```json
{
  "name": "string (required)",
  "type": "article|video (required)"
}
```

### List Articles

**GET** `/articles`

**Query Parameters:**
- `limit`, `offset`
- `category_id` (optional)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "articles": [
      {
        "id": "string",
        "title": "string",
        "slug": "string",
        "excerpt": "string",
        "hero_images": {
          "16x9": "url",
          "1x1": "url",
          "4x3": "url"
        },
        "status": "draft|published",
        "view_count": 0,
        "attachment_count": 0,
        "created_at": "string",
        "published_at": "string"
      }
    ]
  }
}
```

### Get Article by Slug

**GET** `/articles/:slug`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "article": {
      "id": "string",
      "title": "string",
      "slug": "string",
      "content": "html_string",
      "excerpt": "string",
      "category_id": "string",
      "hero_images": { ... },
      "status": "string",
      "view_count": 0,
      "attachment_count": 0,
      "created_at": "string",
      "published_at": "string",
      "updated_at": "string"
    }
  }
}
```

### Create Article

**POST** `/articles` (Admin+)

**Request Body:**
```json
{
  "title": "string (required)",
  "content": "html_string (required)",
  "excerpt": "string (optional)",
  "category_id": "string (optional)",
  "status": "draft|published (optional, default draft)"
}
```

### Update Article

**PUT** `/articles/:id` (Admin+)

**Request Body:** Same as create

### Delete Article

**DELETE** `/articles/:id` (Admin+)

### List Videos

**GET** `/videos`

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "videos": [
      {
        "id": "string",
        "youtube_url": "string",
        "youtube_id": "string",
        "title": "string",
        "description": "string",
        "thumbnail_url": "string",
        "duration": "string",
        "category_id": "string",
        "status": "published",
        "attachment_count": 0,
        "created_at": "string"
      }
    ]
  }
}
```

### Create Video

**POST** `/videos` (Admin+)

**Request Body:**
```json
{
  "youtube_url": "string (required)",
  "category_id": "string (optional)"
}
```

### Delete Video

**DELETE** `/videos/:id` (Admin+)

### List All Content

**GET** `/content`

Returns combined articles and videos.

### Get Popular Content

**GET** `/content/popular`

Returns content sorted by attachment count.

### Upload Image

**POST** `/upload/image` (Admin+)

**Content-Type:** `multipart/form-data`

**Field:** `image` (file upload)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "urls": {
      "16x9": "url",
      "1x1": "url",
      "4x3": "url"
    }
  }
}
```

---

## User Management Endpoints

### List Users

**GET** `/users` (Superadmin only)

### Get User

**GET** `/users/:id`

### Update User Role

**PUT** `/users/:id/role` (Superadmin only)

**Request Body:**
```json
{
  "role": "admin|volunteer"
}
```

### Delete User

**DELETE** `/users/:id` (Superadmin only)

---

## Analytics Endpoints

### Dashboard Stats

**GET** `/dashboard/stats` (Admin+)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "total_patients": 0,
    "total_reminders": 0,
    "completed_reminders": 0,
    "pending_reminders": 0,
    "failed_deliveries": 0,
    "content_counts": {
      "articles": 0,
      "videos": 0
    }
  }
}
```

### Content Analytics

**GET** `/analytics/content` (Admin+)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "content": [
      {
        "id": "string",
        "type": "article|video",
        "title": "string",
        "attachment_count": 0
      }
    ]
  }
}
```

### Sync Attachment Counts

**POST** `/analytics/content/sync` (Admin+)

### Delivery Analytics

**GET** `/analytics/delivery` (Admin+)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "total_sent": 0,
    "delivered": 0,
    "read": 0,
    "failed": 0,
    "success_rate": 0.0
  }
}
```

### Failed Deliveries List

**GET** `/analytics/failed-deliveries` (Admin+)

**Query Parameters:**
- `limit`, `offset`
- `date_from`, `date_to` (optional)

### Export Failed Deliveries

**GET** `/analytics/failed-deliveries/export` (Admin+)

Returns CSV file.

### Failed Delivery Detail

**GET** `/analytics/failed-deliveries/:id` (Admin+)

---

## Health & SSE Endpoints

### Basic Health Check

**GET** `/health`

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "string"
}
```

### Detailed Health Check

**GET** `/health/detailed` (Admin+)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "status": "healthy|degraded",
    "components": {
      "gowa": "connected|disconnected",
      "circuit_breaker": "closed|open"
    }
  }
}
```

### SSE Delivery Status Stream

**GET** `/sse/delivery-status?token=<jwt_token>`

**Response:** SSE stream (text/event-stream)

---

## Webhook Endpoints

### GOWA Webhook

**POST** `/webhook/gowa`

**Headers:**
- `X-Webhook-Signature`: HMAC-SHA256 signature

**Request Body:**
```json
{
  "messageId": "string",
  "status": "delivered|read|failed",
  "timestamp": "string"
}
```

---

## HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | OK |
| 201 | Created |
| 400 | Bad Request |
| 401 | Unauthorized (missing/invalid token) |
| 403 | Forbidden (insufficient permissions) |
| 404 | Not Found |
| 500 | Internal Server Error |
| 503 | Service Unavailable (circuit breaker open) |

---

_Generated using BMAD Method `document-project` workflow_
