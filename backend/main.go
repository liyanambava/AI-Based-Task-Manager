package main

import (
	"log"
	"net/http"
	"time"

	"task-manager/config"
	"task-manager/controllers"
	"task-manager/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// ✅ CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ✅ Connect to MongoDB
	config.ConnectDB()

	// ✅ Authentication Route
	router.POST("/login", func(c *gin.Context) {
		token, err := controllers.GenerateToken("testuser")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	// ✅ Task Routes
	routes.SetupRoutes(router) // This now includes WebSocket setup

	// ✅ Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}
