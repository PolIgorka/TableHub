package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/tablehub/pkg/logger"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Type 	 string `yaml:"type"`
}

type ManagedDB struct {
	Manager *bun.DB
	logger  *slog.Logger
}

func New(config DBConfig) *ManagedDB {
	logger := logger.New(fmt.Sprintf("database.%s", config.Name))
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())
	if config.Type == "debug" {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	err := db.PingContext(context.Background())
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	logger.Info("database connection established")

	return &ManagedDB{Manager: db, logger: logger}
}

func (s *ManagedDB) Close(ctx context.Context) error {
	s.logger.Info("closing database connection")
	return s.Manager.Close()
}
