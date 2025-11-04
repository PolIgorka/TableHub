package storage

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"

	db "github.com/tablehub/pkg/database"
)

type TablesTableDTO struct {
	bun.BaseModel `bun:"table:tables"`

	TableID uuid.UUID `bun:"table_id,type:uuid,pk"`
	Name    string    `bun:"name,type:varchar(32)"`
	Columns []string   `bun:"columns,type:jsonb"`
	Rights  []*RightsTableDTO `bun:"rel:has-many,join:table_id=table_id"`
}

type TablesStorage interface {}

type tablesStorageImpl struct {
	db *db.ManagedDB
}

func NewTablesStorage(db *db.ManagedDB) TablesStorage {
	return &tablesStorageImpl{db: db}
}
