package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/disintegration/imaging"

	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/utils"
)

// ContentStore holds all content stores
type ContentStore struct {
	Categories *models.CategoryStore
	Articles   *models.ArticleStore
	Videos     *models.VideoStore
}

// NewContentStore creates a new content store
func NewContentStore() *ContentStore {
	return &ContentStore{
		Categories: models.NewCategoryStore(),
		Articles:   models.NewArticleStore(),
		Videos:     models.NewVideoStore(),
	}
}

// Data file paths
const (
	categoriesDataFile = "data/categories.json"
	articlesDataFile   = "data/articles.json"
	videosDataFile     = "data/videos.json"
	uploadsDir         = "uploads"
)

// LoadContentData loads all content data from JSON files
func (cs *ContentStore) LoadContentData() {
	cs.loadCategories()
	cs.loadArticles()
	cs.loadVideos()
}

// SaveContentData saves all content data to JSON files
func (cs *ContentStore) SaveContentData() {
	go func() {
		cs.saveCategories()
		cs.saveArticles()
		cs.saveVideos()
	}()
}

func (cs *ContentStore) loadCategories() {
	data, err := os.ReadFile(categoriesDataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}

	var categories map[string]*models.Category
	if err := json.Unmarshal(data, &categories); err != nil {
		return
	}

	cs.Categories.Mu.Lock()
	cs.Categories.Categories = categories
	cs.Categories.ByType = make(map[models.CategoryType][]string)
	for id, cat := range categories {
		cs.Categories.ByType[cat.Type] = append(cs.Categories.ByType[cat.Type], id)
	}
	cs.Categories.Mu.Unlock()
}

func (cs *ContentStore) saveCategories() {
	cs.Categories.Mu.RLock()
	categories := make(map[string]*models.Category)
	for k, v := range cs.Categories.Categories {
		categories[k] = v
	}
	cs.Categories.Mu.RUnlock()

	data, err := json.MarshalIndent(categories, "", "  ")
	if err != nil {
		return
	}

	tmpFile := categoriesDataFile + ".tmp"
	err = os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		return
	}
	os.Rename(tmpFile, categoriesDataFile)
}

func (cs *ContentStore) loadArticles() {
	data, err := os.ReadFile(articlesDataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}

	var articles map[string]*models.Article
	if err := json.Unmarshal(data, &articles); err != nil {
		return
	}

	cs.Articles.Mu.Lock()
	cs.Articles.Articles = articles
	cs.Articles.BySlug = make(map[string]string)
	cs.Articles.ByCategory = make(map[string][]string)
	for id, art := range articles {
		cs.Articles.BySlug[art.Slug] = id
		cs.Articles.ByCategory[art.CategoryID] = append(cs.Articles.ByCategory[art.CategoryID], id)
	}
	cs.Articles.Mu.Unlock()
}

func (cs *ContentStore) saveArticles() {
	cs.Articles.Mu.RLock()
	articles := make(map[string]*models.Article)
	for k, v := range cs.Articles.Articles {
		articles[k] = v
	}
	cs.Articles.Mu.RUnlock()

	data, err := json.MarshalIndent(articles, "", "  ")
	if err != nil {
		return
	}

	tmpFile := articlesDataFile + ".tmp"
	err = os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		return
	}
	os.Rename(tmpFile, articlesDataFile)
}

func (cs *ContentStore) loadVideos() {
	data, err := os.ReadFile(videosDataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}

	var videos map[string]*models.Video
	if err := json.Unmarshal(data, &videos); err != nil {
		return
	}

	cs.Videos.Mu.Lock()
	cs.Videos.Videos = videos
	cs.Videos.ByCategory = make(map[string][]string)
	for id, vid := range videos {
		cs.Videos.ByCategory[vid.CategoryID] = append(cs.Videos.ByCategory[vid.CategoryID], id)
	}
	cs.Videos.Mu.Unlock()
}

func (cs *ContentStore) saveVideos() {
	cs.Videos.Mu.RLock()
	videos := make(map[string]*models.Video)
	for k, v := range cs.Videos.Videos {
		videos[k] = v
	}
	cs.Videos.Mu.RUnlock()

	data, err := json.MarshalIndent(videos, "", "  ")
	if err != nil {
		return
	}

	tmpFile := videosDataFile + ".tmp"
	err = os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		return
	}
	os.Rename(tmpFile, videosDataFile)
}

