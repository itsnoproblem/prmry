package templates

import (
	"embed"
	"github.com/pkg/errors"
	"html/template"
)

var (
	//go:embed "*"
	tpl embed.FS
)

func Parse() (*template.Template, error) {
	tpl, err := template.ParseFS(tpl, "*.gohtml",
		"home/*.gohtml",
		"interactions/*.gohtml",
		"profile/*.gohtml",
	)
	if err != nil {
		return nil, errors.Wrap(err, "templates.Parse")
	}

	return tpl, nil
}
