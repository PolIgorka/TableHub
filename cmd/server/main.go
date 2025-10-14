package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/samber/lo"

	"github.com/dbaas/internal/server"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}
	config := lo.Must(LoadConfig(configPath))

	srv := server.New(config.Server)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	lo.Must0(srv.ListenAndServe(ctx))
}
