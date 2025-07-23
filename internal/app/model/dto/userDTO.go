package dto

import (
	"time"
)

// UserDTO Client version of User
type UserDTO struct {
	Id       string      `json:"Id"`
	Name     string      `json:"Name"`
	Role     UserRoleDTO `json:"Role"`
	DateTime time.Time   `json:"DateTime"`
}

// UserRoleDTO Client version of UserRole
type UserRoleDTO struct {
	Name        string   `json:"Name"`
	IsAdmin     bool     `json:"IsAdmin"`
	Permissions []string `json:"Permissions"`
}
