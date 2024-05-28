package testStore

import (
	"Chat/internal/app/model"
	"Chat/internal/app/store"
)

// Store ...
type Store struct {
	userRepository *UserRepository
}

func New() *Store {
	return &Store{}
}

// User ...
func (store *Store) User() store.UserRepository {
	if store.userRepository != nil {
		return store.userRepository
	}

	store.userRepository = &UserRepository{
		store: store,
		users: make(map[int]*model.User),
	}

	return store.userRepository
}
