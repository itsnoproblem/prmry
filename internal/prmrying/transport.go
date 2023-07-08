package prmrying

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/htmx"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

func RouteHandler(renderer Renderer) func(chi.Router) {
	homeEndpoint := htmx.NewEndpoint(
		makeHomeEndpoint(),
		decodeEmptyRequest,
		formatHomeResponse,
		auth.NotRequired,
	)

	termsEndpoint := htmx.NewEndpoint(
		makeTermsEndpoint(),
		decodeEmptyRequest,
		formatTermsResponse,
		auth.NotRequired,
	)

	privacyEndpoint := htmx.NewEndpoint(
		makePrivacyEndpoint(),
		decodeEmptyRequest,
		formatPrivacyResponse,
		auth.NotRequired,
	)

	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", htmx.MakeHandler(homeEndpoint, renderer))
			r.Get("/legal/terms", htmx.MakeHandler(termsEndpoint, renderer))
			r.Get("/legal/privacy", htmx.MakeHandler(privacyEndpoint, renderer))
		})
	}
}

func decodeEmptyRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}
