package accounting

import (
	"context"
	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/pkg/errors"
	"time"
)

type UserRepo interface {
	FindAPIKeysForUser(ctx context.Context, userID string) ([]auth.APIKey, error)
	UpdateAPIKeyName(ctx context.Context, userID, keyID, name string) error
	InsertAPIKey(ctx context.Context, userID string, key auth.APIKey) error
	DeleteAPIKey(ctx context.Context, userID, keyID string) error
}

type service struct {
	userRepo UserRepo
}

func NewService(userRepo UserRepo) *service {
	return &service{
		userRepo: userRepo,
	}
}

func (s *service) GetAPIKeys(ctx context.Context, userID string) ([]auth.APIKey, error) {
	apiKeys, err := s.userRepo.FindAPIKeysForUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "accounting.GetAPIKeys")
	}

	return apiKeys, nil
}

func (s *service) SetAPIKeyName(ctx context.Context, keyID, name string) error {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return errors.New("accounting.SetAPIKeyName: missing user")
	}

	if err := s.userRepo.UpdateAPIKeyName(ctx, user.ID, keyID, name); err != nil {
		return errors.Wrap(err, "accounting.SetAPIKeyName")
	}

	return nil
}

func (s *service) CreateAPIKey(ctx context.Context) (auth.APIKey, error) {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return auth.APIKey{}, errors.New("accounting.InsertAPIKey: missing user")
	}

	key, err := auth.GenerateAPIKey()
	if err != nil {
		return auth.APIKey{}, errors.Wrap(err, "accounting.InsertAPIKey")
	}

	apiKey := auth.APIKey{
		Key:       key,
		Name:      "Untitled",
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.InsertAPIKey(ctx, user.ID, apiKey); err != nil {
		return auth.APIKey{}, errors.Wrap(err, "accounting.InsertAPIKey")
	}

	return apiKey, nil
}

func (s *service) DeleteAPIKey(ctx context.Context, keyID string) error {
	user := auth.UserFromContext(ctx)
	if user == nil {
		return errors.New("accounting.DeleteAPIKey: missing user")
	}

	if err := s.userRepo.DeleteAPIKey(ctx, user.ID, keyID); err != nil {
		return errors.Wrap(err, "accounting.DeleteAPIKey")
	}

	return nil
}
