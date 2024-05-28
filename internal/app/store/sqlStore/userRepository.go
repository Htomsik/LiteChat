package sqlStore

import (
	"Chat/internal/app/model"
	"database/sql"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Add ...
func (repository *UserRepository) Add(user *model.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := user.BeforeAdd(); err != nil {
		return err
	}

	return repository.store.db.QueryRow(
		"INSERT INTO users(email,encryptedPassword) values ($1,$2) RETURNING id",
		user.Email,
		user.EncryptedPassword,
	).Scan(&user.ID)
}

// FindByEmail ...
func (repository *UserRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}

	err := repository.store.db.QueryRow(
		"SELECT id,email,encryptedPassword FROM users WHERE email = $1",
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.EncryptedPassword,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			err = model.RecordNotFound
		}
	}

	return user, err
}

// Find ...
func (repository *UserRepository) Find(id int) (*model.User, error) {
	user := &model.User{}

	err := repository.store.db.QueryRow(
		"SELECT id,email,encryptedPassword, active FROM users WHERE id = $1",
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.EncryptedPassword,
		&user.Active,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			err = model.RecordNotFound
		}
	}

	return user, err
}

// Deactivate ...
func (repository *UserRepository) Deactivate(id int) error {

	err := repository.store.db.QueryRow(
		"UPDATE users SET	active = $1 WHERE id = $2 RETURNING id",
		false,
		id,
	).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			err = model.RecordNotFound
		}
	}

	return err
}

// Activate ...
func (repository *UserRepository) Activate(id int) error {

	err := repository.store.db.QueryRow(
		"UPDATE users SET	active = $1 WHERE id = $2 RETURNING id",
		true,
		id,
	).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			err = model.RecordNotFound
		}
	}

	return err
}
