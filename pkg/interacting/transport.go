package interacting

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/api"
)

func RouteHandler(svc service, renderer *Renderer) func(chi.Router) {
	listInteractionsEndpoint := api.NewEndpoint(
		makeListInteractionsEndpoint(svc),
		decodeEmptyRequest,
		formatInteractionSummaries,
	)

	getInteractionHandler := api.NewEndpoint(
		makeGetInteractionEndpoint(svc),
		decodeGetInteractionRequest,
		formatGetInteractionResponse,
	)

	addRoutes := func(r chi.Router) {
		r.Get("/", api.MakeHandler(listInteractionsEndpoint))
		r.Get("/{id}", api.MakeHandler(getInteractionHandler))
	}

	return func(r chi.Router) {
		r.Route("/interactions", addRoutes)
	}
}

func decodeEmptyRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetInteractionRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	id := chi.URLParam(request, "id")
	return getInteractionRequest{
		ID: id,
	}, nil
}
