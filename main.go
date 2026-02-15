package main

import (
	"net/http"

	"github.com/ximofam/go-realtime-chat/chat"
)

func main() {
	chatServer := chat.NewServer()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("POST /rooms", chatServer.CreateRoom)
	http.HandleFunc("GET /rooms/{id}", chatServer.JoinRoom)
	http.HandleFunc("GET /rooms/{id}/users", chatServer.GetUsersOfRoom)
	http.HandleFunc("GET /rooms/{id}/exists", chatServer.IsExistRoom)

	http.ListenAndServe(":8080", nil)

}