// EnsureUploadsDir creates the uploads directory if it doesn't exist
func EnsureUploadsDir() error {
	return os.MkdirAll(uploadsDir, 0755)
}

// GenerateID generates a unique ID
func GenerateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	now := time.Now().UnixNano()
	for i := range b {
		b[i] = letters[(now+int64(i))%int64(len(letters))]
	}
	return string(b)
}

// Category Handlers

// ListCategories returns all categories
func (cs *ContentStore) ListCategories(c *gin.Context) {
	cs.Categories.Mu.RLock()
	categories := make([]*models.Category, 0, len(cs.Categories.Categories))
	for _, cat := range cs.Categories.Categories {
		categories = append(categories, cat)
	}
	cs.Categories.Mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// CreateCategory creates a new category
func (cs *ContentStore) CreateCategory(c *gin.Context) {
	var req struct {
		Name string               `json:"name" binding:"required"`
		Type models.CategoryType  `json:"type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and type are required"})
		return
	}

	if req.Type != models.CategoryTypeArticle && req.Type != models.CategoryTypeVideo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'article' or 'video'"})
		return
	}

	cs.Categories.Mu.Lock()
	// Check for duplicate name
	for _, cat := range cs.Categories.Categories {
		if strings.EqualFold(cat.Name, req.Name) && cat.Type == req.Type {
			cs.Categories.Mu.Unlock()
			c.JSON(http.StatusConflict, gin.H{"error": "category with this name already exists"})
			return
		}
	}

	category := &models.Category{
		ID:        GenerateID(),
		Name:      req.Name,
		Type:      req.Type,
		CreatedAt: models.Now(),
	}
	cs.Categories.Categories[category.ID] = category
	cs.Categories.ByType[category.Type] = append(cs.Categories.ByType[category.Type], category.ID)
	cs.Categories.Mu.Unlock()

	cs.saveCategories()
	c.JSON(http.StatusCreated, category)
}

// Article Handlers

// ListArticles returns articles, optionally filtered by category
// Use ?all=true to include drafts (for CMS dashboard)
func (cs *ContentStore) ListArticles(c *gin.Context) {
	categoryID := c.Query("category")
	includeAll := c.Query("all") == "true"

	cs.Articles.Mu.RLock()
	articles := make([]*models.Article, 0)

	for _, art := range cs.Articles.Articles {
		// Only show published articles unless all=true
		if !includeAll && art.Status != models.ArticleStatusPublished {
			continue
		}

		// Filter by category if specified
		if categoryID != "" && art.CategoryID != categoryID {
			continue
		}

		articles = append(articles, art)
	}
	cs.Articles.Mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

// GetArticle returns an article by slug
func (cs *ContentStore) GetArticle(c *gin.Context) {
	slug := c.Param("slug")

	cs.Articles.Mu.RLock()
	articleID, exists := cs.Articles.BySlug[slug]
	var article *models.Article
	if exists {
		article = cs.Articles.Articles[articleID]
	}
	cs.Articles.Mu.RUnlock()

	if !exists || article == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	// Increment view count
	cs.Articles.Mu.Lock()
	if art, ok := cs.Articles.Articles[articleID]; ok {
		art.ViewCount++
	}
	cs.Articles.Mu.Unlock()
	cs.saveArticles()

	c.JSON(http.StatusOK, article)
}

// CreateArticle creates a new article
func (cs *ContentStore) CreateArticle(c *gin.Context) {
	var req struct {
		Title      string  `json:"title" binding:"required"`
		Excerpt    string  `json:"excerpt"`
		Content    string  `json:"content"`
		CategoryID string  `json:"category_id"`
		Status     string  `json:"status"`
		HeroImages struct {
			Hero16x9 string `json:"hero_16x9"`
			Hero1x1  string `json:"hero_1x1"`
			Hero4x3  string `json:"hero_4x3"`
		} `json:"hero_images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	// Generate slug
	slug := models.GenerateSlug(req.Title)

	// Ensure unique slug
	cs.Articles.Mu.RLock()
	originalSlug := slug
	counter := 1
	for {
		if _, exists := cs.Articles.BySlug[slug]; !exists {
			break
		}
		slug = fmt.Sprintf("%s-%d", originalSlug, counter)
		counter++
	}
	cs.Articles.Mu.RUnlock()

	status := models.ArticleStatusDraft
	if req.Status == "published" {
		status = models.ArticleStatusPublished
	}

	authorID := c.GetString("userID")

	article := &models.Article{
		ID:         GenerateID(),
		Title:      req.Title,
		Slug:       slug,
		Excerpt:    req.Excerpt,
		Content:    req.Content,
		AuthorID:   authorID,
		CategoryID: req.CategoryID,
		HeroImages: models.HeroImages{
			Hero16x9: req.HeroImages.Hero16x9,
			Hero1x1:  req.HeroImages.Hero1x1,
			Hero4x3:  req.HeroImages.Hero4x3,
		},
		Status:    status,
		Version:   1,
		ViewCount: 0,
		CreatedAt: models.Now(),
		UpdatedAt: models.Now(),
	}

	if status == models.ArticleStatusPublished {
		article.PublishedAt = models.Now()
	}

	cs.Articles.Mu.Lock()
	cs.Articles.Articles[article.ID] = article
	cs.Articles.BySlug[article.Slug] = article.ID
	if article.CategoryID != "" {
		cs.Articles.ByCategory[article.CategoryID] = append(cs.Articles.ByCategory[article.CategoryID], article.ID)
	}
	cs.Articles.Mu.Unlock()

	cs.saveArticles()
	c.JSON(http.StatusCreated, article)
}

// UpdateArticle updates an existing article
func (cs *ContentStore) UpdateArticle(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title      string  `json:"title"`
		Excerpt    string  `json:"excerpt"`
		Content    string  `json:"content"`
		CategoryID string  `json:"category_id"`
		Status     string  `json:"status"`
		HeroImages struct {
			Hero16x9 string `json:"hero_16x9"`
			Hero1x1  string `json:"hero_1x1"`
			Hero4x3  string `json:"hero_4x3"`
		} `json:"hero_images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cs.Articles.Mu.Lock()
	article, exists := cs.Articles.Articles[id]
	if !exists {
		cs.Articles.Mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	// Update fields
	if req.Title != "" {
		newSlug := models.GenerateSlug(req.Title)
		if newSlug != article.Slug {
			// Check if new slug is already taken
			if _, slugExists := cs.Articles.BySlug[newSlug]; slugExists {
				cs.Articles.Mu.Unlock()
				c.JSON(http.StatusConflict, gin.H{"error": "slug already exists"})
				return
			}
			delete(cs.Articles.BySlug, article.Slug)
			article.Slug = newSlug
			cs.Articles.BySlug[newSlug] = id
		}
		article.Title = req.Title
	}
	if req.Excerpt != "" {
		article.Excerpt = req.Excerpt
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.CategoryID != "" {
		article.CategoryID = req.CategoryID
	}
	if req.Status != "" {
		newStatus := models.ArticleStatus(req.Status)
		if newStatus == models.ArticleStatusPublished && article.Status != models.ArticleStatusPublished {
			article.PublishedAt = models.Now()
		}
		article.Status = newStatus
	}
	if req.HeroImages.Hero16x9 != "" {
		article.HeroImages.Hero16x9 = req.HeroImages.Hero16x9
	}
	if req.HeroImages.Hero1x1 != "" {
		article.HeroImages.Hero1x1 = req.HeroImages.Hero1x1
	}
	if req.HeroImages.Hero4x3 != "" {
		article.HeroImages.Hero4x3 = req.HeroImages.Hero4x3
	}

	article.Version++
	article.UpdatedAt = models.Now()
	cs.Articles.Mu.Unlock()

	cs.saveArticles()
	c.JSON(http.StatusOK, article)
}

// DeleteArticle deletes an article
func (cs *ContentStore) DeleteArticle(c *gin.Context) {
	id := c.Param("id")

	cs.Articles.Mu.Lock()
	article, exists := cs.Articles.Articles[id]
	if !exists {
		cs.Articles.Mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	delete(cs.Articles.Articles, id)
	delete(cs.Articles.BySlug, article.Slug)
	cs.Articles.Mu.Unlock()

	cs.saveArticles()
	c.JSON(http.StatusOK, gin.H{"message": "article deleted"})
}

// Video Handlers

// ListVideos returns published videos, optionally filtered by category
func (cs *ContentStore) ListVideos(c *gin.Context) {
	categoryID := c.Query("category")

	cs.Videos.Mu.RLock()
	videos := make([]*models.Video, 0)

	for _, vid := range cs.Videos.Videos {
		// Only show published videos
		if vid.Status != models.VideoStatusPublished {
			continue
		}

		// Filter by category if specified
		if categoryID != "" && vid.CategoryID != categoryID {
			continue
		}

		videos = append(videos, vid)
	}
	cs.Videos.Mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{"videos": videos})
}

// CreateVideo adds a new YouTube video
func (cs *ContentStore) CreateVideo(c *gin.Context) {
	var req struct {
		YouTubeURL string `json:"youtube_url" binding:"required"`
		CategoryID string `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "youtube_url is required"})
		return
	}

	// Validate YouTube URL
	videoID, err := utils.ExtractYouTubeID(req.YouTubeURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid YouTube URL"})
		return
	}

	// Check if video already exists
	cs.Videos.Mu.RLock()
	for _, vid := range cs.Videos.Videos {
		if vid.YouTubeID == videoID {
			cs.Videos.Mu.RUnlock()
			c.JSON(http.StatusConflict, gin.H{"error": "video already exists"})
			return
		}
	}
	cs.Videos.Mu.RUnlock()

	// Fetch metadata from noembed
	metadata, err := utils.FetchYouTubeMetadata(req.YouTubeURL)
	if err != nil {
		// Use fallback data
		metadata = &utils.YouTubeMetadata{
			Title:       "YouTube Video",
			Description: "",
			Thumbnail:   utils.GetYouTubeThumbnailURL(videoID, "high"),
		}
	}

	thumbnailURL := metadata.Thumbnail
	if thumbnailURL == "" {
		thumbnailURL = utils.GetYouTubeThumbnailURL(videoID, "high")
	}

	video := &models.Video{
		ID:           GenerateID(),
		YouTubeURL:   req.YouTubeURL,
		YouTubeID:    videoID,
		Title:        metadata.Title,
		Description:  metadata.Description,
		ChannelName:  metadata.AuthorName,
		ThumbnailURL: thumbnailURL,
		Duration:     metadata.Duration,
		CategoryID:   req.CategoryID,
		Status:       models.VideoStatusPublished,
		ViewCount:    0,
		CreatedAt:    models.Now(),
		UpdatedAt:    models.Now(),
	}

	cs.Videos.Mu.Lock()
	cs.Videos.Videos[video.ID] = video
	if video.CategoryID != "" {
		cs.Videos.ByCategory[video.CategoryID] = append(cs.Videos.ByCategory[video.CategoryID], video.ID)
	}
	cs.Videos.Mu.Unlock()

	cs.saveVideos()
	c.JSON(http.StatusCreated, video)
}

