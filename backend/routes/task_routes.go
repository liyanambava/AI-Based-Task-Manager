package routes

import (
	"net/http"
	"task-manager/controllers"
	"task-manager/websocket" // Import the websocket package

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all routes
func SetupRoutes(router *gin.Engine) {
	// Task routes
	router.POST("/tasks", controllers.CreateTask)
	router.GET("/tasks", controllers.GetTasks)

	// WebSocket route
	router.GET("/ws", gin.WrapF(http.HandlerFunc(websocket.HandleConnections))) // Wrap WebSocket handler to work with Gin
}
