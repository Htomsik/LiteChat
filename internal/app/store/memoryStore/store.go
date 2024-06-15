package memoryStore

import (
	"Chat/internal/app/model"
	"Chat/internal/app/store"
)

// Store memory storage
type Store struct {
	hubRepository *HubRepository
}

// New create new store
func New() *Store {
	return &Store{}
}

// Hub get hub repository
func (store *Store) Hub() store.HubRepository {
	if store.hubRepository != nil {
		return store.hubRepository
	}

	store.hubRepository = &HubRepository{
		store: store,
		hubs:  make(map[string]*model.Hub),
	}

	return store.hubRepository
}
