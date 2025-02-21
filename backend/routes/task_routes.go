package routes

import (
	"task-manager/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all routes
func SetupRoutes(router *gin.Engine) {
	// Task routes
	router.POST("/tasks", controllers.CreateTask)
	router.GET("/tasks", controllers.GetTasks)
	router.PUT("/tasks/:id/status", controllers.UpdateTaskStatus)
	router.DELETE("/tasks/:id", controllers.DeleteTask)
}
