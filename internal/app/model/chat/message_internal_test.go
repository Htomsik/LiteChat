package chat

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMessage_ClearPrivacy(t *testing.T) {
	// Arrange
	user1 := NewChatUser("user1")
	user2 := NewChatUser("user2")
	users := []*User{user1, user2}

	msg := NewSystemMessage(TypeUsersList, users)

	// Act
	ok := msg.ClearPrivacy(user1)

	// Assert
	assert.True(t, ok)

	cleared, ok := msg.Message.([]User)

	assert.True(t, ok)
	assert.Equal(t, user1.Id, cleared[0].Id)
	assert.Equal(t, uuid.Nil, cleared[1].Id)
}
