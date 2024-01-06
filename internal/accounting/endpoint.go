package accounting

import (
	"context"
	"fmt"
	"github.com/itsnoproblem/prmry/internal/components/profile"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Service interface {
	CreateAPIKey(ctx context.Context) (auth.APIKey, error)
	SetAPIKeyName(ctx context.Context, keyID, name string) error
	GetAPIKeys(ctx context.Context, userID string) ([]auth.APIKey, error)
	DeleteAPIKey(ctx context.Context, keyID string) error
}

type accountResponse struct {
	Provider  string
	Email     string
	Name      string
	AvatarURL string
	APIKeys   []auth.APIKey
}

func makeGetAccountEndpoint(s Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user := auth.UserFromContext(ctx)
		if user == nil {

			return nil, fmt.Errorf("missing user")
		}

		apiKeys, err := s.GetAPIKeys(ctx, user.ID)
		if err != nil {
			return nil, errors.Wrap(err, "makeGetAccountEndpoint")
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

type updateAPIKeyRequest struct {
	KeyID string
	Name  string
}

type updateAPIKeyResponse struct {
	Name string
}

func makeUpdateAPIKeyEndpoint(s Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(updateAPIKeyRequest)
		if !ok {
			return nil, fmt.Errorf("makeUpdateAPIKeyEndpoint: failed to parse request")
		}

		if err := s.SetAPIKeyName(ctx, req.KeyID, req.Name); err != nil {
			return nil, errors.Wrap(err, "makeUpdateAPIKeyEndpoint")
		}

		return nil, nil
	}
}

func makeCreateAPIKeyEndpoint(s Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		newKey, err := s.CreateAPIKey(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "makeCreateAPIKeyEndpoint")
		}

		return profile.APIKeyView{
			Name: newKey.Name,
			Key:  newKey.Key,
		}, nil
	}
}

type deleteAPIKeyRequest struct {
	KeyID string
}

func makeDeleteAPIKeyEndpoint(s Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(deleteAPIKeyRequest)
		if !ok {
			return nil, fmt.Errorf("makeDeleteAPIKeyEndpoint: failed to parse request")
		}

		if err := s.DeleteAPIKey(ctx, req.KeyID); err != nil {
			return nil, errors.Wrap(err, "makeDeleteAPIKeyEndpoint")
		}

		return nil, nil
	}
}
