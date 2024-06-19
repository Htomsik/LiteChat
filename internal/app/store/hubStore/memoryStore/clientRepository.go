package memoryStore

import (
	"Chat/internal/app/model"
	"Chat/internal/app/model/chat"
	"Chat/internal/app/model/client"
	"errors"
	"fmt"
	"github.com/google/uuid"
)

// ClientRepository storage of chat users
type ClientRepository struct {
	store   *HubStore
	clients map[uuid.UUID]*client.Client
}

// GetClientsOriginalName returning map of users by original usernames
func (repository *ClientRepository) GetClientsOriginalName() map[string][]*client.Client {
	clients := make(map[string][]*client.Client)

	for _, cl := range repository.clients {

		clientsByName := make([]*client.Client, 0)
		originalName := cl.User.OriginalName()

		// If clients have array of users by this username
		// just add
		// if not create new
		if clientOrig, ok := clients[originalName]; ok {
			clientsByName = append(clientsByName, clientOrig...)
		} else {
			clientsByName = append(clientsByName, cl)
		}

		clients[originalName] = clientsByName
	}

	return clients
}

// Find search client by guid
func (repository *ClientRepository) Find(id uuid.UUID) (*client.Client, error) {
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
func (repository *ClientRepository) Add(client *client.Client) (string, error) {

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
func (repository *ClientRepository) All() (map[uuid.UUID]*client.Client, error) {
	if repository.clients == nil {
		return nil, model.ErrorRecordNotFound
	}

	return repository.clients, nil
}

// AllUsers get all chat users
func (repository *ClientRepository) AllUsers() ([]*chat.User, error) {

	if repository.clients == nil {
		return nil, model.ErrorRecordNotFound
	}

	users := make([]*chat.User, 0)
	for _, client := range repository.clients {
		users = append(users, client.User)
	}

	return users, nil
}
