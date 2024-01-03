package accounting

import (
	"context"
	"fmt"
	"time"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/components/profile"
)

func formatAccountResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(accountResponse)
	if !ok {
		return nil, fmt.Errorf("formatHomeResponse: failed to parse response")
	}

	cmp := profile.ProfileView{
		Provider:  res.Provider,
		Email:     res.Email,
		Name:      res.Name,
		AvatarURL: res.AvatarURL,
		APIKeys:   make([]profile.APIKeyView, len(res.APIKeys)),
	}

	for i, key := range res.APIKeys {
		cmp.APIKeys[i] = profile.APIKeyView{
			Name:      key.Name,
			Key:       key.Key,
			CreatedAt: key.CreatedAt.Format(time.DateOnly),
		}
	}

	cmp.SetUser(auth.UserFromContext(ctx))

	fullPage := profile.ProfilePage(cmp)
	fragment := profile.Profile(cmp)
	cmp.SetTemplates(fullPage, fragment)

	return &cmp, nil
}
