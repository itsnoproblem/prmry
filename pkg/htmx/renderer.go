package htmx

import (
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/auth"
	"html/template"
	"net/http"
)

type renderer struct {
	tpl *template.Template
}

func NewRenderer(tpl *template.Template) *renderer {
	return &renderer{
		tpl: tpl,
	}
}

func (rnd *renderer) RenderComponent(w http.ResponseWriter, r *http.Request, fullPageTemplate, fragmentTemplate string, cmp Component) error {
	usr := auth.UserFromContext(r.Context())
	if cmp.IsLocked() && usr == nil {
		rnd.Unauthorized(w, r)
	}

	cmp.SetUser(usr)

	if IsHXRequest(r) {
		return rnd.tpl.ExecuteTemplate(w, fragmentTemplate, cmp)
	} else {
		return rnd.tpl.ExecuteTemplate(w, fullPageTemplate, cmp)
	}
}

type ErrorView struct {
	Error string
	BaseComponent
}

func (rnd *renderer) RenderError(w http.ResponseWriter, r *http.Request, err error) {
	cmp := ErrorView{
		Error: err.Error(),
	}
	_ = rnd.RenderComponent(w, r, "page-error", "error", &cmp)
}

func (rnd *renderer) Unauthorized(w http.ResponseWriter, r *http.Request) {
	if IsHXRequest(r) {
		cmp := BaseComponent{
			IsOutOfBand: true,
		}

		_ = rnd.tpl.ExecuteTemplate(w, "top-navigation", cmp)
		//_ = rnd.tpl.ExecuteTemplate(w, "fragment-login", nil)
	} else {
		url := "/"
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}
