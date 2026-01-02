# Deployment Guide - PRIMA

**Generated:** January 2, 2026  
**Project:** PRIMA Healthcare Volunteer Dashboard  
**Stack:** Go/Gin Backend + Svelte 5/Vite Frontend

---

## Deployment Overview

PRIMA can be deployed in various configurations:

1. **Single Server** - Backend + Frontend on one machine
2. **Separate Servers** - Backend and Frontend on different machines
3. **Containerized** - Docker/Docker Compose
4. **Cloud** - AWS, GCP, Azure, DigitalOcean

This guide covers production-ready deployment with Nginx reverse proxy.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Backend Deployment](#backend-deployment)
3. [Frontend Deployment](#frontend-deployment)
4. [Nginx Configuration](#nginx-configuration)
5. [SSL/TLS Setup](#ssltls-setup)
6. [Process Management](#process-management)
7. [Database Backup](#database-backup)
8. [Monitoring](#monitoring)
9. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Server Requirements

**Minimum:**

- CPU: 1 core
- RAM: 512 MB
- Disk: 5 GB
- OS: Ubuntu 20.04+ / Debian 11+ / CentOS 8+

**Recommended:**

- CPU: 2 cores
- RAM: 2 GB
- Disk: 20 GB (includes uploads)
- OS: Ubuntu 22.04 LTS

### Software

- **Go:** 1.21+ (for backend)
- **Nginx:** 1.18+ (reverse proxy)
- **Certbot:** (for SSL certificates)
- **Git:** (for deployment)

### Domain & DNS

- Domain name (e.g., `prima.example.com`)
- DNS A record pointing to server IP

---

## Backend Deployment

### 1. Prepare Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

### 2. Clone & Build

```bash
# Create app directory
sudo mkdir -p /opt/prima
sudo chown $USER:$USER /opt/prima

# Clone repository
cd /opt/prima
git clone <repository-url> .

# Build backend
cd backend
go build -o prima-backend main.go
```

### 3. Configure

**Create production config:**

```bash
cd /opt/prima/backend
cp config.example.yaml config.yaml
nano config.yaml
```

**`config.yaml`:**

```yaml
server:
  port: 8080
  cors_origin: "https://prima.example.com" # Frontend URL

gowa:
  endpoint: "http://localhost:3000" # Or external GOWA URL
  username: "admin"
  password: "CHANGE_THIS_PASSWORD"
  webhook_secret: "STRONG_RANDOM_SECRET"
  timeout: 30

scheduler:
  interval: 60

quiet_hours:
  start: 20
  end: 8

logging:
  level: "info"
  format: "json"

disclaimer:
  enabled: true
  version: "1.0"
  text: "PRIMA adalah aplikasi..."
```

**Generate secure secrets:**

```bash
# JWT secret
openssl rand -base64 32 > data/jwt_secret.txt

# GOWA webhook secret (add to config.yaml)
openssl rand -base64 32
```

### 4. Create System User

```bash
# Create dedicated user (no login)
sudo useradd -r -s /bin/false prima

# Set ownership
sudo chown -R prima:prima /opt/prima/backend
```

### 5. Create Systemd Service

**Create service file:**

```bash
sudo nano /etc/systemd/system/prima-backend.service
```

**`/etc/systemd/system/prima-backend.service`:**

```ini
[Unit]
Description=PRIMA Backend API
After=network.target

[Service]
Type=simple
User=prima
Group=prima
WorkingDirectory=/opt/prima/backend
ExecStart=/opt/prima/backend/prima-backend
Restart=on-failure
RestartSec=5s

# Environment variables (optional, can use .env file instead)
Environment="PORT=8080"

# Logging
StandardOutput=append:/var/log/prima/backend.log
StandardError=append:/var/log/prima/backend.error.log

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

**Create log directory:**

```bash
sudo mkdir -p /var/log/prima
sudo chown prima:prima /var/log/prima
```

**Enable and start service:**

```bash
sudo systemctl daemon-reload
sudo systemctl enable prima-backend
sudo systemctl start prima-backend
sudo systemctl status prima-backend
```

**Check logs:**

```bash
sudo journalctl -u prima-backend -f
# Or
sudo tail -f /var/log/prima/backend.log
```

### 6. Verify Backend

```bash
curl http://localhost:8080/api/health
# {"status":"ok","timestamp":"2026-01-02T10:00:00Z"}
```

---

## Frontend Deployment

### 1. Build Frontend Locally

**On development machine:**

```bash
cd frontend

# Update API URL in code or use .env
echo 'VITE_API_URL=https://prima.example.com/api' > .env

# Build
bun run build

# Output: frontend/dist/
```

### 2. Transfer to Server

**Option A: SCP**

```bash
scp -r dist/ user@server:/opt/prima/frontend/
```

**Option B: rsync**

```bash
rsync -avz --delete dist/ user@server:/opt/prima/frontend/dist/
```

**Option C: Build on server** (requires Bun/Node.js on server)

```bash
ssh user@server
cd /opt/prima/frontend
bun install
bun run build
```

### 3. Set Permissions

```bash
sudo chown -R www-data:www-data /opt/prima/frontend/dist
sudo chmod -R 755 /opt/prima/frontend/dist
```

---

## Nginx Configuration

### 1. Install Nginx

```bash
sudo apt install nginx -y
sudo systemctl enable nginx
sudo systemctl start nginx
```

### 2. Create Site Configuration

```bash
sudo nano /etc/nginx/sites-available/prima
```

**`/etc/nginx/sites-available/prima`:**

```nginx
# Redirect HTTP to HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name prima.example.com;

    # Let's Encrypt challenge
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS server
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name prima.example.com;

    # SSL certificates (will be added by Certbot)
    # ssl_certificate /etc/letsencrypt/live/prima.example.com/fullchain.pem;
    # ssl_certificate_key /etc/letsencrypt/live/prima.example.com/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;

    # Frontend (static files)
    root /opt/prima/frontend/dist;
    index index.html;

    # SPA routing - serve index.html for all non-API routes
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Backend API proxy
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # CORS headers (if needed, but backend already handles CORS)
        # add_header Access-Control-Allow-Origin "https://prima.example.com" always;
    }

    # SSE endpoint (disable buffering for real-time updates)
    location /api/sse/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Connection '';
        chunked_transfer_encoding off;
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 86400s;  # 24 hours

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Uploaded files
    location /uploads/ {
        alias /opt/prima/backend/uploads/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Static assets (Vite output)
    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # index.html should not be cached (for SPA updates)
    location = /index.html {
        add_header Cache-Control "no-cache, no-store, must-revalidate";
        add_header Pragma "no-cache";
        add_header Expires "0";
    }

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/json application/xml+rss application/rss+xml font/truetype font/opentype application/vnd.ms-fontobject image/svg+xml;
}
```

### 3. Enable Site

```bash
# Create symlink
sudo ln -s /etc/nginx/sites-available/prima /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

---

## SSL/TLS Setup

### Using Let's Encrypt (Certbot)

**1. Install Certbot:**

```bash
sudo apt install certbot python3-certbot-nginx -y
```

**2. Obtain Certificate:**

```bash
sudo certbot --nginx -d prima.example.com
```

**Follow prompts:**

- Enter email
- Agree to terms
- Choose redirect HTTP to HTTPS (option 2)

**Certbot will automatically:**

- Obtain SSL certificate
- Update Nginx configuration
- Set up auto-renewal

**3. Verify Auto-Renewal:**

```bash
sudo certbot renew --dry-run
```

**4. Check Certificate:**

```bash
sudo certbot certificates
```

### Manual SSL Certificate

If using custom certificate:

```nginx
ssl_certificate /path/to/fullchain.pem;
ssl_certificate_key /path/to/privkey.pem;
```

---

## Process Management

### Systemd (Recommended)

**Check status:**

```bash
sudo systemctl status prima-backend
```

**Start/Stop/Restart:**

```bash
sudo systemctl start prima-backend
sudo systemctl stop prima-backend
sudo systemctl restart prima-backend
```

**View logs:**

```bash
sudo journalctl -u prima-backend -f  # Follow logs
sudo journalctl -u prima-backend --since today  # Today's logs
sudo journalctl -u prima-backend --since "1 hour ago"
```

### Alternative: PM2 (Node.js)

If backend compiled with Node.js wrapper:

```bash
# Install PM2
npm install -g pm2

# Start backend
pm2 start prima-backend --name prima-backend

# Start on boot
pm2 startup
pm2 save

# Monitor
pm2 status
pm2 logs prima-backend
```

---

## Database Backup

### Manual Backup

**Backup data directory:**

```bash
# Create backup
sudo tar -czf /opt/backups/prima-data-$(date +%Y%m%d-%H%M%S).tar.gz \
  /opt/prima/backend/data/

# Create backup directory if not exists
sudo mkdir -p /opt/backups
```

### Automated Backup (Cron)

**Create backup script:**

```bash
sudo nano /opt/prima/scripts/backup.sh
```

**`/opt/prima/scripts/backup.sh`:**

```bash
#!/bin/bash

BACKUP_DIR="/opt/backups"
DATA_DIR="/opt/prima/backend/data"
DATE=$(date +%Y%m%d-%H%M%S)
BACKUP_FILE="$BACKUP_DIR/prima-data-$DATE.tar.gz"

# Create backup directory if not exists
mkdir -p $BACKUP_DIR

# Create backup
tar -czf $BACKUP_FILE $DATA_DIR

# Keep only last 30 backups
ls -t $BACKUP_DIR/prima-data-*.tar.gz | tail -n +31 | xargs -r rm

echo "Backup completed: $BACKUP_FILE"
```

**Make executable:**

```bash
sudo chmod +x /opt/prima/scripts/backup.sh
```

**Add to cron (daily at 2 AM):**

```bash
sudo crontab -e
```

**Add line:**

```cron
0 2 * * * /opt/prima/scripts/backup.sh >> /var/log/prima/backup.log 2>&1
```

### Restore from Backup

```bash
# Stop backend
sudo systemctl stop prima-backend

# Restore data
sudo tar -xzf /opt/backups/prima-data-20260102-020000.tar.gz -C /

# Restart backend
sudo systemctl start prima-backend
```

---

## Monitoring

### Health Check Endpoint

**Monitor backend health:**

```bash
curl https://prima.example.com/api/health
```

**Detailed health (requires auth):**

```bash
curl https://prima.example.com/api/health/detailed \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Log Monitoring

**Backend logs:**

```bash
sudo tail -f /var/log/prima/backend.log
```

**Nginx access logs:**

```bash
sudo tail -f /var/log/nginx/access.log
```

**Nginx error logs:**

```bash
sudo tail -f /var/log/nginx/error.log
```

### Uptime Monitoring

**Use external service:**

- Uptime Robot: https://uptimerobot.com/
- Pingdom: https://www.pingdom.com/
- StatusCake: https://www.statuscake.com/

**Monitor:**

- `https://prima.example.com/api/health` (every 5 minutes)
- Alert on non-200 response

### Resource Monitoring

**Install htop:**

```bash
sudo apt install htop -y
htop
```

**Check disk usage:**

```bash
df -h
du -sh /opt/prima/backend/uploads/
```

**Check memory:**

```bash
free -h
```

---

## Troubleshooting

### Backend Not Starting

**Check logs:**

```bash
sudo journalctl -u prima-backend -n 50 --no-pager
```

**Common issues:**

1. Port 8080 already in use

   ```bash
   sudo lsof -i :8080
   sudo kill -9 <PID>
   ```

2. Permission denied on data files

   ```bash
   sudo chown -R prima:prima /opt/prima/backend/data/
   ```

3. Missing config file
   ```bash
   ls -la /opt/prima/backend/config.yaml
   ```

### Frontend Not Loading

**Check Nginx:**

```bash
sudo nginx -t
sudo systemctl status nginx
```

**Check file permissions:**

```bash
ls -la /opt/prima/frontend/dist/
sudo chown -R www-data:www-data /opt/prima/frontend/dist/
```

**Check Nginx logs:**

```bash
sudo tail -f /var/log/nginx/error.log
```

### SSL Certificate Issues

**Check certificate:**

```bash
sudo certbot certificates
```

**Renew manually:**

```bash
sudo certbot renew --force-renewal
sudo systemctl reload nginx
```

### CORS Errors

**Update backend config:**

```yaml
server:
  cors_origin: "https://prima.example.com"
```

**Restart backend:**

```bash
sudo systemctl restart prima-backend
```

### GOWA Connection Failed

**Check GOWA service:**

```bash
curl http://localhost:3000/health
```

**Check circuit breaker state:**

```bash
curl https://prima.example.com/api/health/detailed \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Security Checklist

- [ ] Change default superadmin password
- [ ] Use strong JWT secret (32+ random characters)
- [ ] Use HTTPS (SSL certificate via Let's Encrypt)
- [ ] Set secure GOWA webhook secret
- [ ] Enable firewall (ufw/iptables)
- [ ] Disable root SSH login
- [ ] Use SSH keys (disable password auth)
- [ ] Keep system updated (`sudo apt update && sudo apt upgrade`)
- [ ] Set up automated backups
- [ ] Enable fail2ban (block brute force attempts)
- [ ] Monitor logs for suspicious activity
- [ ] Limit backend CORS to frontend domain only
- [ ] Use strong passwords for all services
- [ ] Set up monitoring/alerting

---

## Firewall Configuration

**Using ufw (Uncomplicated Firewall):**

```bash
# Enable firewall
sudo ufw enable

# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Check status
sudo ufw status
```

**Backend port 8080 should NOT be exposed (Nginx proxies it).**

---

## Updating Application

### Backend Update

```bash
# Stop backend
sudo systemctl stop prima-backend

# Backup current version
sudo cp /opt/prima/backend/prima-backend /opt/prima/backend/prima-backend.bak

# Pull latest code
cd /opt/prima
git pull origin main

# Rebuild backend
cd backend
go build -o prima-backend main.go

# Restart backend
sudo systemctl start prima-backend

# Check logs
sudo journalctl -u prima-backend -f
```

### Frontend Update

```bash
# Build locally (on dev machine)
cd frontend
git pull origin main
bun install
bun run build

# Transfer to server
rsync -avz --delete dist/ user@server:/opt/prima/frontend/dist/

# No restart needed (static files)
```

### Rollback

**Backend:**

```bash
sudo systemctl stop prima-backend
sudo cp /opt/prima/backend/prima-backend.bak /opt/prima/backend/prima-backend
sudo systemctl start prima-backend
```

**Frontend:**

```bash
# Restore from previous backup/deployment
```

---

## Production Optimization

### Gzip Compression

Already enabled in Nginx config above. Reduces transfer size by ~70%.

### Caching Strategy

**Static assets:**

- `Cache-Control: public, immutable` (1 year)

**index.html:**

- `Cache-Control: no-cache` (always check for updates)

**API responses:**

- No caching (dynamic data)

### CDN (Optional)

**Use Cloudflare or AWS CloudFront:**

- Cache static assets globally
- DDoS protection
- Free SSL

**Update DNS:**

- Point domain to CDN
- CDN forwards to origin server

---

## Docker Deployment (Alternative)

### Dockerfile - Backend

**`backend/Dockerfile`:**

```dockerfile
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o prima-backend main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/prima-backend .
COPY --from=builder /app/config.yaml .
EXPOSE 8080
CMD ["./prima-backend"]
```

### Dockerfile - Frontend

**`frontend/Dockerfile`:**

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json bun.lockb ./
RUN npm install -g bun && bun install
COPY . .
RUN bun run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### Docker Compose

**`docker-compose.yml`:**

```yaml
version: "3.8"

services:
  backend:
    build: ./backend
    ports:
      - "8080:8080"
    volumes:
      - ./backend/data:/root/data
      - ./backend/uploads:/root/uploads
    environment:
      - PORT=8080
      - CORS_ORIGIN=http://localhost
    restart: unless-stopped

  frontend:
    build: ./frontend
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped
```

**Deploy:**

```bash
docker-compose up -d
```

---

**Summary:** You now have a production-ready PRIMA deployment with Nginx reverse proxy, SSL/TLS, automated backups, and monitoring. Adjust configurations based on your specific infrastructure and requirements.
