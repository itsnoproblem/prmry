package accounting

import (
	"context"
	"fmt"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Service interface {
	GetAPIKeys(ctx context.Context, userID string) ([]auth.APIKey, error)
}

type accountResponse struct {
	Provider  string
	Email     string
	Name      string
	AvatarURL string
	APIKeys   []auth.APIKey
}

func makeAccountEndpoint(s Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user := auth.UserFromContext(ctx)
		if user == nil {

			return nil, fmt.Errorf("missing user")
		}

		apiKeys, err := s.GetAPIKeys(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "makeAccountEndpoint")
		}

		return accountResponse{
			Provider:  user.Provider,
			Email:     user.Email,
			Name:      user.Name,
			AvatarURL: user.AvatarURL,
			APIKeys:   apiKeys,
		}, nil
	}
}
