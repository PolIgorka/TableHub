package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	db "github.com/tablehub/pkg/database"
)

type UserRight struct {
	bun.BaseModel `bun:"table:user_rights"`

	Username         string   `bun:"username,type:varchar(255),pk"`
	PasswordHash     string   `bun:"password_hash,type:varchar(255)"`
	AvailableColumns []string `bun:"available_columns,type:jsonb"`
}

type UserRightsStorage interface {
	GetUserByLogin(ctx context.Context, login string, password []byte) (*UserRight, error)
	CreateUser(ctx context.Context, login string, password []byte) error
}

type userRightsStorage struct {
	db *db.ManagedDB
}

func NewUserRightsStorage(db *db.ManagedDB) UserRightsStorage {
	return &userRightsStorage{db: db}
}

func (storage *userRightsStorage) GetUserByLogin(ctx context.Context, login string, password []byte) (*UserRight, error) {
	user := new(UserRight)
	err := storage.db.Manager.NewSelect().Model(user).Where("username = ?", login).Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), password); err != nil {
		return nil, ErrIncorrectPassword
	}

	return user, nil
}

func (storage *userRightsStorage) CreateUser(ctx context.Context, login string, password []byte) error {
	passwordHash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := &UserRight{
		Username:         login,
		PasswordHash:     string(passwordHash),
		AvailableColumns: []string{},
	}
	_, err = storage.db.Manager.NewInsert().Model(newUser).Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return ErrUserAlreadyExists
		}
		return err
	}

	return nil
}
