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
	var usr *auth.User

	if cmp.IsLocked() {
		usr = auth.UserFromContext(r.Context())
		if usr == nil && r != nil && r.URL.Path != "/" {
			rnd.Unauthorized(w, r)
			return nil
		}
		cmp.SetUser(usr)
	}

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

	_ = rnd.tpl.ExecuteTemplate(w, "top-navigation", cmp)
}

func (rnd *renderer) Unauthorized(w http.ResponseWriter, r *http.Request) {
	if IsHXRequest(r) {
		cmp := BaseComponent{
			IsOutOfBand: true,
		}

		_ = rnd.tpl.ExecuteTemplate(w, "top-navigation", cmp)
		//_ = rnd.tpl.ExecuteTemplate(w, "fragment-login", nil)
	} else {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}
