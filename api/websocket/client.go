package websocket

import (
	"log"
	"net/http"
	"time"

	cred "../models"
	rp "../models/responses"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type client struct {
	ws   *websocket.Conn
	send chan []*rp.ResponsePlayer
	ID   string
}

var creds = cred.CredentialsSlice{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	params := mux.Vars(r)
	if len(params) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &client{
		send: make(chan []*rp.ResponsePlayer, 0),
		ws:   ws,
		ID:   params["id"],
	}

	h.register <- c

	go c.writePump()
	c.readPump()
}

func (c *client) readPump() {
	defer func() {
		h.unregister <- c

		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// for {
	// 	_, message, err := c.ws.ReadMessage()
	// 	if err != nil {
	// 		break
	// 	}
	// 	//Broadcast UserList
	// 	resultPlayer := db.GetActivePlayer(c.ID)

	// 	// h.broadcast <- resultPlayer
	// }
}

func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	// for {
	// 	select {
	// 	case message, ok := <-c.send:
	// 		if !ok {
	// 			c.write(websocket.CloseMessage, []byte{})
	// 			return
	// 		}
	// 		// if err := c.write(websocket.TextMessage, message); err != nil {
	// 		// 	return
	// 		// }
	// 	case <-ticker.C:
	// 		if err := c.write(websocket.PingMessage, []byte{}); err != nil {
	// 			return
	// 		}
	// 	}
	// }
}

func (c *client) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}
