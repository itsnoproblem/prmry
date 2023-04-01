package authorizing

import (
	"context"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/user"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UserRepository interface {
	InsertUser(ctx context.Context, usr user.User) error
	DeleteUser(ctx context.Context, id string) error
	SaveUserFromOAuth(ctx context.Context, usr user.User, oauthProvider, providerUserID string) error
	FindUserViaOAuth(ctx context.Context, provider, providerUserID string) (user.User, bool, error)
	FindUserByEmail(ctx context.Context, email string) (user.User, bool, error)
	ExistsUserViaOAuth(ctx context.Context, provider, providerUserID string) (bool, error)
}

type service struct {
	userRepo UserRepository
}

func NewService(userRepo UserRepository) service {
	return service{
		userRepo: userRepo,
	}
}

func (s service) CreateUser(ctx context.Context, usr user.User) (id string, err error) {
	usr.ID = uuid.NewString()
	if err = s.userRepo.InsertUser(ctx, usr); err != nil {
		return "", errors.Wrap(err, "authService.CreateUser")
	}

	return usr.ID, nil
}

func (s service) DeleteUser(ctx context.Context, id string) error {
	if err := s.userRepo.DeleteUser(ctx, id); err != nil {
		return errors.Wrap(err, "authService.DeleteUser")
	}

	return nil
}

func (s service) UserExistsForOAuthProvider(ctx context.Context, provider, providerUserID string) (bool, error) {
	exists, err := s.userRepo.ExistsUserViaOAuth(ctx, provider, providerUserID)
	if err != nil {
		return false, errors.Wrap(err, "authService.UserExistsForOAuthProvider")
	}

	return exists, nil
}

func (s service) SaveUserWithOAuthConnection(ctx context.Context, usr user.User, provider, providerUserID string) error {
	if err := s.userRepo.SaveUserFromOAuth(ctx, usr, provider, providerUserID); err != nil {
		return errors.Wrap(err, "authService.SaveUserWithOAuthConnection")
	}

	return nil
}

func (s service) GetUserByProvider(ctx context.Context, provider, providerUserID string) (usr user.User, exists bool, err error) {
	usr, exists, err = s.userRepo.FindUserViaOAuth(ctx, provider, providerUserID)
	if err != nil {
		return user.User{}, false, errors.Wrap(err, "authService.GetUserByProvider")
	}

	return usr, exists, nil
}

func (s service) GetUserByEmail(ctx context.Context, email string) (usr user.User, exists bool, err error) {
	usr, exists, err = s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return user.User{}, false, errors.Wrap(err, "authService.GetUserByEmail")
	}

	return usr, exists, nil
}
