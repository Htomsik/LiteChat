package model

import (
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

// GetAllUsers returning all users in chat
func (hub *Hub) GetAllUsers() []*ChatUser {
	chatUsers := make([]*ChatUser, 0)

	for _, client := range hub.clients {
		chatUsers = append(chatUsers, client.User)
	}

	return chatUsers
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

	if message.Type == Message {
		// Todo придумать оптимизацию
		if len(hub.messages) == 50 {
			hub.messages = hub.messages[1:]
		}
		hub.messages = append(hub.messages, message)
	}

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
			msg := hub.NewSystemMessage(UsersList)
			hub.sendMessageAll(msg)

		// Client disconnect
		case client := <-hub.unregister:

			if _, ok := hub.clients[client.User.Name]; ok {
				delete(hub.clients, client.User.Name) // Delete from hub
				close(client.sendMessage)             // Close
			}

			// Send message about disconnected
			msg := hub.NewSystemMessage(UsersList)
			hub.sendMessageAll(msg)

		// Retranslate to other clients
		case message, _ := <-hub.msgRetranslator:
			hub.sendMessageAll(message)
		}
	}
}
