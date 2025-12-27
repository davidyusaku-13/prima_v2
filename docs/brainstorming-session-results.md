# Brainstorming Session Results

**Date:** 2025-12-27
**Project:** PRIMA - Healthcare Volunteer Dashboard Expansion
**Topic:** Berita (News), Video Edukasi (Educational Videos), and CMS

---

## Executive Summary

This brainstorming session explored the addition of new public-facing features to PRIMA: a health news section (Berita), an educational video gallery (Video Edukasi), and a content management system (CMS) for administrators.

### Session Goals
- Define features for Berita and Video Edukasi pages
- Design CMS capabilities for content management
- Plan implementation approach with existing Go/Gin + Svelte 5 stack

### Techniques Used
- Progressive Flow (Features → Content Types → CMS → Architecture → Implementation)
- Divergent Thinking (generate ideas freely)
- Convergent Thinking (prioritize and categorize)
- Synthesis (concrete implementation plan)

### Key Outcomes
- Defined comprehensive feature set for both content types
- Established clear RBAC for content management
- Created detailed implementation plan for backend and frontend
- Completed implementation of all planned features

---

## Session Context

### Initial Questions & Answers

| Question | Answer |
|----------|--------|
| What are we brainstorming about? | Adding Berita, Video Edukasi pages and CMS |
| Target audience | Patients and public (external-facing) |
| Budget constraints | Free only (no paid services) |
| Backend approach | Extending existing JSON file persistence |

### Warm-up Insights
User identified that effective health content:
- Has good format and structure
- Is easy for regular users to understand
- Presentation style has room for improvement in existing content

---

## Technique Sections

### Technique 1: Berita (News/Articles) Features

**Generated Ideas:**
- Categories & tags organization
- Latest/Popular sorting options
- Adjustable text size settings
- Dark/light mode for reading
- Bookmark articles for later
- Author info and publishing dates
- WhatsApp sharing (important for Indonesian users)
- Other social sharing options
- Related content suggestions
- Table of contents for longer articles

**Key Decisions:**
- Draft/Publish workflow (admins see drafts only)
- Version history for rollback capability
- Hero images with 3 standardized aspect ratios: 1920×1080 (16:9), 1080×1080 (1:1), 1600×1200 (4:3)
- Image handling: both upload and URL paste options

### Technique 2: Video Edukasi (Educational Videos) Features

**Generated Ideas:**
- YouTube as hosting platform (zero bandwidth costs)
- Category/playlist organization
- Video metadata display (title, description, duration, channel)
- Player controls (speed, subtitles)
- Progress tracking for users

**Key Decisions:**
- YouTube embed approach: paste URL, auto-fetch metadata
- Backend fetches from noembed API to get: Title, Description, Channel Name, Thumbnail, Duration
- No comments or user-generated content initially

### Technique 3: CMS (Content Management System) Features

**Generated Ideas:**
- Dashboard with content overview stats
- CRUD operations for articles and videos
- Draft/publish workflow
- User role management for content creation
- Media library for images

**Key Decisions:**
- Dashboard stats: total articles, videos, views, most popular content, recent activity
- Media storage: local folder in backend (`backend/uploads/`)
- Auto-resize uploaded images to 3 standard aspect ratios
- Activity log: track who published what and when

---

## Access Control Matrix

| Role | Create Content | Publish Content | Read Public Content | Access CMS |
|------|---------------|-----------------|---------------------|------------|
| Superadmin | ✅ | ✅ | ✅ | ✅ |
| Admin | ✅ | ✅ | ✅ | ✅ |
| Volunteer | ❌ | ❌ | ✅ | ❌ |
| Public (Unauthenticated) | ❌ | ❌ | ✅ | ❌ |

---

## Idea Categorization

### Immediate Opportunities (Phase 1 - Completed)

**Core Pages & Navigation:**
- `Berita` page - Article list and detail views
- `Video Edukasi` page - YouTube embed video gallery
- Navigation links in sidebar and bottom nav
- Public access for both pages

