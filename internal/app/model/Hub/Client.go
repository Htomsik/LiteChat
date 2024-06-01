package Hub

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Client middleware around hub and websocket
type Client struct {
	hub         *Hub
	connection  *websocket.Conn
	logger      *logrus.Logger
	sendMessage chan []byte
}

// RegisterToHub register clint to hub
func (client *Client) RegisterToHub() {
	client.hub.register <- client
}

// NewClient create new client
func NewClient(hub *Hub, connection *websocket.Conn) *Client {

	return &Client{
		hub:         hub,
		connection:  connection,
		logger:      hub.logger,
		sendMessage: make(chan []byte, 256),
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
		client.hub.msgRetranslator <- message
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
		case message, ok := <-client.sendMessage:
			// Check hub closed
			if !ok {
				client.connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// get next author and send message from him
			writer, err := client.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				client.logger.Infof("Error when get next writer %v", err)
				return
			} else {
				writer.Write(message)
			}

			if err := writer.Close(); err != nil {
				client.logger.Infof("Error when close writer %v", err)
				return
			}

		}
	}
}
