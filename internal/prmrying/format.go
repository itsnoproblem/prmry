package prmrying

import (
	"context"
	"fmt"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/components/home"
	"github.com/itsnoproblem/prmry/internal/components/legal"
)

func formatHomeResponse(ctx context.Context, response interface{}) (components.Component, error) {
	res, ok := response.(homeResponse)
	if !ok {
		return nil, fmt.Errorf("formatHomeResponse: failed to parse response")
	}

	cmp := home.HomeView{
		Providers:    res.OAuthProviderKeys,
		ProvidersMap: res.OAuthProvidersByName,
	}
	cmp.SetUser(auth.UserFromContext(ctx))

	fullPage := home.HomePage(cmp)
	fragment := home.HomeFragment(cmp)
	cmp.SetTemplates(fullPage, fragment)

	return &cmp, nil
}

func formatTermsResponse(ctx context.Context, response interface{}) (components.Component, error) {
	cmp := components.BaseComponent{}
	cmp.SetUser(auth.UserFromContext(ctx))
	cmp.SetTemplates(legal.TermsPage(&cmp), legal.TermsOfService())
	return &cmp, nil
}

func formatPrivacyResponse(ctx context.Context, response interface{}) (components.Component, error) {
	cmp := components.BaseComponent{}
	cmp.SetUser(auth.UserFromContext(ctx))
	cmp.SetTemplates(legal.PrivacyPage(&cmp), legal.PrivacyPolicy())
	return &cmp, nil
}
