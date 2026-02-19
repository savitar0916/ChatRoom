package main

import (
    "log"
    "net/http"
    "ChatRoom/router"
)

func main() {
    r := router.NewRouter()

    log.Println("HTTP server started on :8000")
    err := http.ListenAndServe(":8000", r)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
