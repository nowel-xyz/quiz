// web/hub.go
package web

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
    UserID      string
    Conn        *websocket.Conn
    Lobbies     map[string]bool  // tracks which lobbies this client is in
    CurrentPath string            // ← track what page the user is on

}

type Hub struct {
    clients    map[*Client]bool
    register   chan *Client
    unregister chan *Client
    broadcast  chan []byte
}

var HubInstance = NewHub()

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan []byte),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true

        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                client.Conn.Close()
            }

        case message := <-h.broadcast:
            // broadcast to _all_ clients
            for client := range h.clients {
                if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
                    log.Println("write error:", err)
                    h.unregister <- client
                }
            }
        }
    }
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister queues a client for unregistration.
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// SendUpdate sends to every connected client
func (h *Hub) SendUpdate(msg interface{}) {
    b, _ := json.Marshal(msg)
    h.broadcast <- b
}

func (h *Hub) BroadcastToLobby(lobbyID string, msg interface{}) {
	var payloadToSend []byte

	if m, ok := msg.(map[string]interface{}); ok {
		if m["type"] == "member-list" {
			// Safely extract path
			path, ok := m["path"].(string)
			log.Printf("Hub: member-list update for path %s", path)
			if !ok {
				log.Println("Hub: member-list update missing valid 'path'")
				return
			}

			// Render the member-list template
			var buf bytes.Buffer
			tmpl := template.Must(
				template.New("member-list").
					ParseFiles(
						"./views/lobby/index.html", // Assuming this contains a {{define "member-list"}} section
					),
			)
			if err := tmpl.ExecuteTemplate(&buf, "member-list", m); err != nil {
				log.Println("Hub: template render error:", err)
				return
			}

			// Wrap into WS payload
			payload := map[string]string{
				"action": "update",
				"id":     "member-list",
				"html":   buf.String(),
				"path":   path,
			}
			var err error
			payloadToSend, err = json.Marshal(payload)
			if err != nil {
				log.Println("Hub: marshal error:", err)
				return
			}
		}
	}

	// Fallback: send raw if no special handling
	if payloadToSend == nil {
		var err error
		payloadToSend, err = json.Marshal(msg)
		if err != nil {
			log.Println("Hub: marshal error:", err)
			return
		}
	}

	// Broadcast only to users on the correct page
	for client := range h.clients {
		log.Printf("Hub: broadcasting to client %s for lobby %s user path %s", client.UserID, lobbyID, client.CurrentPath)
		if client.Lobbies[lobbyID] && client.CurrentPath == "/lobby/"+lobbyID {
			if err := client.Conn.WriteMessage(websocket.TextMessage, payloadToSend); err != nil {
				log.Printf("Hub: write error to client %s: %v", client.UserID, err)
				h.unregister <- client
			}
		}
	}
}

// AddUserToLobby subscribes all active connections of userID to lobbyID and logs activity.
func (h *Hub) AddUserToLobby(userID, lobbyID string) {
	log.Printf("Hub: AddUserToLobby user=%s lobby=%s", userID, lobbyID)
	for client := range h.clients {
		if client.UserID == userID {
			log.Printf("Hub: adding client %s to lobby %s", client.UserID, lobbyID)
			client.Lobbies[lobbyID] = true
		}
	}
}
