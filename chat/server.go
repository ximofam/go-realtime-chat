package chat

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/ximofam/go-realtime-chat/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	autoIncrementID int
	rooms           map[string]*Room
	mu              sync.Mutex
}

func NewServer() *Server {
	s := &Server{
		autoIncrementID: 1,
		rooms:           make(map[string]*Room),
	}

	return s
}

func (s *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	roomID := uuid.New()

	s.mu.Lock()
	s.rooms[roomID.String()] = NewRoom(roomID.String(), s)
	s.mu.Unlock()

	utils.WriteJSON(w, 201, map[string]any{
		"room_id": roomID,
	})
}

func (s *Server) JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("id")
	s.mu.Lock()
	room, ok := s.rooms[roomID]
	s.mu.Unlock()
	if !ok {
		http.Error(w, fmt.Sprintf("Not exists room with id: %s", roomID), 400)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		username = fmt.Sprintf("Guest%d", s.autoIncrementID)
		s.autoIncrementID++
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to shakehand: %v", err), 500)
		return
	}

	user := &User{
		ID:       uuid.New().String(),
		Username: username,
	}

	client := NewClient(room, conn, user)

	go client.readPump()
	go client.writePump()
}

// GET /rooms/{id}/users
func (s *Server) GetUsersOfRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("id")
	if roomID == "" {
		http.Error(w, "Invalid or missing room id", 400)
		return
	}

	s.mu.Lock()
	room, ok := s.rooms[roomID]
	s.mu.Unlock()

	if !ok {
		http.Error(w, "Invalid room id", 400)
		return
	}

	res := room.GetAllUsers()

	utils.WriteJSON(w, 200, res)
}

// GET /rooms/{id}/exists
func (s *Server) IsExistRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.PathValue("id")
	if roomID == "" {
		http.Error(w, "Invalid or missing room id", 400)
		return
	}

	s.mu.Lock()
	_, ok := s.rooms[roomID]
	s.mu.Unlock()

	exists := false

	if ok {
		exists = true
	}

	utils.WriteJSON(w, 200, map[string]any{"exists": exists})
}

func (s *Server) closeRoom(roomID string) {
	if room, ok := s.rooms[roomID]; ok {
		room.Close()
		delete(s.rooms, roomID)
	}
}

func (s *Server) CurrRooms() int {
	return len(s.rooms)
}
