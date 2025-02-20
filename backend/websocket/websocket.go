package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocket Upgrader to upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

var clients = make(map[*websocket.Conn]bool) // Keeps track of all active WebSocket clients
var mu sync.Mutex                            // Mutex to handle concurrent access to the clients map

// Handle WebSocket Connections
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Add the new connection to the clients map
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	// Keep the connection alive and listen for incoming messages (optional)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			mu.Lock()
			delete(clients, conn) // Remove client if connection is closed
			mu.Unlock()
			break
		}
	}
}

// Broadcast a message to all connected clients
func Broadcast(message string) {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			client.Close()
			delete(clients, client) // Remove broken connection
		}
	}
}
