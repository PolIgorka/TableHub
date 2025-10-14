package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/dbaas/internal/logger"
)

type Server struct {
	httpServer *http.Server
	config     ServerConfig
	logger     *slog.Logger
}

func New(config ServerConfig) *Server {
	mux := http.NewServeMux()
	logger := logger.New("http-server")

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: loggingMiddleware(mux, logger),
	}

	logger.Info("server instance created")

	return &Server{
		httpServer: srv,
		config:     config,
		logger:     logger,
	}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
    ln, err := net.Listen("tcp", s.httpServer.Addr)
    if err != nil {
        s.logger.Error("failed to listen", "addr", s.httpServer.Addr, "error", err)
        return err
    }

    s.logger.Info("listening", "addr", ln.Addr().String())

    go func() {
        <-ctx.Done()
        s.logger.Info("context canceled, shutting down")
        _ = s.httpServer.Shutdown(context.Background())
    }()

    return s.httpServer.Serve(ln)
}


func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server gracefully")
	return s.httpServer.Shutdown(ctx)
}
