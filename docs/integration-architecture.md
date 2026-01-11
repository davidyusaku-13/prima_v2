# Integration Architecture - PRIMA

**Generated:** January 11, 2026 (Updated)
**Project:** Healthcare Volunteer Dashboard
**Architecture Type:** Multi-Part (Backend API + Frontend SPA + External Services)
**Scan Type:** Exhaustive Rescan

---

## System Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         PRIMA System                            │
│                                                                 │
│  ┌───────────────┐         ┌──────────────────┐               │
│  │   Frontend    │◄───────►│     Backend      │               │
│  │  Svelte 5 +   │  REST   │   Go/Gin API     │               │
│  │     Vite      │   API   │   (Port 8080)    │               │
│  │  (Port 5173)  │         │                  │               │
│  └───────┬───────┘         └────────┬─────────┘               │
│          │                          │                          │
│          │ SSE                      │ HTTP + Circuit Breaker   │
│          │ (Real-time)              │                          │
│          │                          ▼                          │
│          │                  ┌───────────────┐                 │
│          └─────────────────►│     GOWA      │                 │
│                              │   WhatsApp    │                 │
│                              │   Gateway     │                 │
│                              │  (Port 3000)  │                 │
│                              └───────┬───────┘                 │
│                                      │                          │
│                                      │ Webhook (HMAC)          │
│                                      │                          │
│                                      ▼                          │
│                              ┌───────────────┐                 │
│                              │   Backend     │                 │
│                              │   Webhook     │                 │
│                              │   Endpoint    │                 │
│                              └───────────────┘                 │
└─────────────────────────────────────────────────────────────────┘
                                     │
                                     │ API Request
                                     ▼
                              ┌───────────────┐
                              │  noembed.com  │
                              │   YouTube     │
                              │   Metadata    │
                              └───────────────┘
```

---

## Integration Points

### 1. Frontend ↔ Backend (REST API)

**Protocol:** HTTP/HTTPS  
**Authentication:** JWT Bearer Token  
**Data Format:** JSON  
**CORS:** Enabled for `http://localhost:5173` (dev)

#### Communication Pattern

**Request:**

```http
GET /api/patients HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json
```

**Response:**

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "patients": [...]
}
```

#### API Categories

| Category      | Base Path                      | Authentication      | Public Endpoints |
| ------------- | ------------------------------ | ------------------- | ---------------- |
| **Auth**      | `/api/auth`                    | No (login/register) | 2                |
| **Patients**  | `/api/patients`                | Yes                 | 0                |
| **Reminders** | `/api/patients/:id/reminders`  | Yes                 | 0                |
| **Content**   | `/api/articles`, `/api/videos` | Partial             | 9 public         |
| **Analytics** | `/api/analytics`               | Yes (admin)         | 0                |
| **Health**    | `/api/health`                  | Partial             | 1 public         |
| **Config**    | `/api/config`                  | No                  | 2                |

**Total Endpoints:** 54 (see [API Contracts](./api-contracts-backend.md))

#### Error Handling

**Backend:**

```go
c.JSON(http.StatusBadRequest, gin.H{
    "error": "Invalid phone number",
    "code": "INVALID_PHONE"
})
```

**Frontend:**

```javascript
try {
  const response = await api.get("/api/patients");
  return response.data;
} catch (error) {
  if (error.response?.status === 401) {
    // Redirect to login
  }
  showNotification(error.message, "error");
}
```

---

### 2. Frontend ↔ Backend (Server-Sent Events)

**Protocol:** HTTP with `text/event-stream`  
**Authentication:** JWT in query parameter  
**Usage:** Real-time delivery status updates

#### Connection Flow

```javascript
// Frontend
const token = localStorage.getItem("token");
const eventSource = new EventSource(`/api/sse/delivery-status?token=${token}`);

eventSource.addEventListener("connection.established", (e) => {
  console.log("SSE connected:", JSON.parse(e.data));
});

eventSource.addEventListener("delivery.status.updated", (e) => {
  const { reminder_id, status } = JSON.parse(e.data);
  updateReminderUI(reminder_id, status);
});

eventSource.addEventListener("delivery.failed", (e) => {
  const { reminder_id, error } = JSON.parse(e.data);
  showErrorNotification(error);
});

