package store

import (
	"Chat/internal/app/model"
)

// UserRepository ...
type UserRepository interface {
	Add(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	Find(id int) (*model.User, error)
	Deactivate(id int) error
	Activate(id int) error
}
