package flowing

import (
	"context"
	"encoding/json"
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
	listFlowsEndpoint := internalhttp.NewJSONEndpoint(
		makeListFlowsEndpoint(svc),
		decodeEmptyRequest,
		formatListFlowsAPIResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Get("/api/flows", api.MakeHandler(listFlowsEndpoint, renderer))
	}
}

type getFlowRequest struct {
	ID string
}

func decodeGetFlowRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return getFlowRequest{
		ID: chi.URLParam(r, "flowID"),
	}, nil
}
