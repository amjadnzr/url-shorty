package models

import (
	"errors"
	"net/mail"
)

type User struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"-"`
	Token        string `json:"token"`
}

func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("user name is required")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	if _, err := mail.ParseAddress(u.Email); err != nil {
		return err
	}
	return nil
}
