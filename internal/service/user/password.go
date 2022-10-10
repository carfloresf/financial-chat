package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	GenerateHash(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

type PasswordGenerator struct {
	pepper string
}

func NewPasswordGenerator(pepper string) *PasswordGenerator {
	return &PasswordGenerator{pepper: pepper}
}

func (u *PasswordGenerator) GenerateHash(password string) (string, error) {
	if password == "" {
		return "", ErrPasswordEmpty
	}

	pwBytes := []byte(password + u.pepper)

	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

func (u *PasswordGenerator) CompareHashAndPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password+u.pepper))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrPasswordIncorrect
		}

		return err
	}

	return nil
}
