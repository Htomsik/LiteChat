package memoryStore

import (
	"Chat/internal/app/store/serverStore"
	"github.com/sirupsen/logrus"
)

// ServerStore memory storage
type ServerStore struct {
	hubRepository *HubRepository
	logger        *logrus.Logger
}

// New create new store
func New() *ServerStore {
	return &ServerStore{
		logger: logrus.New(),
	}
}

// Hub get hub repository
func (store *ServerStore) Hub() serverStore.HubRepository {
	if store.hubRepository != nil {
		return store.hubRepository
	}

	store.hubRepository = NewHubRepository(store)

	return store.hubRepository
}
