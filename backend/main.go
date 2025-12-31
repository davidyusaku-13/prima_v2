package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/davidyusaku-13/prima_v2/config"
	"github.com/davidyusaku-13/prima_v2/handlers"
	"github.com/davidyusaku-13/prima_v2/models"
	"github.com/davidyusaku-13/prima_v2/services"
	"github.com/davidyusaku-13/prima_v2/utils"
)

const (
	dataFile           = "data/patients.json"
	usersDataFile      = "data/users.json"
	jwtSecretFile      = "data/jwt_secret.txt"
	categoriesDataFile = "data/categories.json"
	articlesDataFile   = "data/articles.json"
	videosDataFile     = "data/videos.json"
)

const tokenExpiry = 24 * 7 * time.Hour // 1 week

// Get environment variable with default
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Load .env file if it exists
func loadEnvFile() {
	file, err := os.ReadFile(".env")
	if err != nil {
		return
	}

	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")

		if _, exists := os.LookupEnv(key); !exists {
			os.Setenv(key, value)
		}
	}
}

// Get or create JWT secret
func getJWTSecret() string {
	data, err := os.ReadFile(jwtSecretFile)
	if err == nil {
		return strings.TrimSpace(string(data))
	}

	// Generate new secret
	secret := base64.StdEncoding.EncodeToString([]byte(time.Now().Format("2006-01-02T15:04:05Z07:00")))
	os.WriteFile(jwtSecretFile, []byte(secret), 0644)
	return secret
}

// Hash password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.StdEncoding.EncodeToString(hash[:])
}

// Verify password
func verifyPassword(password, hash string) bool {
	return hashPassword(password) == hash
}

// Role constants
type Role string

const (
	RoleSuperadmin Role = "superadmin"
	RoleAdmin      Role = "admin"
	RoleVolunteer  Role = "volunteer"
)

// JWT claims
type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Role     Role   `json:"role"`
	jwt.RegisteredClaims
}