**CMS Core:**
- Admin-only dashboard page
- Backend API endpoints for articles and videos CRUD
- YouTube URL parsing and metadata fetch
- Image upload with auto-resize to 3 formats

**Data Models:**
```go
type Category struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Type string `json:"type"` // "article" or "video"
    CreatedAt string `json:"created_at"`
}

type Article struct {
    ID, Title, Slug, Excerpt, Content string
    AuthorID, CategoryID string
    HeroImages struct {
        Hero16x9, Hero1x1, Hero4x3 string
    }
    Status string // "draft" or "published"
    Version, ViewCount int
    CreatedAt, PublishedAt, UpdatedAt string
}

type Video struct {
    ID, YouTubeURL, YouTubeID string
    Title, Description, ChannelName string
    ThumbnailURL, Duration string
    CategoryID, Status string
    ViewCount int
    CreatedAt, UpdatedAt string
}
```

### Future Innovations (Phase 2)

- Version history UI (view diffs, rollback to previous versions)
- Article bookmarks (localStorage for public users)
- WhatsApp sharing integration
- Dark/light mode for article reading
- Image URL paste option in addition to upload

### Moonshots (Ambitious)

- Comments/discussion on articles
- Video completion tracking with quizzes
- User-generated content submissions
- Email newsletter integration
- Search functionality across content

---

## Action Planning

### Priority 1: Backend Foundation

**Files Created:**
- `backend/models/content.go` - Category, Article, Video models with mutex protection
- `backend/utils/youtube.go` - YouTube URL parsing and metadata fetch via noembed API
- `backend/handlers/content.go` - CRUD handlers for all content types and image upload

**API Endpoints:**
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/categories` | Public | List all categories |
| GET | `/api/categories/:type` | Public | Get categories by type |
| POST | `/api/categories` | Admin+ | Create category |
| GET | `/api/articles` | Public | List published articles |
| GET | `/api/articles/:slug` | Public | Get article by slug |
| POST | `/api/articles` | Admin+ | Create article |
| PUT | `/api/articles/:id` | Admin+ | Update article |
| DELETE | `/api/articles/:id` | Admin+ | Delete article |
| GET | `/api/videos` | Public | List published videos |
| POST | `/api/videos` | Admin+ | Add YouTube video |
| DELETE | `/api/videos/:id` | Admin+ | Delete video |
| POST | `/api/upload/image` | Admin+ | Upload and resize image |
| GET | `/api/dashboard/stats` | Admin+ | Get CMS statistics |

**Data Files:**
- `backend/data/categories.json`
- `backend/data/articles.json`
- `backend/data/videos.json`

### Priority 2: Frontend Implementation

**New Components:**
| Component | Path | Description |
|-----------|------|-------------|
| `ArticleCard.svelte` | `src/lib/components/` | Article preview card with hero image |
| `VideoCard.svelte` | `src/lib/components/` | Video thumbnail with play overlay |
| `ImageUploader.svelte` | `src/lib/components/` | Drag & drop upload with preview |
| `DashboardStats.svelte` | `src/lib/components/` | Stats display (articles, videos, views) |
| `ActivityLog.svelte` | `src/lib/components/` | Recent activity feed |
| `VideoModal.svelte` | `src/lib/components/` | YouTube embed modal |

**New Views:**
| View | Path | Access | Description |
|------|------|--------|-------------|
| `BeritaView.svelte` | `src/lib/views/` | Public | Health news list with filtering |
| `BeritaDetailView.svelte` | `src/lib/views/` | Public | Full article reader |
| `VideoEdukasiView.svelte` | `src/lib/views/` | Public | YouTube video gallery |
| `CMSDashboardView.svelte` | `src/lib/views/` | Admin+ | Admin dashboard |
| `ArticleEditorView.svelte` | `src/lib/views/` | Admin+ | Create/edit articles |
| `VideoManagerView.svelte` | `src/lib/views/` | Admin+ | Add YouTube videos |

**API Functions Added:**
- `fetchArticles()`, `fetchArticle(slug)`, `createArticle()`, `updateArticle()`, `deleteArticle()`
- `fetchVideos()`, `createVideo()`, `updateVideo()`, `deleteVideo()`
- `fetchCategories()`
- `fetchDashboardStats()`, `fetchActivityLog()`
- `uploadImage()`

### Priority 3: Testing & Fixes

**Issues Resolved:**
1. Fixed endpoint paths in `api.js` (`/api/cms/stats` → `/api/dashboard/stats`)
2. Fixed data format to match backend (`category_id`, `hero_images.hero_16x9`, etc.)
3. Fixed component props for proper data display
4. Added token prop passing to editor views
5. Resolved accessibility warnings (aria-labels, label associations, tabindex)

---

## Implementation Order Completed

1. ✅ Backend: Categories & Video CRUD + YouTube fetch
2. ✅ Backend: Article CRUD + Image upload/resize
3. ✅ Frontend: Video Edukasi page (simplest - just embed)
4. ✅ Frontend: Berita list + detail pages
5. ✅ Frontend: CMS dashboard + article editor
6. ✅ Frontend: Navigation updates
7. ✅ Testing: Public access + admin workflows
8. ✅ Accessibility fixes

---

## Reflection & Follow-up

### What Worked Well
- Progressive technique flow helped organize thoughts systematically
- Clear requirements emerged early (YouTube for videos, local storage for images)
- Role-based access control was straightforward to define
- Implementation completed in single session

### Areas for Further Exploration
- Search functionality across articles and videos
- Category management UI (create/edit categories)
- Version history UI for articles
- User engagement features (bookmarks, comments)

### Questions for Future Sessions
- Should volunteers be able to contribute drafts for admin review?
- Should there be a content approval workflow?
- How to handle content moderation for public submissions?
- Integration with existing patient management features?

---

## Running the Application

```bash
# Terminal 1 - Backend (port 8080)
cd backend && go run main.go

