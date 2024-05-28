package testStore

import (
	"Chat/internal/app/model"
)

// UserRepository ...
type UserRepository struct {
	store *Store
	users map[int]*model.User
}

// Add ...
func (repository *UserRepository) Add(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeAdd(); err != nil {
		return err
	}

	user.ID = len(repository.users) + 1

	user.Active = true

	repository.users[user.ID] = user

	return nil
}

// FindByEmail ...
func (repository *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, elem := range repository.users {
		if elem.Email == email {
			return elem, nil
		}
	}

	return nil, model.RecordNotFound
}

// Find ...
func (repository *UserRepository) Find(id int) (*model.User, error) {
	user, exist := repository.users[id]

	if !exist {
		return nil, model.RecordNotFound
	}

	return user, nil
}

// Deactivate ...
func (repository *UserRepository) Deactivate(id int) error {
	user, exist := repository.users[id]

	if !exist {
		return model.RecordNotFound
	}

	user.Active = false

	return nil
}

// Activate ...
func (repository *UserRepository) Activate(id int) error {
	user, exist := repository.users[id]

	if !exist {
		return model.RecordNotFound
	}

	user.Active = true

	return nil
}
