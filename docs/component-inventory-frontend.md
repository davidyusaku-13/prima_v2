# PRIMA Frontend Component Inventory

**Date:** 2026-01-13

## Overview

This document catalogs all Svelte components in the PRIMA frontend, organized by functional area. The frontend follows Svelte 5 patterns with runes (`$state`, `$derived`, `$effect`) for reactivity.

## Component Statistics

| Category | Count |
|----------|-------|
| Total Components | ~50 |
| Views (Pages) | 13 |
| Modal Components | 12 |
| Reusable UI Components | 25+ |
| Analytics Widgets | 3 |

## Components by Functional Area

### Authentication Components

| Component | File | Purpose |
|-----------|------|---------|
| **LoginScreen** | `views/LoginScreen.svelte` | Login and registration form with password strength validation |

### Navigation Components

| Component | File | Purpose |
|-----------|------|---------|
| **Sidebar** | `components/Sidebar.svelte` | Desktop navigation sidebar with role-based visibility |
| **BottomNav** | `components/BottomNav.svelte` | Mobile bottom navigation with admin dropdown |

### Dashboard Components

| Component | File | Purpose |
|-----------|------|---------|
| **DashboardView** | `views/DashboardView.svelte` | Main dashboard with stats, upcoming reminders, recent patients |
| **DashboardStats** | `components/DashboardStats.svelte` | Statistics display cards |

### Patient Management Components

| Component | File | Purpose |
|-----------|------|---------|
| **PatientsView** | `views/PatientsView.svelte` | Two-pane patient list and detail view |
| **PatientListPane** | `components/patients/PatientListPane.svelte` | Searchable patient list with selection |
| **PatientDetailPane** | `components/patients/PatientDetailPane.svelte` | Patient details with reminders and history tabs |
| **PatientDetailView** | `views/patients/PatientDetailView.svelte` | Modal overlay version of detail view |
| **PatientModal** | `components/PatientModal.svelte` | Add/edit patient form with validation |

### Reminder Components

| Component | File | Purpose |
|-----------|------|---------|
| **ReminderModal** | `components/ReminderModal.svelte` | Create/edit reminders with content attachment |
| **SendReminderModal** | `components/SendReminderModal.svelte` | Confirm and send with WhatsApp preview |
| **ReminderListTab** | `components/patients/ReminderListTab.svelte` | List of active reminders |
| **ReminderHistoryView** | `views/patients/ReminderHistoryView.svelte` | Paginated reminder history with cancel capability |
| **CancelConfirmationModal** | `components/reminders/CancelConfirmationModal.svelte` | Cancel scheduled reminder confirmation |

### Content Attachment Components

| Component | File | Purpose |
|-----------|------|---------|
| **ContentPickerModal** | `components/content/ContentPickerModal.svelte` | Select articles/videos to attach to reminders |
| **ContentPreviewPanel** | `components/content/ContentPreviewPanel.svelte` | Preview selected content |
| **ContentChip** | `components/content/ContentChip.svelte` | Display attached content as chips |
| **ContentDisclaimer** | `components/content/ContentDisclaimer.svelte` | Health content disclaimer |

### WhatsApp Preview

| Component | File | Purpose |
|-----------|------|---------|
| **WhatsAppPreview** | `components/whatsapp/WhatsAppPreview.svelte` | Preview WhatsApp message format |

### CMS Components (Admin)

| Component | File | Purpose |
|-----------|------|---------|
| **CMSDashboardView** | `views/CMSDashboardView.svelte` | Content management dashboard |
| **ArticleEditorView** | `views/ArticleEditorView.svelte` | Create/edit articles with Quill rich text editor |
| **VideoManagerView** | `views/VideoManagerView.svelte` | Add YouTube videos |
| **ArticleCard** | `components/ArticleCard.svelte` | Article display card |
| **VideoCard** | `components/VideoCard.svelte` | Video display card |
| **ContentListItem** | `components/ContentListItem.svelte` | List item for CMS dashboard |
| **QuillEditor** | `components/QuillEditor.svelte` | Rich text editor wrapper |
| **ImageUploader** | `components/ImageUploader.svelte` | Hero image upload component |

### Content Display Components

| Component | File | Purpose |
|-----------|------|---------|
| **BeritaView** | `views/BeritaView.svelte` | Health news articles list |
| **BeritaDetailView** | `views/BeritaDetailView.svelte` | Single article view |
| **VideoEdukasiView** | `views/VideoEdukasiView.svelte` | Educational videos list |
| **VideoModal** | `components/VideoModal.svelte` | Video playback modal |

### User Management (Superadmin)

| Component | File | Purpose |
|-----------|------|---------|
| **UsersView** | `views/UsersView.svelte` | User management list |
| **UserModal** | `components/UserModal.svelte` | Create/edit users, change roles |

### Analytics Components

| Component | File | Purpose |
|-----------|------|---------|
| **FailedDeliveriesView** | `views/analytics/FailedDeliveriesView.svelte` | Failed delivery list with filtering and CSV export |
| **CmsAnalyticsView** | `views/cms/CmsAnalyticsView.svelte` | Analytics page wrapper |
| **DeliveryAnalyticsWidget** | `components/analytics/DeliveryAnalyticsWidget.svelte` | Delivery statistics chart |
| **ContentAnalyticsWidget** | `components/analytics/ContentAnalyticsWidget.svelte` | Content attachment analytics |
| **FailedDeliveryCard** | `components/analytics/FailedDeliveryCard.svelte` | Individual failed delivery display |

### Health Monitoring

| Component | File | Purpose |
|-----------|------|---------|
| **SystemHealthWidget** | `components/health/SystemHealthWidget.svelte` | Backend health status display |

