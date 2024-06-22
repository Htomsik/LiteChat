package chat

import (
	"github.com/google/uuid"
	"time"
)

// User user info in chat
type User struct {
	Id           uuid.UUID `json:"Id"`
	Name         string    `json:"Name"`
	DateTime     time.Time `json:"DateTime"`
	originalName string
}

// NewChatUser generate new user
func NewChatUser(name string) *User {
	return &User{
		Id:           uuid.New(),
		Name:         name,
		DateTime:     time.Now(),
		originalName: name,
	}
}

// OriginalName returned original name of user
func (user *User) OriginalName() string {
	return user.originalName
}
