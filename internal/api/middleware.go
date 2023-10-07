package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/itsnoproblem/prmry/internal/auth"
)

const (
	AuthContextKey   = "X-API-KEY"
	AuthErrorMessage = "API key is invalid"
)

type AuthProvider interface {
	FindUserByAPIKey(ctx context.Context, apiKey string) (usr auth.User, exists bool, err error)
}

func IsAPIPath(path string) bool {
	return strings.HasPrefix(path, "/api")
}

func Middleware(userProvider AuthProvider, rndr Renderer) func(http.Handler) http.Handler {
	handlerMaker := func(h http.Handler) http.Handler {
		handler := func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(AuthContextKey)
			usr, exists, err := userProvider.FindUserByAPIKey(r.Context(), apiKey)
			if err != nil {
				rndr.RenderError(w, r, err)
				return
			}

			if !exists && IsAPIPath(r.URL.Path) {
				rndr.RenderError(w, r, fmt.Errorf(AuthErrorMessage))
				return
			} else if exists {
				ctx := context.WithValue(r.Context(), auth.ContextKey, usr)
				r = r.WithContext(ctx)
			}

			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(handler)
	}
	return handlerMaker
}
