package accounting

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
	accountEndpoint := htmx.NewEndpoint(
		makeAccountEndpoint(),
		decodeEmptyRequest,
		formatAccountResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/account", htmx.MakeHandler(accountEndpoint, renderer))
		})
	}
}

func decodeEmptyRequest(_ context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}
