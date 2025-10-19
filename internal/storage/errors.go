package storage

import (
	"errors"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserAlreadyExists = errors.New("user already exists")
)
