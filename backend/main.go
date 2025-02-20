package main

import (
	"log"
	"net/http"

	"task-manager/config"
	"task-manager/controllers"
	"task-manager/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Connect to MongoDB
	config.ConnectDB()

	// Authentication Route
	router.POST("/login", func(c *gin.Context) {
		token, err := controllers.GenerateToken("testuser")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	// Task Routes
	routes.SetupRoutes(router)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}
