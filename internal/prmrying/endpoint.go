package prmrying

import (
	"context"

	"github.com/itsnoproblem/prmry/internal/auth"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type homeResponse struct {
	OAuthProviderKeys    []string
	OAuthProvidersByName map[string]string
}

func makeHomeEndpoint() internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		keys, providersByName := auth.Providers()
		return homeResponse{
			OAuthProviderKeys:    keys,
			OAuthProvidersByName: providersByName,
		}, nil
	}
}

func makeTermsEndpoint() internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

func makePrivacyEndpoint() internalhttp.HandlerFunc {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}
