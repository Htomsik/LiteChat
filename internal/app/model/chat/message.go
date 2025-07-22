package chat

import (
	"Chat/internal/app/model/dto"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// Message letter from chat
type Message struct {
	Type         MessageType `json:"type"`
	User         *User       `json:"user"`
	Message      any         `json:"message"` // TODO Переделать на интерфейс + Разные типы контента
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
		User:     &User{Name: SystemUser},
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

func NewMessage(user *User, message any) Message {
	return Message{
		Type:     TypeMessage,
		User:     user,
		Message:  message,
		DateTime: time.Now(),
	}
}

func (msg *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(dto.MessageDTO{
		Type:     msg.Type.String(),
		User:     msg.User.Name,
		Message:  msg.Message,
		DateTime: msg.DateTime,
	})
}

// ToJson converting message to json
func (msg *Message) ToJson() string {
	byteMessage, _ := json.Marshal(msg)
	return string(byteMessage)
}

// ToByteArray converting message to json and byte array
func (msg *Message) ToByteArray() ([]byte, error) {
	byteMessage, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return byteMessage, nil
}
