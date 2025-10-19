package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/dbaas/internal/server/routes"
	"github.com/dbaas/internal/storage"
	"github.com/dbaas/pkg/logger"
)

type AppServer struct {
	router chi.Router
	server *http.Server
	logger *slog.Logger
}

func New(cfg ServerConfig, db storage.UserRightsStorage) *AppServer {
	logger := logger.New("http-server")

	r := chi.NewRouter()

	// --- Middleware ---
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(loggingMiddleware(logger))
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(contentTypeJsonMiddleware)

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(routes.LoginRoutes(db, logger))
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &AppServer{
		router: r,
		server: srv,
		logger: logger,
	}
}

func (s *AppServer) Run() error {
	s.logger.Info("starting server", slog.String("addr", s.server.Addr))
	return s.server.ListenAndServe()
}

func (s *AppServer) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	return s.server.Shutdown(ctx)
}
