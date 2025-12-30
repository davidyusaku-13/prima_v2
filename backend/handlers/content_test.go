package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/davidyusaku-13/prima_v2/models"
)

func TestListAllContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a new content store
	cs := NewContentStore()

	// Add test data
	cs.Categories.Mu.Lock()
	cs.Categories.Categories = map[string]*models.Category{
		"cat-1": {ID: "cat-1", Name: "Health", Type: models.CategoryTypeArticle, CreatedAt: models.Now()},
		"cat-2": {ID: "cat-2", Name: "Nutrition", Type: models.CategoryTypeVideo, CreatedAt: models.Now()},
	}
	cs.Categories.ByType = map[models.CategoryType][]string{
		models.CategoryTypeArticle: {"cat-1"},
		models.CategoryTypeVideo:   {"cat-2"},
	}
	cs.Categories.Mu.Unlock()

	cs.Articles.Mu.Lock()
	cs.Articles.Articles = map[string]*models.Article{
		"art-1": {
			ID:         "art-1",
			Title:      "Healthy Eating Guide",
			Slug:       "healthy-eating-guide",
			Excerpt:    "A guide to eating healthy",
			CategoryID: "cat-1",
			Status:     models.ArticleStatusPublished,
			CreatedAt:  models.Now(),
		},
		"art-2": {
			ID:         "art-2",
			Title:      "Exercise Tips",
			Slug:       "exercise-tips",
			Excerpt:    "Tips for better exercise",
			CategoryID: "cat-1",
			Status:     models.ArticleStatusDraft, // Should not be returned
			CreatedAt:  models.Now(),
		},
	}
	cs.Articles.BySlug = map[string]string{"art-1": "art-1", "art-2": "art-2"}
	cs.Articles.ByCategory = map[string][]string{"cat-1": {"art-1", "art-2"}}
	cs.Articles.Mu.Unlock()

	cs.Videos.Mu.Lock()
	cs.Videos.Videos = map[string]*models.Video{
		"vid-1": {
			ID:           "vid-1",
			Title:        "Morning Yoga",
			YouTubeID:    "abc123",
			CategoryID:   "cat-2",
			Status:       models.VideoStatusPublished,
			ThumbnailURL: "https://example.com/thumb.jpg",
			CreatedAt:    models.Now(),
		},
	}
	cs.Videos.ByCategory = map[string][]string{"cat-2": {"vid-1"}}
	cs.Videos.Mu.Unlock()

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		checkFunc      func(t *testing.T, response map[string]interface{})
	}{
		{
			name:           "returns all content by default",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, response map[string]interface{}) {
				articles := response["articles"].([]interface{})
				videos := response["videos"].([]interface{})

				if len(articles) != 1 {
					t.Errorf("Expected 1 published article, got %d", len(articles))
				}
				if len(videos) != 1 {
					t.Errorf("Expected 1 published video, got %d", len(videos))
				}
			},
		},
		{
			name:           "filters by type=article",
			queryParams:    "?type=article",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, response map[string]interface{}) {
				articles := response["articles"].([]interface{})
				videos := response["videos"].([]interface{})

				if len(articles) != 1 {
					t.Errorf("Expected 1 article, got %d", len(articles))
				}
				if len(videos) != 0 {
					t.Errorf("Expected 0 videos, got %d", len(videos))
				}
			},
		},
		{
			name:           "filters by type=video",
			queryParams:    "?type=video",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, response map[string]interface{}) {
				articles := response["articles"].([]interface{})
				videos := response["videos"].([]interface{})

				if len(articles) != 0 {
					t.Errorf("Expected 0 articles, got %d", len(articles))
				}
				if len(videos) != 1 {
					t.Errorf("Expected 1 video, got %d", len(videos))
				}
			},
		},
		{
			name:           "filters by category",
			queryParams:    "?category=cat-1",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, response map[string]interface{}) {
				articles := response["articles"].([]interface{})
				videos := response["videos"].([]interface{})

				if len(articles) != 1 {
					t.Errorf("Expected 1 article in cat-1, got %d", len(articles))
				}
				if len(videos) != 0 {
					t.Errorf("Expected 0 videos (wrong category), got %d", len(videos))
				}
			},
		},
		{
			name:           "filters by type and category",
			queryParams:    "?type=all&category=cat-2",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, response map[string]interface{}) {
				articles := response["articles"].([]interface{})
				videos := response["videos"].([]interface{})

				if len(articles) != 0 {
					t.Errorf("Expected 0 articles (wrong category), got %d", len(articles))
				}
				if len(videos) != 1 {
					t.Errorf("Expected 1 video in cat-2, got %d", len(videos))
				}
			},
		},
		{
			name:           "type=all returns everything",
			queryParams:    "?type=all",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, response map[string]interface{}) {
				articles := response["articles"].([]interface{})
				videos := response["videos"].([]interface{})

				if len(articles) != 1 {
					t.Errorf("Expected 1 article, got %d", len(articles))
				}
				if len(videos) != 1 {
					t.Errorf("Expected 1 video, got %d", len(videos))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/content"+tt.queryParams, nil)

			cs.ListAllContent(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, response)
			}
		})
	}
}

