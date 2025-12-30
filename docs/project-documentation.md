# PRIMA - Healthcare Volunteer Dashboard

A full-stack web application for healthcare volunteers to manage patients and their care reminders with WhatsApp notifications.

## Overview

PRIMA is a volunteer management dashboard designed to help healthcare volunteers track patients and manage their care reminders. The application provides a clean, intuitive interface for maintaining patient records and scheduling reminders with automatic WhatsApp notifications via GOWA.

## Tech Stack

### Backend
- **Language**: Go 1.25+
- **Framework**: Gin (Go web framework)
- **CORS**: gin-contrib/cors
- **Storage**: JSON file persistence with in-memory cache
- **WhatsApp**: GOWA (Go WhatsApp Web Multi-Device)

### Frontend
- **Framework**: Svelte 5 + Vite
- **Styling**: Tailwind CSS 4
- **Package Manager**: Bun

## Project Structure

```
prima_v2/
├── backend/
│   ├── main.go              # Backend entry point and API handlers
│   ├── .env.example         # Environment variables template
│   ├── .env                 # Actual environment variables (gitignored)
│   ├── data/
│   │   ├── patients.json    # Patient data persistence
│   │   ├── users.json       # User accounts persistence
│   │   └── jwt_secret.txt   # JWT signing key
│   ├── go.mod               # Go module definition
│   └── go.sum               # Go dependencies
├── frontend/
│   ├── src/
│   │   ├── main.js          # Svelte app entry point
│   │   ├── App.svelte       # Main dashboard component
│   │   ├── app.css          # Tailwind CSS imports
│   │   └── assets/
│   │       └── svelte.svg   # Svelte logo
│   ├── index.html           # HTML entry point
│   ├── package.json         # Frontend dependencies
│   └── vite.config.js       # Vite configuration
├── docs/
│   ├── project-documentation.md
│   └── brownfield-architecture.md
└── CLAUDE.md                # Claude Code guidance
```

## Getting Started

### Prerequisites

- Go 1.25 or later
- Bun (for frontend)
- A web browser
- GOWA server running (for WhatsApp notifications)

### Environment Configuration

1. Copy the environment template:
```bash
cd backend
cp .env.example .env
```

2. Edit `.env` with your GOWA credentials:
```env
# GOWA Configuration (Go WhatsApp)
GOWA_ENDPOINT=http://localhost:3000
GOWA_USER=admin
GOWA_PASS=your_password_here
```

### Running the Application

**1. Start GOWA server** (WhatsApp bridge)
```bash
# Make sure GOWA is running on port 3000
```

**2. Start Backend (port 8080)**
```bash
cd backend
go run main.go
```

**3. Start Frontend (port 5173)**
```bash
cd frontend
bun run dev
```

Access the application at: `http://localhost:5173`

## WhatsApp Integration

### Setup

1. Ensure GOWA server is running
2. Configure `.env` with GOWA credentials
3. Add patients with valid Indonesian phone numbers

### Supported Phone Formats

| Input Format | Converted To |
|-------------|--------------|
| `081234567890` | `6281234567890@s.whatsapp.net` |
| `+6281234567890` | `6281234567890@s.whatsapp.net` |
| `6281234567890` | `6281234567890@s.whatsapp.net` |
| `81234567890` | `6281234567890@s.whatsapp.net` |

### Automatic Notifications

- Reminders are checked every 1 minute
- When a reminder is due (within 5 minutes of due time), a WhatsApp message is automatically sent
- Message includes: reminder title, description, priority, and recurrence info
- Once notified, a reminder won't send again until marked incomplete

## API Endpoints

### Health Check
| Method | Endpoint | Response |
|--------|----------|----------|
| GET | `/api/health` | `{"status": "ok"}` |

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register new user (role: volunteer) |
| POST | `/api/auth/login` | Login, returns JWT token |
| GET | `/api/auth/me` | Get current user info |

### Patients
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/patients` | List all patients |
| POST | `/api/patients` | Create a new patient |
| GET | `/api/patients/:id` | Get patient by ID |
| PUT | `/api/patients/:id` | Update patient |
| DELETE | `/api/patients/:id` | Delete patient |

### Reminders
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/patients/:id/reminders` | Create a reminder |
| PUT | `/api/patients/:id/reminders/:reminderId` | Update a reminder |
| POST | `/api/patients/:id/reminders/:reminderId/toggle` | Toggle completion |
| DELETE | `/api/patients/:id/reminders/:reminderId` | Delete a reminder |

### User Management (Superadmin Only)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users` | List all users |
| PUT | `/api/users/:id/role` | Update user role |
| DELETE | `/api/users/:id` | Delete a user |

## Data Models

### User
```json
{
  "id": "string",
  "username": "string (unique)",
  "fullName": "string",
  "password": "string (SHA256 hashed)",
  "role": "superadmin | admin | volunteer",
  "createdAt": "RFC3339 timestamp"
}
```

