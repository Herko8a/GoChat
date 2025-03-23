package chat

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Hub maintains the registry of all rooms
type Hub struct {
	Rooms map[string]*Room
	Mutex sync.Mutex
}

// Create a new hub
func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

// GetOrCreateRoom retrieves an existing room or creates a new one
func (h *Hub) GetOrCreateRoom(c *websocket.Conn) *Room {

	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	// Get room name
	roomName := c.Query("room")
	if roomName == "" {
		roomName = "general"
	}

	// Check if the room already exists
	if room, ok := h.Rooms[roomName]; ok {
		log.Println("Room found: ", roomName)
		return room
	}

	log.Println("Creating new room: ", roomName)
	room := NewRoom(roomName)
	h.Rooms[roomName] = room
	go room.Run()

	return room
}
