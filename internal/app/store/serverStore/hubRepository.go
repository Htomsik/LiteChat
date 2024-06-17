package serverStore

import "Chat/internal/app/model"

// HubRepository chats storage
type HubRepository interface {
	Add(hub *model.Hub) error
	Find(id string) (*model.Hub, error)
	Remove(id string) error
}
