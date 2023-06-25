package components

import (
	"net/http"

	"github.com/a-h/templ"

	errorcmp "github.com/itsnoproblem/prmry/pkg/components/error"
	"github.com/itsnoproblem/prmry/pkg/htmx"
)

type renderer struct{}

func NewRenderer() *renderer {
	return &renderer{}
}

func (rnd *renderer) RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error {
	if htmx.IsHXRequest(r.Context()) {
		return fragment.Render(r.Context(), w)
	}

	return fullPage.Render(r.Context(), w)
}

func (rnd *renderer) RenderError(w http.ResponseWriter, r *http.Request, err error) {
	view := errorcmp.ErrorView{Error: err.Error()}
	//page := ErrorPage(view)
	frag := errorcmp.Error(view)

	frag.Render(r.Context(), w)
}

func (rnd *renderer) Unauthorized(w http.ResponseWriter, r *http.Request) {
	if htmx.IsHXRequest(r.Context()) {
		htmx.Redirect(w, "/")
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
