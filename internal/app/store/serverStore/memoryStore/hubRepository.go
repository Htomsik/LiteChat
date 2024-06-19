package memoryStore

import (
	"Chat/internal/app/model"
	"Chat/internal/app/model/constants"
	"errors"
)

// HubRepository memory hub storage
type HubRepository struct {
	store *ServerStore
	hubs  map[string]*model.Hub
}

// Add create new chat hub
func (repository *HubRepository) Add(hub *model.Hub) error {

	// Todo add validation
	hubFind, err := repository.Find(hub.Id)

	// if error is not record not found
	if err != nil && !errors.Is(err, constants.ErrorRecordNotFound) {
		return err
	}

	// if hub exists don't create new
	if hubFind != nil {
		return constants.ErrorAlreadyExists
	}

	repository.hubs[hub.Id] = hub

	return nil
}

// Find search exists chat hub
func (repository *HubRepository) Find(id string) (*model.Hub, error) {

	if hub, ok := repository.hubs[id]; ok {
		return hub, nil
	}

	return nil, constants.ErrorRecordNotFound
}

// Remove delete exists chat hub
func (repository *HubRepository) Remove(id string) error {

	if _, err := repository.Find(id); err != nil {
		return err
	}

	delete(repository.hubs, id)

	return nil
}
