package memoryStore

import (
	"Chat/internal/app/model/client"
	"Chat/internal/app/store/hubStore"
	"github.com/google/uuid"
)

type HubStore struct {
	clientRepository *ClientRepository
}

// New create new store
func New() *HubStore {
	return &HubStore{}
}

// Client get client repository
func (store *HubStore) Client() hubStore.ClientRepository {
	if store.clientRepository != nil {
		return store.clientRepository
	}

	store.clientRepository = &ClientRepository{
		store:   store,
		clients: make(map[uuid.UUID]*client.Client),
	}

	return store.clientRepository
}