eventSource.onerror = () => {
  // Reconnect logic
  setTimeout(() => eventSource.close(), 5000);
};
```

```go
// Backend (handlers/sse.go)
func (h *SSEHandler) HandleDeliveryStatusSSE(c *gin.Context) {
    // Validate JWT from query param
    // Add client to connection pool
    // Send connection.established event

    for {
        select {
        case event := <-h.broadcast:
            // Send event to client
        case <-c.Request.Context().Done():
            // Client disconnected
            return
        }
    }
}
```

#### Event Types

1. **connection.established**

   - Sent immediately after SSE connection
   - Confirms successful authentication

2. **delivery.status.updated**

   - Triggered by GOWA webhook or manual send
   - Updates reminder delivery status

3. **delivery.failed**
   - Triggered when delivery fails after retries
   - Contains error details

---

### 3. Backend ↔ GOWA (WhatsApp Gateway)

**Protocol:** HTTP/HTTPS  
**Authentication:** Basic Auth (username:password)  
**Circuit Breaker:** Enabled (5 failures → 5min cooldown)  
**Retry Strategy:** Exponential backoff (1s, 5s, 30s, 2m, 10m)

#### Send Message Flow

```go
// services/gowa.go
func (c *GOWAClient) SendMessage(phone, message string) (string, error) {
    // Check circuit breaker state
    if c.circuitBreaker.IsOpen() {
        return "", ErrCircuitBreakerOpen
    }

    // Prepare request
    payload := map[string]string{
        "phone":   phone,
        "message": message,
    }

    // Send HTTP POST with Basic Auth
    resp, err := c.httpClient.Post(
        c.config.GOWA.Endpoint + "/send/message",
        payload,
    )

    if err != nil {
        c.circuitBreaker.RecordFailure()
        return "", err
    }

    c.circuitBreaker.RecordSuccess()
    return resp.MessageID, nil
}
```

#### Request

```http
POST /send/message HTTP/1.1
Host: localhost:3000
Authorization: Basic dXNlcjpwYXNzd29yZA==
Content-Type: application/json

{
  "phone": "628123456789",
  "message": "🔔 Pengingat: Minum Obat Pagi\n\nWaktunya minum obat rutin Anda.\n\n📚 Artikel: Pentingnya Konsistensi Minum Obat\nhttps://prima.app/articles/konsistensi-obat"
}
```

#### Response

```json
{
  "status": "success",
  "message_id": "gowa-msg-550e8400-e29b-41d4-a716-446655440000"
}
```

#### Circuit Breaker States

```
┌─────────────┐
│   CLOSED    │ ◄─── Normal operation
│  (Normal)   │
└──────┬──────┘
       │ 5 consecutive failures
       ▼
┌─────────────┐
│    OPEN     │ ◄─── No requests allowed
│ (Failing)   │
└──────┬──────┘
       │ After 5 minutes
       ▼
┌─────────────┐
│  HALF-OPEN  │ ◄─── Test with 1 request
│  (Testing)  │
└──────┬──────┘
       │ Success → CLOSED
       │ Failure → OPEN
```

---

### 4. GOWA ↔ Backend (Webhook Callbacks)

**Protocol:** HTTP POST  
**Authentication:** HMAC-SHA256 signature  
**Header:** `X-Webhook-Signature`  
**Endpoint:** `POST /api/webhook/gowa`

#### Webhook Flow

```
┌──────────┐                           ┌──────────┐
│   GOWA   │                           │ Backend  │
└────┬─────┘                           └─────┬────┘
     │                                       │
     │ 1. Message delivered/read             │
     ├──────────────────────────────────────►│
     │   POST /api/webhook/gowa              │
     │   X-Webhook-Signature: hmac-sha256    │
     │   Body: {event, message}              │
     │                                       │
     │                                  2. Validate HMAC
     │                                       │
     │                                  3. Update reminder
     │                                       │
     │                                  4. Broadcast SSE
     │                                       │
     │◄──────────────────────────────────────│
     │   200 OK                              │
     │                                       │
