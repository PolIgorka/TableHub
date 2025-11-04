package storage

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrLoginTooLong    = errors.New("login is too long")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserAlreadyExists = errors.New("user already exists")
)
