package websocket

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Store the WebSocket connections
var clients = make(map[*websocket.Conn]bool) // Keep track of connected clients
var broadcast = make(chan string)            // Channel to broadcast messages

// Mutex to avoid concurrent access to the clients map
var mutex sync.Mutex

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

// Handle WebSocket connections
func HandleConnections(c *gin.Context) {
	// Upgrade initial GET request to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Register the new client
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	// Listen for incoming WebSocket messages and handle broadcasts
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn) // Unregister client on error
			mutex.Unlock()
			break
		}

		// Broadcast the received message back to the sender (echo)
		if err := conn.WriteMessage(messageType, msg); err != nil {
			mutex.Lock()
			delete(clients, conn) // Unregister client on error
			mutex.Unlock()
			break
		}
	}

	// If the connection is closed, unregister it
	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()
}

// Broadcast a message to all WebSocket clients
func BroadcastTaskCreated(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Broadcast to all clients
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("Error broadcasting task creation:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// Broadcast task status update to all WebSocket clients
func BroadcastTaskUpdated(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	// Broadcast to all clients
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("Error broadcasting task update:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func BroadcastMessage(message string) {
	mutex.Lock()
	defer mutex.Unlock()

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("❌ WebSocket Send Error:", err)
			client.Close()
			delete(clients, client)
		}
	}
}
