package interacting

import (
	"context"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/htmx"
)

type Renderer interface {
	Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
	RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error
	Unauthorized(w http.ResponseWriter, r *http.Request)
}

func RouteHandler(svc interactingService, flowSvc flowService, renderer Renderer) func(chi.Router) {
	chatPromptEndpoint := htmx.NewEndpoint(
		makeChatPromptEndpoint(flowSvc),
		decodeChatPromptRequest,
		formatChatPromptResponse,
		auth.Required,
	)

	createInteractionEndpoint := htmx.NewEndpoint(
		makeCreateInteractionEndpoint(svc),
		decodeCreateInteractionRequest,
		formatGetInteractionResponse,
		auth.Required,
	)

	updateChatConsoleEndpoint := htmx.NewEndpoint(
		makeChatPromptEndpoint(flowSvc),
		decodeChatPromptRequest,
		formatChatPromptResponse,
		auth.Required,
	)

	listInteractionsEndpoint := htmx.NewEndpoint(
		makeListInteractionsEndpoint(svc),
		decodeEmptyRequest,
		formatInteractionSummaries,
		auth.Required,
	)

	getInteractionEndpoint := htmx.NewEndpoint(
		makeGetInteractionEndpoint(svc),
		decodeGetInteractionRequest,
		formatGetInteractionResponse,
		auth.Required,
	)

	return func(r chi.Router) {
		r.Route("/interactions", func(r chi.Router) {
			r.Post("/", htmx.MakeHandler(createInteractionEndpoint, renderer))
			r.Put("/", htmx.MakeHandler(updateChatConsoleEndpoint, renderer))
			r.Get("/", htmx.MakeHandler(listInteractionsEndpoint, renderer))
			r.Get("/chat", htmx.MakeHandler(chatPromptEndpoint, renderer))
			r.Get("/{id}", htmx.MakeHandler(getInteractionEndpoint, renderer))
		})
	}
}

func decodeEmptyRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeChatPromptRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	flowParams := request.PostFormValue("flowParams")
	selectedFlow := request.PostFormValue("flowSelector")
	log.Println("flowParams", flowParams)

	return chatPromptRequest{
		SelectedFlow: selectedFlow,
		//FlowParams:
	}, nil
}

func decodeCreateInteractionRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	inputParams := make(map[string]string)

	paramNamesField := request.PostFormValue("flowParamNames")
	inputParamNames, err := htmx.StringOrSlice(paramNamesField)
	if err != nil {
		return nil, errors.Wrap(err, "decodeCreateInteractionRequest")
	}

	flowParamsField := request.PostFormValue("flowParams")
	flowParams, err := htmx.StringOrSlice(flowParamsField)
	if err != nil {
		return nil, errors.Wrap(err, "decodeCreateInteractionRequest")
	}

	if len(flowParams) != len(inputParamNames) {
		return nil, errors.New("decodeCreateInteractionRequest: invalid number of params")
	}

	for i, param := range inputParamNames {
		inputParams[param] = flowParams[i]
	}

	return createInteractionRequest{
		FlowID: request.PostFormValue("flowSelector"),
		Input:  request.PostFormValue("prompt"),
		Params: inputParams,
	}, nil
}

func decodeGetInteractionRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	id := chi.URLParam(request, "id")
	return getInteractionRequest{
		ID: id,
	}, nil
}
