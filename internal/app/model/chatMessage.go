package model

import (
	"encoding/json"
	"time"
)

// ChatMessage message from chat
type ChatMessage struct {
	Type     ChatMessageType `json:"type"`
	User     string          `json:"user"`
	Message  any             `json:"message"`
	DateTime time.Time       `json:"dateTime"`
}

// NewSystemMessage message from chat
func (hub *Hub) NewSystemMessage(msgType ChatMessageType) ChatMessage {

	message := ChatMessage{
		User:     SystemUser,
		Type:     msgType,
		DateTime: time.Now(),
	}

	switch msgType {

	case UsersList:
		message.Message = hub.GetAllUsers()
	default:
		message.Message = ""
	}

	return message
}

func NewMessage(user string, message any) ChatMessage {
	return ChatMessage{
		Type:     Message,
		User:     user,
		Message:  message,
		DateTime: time.Now(),
	}
}

// ToJson converting message to json
func (msg *ChatMessage) ToJson() string {
	byteMessage, _ := json.Marshal(msg)
	return string(byteMessage)
}

// ToByteArray converting message to json and byte array
func (msg *ChatMessage) ToByteArray() []byte {
	byteMessage, _ := json.Marshal(msg)
	return byteMessage
}
