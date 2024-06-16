package model

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// ChatMessage message from chat
type ChatMessage struct {
	Type         ChatMessageType `json:"type"`
	User         string          `json:"user"`
	Message      any             `json:"message"`
	DateTime     time.Time       `json:"dateTime"`
	clearPrivacy bool
}

// ClearPrivacy clear privacy for current user
func (msg *ChatMessage) ClearPrivacy(client *Client) bool {

	if !msg.clearPrivacy {
		return true
	}

	switch msg.Type {

	case UsersList:

		// Clear userData
		var clearUsers = make([]ChatUser, 0)

		if value, ok := msg.Message.([]ChatUser); !ok {
			return ok
		} else {
			for _, user := range value {
				if user.Id != client.User.Id {
					user.Id = uuid.Nil
				}
				clearUsers = append(clearUsers, user)
			}
		}
		msg.Message = clearUsers

	default:
		return false
	}

	return true
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
		message.clearPrivacy = true
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
