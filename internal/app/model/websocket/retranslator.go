package websocket

import (
	"Chat/internal/app/model/chat"
	"github.com/sirupsen/logrus"
)

// Retranslator chans of our commands
type Retranslator struct {
	Message    chan chat.Message
	Register   chan *Client
	Unregister chan *Client
	logger     *logrus.Logger
}

// NewCommands create new command
func NewCommands(logger *logrus.Logger) *Retranslator {
	return &Retranslator{
		Message:    make(chan chat.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		logger:     logger,
	}
}