```

#### Webhook Payload

```json
{
  "event": "message.ack",
  "message": {
    "id": "gowa-msg-550e8400-e29b-41d4-a716-446655440000",
    "status": "delivered",
    "timestamp": "2026-01-02T10:00:10Z"
  }
}
```

#### HMAC Validation

```go
// utils/hmac.go
func ValidateHMAC(secret, payload, signature string) bool {
    mac := hmac.New(sha256.New, []byte(secret))
    mac.Write([]byte(payload))
    expectedMAC := hex.EncodeToString(mac.Sum(nil))
    return hmac.Equal([]byte(signature), []byte(expectedMAC))
}

// handlers/webhook.go
func (h *WebhookHandler) HandleGOWAWebhook(c *gin.Context) {
    signature := c.GetHeader("X-Webhook-Signature")
    if signature == "" {
        c.JSON(401, gin.H{"error": "Missing signature"})
        return
    }

    body, _ := io.ReadAll(c.Request.Body)
    if !utils.ValidateHMAC(h.config.GOWA.WebhookSecret, string(body), signature) {
        c.JSON(401, gin.H{"error": "Invalid signature"})
        return
    }

    // Process webhook
}
```

#### Status Mapping

| GOWA Status | PRIMA Status | Description                     |
| ----------- | ------------ | ------------------------------- |
| `sent`      | `sent`       | Message sent to WhatsApp server |
| `delivered` | `delivered`  | Delivered to recipient's phone  |
| `read`      | `read`       | Opened by recipient             |
| `failed`    | `failed`     | Delivery failed                 |

---

### 5. Backend ↔ YouTube Metadata API

**Protocol:** HTTPS  
**Service:** noembed.com (free YouTube oEmbed proxy)  
**Authentication:** None  
**Rate Limit:** Reasonable use

#### Metadata Fetch Flow

```
┌──────────┐                           ┌──────────────┐
│  Admin   │                           │   Backend    │
└────┬─────┘                           └──────┬───────┘
     │                                        │
     │ POST /api/videos                       │
     │ {youtube_url}                          │
     ├───────────────────────────────────────►│
     │                                        │
     │                                   1. Extract video ID
     │                                        │
     │                                   2. Fetch metadata
     │                                        │
     │                                        ▼
     │                                ┌──────────────┐
     │                                │ noembed.com  │
     │                                └──────┬───────┘
     │                                        │
     │                                   3. Parse response
     │                                        │
     │                                   4. Save video
     │                                        │
     │◄───────────────────────────────────────│
     │   201 Created                          │
     │   {video object}                       │
```

#### API Request

```http
GET /embed?url=https://www.youtube.com/watch?v=VIDEO_ID HTTP/1.1
Host: noembed.com
```

#### API Response

```json
{
  "title": "Cara Cuci Tangan yang Benar - WHO",
  "author_name": "Kementerian Kesehatan RI",
  "thumbnail_url": "https://i.ytimg.com/vi/VIDEO_ID/maxresdefault.jpg",
  "width": 1920,
  "height": 1080,
  "html": "<iframe...></iframe>"
}
```

#### Error Handling

```go
// utils/youtube.go
func FetchYouTubeMetadata(url string) (*models.YouTubeMetadata, error) {
    resp, err := http.Get("https://noembed.com/embed?url=" + url)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch metadata: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("noembed returned status %d", resp.StatusCode)
    }

    var metadata models.YouTubeMetadata
    if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
        return nil, fmt.Errorf("failed to parse metadata: %w", err)
    }

    return &metadata, nil
}
```

---

## Data Flow Diagrams

### Reminder Send Flow

```
┌──────────┐     ┌──────────┐     ┌──────────────┐     ┌────────┐     ┌──────────┐
│Volunteer │     │ Frontend │     │   Backend    │     │  GOWA  │     │ Patient  │
└────┬─────┘     └────┬─────┘     └──────┬───────┘     └───┬────┘     └────┬─────┘
     │                │                   │                 │               │
     │ 1. Click "Send"│                   │                 │               │
     ├───────────────►│                   │                 │               │
     │                │                   │                 │               │
     │                │ 2. POST /api/.../ │                 │               │
     │                │    reminders/:id/ │                 │               │
     │                │    send           │                 │               │
     │                ├──────────────────►│                 │               │
     │                │                   │                 │               │
     │                │              3. Check quiet hours   │               │
     │                │                   │                 │               │
     │                │              4. Format message      │               │
     │                │                   │                 │               │
     │                │              5. Check circuit breaker               │
     │                │                   │                 │               │
     │                │                   │ 6. POST /send/  │               │
     │                │                   │    message      │               │
     │                │                   ├────────────────►│               │
     │                │                   │                 │               │
     │                │                   │                 │ 7. Send to    │
     │                │                   │                 │    WhatsApp   │
     │                │                   │                 ├──────────────►│
     │                │                   │                 │               │
     │                │                   │ 8. Return       │               │
     │                │                   │    message_id   │               │
     │                │                   │◄────────────────│               │
     │                │                   │                 │               │
     │                │              9. Update status=sent  │               │
     │                │                   │                 │               │
     │                │             10. Broadcast SSE       │               │
     │                │◄─────────────────┼─────────────────┘               │
     │                │  (SSE event)      │                                 │
     │                │                   │                                 │
     │◄───────────────│ 11. UI update     │                                 │
     │  "Sent"        │                   │                                 │
     │                │                   │                                 │
     │                │                   │ 12. Webhook:    │               │
     │                │                   │     delivered   │               │
     │                │                   │◄────────────────┼───────────────┘
     │                │                   │                 │  (WhatsApp ack)
     │                │              13. Update status      │
     │                │                   │                 │
     │                │             14. Broadcast SSE       │
     │                │◄─────────────────┘                 │
     │                │  (delivery.status.updated)          │
     │                │                                     │
     │◄───────────────│ 15. UI update "Delivered"          │
     │                │                                     │
