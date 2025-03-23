package chat

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

// Cliente representa una conexión de chat
type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	Room     *Room
	Username string
}

func NewClient(c *websocket.Conn, room *Room) *Client {

	username := c.Query("username")
	if username == "" {
		username = "anonymous"
	}

	return &Client{
		Conn:     c,
		Send:     make(chan []byte, 256),
		Room:     room,
		Username: username,
	}

}

// Escritor para enviar mensajes al cliente
func (c *Client) WritePump() {

	// Cerrar la conexion cuando el cliente se pierda
	defer func() {
		c.Conn.Close()
	}()

	// Ciclo para esperar mensajes para el cliente
	for {

		log.Println("Esperando mensajes para enviar al usuario: " + c.Username)

		// Se quea esperando hasta que se recibe un mensaje
		message, ok := <-c.Send

		log.Println("Mensaje recibido para el usuario: "+c.Username+", ", string(message))

		if !ok {
			// La sala se cerro
			log.Println(c.Username + ": la sala se ha cerrado.")
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		log.Println("Enviando mensaje al cliente: ", c.Username)
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Ocurrio un error al mandar mensaje al cliente: "+c.Username+", ", err.Error())
			return
		}

		log.Println("Se envio el menaje al cliente: ", c.Username)
	}

}

// Lector para recibir mensajes del cliente
func (c *Client) ReadPump() {

	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	for {

		log.Println("Esperando nuevos mensajes del cliente: ", c.Username)

		_, message, err := c.Conn.ReadMessage()

		log.Println("Mensaje recibido del cliente: ", c.Username)

		if err != nil {
			log.Printf(c.Username+": Se desconecto de la escucha, error: %v", err)
			break
		}

		log.Println("Mensaje escrito por el cliente: "+c.Username+", ", string(message))

		formattedMsg := []byte(c.Username + ": " + string(message))

		log.Println(c.Username+": Se va mandar mensaje al Broadcast del canal: ", c.Room.Name)

		c.Room.Broadcast <- formattedMsg

		log.Println(c.Username+": Se mando mensaje al Broadcast del canal: ", c.Room.Name)

	}

}