// DeleteVideo deletes a video
func (cs *ContentStore) DeleteVideo(c *gin.Context) {
	id := c.Param("id")

	cs.Videos.Mu.Lock()
	_, exists := cs.Videos.Videos[id]
	if !exists {
		cs.Videos.Mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	delete(cs.Videos.Videos, id)
	cs.Videos.Mu.Unlock()

	cs.saveVideos()
	c.JSON(http.StatusOK, gin.H{"message": "video deleted"})
}

// Image Upload Handler

// UploadImage handles image upload with resizing
func (cs *ContentStore) UploadImage(c *gin.Context) {
	// Ensure uploads directory exists
	if err := EnsureUploadsDir(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create uploads directory"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}
	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only JPEG, PNG, and WebP images are allowed"})
		return
	}

	// Read image data
	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image"})
		return
	}

	// Decode image
	img, err := imaging.Decode(strings.NewReader(string(imageData)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image format"})
		return
	}

	// Generate unique filename prefix
	filenamePrefix := GenerateID()

	// Resize and save images in different aspect ratios
	var heroImages models.HeroImages

	// 16:9 - 1920x1080
	img16x9 := imaging.Resize(img, 1920, 1080, imaging.Lanczos)
	filename16x9 := filepath.Join(uploadsDir, filenamePrefix+"_16x9.jpg")
	if err := imaging.Save(img16x9, filename16x9, imaging.JPEGQuality(85)); err == nil {
		heroImages.Hero16x9 = "/" + strings.ReplaceAll(filename16x9, "\\", "/")
	}

	// 1:1 - 1080x1080
	img1x1 := imaging.Resize(img, 1080, 1080, imaging.Lanczos)
	filename1x1 := filepath.Join(uploadsDir, filenamePrefix+"_1x1.jpg")
	if err := imaging.Save(img1x1, filename1x1, imaging.JPEGQuality(85)); err == nil {
		heroImages.Hero1x1 = "/" + strings.ReplaceAll(filename1x1, "\\", "/")
	}

	// 4:3 - 1600x1200
	img4x3 := imaging.Resize(img, 1600, 1200, imaging.Lanczos)
	filename4x3 := filepath.Join(uploadsDir, filenamePrefix+"_4x3.jpg")
	if err := imaging.Save(img4x3, filename4x3, imaging.JPEGQuality(85)); err == nil {
		heroImages.Hero4x3 = "/" + strings.ReplaceAll(filename4x3, "\\", "/")
	}

	c.JSON(http.StatusOK, gin.H{
		"images": heroImages,
	})
}

// Dashboard Stats Handler

// GetDashboardStats returns statistics for the CMS dashboard
func (cs *ContentStore) GetDashboardStats(c *gin.Context) {
	stats := gin.H{}

	// Category counts by type
	cs.Categories.Mu.RLock()
	articleCategories := 0
	videoCategories := 0
	for _, cat := range cs.Categories.Categories {
		if cat.Type == models.CategoryTypeArticle {
			articleCategories++
		} else if cat.Type == models.CategoryTypeVideo {
			videoCategories++
		}
	}
	cs.Categories.Mu.RUnlock()

	// Article counts by status
	cs.Articles.Mu.RLock()
	publishedArticles := 0
	draftArticles := 0
	totalViews := 0
	for _, art := range cs.Articles.Articles {
		if art.Status == models.ArticleStatusPublished {
			publishedArticles++
		} else {
			draftArticles++
		}
		totalViews += art.ViewCount
	}
	cs.Articles.Mu.RUnlock()

	// Video counts
	cs.Videos.Mu.RLock()
	totalVideos := len(cs.Videos.Videos)
	cs.Videos.Mu.RUnlock()

	stats["categories"] = gin.H{
		"articles": articleCategories,
		"videos":   videoCategories,
	}
	stats["articles"] = gin.H{
		"published": publishedArticles,
		"drafts":    draftArticles,
		"total":     publishedArticles + draftArticles,
	}
	stats["videos"] = gin.H{
		"total": totalVideos,
	}
	stats["total_views"] = gin.H{
		"articles": totalViews,
	}

	c.JSON(http.StatusOK, stats)
}

// GetCategoriesByType returns categories filtered by type
func (cs *ContentStore) GetCategoriesByType(c *gin.Context) {
	categoryType := models.CategoryType(c.Param("type"))

	if categoryType != models.CategoryTypeArticle && categoryType != models.CategoryTypeVideo {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'article' or 'video'"})
		return
	}

	cs.Categories.Mu.RLock()
	categories := make([]*models.Category, 0)
	for _, cat := range cs.Categories.Categories {
		if cat.Type == categoryType {
			categories = append(categories, cat)
		}
	}
	cs.Categories.Mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// IncrementVideoView increments the view count of a video
func (cs *ContentStore) IncrementVideoView(c *gin.Context) {
	id := c.Param("id")

	cs.Videos.Mu.Lock()
	video, exists := cs.Videos.Videos[id]
	if !exists {
		cs.Videos.Mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}
	video.ViewCount++
	cs.Videos.Mu.Unlock()

	cs.saveVideos()
	c.JSON(http.StatusOK, gin.H{"view_count": video.ViewCount})
}

// StrToInt converts string to int, returns 0 if invalid
func StrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