func TestListAllContent_EmptyStore(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := NewContentStore()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/content", nil)

	cs.ListAllContent(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	articles := response["articles"].([]interface{})
	videos := response["videos"].([]interface{})

	if len(articles) != 0 {
		t.Errorf("Expected 0 articles, got %d", len(articles))
	}
	if len(videos) != 0 {
		t.Errorf("Expected 0 videos, got %d", len(videos))
	}
}

func TestListAllContent_OnlyPublishedContent(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := NewContentStore()

	// Add only draft articles
	cs.Articles.Mu.Lock()
	cs.Articles.Articles = map[string]*models.Article{
		"art-1": {
			ID:        "art-1",
			Title:     "Draft Article",
			Slug:      "draft-article",
			Status:    models.ArticleStatusDraft,
			CreatedAt: models.Now(),
		},
	}
	cs.Articles.BySlug = map[string]string{"art-1": "art-1"}
	cs.Articles.Mu.Unlock()

	// Videos don't have draft status - they default to published
	// So we'll test with an empty store
	cs.Videos.Mu.Lock()
	cs.Videos.Videos = map[string]*models.Video{}
	cs.Videos.Mu.Unlock()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/content", nil)

	cs.ListAllContent(c)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	articles := response["articles"].([]interface{})
	videos := response["videos"].([]interface{})

	// Should return empty since only drafts exist
	if len(articles) != 0 {
		t.Errorf("Expected 0 articles (only drafts), got %d", len(articles))
	}
	if len(videos) != 0 {
		t.Errorf("Expected 0 videos, got %d", len(videos))
	}
}

func TestGetAuthorName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := NewContentStore()

	// Test empty userStore
	if name := cs.GetAuthorName("user-1"); name != "" {
		t.Errorf("Expected empty name for nil userStore, got %s", name)
	}

	// Test empty authorID
	if name := cs.GetAuthorName(""); name != "" {
		t.Errorf("Expected empty name for empty authorID, got %s", name)
	}

	// Set up user store
	userStore := map[string]*UserInfo{
		"user-1": {ID: "user-1", FullName: "Dr. John Doe"},
		"user-2": {ID: "user-2", FullName: "Nurse Jane"},
	}
	cs.SetUserStore(userStore)

	// Test with existing user
	if name := cs.GetAuthorName("user-1"); name != "Dr. John Doe" {
		t.Errorf("Expected 'Dr. John Doe', got %s", name)
	}

	// Test with another existing user
	if name := cs.GetAuthorName("user-2"); name != "Nurse Jane" {
		t.Errorf("Expected 'Nurse Jane', got %s", name)
	}

	// Test with non-existing user
	if name := cs.GetAuthorName("user-3"); name != "" {
		t.Errorf("Expected empty name for non-existing user, got %s", name)
	}
}

