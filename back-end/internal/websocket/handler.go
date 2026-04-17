package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Restrict in production
	},
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	boardID := r.URL.Query().Get("board_id")
	if boardID == "" {
		http.Error(w, "board_id is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		BoardID: boardID,
		Conn:    conn,
		Send:    make(chan []byte, 256),
	}

	hub.register <- client

	go client.WritePump()
	go client.ReadPump(hub)
}
