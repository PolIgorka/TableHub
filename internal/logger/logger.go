package logger

import (
	"log/slog"
	"os"
)

func New(loggerName string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})

	return slog.New(handler).With("logger-name", loggerName)
}
