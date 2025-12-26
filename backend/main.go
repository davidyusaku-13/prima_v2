package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
)

type Message struct {
    Text string `json:"text"`
}

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

    router.Run(":8080") // Backend runs on port 8080
}

func getItems(c *gin.Context) {
    items := []string{"Item 1", "Item 2", "Item 3"}
    c.JSON(http.StatusOK, gin.H{"items": items})
}

func createItem(c *gin.Context) {
    var msg Message
    if err := c.BindJSON(&msg); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "Item created", "data": msg})
}