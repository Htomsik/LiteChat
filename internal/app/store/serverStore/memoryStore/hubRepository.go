package memoryStore

import (
	"Chat/internal/app/model"
	"Chat/internal/app/model/constants"
	"errors"
	"github.com/sirupsen/logrus"
)

// HubRepository memory hub storage
type HubRepository struct {
	hubs       map[string]*model.Hub
	hubDeleted chan string
	logger     *logrus.Logger
}

// NewHubRepository create new hub repository
func NewHubRepository(store *ServerStore) *HubRepository {
	hubRepository := &HubRepository{
		hubs:       make(map[string]*model.Hub),
		hubDeleted: make(chan string),
		logger:     store.logger,
	}

	go hubRepository.Run()

	return hubRepository
}

// Run handler for chans
func (repository *HubRepository) Run() {
	for {
		select {
		case hubId := <-repository.hubDeleted:
			err := repository.Remove(hubId)
			if err != nil {
				repository.logger.Error(err)
				return
			}
		}
	}
}

// Create new hub and place in repository
func (repository *HubRepository) Create(hubId string) (*model.Hub, error) {
	newHub := model.HewHub(hubId, repository.logger, repository.hubDeleted)

	err := repository.Add(newHub)

	if err != nil {
		return nil, err
	}

	return newHub, nil
}

// Add create new chat hub
func (repository *HubRepository) Add(hub *model.Hub) error {
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
	hub, err := repository.Find(id)
	if err != nil {
		return err
	}

	hub.Close()
	delete(repository.hubs, id)

	return nil
}
