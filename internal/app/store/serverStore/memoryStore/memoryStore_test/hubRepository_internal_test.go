package memoryStore_test

import (
	"Chat/internal/app/model"
	"Chat/internal/app/model/chat"
	"Chat/internal/app/store/serverStore/memoryStore"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	hubCfg = &model.HubConfig{

		AdminRole: &chat.UserRole{
			Name:    "Admin",
			IsAdmin: true,
		},

		DefaultRole: &chat.UserRole{
			Name:        "User",
			IsAdmin:     false,
			Permissions: []chat.RolePermission{0},
		},

		Roles: nil,
	}
)

// TestMain Test Startup
func TestMain(m *testing.M) {

	// init roles
	if hubCfg.Roles == nil {
		hubCfg.Roles = []*chat.UserRole{}
		hubCfg.Roles = append(hubCfg.Roles, hubCfg.DefaultRole)
		hubCfg.Roles = append(hubCfg.Roles, hubCfg.AdminRole)
	}

	// Exit from tests
	exitVal := m.Run()
	os.Exit(exitVal)
}

// TestHubRepository_Create create exists and not exists hub
func TestHubRepository_Create(t *testing.T) {
	// Arrange
	idUniq := "new"
	idExists := "exists"

	st := memoryStore.New()

	st.Hub().Create(idExists, hubCfg)

	// Act
	hubUniq, errUniq := st.Hub().Create(idUniq, hubCfg)
	hubExists, errExists := st.Hub().Create(idExists, hubCfg)

	// Assert
	assert.NotNil(t, hubUniq)
	assert.NoError(t, errUniq)

	assert.Nil(t, hubExists)
	assert.Error(t, errExists)
}

// TestHubRepository_Add add exists and not exists hub
func TestHubRepository_Add(t *testing.T) {
	// Arrange
	idUniq := "new"
	idExists := "exists"

	st := memoryStore.New()

	hubExists := model.HewHub(idUniq, logrus.New(), make(chan string), hubCfg)
	hubUniq := model.HewHub(idExists, logrus.New(), make(chan string), hubCfg)

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

	hubExists := model.HewHub(idExists, logrus.New(), make(chan string), hubCfg)
	hubNotExists := model.HewHub(idNotExists, logrus.New(), make(chan string), hubCfg)

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

	addedHub := model.HewHub(idAdded, logrus.New(), make(chan string), hubCfg)
	st.Hub().Add(addedHub)

	removedHub := model.HewHub(idRemoved, logrus.New(), make(chan string), hubCfg)
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
