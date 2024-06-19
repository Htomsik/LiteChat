package model

import (
	"Chat/internal/app/model/chat"
	"Chat/internal/app/model/client"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	Id       string
	clients  map[uuid.UUID]*client.Client // All users connected to chat
	messages []chat.Message
	Commands *client.Retranslator
	logger   *logrus.Logger
}

// HewHub create new hub
func HewHub(id string, logger *logrus.Logger) *Hub {
	return &Hub{
		Id:       id,
		clients:  make(map[uuid.UUID]*client.Client),
		messages: make([]chat.Message, 0),
		Commands: client.NewCommands(logger),
		logger:   logger,
	}
}

// GetAllUsers returning all users in chat
func (hub *Hub) GetAllUsers() []chat.User {
	chatUsers := make([]chat.User, 0)

	for _, client := range hub.clients {
		chatUsers = append(chatUsers, *client.User)
	}

	return chatUsers
}

// TODO сделать репозиторий управления юзерами

// CountUsersByOriginalName find first client with original name
func (hub *Hub) CountUsersByOriginalName(originalName string) int {
	count := 0

	for _, client := range hub.clients {
		if client.User.OriginalName() == originalName {
			count++
		}
	}

	return count
}

// sendMessageAll send message to all users in hub
func (hub *Hub) sendMessageAll(message chat.Message) {

	if message.Type == chat.TypeMessage {
		// Todo придумать оптимизацию
		if len(hub.messages) == 50 {
			hub.messages = hub.messages[1:]
		}
		hub.messages = append(hub.messages, message)
	}

	for _, client := range hub.clients {
		localMessage := message
		if localMessage.ClearPrivacy(client.User) {
			client.SendMessage <- localMessage
		} else {
			hub.logger.Warnf("Can't clear client priuvacy")
		}
	}
}

// clientConnected operations when client connecting first time
func (hub *Hub) clientConnected(client *client.Client) {
	for _, message := range hub.messages {
		client.SendMessage <- message
	}
}

// Run running hub
func (hub *Hub) Run() {
	for {
		select {

		// Client connect
		case client := <-hub.Commands.Register:
			hub.clients[client.User.Id] = client

			// Check is user with same originalName is connected
			// if yes change name +1
			if count := hub.CountUsersByOriginalName(client.User.Name); count > 1 {
				client.User.Name = fmt.Sprintf("%v[%v]", client.User.Name, count)

				// notify about changing username
				msg := chat.NewSystemMessage(chat.TypeUserNameChanged, client.User.Name)
				client.SendMessage <- msg
			}

			hub.clientConnected(client)

			// Send message about connected
			msgAll := chat.NewSystemMessage(chat.TypeUsersList, hub.GetAllUsers())
			hub.sendMessageAll(msgAll)

		// Client disconnect
		case client := <-hub.Commands.Unregister:

			if _, ok := hub.clients[client.User.Id]; ok {
				delete(hub.clients, client.User.Id) // Delete from hub
				close(client.SendMessage)           // Close
			}

			// Send message about disconnected
			msg := chat.NewSystemMessage(chat.TypeUsersList, hub.GetAllUsers())
			hub.sendMessageAll(msg)

		// Retranslate to other clients
		case message, _ := <-hub.Commands.Message:
			hub.sendMessageAll(message)
		}
	}
}
