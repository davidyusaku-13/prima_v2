package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/golang-jwt/jwt/v5"
)

const (
	dataFile       = "data/patients.json"
	usersDataFile  = "data/users.json"
	jwtSecretFile  = "data/jwt_secret.txt"
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

// GOWA Configuration
func getGowaEndpoint() string { return getEnv("GOWA_ENDPOINT", "http://localhost:3000") }
func getGowaUser() string     { return getEnv("GOWA_USER", "admin") }
func getGowaPass() string     { return getEnv("GOWA_PASS", "password123") }

// Format phone number for WhatsApp (Indonesian format: 08xx -> 628xx)
func formatWhatsAppNumber(phone string) string {
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, " ", "")

	if idx := strings.Index(phone, "@"); idx != -1 {
		phone = phone[:idx]
	}

	if strings.HasPrefix(phone, "08") {
		phone = "62" + phone[1:]
	} else if strings.HasPrefix(phone, "+62") {
		phone = "62" + phone[3:]
	} else if strings.HasPrefix(phone, "8") {
		phone = "62" + phone
	}

	return phone + "@s.whatsapp.net"
}

// Models
type Recurrence struct {
	Frequency  string `json:"frequency"`
	Interval   int    `json:"interval"`
	DaysOfWeek []int  `json:"daysOfWeek"`
	EndDate    string `json:"endDate,omitempty"`
}

type Reminder struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     string     `json:"dueDate,omitempty"`
	Priority    string     `json:"priority"`
	Completed   bool       `json:"completed"`
	Recurrence  Recurrence `json:"recurrence"`
	Notified    bool       `json:"notified"`
}

type Patient struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Email     string     `json:"email,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	Reminders []*Reminder `json:"reminders,omitempty"`
	CreatedBy string     `json:"createdBy,omitempty"`
}

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

type PatientStore struct {
	mu       sync.RWMutex
	patients map[string]*Patient
}

type UserStore struct {
	mu     sync.RWMutex
	users  map[string]*User
	byName map[string]string // username -> userID
}

var (
	store    = PatientStore{patients: make(map[string]*Patient)}
	userStore = UserStore{users: make(map[string]*User), byName: make(map[string]string)}
)

func main() {
	// Ensure data directory exists
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Load .env file
	loadEnvFile()

	// Load existing data
	loadData()
	loadUsers()

	// Create default superadmin if not exists
	createDefaultSuperadmin()

	// Start reminder checker goroutine
	go checkReminders()

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check (public)
	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth routes (public)
	router.POST("/api/auth/register", register)
	router.POST("/api/auth/login", login)

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

		// Reminder routes
		api.POST("/patients/:id/reminders", createReminder)
		api.PUT("/patients/:id/reminders/:reminderId", updateReminder)
		api.POST("/patients/:id/reminders/:reminderId/toggle", toggleReminder)
		api.DELETE("/patients/:id/reminders/:reminderId", deleteReminder)

		// User management routes (superadmin only)
		api.GET("/users", requireRole(RoleSuperadmin), getUsers)
		api.PUT("/users/:id/role", requireRole(RoleSuperadmin), updateUserRole)
		api.DELETE("/users/:id", requireRole(RoleSuperadmin), deleteUser)
	}

	router.Run(":8080")
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

// Data persistence
func loadData() {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}

	var patients map[string]*Patient
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
		patients := make(map[string]*Patient)
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

// WhatsApp
func sendWhatsAppMessage(phone, message string) error {
	endpoint := getGowaEndpoint() + "/send/message"

	payload := map[string]interface{}{
		"phone":   phone,
		"message": message,
	}

	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	auth := getGowaUser() + ":" + getGowaPass()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	return nil
}

// Reminder checker
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

				dueTime, err := time.ParseInLocation("2006-01-02T15:04", reminder.DueDate, localLoc)
				if err != nil {
					dueTime, err = time.ParseInLocation("2006-01-02T15:04:05", reminder.DueDate, localLoc)
					if err != nil {
						continue
					}
				}

				if !now.Before(dueTime) && now.Before(dueTime.Add(5*time.Minute)) {
					whatsappPhone := formatWhatsAppNumber(patient.Phone)
					message := fmt.Sprintf("Reminder: %s\n\n", reminder.Title)
					if reminder.Description != "" {
						message += fmt.Sprintf("Description: %s\n\n", reminder.Description)
					}
					message += fmt.Sprintf("Priority: %s\n", reminder.Priority)
					if reminder.Recurrence.Frequency != "none" {
						message += fmt.Sprintf("Repeats: %s", reminder.Recurrence.Frequency)
					}

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

			sendWhatsAppMessage(n.phone, message)
		}
	}
}

// Patient handlers
func getPatients(c *gin.Context) {
	userID := c.GetString("userID")
	role := c.GetString("role")

	store.mu.RLock()
	patients := make([]*Patient, 0, len(store.patients))
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

	if len(req.Phone) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid phone number"})
		return
	}

	userID := c.GetString("userID")

	store.mu.Lock()
	patient := &Patient{
		ID:        generateID(),
		Name:      req.Name,
		Phone:     req.Phone,
		Email:     req.Email,
		Notes:     req.Notes,
		Reminders: make([]*Reminder, 0),
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

// Reminder handlers
func createReminder(c *gin.Context) {
	patientID := c.Param("id")
	userID := c.GetString("userID")
	role := c.GetString("role")

	var req struct {
		Title       string     `json:"title" binding:"required"`
		Description string     `json:"description"`
		DueDate     string     `json:"dueDate"`
		Priority    string     `json:"priority"`
		Recurrence  Recurrence `json:"recurrence"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.mu.Lock()
	patient, exists := store.patients[patientID]
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

	reminder := &Reminder{
		ID:          generateID(),
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
		Completed:   false,
		Recurrence:  req.Recurrence,
		Notified:    false,
	}
	patient.Reminders = append(patient.Reminders, reminder)
	store.mu.Unlock()

	saveData()
	c.JSON(http.StatusCreated, reminder)
}

