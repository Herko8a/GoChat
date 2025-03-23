package chat

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

// Client represents a chat connection
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

// Writer to send messages to the client
func (c *Client) WritePump() {

	// Close the connection when the client is lost
	defer func() {
		c.Conn.Close()
	}()

	// Loop to wait for messages for the client
	for {

		log.Println("Waiting for messages to send to user: " + c.Username)

		// Wait until a message is received
		message, ok := <-c.Send

		log.Println("Message received for user: "+c.Username+", ", string(message))

		if !ok {
			// The room was closed
			log.Println(c.Username + ": the room has been closed.")
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		log.Println("Sending message to client: ", c.Username)
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("An error occurred while sending a message to client: "+c.Username+", ", err.Error())
			return
		}

		log.Println("Message sent to client: ", c.Username)
	}

}

// Reader to receive messages from the client
func (c *Client) ReadPump() {

	defer func() {
		c.Room.Unregister <- c
		c.Conn.Close()
	}()

	for {

		log.Println("Waiting for new messages from client: ", c.Username)

		_, message, err := c.Conn.ReadMessage()

		log.Println("Message received from client: ", c.Username)

		if err != nil {
			log.Printf(c.Username+": Disconnected from listening, error: %v", err)
			break
		}

		log.Println("Message written by client: "+c.Username+", ", string(message))

		formattedMsg := []byte(c.Username + ": " + string(message))

		log.Println(c.Username+": Sending message to Broadcast of channel: ", c.Room.Name)

		c.Room.Broadcast <- formattedMsg

		log.Println(c.Username+": Message sent to Broadcast of channel: ", c.Room.Name)

	}

}
