package funneling

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/itsnoproblem/prmry/internal/api"
	"github.com/itsnoproblem/prmry/internal/auth"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type JSONRenderer interface {
	Render(w http.ResponseWriter, r *http.Request, data json.RawMessage) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

func JSONRouteHandler(svc Service, renderer JSONRenderer) func(chi.Router) {
	executeFunnelEndpoint := internalhttp.NewJSONEndpoint(
		makeExecuteFunnelEndpoint(svc),
		decodeFunnelRequest,
		formatFunnelResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Post("/api/funnels/{path}", api.MakeHandler(executeFunnelEndpoint, renderer))
	}
}

func decodeFunnelRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req executeFunnelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "decodeFunnelRequest: failed to decode request body")
	}

	req.Path = chi.URLParam(r, "path")
	return req, nil
}