func TestAddUserToStore(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := NewContentStore()

	// Add initial user
	user1 := &UserInfo{ID: "user-1", FullName: "Dr. John Doe"}
	cs.AddUserToStore(user1)

	if name := cs.GetAuthorName("user-1"); name != "Dr. John Doe" {
		t.Errorf("Expected 'Dr. John Doe', got %s", name)
	}

	// Add another user
	user2 := &UserInfo{ID: "user-2", FullName: "Nurse Jane"}
	cs.AddUserToStore(user2)

	if name := cs.GetAuthorName("user-2"); name != "Nurse Jane" {
		t.Errorf("Expected 'Nurse Jane', got %s", name)
	}

	// Verify user1 still exists
	if name := cs.GetAuthorName("user-1"); name != "Dr. John Doe" {
		t.Errorf("Expected 'Dr. John Doe' to still exist, got %s", name)
	}

	// Test nil user (should not panic)
	cs.AddUserToStore(nil)
}

func TestListAllContent_WithAuthorNames(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := NewContentStore()

	// Set up user store
	userStore := map[string]*UserInfo{
		"author-1": {ID: "author-1", FullName: "Dr. John Doe"},
		"author-2": {ID: "author-2", FullName: "Dr. Jane Smith"},
	}
	cs.SetUserStore(userStore)

	// Add test articles with authors
	cs.Articles.Mu.Lock()
	cs.Articles.Articles = map[string]*models.Article{
		"art-1": {
			ID:         "art-1",
			Title:      "Article by Dr. John",
			Slug:       "article-1",
			Excerpt:    "Content 1",
			CategoryID: "cat-1",
			Status:     models.ArticleStatusPublished,
			AuthorID:   "author-1",
			CreatedAt:  models.Now(),
		},
		"art-2": {
			ID:         "art-2",
			Title:      "Article by Dr. Jane",
			Slug:       "article-2",
			Excerpt:    "Content 2",
			CategoryID: "cat-1",
			Status:     models.ArticleStatusPublished,
			AuthorID:   "author-2",
			CreatedAt:  models.Now(),
		},
		"art-3": {
			ID:         "art-3",
			Title:      "Article with no author",
			Slug:       "article-3",
			Excerpt:    "Content 3",
			CategoryID: "cat-1",
			Status:     models.ArticleStatusPublished,
			AuthorID:   "", // No author
			CreatedAt:  models.Now(),
		},
	}
	cs.Articles.BySlug = map[string]string{"art-1": "art-1", "art-2": "art-2", "art-3": "art-3"}
	cs.Articles.ByCategory = map[string][]string{"cat-1": {"art-1", "art-2", "art-3"}}
	cs.Articles.Mu.Unlock()

	cs.Videos.Mu.Lock()
	cs.Videos.Videos = map[string]*models.Video{}
	cs.Videos.Mu.Unlock()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/content", nil)

	cs.ListAllContent(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	articles := response["articles"].([]interface{})

	if len(articles) != 3 {
		t.Errorf("Expected 3 articles, got %d", len(articles))
	}

	// Check author names in response
	authorNames := make(map[string]string)
	for _, art := range articles {
		a := art.(map[string]interface{})
		id := a["id"].(string)
		authorName, _ := a["authorName"].(string)
		authorNames[id] = authorName
	}

	if authorNames["art-1"] != "Dr. John Doe" {
		t.Errorf("Expected 'Dr. John Doe' for art-1, got %s", authorNames["art-1"])
	}
	if authorNames["art-2"] != "Dr. Jane Smith" {
		t.Errorf("Expected 'Dr. Jane Smith' for art-2, got %s", authorNames["art-2"])
	}
	if authorNames["art-3"] != "" {
		t.Errorf("Expected empty authorName for art-3, got %s", authorNames["art-3"])
	}
}

func TestGetContentAnalytics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := NewContentStore()

	// Set up test data with attachment counts
	cs.Articles.Mu.Lock()
	cs.Articles.Articles = map[string]*models.Article{
		"art-1": {
			ID:              "art-1",
			Title:           "Most Popular Article",
			Slug:            "most-popular",
			Status:          models.ArticleStatusPublished,
			AttachmentCount: 25,
			CreatedAt:       models.Now(),
		},
		"art-2": {
			ID:              "art-2",
			Title:           "Second Article",
			Slug:            "second-article",
			Status:          models.ArticleStatusPublished,
			AttachmentCount: 15,
			CreatedAt:       models.Now(),
		},
		"art-3": {
			ID:              "art-3",
			Title:           "Article Never Attached",
			Slug:            "never-attached",
			Status:          models.ArticleStatusPublished,
			AttachmentCount: 0,
			CreatedAt:       models.Now(),
		},
	}
	cs.Articles.BySlug = map[string]string{
		"art-1": "art-1",
		"art-2": "art-2",
		"art-3": "art-3",
	}
	cs.Articles.Mu.Unlock()

	cs.Videos.Mu.Lock()
	cs.Videos.Videos = map[string]*models.Video{
		"vid-1": {
			ID:              "vid-1",
			Title:           "Most Popular Video",
			YouTubeID:       "abc123",
			Status:          models.VideoStatusPublished,
			AttachmentCount: 30,
			CreatedAt:       models.Now(),
		},
		"vid-2": {
			ID:              "vid-2",
			Title:           "Video Never Attached",
			YouTubeID:       "def456",
			Status:          models.VideoStatusPublished,
			AttachmentCount: 0,
			CreatedAt:       models.Now(),
		},
	}
	cs.Videos.Mu.Unlock()

	tests := []struct {
		name           string
		role           string
		expectedStatus int
		checkFunc      func(t *testing.T, response map[string]interface{})
	}{
		{
			name:           "admin can access analytics",
			role:           "admin",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, response map[string]interface{}) {
				data := response["data"].(map[string]interface{})

				articles := data["articles"].([]interface{})
				if len(articles) != 3 {
					t.Errorf("Expected 3 articles, got %d", len(articles))
				}

				videos := data["videos"].([]interface{})
				if len(videos) != 2 {
					t.Errorf("Expected 2 videos, got %d", len(videos))
				}

				topContent := data["topContent"].([]interface{})
				if len(topContent) != 5 {
					t.Errorf("Expected 5 top content items, got %d", len(topContent))
				}
			},
		},
		{
			name:           "superadmin can access analytics",
			role:           "superadmin",
			expectedStatus: http.StatusOK,
			checkFunc: func(t *testing.T, response map[string]interface{}) {
				// Same as admin
				data := response["data"].(map[string]interface{})
				topContent := data["topContent"].([]interface{})
				if len(topContent) != 5 {
					t.Errorf("Expected 5 top content items, got %d", len(topContent))
				}
			},
		},
		{
			name:           "volunteer cannot access analytics",
			role:           "volunteer",
			expectedStatus: http.StatusForbidden,
			checkFunc:      nil,
		},
		{
			name:           "no role returns forbidden",
			role:           "",
			expectedStatus: http.StatusForbidden,
			checkFunc:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/analytics/content", nil)

			if tt.role != "" {
				c.Set("role", tt.role)
			}

			cs.GetContentAnalytics(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkFunc != nil {
				var response map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				tt.checkFunc(t, response)
			}
		})
	}
}

