package chat

import (
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

// Room representa una sala de chat
type Room struct {
	Name       string
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Mutex      sync.Mutex
}

// NewRoom crea una nueva sala
func NewRoom(name string) *Room {
	return &Room{
		Name:       name,
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte, 256),
	}
}

// Add new user
func (r *Room) AddClient(c *websocket.Conn) {

	// Se crea el cliente
	client := NewClient(c, r)

	// Registramos al cliente en la sala
	r.Register <- client

	// Goroutine para enviar mensajes al cliente
	go client.WritePump()

	// Gorutine para recibir mensajes del cliente
	client.ReadPump()
}

// Run inicia el bucle de la sala
func (r *Room) Run() {

	for {

		log.Println("SYSTEM [" + r.Name + "]: Toy ciclado")

		select {
		case client := <-r.Register:
			r.Mutex.Lock()
			r.Clients[client] = true
			r.Mutex.Unlock()
			msg := []byte("SYSTEM: " + client.Username + " se ha unido a la sala [" + r.Name + "]")
			log.Println(string(msg))
			r.Broadcast <- msg
		case client := <-r.Unregister:
			r.Mutex.Lock()
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.Send)
			}
			r.Mutex.Unlock()
			msg := []byte("SYSTEM: " + client.Username + " ha salido de la sala [" + r.Name + "]")
			log.Println(string(msg))
			r.Broadcast <- msg
		case message := <-r.Broadcast:
			log.Print("SYSTEM ["+r.Name+"]: Se va repartir mensaje: ", string(message))
			r.Mutex.Lock()
			for client := range r.Clients {
				select {
				case client.Send <- message:
				default:
					log.Println("SYSTEM [" + r.Name + "]: Se cierra conexión del cliente, " + client.Username)
					close(client.Send)
					delete(r.Clients, client)
				}
			}
			r.Mutex.Unlock()
		}

	}
}
