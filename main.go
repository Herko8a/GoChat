package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/herko8a/gochat/chat"
)

func main() {

	// Define a Chat Room Hub
	hub := chat.NewHub()

	app := fiber.New(fiber.Config{
		AppName: "Go WebSocket Chat",
	})

	// Configure middleware for WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {

		// Check if the connection is a WebSocket
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired

	})

	// Serve static files
	app.Static("/", "./static")

	// WebSocket endpoint
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {

		// Create or get the room in the Hub
		room := hub.GetOrCreateRoom(c)

		// Add client to the room and start listening
		room.AddClient(c)

	}))

	log.Fatal(app.Listen(":8080"))
}
