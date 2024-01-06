package accounting

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/htmx"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

func RouteHandler(svc Service, renderer Renderer) func(chi.Router) {
	getAccountEndpoint := internalhttp.NewHTMXEndpoint(
		makeGetAccountEndpoint(svc),
		decodeEmptyRequest,
		formatAccountResponse,
		auth.Required,
	)

	createAPIKeyEndpoint := internalhttp.NewHTMXEndpoint(
		makeCreateAPIKeyEndpoint(svc),
		decodeEmptyRequest,
		formatAPIKeyResponse,
		auth.Required,
	)

	updateAPIKeyEndpoint := internalhttp.NewHTMXEndpoint(
		makeUpdateAPIKeyEndpoint(svc),
		decodeUpdateAPIKeyRequest,
		formatUpdateAPIKeyResponse,
		auth.Required,
	)

	deleteAPIKeyEndpoint := internalhttp.NewHTMXEndpoint(
		makeDeleteAPIKeyEndpoint(svc),
		decodeDeleteAPIKeyRequest,
		formatDeleteAPIKeyResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/account", htmx.MakeHandler(getAccountEndpoint, renderer))
			r.Post("/account/api-keys", htmx.MakeHandler(createAPIKeyEndpoint, renderer))
			r.Put("/account/api-keys/{keyID}", htmx.MakeHandler(updateAPIKeyEndpoint, renderer))
			r.Delete("/account/api-keys/{keyID}", htmx.MakeHandler(deleteAPIKeyEndpoint, renderer))
		})
	}
}

func decodeEmptyRequest(_ context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeUpdateAPIKeyRequest(_ context.Context, request *http.Request) (interface{}, error) {
	return updateAPIKeyRequest{
		KeyID: chi.URLParam(request, "keyID"),
		Name:  request.FormValue("keyName"),
	}, nil
}

func decodeDeleteAPIKeyRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return deleteAPIKeyRequest{
		KeyID: chi.URLParam(request, "keyID"),
	}, nil
}
