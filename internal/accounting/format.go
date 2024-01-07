package accounting

import (
	"context"
	"fmt"
	"time"

	"github.com/a-h/templ"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/components/profile"
	"github.com/itsnoproblem/prmry/internal/components/success"
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
			CreatedAt: key.CreatedAt.Format(time.DateTime),
		}
	}

	cmp.SetUser(auth.UserFromContext(ctx))

	fullPage := profile.ProfilePage(cmp)
	fragment := profile.Profile(cmp)
	cmp.SetTemplates(fullPage, fragment)

	return &cmp, nil
}

func formatAPIKeyResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(profile.APIKeyView)
	if !ok {
		return nil, fmt.Errorf("formatAPIKeyResponse: failed to parse response")
	}

	cmp := profile.APIKeyView{
		Name: res.Name,
		Key:  res.Key,
	}

	fragment := profile.APIKey(cmp)
	cmp.SetTemplates(fragment, fragment)

	return &cmp, nil
}

func formatUpdateAPIKeyResponse(ctx context.Context, _ interface{}) (components.Component, error) {
	cmp := profile.APIKeySuccessView{}
	cmp.SetTemplates(profile.APIKeySuccess(cmp), profile.APIKeySuccess(cmp))
	return &cmp, nil
}

func formatDeleteAPIKeyResponse(_ context.Context, _ interface{}) (components.Component, error) {
	cmp := components.BaseComponent{}
	cmp.SetTemplates(templ.NopComponent, templ.NopComponent)
	return &cmp, nil
}

func formatUpdateProfileResponse(_ context.Context, _ interface{}) (components.Component, error) {
	cmp := success.SuccessView{
		Message: "Profile updated",
	}
	cmp.SetTemplates(success.Success(cmp), success.Success(cmp))
	return &cmp, nil
}
