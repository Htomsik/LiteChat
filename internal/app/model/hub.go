package model

import (
	"Chat/internal/app/model/chat"
	"Chat/internal/app/model/websocket"
	"Chat/internal/app/store/hubStore"
	"Chat/internal/app/store/hubStore/memoryStore"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	hubCfgPath string
)

func init() {
	flag.StringVar(&hubCfgPath, "HubConfig-path", "configs/chatConfig.toml", "path to cfg file")
}

type Hub struct {
	Id         string
	store      hubStore.HubStore
	admin      *websocket.Client
	messages   []chat.Message
	Commands   *websocket.Retranslator
	logger     *logrus.Logger
	config     *HubConfig
	hubDeleted chan string
}

// SetAdmin change admin
func (hub *Hub) SetAdmin(client *websocket.Client) {
	// Delete old admin
	if hub.admin != nil {
		hub.admin.User.Role = hub.config.DefaultRole.Name
	}

	// Set new admin
	client.User.Role = hub.config.AdminRole.Name
	hub.admin = client
}

// userListChanged send all clients new userlist
func (hub *Hub) userListChanged() {

	users, err := hub.store.Client().AllUsers()
	if err != nil {
		hub.logger.Error(err)
	}
	msgAll := chat.NewSystemMessage(chat.TypeUsersList, users)
	hub.sendMessageAll(msgAll)
}

// HewHub create new hub
func HewHub(id string, logger *logrus.Logger, hubDeleted chan string) *Hub {

	cfg := NewHubConfig()
	if _, err := toml.DecodeFile(hubCfgPath, cfg); err != nil {
		logger.Fatal(err)
		return nil
	}

	// if admin change def/adm role it allow not lose roles
	cfg.Roles = append(cfg.Roles, cfg.AdminRole)
	cfg.Roles = append(cfg.Roles, cfg.DefaultRole)

	return &Hub{
		Id:         id,
		store:      memoryStore.New(),
		messages:   make([]chat.Message, 0),
		Commands:   websocket.NewCommands(logger),
		logger:     logger,
		config:     cfg,
		hubDeleted: hubDeleted,
	}
}

// Close delete all connections
func (hub *Hub) Close() {
	clients, err := hub.store.Client().All()

	// Close connections to all users
	if err != nil {
		hub.logger.Error(err)
	}
	if len(clients) > 0 {
		for _, client := range clients {
			client.Disconnect()
		}
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
			hub.logger.Warnf("Can't clear message privacy")
		}
	}
}

// clientConnected operations when websocket connecting first time
func (hub *Hub) clientConnected(client *websocket.Client) {
	// Change role if he not admin
	if client.User.Role != hub.config.AdminRole.Name {
		client.User.Role = hub.config.DefaultRole.Name
	}

	// Send him all messages
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
			hub.userListChanged()

		// Client disconnect
		case client := <-hub.Commands.Unregister:

			// Change admin on first connected if admin is disconnected
			if client.User.Role == hub.config.AdminRole.Name {
				newAdmin, err := hub.store.Client().FirstConnected(client.User.Id)
				if err != nil {
					hub.logger.Error(err)
					continue
				}
				hub.SetAdmin(newAdmin)
			}

			err := hub.store.Client().Remove(client.User.Id)
			if err != nil {
				hub.logger.Error(err)
				continue
			}
			hub.userListChanged()

			clients, err := hub.store.Client().All()
			if err != nil {
				hub.logger.Error(err)
			}

			// Delete hub if zero clients
			if len(clients) == 0 {
				hub.hubDeleted <- hub.Id
			}

		// Retranslate to other clients
		case message, _ := <-hub.Commands.Message:
			hub.sendMessageAll(message)
		}
	}
}
