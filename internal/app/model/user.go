package model

import (
	val "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	ID                int    `json:"id"`
	Email             string `json:"email"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
	Active            bool   `json:"-"`
}

// UserShort short version for request
type UserShort struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate ...
func (user *User) Validate() error {
	return val.ValidateStruct(
		user,
		val.Field(&user.Email, val.Required, is.Email),
		val.Field(&user.Password, val.By(requiredIf(user.EncryptedPassword == "")), val.Length(6, 100)),
	)
}

// BeforeAdd ...
func (user *User) BeforeAdd() error {
	if len(user.Password) > 0 {
		encryptedPassword, err := encryptString(user.Password)

		if err != nil {
			return err
		}

		user.EncryptedPassword = encryptedPassword
	}

	return nil
}

// ClearPrivate cleaning private data before push
func (user *User) ClearPrivate() {
	user.Password = ""
}

// ComparePassword compare hashed password in the database and passed password
func (user *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password)) == nil
}

// encryptString
func encryptString(text string) (string, error) {
	encryptedText, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(encryptedText), err
}