### Status Indicators

| Component | File | Purpose |
|-----------|------|---------|
| **DeliveryStatusBadge** | `components/delivery/DeliveryStatusBadge.svelte` | Badge with status icon and color |
| **DeliveryStatusFilter** | `components/delivery/DeliveryStatusFilter.svelte` | Filter tabs for delivery status |
| **FailedReminderBadge** | `components/indicators/FailedReminderBadge.svelte` | Badge showing failed reminder count |
| **QuietHoursHint** | `components/indicators/QuietHoursHint.svelte` | Indicator when in quiet hours |

### UI Components

| Component | File | Purpose |
|-----------|------|---------|
| **Toast** | `components/ui/Toast.svelte` | Global toast notification display |
| **EmptyState** | `components/ui/EmptyState.svelte` | Empty state placeholder |
| **ConfirmModal** | `components/ConfirmModal.svelte` | Generic confirmation dialog |
| **ProfileModal** | `components/ProfileModal.svelte` | User profile with locale switcher |
| **PhoneEditModal** | `components/PhoneEditModal.svelte` | Update patient phone after failed delivery |
| **VideoEditModal** | `components/VideoEditModal.svelte` | Edit video metadata |
| **ActivityLog** | `components/ActivityLog.svelte` | Activity history display |

## Component Hierarchy

```
App.svelte
├── Sidebar.svelte
├── BottomNav.svelte
├── Toast.svelte
└── Views (conditional rendering)
    ├── LoginScreen.svelte
    ├── DashboardView.svelte
    │   └── DashboardStats.svelte
    ├── PatientsView.svelte
    │   ├── PatientListPane.svelte
    │   └── PatientDetailPane.svelte
    │       └── ReminderListTab.svelte
    ├── UsersView.svelte
    ├── CMSDashboardView.svelte
    │   ├── ContentListItem.svelte
    │   ├── ArticleCard.svelte
    │   └── VideoCard.svelte
    ├── ArticleEditorView.svelte
    │   ├── QuillEditor.svelte
    │   └── ImageUploader.svelte
    ├── VideoManagerView.svelte
    ├── BeritaView.svelte
    ├── BeritaDetailView.svelte
    ├── VideoEdukasiView.svelte
    │   └── VideoModal.svelte
    ├── analytics/FailedDeliveriesView.svelte
    │   ├── DeliveryStatusFilter.svelte
    │   └── FailedDeliveryCard.svelte
    └── cms/CmsAnalyticsView.svelte
        ├── DeliveryAnalyticsWidget.svelte
        └── ContentAnalyticsWidget.svelte

Modals (overlay)
├── PatientModal.svelte
├── ReminderModal.svelte
│   ├── ContentPickerModal.svelte
│   │   └── ContentPreviewPanel.svelte
│   └── WhatsAppPreview.svelte
├── SendReminderModal.svelte
├── ContentPickerModal.svelte
├── PhoneEditModal.svelte
├── UserModal.svelte
├── ConfirmModal.svelte
├── ProfileModal.svelte
├── VideoEditModal.svelte
├── CancelConfirmationModal.svelte
└── VideoModal.svelte

Indicators (inline)
├── DeliveryStatusBadge.svelte
├── FailedReminderBadge.svelte
├── QuietHoursHint.svelte
└── ContentDisclaimer.svelte
```

## Design System Elements

### Reusable UI Patterns

| Pattern | Components | Usage |
|---------|------------|-------|
| **Modal Wrapper** | All *Modal.svelte | Consistent modal layout with header, body, footer |
| **Card Layout** | ArticleCard, VideoCard, FailedDeliveryCard | Image, title, description, actions |
| **Status Badge** | DeliveryStatusBadge | Color-coded status indicators |
| **Filter Tabs** | DeliveryStatusFilter | Tab-based filtering |
| **Toast Notifications** | Toast.svelte + toastStore | Success, error, info, warning |

### Tailwind CSS Usage

The frontend uses Tailwind CSS 4 for styling:
- Utility-first approach
- Custom theme configuration in CSS
- Responsive design with mobile-first breakpoints
- Dark mode support (if configured)

## State Management Components

### Stores (State Management)

| Store | File | Type | Purpose |
|-------|------|------|---------|
| `auth` | `stores/auth.js` | Legacy Svelte | Authentication state |
| `deliveryStore` | `stores/delivery.svelte.js` | Runes Class | Real-time delivery statuses |
| `toastStore` | `stores/toast.svelte.js` | Runes Class | Toast notifications |

### Services

| Service | File | Purpose |
|---------|------|---------|
| `sseService` | `services/sse.js` | Server-Sent Events connection |

## Component Testing

### Test Files Structure

| Component | Test File | Coverage |
|-----------|-----------|----------|
| ContentPickerModal | `ContentPickerModal.*.test.js` | Rendering, selection, interaction |
| DeliveryStatusBadge | `DeliveryStatusBadge.test.js` | Rendering |
| DeliveryStatusFilter | `DeliveryStatusFilter.test.js` | Rendering |
| FailedReminderBadge | `FailedReminderBadge.test.js` | Rendering |
| Toast | `Toast.test.js` | Rendering, behavior |
| DeliveryStore | `delivery.test.js` | State management |
| ToastStore | `toast.test.js` | State management |
| PatientDetailPane | `PatientDetailPane.test.js` | Rendering |
| PatientListPane | `PatientListPane.test.js` | Rendering |
| ReminderHistoryView | `ReminderHistoryView.test.js` | Rendering |
| analytics-navigation | `analytics-navigation.test.js` | Navigation |

---

_Generated using BMAD Method `document-project` workflow_
