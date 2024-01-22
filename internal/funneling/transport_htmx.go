package funneling

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
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
	funnelBuilderEndpoint := internalhttp.NewHTMXEndpoint(
		makeFunnelBuilderEndpoint(svc),
		decodeEmptyRequest,
		formatFunnelBuilderResponse,
		auth.Required,
	)

	//createFunnelEndpoint := internalhttp.NewHTMXEndpoint(
	//	makeCreateFunnelEndpoint(svc),
	//	decodeCreateFunnelRequest,
	//	formatCreateFunnelResponse,
	//	auth.Required,
	//)
	//
	//listFunnelsEndpoint := internalhttp.NewHTMXEndpoint(
	//	makeListFunnelsEndpoint(svc),
	//	decodeEmptyRequest,
	//	formatListFunnelsResponse,
	//	auth.Required,
	//)

	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/funnels/edit", htmx.MakeHandler(funnelBuilderEndpoint, renderer))
			r.Get("/funnels/edit/{id}", htmx.MakeHandler(funnelBuilderEndpoint, renderer))
			//r.Post("/funnels", htmx.MakeHandler(createFunnelEndpoint, renderer))
			//r.Get("/funnels", htmx.MakeHandler(listFunnelsEndpoint, renderer))
		})
	}
}

func decodeEmptyRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeCreateFunnelRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	var req CreateFunnelRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "funneling.decodeCreateFunnelRequest")
	}

	return req, nil
}
