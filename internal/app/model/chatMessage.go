package model

import "encoding/json"

const (
	SystemUser = "System"
)

// ChatMessage message from chat
type ChatMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

// NewSystemMessage message from chat
func NewSystemMessage(text string) ChatMessage {
	return ChatMessage{
		User:    SystemUser,
		Message: text,
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
