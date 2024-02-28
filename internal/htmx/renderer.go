package htmx

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/pkg/errors"
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
	debug := r.URL.Query().Get("debug")

	if debug != "" {
		err := writeDebug(w, cmp)
		if err != nil {
			return errors.Wrap(err, "RenderTemplComponent")
		}
		return nil
	}

	if IsHXRequest(ctx) {
		tmpl := cmp.GetFragmentTemplate()
		if tmpl == nil {
			return errors.New("htmx.Render: missing fragment template")
		}

		allErrors := cmp.GetErrors()
		if len(allErrors) > 0 {
			errs := make([]string, len(allErrors))
			for i, err := range allErrors {
				errs[i] = err.Error
				errTempl := components.Error(err)
				errTempl.Render(r.Context(), newHTMLWriter(w))
			}
			w.Header().Set("HX-Trigger", "error")
			w.Header().Set("HX-Trigger-Error", strings.Join(errs, ","))
			return tmpl.Render(r.Context(), newHTMLWriter(w))
		}

		return tmpl.Render(r.Context(), newHTMLWriter(w))
	}

	tmpl := cmp.GetFullTemplate()
	if tmpl == nil {
		return errors.New("htmx.Render: missing full template")
	}

	return tmpl.Render(r.Context(), newHTMLWriter(w))
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

	//statusCode := view.Code
	//statusCode := http.StatusInternalServerError
	//w.WriteHeader(statusCode)

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
	gohtml.Condense = true
	htmlWriter := gohtml.NewWriter(w)
	htmlWriter.SetLastElement(">")
	return htmlWriter
}

func writeDebug(w http.ResponseWriter, cmp components.Component) error {
	response, err := json.Marshal(cmp)
	if err != nil {
		return errors.Wrap(err, "writeDebug")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

	return nil
}
