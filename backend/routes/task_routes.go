package routes

import (
	"task-manager/controllers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes initializes all routes
func SetupRoutes(router *gin.Engine) {
	router.POST("/tasks", controllers.CreateTask)
	router.GET("/tasks", controllers.GetTasks)
}
