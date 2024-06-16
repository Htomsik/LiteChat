package model

import "github.com/google/uuid"

// ChatUser user info in chat
type ChatUser struct {
	Id           uuid.UUID `json:"Id"`
	Name         string    `json:"Name"`
	originalName string
}

// NewChatUser generate new user
func NewChatUser(name string) *ChatUser {
	return &ChatUser{
		Id:           uuid.New(),
		Name:         name,
		originalName: name,
	}
}
