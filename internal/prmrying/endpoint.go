package prmrying

import (
	"context"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/htmx"
)

type homeResponse struct {
	OAuthProviderKeys    []string
	OAuthProvidersByName map[string]string
}

func makeHomeEndpoint() htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		keys, providersByName := auth.Providers()
		return homeResponse{
			OAuthProviderKeys:    keys,
			OAuthProvidersByName: providersByName,
		}, nil
	}
}

func makeTermsEndpoint() htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

func makePrivacyEndpoint() htmx.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}
