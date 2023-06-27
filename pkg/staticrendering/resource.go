package staticrendering

import (
	"github.com/itsnoproblem/prmry/pkg/components"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/itsnoproblem/prmry/pkg/components/legal"
)

type Renderer interface {
	RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error
	RenderError(w http.ResponseWriter, r *http.Request, err error)
}

type Resource struct {
	renderer Renderer
}

func NewResource(renderer Renderer) Resource {
	return Resource{
		renderer: renderer,
	}
}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/terms", rs.Terms)
	r.Get("/privacy", rs.Privacy)
	return r
}

func (rs Resource) Terms(w http.ResponseWriter, r *http.Request) {
	cmp := components.BaseComponent{}
	if err := rs.renderer.RenderTemplComponent(w, r, legal.TermsPage(&cmp), legal.TermsOfService()); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

func (rs Resource) Privacy(w http.ResponseWriter, r *http.Request) {
	cmp := components.BaseComponent{}
	if err := rs.renderer.RenderTemplComponent(w, r, legal.PrivacyPage(&cmp), legal.PrivacyPolicy()); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}
