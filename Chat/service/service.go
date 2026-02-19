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
	broadcast = make(chan Message)
	messages  []Message
	mu        sync.Mutex
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
	defer mu.Unlock()
	messages = append(messages, msg)
	broadcast <- msg
}

func handleMessages() {
	for {
		msg := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