```

---

### CMS Article Creation Flow

```
┌──────┐     ┌──────────┐     ┌──────────────┐
│Admin │     │ Frontend │     │   Backend    │
└──┬───┘     └────┬─────┘     └──────┬───────┘
   │              │                   │
   │ 1. Upload    │                   │
   │    image     │                   │
   ├─────────────►│                   │
   │              │ 2. POST /api/     │
   │              │    upload/image   │
   │              ├──────────────────►│
   │              │                   │
   │              │              3. Resize to 3
   │              │                 aspect ratios
   │              │                   │
   │              │              4. Save to uploads/
   │              │                   │
   │              │ 5. Return URLs    │
   │              │◄──────────────────│
   │              │                   │
   │ 2. Create    │                   │
   │    article   │                   │
   ├─────────────►│                   │
   │              │ 6. POST /api/     │
   │              │    articles       │
   │              │    {title, content│
   │              │     hero_images}  │
   │              ├──────────────────►│
   │              │                   │
   │              │              7. Generate slug
   │              │                   │
   │              │              8. Save to articles.json
   │              │                   │
   │              │ 9. Return article │
   │              │◄──────────────────│
   │              │                   │
   │◄─────────────│ 10. Show success  │
   │  "Published" │                   │
```

---

## Security Measures

### 1. Authentication

**JWT Token:**

- Algorithm: HS256
- Expiry: 7 days
- Secret: Stored in `data/jwt_secret.txt`
- Claims: `userId`, `username`, `role`

**Password:**

- Hash: SHA256
- Encoding: Base64
- Storage: `users.json`

### 2. Authorization

**Role-Based Access Control (RBAC):**

| Role         | Permissions                            |
| ------------ | -------------------------------------- |
| `superadmin` | All operations + user management       |
| `admin`      | CMS management + analytics             |
| `volunteer`  | Patient management (own patients only) |

**Enforcement:**

```go
func requireRole(allowedRoles ...Role) gin.HandlerFunc {
    return func(c *gin.Context) {
        claims := c.MustGet("claims").(*Claims)

        for _, role := range allowedRoles {
            if claims.Role == role {
                c.Next()
                return
            }
        }

        c.JSON(403, gin.H{"error": "insufficient permissions"})
        c.Abort()
    }
}
```

### 3. CORS

**Configuration:**

```yaml
server:
  cors_origin: "http://localhost:5173"
