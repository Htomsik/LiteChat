package sqlStore_test

import (
	"Chat/internal/app/model"
	"Chat/internal/app/store/sqlStore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Add(t *testing.T) {
	// Arrange
	db, teardown := sqlStore.TestDb(t, databaseType, databaseURL)
	defer teardown("users")

	// Act
	st := sqlStore.New(db)
	user := model.TestUser(t)
	err := st.User().Add(user)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func TestUserRepository_FindByEmailNotAdded(t *testing.T) {
	// Arrange
	db, teardown := sqlStore.TestDb(t, databaseType, databaseURL)
	defer teardown("users")
	st := sqlStore.New(db)

	// Act
	_, err := st.User().FindByEmail("user@ex.com")

	// Assert
	assert.EqualError(t, err, model.RecordNotFound.Error())
}

func TestUserRepository_FindByEmailAdded(t *testing.T) {
	// Arrange
	db, teardown := sqlStore.TestDb(t, databaseType, databaseURL)
	defer teardown("users")
	st := sqlStore.New(db)
	user := model.TestUser(t)

	// Act
	st.User().Add(user)
	user, err := st.User().FindByEmail(user.Email)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserRepository_FindNotAdded(t *testing.T) {
	// Arrange
	db, teardown := sqlStore.TestDb(t, databaseType, databaseURL)
	defer teardown("users")
	st := sqlStore.New(db)

	// Act
	_, err := st.User().Find(0)

	// Assert
	assert.EqualError(t, err, model.RecordNotFound.Error())
}

func TestUserRepository_Find(t *testing.T) {
	// Arrange
	db, teardown := sqlStore.TestDb(t, databaseType, databaseURL)
	defer teardown("users")
	st := sqlStore.New(db)
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
	db, teardown := sqlStore.TestDb(t, databaseType, databaseURL)
	defer teardown("users")
	st := sqlStore.New(db)
	user := model.TestUser(t)

	// Assert
	st.User().Add(user)
	err := st.User().Deactivate(user.ID)

	// Assert
	assert.NoError(t, err)
}

func TestUserRepository_Activate(t *testing.T) {
	// Arrange
	db, teardown := sqlStore.TestDb(t, databaseType, databaseURL)
	defer teardown("users")
	st := sqlStore.New(db)
	user := model.TestUser(t)

	// Assert
	st.User().Add(user)
	err := st.User().Activate(user.ID)

	// Assert
	assert.NoError(t, err)
}
