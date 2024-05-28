package model

import (
	"testing"
)

// TestUser ...
func TestUser(_ *testing.T) *User {
	return &User{
		Email:    "user@ex.com",
		Password: "TestPassword",
		Active:   true,
	}
}
