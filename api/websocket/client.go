package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/rithikjain/LiveQnA/pkg/question"
	"log"
	"net/http"
)

type Client struct {
	hub     *Hub
	conn    *websocket.Conn
	inbound chan *question.Question
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Listen for closing
func (c *Client) readWS() {
	defer func() {
		c.hub.Unregister <- c
		_ = c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// Function for writing to the websocket
func (c *Client) writeWS() {
	for {
		select {
		case que, ok := <-c.inbound:
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(que)
			if err != nil {
				log.Println("Error writing to client")
				_ = c.conn.Close()
				break
			}
		}
	}
}

// Serve the websocket
func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		hub:     hub,
		conn:    conn,
		inbound: make(chan *question.Question),
	}
	client.hub.Register <- client

	go client.writeWS()
	go client.readWS()
}
