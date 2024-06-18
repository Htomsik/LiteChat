package memoryStore

import (
	"Chat/internal/app/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

// ClientRepository storage of chat users
type ClientRepository struct {
	store   *HubStore
	clients map[uuid.UUID]*model.Client
}

// GetClientsOriginalName returning map of users by original usernames
func (repository *ClientRepository) GetClientsOriginalName() map[string][]*model.Client {
	clients := make(map[string][]*model.Client)

	for _, client := range repository.clients {
		clientsByName := make([]*model.Client, 0)
		originalName := client.User.OriginalName()

		// If clients have array of users by this username
		// just add
		// if not create new
		if clientOrig, ok := clients[originalName]; ok {
			clientsByName = append(clientsByName, clientOrig...)
		} else {
			clientsByName = append(clientsByName, client)
		}

		clients[originalName] = clientsByName
	}

	return clients
}

// Find search client by guid
func (repository *ClientRepository) Find(id uuid.UUID) (*model.Client, error) {
	if client, ok := repository.clients[id]; ok {
		return client, nil
	}

	return nil, model.ErrorRecordNotFound
}

// CountByOriginalName count clients
func (repository *ClientRepository) CountByOriginalName(originalName string) (int, error) {
	clients := repository.GetClientsOriginalName()
	if client, ok := clients[originalName]; ok {
		return len(client), nil
	}

	return 0, model.ErrorRecordNotFound
}

// Add returned new userName if its changed
func (repository *ClientRepository) Add(client *model.Client) (string, error) {

	clientFind, err := repository.Find(client.User.Id)
	if err != nil && !errors.Is(err, model.ErrorRecordNotFound) {
		return "", err
	}

	// if client exists return error
	if clientFind != nil {
		return "", model.ErrorAlreadyExists
	}

	count, err := repository.CountByOriginalName(client.User.OriginalName())
	if err != nil && !errors.Is(err, model.ErrorRecordNotFound) {
		return "", err
	}

	// change client name if he not first
	if count > 0 {
		client.User.Name = fmt.Sprintf("%v[%v]", client.User.Name, count)
	}

	// Add to repo
	repository.clients[client.User.Id] = client

	return client.User.Name, nil
}

// Remove delete user by id
func (repository *ClientRepository) Remove(id uuid.UUID) error {

	if _, ok := repository.clients[id]; !ok {
		return model.ErrorRecordNotFound
	}

	delete(repository.clients, id)

	return nil
}

// All get all clients
func (repository *ClientRepository) All() (map[uuid.UUID]*model.Client, error) {
	if repository.clients == nil {
		return nil, model.ErrorRecordNotFound
	}

	return repository.clients, nil
}

// AllUsers get all chat users
func (repository *ClientRepository) AllUsers() ([]*model.ChatUser, error) {

	if repository.clients == nil {
		return nil, model.ErrorRecordNotFound
	}

	users := make([]*model.ChatUser, 0)
	for _, client := range repository.clients {
		users = append(users, client.User)
	}

	return users, nil
}
