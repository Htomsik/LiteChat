package Hub

import (
	"github.com/sirupsen/logrus"
)

type Hub struct {
	clients         map[*Client]bool
	msgRetranslator chan []byte // listen message from client
	register        chan *Client
	unregister      chan *Client
	logger          *logrus.Logger
}

// HewHub create new hub
func HewHub(logger *logrus.Logger) *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		msgRetranslator: make(chan []byte),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		logger:          logger,
	}
}

// Run running hub
func (hub *Hub) Run() {
	for {
		select {

		// Client connect
		case client := <-hub.register:
			hub.clients[client] = true

		// Client disconnect
		case client := <-hub.unregister:

			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client) // Delete from hub
				close(client.sendMessage)   // Close
			}

		// Retranslate to other clients
		case message := <-hub.msgRetranslator:
			for client := range hub.clients {
				client.sendMessage <- message
			}
		}
	}
}
