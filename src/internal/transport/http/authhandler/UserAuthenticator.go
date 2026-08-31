package authhandler

import (
	"WebTic-tac-toe2/internal/service/jwtauth"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	service jwtauth.JWTService
}

func NewUserAuthenticator(s jwtauth.JWTService) *UserAuthenticator {
	return &UserAuthenticator{service: s}
}

type contextKey string

const UserUUIDKey contextKey = "userUUID"

func (ua *UserAuthenticator) Middleware(f http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		fmt.Printf("[DEBUG] Входящий запрос: %s %s\n", r.Method, path)

		if !strings.HasPrefix(path, "/api/game/") {
			f.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, "Bearer auth required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == "" {
			http.Error(w, "Empty token", http.StatusUnauthorized)
			return
		}

		err := ua.service.IsValidAccess(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Trouble with token, err: %q", err), http.StatusUnauthorized)
			return
		}

		uUUID, err := ua.service.GetUUIDfromAccessToken(token)
		if err != nil {
			http.Error(w, fmt.Sprintf("Trouble with getting UUID from token, err: %q", err), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), UserUUIDKey, uUUID)
		f.ServeHTTP(w, r.WithContext(ctx))
	})
}
