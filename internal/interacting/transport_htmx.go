package interacting

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/htmx"
	internalhttp "github.com/itsnoproblem/prmry/internal/http"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
	RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error
	Unauthorized(w http.ResponseWriter, r *http.Request)
}

func HTMXRouteHandler(svc interactingService, flowSvc flowService, renderer Renderer) func(chi.Router) {
	chatPromptEndpoint := internalhttp.NewHTMXEndpoint(
		makeChatPromptEndpoint(flowSvc),
		decodeChatPromptRequest,
		formatChatPromptResponse,
		auth.Required,
	)

	createInteractionEndpoint := internalhttp.NewHTMXEndpoint(
		makeCreateInteractionEndpoint(svc),
		decodeCreateInteractionRequest,
		formatGetInteractionResponse,
		auth.Required,
	)

	listInteractionsEndpoint := internalhttp.NewHTMXEndpoint(
		makeListInteractionsEndpoint(svc),
		decodeEmptyRequest,
		formatInteractionSummaries,
		auth.Required,
	)

	getInteractionEndpoint := internalhttp.NewHTMXEndpoint(
		makeGetInteractionEndpoint(svc),
		decodeGetInteractionRequest,
		formatGetInteractionResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Route("/interactions", func(r chi.Router) {
			r.Post("/", htmx.MakeHandler(createInteractionEndpoint, renderer))
			r.Get("/", htmx.MakeHandler(listInteractionsEndpoint, renderer))
			r.Get("/chat", htmx.MakeHandler(chatPromptEndpoint, renderer))
			r.Put("/chat", htmx.MakeHandler(chatPromptEndpoint, renderer))
			r.Get("/{id}", htmx.MakeHandler(getInteractionEndpoint, renderer))
		})
	}
}

func decodeEmptyRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}

type chatPromptRequest struct {
	SelectedFlow string      `json:"flowSelector"`
	FlowParams   interface{} `json:"flowParams"`
}

func decodeChatPromptRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	var req chatPromptRequest

	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "decodeFlowBuilderRequest")
	}

	if len(body) > 0 {
		if err = json.Unmarshal(body, &req); err != nil {
			return nil, errors.Wrap(err, "decodeFlowBuilderRequest")
		}
	}

	inputParams := make(map[string]string)

	inputParamNames, err := htmx.StringOrSlice(req.FlowParams)
	if err != nil {
		return nil, errors.Wrap(err, "decodeChatPromptRequest")
	}

	flowParams, err := htmx.StringOrSlice(req.FlowParams)
	if err != nil {
		return nil, errors.Wrap(err, "decodeChatPromptRequest")
	}

	if len(flowParams) != len(inputParamNames) {
		return nil, errors.New("decodeChatPromptRequest: invalid number of params")
	}

	for i, param := range inputParamNames {
		inputParams[param] = flowParams[i]
	}

	return chatPrompt{
		SelectedFlow: req.SelectedFlow,
		InputParams:  inputParams,
	}, nil
}

type createInteractionRequest struct {
	FlowID         string      `json:"flowSelector"`
	Input          string      `json:"prompt"`
	FlowParamNames interface{} `json:"flowParamNames"`
	FlowParams     interface{} `json:"flowParams"`
}

func (c createInteractionRequest) Params() map[string]string {
	params := make(map[string]string)

	paramNames, err := htmx.StringOrSlice(c.FlowParamNames)
	if err != nil {
		return params
	}

	paramValues, err := htmx.StringOrSlice(c.FlowParams)
	if err != nil {
		return params
	}

	if len(paramValues) != len(paramNames) {
		return params
	}

	for i, key := range paramNames {
		params[key] = paramValues[i]
	}

	return params
}

func decodeCreateInteractionRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	inputParams := make(map[string]string)

	var req createInteractionRequest
	body, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "decodeCreateInteractionRequest")
	}

	if err := json.Unmarshal(body, &req); err != nil {
		return nil, errors.Wrap(err, "decodeCreateInteractionRequest")
	}

	paramNames, err := htmx.StringOrSlice(req.FlowParamNames)
	if err != nil {
		return nil, errors.Wrap(err, "decodeCreateInteractionRequest")
	}

	paramValues, err := htmx.StringOrSlice(req.FlowParams)
	if err != nil {
		return nil, errors.Wrap(err, "decodeCreateInteractionRequest")
	}

	if len(paramValues) != len(paramNames) {
		return nil, errors.New("decodeCreateInteractionRequest: invalid number of params")
	}

	for i, key := range paramNames {
		inputParams[key] = paramValues[i]
	}

	return Input{
		FlowID:       req.FlowID,
		InputMessage: req.Input,
		Params:       inputParams,
	}, nil
}

func decodeGetInteractionRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	id := chi.URLParam(request, "id")
	return getInteractionRequest{
		ID: id,
	}, nil
}
