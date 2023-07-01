package profiling

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components/home"
	"github.com/itsnoproblem/prmry/internal/components/profile"
)

type Renderer interface {
	RenderError(w http.ResponseWriter, r *http.Request, err error)
	RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error
	Unauthorized(w http.ResponseWriter, r *http.Request)
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
	r.Get("/", rs.Home)
	r.Get("/profile", rs.GetProfile)
	return r
}

func (rs Resource) Home(w http.ResponseWriter, r *http.Request) {
	keys, providersByName := auth.Providers()
	cmp := home.HomeView{
		Providers:    keys,
		ProvidersMap: providersByName,
	}

	if r == nil {
		rs.renderer.RenderError(w, r, fmt.Errorf("empty request"))
		return
	}

	if err := cmp.Lock(r); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}

	if err := rs.renderer.RenderTemplComponent(w, r, home.HomePage(cmp), home.HomeFragment(cmp)); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}

func (rs Resource) GetProfile(w http.ResponseWriter, r *http.Request) {
	cmp := profile.ProfileView{}
	if err := cmp.Lock(r); err != nil {
		rs.renderer.RenderError(w, r, err)
	}

	page := profile.ProfilePage(cmp)
	fragment := profile.Profile(cmp)

	if err := rs.renderer.RenderTemplComponent(w, r, page, fragment); err != nil {
		rs.renderer.RenderError(w, r, err)
		return
	}
}
