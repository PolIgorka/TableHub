package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"

	db "github.com/tablehub/pkg/database"
)

type UsersTableDTO struct {
	bun.BaseModel `bun:"table:users"`

	UserID		 uuid.UUID 		   `bun:"user_id,type:uuid,pk,default:gen_random_uuid()"`
	Login		 string			   `bun:"login,type:varchar(32)"`
	PasswordHash string			   `bun:"password_hash,type:varchar(255)"`
	Rights		 []*RightsTableDTO `bun:"rel:has-many,join:user_id=user_id"`
}

type UsersStorage interface {
	CheckUserExist(ctx context.Context, login string, password string) error
	CreateUser(ctx context.Context, login string, password string) error
}

type userStorageImpl struct {
	db *db.ManagedDB
}

func NewUsersStorage(db *db.ManagedDB) UsersStorage {
	return &userStorageImpl{db: db}
}

func (storage *userStorageImpl) CheckUserExist(ctx context.Context, login string, password string) error {
	user := new(UsersTableDTO)
	err := storage.db.Manager.NewSelect().Model(user).Where("login = ?", login).Scan(ctx)

	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return ErrIncorrectPassword
	}

	return nil
}

func (storage *userStorageImpl) CreateUser(ctx context.Context, login string, password string) error {
	if len(login) > 32 {
		return ErrLoginTooLong
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	newUser := &UsersTableDTO{
		Login:         login,
		PasswordHash:     string(passwordHash),
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
