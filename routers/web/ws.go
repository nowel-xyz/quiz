package web

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
	"github.com/nowel-xyz/quiz/database/models"
	"github.com/nowel-xyz/quiz/routers/middleware"
	"github.com/nowel-xyz/quiz/service/lobby"
	"github.com/nowel-xyz/quiz/utils/web"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

// SetupWebSocket configures the /ws endpoint with auth, heartbeat, and hub integration.
func SetupWebSocket(app *fiber.App) {
	go web.HubInstance.Run()

	app.Use("/ws", middleware.RequireAuth())

	app.Get("/ws", fiberws.New(func(wsConn *fiberws.Conn) {
		raw := wsConn.Locals("user")
		user, ok := raw.(models.User)
		if !ok {
			log.Println("ws: missing user in Locals")
			wsConn.Close()
			return
		}

		// Fetch all lobby IDs for the user from Redis
		lobbyIDs, err := service_lobby.GetLobbyIDsForUser(context.Background(), user)
		if err != nil {
			log.Println("ws: failed to load user lobbies:", err)
			wsConn.Close()
			return
		}

		initialPath := wsConn.Query("path")
		client := &web.Client{
			UserID:      user.ID.Hex(),
			Conn:        wsConn,
			Lobbies:    make(map[string]bool),
			CurrentPath: initialPath,
		}
		for _, lid := range lobbyIDs {
			client.Lobbies[lid] = true
		}

		web.HubInstance.Register(client)
		defer web.HubInstance.Unregister(client)

		// WebSocket-level heartbeat (ping/pong frames)
		wsConn.SetReadDeadline(time.Now().Add(pongWait))
		wsConn.SetPongHandler(func(string) error {
			return wsConn.SetReadDeadline(time.Now().Add(pongWait))
		})
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		go func() {
			for range ticker.C {
				if err := wsConn.WriteMessage(fiberws.PingMessage, nil); err != nil {
					return
				}
			}
		}()

		// Send welcome message
		welcome := map[string]interface{}{ "type": "welcome", "user": user.Username }
		if b, err := json.Marshal(welcome); err == nil {
			wsConn.WriteMessage(fiberws.TextMessage, b)
		}

		// Main read loop
		for {
			_, msg, err := wsConn.ReadMessage()
			if err != nil {
				break
			}

			var req struct {
				Action  string `json:"action"`
				LobbyID string `json:"lobbyID"`
				Path    string `json:"path,omitempty"`
			}
			if err := json.Unmarshal(msg, &req); err != nil {
				continue
			}

			switch req.Action {

			case "ping":
				pong := map[string]string{"type": "pong"}
				if b, err := json.Marshal(pong); err == nil {
					wsConn.WriteMessage(fiberws.TextMessage, b)
				}
			case "setLocation":
				if req.Path != "" {
					client.CurrentPath = req.Path
					log.Printf("Client %s set path to %s", client.UserID, req.Path)
				}
			}
		}
	}))
}
