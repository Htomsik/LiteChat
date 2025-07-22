package chat

import (
	"Chat/internal/app/model/dto"
	"encoding/json"
	"errors"
	"fmt"
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

// HavePermission check access to functional
func (user *User) HavePermission(permission RolePermission) bool {
	if permission == PermissionNone {
		return true
	}

	if user.Role.IsAdmin {
		return true
	}

	for _, perm := range user.Role.Permissions {
		if perm == permission {
			return true
		}
	}

	return false
}

func (user *User) MarshalJSON() ([]byte, error) {

	if !user.Role.IsAdmin && (user.Role == nil || user.Role.Permissions == nil) {
		return nil, errors.New(fmt.Sprintf("user %v permission cannot be nil", user.Name))
	}

	permissions := make([]string, len(user.Role.Permissions))
	for i, perm := range user.Role.Permissions {
		permissions[i] = perm.String()
	}

	return json.Marshal(dto.UserDTO{
		Id:   user.Id,
		Name: user.Name,
		Role: dto.UserRoleDTO{
			Name:        user.Role.Name,
			IsAdmin:     user.Role.IsAdmin,
			Permissions: permissions,
		},
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
