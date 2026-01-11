# Error Codes Reference - PRIMA

**Centralized reference for all API error codes**

---

## Error Response Format

All API errors return this structure:

```json
{
  "error": "Human-readable error message",
  "code": "ERROR_CODE",
  "details": "Additional context (optional)"
}
```

---

## HTTP Status Codes

| Status | Meaning | When Used |
|--------|---------|-----------|
| 400 | Bad Request | Invalid input, validation failure |
| 401 | Unauthorized | Missing or invalid token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Duplicate resource |
| 500 | Internal Server Error | Unexpected server error |
| 503 | Service Unavailable | External service down |

---

## Authentication Errors (401)

| Code | Message | Cause | Fix |
|------|---------|-------|-----|
| `MISSING_TOKEN` | Missing authorization header | No Bearer token in request | Add `Authorization: Bearer <token>` header |
| `INVALID_TOKEN` | Invalid or expired token | Token malformed or expired | Re-login to get new token |
| `TOKEN_EXPIRED` | Token has expired | JWT past 7-day expiry | Re-login to get new token |
| `UNAUTHORIZED` | Invalid credentials | Wrong username/password | Check credentials |

---

## Authorization Errors (403)

| Code | Message | Cause | Fix |
|------|---------|-------|-----|
| `FORBIDDEN` | Access denied | Role insufficient for action | Contact admin for role upgrade |
| `NOT_OWNER` | Cannot access other user's data | Volunteer accessing another's patient | Only access own patients |
| `CANNOT_MODIFY_SELF` | Cannot change own role | Superadmin trying to change own role | Have another superadmin make change |
| `CANNOT_DELETE_SELF` | Cannot delete yourself | User trying to self-delete | Have another superadmin delete |

---

## Validation Errors (400)

| Code | Message | Cause | Fix |
|------|---------|-------|-----|
| `VALIDATION_ERROR` | Invalid input | Generic validation failure | Check request body against schema |
| `INVALID_PHONE` | Invalid phone number format | Phone doesn't match Indonesian format | Use format: 081234567890 or +6281234567890 |
| `INVALID_EMAIL` | Invalid email format | Email validation failed | Use valid email format |
| `REQUIRED_FIELD` | Field is required | Missing required field | Include required field in request |
| `INVALID_DATE` | Invalid date format | Date parsing failed | Use ISO 8601: 2026-01-15T10:00:00Z |
| `TITLE_TOO_LONG` | Title exceeds maximum length | Title > 200 characters | Shorten title |
| `MAX_ATTACHMENTS_EXCEEDED` | Maximum 3 attachments allowed | Reminder has > 3 attachments | Remove excess attachments |
| `INVALID_ATTACHMENT` | Content ID not found | Attachment references non-existent content | Use valid article/video ID |
| `INVALID_PRIORITY` | Invalid priority value | Priority not low/medium/high | Use: low, medium, or high |
| `INVALID_FREQUENCY` | Invalid recurrence frequency | Frequency not recognized | Use: none, daily, weekly, monthly |

---

## Resource Errors (404)

| Code | Message | Cause | Fix |
|------|---------|-------|-----|
| `NOT_FOUND` | Resource not found | Generic not found | Check resource ID |
| `PATIENT_NOT_FOUND` | Patient not found | Patient ID doesn't exist | Verify patient ID |
| `REMINDER_NOT_FOUND` | Reminder not found | Reminder ID doesn't exist | Verify reminder ID |
| `ARTICLE_NOT_FOUND` | Article not found | Article slug/ID doesn't exist | Verify article slug |
| `VIDEO_NOT_FOUND` | Video not found | Video ID doesn't exist | Verify video ID |
| `CATEGORY_NOT_FOUND` | Category not found | Category ID doesn't exist | Verify category ID |
| `USER_NOT_FOUND` | User not found | User ID doesn't exist | Verify user ID |

---

## Conflict Errors (409)

| Code | Message | Cause | Fix |
|------|---------|-------|-----|
| `CONFLICT` | Resource already exists | Generic conflict | Check for duplicates |
| `USERNAME_EXISTS` | Username already taken | Registration with existing username | Choose different username |
| `VIDEO_EXISTS` | Video already added | YouTube URL already in system | Video already exists, no action needed |
| `ALREADY_SENDING` | Reminder is already being sent | Concurrent send attempt | Wait for current send to complete |

