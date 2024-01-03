package accounting

import (
	"context"
	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/pkg/errors"
)

type UserRepo interface {
	FindAPIKeysForUser(ctx context.Context, userID string) ([]auth.APIKey, error)
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
