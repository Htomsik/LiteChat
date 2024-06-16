package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	Id       string
	clients  map[string]*Client // All users connected to chat
	messages []ChatMessage

	msgRetranslator chan ChatMessage // listen message from client
	register        chan *Client
	unregister      chan *Client

	logger *logrus.Logger
}

// HewHub create new hub
func HewHub(id string, logger *logrus.Logger) *Hub {
	return &Hub{
		Id:       id,
		clients:  make(map[string]*Client),
		messages: make([]ChatMessage, 0),

		msgRetranslator: make(chan ChatMessage),
		register:        make(chan *Client),
		unregister:      make(chan *Client),

		logger: logger,
	}
}

// FindClient find client by userName
func (hub *Hub) FindClient(userName string) *Client {
	client, ok := hub.clients[userName]
	if !ok {
		return nil
	}
	return client
}

// sendMessageAll send message to all users in hub
func (hub *Hub) sendMessageAll(message ChatMessage) {

	// Todo придумать оптимизацию
	if len(hub.messages) == 50 {
		hub.messages = hub.messages[1:]
	}
	hub.messages = append(hub.messages, message)

	for _, client := range hub.clients {
		client.sendMessage <- message
	}
}

// clientConnected operations when client connecting first time
func (hub *Hub) clientConnected(client *Client) {
	for _, message := range hub.messages {
		client.sendMessage <- message
	}
}

// Run running hub
func (hub *Hub) Run() {
	for {
		select {

		// Client connect
		case client := <-hub.register:
			hub.clients[client.User.Name] = client

			hub.clientConnected(client)

			// Send message about connected
			msg := NewSystemMessage(fmt.Sprintf("%v connected to chat", client.User.Name))
			hub.sendMessageAll(msg)

		// Client disconnect
		case client := <-hub.unregister:

			// Send message about connected
			msg := NewSystemMessage(fmt.Sprintf("%v disconnected from chat", client.User.Name))
			hub.sendMessageAll(msg)

			if _, ok := hub.clients[client.User.Name]; ok {
				delete(hub.clients, client.User.Name) // Delete from hub
				close(client.sendMessage)             // Close
			}

		// Retranslate to other clients
		case message, _ := <-hub.msgRetranslator:
			hub.sendMessageAll(message)
		}
	}
}
