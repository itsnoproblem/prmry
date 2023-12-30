package accounting

//import (
//	"context"
//	"fmt"
//	"github.com/itsnoproblem/prmry/internal/auth"
//)
//
//type UserRepo interface {
//	FindUserViaOAuth(ctx context.Context, provider, providerUserID string)
//}
//
//type service struct {
//	userRepo UserRepo
//}
//
//func NewService(userRepo UserRepo) *service {
//	return &service{
//		userRepo: userRepo,
//	}
//}
//
//func (s *service) Account(ctx context.Context) (accountResponse, error) {
//	user := auth.UserFromContext(ctx)
//	if user == nil {
//		return accountResponse{}, fmt.Errorf("missing user")
//	}
//
//
//	return accountResponse{
//		Provider:  user.Provider,
//		Email:     user.Email,
//		Name:      user.Name,
//		AvatarURL: user.AvatarURL,
//	}, nil
//}
