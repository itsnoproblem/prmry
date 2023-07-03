package components

import (
	"net/http"

	"github.com/a-h/templ"

	"github.com/itsnoproblem/prmry/internal/htmx"
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
	view := ErrorView{Error: err.Error()}
	ctx := r.Context()

	if htmx.IsHXRequest(ctx) {
		Error(view).Render(ctx, w)
		return
	}

	ErrorPage(view).Render(ctx, w)

}

func (rnd *renderer) Unauthorized(w http.ResponseWriter, r *http.Request) {
	if htmx.IsHXRequest(r.Context()) {
		htmx.Redirect(w, "/")
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
