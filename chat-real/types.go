package main

import (
	"encoding/json"
	// "fmt"
	// "sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID string

	// Hub object
	hub *Hub
	// mu  sync.Mutex

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	clients    map[string]Client
	broadcast  chan []byte
	register   chan Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan Client),
		unregister: make(chan *Client),
		clients:    make(map[string]Client),
	}
}

type Message struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

func (manager *Hub) send(message []byte, ignore Client) {
	for _, conn := range manager.clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}

type M struct {
	Content  string
	ID       string
	connects []Client
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:

			var c []Client

			for _, y := range h.clients {
				c = append(c, y)
			}

			h.clients[client.ID] = client
			value, _ := json.Marshal(&M{Content: "New Socket connected", ID: client.ID, connects: c})
			h.send(value, client)

		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
			}

		case message := <-h.broadcast:
			var f Message
			_ = json.Unmarshal(message, &f)
			if client, ok := h.clients[f.ID]; ok {
				select {
				case client.send <- []byte(f.Data):
				default:
					close(client.send)
					delete(h.clients, client.ID)
				}
			}
		}

	}
}
