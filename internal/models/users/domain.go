package users

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Text *string
	Hash []byte
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Pwd       Password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

// Payloads

type CreateUserPayload struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Pwd   string `json:"password"`
}

type UpdateUserPayload struct {
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	Pwd       *string `json:"password"`
	Activated *bool   `json:"activated"`
}

// Password

func (p *Password) Set(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	if err != nil {
		return err
	}

	p.Hash = hash
	p.Text = &pwd

	return nil
}

func (p *Password) Check(pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(pwd))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
