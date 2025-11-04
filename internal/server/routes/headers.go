package routes

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

var (
	BasicAuthPrefix = "Basic "
	AuthHeader      = "Authorization"
	ErrMissingAuth  = errors.New("missing basic authorization header")
	ErrInvalidAuth  = errors.New("invalid auth format")
)

func extractBasicAuth(r *http.Request) (string, string, error) {
	authHeader := r.Header.Get(AuthHeader)
	if authHeader == "" || !strings.HasPrefix(authHeader, BasicAuthPrefix) {
		return "", "", ErrMissingAuth
	}
	encoded := strings.TrimPrefix(authHeader, BasicAuthPrefix)
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", err
	}

	parts := strings.SplitN(string(decodedBytes), ":", 2)
	if len(parts) != 2 {
		return "", "", ErrInvalidAuth
	}

	return parts[0], parts[1], nil
}
