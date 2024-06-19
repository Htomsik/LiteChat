package client

import (
	"Chat/internal/app/model/chat"
	"github.com/gorilla/websocket"
)

// Client middleware around hub and client
type Client struct {
	User        *chat.User
	connection  *websocket.Conn
	SendMessage chan chat.Message
	commands    *Retranslator
}

// RegisterToHub register clint to hub
func (client *Client) RegisterToHub() {
	client.commands.Register <- client
}

// NewClient create new client
func NewClient(commands *Retranslator, connection *websocket.Conn, user *chat.User) *Client {

	return &Client{
		User:        user,
		connection:  connection,
		SendMessage: make(chan chat.Message),
		commands:    commands,
	}
}

// WriteToHub write client messages and pump it to hub
func (client *Client) WriteToHub() {
	defer func() {
		client.commands.Unregister <- client
		client.connection.Close()
	}()

	for {
		_, message, err := client.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				client.commands.logger.Infof("Error when read connection from hub %v", err)
			}
			break
		}

		stringMessage := string(message[:])
		jsonMessage := chat.NewMessage(client.User.Name, stringMessage)

		client.commands.Message <- jsonMessage
	}

}

// ReadFromHub write client hub messages and pump it to client
func (client *Client) ReadFromHub() {

	// TODO добавить время обработки запроса
	defer func() {
		err := client.connection.Close()

		if err != nil {
			client.commands.logger.Infof("Error when close websocet connection %v", err)
		}
	}()

	for {
		select {
		// listening client hub messages
		case chatMessage, ok := <-client.SendMessage:
			// Check hub closed
			if !ok {
				client.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// get next author and send chatMessage from him
			writer, err := client.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				client.commands.logger.Infof("Error when get next writer %v", err)
				return
			} else {
				writer.Write(chatMessage.ToByteArray())
			}

			if err := writer.Close(); err != nil {
				client.commands.logger.Infof("Error when close writer %v", err)
				return
			}

		}
	}
}
