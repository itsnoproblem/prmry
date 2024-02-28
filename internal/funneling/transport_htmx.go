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
	editFunnelEndpoint := internalhttp.NewHTMXEndpoint(
		makeEditFunnelEndpoint(svc),
		decodeFunnelIDRequest,
		formatFunnelBuilderResponse,
		auth.Required,
	)

	funnelBuilderEndpoint := internalhttp.NewHTMXEndpoint(
		makeFunnelBuilderEndpoint(svc),
		decodeFunnelBuilderRequest,
		formatFunnelBuilderResponse,
		auth.Required,
	)

	saveFunnelEndpoint := internalhttp.NewHTMXEndpoint(
		makeSaveFunnelEndpoint(svc),
		decodeFunnelBuilderRequest,
		formatFunnelBuilderResponse,
		auth.Required,
	)

	listFunnelsEndpoint := internalhttp.NewHTMXEndpoint(
		makeListFunnelsEndpoint(svc),
		decodeEmptyRequest,
		formatListFunnelsResponse,
		auth.Required,
	)

	searchFlowsEndpoint := internalhttp.NewHTMXEndpoint(
		makeSearchFlowsEndpoint(svc),
		decodeSearchFlowsRequest,
		formatSearchFlowsResponse,
		auth.Required,
	)

	addFlowToFunnelEndpoint := internalhttp.NewHTMXEndpoint(
		makeAddFlowToFunnelEndpoint(svc),
		decodeAddFlowToFunnelRequest,
		formatAddFlowToFunnelResponse,
		auth.Required,
	)

	deleteFlowFromFunnelEndpoint := internalhttp.NewHTMXEndpoint(
		makeRemoveFlowFromFunnelEndpoint(svc),
		decodeRemoveFlowFromFunnelRequest,
		formatRemoveFlowFromFunnelResponse,
		auth.Required,
	)

	createFunnelEndpoint := internalhttp.NewHTMXEndpoint(
		makeCreateFunnelEndpoint(svc),
		decodeCreateFunnelRequest,
		formatRedirectResponse,
		auth.Required,
	)

	deleteFunnelEndpoint := internalhttp.NewHTMXEndpoint(
		makeDeleteFunnelEndpoint(svc),
		decodeFunnelIDRequest,
		formatRedirectResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Group(func(r chi.Router) {
			// funnel list, forms
			r.Get("/funnels", htmx.MakeHandler(listFunnelsEndpoint, renderer))
			r.Get("/funnels/new", htmx.MakeHandler(funnelBuilderEndpoint, renderer))
			r.Post("/funnels", htmx.MakeHandler(createFunnelEndpoint, renderer))
			r.Get("/funnels/{id}", htmx.MakeHandler(editFunnelEndpoint, renderer))
			// create / update
			r.Post("/funnels", htmx.MakeHandler(saveFunnelEndpoint, renderer))
			r.Put("/funnels/{id}", htmx.MakeHandler(saveFunnelEndpoint, renderer))
			r.Delete("/funnels/{id}", htmx.MakeHandler(deleteFunnelEndpoint, renderer))
			// funnel flows
			r.Post("/funnels/search-flows", htmx.MakeHandler(searchFlowsEndpoint, renderer))
			r.Post("/funnels/{id}/flows/{flowId}", htmx.MakeHandler(addFlowToFunnelEndpoint, renderer))
			r.Delete("/funnels/{id}/flows/{flowId}", htmx.MakeHandler(deleteFlowFromFunnelEndpoint, renderer))
		})
	}
}

func decodeEmptyRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeCreateFunnelRequest(_ context.Context, request *http.Request) (interface{}, error) {
	var req createFunnelRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "funneling.decodeCreateFunnelRequest")
	}

	return req, nil
}

func decodeFunnelIDRequest(_ context.Context, request *http.Request) (interface{}, error) {
	return funnelIDRequest{
		ID: chi.URLParam(request, "id"),
	}, nil
}

func decodeFunnelBuilderRequest(_ context.Context, request *http.Request) (interface{}, error) {
	var req funnelBuilderRequest

	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		if request.Method != http.MethodGet {
			return req, errors.Wrap(err, "funneling.decodeFunnalBuilderRequest")
		}
	}

	req.ID = chi.URLParam(request, "id")
	return req, nil
}

func decodeSearchFlowsRequest(_ context.Context, request *http.Request) (interface{}, error) {
	var req searchFlowsRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "funneling.decodeSearchFlowsRequest")
	}

	return req, nil
}

func decodeAddFlowToFunnelRequest(_ context.Context, request *http.Request) (interface{}, error) {
	return addFlowToFunnelRequest{
		FunnelID: chi.URLParam(request, "id"),
		FlowID:   chi.URLParam(request, "flowId"),
	}, nil
}

func decodeRemoveFlowFromFunnelRequest(_ context.Context, request *http.Request) (interface{}, error) {
	return removeFlowFromFunnelRequest{
		FunnelID: chi.URLParam(request, "id"),
		FlowID:   chi.URLParam(request, "flowId"),
	}, nil
}
