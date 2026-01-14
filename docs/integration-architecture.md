# PRIMA Integration Architecture

**Date:** 2026-01-13

## Overview

PRIMA is a multi-part project with two main components:
- **Backend:** Go/Gin API server (port 8080)
- **Frontend:** Svelte 5 + Vite SPA (port 5173)

This document describes how these parts integrate with each other and with external services.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Browser                           │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              PRIMA Frontend (Svelte 5 + Vite)           │   │
│  │                                                          │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐  │   │
│  │  │  Stores  │  │ Services │  │   Components/Views   │  │   │
│  │  │ auth     │  │ SSE      │  │   Patients, CMS      │  │   │
│  │  │ delivery │  └────┬─────┘  └──────────────────────┘  │   │
│  │  └──────────┘       │                                  │   │
│  └────────────────────┼──────────────────────────────────┘   │
│                       │                                      │
└───────────────────────┼──────────────────────────────────────┘
                        │ HTTP/REST + SSE
                        │
┌───────────────────────┼──────────────────────────────────────┐
│                       ▼                                      │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │              PRIMA Backend (Go + Gin)                    │   │
│  │                                                          │   │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐  │   │
│  │  │ Handlers │  │ Services │  │   Data Stores        │  │   │
│  │  │ API      │  │ GOWA     │  │   (JSON Files)       │  │   │
│  │  │ SSE      │  │ Scheduler│  │                     │  │   │
│  │  └──────────┘  └──────────┘  └──────────────────────┘  │   │
│  │                                                          │   │
│  └─────────────────────────────────────────────────────────┘   │
│                       │                                      │
│              ┌────────┼────────┐                             │
│              ▼        ▼        ▼                             │
│       ┌──────────┐ ┌───────┐ ┌──────────────────────────┐  │
│       │   GOWA   │ │ JWT   │ │      noembed.com         │  │
│       │ (WA)     │ │ Auth  │ │    (YouTube)             │  │
│       └──────────┘ └───────┘ └──────────────────────────┘  │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Frontend → Backend Integration

### REST API Communication

| Aspect | Details |
|--------|---------|
| **Protocol** | HTTP/HTTPS |
| **Base URL** | `http://localhost:8080/api` |
| **Format** | JSON |
| **Authentication** | JWT Bearer token in `Authorization` header |

### API Call Flow

```
Frontend Component
       │
       ▼ (fetch/axios with JWT)
┌──────────────────┐
│  api.js wrapper  │  ← Adds auth header, handles errors
└────────┬─────────┘
         │
         ▼ (REST call)
┌──────────────────┐
│  Gin Handlers    │  ← Routes to appropriate handler
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│   Services/      │  ← Business logic
│   Models         │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│  JSON Files      │  ← Persistence
└──────────────────┘
```

### Authentication Flow

```
1. User submits login form
2. Frontend POST /api/auth/login with credentials
3. Backend validates, returns JWT token
4. Frontend stores token in localStorage
5. Subsequent requests include: Authorization: Bearer <token>
6. Backend middleware validates token, extracts user role
7. Role-based access control applied
```

### Real-Time Updates (SSE)

```
Frontend                              Backend
   │                                     │
   │── GET /api/sse/delivery-status ───►│
   │    ?token=<jwt_token>              │
   │◄───────────────────────────────────│ 200 OK (SSE stream)
   │                                     │
   │              ┌─────────────────►   │ Webhook from GOWA
   │              │                     │
   │◄─────────────┤ "delivery.status.updated" event
   │              │ { reminder_id, status }
   │              │
   │   (deliveryStore updates UI)       │
```

**SSE Events:**

| Event | Payload | Purpose |
|-------|---------|---------|
| `connection.established` | `{ event: 'connected' }` | SSE connected |
| `connection.status` | `{ status: 'connected'/'disconnected' }` | Connection status |
| `delivery.status.updated` | `{ reminder_id, status, ... }` | Real-time delivery update |
| `delivery.failed` | `{ reminder_id, error }` | Failed delivery notification |

## Backend → GOWA Integration

### WhatsApp Message Sending

```
Backend Handler                         GOWA Server
      │                                      │
      │── POST /send/message ──────────────►│
      │  Headers:                            │
      │  - Authorization: Basic <creds>      │
      │  Body:                               │
      │  {                                   │
      │    "phone": "628123456789",          │
      │    "message": "Reminder: ..."        │
      │  }                                   │
      │◄─────────────────────────────────────│ 200 OK
      │  Body:                               │
      │  { "messageId": "abc123" }           │
```

