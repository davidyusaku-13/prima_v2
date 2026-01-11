# User Guide - PRIMA

**A guide for healthcare volunteers and administrators**

---

## What is PRIMA?

PRIMA is a Healthcare Volunteer Dashboard that helps you:
- Manage patient information
- Send medication and appointment reminders via WhatsApp
- Access health education content (articles and videos)

---

## Roles and Permissions

| Role | Can Do |
|------|--------|
| **Volunteer** | Manage own patients, send reminders |
| **Admin** | All volunteer permissions + manage content (articles, videos) |
| **Superadmin** | All admin permissions + manage users |

---

## Logging In

1. Open PRIMA in your browser
2. Enter your username and password
3. Click **Login**

> **Forgot password?** Contact your administrator.

---

## Managing Patients

### Adding a Patient

1. Click **"Patients"** in the sidebar
2. Click the **"Add Patient"** button
3. Fill in the form:
   - **Name** (required): Patient's full name
   - **Phone** (required): WhatsApp number (e.g., 081234567890)
   - **Email** (optional): Patient's email
   - **Notes** (optional): Any relevant notes
4. Click **Save**

### Editing a Patient

1. Click on a patient's name to open their details
2. Click the **Edit** button
3. Make your changes
4. Click **Save**

### Deleting a Patient

1. Click on a patient's name
2. Click the **Delete** button
3. Confirm the deletion

> **Note:** Deleting a patient also deletes all their reminders.

---

## Managing Reminders

### Creating a Reminder

1. Open a patient's details
2. Click **"Add Reminder"**
3. Fill in the form:
   - **Title** (required): What the reminder is about
   - **Description** (optional): Additional details
   - **Due Date** (required): When to send the reminder
   - **Priority**: Low, Medium, or High
   - **Recurrence**: One-time, Daily, Weekly, or Monthly
4. Optionally attach educational content (articles or videos)
5. Click **Save**

### Sending a Reminder

1. Open a patient's details
2. Find the reminder you want to send
3. Click the **Send** button
4. The reminder will be sent via WhatsApp

> **Quiet Hours:** Messages are not sent between 8 PM and 8 AM. They will be automatically scheduled for 8 AM the next day.

### Reminder Status

| Status | Meaning |
|--------|---------|
| **Pending** | Not yet sent |
| **Scheduled** | Queued for quiet hours |
| **Sending** | Currently being sent |
| **Sent** | Sent to WhatsApp |
| **Delivered** | Received by patient |
| **Read** | Opened by patient |
| **Failed** | Could not be delivered |

### Canceling a Reminder

1. Find the reminder
2. Click the **Cancel** button
3. Confirm the cancellation

---

## Content Management (Admins Only)

### Adding an Article

1. Go to **CMS Dashboard**
2. Click **"New Article"**
3. Write your article using the rich text editor
4. Upload a hero image
5. Select a category
6. Click **Publish** (or Save as Draft)

### Adding a Video

1. Go to **CMS Dashboard**
2. Click **"Add Video"**
3. Paste a YouTube URL
4. Select a category
5. Click **Add**

> The system automatically fetches the video title and thumbnail from YouTube.

---

## Dashboard Statistics (Admins Only)

The dashboard shows:
- **Total Patients**: Number of patients in the system
- **Total Reminders**: Number of reminders created
- **Delivery Success Rate**: Percentage of reminders successfully delivered
- **Content Stats**: Article and video counts

---

## Troubleshooting

### My reminder failed to send

**Possible causes:**
1. **Invalid phone number** - Check the patient's phone number is correct
2. **WhatsApp not registered** - The number may not have WhatsApp
3. **GOWA service down** - Contact your administrator

**What to do:**
1. Check the error message in the reminder details
2. Fix the issue (e.g., correct phone number)
3. Click **Retry** to send again

### I can't see my patients

- **Volunteers** only see patients they created
- Ask an **Admin** if you need access to other patients

### The page won't load

1. Check your internet connection
2. Try refreshing the page (Ctrl+R or Cmd+R)
3. Clear your browser cache
4. Contact your administrator if the problem persists

---

## Language Settings

PRIMA supports:
- **English**
- **Indonesian (Bahasa Indonesia)**

To change language:
1. Click on your profile icon
2. Select your preferred language

---

## Getting Help

For technical issues, contact your system administrator.

For questions about using PRIMA, refer to this guide or ask your supervisor.

---

**Document Version:** January 11, 2026
