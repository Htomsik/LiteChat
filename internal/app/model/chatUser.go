package model

// ChatUser user info in chat
type ChatUser struct {
	Name string
}

// NewChatUser generate new user
func NewChatUser(name string) *ChatUser {
	return &ChatUser{
		Name: name,
	}
}
