package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/lo"

	"github.com/tablehub/internal/server"
	"github.com/tablehub/internal/storage"
	db "github.com/tablehub/pkg/database"
)

func main() {
	ctx := context.Background()
	config := lo.Must(LoadConfig())

	rightsDB := db.New(config.DB.RightsConfig)
	defer rightsDB.Close(ctx)
	tablesDB := db.New(config.DB.TablesConfig)
	defer tablesDB.Close(ctx)

	userStorage := storage.NewUserRightsStorage(rightsDB)

	app := server.New(config.Server, userStorage)

	go func() {
		if err := app.Run(); err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	app.Shutdown(ctx)
}
