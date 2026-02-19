package router

import (
	"ChatRoom/service"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ws", handleConnections)
	r.HandleFunc("/messages", service.GetMessages).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	return r
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	service.RegisterClient(ws)

	for {
		var msg service.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			service.UnregisterClient(ws)
			break
		}
		service.BroadcastMessage(msg)
	}
}
