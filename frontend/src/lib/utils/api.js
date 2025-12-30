const API_URL = 'http://localhost:8080/api';

function getHeaders(token) {
  const headers = { 'Content-Type': 'application/json' };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  return headers;
}

// Auth
export async function login(username, password) {
  const res = await fetch(`${API_URL}/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Login failed');
  return data;
}

export async function register(username, password, fullName) {
  const res = await fetch(`${API_URL}/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password, fullName })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Registration failed');
  return data;
}

export async function fetchUser(token) {
  const res = await fetch(`${API_URL}/auth/me`, {
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Unauthorized');
  return res.json();
}

// Patients
export async function fetchPatients(token) {
  const res = await fetch(`${API_URL}/patients`, {
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Failed to fetch patients');
  const data = await res.json();
  return data.patients || [];
}

export async function savePatient(token, patient, editingId = null) {
  const method = editingId ? 'PUT' : 'POST';
  const url = editingId
    ? `${API_URL}/patients/${editingId}`
    : `${API_URL}/patients`;
  const res = await fetch(url, {
    method,
    headers: getHeaders(token),
    body: JSON.stringify(patient)
  });
  if (!res.ok) throw new Error('Failed to save patient');
  return res.json();
}

export async function deletePatient(token, id) {
  const res = await fetch(`${API_URL}/patients/${id}`, {
    method: 'DELETE',
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Failed to delete patient');
}

// Reminders
export async function saveReminder(token, patientId, reminder, editingId = null) {
  const method = editingId ? 'PUT' : 'POST';
  const url = editingId
    ? `${API_URL}/patients/${patientId}/reminders/${editingId}`
    : `${API_URL}/patients/${patientId}/reminders`;
  const res = await fetch(url, {
    method,
    headers: getHeaders(token),
    body: JSON.stringify(reminder)
  });
  if (!res.ok) throw new Error('Failed to save reminder');
  return res.json();
}

export async function toggleReminder(token, patientId, reminderId) {
  const res = await fetch(`${API_URL}/patients/${patientId}/reminders/${reminderId}/toggle`, {
    method: 'POST',
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Failed to toggle reminder');
}

export async function deleteReminder(token, patientId, reminderId) {
  const res = await fetch(`${API_URL}/patients/${patientId}/reminders/${reminderId}`, {
    method: 'DELETE',
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Failed to delete reminder');
}

export async function sendReminder(token, patientId, reminderId) {
  const res = await fetch(`${API_URL}/patients/${patientId}/reminders/${reminderId}/send`, {
    method: 'POST',
    headers: getHeaders(token)
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Failed to send reminder');
  return data;
}

export async function retryReminder(token, reminderId) {
  const res = await fetch(`${API_URL}/reminders/${reminderId}/retry`, {
    method: 'POST',
    headers: getHeaders(token)
  });
  const data = await res.json();
  if (!res.ok) {
    const error = new Error(data.error || 'Failed to retry reminder');
    error.code = data.code;
    throw error;
  }
  return data.data;
}

// Users (superadmin)
export async function fetchUsers(token) {
  const res = await fetch(`${API_URL}/users`, {
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Failed to fetch users');
  const data = await res.json();
  return data.users || [];
}

export async function updateUserRole(token, userId, role) {
  const res = await fetch(`${API_URL}/users/${userId}/role`, {
    method: 'PUT',
    headers: getHeaders(token),
    body: JSON.stringify({ role })
  });
  if (!res.ok) throw new Error('Failed to update user role');
}

export async function deleteUser(token, userId) {
  const res = await fetch(`${API_URL}/users/${userId}`, {
    method: 'DELETE',
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Failed to delete user');
}

export async function registerUser(token, username, password) {
  const res = await fetch(`${API_URL}/auth/register`, {
    method: 'POST',
    headers: { ...getHeaders(token), 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password })
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data.error || 'Failed to create user');
  return data;
}

// CMS - Articles
export async function fetchArticles(token = null, category = null, all = false) {
  const params = new URLSearchParams();
  if (category) params.append('category', category);
  if (all) params.append('all', 'true');
  const queryString = params.toString();
  const url = queryString ? `${API_URL}/articles?${queryString}` : `${API_URL}/articles`;
  const res = await fetch(url, { headers: getHeaders(token) });
  if (!res.ok) throw new Error('Failed to fetch articles');
  const data = await res.json();
  return data.articles || [];
}

export async function fetchArticle(token = null, id) {
  const res = await fetch(`${API_URL}/articles/${id}`, { headers: getHeaders(token) });
  if (!res.ok) throw new Error('Failed to fetch article');
  return res.json();
}

export async function createArticle(token, article) {
  const res = await fetch(`${API_URL}/articles`, {
    method: 'POST',
    headers: getHeaders(token),
    body: JSON.stringify(article)
  });
  if (!res.ok) throw new Error('Failed to create article');
  return res.json();
}

export async function updateArticle(token, id, article) {
  const res = await fetch(`${API_URL}/articles/${id}`, {
    method: 'PUT',
    headers: getHeaders(token),
    body: JSON.stringify(article)
  });
  if (!res.ok) throw new Error('Failed to update article');
  return res.json();
}

export async function deleteArticle(token, id) {
  const res = await fetch(`${API_URL}/articles/${id}`, {
    method: 'DELETE',
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Failed to delete article');
}

// CMS - Videos
export async function fetchVideos(token = null, category = null) {
  const url = category
    ? `${API_URL}/videos?category=${encodeURIComponent(category)}`
    : `${API_URL}/videos`;
  const res = await fetch(url, { headers: getHeaders(token) });
  if (!res.ok) throw new Error('Failed to fetch videos');
  const data = await res.json();
  return data.videos || [];
}

export async function createVideo(token, video) {
  const res = await fetch(`${API_URL}/videos`, {
    method: 'POST',
    headers: getHeaders(token),
    body: JSON.stringify(video)
  });
  if (!res.ok) throw new Error('Failed to create video');
  return res.json();
}

export async function updateVideo(token, id, video) {
  const res = await fetch(`${API_URL}/videos/${id}`, {
    method: 'PUT',
    headers: getHeaders(token),
    body: JSON.stringify(video)
  });
  if (!res.ok) throw new Error('Failed to update video');
  return res.json();
}

export async function deleteVideo(token, id) {
  const res = await fetch(`${API_URL}/videos/${id}`, {
    method: 'DELETE',
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Failed to delete video');
}

// CMS - Categories
export async function fetchCategories(token = null) {
  const res = await fetch(`${API_URL}/categories`, { headers: getHeaders(token) });
  if (!res.ok) throw new Error('Failed to fetch categories');
  const data = await res.json();
  return data.categories || [];
}

// CMS - All Content (articles + videos combined, for content picker)
export async function fetchAllContent(token = null, type = 'all', category = null) {
  const params = new URLSearchParams();
  if (type && type !== 'all') params.append('type', type);
  if (category) params.append('category', category);
  const queryString = params.toString();
  const url = queryString ? `${API_URL}/content?${queryString}` : `${API_URL}/content`;
  const res = await fetch(url, { headers: getHeaders(token) });
  if (!res.ok) throw new Error('Failed to fetch content');
  const data = await res.json();
  return { articles: data.articles || [], videos: data.videos || [] };
}

// CMS - Popular Content (for content suggestions)
export async function fetchPopularContent(limit = 5, token = null) {
  try {
    const res = await fetch(`${API_URL}/content/popular?limit=${limit}`, {
      headers: getHeaders(token)
    });
    if (!res.ok) {
      // Return empty array on error for graceful degradation
      return [];
    }
    const data = await res.json();
    return data.content || [];
  } catch (e) {
    // Return empty array on network error for graceful degradation
    console.warn('Failed to fetch popular content:', e.message);
    return [];
  }
}

// CMS - Increment Attachment Count
export async function incrementAttachmentCount(contentType, contentId, token = null) {
  const res = await fetch(`${API_URL}/content/${contentType}/${contentId}/increment-attachment`, {
    method: 'POST',
    headers: getHeaders(token)
  });
  if (!res.ok) {
    const data = await res.json();
    throw new Error(data.error || 'Failed to increment attachment count');
  }
  return res.json();
}

// CMS - Dashboard Stats
export async function fetchDashboardStats(token) {
  const res = await fetch(`${API_URL}/dashboard/stats`, {
    headers: getHeaders(token)
  });
  if (!res.ok) throw new Error('Failed to fetch dashboard stats');
  return res.json();
}

// CMS - Activity Log (returns empty for now - backend doesn't have this endpoint)
export async function fetchActivityLog(token) {
  // This endpoint doesn't exist in backend yet
  return [];
}

// CMS - Image Upload
export async function uploadImage(token, file) {
  const formData = new FormData();
  formData.append('image', file);
  const res = await fetch(`${API_URL}/upload/image`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}` },
    body: formData
  });
  if (!res.ok) throw new Error('Failed to upload image');
  const data = await res.json();
  return data.images;
}

// Config - Disclaimer
export async function fetchDisclaimer({ signal } = {}) {
  const res = await fetch(`${API_URL}/config/disclaimer`, signal ? { signal } : {});
  if (!res.ok) throw new Error('Failed to fetch disclaimer config');
  const data = await res.json();
  return data.data;
}

// Config - Quiet Hours
export async function fetchQuietHoursConfig() {
  const res = await fetch(`${API_URL}/config/quiet-hours`);
  if (!res.ok) throw new Error('Failed to fetch quiet hours config');
  const data = await res.json();
  return data.data;
}

// Analytics - Content Attachment Statistics
export async function getContentAnalytics(token) {
  const res = await fetch(`${API_URL}/analytics/content`, {
    headers: getHeaders(token)
  });
  if (!res.ok) {
    const data = await res.json();
    throw new Error(data.error || 'Failed to fetch content analytics');
  }
  const data = await res.json();
  return data.data;
}

/**
 * Check if current time is within quiet hours
 * @param {Object} config - Quiet hours config {start_hour, end_hour, timezone}
 * @returns {boolean} True if currently in quiet hours
 */
export function isQuietHours(config) {
  if (!config) return false;

  // Get current hour in the configured timezone
  // WIB = UTC+7, WITA = UTC+8, WIT = UTC+9
  const tzOffsets = { 'WIB': 7, 'WITA': 8, 'WIT': 9 };
  const offset = tzOffsets[config.timezone] || 7;

  const now = new Date();
  const utcHour = now.getUTCHours();
  const localHour = (utcHour + offset) % 24;

  // Quiet hours span midnight (e.g., 21:00 - 06:00)
  return localHour >= config.start_hour || localHour < config.end_hour;
}

/**
 * Get the next active time (when quiet hours end)
 * @param {Object} config - Quiet hours config {start_hour, end_hour, timezone}
 * @returns {string} ISO 8601 string of next active time
 */
export function getNextActiveTime(config) {
  if (!config) return null;

  const tzOffsets = { 'WIB': 7, 'WITA': 8, 'WIT': 9 };
  const offset = tzOffsets[config.timezone] || 7;

  const now = new Date();
  const utcHour = now.getUTCHours();
  const localHour = (utcHour + offset) % 24;

  // Calculate next active time (end_hour in local timezone)
  let nextActive = new Date(now);
  nextActive.setUTCMinutes(0, 0, 0);

  // Set to end_hour in local timezone (convert to UTC)
  const targetUTCHour = (config.end_hour - offset + 24) % 24;
  nextActive.setUTCHours(targetUTCHour);

  // If we're past end_hour today, move to tomorrow
  if (localHour >= config.end_hour) {
    nextActive.setUTCDate(nextActive.getUTCDate() + 1);
  }

  return nextActive.toISOString();
}
