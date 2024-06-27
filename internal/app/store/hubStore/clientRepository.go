package hubStore

import (
	"Chat/internal/app/model/chat"
	"Chat/internal/app/model/websocket"
	"github.com/google/uuid"
)

// ClientRepository storage of chat users
type ClientRepository interface {
	All() (map[uuid.UUID]*websocket.Client, error)
	AllUsers() ([]*chat.User, error)
	Find(id uuid.UUID) (*websocket.Client, error)
	FirstConnected(excludedGuid uuid.UUID) (*websocket.Client, error)
	CountByOriginalName(name string) (int, error)
	Add(client *websocket.Client) (string, error)
	Remove(guid uuid.UUID) error
}
