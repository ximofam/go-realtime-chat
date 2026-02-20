package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ximofam/go-realtime-chat/chat"
)

func main() {
	chatServer := chat.NewServer()

	// ===== serve static files =====
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// ===== serve html =====
	http.Handle("/", http.FileServer(http.Dir("./web/templates")))

	http.HandleFunc("POST /rooms", chatServer.CreateRoom)
	http.HandleFunc("GET /rooms/{id}", chatServer.JoinRoom)
	http.HandleFunc("GET /rooms/{id}/users", chatServer.GetUsersOfRoom)
	http.HandleFunc("GET /rooms/{id}/exists", chatServer.IsExistRoom)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port: %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("Failed to run chat: %v", err)
	}
}
