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
    UserID   string
    Conn     *websocket.Conn
    Lobbies  map[string]bool  // tracks which lobbies this client is in
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

// BroadcastToLobby sends a message only to clients subscribed to lobbyID.
func (h *Hub) BroadcastToLobby(lobbyID string, msg interface{}) {
	var payloadToSend []byte

	// Detect if this is a member-list update
	if m, ok := msg.(map[string]interface{}); ok {
		if m["type"] == "member-list" {
			// Render the member-list template
			var buf bytes.Buffer
			tmpl := template.Must(
				template.New("member-list").
					ParseFiles(
						"./views/lobby/index.html",
					),
			)
			err := tmpl.ExecuteTemplate(&buf, "member-list", m)
			if err != nil {
				log.Println("Hub: template render error:", err)
				return
			}

			// Wrap in a message with type + html
			payload := map[string]string{
				"action": "update",
				"id":     "member-list",
				"html":   buf.String(),
			}
			payloadToSend, err = json.Marshal(payload)
			if err != nil {
				log.Println("Hub: marshal error:", err)
				return
			}
		}
	}

	// If we didn't handle it as member-list, send raw JSON
	if payloadToSend == nil {
		var err error
		payloadToSend, err = json.Marshal(msg)
		if err != nil {
			log.Println("Hub: marshal error:", err)
			return
		}
	}

	// Broadcast to clients in the lobby
	for client := range h.clients {
		if client.Lobbies[lobbyID] {
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