// Generate JWT token
func generateToken(userID, username string, role Role) (string, error) {
	secret := getJWTSecret()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Verify JWT token
func verifyToken(tokenString string) (*Claims, error) {
	secret := getJWTSecret()
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}


// Patient and Reminder types are now in models/patient.go

// User is stored in memory and file (includes password)
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"fullName,omitempty"`
	Password  string `json:"password"`
	Role      Role   `json:"role"`
	CreatedAt string `json:"createdAt"`
}

// UserResponse is used for API responses (excludes password)
type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FullName  string `json:"fullName,omitempty"`
	Role      Role   `json:"role"`
	CreatedAt string `json:"createdAt"`
}

// MarshalJSON for User - excludes password from API responses
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		FullName:  u.FullName,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	})
}

// PatientStore wraps models.PatientStore for backward compatibility
type PatientStore struct {
	mu       sync.RWMutex
	patients map[string]*models.Patient
}

type UserStore struct {
	mu     sync.RWMutex
	users  map[string]*User
	byName map[string]string // username -> userID
}

var (
	store            = PatientStore{patients: make(map[string]*models.Patient)}
	userStore        = UserStore{users: make(map[string]*User), byName: make(map[string]string)}
	contentStore     *handlers.ContentStore
	appConfig        *config.Config
	appLogger        *slog.Logger
	gowaClient       *services.GOWAClient
	reminderHandler  *handlers.ReminderHandler
	patientStore     *models.PatientStore
	scheduler        *services.ReminderScheduler
	webhookHandler   *handlers.WebhookHandler
	sseHandler       *handlers.SSEHandler
	analyticsHandler *handlers.AnalyticsHandler
	healthHandler    *handlers.HealthHandler
)

func main() {
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Load .env file
	loadEnvFile()

	// Load configuration from YAML
	appConfig = config.LoadOrDefault("config.yaml")

	// Initialize structured logger
	utils.InitDefaultLogger(appConfig.Logging.Level, appConfig.Logging.Format)
	appLogger = utils.DefaultLogger

	// Load existing data
	loadData()
	loadUsers()

	// Initialize patient store with shared models
	patientStore = models.NewPatientStore(saveData)
	patientStore.Patients = store.patients

	// Initialize GOWA client with circuit breaker
	gowaClient = services.NewGOWAClientFromConfig(appConfig, appLogger)

	// Initialize and load content store (before reminder handler)
	contentStore = handlers.NewContentStore()
	contentStore.LoadContentData()

	// Set user store for author name resolution
	userMap := make(map[string]*handlers.UserInfo)
	for id, user := range userStore.users {
		userMap[id] = &handlers.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		}
	}
	contentStore.SetUserStore(userMap)

	// Initialize reminder handler with new architecture (after contentStore)
	reminderHandler = handlers.NewReminderHandler(
		patientStore,
		appConfig,
		gowaClient,
		appLogger,
		generateID,
		contentStore, // Pass contentStore for attachment validation
	)

	// Create default superadmin if not exists
	createDefaultSuperadmin()

	// Initialize and start reminder scheduler for quiet hours
	scheduler = services.NewReminderScheduler(patientStore, gowaClient, appConfig, appLogger)
	scheduler.Start()

	// Initialize webhook handler for GOWA delivery status updates
	webhookHandler = handlers.NewWebhookHandler(patientStore, appConfig, appLogger)

	// Initialize SSE handler for real-time delivery status updates
	sseHandler = handlers.NewSSEHandler(appConfig, appLogger)

	// Connect SSE handler to webhook handler for broadcasting
	webhookHandler.SetSSEHandler(sseHandler)

	// Connect SSE handler to reminder handler for manual send broadcasts
	reminderHandler.SetSSEHandler(sseHandler)

	// Connect SSE handler to scheduler for auto-send broadcasts
	scheduler.SetSSEHandler(sseHandler)

	// Initialize analytics handler for delivery statistics
	analyticsHandler = handlers.NewAnalyticsHandler(patientStore)

	// Initialize health handler for system health monitoring
	healthHandler = handlers.NewHealthHandler(patientStore, gowaClient)

	// Start health check goroutine (runs every 60 seconds)
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if gowaClient != nil {
				connected := gowaClient.IsAvailable()
				healthHandler.UpdateGOWAPing(connected)
			}
		}
	}()

	// Initial health check
	go func() {
		if gowaClient != nil {
			connected := gowaClient.IsAvailable()
			healthHandler.UpdateGOWAPing(connected)
		}
	}()

	// Start reminder checker goroutine (DISABLED - replaced by ReminderScheduler auto-send)
	// The ReminderScheduler now handles all reminder sending: scheduled, retry, and auto-send
	// go checkReminders()

	router := gin.Default()

	// Configure CORS using config
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{appConfig.Server.CORSOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Serve uploaded files
	router.Static("/uploads", "./uploads")

	// Health check (public)
	router.GET("/api/health", healthHandler.GetHealth)
	// Health check detailed (admin only)
	router.GET("/api/health/detailed", authMiddleware(), healthHandler.GetHealthDetailed)

	// GOWA webhook endpoint (public - no auth, uses HMAC validation)
	router.POST("/api/webhook/gowa", webhookHandler.HandleGOWAWebhook)

	// Config endpoints (public)
	router.GET("/api/config/disclaimer", getDisclaimerConfig)
	router.GET("/api/config/quiet-hours", getQuietHoursConfig)

	// Auth routes (public)
	router.POST("/api/auth/register", register)
	router.POST("/api/auth/login", login)

	// Public content routes (no auth required)
	contentPublic := router.Group("/api")
	{
		// Categories
		contentPublic.GET("/categories", contentStore.ListCategories)
		contentPublic.GET("/categories/:type", contentStore.GetCategoriesByType)

		// Articles
		contentPublic.GET("/articles", contentStore.ListArticles)
		contentPublic.GET("/articles/:slug", contentStore.GetArticle)

		// Videos
		contentPublic.GET("/videos", contentStore.ListVideos)

		// Combined content listing (for content picker)
		contentPublic.GET("/content", contentStore.ListAllContent)

		// Popular content (for content suggestions)
		contentPublic.GET("/content/popular", contentStore.GetPopularContent)
	}

	// SSE endpoint (uses query parameter auth, not header auth)
	// Must be outside protected group because EventSource API doesn't support custom headers
	router.GET("/api/sse/delivery-status", sseAuthMiddleware(), sseHandler.HandleDeliveryStatusSSE)

	// Protected routes
	api := router.Group("/api")
	api.Use(authMiddleware())
	{
		// Current user
		api.GET("/auth/me", getCurrentUser)

		// Patient routes
		api.GET("/patients", getPatients)
		api.POST("/patients", createPatient)
		api.GET("/patients/:id", getPatient)
		api.PUT("/patients/:id", updatePatient)
		api.DELETE("/patients/:id", deletePatient)

		// Reminder routes - using new handler with circuit breaker
		api.POST("/patients/:id/reminders", reminderHandler.Create)
		api.PUT("/patients/:id/reminders/:reminderId", reminderHandler.Update)
		api.POST("/patients/:id/reminders/:reminderId/toggle", reminderHandler.Toggle)
		api.DELETE("/patients/:id/reminders/:reminderId", reminderHandler.Delete)
		api.POST("/patients/:id/reminders/:reminderId/send", reminderHandler.Send)
		api.GET("/reminders/:id/status", reminderHandler.GetReminderStatus)
		api.POST("/reminders/:id/retry", reminderHandler.RetryReminder)

		// User management routes (superadmin only)
		api.GET("/users", requireRole(RoleSuperadmin), getUsers)
		api.GET("/users/:id", getUserByID) // Get user by ID (for author lookup)
		api.PUT("/users/:id/role", requireRole(RoleSuperadmin), updateUserRole)
		api.DELETE("/users/:id", requireRole(RoleSuperadmin), deleteUser)

		// CMS routes (admin+)
		// Categories
		api.POST("/categories", requireRole(RoleAdmin, RoleSuperadmin), contentStore.CreateCategory)

		// Articles
		api.POST("/articles", requireRole(RoleAdmin, RoleSuperadmin), contentStore.CreateArticle)
		api.PUT("/articles/:id", requireRole(RoleAdmin, RoleSuperadmin), contentStore.UpdateArticle)
		api.DELETE("/articles/:id", requireRole(RoleAdmin, RoleSuperadmin), contentStore.DeleteArticle)

		// Videos
		api.POST("/videos", requireRole(RoleAdmin, RoleSuperadmin), contentStore.CreateVideo)
		api.DELETE("/videos/:id", requireRole(RoleAdmin, RoleSuperadmin), contentStore.DeleteVideo)

		// Content attachment count
		api.POST("/content/:type/:id/increment-attachment", contentStore.IncrementAttachmentCount)

		// Upload
		api.POST("/upload/image", requireRole(RoleAdmin, RoleSuperadmin), contentStore.UploadImage)

		// Dashboard stats
		api.GET("/dashboard/stats", requireRole(RoleAdmin, RoleSuperadmin), contentStore.GetDashboardStats)

		// Analytics - Content attachment statistics
		api.GET("/analytics/content", requireRole(RoleAdmin, RoleSuperadmin), contentStore.GetContentAnalytics)

		// Analytics - Delivery statistics
		api.GET("/analytics/delivery", requireRole(RoleAdmin, RoleSuperadmin), analyticsHandler.GetDeliveryAnalytics)

		// Analytics - Failed deliveries
		api.GET("/analytics/failed-deliveries", requireRole(RoleAdmin, RoleSuperadmin), analyticsHandler.GetFailedDeliveries)
		api.GET("/analytics/failed-deliveries/export", requireRole(RoleAdmin, RoleSuperadmin), analyticsHandler.ExportFailedDeliveries)
		api.GET("/analytics/failed-deliveries/:id", requireRole(RoleAdmin, RoleSuperadmin), analyticsHandler.GetFailedDeliveryDetail)
	}

	// Create HTTP server for graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.Server.Port),
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("🚀 Server started on port %d", appConfig.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	// Stop the scheduler first
	if scheduler != nil {
		appLogger.Info("Stopping reminder scheduler...")
		scheduler.Stop()
		appLogger.Info("Reminder scheduler stopped")
	}

	// Close all SSE connections before shutting down HTTP server
	if sseHandler != nil {
		appLogger.Info("Closing SSE connections...")
		sseHandler.Shutdown()
	}

	// Create context with timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server exited gracefully")
}

// Auth middleware
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}

		// Extract "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := verifyToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", string(claims.Role))
		c.Next()
	}
}

// SSE auth middleware - accepts token from query parameter
// EventSource API doesn't support custom headers, so token must be in query string
func sseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get token from query parameter first (for SSE)
		token := c.Query("token")
		if token == "" {
			// Fallback to Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.SSEvent("error", `{"error": "authorization required"}`)
				c.Writer.Flush()
				c.Abort()
				return
			}

			// Extract "Bearer <token>"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.SSEvent("error", `{"error": "invalid authorization format"}`)
				c.Writer.Flush()
				c.Abort()
				return
			}
			token = parts[1]
		}

		claims, err := verifyToken(token)
		if err != nil {
			c.SSEvent("error", `{"error": "invalid or expired token"}`)
			c.Writer.Flush()
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", string(claims.Role))
		c.Next()
	}
}

// Role check middleware
func requireRole(allowedRoles ...Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
			c.Abort()
			return
		}

		userRole := Role(roleVal.(string))
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
	}
}

// Auth handlers
func register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=3,max=30"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"fullName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: username must be 3-30 characters, password must be at least 6 characters"})
		return
	}

	userStore.mu.Lock()
	if _, exists := userStore.byName[req.Username]; exists {
		userStore.mu.Unlock()
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	// New users always become volunteers
	user := &User{
		ID:        generateID(),
		Username:  req.Username,
		FullName:  req.FullName,
		Password:  hashPassword(req.Password),
		Role:      RoleVolunteer,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	userStore.users[user.ID] = user
	userStore.byName[user.Username] = user.ID
	userStore.mu.Unlock()

	saveUsers()

	// Sync new user to contentStore for author attribution
	if contentStore != nil {
		contentStore.AddUserToStore(&handlers.UserInfo{
			ID:        user.ID,
			Username:  user.Username,
			FullName:  user.FullName,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		})
	}

	// Generate token for immediate login
	token, _ := generateToken(user.ID, user.Username, user.Role)

	c.JSON(http.StatusCreated, gin.H{
		"message":   "user registered successfully",
		"userId":    user.ID,
		"username":  user.Username,
		"fullName":  user.FullName,
		"role":      user.Role,
		"token":     token,
		"expiresAt": time.Now().Add(tokenExpiry).Format(time.RFC3339),
	})
}

func login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userStore.mu.RLock()
	userID, exists := userStore.byName[req.Username]
	if !exists {
		userStore.mu.RUnlock()
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	user := userStore.users[userID]
	userStore.mu.RUnlock()

	if !verifyPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := generateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"userId":    user.ID,
		"username":  user.Username,
		"fullName":  user.FullName,
		"role":      user.Role,
		"expiresAt": time.Now().Add(tokenExpiry).Format(time.RFC3339),
	})
}

func getCurrentUser(c *gin.Context) {
	userID := c.GetString("userID")
	username := c.GetString("username")
	role := c.GetString("role")

	userStore.mu.RLock()
	user := userStore.users[userID]
	userStore.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"userId":   userID,
		"username": username,
		"fullName": user.FullName,
		"role":     role,
	})
}

// getUserByID returns user details by ID (for author lookup in content)
func getUserByID(c *gin.Context) {
	id := c.Param("id")

	userStore.mu.RLock()
	user, exists := userStore.users[id]
	userStore.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"fullName":  user.FullName,
		"role":      user.Role,
		"createdAt": user.CreatedAt,
	})
}

// Data persistence
func loadData() {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}

	var patients map[string]*models.Patient
	if err := json.Unmarshal(data, &patients); err != nil {
		return
	}

	store.mu.Lock()
	store.patients = patients
	store.mu.Unlock()
}

func saveData() {
	go func() {
		store.mu.RLock()
		patients := make(map[string]*models.Patient)
		for k, v := range store.patients {
			patients[k] = v
		}
		store.mu.RUnlock()

		data, err := json.MarshalIndent(patients, "", "  ")
		if err != nil {
			return
		}

		tmpFile := dataFile + ".tmp"
		err = os.WriteFile(tmpFile, data, 0644)
		if err != nil {
			return
		}
		os.Rename(tmpFile, dataFile)
	}()
}

func loadUsers() {
	data, err := os.ReadFile(usersDataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}

	var users map[string]*User
	if err := json.Unmarshal(data, &users); err != nil {
		return
	}

	userStore.mu.Lock()
	userStore.users = users
	userStore.byName = make(map[string]string)
	for id, user := range users {
		userStore.byName[user.Username] = id
	}
	userStore.mu.Unlock()
}

func saveUsers() {
	userStore.mu.RLock()
	// Create a map with all fields including password (bypass MarshalJSON)
	users := make(map[string]map[string]interface{})
	for id, u := range userStore.users {
		users[id] = map[string]interface{}{
			"id":        u.ID,
			"username":  u.Username,
			"fullName":  u.FullName,
			"password":  u.Password,
			"role":      u.Role,
			"createdAt": u.CreatedAt,
		}
	}
	userStore.mu.RUnlock()

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return
	}

	tmpFile := usersDataFile + ".tmp"
	err = os.WriteFile(tmpFile, data, 0644)
	if err != nil {
		return
	}
	os.Rename(tmpFile, usersDataFile)
}

// sendWhatsAppMessage sends a WhatsApp message using the GOWA client with circuit breaker
func sendWhatsAppMessage(phone, message string) error {
	if gowaClient == nil {
		return fmt.Errorf("GOWA client not initialized")
	}

	_, err := gowaClient.SendMessage(phone, message)
	return err
}

// Reminder checker - uses structured logging and new phone validation
func checkReminders() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	localLoc := time.Local

	for range ticker.C {
		store.mu.RLock()
		now := time.Now().In(localLoc)

		// Collect reminders to notify (avoid holding lock during HTTP call)
		var toNotify []struct {
			phone    string
			title    string
			desc     string
			priority string
			recFreq  string
		}

		for _, patient := range store.patients {
			if patient.Phone == "" {
				continue
			}

			for _, reminder := range patient.Reminders {
				if reminder.Completed || reminder.Notified || reminder.DueDate == "" {
					continue
				}

				// Skip if already sent via new delivery system (manual send or scheduled)
				if reminder.DeliveryStatus != "" && reminder.DeliveryStatus != models.DeliveryStatusPending {
					continue
				}

				dueTime, err := time.ParseInLocation("2006-01-02T15:04", reminder.DueDate, localLoc)
				if err != nil {
					dueTime, err = time.ParseInLocation("2006-01-02T15:04:05", reminder.DueDate, localLoc)
					if err != nil {
						continue
					}
				}

				if !now.Before(dueTime) && now.Before(dueTime.Add(5*time.Minute)) {
					// Use new phone validation utility
					whatsappPhone := utils.FormatWhatsAppNumber(patient.Phone)

					toNotify = append(toNotify, struct {
						phone    string
						title    string
						desc     string
						priority string
						recFreq  string
					}{whatsappPhone, reminder.Title, reminder.Description, reminder.Priority, reminder.Recurrence.Frequency})
				}
			}
		}
		store.mu.RUnlock()

		// Send notifications without holding any locks
		for _, n := range toNotify {
			message := fmt.Sprintf("Reminder: %s\n\n", n.title)
			if n.desc != "" {
				message += fmt.Sprintf("Description: %s\n\n", n.desc)
			}
			message += fmt.Sprintf("Priority: %s\n", n.priority)
			if n.recFreq != "none" {
				message += fmt.Sprintf("Repeats: %s", n.recFreq)
			}

			if err := sendWhatsAppMessage(n.phone, message); err != nil {
				appLogger.Error("Failed to send WhatsApp reminder",
					"phone", utils.MaskPhone(n.phone),
					"error", err.Error(),
				)
			} else {
				appLogger.Info("WhatsApp reminder sent",
					"phone", utils.MaskPhone(n.phone),
					"title", n.title,
				)
			}
		}
	}
}

// Patient handlers
func getPatients(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")

	store.mu.RLock()
	patients := make([]*models.Patient, 0, len(store.patients))
	for _, p := range store.patients {
		// Superadmin and Admin can see all patients
		// Volunteers can only see patients they created
		if role == string(RoleVolunteer) {
			if p.CreatedBy == userID {
				patients = append(patients, p)
			}
		} else {
			patients = append(patients, p)
		}
	}
	store.mu.RUnlock()
	c.JSON(http.StatusOK, gin.H{"patients": patients})
}

func createPatient(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Phone string `json:"phone" binding:"required"`
		Email string `json:"email"`
		Notes string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use new phone validation utility
	phoneResult := utils.ValidatePhoneNumber(req.Phone)
	if !phoneResult.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": phoneResult.ErrorMessage})
		return
	}

	userID := c.GetString("userID")

	store.mu.Lock()
	patient := &models.Patient{
		ID:        generateID(),
		Name:      req.Name,
		Phone:     phoneResult.Normalized, // Store normalized phone
		Email:     req.Email,
		Notes:     req.Notes,
		Reminders: make([]*models.Reminder, 0),
		CreatedBy: userID,
	}
	store.patients[patient.ID] = patient
	store.mu.Unlock()

	saveData()
	c.JSON(http.StatusCreated, patient)
}

func getPatient(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	role := c.GetString("role")

	store.mu.RLock()
	patient, exists := store.patients[id]
	store.mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	// Check access for volunteers
	if role == string(RoleVolunteer) && patient.CreatedBy != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func updatePatient(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	role := c.GetString("role")

	var req struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
		Notes string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.mu.Lock()
	patient, exists := store.patients[id]
	if !exists {
		store.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	// Check access for volunteers
	if role == string(RoleVolunteer) && patient.CreatedBy != userID {
		store.mu.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	if req.Name != "" {
		patient.Name = req.Name
	}
	if req.Phone != "" {
		patient.Phone = req.Phone
	}
	if req.Email != "" {
		patient.Email = req.Email
	}
	patient.Notes = req.Notes
	store.mu.Unlock()

	saveData()
	c.JSON(http.StatusOK, patient)
}

func deletePatient(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	role := c.GetString("role")

	store.mu.Lock()
	patient, exists := store.patients[id]
	if !exists {
		store.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	// Check access for volunteers
	if role == string(RoleVolunteer) && patient.CreatedBy != userID {
		store.mu.Unlock()
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		return
	}

	delete(store.patients, id)
	store.mu.Unlock()

	saveData()
	c.JSON(http.StatusOK, gin.H{"message": "patient deleted"})
}

// Reminder handlers are now in handlers/reminder.go

// User management handlers
func getUsers(c *gin.Context) {
	userStore.mu.RLock()
	users := make([]gin.H, 0, len(userStore.users))
	for id, user := range userStore.users {
		users = append(users, gin.H{
			"id":        id,
			"username":  user.Username,
			"fullName":  user.FullName,
			"role":      user.Role,
			"createdAt": user.CreatedAt,
		})
	}
	userStore.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func updateUserRole(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("userID")

	// Cannot change own role
	if userID == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot change your own role"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role is required"})
		return
	}

	// Validate role
	if req.Role != string(RoleAdmin) && req.Role != string(RoleVolunteer) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	userStore.mu.Lock()
	user, exists := userStore.users[userID]
	if !exists {
		userStore.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.Role = Role(req.Role)
	userStore.mu.Unlock()
	saveUsers()

	c.JSON(http.StatusOK, gin.H{
		"message": "role updated successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func deleteUser(c *gin.Context) {
	userID := c.Param("id")
	currentUserID := c.GetString("userID")

	// Cannot delete yourself
	if userID == currentUserID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete yourself"})
		return
	}

	userStore.mu.Lock()
	user, exists := userStore.users[userID]
	if !exists {
		userStore.mu.Unlock()
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	delete(userStore.users, userID)
	delete(userStore.byName, user.Username)
	userStore.mu.Unlock()
	saveUsers()

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

func createDefaultSuperadmin() {
	userStore.mu.Lock()
	defer userStore.mu.Unlock()

	if _, exists := userStore.byName["superadmin"]; exists {
		return // Superadmin already exists
	}

	user := &User{
		ID:        "superadmin",
		Username:  "superadmin",
		FullName:  "System Administrator",
		Password:  hashPassword("superadmin"),
		Role:      RoleSuperadmin,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	userStore.users[user.ID] = user
	userStore.byName[user.Username] = user.ID

	// Save with password (bypass MarshalJSON which excludes password for API responses)
	users := map[string]map[string]interface{}{
		user.ID: {
			"id":        user.ID,
			"username":  user.Username,
			"fullName":  user.FullName,
			"password":  user.Password,
			"role":      user.Role,
			"createdAt": user.CreatedAt,
		},
	}
	data, _ := json.MarshalIndent(users, "", "  ")
	os.WriteFile(usersDataFile, data, 0644)
}

func generateID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}

// getDisclaimerConfig returns the disclaimer configuration for frontend sync
func getDisclaimerConfig(c *gin.Context) {
	enabled := false
	if appConfig.Disclaimer.Enabled != nil {
		enabled = *appConfig.Disclaimer.Enabled
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"text":    appConfig.Disclaimer.Text,
			"enabled": enabled,
		},
	})
}

// getQuietHoursConfig returns the quiet hours configuration for frontend
func getQuietHoursConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"start_hour": appConfig.QuietHours.GetStartHour(),
			"end_hour":   appConfig.QuietHours.GetEndHour(),
			"timezone":   appConfig.QuietHours.Timezone,
		},
	})
}