func TestGetContentAnalytics_TopContentSorted(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := NewContentStore()

	// Set up test data - more than 10 items
	cs.Articles.Mu.Lock()
	for i := 1; i <= 12; i++ {
		cs.Articles.Articles[fmt.Sprintf("art-%d", i)] = &models.Article{
			ID:              fmt.Sprintf("art-%d", i),
			Title:           fmt.Sprintf("Article %d", i),
			Slug:            fmt.Sprintf("article-%d", i),
			Status:          models.ArticleStatusPublished,
			AttachmentCount: i, // 1 to 12
			CreatedAt:       models.Now(),
		}
	}
	cs.Articles.BySlug = make(map[string]string)
	for i := 1; i <= 12; i++ {
		cs.Articles.BySlug[fmt.Sprintf("article-%d", i)] = fmt.Sprintf("art-%d", i)
	}
	cs.Articles.Mu.Unlock()

	cs.Videos.Mu.Lock()
	for i := 1; i <= 5; i++ {
		cs.Videos.Videos[fmt.Sprintf("vid-%d", i)] = &models.Video{
			ID:              fmt.Sprintf("vid-%d", i),
			Title:           fmt.Sprintf("Video %d", i),
			YouTubeID:       fmt.Sprintf("yt%d", i),
			Status:          models.VideoStatusPublished,
			AttachmentCount: i + 20, // 21 to 25
			CreatedAt:       models.Now(),
		}
	}
	cs.Videos.Mu.Unlock()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/analytics/content", nil)
	c.Set("role", "admin")

	cs.GetContentAnalytics(c)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data := response["data"].(map[string]interface{})
	topContent := data["topContent"].([]interface{})

	// Should only return top 10
	if len(topContent) != 10 {
		t.Errorf("Expected 10 top content items, got %d", len(topContent))
	}

	// Check sorting - should be highest first (Video 5 with 25, then Video 4 with 24, etc.)
	expectedOrder := []string{"vid-5", "vid-4", "vid-3", "vid-2", "vid-1", "art-12", "art-11", "art-10", "art-9", "art-8"}
	for i, item := range topContent {
		itemMap := item.(map[string]interface{})
		if itemMap["id"] != expectedOrder[i] {
			t.Errorf("Expected item %d to be %s, got %s", i+1, expectedOrder[i], itemMap["id"])
		}
	}
}

