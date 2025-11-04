package routes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5"

	db "github.com/tablehub/internal/storage"
	logs "github.com/tablehub/pkg/logger"
)

type AuthRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}


func LoginRoutes(usersStorage db.UsersStorage, logger *slog.Logger) func(chi.Router) {
	return func(r chi.Router) {
		r.Post("/login", loginHandler(usersStorage, logger))
		r.Post("/register", registerHandler(usersStorage, logger))
	}
}

func loginHandler(storage db.UsersStorage, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logs.Inject(logger, ctx)

		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Failed to read request body", "error", err)
			http.Error(w, CodeInternalServerError, http.StatusInternalServerError)
			return
		}
		parsedBody := new(AuthRequestBody)
		err = json.Unmarshal(rawBody, parsedBody)
		if err != nil {
			logger.Error("Failed to unmarshal request body", "error", err)
			http.Error(w, CodeInternalServerError, http.StatusInternalServerError)
			return
		}

		err = storage.CheckUserExist(ctx, parsedBody.Login, parsedBody.Password)
		switch {
		case errors.Is(err, db.ErrUserNotFound):
			http.Error(w, CodeUnauthorized, http.StatusUnauthorized)
			return
		case errors.Is(err, db.ErrIncorrectPassword):
			http.Error(w, CodeIncorrectPassword, http.StatusUnauthorized)
			return
		case err != nil:
			logger.Error("Failed to login", "error", err.Error())
			http.Error(w, CodeInternalServerError, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func registerHandler(storage db.UsersStorage, logger *slog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := logs.Inject(logger, ctx)

		rawBody, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Failed to read request body", "error", err)
			http.Error(w, CodeInternalServerError, http.StatusInternalServerError)
			return
		}
		parsedBody := new(AuthRequestBody)
		err = json.Unmarshal(rawBody, parsedBody)
		if err != nil {
			logger.Error("Failed to unmarshal request body", "error", err)
			http.Error(w, CodeInternalServerError, http.StatusInternalServerError)
			return
		}

		err = storage.CreateUser(ctx, parsedBody.Login, parsedBody.Password)
		switch {
		case errors.Is(err, db.ErrUserAlreadyExists):
			http.Error(w, CodeUserExists, http.StatusUnprocessableEntity)
			return
		case errors.Is(err, db.ErrLoginTooLong):
			http.Error(w, CodeLongLogin, http.StatusUnprocessableEntity)
			return
		case err != nil:
			logger.Error("Failed to create user", "error", err.Error())
			http.Error(w, CodeInternalServerError, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
