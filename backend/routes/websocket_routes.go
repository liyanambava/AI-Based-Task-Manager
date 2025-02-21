package routes

import (
	"task-manager/websocket" // Import WebSocket package

	"github.com/gin-gonic/gin"
)

// ✅ Setup WebSocket route
func SetupWebSocketRoutes(router *gin.Engine) {
	router.GET("/ws", websocket.HandleConnections)
}
