package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func HandleDump(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	// Create a temporary file to store the messages
	tmpFile, err := os.Create("chat_dump.json")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()

	// Write messages to the temporary file
	json.NewEncoder(tmpFile).Encode(messages)

	// Execute the Python script
	cmd := exec.Command("python3", "dump_chat.py")
	err = cmd.Run()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Serve the file as a download
	w.Header().Set("Content-Disposition", "attachment; filename=chat_dump.json")
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "chat_dump.json")
}
