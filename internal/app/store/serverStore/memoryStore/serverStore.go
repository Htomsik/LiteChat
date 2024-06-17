package memoryStore

import (
	"Chat/internal/app/model"
	"Chat/internal/app/store/serverStore"
)

// ServerStore memory storage
type ServerStore struct {
	hubRepository *HubRepository
}

// New create new store
func New() *ServerStore {
	return &ServerStore{}
}

// Hub get hub repository
func (store *ServerStore) Hub() serverStore.HubRepository {
	if store.hubRepository != nil {
		return store.hubRepository
	}

	store.hubRepository = &HubRepository{
		store: store,
		hubs:  make(map[string]*model.Hub),
	}

	return store.hubRepository
}