### Patient
```json
{
  "id": "string",
  "name": "string (required)",
  "phone": "string (required, Indonesian format)",
  "email": "string (optional)",
  "notes": "string (optional)",
  "reminders": ["Reminder array"]
}
```

### Reminder
```json
{
  "id": "string",
  "title": "string",
  "description": "string (optional)",
  "dueDate": "string (ISO datetime, e.g., 2025-12-27T14:30)",
  "priority": "low | medium | high",
  "completed": "boolean",
  "notified": "boolean",
  "recurrence": {
    "frequency": "none | daily | weekly | monthly | yearly",
    "interval": "number (default: 1)",
    "daysOfWeek": "array of numbers (0=Sunday, 6=Saturday)",
    "endDate": "string (optional)"
  }
}
```

### Recurrence Configuration

| Field | Description |
|-------|-------------|
| `frequency` | How often the reminder repeats |
| `interval` | Repeat every N units (e.g., 2 weeks) |
| `daysOfWeek` | Days of week for weekly recurrence [0,1,2,3,4,5,6] |
| `endDate` | When to stop recurring (optional) |

## Features

### Dashboard View
- **Stats Overview**: Total patients, total reminders, completed count, pending count
- **Upcoming Reminders**: Next 5 pending reminders sorted by due date
- **Recent Patients**: List of recently added patients with reminder counts
- **Recurrence Indicators**: Visual badge showing recurring reminders

### Patients View
- **Patient List**: All patients with search functionality
- **Patient Details**: Name, phone, email, and notes
- **Reminders Management**: View, add, edit, complete, and delete reminders per patient
- **Priority Levels**: Low (green), Medium (amber), High (red) with color coding
- **Recurrence Support**: Set reminders to repeat daily, weekly, monthly, or yearly

### Reminder Management
- Create reminders with title, description, due date, and priority
- Set custom recurring schedules
- Mark reminders as complete/incomplete
- Delete reminders
- Visual priority and recurrence indicators

### WhatsApp Notifications
- Automatic reminders when due
- Message includes all reminder details
- Only sends once per reminder occurrence

### Role-Based Access Control (RBAC)
Three user roles with different permission levels:

| Role | Permissions |
|------|-------------|
| **superadmin** | Full access, user management, all patients |
| **admin** | View all patients, cannot manage users |
| **volunteer** | Only see patients they created |

### User Management (Superadmin Only)
- View all registered users in a table
- Add new users with automatic volunteer role
- Edit user roles (admin or volunteer)
- Delete users (cannot delete yourself)
- User avatar with initials, role badges, creation date

## UI Components

### Layout
- **Sidebar**: Navigation with Dashboard and Patients views, user info
- **Header**: Search bar and "Add Patient" button
- **Main Content**: Dynamic content based on current view

### Modals
- **Auth Modal**: Login/Register forms with password strength indicator
- **Patient Modal**: Add/Edit patient form (name, phone required)
- **Reminder Modal**: Create/Edit reminder form with priority and recurrence options
- **User Modal**: Add new user or edit user role (superadmin only)
- **Confirm Modal**: Custom dialog for delete confirmations with warning icon

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `GOWA_ENDPOINT` | `http://localhost:3000` | GOWA server URL |
| `GOWA_USER` | `admin` | GOWA basic auth username |
| `GOWA_PASS` | `password123` | GOWA basic auth password |

### CORS
The backend is configured to accept requests only from `http://localhost:5173` (the Vite dev server).

### Data Storage
Data is persisted to `data/patients.json` as JSON. The backend:
- Loads existing data from file on startup
- Saves data after every write operation (create/update/delete)
- Background scheduler runs every 1 minute to check reminders

## Development

### Frontend
- Svelte 5 with reactive state management
- Tailwind CSS for styling
- Vite for fast development and building

### Backend
- Gin router for handling HTTP requests
- RWMutex for thread-safe in-memory storage
- Background goroutine for reminder checking (1-minute interval)
- GOWA API integration for WhatsApp messaging

## Troubleshooting

### WhatsApp Not Sending
1. Check GOWA server is running
2. Verify `.env` credentials are correct
3. Ensure phone number is in Indonesian format (08xx or 628xx)
4. Check backend logs for debug information

### Reminders Not Triggering
1. Verify system timezone matches local time
2. Check due date format is correct (2025-12-27T14:30)
3. Ensure reminder is not already marked as notified
4. Check backend logs for `[DEBUG]` messages

## Future Enhancements

Potential improvements for the application:
- Patient categories or groups
- Bulk operations
- Data export functionality (CSV/PDF)
- Mobile app version
- Offline support with service workers
- WhatsApp message status tracking
- Customizable notification templates
- Two-factor authentication
- Audit logging
- Patient notes history/versions

## License

This project is proprietary software.