### GOWA Webhook

```
GOWA Server                          Backend
    │                                   │
    │── POST /api/webhook/gowa ────────►│
    │  Headers:                         │
    │  - X-Webhook-Signature: <hmac>    │
    │  Body:                            │
    │  {                                │
    │    "messageId": "abc123",         │
    │    "status": "delivered",         │
    │    "timestamp": "..."             │
    │  }                                │
    │◄───────────────────────────────────│ 200 OK
```

**Webhook Flow:**
1. GOWA validates HMAC signature
2. Backend updates reminder status (sent → delivered/read/failed)
3. Backend broadcasts SSE event to connected frontend clients
4. Frontend deliveryStore updates in real-time

## Backend → External Services

### YouTube Metadata (noembed.com)

```
Backend                         noembed.com
    │                               │
    │── GET /embed?url=<yt_url> ───►│
    │◄───────────────────────────────│ 200 OK
    │  {                             │
    │    "title": "...",             │
    │    "description": "...",       │
    │    "thumbnail_url": "...",     │
    │    "duration": "..."           │
    │  }                             │
```

## Data Flow Diagrams

### Patient Management Flow

```
┌─────────┐     ┌─────────┐     ┌─────────┐     ┌─────────┐
│ Frontend│────▶│  API    │────▶│ Handler │────▶│  Data   │
│ Patient │     │ Layer   │     │         │     │  Store  │
│  View   │◀────│         │◀────│         │◀────│         │
└─────────┘     └─────────┘     └─────────┘     └─────────┘
     │                                                  │
     │◀─ SSE Updates ───────────────────────────────────┘
     │     (delivery status changes)
```

### Reminder Sending Flow

```
User clicks "Send Reminder"
         │
         ▼
┌─────────────────┐
│ 1. Validate     │
│ 2. Check hours  │
│ 3. Check phone  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐     ┌─────────────┐
│ POST /send      │────▶│ Validate    │
│                  │     │ phone format│
└────────┬────────┘     └──────┬──────┘
         │                     │
         ▼                     ▼
┌─────────────────┐     ┌─────────────┐
│ 4. Send to GOWA │     │ Normalize   │
│                 │     │ 628xxxxx    │
└────────┬────────┘     └─────────────┘
         │
         ▼
┌─────────────────┐
│ 5. Update status│
│    sent/failed  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ 6. Broadcast    │
│    SSE event    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Frontend updates│
│ delivery badge  │
└─────────────────┘
```

### CMS Content Flow

```
Admin creates article
         │
         ▼
┌─────────────────┐     ┌─────────────┐
│ POST /articles  │────▶│ Validate    │
│                 │     │ data        │
└────────┬────────┘     └──────┬──────┘
         │                     │
         ▼                     ▼
┌─────────────────┐     ┌─────────────┐
│ Upload images   │     │ Create      │
│ (if any)        │     │ slug        │
└────────┬────────┘     └─────────────┘
         │
         ▼
┌─────────────────┐
│ Save to JSON    │
│ (articles.json) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Return created  │
│ article         │
└─────────────────┘
```

## Integration Points Summary

| From | To | Protocol | Data | Auth |
|------|-----|----------|------|------|
| Frontend | Backend API | HTTP REST | JSON | JWT Bearer |
| Frontend | Backend SSE | SSE | Event stream | JWT query param |
| Backend | GOWA | HTTP REST | JSON | Basic Auth |
| GOWA | Backend | HTTP webhook | JSON | HMAC signature |
| Backend | noembed.com | HTTP REST | JSON | None (public) |

## Security Considerations

| Aspect | Implementation |
|--------|----------------|
| API Auth | JWT tokens (7-day expiry) |
| SSE Auth | Token in query parameter |
| Webhook Auth | HMAC-SHA256 signature validation |
| CORS | Configurable origin whitelist |
| PII Masking | Phone numbers masked in logs |
| Password Storage | SHA256 hashing |

## Port Assignments

| Service | Port | Protocol |
|---------|------|----------|
| Backend API | 8080 | HTTP |
| Frontend Dev | 5173 | HTTP |
| GOWA Server | 3000 | HTTP |

---

_Generated using BMAD Method `document-project` workflow_
