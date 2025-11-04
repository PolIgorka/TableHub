package storage

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	db "github.com/tablehub/pkg/database"
)

type RightsTableDTO struct {
	bun.BaseModel `bun:"table:rights"`

	RightID          uuid.UUID 		 `bun:"right_id,type:uuid,pk"`
	UserID           uuid.UUID 		 `bun:"user_id,type:uuid"`
	TableID          uuid.UUID 		 `bun:"table_id,type:uuid"`
	RightsMask       int       		 `bun:"rights_mask,type:integer"`
	AvailableColumns []string        `bun:"avaliable_columns,type:jsonb"`
	User             *UsersTableDTO  `bun:"rel:belongs-to,join:user_id=user_id"`
	Table            *TablesTableDTO `bun:"rel:belongs-to,join:table_id=table_id"`
}

type RightsStorage interface {}

type rightsStorageImpl struct {
	db *db.ManagedDB
}

func NewUserRightsStorage(db *db.ManagedDB) RightsStorage {
	return &rightsStorageImpl{db: db}
}
