package model

import (
	"Chat/internal/app/model/chat"
	"Chat/internal/app/model/websocket"
	"Chat/internal/app/store/hubStore"
	"Chat/internal/app/store/hubStore/memoryStore"
	"github.com/sirupsen/logrus"
)

type Hub struct {
	Id       string
	store    hubStore.HubStore
	messages []chat.Message
	Commands *websocket.Retranslator
	logger   *logrus.Logger
}

// HewHub create new hub
func HewHub(id string, logger *logrus.Logger) *Hub {
	return &Hub{
		Id:       id,
		store:    memoryStore.New(),
		messages: make([]chat.Message, 0),
		Commands: websocket.NewCommands(logger),
		logger:   logger,
	}
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

	clients, err := hub.store.Client().All()

	if err != nil {
		hub.logger.Error(err)
		return
	}

	for _, cl := range clients {

		localMessage := message
		if localMessage.ClearPrivacy(cl.User) {
			cl.SendMessage <- localMessage
		} else {
			hub.logger.Warnf("Can't clear websocket priuvacy")
		}
	}
}

// clientConnected operations when websocket connecting first time
func (hub *Hub) clientConnected(client *websocket.Client) {
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

			originName := client.User.OriginalName()

			newName, err := hub.store.Client().Add(client)

			if err != nil {
				hub.logger.Error(err)
				continue
			}

			if newName != originName {
				msg := chat.NewSystemMessage(chat.TypeUserNameChanged, client.User.Name)
				client.SendMessage <- msg
			}

			hub.clientConnected(client)

			// Send message about connected
			users, err := hub.store.Client().AllUsers()
			if err != nil {
				hub.logger.Error(err)
				continue
			}
			msgAll := chat.NewSystemMessage(chat.TypeUsersList, users)
			hub.sendMessageAll(msgAll)

		// Client disconnect
		case client := <-hub.Commands.Unregister:

			err := hub.store.Client().Remove(client.User.Id)
			if err != nil {
				hub.logger.Error(err)
				continue
			}

			// Send message about disconnected
			users, err := hub.store.Client().AllUsers()
			if err != nil {
				hub.logger.Error(err)
				continue
			}
			msg := chat.NewSystemMessage(chat.TypeUsersList, users)
			hub.sendMessageAll(msg)

		// Retranslate to other clients
		case message, _ := <-hub.Commands.Message:
			hub.sendMessageAll(message)
		}
	}
}
