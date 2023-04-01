package authorizing

import (
	"fmt"
	"github.com/itsnoproblem/mall-fountain-cop-bot/pkg/templates"
	"html/template"
	"io"
)

type Renderer struct {
	templ *template.Template
}

func NewRenderer() (*Renderer, error) {
	tpl, err := templates.Parse()
	if err != nil {
		return nil, fmt.Errorf("authorizing.NewRenderer: %s", err)
	}

	return &Renderer{
		templ: tpl,
	}, nil
}

func (r Renderer) RenderLoginSuccess(w io.Writer, userView UserView) error {
	return r.templ.ExecuteTemplate(w, "page-login-success.gohtml", userView)
}
