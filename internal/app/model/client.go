package model

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Client middleware around hub and websocket
type Client struct {
	User        *ChatUser
	hub         *Hub
	connection  *websocket.Conn
	logger      *logrus.Logger
	sendMessage chan ChatMessage
}

// RegisterToHub register clint to hub
func (client *Client) RegisterToHub() {
	client.hub.register <- client
}

// NewClient create new client
func NewClient(hub *Hub, connection *websocket.Conn, user *ChatUser) *Client {

	return &Client{
		User:        user,
		hub:         hub,
		connection:  connection,
		logger:      hub.logger,
		sendMessage: make(chan ChatMessage),
	}
}

// WriteToHub write websocket messages and pump it to hub
func (client *Client) WriteToHub() {
	defer func() {
		client.hub.unregister <- client
		client.connection.Close()
	}()

	for {
		_, message, err := client.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				client.logger.Infof("Error when read connection from hub %v", err)
			}
			break
		}

		stringMessage := string(message[:])

		jsonMessage := ChatMessage{
			User:    client.User.Name,
			Message: stringMessage,
		}

		client.hub.msgRetranslator <- jsonMessage
	}

}

// ReadFromHub write client hub messages and pump it to websocket
func (client *Client) ReadFromHub() {

	// TODO добавить время обработки запроса
	defer func() {
		err := client.connection.Close()

		if err != nil {
			client.logger.Infof("Error when close websocet connection %v", err)
		}
	}()

	for {
		select {
		// listening client hub messages
		case chatMessage, ok := <-client.sendMessage:
			// Check hub closed
			if !ok {
				client.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// get next author and send chatMessage from him
			writer, err := client.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				client.logger.Infof("Error when get next writer %v", err)
				return
			} else {
				writer.Write(chatMessage.ToByteArray())
			}

			if err := writer.Close(); err != nil {
				client.logger.Infof("Error when close writer %v", err)
				return
			}

		}
	}
}
