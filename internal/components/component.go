package components

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/pkg/errors"

	"github.com/itsnoproblem/prmry/internal/auth"
)

type Component interface {
	User() *auth.User
	Lock(r *http.Request) error
	IsOOB() bool
	SetUser(u *auth.User)
	SetTemplates(full, fragment templ.Component)
	GetFullTemplate() templ.Component
	GetFragmentTemplate() templ.Component
	GetErrors() []ErrorView
}

type BaseComponent struct {
	IsOutOfBand   bool
	templFullPage templ.Component
	templFragment templ.Component
	user          *auth.User
	requiresAuth  bool
	Errors        []ErrorView
}

func (c *BaseComponent) SetUser(u *auth.User) {
	c.user = u
}

func (c *BaseComponent) GetErrors() []ErrorView {
	return c.Errors
}

func (c *BaseComponent) SetTemplates(full, fragment templ.Component) {
	c.templFullPage = full
	c.templFragment = fragment
}

func (c *BaseComponent) GetFullTemplate() templ.Component {
	return c.templFullPage
}

func (c *BaseComponent) GetFragmentTemplate() templ.Component {
	return c.templFragment
}

// Lock checks that an authenticated user exists in the request context
func (c *BaseComponent) Lock(r *http.Request) error {
	usr := auth.UserFromContext(r.Context())
	if usr == nil && r != nil && r.URL.String() != "/" {
		return fmt.Errorf("Unauthorized")
	}

	c.user = usr
	return nil
}

// User returns a pointer to the user accessing the component
func (c *BaseComponent) User() *auth.User {
	return c.user
}

// IsOOB indicates whether the component should be rendered [out of band].
//
// [out of band]: https://htmx.org/attributes/hx-swap-oob/
func (c *BaseComponent) IsOOB() bool {
	return c.IsOutOfBand
}

func (c *BaseComponent) IsAuthenticated() bool {
	return c.user != nil
}

func NewlineToBR(html string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		//html = templ.EscapeString(html)
		html = strings.Replace(html, "\n", "<br/>", -1)
		if _, err := io.WriteString(w, html); err != nil {
			return errors.Wrap(err, "components.NewlineToBR")
		}
		return nil
	})
}

func requestIsMissingAuthentication(usr *auth.User, r *http.Request) bool {
	return usr == nil && r != nil && r.URL.Path != "/"
}
