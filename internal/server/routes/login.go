package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"

	db "github.com/tablehub/internal/storage"
	logs "github.com/tablehub/pkg/logger"
)

type LoginResponseBody struct {
	AvailableColumns []string `json:"available_columns"`
}

func LoginRoutes(usersStorage db.UserRightsStorage, logger *slog.Logger) func(chi.Router) {
	return func(r chi.Router) {
		r.Post("/login", loginHandler(usersStorage, logger))
		r.Post("/register", registerHandler(usersStorage, logger))
	}
}

func loginHandler(storage db.UserRightsStorage, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logs.Inject(logger, ctx)

		username, password, err := extractBasicAuth(r)
		if err != nil {
			http.Error(w, CodeUnauthorized, http.StatusUnauthorized)
			return
		}

		user, err := storage.GetUserByLogin(ctx, username, []byte(password))
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			http.Error(w, CodeUnauthorized, http.StatusUnauthorized)
			return
		case errors.Is(err, db.ErrIncorrectPassword):
			http.Error(w, CodeIncorrectLogin, http.StatusUnauthorized)
			return
		case err != nil:
			logger.Error("failed to login", "error", err.Error())
			http.Error(w, CodeInternalServerError, http.StatusInternalServerError)
			return
		}

		response := LoginResponseBody{
			AvailableColumns: user.AvailableColumns,
		}
		responseBodyJson, err := json.Marshal(response)
		if err != nil {
			logger.Error("Failed to serialize response body", "error", err)
			http.Error(w, CodeInternalServerError, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(responseBodyJson)
	}
}

func registerHandler(storage db.UserRightsStorage, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logs.Inject(logger, ctx)

		username, password, err := extractBasicAuth(r)
		if err != nil {
			http.Error(w, CodeUnauthorized, http.StatusUnauthorized)
			return
		}

		err = storage.CreateUser(ctx, username, []byte(password))
		switch {
		case errors.Is(err, db.ErrUserAlreadyExists):
			http.Error(w, CodeUserExists, http.StatusUnprocessableEntity)
			return
		case err != nil:
			logger.Error("failed to create user", "error", err.Error())
			http.Error(w, CodeInternalServerError, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
