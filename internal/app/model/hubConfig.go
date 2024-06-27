package model

import "Chat/internal/app/model/chat"

// HubConfig Base configuration for hub
type HubConfig struct {
	AdminRole   *chat.UserRole   `toml:"adminRole"`
	DefaultRole *chat.UserRole   `toml:"defaultRole"`
	Roles       []*chat.UserRole `toml:"roles"`
}

// NewHubConfig new cfg for hub
func NewHubConfig() *HubConfig {
	return &HubConfig{}
}
