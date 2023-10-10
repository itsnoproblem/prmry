package interacting

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/itsnoproblem/prmry/internal/api"
	"github.com/itsnoproblem/prmry/internal/auth"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
	"github.com/pkg/errors"
	"net/http"
)

type JSONRenderer interface {
	Render(w http.ResponseWriter, r *http.Request, data json.RawMessage) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

func JSONRouteHandler(svc interactingService, jsonRenderer JSONRenderer) func(chi.Router) {
	createInteractionAPIEndpoint := internalhttp.NewJSONEndpoint(
		makeCreateInteractionEndpoint(svc),
		decodeCreateInteractionAPIRequest,
		formatInteractionAPIResponse,
		auth.Required,
	)

	getInteractionAPIEndpoint := internalhttp.NewJSONEndpoint(
		makeGetInteractionEndpoint(svc),
		decodeGetInteractionAPIRequest,
		formatInteractionAPIResponse,
		auth.Required,
	)

	getInteractionsAPIEndpoint := internalhttp.NewJSONEndpoint(
		makeListInteractionsEndpoint(svc),
		decodeEmptyRequest,
		formatInteractionSummariesAPIResponse,
		auth.NotRequired,
	)

	executeFlowAPIEndpoint := internalhttp.NewJSONEndpoint(
		makeExecuteFlowEndpoint(svc),
		decodeExecuteFlowAPIRequest,
		formatExecuteFlowAPIResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Post("/api/interactions", api.MakeHandler(createInteractionAPIEndpoint, jsonRenderer))
		r.Get("/api/interactions", api.MakeHandler(getInteractionsAPIEndpoint, jsonRenderer))
		r.Get("/api/interactions/{id}", api.MakeHandler(getInteractionAPIEndpoint, jsonRenderer))

		r.Post("/api/flows/{flowID}/execute", api.MakeHandler(executeFlowAPIEndpoint, jsonRenderer))
	}
}

// private

type createInteractionAPIRequest struct {
	FlowID     string            `json:"flowID"`
	FlowParams map[string]string `json:"flowParams"`
	Message    string            `json:"message"`
}

func decodeCreateInteractionAPIRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req createInteractionAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "decodeCreateInteractionAPIRequest")
	}

	return Input{
		FlowID:       req.FlowID,
		InputMessage: req.Message,
		Params:       req.FlowParams,
	}, nil
}

type executeFlowAPIRequest struct {
	FlowID     string
	FlowParams map[string]string `json:"params"`
	Message    string            `json:"message"`
}

func decodeExecuteFlowAPIRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req executeFlowAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "decodeExecuteFlowAPIRequest")
	}

	return Input{
		FlowID:       chi.URLParam(r, "flowID"),
		InputMessage: req.Message,
		Params:       req.FlowParams,
	}, nil
}

func decodeGetInteractionAPIRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return getInteractionRequest{
		ID: chi.URLParam(r, "id"),
	}, nil
}
