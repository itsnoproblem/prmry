package auth

import (
	"context"
	"net/http"
)

func TestModeUser() User {
	return User{
		ID:         "user-123",
		Email:      "user@examples.com",
		Name:       "Test User",
		Nickname:   "testy",
		AvatarURL:  "https://example.com/img/user-123.jpg",
		Provider:   "github",
		ProviderID: "123456789",
	}
}

func TestUserMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ContextKey, TestModeUser())
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
