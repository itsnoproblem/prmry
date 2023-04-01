package auth

import (
	"context"
	"sort"

	"github.com/markbates/goth"
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
	u, isUser := userFromContext.(goth.User)
	if !isUser || u.UserID == "" {
		return nil
	}

	return &User{
		ID:        u.UserID,
		Name:      u.Name,
		Nickname:  u.NickName,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
		Provider:  u.Provider,
	}
}
