package interacting

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/htmx"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/interaction"
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
	RenderComponent(w http.ResponseWriter, r *http.Request, fullPageTemplate, fragmentTemplate string, cmp htmx.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

type Resource struct {
	renderer Renderer
	service  Service
}

func NewResource(tpl *template.Template, svc Service) *Resource {
	return &Resource{
		renderer: htmx.NewRenderer(tpl),
		service:  svc,
	}
}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()

	// Get all interactions
	r.Get("/", rs.List)

	// Get an interaction by ID
	r.Route(fmt.Sprintf("/{%s}", paramNameID), func(r chi.Router) {
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

	cmp := ChatResponse{
		Interaction: DetailView(ixn),
	}
	cmp.Lock()

	if err = rs.renderer.RenderComponent(w, r, "chat-response.gohtml", "chat-response.gohtml", &cmp); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// Chat - GET /interactions/chat
func (rs Resource) Chat(w http.ResponseWriter, r *http.Request) {
	cmp := ChatControls{
		Personas: []PersonaSelector{
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

	cmp.Lock()
	if err := rs.renderer.RenderComponent(w, r, "page-chat.gohtml", "fragment-chat.gohtml", &cmp); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// Get - GET /interactions/{id} - Read a single interaction by :id.
func (rs Resource) Get(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, paramNameID)
	if id == "" {
		rs.renderer.RenderError(w, r, fmt.Errorf("missing required 'id'"))
		return
	}

	ixn, err := rs.service.Interaction(r.Context(), id)
	if err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
	cmp := DetailView(ixn)
	cmp.Lock()

	if err := rs.renderer.RenderComponent(w, r, "page-interaction.gohtml", "interaction-single.gohtml", &cmp); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

// List - GET /interactions - Read a list of interactions.
func (rs Resource) List(w http.ResponseWriter, r *http.Request) {
	ixns, err := rs.service.Interactions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cmp := ListView(ixns)
	cmp.Lock()
	if err := rs.renderer.RenderComponent(w, r, "page-interactions.gohtml", "interactions-list.gohtml", &cmp); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}
