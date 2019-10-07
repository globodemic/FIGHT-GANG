package websocket

import rp "../models/responses"

type hub struct {
	clients    map[*client]bool
	broadcast  chan []*rp.ResponsePlayer
	register   chan *client
	unregister chan *client

	content []*rp.ResponsePlayer
}

var h = hub{
	broadcast:  make(chan []*rp.ResponsePlayer, 0),
	register:   make(chan *client),
	unregister: make(chan *client),
	clients:    make(map[*client]bool),
	content:    make([]*rp.ResponsePlayer, 0),
}

//NewHub func
func NewHub() hub {
	return h
}

func (h *hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			//When Client registered send to active players list
			c.send <- h.content
			break

		case c := <-h.unregister:
			_, ok := h.clients[c]
			if ok {
				delete(h.clients, c)
				close(c.send)
			}
			break

		case m := <-h.broadcast:
			h.content = m
			h.BroadcastMessage()
			break
		}
	}
}

func (h *hub) BroadcastMessage() {
	for c := range h.clients {
		select {
		case c.send <- h.content:
			break

		// We can't reach the client
		default:
			close(c.send)
			delete(h.clients, c)
		}
	}
}