func updateReminder(c *gin.Context) {
	patientID := c.Param("id")
	reminderID := c.Param("reminderId")
	userID := c.GetString("userID")
	role := c.GetString("role")

	var req struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		DueDate     string     `json:"dueDate"`
		Priority    string     `json:"priority"`
		Recurrence  Recurrence `json:"recurrence"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	store.mu.Lock()
	patient, exists := store.patients[patientID]
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

	for _, r := range patient.Reminders {
		if r.ID == reminderID {
			if req.Title != "" {
				r.Title = req.Title
			}
			r.Description = req.Description
			r.DueDate = req.DueDate
			r.Priority = req.Priority
			r.Recurrence = req.Recurrence
			if req.DueDate != "" && req.DueDate != r.DueDate {
				r.Notified = false
			}
			store.mu.Unlock()
			saveData()
			c.JSON(http.StatusOK, r)
			return
		}
	}
	store.mu.Unlock()
	c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
}

func toggleReminder(c *gin.Context) {
	patientID := c.Param("id")
	reminderID := c.Param("reminderId")
	userID := c.GetString("userID")
	role := c.GetString("role")

	store.mu.Lock()
	patient, exists := store.patients[patientID]
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

	for _, r := range patient.Reminders {
		if r.ID == reminderID {
			r.Completed = !r.Completed
			if !r.Completed {
				r.Notified = false
			}
			store.mu.Unlock()
			saveData()
			c.JSON(http.StatusOK, r)
			return
		}
	}
	store.mu.Unlock()
	c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
}

func deleteReminder(c *gin.Context) {
	patientID := c.Param("id")
	reminderID := c.Param("reminderId")
	userID := c.GetString("userID")
	role := c.GetString("role")

	store.mu.Lock()
	patient, exists := store.patients[patientID]
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

	for i, r := range patient.Reminders {
		if r.ID == reminderID {
			patient.Reminders = append(patient.Reminders[:i], patient.Reminders[i+1:]...)
			store.mu.Unlock()
			saveData()
			c.JSON(http.StatusOK, gin.H{"message": "reminder deleted"})
			return
		}
	}
	store.mu.Unlock()
	c.JSON(http.StatusNotFound, gin.H{"error": "reminder not found"})
}

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