# Terminal 2 - Frontend (port 5173)
cd frontend && bun run dev
```

### Default Access
- **Superadmin:** `superadmin` / `superadmin`
- **Admin:** Create via superadmin CMS
- **Public:** Access Berita and Video Edukasi without login

---

## Files Modified/Created

### Backend
- `backend/models/content.go` (new)
- `backend/utils/youtube.go` (new)
- `backend/handlers/content.go` (new)
- `backend/data/categories.json` (new)
- `backend/data/articles.json` (new)
- `backend/data/videos.json` (new)
- `backend/main.go` (modified - added routes)

### Frontend
- `frontend/src/lib/utils/api.js` (modified - added CMS functions)
- `frontend/src/lib/components/ArticleCard.svelte` (new)
- `frontend/src/lib/components/VideoCard.svelte` (new)
- `frontend/src/lib/components/ImageUploader.svelte` (new)
- `frontend/src/lib/components/DashboardStats.svelte` (new)
- `frontend/src/lib/components/ActivityLog.svelte` (new)
- `frontend/src/lib/components/VideoModal.svelte` (new)
- `frontend/src/lib/views/BeritaView.svelte` (new)
- `frontend/src/lib/views/BeritaDetailView.svelte` (new)
- `frontend/src/lib/views/VideoEdukasiView.svelte` (new)
- `frontend/src/lib/views/CMSDashboardView.svelte` (new)
- `frontend/src/lib/views/ArticleEditorView.svelte` (new)
- `frontend/src/lib/views/VideoManagerView.svelte` (new)
- `frontend/src/App.svelte` (modified - added routes)
- `frontend/src/lib/i18n/en.json` (modified - added translations)
- `frontend/src/lib/i18n/id.json` (modified - added translations)
