package user

import (
	"errors"
	"fmt"

	"github.com/carfloresf/financial-chat/internal/storage"
)

var ErrPasswordIncorrect = errors.New("incorrect password provided")

type User struct {
	db storage.Storage
	ph PasswordHasher
}

func NewUser(ph PasswordHasher, db storage.Storage) *User {
	return &User{
		db: db,
		ph: ph,
	}
}

func (u *User) Register(username, password string) error {
	hash, err := u.ph.GenerateHash(password)
	if err != nil {
		return fmt.Errorf("error generating password hash: %w", err)
	}

	_, err = u.db.CreateUser(username, hash)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

var ErrPasswordEmpty = errors.New("password is empty")

func (u *User) Authenticate(username, password string) error {
	mu, err := u.db.GetUser(username)
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	err = u.ph.CompareHashAndPassword(mu.PasswordHash, password)
	if err != nil {
		return fmt.Errorf("error comparing password: %w", err)
	}

	return err
}