```

**Allowed:**

- Methods: GET, POST, PUT, DELETE, OPTIONS
- Headers: Origin, Content-Type, Authorization
- Credentials: true

### 4. HMAC Webhook Validation

**Algorithm:** HMAC-SHA256  
**Header:** `X-Webhook-Signature`  
**Secret:** From `config.yaml` (`gowa.webhook_secret`)

**Process:**

1. Read request body
2. Compute HMAC with secret
3. Compare with signature header
4. Reject if mismatch

### 5. Data Masking

**Phone Numbers:**

```
628123456789 → 6281234***789
```

**Emails:**

```
john.doe@example.com → j***@example.com
```

**Usage:** Logs, admin analytics views

---

## Performance Considerations

### 1. Concurrent Request Handling

**Backend:**

- Gin default: Concurrent HTTP handling
- `sync.RWMutex` for data stores
- Multiple readers, single writer

**Frontend:**

- Async API calls with `await`
- Loading states during requests
- Debouncing for search inputs

### 2. SSE Connection Management

**Backend:**

```go
type SSEHandler struct {
    clients   map[string]chan Event
    mu        sync.RWMutex
    broadcast chan Event
}
```

**Scaling:**

- Max clients: Unlimited (limited by OS file descriptors)
- Memory per client: ~8KB (channel buffer)
- Reconnection: Automatic on frontend

### 3. Circuit Breaker Benefits

**Without Circuit Breaker:**

- Every failed request waits for timeout (30s)
- Cascading failures overwhelm GOWA
- Backend becomes unresponsive

**With Circuit Breaker:**

- After 5 failures, stop sending for 5 minutes
- Fast-fail for pending reminders (scheduled for retry)
- Backend remains responsive

---

## Monitoring & Observability

### Health Endpoints

**Basic Health Check (Public):**

```
GET /api/health
```

**Detailed Health (Admin):**

```
GET /api/health/detailed
```

**Response:**

```json
{
  "data": {
    "status": "ok",
    "timestamp": "2026-01-02T10:00:00Z",
    "gowa": {
      "connected": true,
      "last_ping": "2026-01-02T09:59:50Z"
    },
    "circuit_breaker": {
      "state": "closed",
      "failure_count": 0
    },
    "queue": {
      "total": 15,
      "scheduled": 8,
      "retrying": 3
    }
  }
}
```

### Logging

**Structured Logging (slog):**

```go
logger.Info("Reminder sent",
    "reminder_id", reminderID,
    "patient_id", patientID,
    "phone", utils.MaskPhone(phone),
    "status", "sent",
)
```

**Log Levels:**

- `DEBUG` - Development details
- `INFO` - Normal operations
- `WARN` - Recoverable errors
- `ERROR` - Critical failures

---

## Deployment Architecture

### Development

```
┌────────────────────────────────────────────┐
│         Local Machine                      │
│                                            │
│  Backend (go run)         Frontend (Vite)  │
│  Port 8080                Port 5173        │
│                                            │
│  GOWA (Docker)                             │
│  Port 3000                                 │
└────────────────────────────────────────────┘
```

### Production (Recommended)

```
┌──────────────────────────────────────────────────┐
│                  Nginx/Caddy                     │
│           (Reverse Proxy + SSL)                  │
└───────────────────┬──────────────────────────────┘
                    │
        ┌───────────┴───────────┐
        │                       │
        ▼                       ▼
┌───────────────┐       ┌───────────────┐
│   Backend     │       │   Frontend    │
│   (Binary)    │       │   (Static)    │
│   Port 8080   │       │   /var/www    │
└───────┬───────┘       └───────────────┘
        │
        │ HTTP
        ▼
┌───────────────┐
│     GOWA      │
│   (Docker)    │
│   Port 3000   │
└───────────────┘
```

**Nginx Configuration:**

```nginx
server {
    listen 80;
    server_name prima.example.com;

    # Frontend (static files)
    location / {
        root /var/www/prima/dist;
        try_files $uri $uri/ /index.html;
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # SSE (disable buffering)
    location /api/sse/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Connection '';
        proxy_http_version 1.1;
        chunked_transfer_encoding off;
        proxy_buffering off;
        proxy_cache off;
    }

    # Backend uploads
    location /uploads/ {
        alias /opt/prima/backend/uploads/;
    }
}
```

---

**Next:** See [Architecture - Backend](./architecture-backend.md) and [Architecture - Frontend](./architecture-frontend.md) for detailed part-specific designs.
