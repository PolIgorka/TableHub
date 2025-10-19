package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/lo"

	"github.com/dbaas/internal/server"
	"github.com/dbaas/internal/storage"
	db "github.com/dbaas/pkg/database"
)

func main() {
	ctx := context.Background()
	config := lo.Must(LoadConfig())

	managedDB := db.New(config.DB)
	defer managedDB.Close(ctx)

	userStorage := storage.NewUserRightsStorage(managedDB)

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
