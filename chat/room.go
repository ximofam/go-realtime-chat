package chat

import (
	"encoding/json"
	"sync"
)

type Room struct {
	id string

	server *Server

	clients map[*Client]struct{}

	join chan *Client

	leave chan *Client

	broadcast chan []byte

	mu sync.RWMutex

	close chan struct{}
}

func NewRoom(id string, server *Server) *Room {
	room := &Room{
		id:        id,
		server:    server,
		clients:   make(map[*Client]struct{}),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan []byte),
		close:     make(chan struct{}),
	}

	go room.run()

	return room
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = struct{}{}

			go func() {
				response := WSResponse{
					Type: TypeConnectUser,
					Data: client.user,
				}

				data, err := json.Marshal(&response)
				if err != nil {
					return
				}

				select {
				case r.broadcast <- data:
				case <-r.close:
				}
			}()
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)

			go func() {
				response := WSResponse{
					Type: TypeDisconnectUser,
					Data: client.user,
				}

				data, err := json.Marshal(&response)
				if err != nil {
					return
				}

				select {
				case r.broadcast <- data:
				case <-r.close:
				}
			}()

			if len(r.clients) == 0 {
				r.server.closeRoom(r.id)
				return
			}
		case msg := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(r.clients, client)
				}
			}
		case <-r.close:
			return
		}
	}
}

func (r *Room) Close() {
	close(r.close)
}

func (r *Room) GetAllUsers() []User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]User, 0, len(r.clients))

	for client := range r.clients {
		res = append(res, *client.user)
	}

	return res
}
