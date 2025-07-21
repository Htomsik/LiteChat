package chat

import (
	"Chat/internal/app/model/dto"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// User user info in chat
type User struct {
	Id           uuid.UUID `json:"Id"`
	Name         string    `json:"Name"`
	Role         *UserRole `json:"Role"`
	DateTime     time.Time `json:"DateTime"` // Connection time
	originalName string
}

func (user *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(dto.UserDTO{
		Id:       user.Id,
		Name:     user.Name,
		Role:     user.Role.Name,
		DateTime: user.DateTime,
	})
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
