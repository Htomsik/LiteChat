package model_test

import (
	"Chat/internal/app/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestUser_Validate ...
func TestUser_Validate(t *testing.T) {

	testCases := []struct {
		name    string
		user    func() *model.User
		IsValid bool
	}{
		{
			name: "valid",
			user: func() *model.User {
				return model.TestUser(t)
			},
			IsValid: true,
		},

		// Email cases
		{
			name: "invalid email",
			user: func() *model.User {
				us := model.TestUser(t)

				us.Email = "NoEmail"

				return us
			},
			IsValid: false,
		},
		{
			name: "empty email",
			user: func() *model.User {
				us := model.TestUser(t)

				us.Email = ""

				return us
			},
			IsValid: false,
		},

		// Password cases
		{
			name: "empty password",
			user: func() *model.User {
				us := model.TestUser(t)

				us.Password = ""

				return us
			},
			IsValid: false,
		},
		{
			name: "short password",
			user: func() *model.User {
				us := model.TestUser(t)

				us.Password = "12345"

				return us
			},
			IsValid: false,
		},
		{
			name: "with encrypted password",
			user: func() *model.User {
				us := model.TestUser(t)

				us.Password = ""

				us.EncryptedPassword = "SomeEncryptedPassword"

				return us
			},
			IsValid: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.IsValid {
				assert.NoError(t, testCase.user().Validate())
			} else {
				assert.Error(t, testCase.user().Validate())
			}
		})
	}
}

// TestUser_BeforeAdd ...
func TestUser_BeforeAdd(t *testing.T) {
	user := model.TestUser(t)

	assert.NoError(t, user.BeforeAdd())
	assert.NotEmpty(t, user.EncryptedPassword)
}
