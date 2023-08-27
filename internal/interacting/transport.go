package interacting

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

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
	return chatPromptRequest{
		SelectedFlow: request.URL.Query().Get("selectedFlow"),
	}, nil
}

func decodeCreateInteractionRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	return createInteractionRequest{
		FlowID: request.PostFormValue("flowSelector"),
		Input:  request.PostFormValue("prompt"),
	}, nil
}

func decodeGetInteractionRequest(ctx context.Context, request *http.Request) (interface{}, error) {
	id := chi.URLParam(request, "id")
	return getInteractionRequest{
		ID: id,
	}, nil
}
