// Code generated by templ@v0.2.282 DO NOT EDIT.

package page

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

// GoExpression
import "github.com/itsnoproblem/prmry/pkg/components"

func TopNavigation(cmp components.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_1 := templ.GetChildren(ctx)
		if var_1 == nil {
			var_1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<nav")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"navbar navbar-expand-lg bg-body-primary mb-3 mt-1\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"container-fluid\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<button")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"navbar-toggler\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" type=\"button\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-toggle=\"collapse\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-target=\"#top-nav\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" aria-controls=\"top-nav\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" aria-expanded=\"false\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" aria-label=\"Toggle navigation\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<span")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"navbar-toggler-icon\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</span>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" id=\"top-nav\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" class=\"collapse navbar-collapse\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-swap-oob=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(components.TrueFalse(cmp.IsOOB())))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// If
		if cmp.User() != nil {
			// TemplElement
			err = UserNavigation(cmp).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		} else {
			// TemplElement
			err = GuestNavigation().Render(ctx, templBuffer)
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</nav>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func GuestNavigation() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_2 := templ.GetChildren(ctx)
		if var_2 == nil {
			var_2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<ul")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"navbar-nav me-auto mb-2 mt-2 mb-lg-0\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<li")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"nav-item\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"nav-link\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" title=\"Documentation\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"#\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-toggle=\"collapse\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-target=\".navbar-collapse.show\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_3 := `Documentation`
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"nav-link\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" title=\"Source Code\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"https://github.com/itsnoproblem/prmry\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_4 := `Source Code`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</ul>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func UserNavigation(cmp components.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_5 := templ.GetChildren(ctx)
		if var_5 == nil {
			var_5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<ul")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"navbar-nav me-auto mb-2 mt-2 mb-lg-0\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<li")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"nav-item dropdown\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" href=\"#\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" class=\"nav-link dropdown-toggle\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" role=\"button\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-toggle=\"dropdown\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" aria-expanded=\"false\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<img")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"img-fluid navbar-avatar\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" alt=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(cmp.User().Name))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" src=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(cmp.User().AvatarURL))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<ul")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"dropdown-menu\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<li>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"dropdown-item\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-get=\"/profile\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-target=\"#content-root\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-push-url=\"true\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-toggle=\"collapse\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-target=\".navbar-collapse.show\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_6 := `Profile`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<li>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"dropdown-item\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=\"#\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-toggle=\"collapse\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-target=\".navbar-collapse.show\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_7 := `Settings`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<li>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<hr")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"dropdown-divider\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<li>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"dropdown-item\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" href=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		var var_8 templ.SafeURL = templ.SafeURL("/auth/logout/"+cmp.User().Provider)
		_, err = templBuffer.WriteString(templ.EscapeString(string(var_8)))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_9 := `Log out`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</ul>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<li")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"nav-item\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"nav-link\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" title=\"Flows\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-get=\"/flows\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-target=\"#content-root\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-push-url=\"true\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-toggle=\"collapse\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-target=\".navbar-collapse.show\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<i")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"fa fa-circle-nodes fa-2x\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</i>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<li")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" nav-class=\"nav-item\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"nav-link\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" title=\"Interactions\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-get=\"/interactions\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-target=\"#content-root\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-push-url=\"true\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-toggle=\"collapse\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-target=\".navbar-collapse.show\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<i")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"fa fa-list fa-2x\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</i>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<li")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"nav-item\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"nav-link\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" title=\"Prompt\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-get=\"/interactions/chat\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-target=\"#content-root\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-push-url=\"true\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-toggle=\"collapse\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" data-bs-target=\".navbar-collapse.show\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<i")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"fa fa-message fa-2x\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" aria-hidden=\"true\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</i>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</li>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</ul>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

