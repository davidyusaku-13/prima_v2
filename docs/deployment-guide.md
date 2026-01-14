# PRIMA Deployment Guide

**Date:** 2026-01-13

## Overview

This guide covers deployment options for the PRIMA Healthcare Volunteer Dashboard. The application consists of a Go backend API and a Svelte 5 + Vite frontend SPA.

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      Load Balancer / Proxy                   │
│                     (nginx, Caddy, etc.)                     │
└───────────────────────┬─────────────────────────────────────┘
                        │
        ┌───────────────┴───────────────┐
        ▼                               ▼
┌─────────────────┐           ┌─────────────────┐
│   Backend       │           │   Frontend      │
│   (Go/Gin)      │           │   (Static)      │
│   Port: 8080    │           │   Port: 80/443  │
└─────────────────┘           └─────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────────────────┐
│                    GOWA Server (External)                    │
│                    Port: 3000                                │
└─────────────────────────────────────────────────────────────┘
```

## Prerequisites

### System Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| CPU | 1 core | 2+ cores |
| RAM | 512 MB | 1 GB+ |
| Storage | 1 GB | 5 GB+ |
| Go | 1.25+ | 1.25+ |
| Node/Bun | Latest | Latest |

### Required External Services

1. **GOWA Server** - WhatsApp gateway service
   - Port: 3000
   - Basic authentication required

## Deployment Options

### Option 1: Docker Compose (Recommended)

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - GOWA_ENDPOINT=http://gowa:3000
      - GOWA_USERNAME=admin
      - GOWA_PASSWORD=password123
    volumes:
      - backend_data:/app/data
    depends_on:
      - gowa
    restart: unless-stopped

  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped

  gowa:
    image: davidyusaku13/gowa:latest
    ports:
      - "3000:3000"
    restart: unless-stopped

volumes:
  backend_data:
```

**Build and run:**

```bash
docker-compose up -d --build
```

### Option 2: Manual Deployment

#### Backend Deployment

```bash
# Build binary
cd backend
go build -o prima-backend main.go

# Create data directory
mkdir -p data uploads

# Copy configuration
cp config.example.yaml config.yaml
# Edit config.yaml for production settings

# Run with systemd (Linux)
# Create /etc/systemd/system/prima-backend.service
```

**Systemd service file:**

```ini
[Unit]
Description=PRIMA Backend API
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/opt/prima
ExecStart=/opt/prima/prima-backend
Restart=on-failure
Environment=GOWA_ENDPOINT=http://localhost:3000

[Install]
WantedBy=multi-user.target
```

#### Frontend Deployment

```bash
# Build production bundle
cd frontend
bun install
bun run build

# Output is in dist/
# Serve with nginx, Caddy, or any static file server
```

**nginx configuration:**

```nginx
server {
    listen 80;
    server_name prima.example.com;
    root /var/www/prima/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## Environment Configuration

### Backend Configuration

Edit `backend/config.yaml`:

```yaml
server:
  port: 8080
  cors_origin: "https://prima.example.com"

gowa:
  endpoint: "http://localhost:3000"
  username: "admin"
  password: "your_secure_password"
  webhook_secret: "your_webhook_secret"

circuit_breaker:
  failure_threshold: 5
  cooldown_minutes: 5

logging:
  level: info
  format: json

disclaimer:
  text: "Peringatan: Informasi ini hanya untuk tujuan edukasi..."

quiet_hours:
  start_hour: 21
  end_hour: 6
  timezone: "Asia/Jakarta"
```

### Frontend Configuration

For production, update the API URL in `frontend/src/lib/utils/api.js`:

```javascript
const API_URL = 'https://api.prima.example.com/api';
```

## Data Persistence

### Backup Strategy

```bash
# Create backup script
#!/bin/bash
BACKUP_DIR="/backups/prima"
DATE=$(date +%Y%m%d_%H%M%S)

# Backup JSON files
tar -czf "$BACKUP_DIR/data_$DATE.tar.gz" /opt/prima/data/

# Keep last 30 days
find "$BACKUP_DIR" -name "*.tar.gz" -mtime +30 -delete
```

### Restore from Backup

```bash
# Restore data files
tar -xzf data_20260101_120000.tar.gz -C /opt/prima/
```

## SSL/TLS Configuration

### Using Let's Encrypt with Certbot

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d prima.example.com

# Auto-renewal
sudo certbot renew --dry-run
```

### Using Caddy (Automatic HTTPS)

```Caddyfile
prima.example.com {
    root * /var/www/prima/dist
    file_server

    handle /api/* {
        reverse_proxy localhost:8080
    }
}
```

## Monitoring

### Health Check Endpoint

```bash
# Basic health check
curl https://api.prima.example.com/api/health

# Detailed health (admin only)
curl -H "Authorization: Bearer <token>" \
  https://api.prima.example.com/api/health/detailed
```

### Log Management

Configure log shipping to your preferred observability platform:

```yaml
# Example: Filebeat configuration
filebeat.inputs:
- type: log
  paths:
    - /var/log/prima/backend.log
  json:
    keys_under_root: true
```

## Security Considerations

1. **JWT Secret**: Generate strong secret in production
   ```bash
   openssl rand -base64 32
   ```

2. **CORS**: Restrict `cors_origin` to your frontend domain

3. **GOWA Credentials**: Use strong passwords

4. **Webhook Secret**: Set secure webhook secret

5. **Rate Limiting**: Implement at proxy level

## Scaling Considerations

### Horizontal Scaling

For high availability, run multiple backend instances:

```yaml
# docker-compose.yml with scaling
services:
  backend:
    deploy:
      replicas: 3
```

**Note:** When scaling, use external data store (Redis, PostgreSQL) instead of JSON files.

### Database Migration Path

Current: JSON file persistence
Future: Migrate to PostgreSQL by:
1. Creating SQL schema
2. Implementing repository pattern
3. Running migration script
4. Updating configuration

## Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| 502 Bad Gateway | Check backend is running on port 8080 |
| CORS errors | Verify `cors_origin` in config.yaml |
| GOWA connection failed | Ensure GOWA server is running |
| SSE not working | Check proxy supports streaming |
| JWT errors | Verify `jwt_secret.txt` exists |

### Logs Location

| Service | Log Location |
|---------|-------------|
| Backend | stdout or `/var/log/prima/backend.log` |
| nginx | `/var/log/nginx/access.log` |
| Docker | `docker-compose logs` |

---

_Generated using BMAD Method `document-project` workflow_
