package accounting

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components/profile"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Service interface {
	GetUser(ctx context.Context, userID string) (auth.User, error)
	UpdateAccountProfile(ctx context.Context, userID, name, email string) error
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

type getAccountRequest struct {
	UserID string
}

func makeGetAccountEndpoint(s Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getAccountRequest)
		// always get latest user info
		user, err := s.GetUser(ctx, req.UserID)
		if err != nil {
			return nil, errors.Wrap(err, "makeGetAccountEndpoint")
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

type profileOptions struct {
	UserID string
	Name   string
	Email  string
}

func (o profileOptions) Validate() error {
	if o.UserID == "" {
		return fmt.Errorf("missing user ID")
	}

	if o.Name == "" {
		return fmt.Errorf("missing name")
	}

	if o.Email == "" {
		return fmt.Errorf("missing email")
	}

	if _, err := mail.ParseAddress(o.Email); err != nil {
		return fmt.Errorf("email address %s is not valid", o.Email)
	}

	return nil
}

func makeUpdateProfileEndpoint(s Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(profileOptions)
		if !ok {
			return nil, fmt.Errorf("makeUpdateProfileEndpoint: failed to parse request")
		}

		if err := req.Validate(); err != nil {
			return nil, errors.Wrap(err, "makeUpdateProfileEndpoint")
		}

		user := auth.UserFromContext(ctx)
		if user == nil {
			return nil, fmt.Errorf("missing user")
		}

		if err := s.UpdateAccountProfile(ctx, user.ID, req.Name, req.Email); err != nil {
			return nil, errors.Wrap(err, "makeUpdateProfileEndpoint")
		}

		return nil, nil
	}
}

type updateAPIKeyOptions struct {
	KeyID string
	Name  string
}

func makeUpdateAPIKeyEndpoint(s Service) internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(updateAPIKeyOptions)
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
