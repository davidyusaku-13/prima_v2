package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

type Message struct {
	Text string `json:"text"`
}

var (
	items []Message
	mu    sync.Mutex
)

func main() {
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Svelte dev server
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Routes
	router.GET("/api/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Go!",
		})
	})

	router.GET("/api/items", getItems)
	router.POST("/api/items", createItem)
	router.DELETE("/api/items", deleteItem)

	router.Run(":8080") // Backend runs on port 8080
}

func getItems(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func createItem(c *gin.Context) {
	var msg Message
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	items = append(items, msg)
	mu.Unlock()
	c.JSON(http.StatusCreated, gin.H{"message": "Item created", "data": msg})
}

func deleteItem(c *gin.Context) {
	var req struct {
		Index int `json:"index"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mu.Lock()
	if req.Index >= 0 && req.Index < len(items) {
		items = append(items[:req.Index], items[req.Index+1:]...)
	}
	mu.Unlock()
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted"})
}
