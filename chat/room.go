package chat

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Room represents a chat room
type Room struct {
	Name       string
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Mutex      sync.Mutex
}

// NewRoom creates a new room
func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte, 256),
	}
}

// Add a new user
func (r *Room) AddClient(c *websocket.Conn) {

	// Create the client
	client := NewClient(c, r)

	// Register the client in the room
	r.Register <- client

	// Goroutine to send messages to the client
	go client.WritePump()

	// Goroutine to receive messages from the client
	client.ReadPump()
}

// Run starts the room loop
func (r *Room) Run() {

	for {

		log.Println("SYSTEM [" + r.Name + "]: Running loop")

		select {
		case client := <-r.Register:
			r.Mutex.Lock()
			r.Clients[client] = true
			r.Mutex.Unlock()
			msg := []byte("SYSTEM: " + client.Username + " has joined the room [" + r.Name + "]")
			log.Println(string(msg))
			r.Broadcast <- msg
		case client := <-r.Unregister:
			r.Mutex.Lock()
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
			r.Mutex.Unlock()
			msg := []byte("SYSTEM: " + client.Username + " has left the room [" + r.Name + "]")
			log.Println(string(msg))
			r.Broadcast <- msg
		case message := <-r.Broadcast:
			log.Print("SYSTEM ["+r.Name+"]: Broadcasting message: ", string(message))
			r.Mutex.Lock()
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					log.Println("SYSTEM [" + r.Name + "]: Closing client connection, " + client.Username)
					close(client.Send)
					delete(r.Clients, client)
				}
			}
			r.Mutex.Unlock()
		}

	}
}
