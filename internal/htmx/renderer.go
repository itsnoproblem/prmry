package htmx

import (
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/yosssi/gohtml"

	"github.com/itsnoproblem/prmry/internal/auth"
	"github.com/itsnoproblem/prmry/internal/components"
)

type renderer struct{}

func NewRenderer() *renderer {
	return &renderer{}
}

func (rnd *renderer) Render(w http.ResponseWriter, r *http.Request, cmp components.Component) error {
	ctx := r.Context()
	cmp.SetUser(auth.UserFromContext(ctx))

	if IsHXRequest(ctx) {
		return cmp.GetFragmentTemplate().Render(r.Context(), newHTMLWriter(w))
	}

	return cmp.GetFullTemplate().Render(r.Context(), newHTMLWriter(w))
}

func (rnd *renderer) RenderTemplComponent(w http.ResponseWriter, r *http.Request, fullPage, fragment templ.Component) error {
	if IsHXRequest(r.Context()) {
		return fragment.Render(r.Context(), newHTMLWriter(w))
	}

	return fullPage.Render(r.Context(), newHTMLWriter(w))
}

func (rnd *renderer) RenderError(w http.ResponseWriter, r *http.Request, err error) {
	view := components.NewErrorView(err.Error(), http.StatusInternalServerError)
	ctx := r.Context()

	w.WriteHeader(view.Code)

	if IsHXRequest(ctx) {
		components.Error(view).Render(ctx, newHTMLWriter(w))
		return
	}

	components.ErrorPage(view).Render(ctx, newHTMLWriter(w))
}

func (rnd *renderer) Unauthorized(w http.ResponseWriter, r *http.Request) {
	if IsHXRequest(r.Context()) {
		Redirect(w, "/")
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func newHTMLWriter(w io.Writer) io.Writer {
	htmlWriter := gohtml.NewWriter(w)
	htmlWriter.SetLastElement(">")
	return htmlWriter
}
