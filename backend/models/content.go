package models

import (
	"strings"
	"sync"
	"time"
)

// CategoryType represents the type of category
type CategoryType string

const (
	CategoryTypeArticle CategoryType = "article"
	CategoryTypeVideo   CategoryType = "video"
)

// Category represents a content category
type Category struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Type      CategoryType `json:"type"`
	CreatedAt string       `json:"created_at"`
}

// ArticleStatus represents the publication status of an article
type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusPublished ArticleStatus = "published"
)

// VideoStatus represents the publication status of a video
type VideoStatus string

const (
	VideoStatusPublished VideoStatus = "published"
)

// HeroImages represents the different hero image sizes for an article
type HeroImages struct {
	Hero16x9 string `json:"hero_16x9"` // 1920x1080
	Hero1x1  string `json:"hero_1x1"`  // 1080x1080
	Hero4x3  string `json:"hero_4x3"`  // 1600x1200
}

// Article represents a news/educational article
type Article struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Slug        string       `json:"slug"`
	Excerpt     string       `json:"excerpt"`
	Content     string       `json:"content"`
	AuthorID    string       `json:"author_id"`
	CategoryID  string       `json:"category_id"`
	HeroImages  HeroImages   `json:"hero_images"`
	Status      ArticleStatus `json:"status"`
	Version     int          `json:"version"`
	ViewCount   int          `json:"view_count"`
	CreatedAt   string       `json:"created_at"`
	PublishedAt string       `json:"published_at"`
	UpdatedAt   string       `json:"updated_at"`
}

// Video represents an educational YouTube video
type Video struct {
	ID           string     `json:"id"`
	YouTubeURL   string     `json:"youtube_url"`
	YouTubeID    string     `json:"youtube_id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	ChannelName  string     `json:"channel_name"`
	ThumbnailURL string     `json:"thumbnail_url"`
	Duration     string     `json:"duration"`
	CategoryID   string     `json:"category_id"`
	Status       VideoStatus `json:"status"`
	ViewCount    int        `json:"view_count"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
}

// YouTubeMetadata represents metadata fetched from YouTube oEmbed/noembed
type YouTubeMetadata struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ChannelName string `json:"author_name"`
	Thumbnail   string `json:"thumbnail_url"`
	Duration    string `json:"duration,omitempty"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	HTML        string `json:"html"`
}

// UploadResponse represents the response after uploading an image
type UploadResponse struct {
	Images HeroImages `json:"images"`
}

// CategoryStore handles category data persistence
type CategoryStore struct {
	Mu          sync.RWMutex
	Categories  map[string]*Category
	ByType      map[CategoryType][]string // type -> category IDs
}

// ArticleStore handles article data persistence
type ArticleStore struct {
	Mu          sync.RWMutex
	Articles    map[string]*Article
	BySlug      map[string]string // slug -> article ID
	ByCategory  map[string][]string // category ID -> article IDs
}

// VideoStore handles video data persistence
type VideoStore struct {
	Mu          sync.RWMutex
	Videos      map[string]*Video
	ByCategory  map[string][]string // category ID -> video IDs
}

// NewCategoryStore creates a new category store
func NewCategoryStore() *CategoryStore {
	return &CategoryStore{
		Categories:  make(map[string]*Category),
		ByType:      make(map[CategoryType][]string),
	}
}

// NewArticleStore creates a new article store
func NewArticleStore() *ArticleStore {
	return &ArticleStore{
		Articles:    make(map[string]*Article),
		BySlug:      make(map[string]string),
		ByCategory:  make(map[string][]string),
	}
}

// NewVideoStore creates a new video store
func NewVideoStore() *VideoStore {
	return &VideoStore{
		Videos:      make(map[string]*Video),
		ByCategory:  make(map[string][]string),
	}
}

// GenerateSlug generates a URL-friendly slug from title
func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove special characters except hyphens
	var result strings.Builder
	for _, ch := range slug {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			result.WriteRune(ch)
		}
	}

	// Remove multiple consecutive hyphens
	slug = result.String()
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// Now returns current time in RFC3339 format
func Now() string {
	return time.Now().Format(time.RFC3339)
}
