package logger

import (
	"context"
	"log/slog"
	"os"
)

func New(loggerName string) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})

	return slog.New(handler).With("logger-name", loggerName)
}

func Inject(logger *slog.Logger, ctx context.Context) *slog.Logger {
	traceID, ok := ctx.Value(TraceID).(string)
	if ok && traceID != "" {
		logger = logger.With(TraceID, traceID)
	}

	return logger
}
