package main

import (
	"log"
	"net/http"

	"encoding/json"

	"os"
	"os/exec"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}
var messages []Message

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/dump", handleDump)

	go handleMessages()

	log.Println("HTTP server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		messages = append(messages, msg)
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleDump(w http.ResponseWriter, r *http.Request) {
	// Create a temporary file to store the messages
	tmpFile, err := os.Create("chat_dump.json")
	if err != nil {
		log.Fatal(err)
	}
	defer tmpFile.Close()

	// Write messages to the temporary file
	json.NewEncoder(tmpFile).Encode(messages)

	// Execute the Python script
	cmd := exec.Command("python3", "dump_chat.py")
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Chat log dumped successfully"})
}
