package routers

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var clients = make(map[*websocket.Conn]bool)

func SetupWebSocket(app *fiber.App) {
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()

		clients[c] = true

		// Example: send an update when a new client connects
		message := map[string]interface{}{
			"type": "update",
			"id":   "quiz-container",
			"html": `<h2>New Quiz Question</h2><p>What's 2 + 2?</p>`,
		}

		msgBytes, _ := json.Marshal(message)
		if err := c.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
			log.Println("Error writing message:", err)
		}
	}))
}
