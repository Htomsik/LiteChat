package memoryStore_test

import (
	"Chat/internal/app/model"
	"Chat/internal/app/store/serverStore/memoryStore"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestHubRepository_Add add exists and not exists hub
func TestHubRepository_Add(t *testing.T) {
	// Arrange
	idUniq := "new"
	idExists := "exists"

	st := memoryStore.New()

	hubExists := &model.Hub{Id: idExists}
	hubUniq := &model.Hub{Id: idUniq}

	st.Hub().Add(hubExists)

	// Act
	errUniq := st.Hub().Add(hubUniq)
	errExists := st.Hub().Add(hubExists)

	// Assert
	assert.NoError(t, errUniq)
	assert.Error(t, errExists)
}

// TestHubRepository_Remove remove exists and not exists hub
func TestHubRepository_Remove(t *testing.T) {
	// Arrange
	idNotExists := "notExists"
	idExists := "exists"

	st := memoryStore.New()

	hubExists := &model.Hub{Id: idExists}
	hubNotExists := &model.Hub{Id: idNotExists}

	st.Hub().Add(hubExists)

	// Act
	errExists := st.Hub().Remove(hubExists.Id)
	errNotExists := st.Hub().Remove(hubNotExists.Id)

	// Assert
	assert.NoError(t, errExists)
	assert.Error(t, errNotExists)
}

// TestHubRepository_Find find added, not exists and deleted hub
func TestHubRepository_Find(t *testing.T) {
	// Arrange
	idAdded := "Added"
	idNotAdded := "NotAdded"
	idRemoved := "Removed"

	st := memoryStore.New()

	addedHub := &model.Hub{Id: idAdded}
	st.Hub().Add(addedHub)

	removedHub := &model.Hub{Id: idRemoved}
	st.Hub().Add(removedHub)
	st.Hub().Remove(removedHub.Id)

	// Act
	added, errAdded := st.Hub().Find(idAdded)
	notAdded, errNoAdded := st.Hub().Find(idNotAdded)
	removed, errRemoved := st.Hub().Find(idRemoved)

	// Assert

	// Added
	assert.NoError(t, errAdded)
	assert.NotNil(t, added)

	// Not Added
	assert.Error(t, errNoAdded)
	assert.Nil(t, notAdded)

	// Deleted
	assert.Error(t, errRemoved)
	assert.Nil(t, removed)
}
