package server

import (
	"log/slog"
	"net/http"
	"time"

	mv "github.com/go-chi/chi/v5/middleware"
	"github.com/urfave/negroni"
)

func loggingMiddleware(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := r.Context().Value(mv.RequestIDKey)

			log.Info("request started",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Any("request_id", requestID),
			)

			lrw := negroni.NewResponseWriter(w)
			next.ServeHTTP(lrw, r)

			log.Info("request finished",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Any("request_id", requestID),
				slog.Int("status", lrw.Status()),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

func contentTypeJsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
