package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Main function
func main() {
	// ✅ Initialize MongoDB (Check if ConnectDB exists in db.go)
	ConnectDB()

	// ✅ Initialize Router
	router := gin.Default()

	// ✅ Public route (Health check)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Task Manager API is running!"})
	})

	// ✅ Authentication routes (Check if GenerateToken exists in auth.go)
	router.POST("/login", GenerateToken)

	// ✅ Task routes (Check if CreateTask & GetTasks exist in routes.go)
	router.POST("/tasks", CreateTask)
	router.GET("/tasks", GetTasks)

	// ✅ Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
