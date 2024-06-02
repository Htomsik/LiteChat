package model

import "encoding/json"

// ChatMessage message from chat
type ChatMessage struct {
	User    string `json:"user"`
	Message string `json:"message"`
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
