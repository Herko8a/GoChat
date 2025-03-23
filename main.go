package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/herko8a/gochat/chat"
)

func main() {

	// Se define un Hub de Salas de Chat
	hub := chat.NewHub()

	app := fiber.New(fiber.Config{
		AppName: "Go WebSocket Chat",
	})

	// Configurar middleware para WebSocket
	app.Use("/ws", func(c *fiber.Ctx) error {

		// Verificar si la conexión es un WebSocket
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired

	})

	// Servir archivos estáticos
	app.Static("/", "./static")

	// Endpoint para WebSocket
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {

		// Se crea u obtiene la sala en el Hub
		room := hub.GetOrCreateRoom(c)

		// Se agrega cliente a la sala y se queda escuchando
		room.AddClient(c)

	}))

	log.Fatal(app.Listen(":8080"))
}
