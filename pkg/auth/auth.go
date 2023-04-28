package auth

import (
	"context"
	"log"
	"sort"
)

const (
	CookieName = "rgb_user"
	ContextKey = "User"
)

type Byte32 []byte

func Providers() (keys []string, providersByName map[string]string) {
	providersByName = map[string]string{
		"github": "GitHub",
		"google": "Google",
	}

	for k := range providersByName {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys, providersByName
}

func UserFromContext(ctx context.Context) *User {
	userFromContext := ctx.Value(ContextKey)
	u, isUser := userFromContext.(User)
	if !isUser || u.ID == "" {
		log.Printf("auth.UserFromContext: failed to cast user")
		return nil
	}

	return &u
}