---

## External Service Errors (503)

| Code | Message | Cause | Fix |
|------|---------|-------|-----|
| `GOWA_UNAVAILABLE` | WhatsApp service unavailable | GOWA server down or circuit breaker open | Wait and retry, or check GOWA status |
| `GOWA_TIMEOUT` | WhatsApp request timed out | GOWA didn't respond in time | Retry later |
| `RETRY_SCHEDULED` | Message queued for retry | Transient GOWA failure | Message will be retried automatically |
| `CIRCUIT_BREAKER_OPEN` | Service temporarily disabled | 5+ consecutive failures | Wait 5 minutes for half-open state |

---

## Reminder Status Errors (400)

| Code | Message | Cause | Fix |
|------|---------|-------|-----|
| `INVALID_STATUS` | Invalid status transition | Action not allowed for current status | Check reminder status first |
| `CANNOT_CANCEL` | Cannot cancel sent reminder | Reminder already sent/delivered | Only pending/scheduled can be cancelled |
| `CANNOT_RETRY` | Cannot retry non-failed reminder | Retry on non-failed reminder | Only failed reminders can be retried |
| `MAX_RETRIES_EXCEEDED` | Maximum retry attempts reached | 5 retry attempts exhausted | Investigate root cause, create new reminder |

---

## File Upload Errors (400)

| Code | Message | Cause | Fix |
|------|---------|-------|-----|
| `INVALID_FILE` | Invalid file format | File type not supported | Use JPEG, PNG, or WebP |
| `FILE_TOO_LARGE` | File exceeds size limit | File > 10MB | Reduce file size |
| `UPLOAD_FAILED` | File upload failed | Server error during upload | Retry upload |

---

## Webhook Errors (401/400)

| Code | Message | Cause | Fix |
|------|---------|-------|-----|
| `MISSING_SIGNATURE` | Missing X-Webhook-Signature | No HMAC signature header | Include signature header |
| `INVALID_SIGNATURE` | Signature verification failed | HMAC doesn't match | Check webhook secret configuration |
| `INVALID_PAYLOAD` | Malformed webhook payload | JSON parsing failed | Check payload format |

---

## Frontend Error Handling

### Recommended Pattern

```javascript
async function handleApiCall() {
  try {
    const result = await api.post('/api/patients', data);
    showSuccess('Patient created');
    return result;
  } catch (error) {
    const { code, error: message } = error.response?.data || {};

    switch (code) {
      case 'INVALID_PHONE':
        showError('Please enter a valid Indonesian phone number');
        break;
      case 'VALIDATION_ERROR':
        showError('Please check your input');
        break;
      case 'UNAUTHORIZED':
        auth.logout();
        window.location.href = '/login';
        break;
      case 'GOWA_UNAVAILABLE':
        showError('WhatsApp service is temporarily unavailable');
        break;
      default:
        showError(message || 'An unexpected error occurred');
    }
  }
}
```

### Error Display Components

| Error Type | Component | Behavior |
|------------|-----------|----------|
| Validation | Inline field error | Show below input |
| Authentication | Toast + redirect | Redirect to login |
| Not Found | Toast | Show message |
| Conflict | Toast | Show message |
| Server Error | Toast | Show generic message |
| GOWA Unavailable | Banner | Show status in UI |

---

## Debugging Tips

### Check the Response

```bash
curl -X POST http://localhost:8080/api/patients \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "phone": "invalid"}' \
  -v
```

### Common Mistakes

| Symptom | Likely Cause | Solution |
|---------|--------------|----------|
| 401 on all requests | Token not being sent | Check Authorization header |
| 403 on admin endpoints | Wrong role | Login as admin/superadmin |
| 400 INVALID_PHONE | Wrong format | Use 08xxx or 628xxx format |
| 503 GOWA_UNAVAILABLE | GOWA not running | Start GOWA or wait for circuit breaker |

---

**Last Updated:** January 11, 2026
