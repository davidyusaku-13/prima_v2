# Quickstart Guide - PRIMA

**Get PRIMA running in 5 minutes**

---

## Prerequisites

| Tool | Version | Check Command |
|------|---------|---------------|
| Go | 1.25+ | `go version` |
| Bun | 1.0+ | `bun --version` |
| Git | Any | `git --version` |

---

## Step 1: Clone and Install (2 minutes)

```bash
# Clone the repository
git clone https://github.com/davidyusaku-13/prima_v2.git
cd prima_v2

# Install backend dependencies
cd backend
go mod download

# Install frontend dependencies
cd ../frontend
bun install
```

---

## Step 2: Start the Services (1 minute)

Open **two terminals**:

**Terminal 1 - Backend (port 8080):**
```bash
cd backend
go run main.go
```

**Terminal 2 - Frontend (port 5173):**
```bash
cd frontend
bun run dev
```

---

## Step 3: Login (30 seconds)

1. Open browser: `http://localhost:5173`
2. Login with default credentials:
   - **Username:** `superadmin`
   - **Password:** `superadmin`

---

## Step 4: First User Journey (1 minute)

### Create a Patient
1. Click **"Patients"** in sidebar
2. Click **"Add Patient"** button
3. Fill in:
   - Name: `Test Patient`
   - Phone: `081234567890`
4. Click **Save**

### Create a Reminder
1. Click on your new patient
2. Click **"Add Reminder"**
3. Fill in:
   - Title: `Test Reminder`
   - Due Date: Tomorrow
4. Click **Save**

---

## Optional: GOWA (WhatsApp Gateway)

To actually send WhatsApp messages, you need GOWA running:

```bash
# GOWA runs on port 3000
# See GOWA-README.md for setup instructions
```

Without GOWA, reminders will show "GOWA Unavailable" but the app functions normally.

---

## Common Issues

### Port Already in Use

```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# macOS/Linux
lsof -i :8080
kill -9 <PID>
```

### CORS Error

Ensure backend is running before frontend. Check `backend/config.yaml`:
```yaml
server:
  cors_origin: "http://localhost:5173"
```

### Frontend Can't Connect

1. Verify backend is running: `curl http://localhost:8080/api/health`
2. Check browser console for errors
3. Restart both services

---

## Next Steps

| Role | Next Document |
|------|---------------|
| **Frontend Developer** | [Development Guide - Frontend](./development-guide-frontend.md) |
| **Backend Developer** | [Development Guide - Backend](./development-guide-backend.md) |
| **DevOps** | [Deployment Guide](./deployment-guide.md) |
| **Architect** | [Architecture Decisions](./architecture-decisions.md) |

---

## Glossary

| Term | Meaning |
|------|---------|
| **GOWA** | Go WhatsApp Gateway - external service for sending messages |
| **Berita** | Indonesian for "News" - health education articles |
| **Video Edukasi** | Indonesian for "Educational Videos" |
| **WIB** | Western Indonesian Time (UTC+7) |
| **Quiet Hours** | 8 PM - 8 AM WIB - no messages sent during this time |

---

**Total Setup Time:** ~5 minutes
