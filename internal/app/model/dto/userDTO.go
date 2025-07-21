package dto

import (
	"github.com/google/uuid"
	"time"
)

// UserDTO Client version of User
type UserDTO struct {
	Id       uuid.UUID `json:"Id"`
	Name     string    `json:"Name"`
	Role     string    `json:"Role"`
	DateTime time.Time `json:"DateTime"`
}
