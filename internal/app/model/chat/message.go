package chat

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// Message letter from chat
type Message struct {
	Type         MessageType `json:"type"`
	User         string      `json:"user"`
	Message      any         `json:"message"`
	DateTime     time.Time   `json:"dateTime"`
	clearPrivacy bool
}

// ClearPrivacy clear privacy for current user
func (msg *Message) ClearPrivacy(chatUser *User) bool {

	if !msg.clearPrivacy {
		return true
	}

	switch msg.Type {

	case TypeUsersList:

		// Clear userData
		var clearUsers = make([]User, 0)

		if value, ok := msg.Message.([]*User); !ok {
			return ok
		} else {
			for _, userLink := range value {
				user := *userLink
				if user.Id != chatUser.Id {
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
func NewSystemMessage(msgType MessageType, data any) Message {

	message := Message{
		User:     SystemUser,
		Type:     msgType,
		DateTime: time.Now(),
		Message:  data,
	}

	switch msgType {

	case TypeUsersList:
		message.clearPrivacy = true

	default:
		if message.Message == nil {
			message.Message = ""
		}
	}

	return message
}

func NewMessage(user string, message any) Message {
	return Message{
		Type:     TypeMessage,
		User:     user,
		Message:  message,
		DateTime: time.Now(),
	}
}

// ToJson converting message to json
func (msg *Message) ToJson() string {
	byteMessage, _ := json.Marshal(msg)
	return string(byteMessage)
}

// ToByteArray converting message to json and byte array
func (msg *Message) ToByteArray() []byte {

	byteMessage, _ := json.Marshal(msg)

	return byteMessage
}
