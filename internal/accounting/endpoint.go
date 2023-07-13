package accounting

import (
	"context"
	"fmt"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/htmx"
)

type accountResponse struct {
	Provider  string
	Email     string
	Name      string
	AvatarURL string
}

func makeAccountEndpoint() htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, fmt.Errorf("missing user")
		}

		return accountResponse{
			Provider:  user.Provider,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
		}, nil
	}
}
