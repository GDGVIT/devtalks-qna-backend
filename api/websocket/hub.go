package websocket

import (
	"github.com/rithikjain/LiveQnA/pkg/question"
	"log"
)

type Hub struct {
	// Stores all the connected clients
	Clients    map[*Client]bool
	Broadcast  chan *question.Question
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan *question.Question),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println("User Connected..")
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.inbound)
				log.Println("User Disconnected..")
			}
		case que := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.inbound <- que:
				default:
					close(client.inbound)
					delete(h.Clients, client)
				}
			}
		}
	}
}
