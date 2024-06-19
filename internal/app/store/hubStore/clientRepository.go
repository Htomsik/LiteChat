package hubStore

import (
	"Chat/internal/app/model/chat"
	"Chat/internal/app/model/client"
	"github.com/google/uuid"
)

// ClientRepository storage of chat users
type ClientRepository interface {
	All() (map[uuid.UUID]*client.Client, error)
	AllUsers() ([]*chat.User, error)
	Find(id uuid.UUID) (*client.Client, error)
	CountByOriginalName(name string) (int, error)
	Add(client *client.Client) (string, error)
	Remove(guid uuid.UUID) error
}
