package service

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Message, 32)
	messages  []Message
	mu        sync.RWMutex
)

func init() {
	go handleMessages()
}

func RegisterClient(ws *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	clients[ws] = true
}

func UnregisterClient(ws *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, ws)
}

func BroadcastMessage(msg Message) {
	mu.Lock()
	messages = append(messages, msg)
	mu.Unlock()

	broadcast <- msg
}

func handleMessages() {
	for {
		msg := <-broadcast
		for _, client := range snapshotClients() {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				UnregisterClient(client)
			}
		}
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	snapshot := append([]Message(nil), messages...)
	mu.RUnlock()

	if snapshot == nil {
		snapshot = make([]Message, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(snapshot); err != nil {
		http.Error(w, "failed to encode messages", http.StatusInternalServerError)
	}
}

func snapshotClients() []*websocket.Conn {
	mu.RLock()
	defer mu.RUnlock()

	connections := make([]*websocket.Conn, 0, len(clients))
	for client := range clients {
		connections = append(connections, client)
	}

	return connections
}
