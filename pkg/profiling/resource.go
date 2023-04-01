package profiling

import (
	"github.com/go-chi/chi/v5"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/auth"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/htmx"
	"html/template"
	"net/http"
)

type Renderer interface {
	RenderComponent(
		w http.ResponseWriter, r *http.Request, fullPagetemplate, fragmentTemplate string, cmp htmx.Component) error
}

type Resource struct {
	renderer Renderer
}

func NewResource(tpl *template.Template) (Resource, error) {
	return Resource{
		renderer: htmx.NewRenderer(tpl),
	}, nil
}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", rs.Home)
	r.Get("/profile", rs.GetProfile)
	return r
}

func (rs Resource) Home(w http.ResponseWriter, r *http.Request) {
	keys, providersByName := auth.Providers()

	home := HomeView{
		Providers:    keys,
		ProvidersMap: providersByName,
	}

	if err := rs.renderer.RenderComponent(w, r, "page-home.gohtml", "fragment-home.gohtml", &home); err != nil {
		writeError(w, err)
		return
	}
}

func (rs Resource) GetProfile(w http.ResponseWriter, r *http.Request) {
	cmp := htmx.BaseComponent{}
	cmp.Lock()
	if err := rs.renderer.RenderComponent(w, r, "page-profile.gohtml", "fragment-profile.gohtml", &cmp); err != nil {
		writeError(w, err)
		return
	}
}

func writeError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
