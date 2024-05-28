package testStore_test

import (
	"Chat/internal/app/model"
	"Chat/internal/app/store/testStore"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Add(t *testing.T) {
	// Arrange
	st := testStore.New()
	user := model.TestUser(t)

	// Act
	err := st.User().Add(user)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_FindByEmailNotAdded(t *testing.T) {
	// Arrange
	st := testStore.New()

	// Act
	_, err := st.User().FindByEmail("user@ex.com")

	// Assert
	assert.EqualError(t, err, model.RecordNotFound.Error())
}

func TestUserRepository_FindByEmailAdded(t *testing.T) {
	// Arrange
	st := testStore.New()
	email := "user@ex.com"
	user := model.TestUser(t)
	user.Email = email

	// Act
	st.User().Add(user)
	user, err := st.User().FindByEmail(email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindNotAdded(t *testing.T) {
	// Arrange
	st := testStore.New()

	// Act
	_, err := st.User().Find(0)

	// Assert
	assert.EqualError(t, err, model.RecordNotFound.Error())
}

func TestUserRepository_FindAdded(t *testing.T) {
	// Arrange
	st := testStore.New()
	user := model.TestUser(t)

	// Act
	st.User().Add(user)
	user, err := st.User().Find(user.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_Deactivate(t *testing.T) {
	// Arrange
	st := testStore.New()
	user := model.TestUser(t)

	// Act
	st.User().Add(user)
	err := st.User().Deactivate(user.ID)
	user, _ = st.User().Find(user.ID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, user.Active)
}

func TestUserRepository_Activate(t *testing.T) {
	// Arrange
	st := testStore.New()
	user := model.TestUser(t)

	// Act
	st.User().Add(user)
	err := st.User().Activate(user.ID)
	user, _ = st.User().Find(user.ID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, user.Active)
}
