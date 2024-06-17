package memoryStore_test

import (
	"Chat/internal/app/model"
	"Chat/internal/app/store/hubStore/memoryStore"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestClientRepository_Add add exists and not exists clients
func TestClientRepository_Add(t *testing.T) {
	// Arrange
	nameNotExists := "new"
	nameExists := "exists"

	st := memoryStore.New()

	clientNotExists := &model.Client{
		User: model.NewChatUser(nameNotExists),
	}

	clientExists := &model.Client{
		User: model.NewChatUser(nameExists),
	}

	st.Client().Add(clientExists)

	// Act
	_, errNotExists := st.Client().Add(clientNotExists)
	_, errExists := st.Client().Add(clientExists)

	// Assert
	assert.NoError(t, errNotExists)
	assert.Error(t, errExists)
}

// TestClientRepository_Remove remove exists and not exists clients
func TestClientRepository_Remove(t *testing.T) {
	// Arrange
	nameNotExists := "notExists"
	nameExists := "exists"

	st := memoryStore.New()

	clientNotExists := &model.Client{
		User: model.NewChatUser(nameNotExists),
	}

	clientExists := &model.Client{
		User: model.NewChatUser(nameExists),
	}

	st.Client().Add(clientExists)

	// Act
	errExists := st.Client().Remove(clientExists.User.Id)
	errNotExists := st.Client().Remove(clientNotExists.User.Id)

	// Assert
	assert.NoError(t, errExists)
	assert.Error(t, errNotExists)
}

// TestClientRepository_Find Find exists and not exists clients
func TestClientRepository_Find(t *testing.T) {
	// Arrange
	nameNotExists := "notExists"
	nameExists := "exists"

	st := memoryStore.New()

	clientNotExists := &model.Client{
		User: model.NewChatUser(nameNotExists),
	}

	clientExists := &model.Client{
		User: model.NewChatUser(nameExists),
	}

	st.Client().Add(clientExists)

	// Act
	clientExistsFind, errExists := st.Client().Find(clientExists.User.Id)
	clientNotExistsFind, errNotExists := st.Client().Find(clientNotExists.User.Id)

	// Assert
	assert.NoError(t, errExists)
	assert.NotNil(t, clientExistsFind)
	assert.Equal(t, clientExistsFind, clientExists)

	assert.Error(t, errNotExists)
	assert.Nil(t, clientNotExistsFind)
}

// TestClientRepository_CountByOriginalName
func TestClientRepository_CountByOriginalName(t *testing.T) {
	// Arrange
	nameNotExists := "notExists"
	nameExists := "exists"

	st := memoryStore.New()

	clientExists := &model.Client{
		User: model.NewChatUser(nameExists),
	}
	st.Client().Add(clientExists)

	// Act
	existsCount, errExists := st.Client().CountByOriginalName(nameExists)
	notExistsCount, errNotExists := st.Client().CountByOriginalName(nameNotExists)

	// Assert
	assert.NoError(t, errExists)
	assert.NotNil(t, existsCount)
	assert.Equal(t, existsCount, 1)

	assert.Error(t, errNotExists)
	assert.Equal(t, notExistsCount, 0)
}
