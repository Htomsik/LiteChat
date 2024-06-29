package serverStore

import (
	"Chat/internal/app/model"
)

// HubRepository chats storage
type HubRepository interface {
	Create(id string, cfg *model.HubConfig) (*model.Hub, error)
	Add(hub *model.Hub) error
	Find(id string) (*model.Hub, error)
	Remove(id string) error
}
