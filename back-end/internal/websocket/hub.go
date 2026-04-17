package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client watching a board.
type Client struct {
	BoardID string
	Conn    *websocket.Conn
	Send    chan []byte
}

// Message to broadcast to all clients watching a board.
type Message struct {
	BoardID string
	Data    []byte
}

// Hub maintains board -> clients mappings.
type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.BoardID] == nil {
				h.clients[client.BoardID] = make(map[*Client]bool)
			}
			h.clients[client.BoardID][client] = true
			h.mu.Unlock()
			log.Printf("WS: client connected to board %s", client.BoardID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.BoardID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
				}
			}
			h.mu.Unlock()
			log.Printf("WS: client disconnected from board %s", client.BoardID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[msg.BoardID]
			h.mu.RUnlock()

			for client := range clients {
				select {
				case client.Send <- msg.Data:
				default:
					close(client.Send)
					h.mu.Lock()
					delete(h.clients[client.BoardID], client)
					h.mu.Unlock()
				}
			}
		}
	}
}

func (h *Hub) Broadcast(boardID string, data []byte) {
	h.broadcast <- Message{BoardID: boardID, Data: data}
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.Conn.Close()
	}()
	for {
		if _, _, err := c.Conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for {
		msg, ok := <-c.Send
		if !ok {
			c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
