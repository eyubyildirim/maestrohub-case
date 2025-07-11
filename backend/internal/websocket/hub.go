package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"tower-defence-engine/internal/game"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
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
				close(client.send)
			}
		}
	}
}

func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request, commandCh chan<- game.ICommand) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}

	client.hub.register <- client

	go client.writePump()
	go client.readPump(commandCh)
}

func (h *Hub) Broadcast(message any) {
	encodedMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshalling broadcast message: %v", err)
		return
	}

	for client := range h.clients {
		select {
		case client.send <- encodedMessage:
		default:
			close(client.send)
			delete(h.clients, client)
		}
	}
}
