package router

import (
	"ChatRoom/chat/service"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/ws", handleConnections)
	r.HandleFunc("/messages", service.GetMessages).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(resolvePublicDir())))
	return r
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		http.Error(w, "websocket upgrade failed", http.StatusBadRequest)
		return
	}
	defer func() {
		service.UnregisterClient(ws)
		ws.Close()
	}()

	service.RegisterClient(ws)

	for {
		var msg service.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		service.BroadcastMessage(msg)
	}
}

func resolvePublicDir() string {
	candidates := []string{
		"./public",
		"./chat/public",
	}

	if _, filename, _, ok := runtime.Caller(0); ok {
		candidates = append(candidates, filepath.Join(filepath.Dir(filename), "..", "public"))
	}

	for _, dir := range candidates {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}

		indexPath := filepath.Join(dir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			return dir
		}
	}

	log.Printf("public directory not found, falling back to %s", candidates[len(candidates)-1])
	return candidates[len(candidates)-1]
}
