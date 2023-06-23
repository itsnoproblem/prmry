package interacting

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/itsnoproblem/prmry/pkg/components/chat"
	"github.com/itsnoproblem/prmry/pkg/components/interactions"
	"github.com/itsnoproblem/prmry/pkg/interaction"
)

const (
	paramNameID = "id"
)

type Service interface {
	Interactions(ctx context.Context) ([]interaction.Summary, error)
	Interaction(ctx context.Context, id string) (interaction.Interaction, error)
	NewInteraction(ctx context.Context, msg string) (interaction.Interaction, error)
}

type Renderer interface {
	RenderError(w http.ResponseWriter, r *http.Request, err error)
	RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error
	Unauthorized(w http.ResponseWriter, r *http.Request)
}

type Resource struct {
	renderer Renderer
	service  Service
}

func NewResource(renderer Renderer, svc Service) *Resource {
	return &Resource{
		renderer: renderer,
		service:  svc,
	}
}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()

	// Get all interactions
	r.Get("/", rs.List)

	// Get an interaction by ID
	r.Route(fmt.Sprintf("/{%s}", "id"), func(r chi.Router) {
		//r.Use(WithIDContext)
		r.Get("/", rs.Get)
	})

	// Create an interaction
	r.Post("/", rs.Create)

	// Get the chat prompt
	r.Get("/chat", rs.Chat)

	return r
}

// Create - POST /interactions - send a prompt and receive the prompt + response
func (rs Resource) Create(w http.ResponseWriter, r *http.Request) {
	ixn, err := rs.service.NewInteraction(r.Context(), r.PostFormValue("prompt"))
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	cmp := chat.ChatResponseView{
		Interaction: chat.NewChatDetailView(ixn),
	}
	if err := cmp.Lock(r); err != nil {
		rs.renderer.Unauthorized(w, r)
		return
	}

	templComponent := chat.ChatResponse(cmp)
	if err = rs.renderer.RenderTemplComponent(w, r, templComponent, templComponent); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// Chat - GET /interactions/chat
func (rs Resource) Chat(w http.ResponseWriter, r *http.Request) {
	cmp := chat.ChatControlsView{
		Personas: []chat.PersonaSelector{
			{
				ID:   "123",
				Name: "No Persona",
			},
			{
				ID:   "234",
				Name: "Sarcastic Cop",
			},
			{
				ID:   "345",
				Name: "Concerned Parent",
			},
		},
	}

	if err := cmp.Lock(r); err != nil {
		rs.renderer.Unauthorized(w, r)
		return
	}

	chatFragment := chat.ChatConsole(cmp)
	chatPage := chat.ChatPage(cmp)

	if err := rs.renderer.RenderTemplComponent(w, r, chatPage, chatFragment); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// Get - GET /interactions/{id} - Read a single interaction by :id.
func (rs Resource) Get(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")
	if id == "" {
		rs.renderer.RenderError(w, r, fmt.Errorf("missing required 'id'"))
		return
	}

	ixn, err := rs.service.Interaction(r.Context(), id)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	cmp := chat.NewChatDetailView(ixn)
	if err := cmp.Lock(r); err != nil {
		rs.renderer.Unauthorized(w, r)
		return
	}

	page := interactions.InteractionDetailPage(cmp)
	fragment := interactions.InteractionDetail(cmp)
	if err := rs.renderer.RenderTemplComponent(w, r, page, fragment); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// List - GET /interactions - Read a list of interactions.
func (rs Resource) List(w http.ResponseWriter, r *http.Request) {
	ixns, err := rs.service.Interactions(r.Context())
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	cmp := chat.NewInteractionListView(ixns)
	if err := cmp.Lock(r); err != nil {
		rs.renderer.Unauthorized(w, r)
		return
	}

	fragment := interactions.InteractionsList(cmp)
	page := interactions.InteractionsListPage(cmp)

	if err := rs.renderer.RenderTemplComponent(w, r, page, fragment); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}
