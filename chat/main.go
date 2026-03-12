package main

import (
	"ChatRoom/chat/router"
	"log"
	"net/http"
)

func main() {
	r := router.NewRouter()

	log.Println("HTTP server started on :8000")
	err := http.ListenAndServe(":8000", r)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
