package chat

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Hub mantiene el registro de todas las salas
type Hub struct {
	Rooms map[string]*Room
	Mutex sync.Mutex
}

// Crea un nuevo hub
func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

// GetOrCreateRoom obtiene una sala existente o crea una nueva
func (h *Hub) GetOrCreateRoom(c *websocket.Conn) *Room {

	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	// Obtener nombre de sala
	roomName := c.Query("room")
	if roomName == "" {
		roomName = "general"
	}

	// Se busca si ya existe la sala
	if room, ok := h.Rooms[roomName]; ok {
		log.Println("Sala encontrada: ", roomName)
		return room
	}

	log.Println("Se crea nueva sala: ", roomName)
	room := NewRoom(roomName)
	h.Rooms[roomName] = room
	go room.Run()

	return room
}