func TestGetContentAnalytics_ZeroAttachments(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := NewContentStore()

	// Set up articles with zero attachment count
	cs.Articles.Mu.Lock()
	cs.Articles.Articles = map[string]*models.Article{
		"art-1": {
			ID:              "art-1",
			Title:           "Never Attached Article",
			Slug:            "never-attached",
			Status:          models.ArticleStatusPublished,
			AttachmentCount: 0,
			CreatedAt:       models.Now(),
		},
	}
	cs.Articles.BySlug = map[string]string{"art-1": "art-1"}
	cs.Articles.Mu.Unlock()

	cs.Videos.Mu.Lock()
	cs.Videos.Videos = map[string]*models.Video{
		"vid-1": {
			ID:              "vid-1",
			Title:           "Never Attached Video",
			YouTubeID:       "xyz123",
			Status:          models.VideoStatusPublished,
			AttachmentCount: 0,
			CreatedAt:       models.Now(),
		},
	}
	cs.Videos.Mu.Unlock()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/analytics/content", nil)
	c.Set("role", "admin")

	cs.GetContentAnalytics(c)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data := response["data"].(map[string]interface{})
	articles := data["articles"].([]interface{})
	videos := data["videos"].([]interface{})
	topContent := data["topContent"].([]interface{})

	if len(articles) != 1 {
		t.Errorf("Expected 1 article, got %d", len(articles))
	}
	if len(videos) != 1 {
		t.Errorf("Expected 1 video, got %d", len(videos))
	}
	if len(topContent) != 2 {
		t.Errorf("Expected 2 top content items (all with 0), got %d", len(topContent))
	}

	// Verify zero counts
	for _, item := range articles {
		a := item.(map[string]interface{})
		if a["attachmentCount"].(float64) != 0 {
			t.Errorf("Expected attachmentCount 0, got %v", a["attachmentCount"])
		}
	}
	for _, item := range videos {
		v := item.(map[string]interface{})
		if v["attachmentCount"].(float64) != 0 {
			t.Errorf("Expected attachmentCount 0, got %v", v["attachmentCount"])
		}
	}
}

func TestGetContentAnalytics_EmptyStore(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cs := NewContentStore()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/analytics/content", nil)
	c.Set("role", "admin")

	cs.GetContentAnalytics(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data := response["data"].(map[string]interface{})
	articles := data["articles"].([]interface{})
	videos := data["videos"].([]interface{})
	topContent := data["topContent"].([]interface{})

	if len(articles) != 0 {
		t.Errorf("Expected 0 articles, got %d", len(articles))
	}
	if len(videos) != 0 {
		t.Errorf("Expected 0 videos, got %d", len(videos))
	}
	if len(topContent) != 0 {
		t.Errorf("Expected 0 top content, got %d", len(topContent))
	}
}
