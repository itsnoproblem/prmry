package accounting

import (
	"context"
	"fmt"

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
	}
	cmp.SetUser(auth.UserFromContext(ctx))

	fullPage := profile.ProfilePage(cmp)
	fragment := profile.Profile(cmp)
	cmp.SetTemplates(fullPage, fragment)

	return &cmp, nil
}
