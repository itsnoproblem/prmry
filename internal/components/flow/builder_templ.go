// Code generated by templ@v0.2.304 DO NOT EDIT.

package flow

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

// GoExpression
import (
	"fmt"

	"github.com/itsnoproblem/prmry/internal/components"
	"github.com/itsnoproblem/prmry/internal/flow"
)

func FlowBuilderPage(view Detail) templ.Component {
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
		// TemplElement
		var_2 := templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
			templBuffer, templIsBuffer := w.(*bytes.Buffer)
			if !templIsBuffer {
				templBuffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templBuffer)
			}
			// TemplElement
			err = FlowBuilder(view).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = components.Page(&view).Render(templ.WithChildren(ctx, var_2), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func FlowBuilder(view Detail) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_3 := templ.GetChildren(ctx)
		if var_3 == nil {
			var_3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<form")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" id=\"flow-builder\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-post=\"/flows\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-target=\"#content-root\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-push-url=\"false\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-ext=\"json-enc\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<input")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" type=\"hidden\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" name=\"id\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" value=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(view.ID))
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
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"row\"")
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
		_, err = templBuffer.WriteString(" class=\"col col-5\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// TemplElement
		err = FlowOptions(view).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"col col-7\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// TemplElement
		err = RuleBuilder(view).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<hr")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"text-info mb-4 mt-4\"")
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
		_, err = templBuffer.WriteString(" class=\"row\"")
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
		_, err = templBuffer.WriteString(" class=\"col col-1\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<input")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"btn btn-primary\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" type=\"submit\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" value=\"Save\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"col col-1\"")
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
		_, err = templBuffer.WriteString(" class=\"btn btn-secondary\"")
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
		_, err = templBuffer.WriteString(" hx-confirm=\"Abandon changes?\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_4 := `Cancel`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"col col-10\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</form>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func FieldSelector(name string, options SortedMap, selected, label string) templ.Component {
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
		_, err = templBuffer.WriteString("<select")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" name=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(name))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" class=\"form-select form-select-md\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" aria-label=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(label))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-put=\"/flows/new/prompt\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<option>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</option>")
		if err != nil {
			return err
		}
		// For
		for _, value := range options.Keys() {
			// If
			if value == selected {
				// Element (standard)
				_, err = templBuffer.WriteString("<option")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" value=")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(value))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(" selected=\"true\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// StringExpression
				var var_6 string = options[value]
				_, err = templBuffer.WriteString(templ.EscapeString(var_6))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
			} else {
				// Element (standard)
				_, err = templBuffer.WriteString("<option")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" value=")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(templ.EscapeString(value))
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
				// StringExpression
				var var_7 string = options[value]
				_, err = templBuffer.WriteString(templ.EscapeString(var_7))
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</option>")
				if err != nil {
					return err
				}
			}
		}
		_, err = templBuffer.WriteString("</select>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<label>")
		if err != nil {
			return err
		}
		// StringExpression
		var var_8 string = label
		_, err = templBuffer.WriteString(templ.EscapeString(var_8))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func RuleBuilder(view Detail) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_9 := templ.GetChildren(ctx)
		if var_9 == nil {
			var_9 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"d-flex justify-content-between\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"text-info\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_10 := `Rules`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<button")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" id=\"add-rule\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-post=\"/flows/new/rules\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-target=\"#content-root\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-push-url=\"false\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" class=\"btn btn-info btn-sm mb-3\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_11 := `New`
		_, err = templBuffer.WriteString(var_11)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</button>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<hr")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"text-info mb-5 mt-0\"")
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
		_, err = templBuffer.WriteString(" id=\"rules-container\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" class=\"container pb-4\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// If
		if len(view.Rules) == 0 {
			// Element (standard)
			_, err = templBuffer.WriteString("<h2")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"text-body-secondary text-body-secondary\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// Text
			var_12 := `No rules yet`
			_, err = templBuffer.WriteString(var_12)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</h2>")
			if err != nil {
				return err
			}
		}
		// For
		for i, rule := range view.Rules {
			// If
			if i > 0 {
				// Element (void)
				_, err = templBuffer.WriteString("<hr")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"text-secondary pb-3\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
			}
			// Whitespace (normalised)
			_, err = templBuffer.WriteString(` `)
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"row flow-rule fade-in\"")
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
			_, err = templBuffer.WriteString(" class=\"col\"")
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
			_, err = templBuffer.WriteString(" class=\"form-floating mb-3\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// TemplElement
			err = FieldSelector("fieldName", view.SupportedFields, rule.Field.Source, "Field Name").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			// If
			if rule.Field.Source == flow.FieldSourceFlow.String() {
				// Element (standard)
				_, err = templBuffer.WriteString("<div")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"form-floating mb-3\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// TemplElement
				err = FieldSelector("selectedFlows", view.AvailableFlowsByID, rule.Field.Value, "Flow").Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"col\"")
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
			_, err = templBuffer.WriteString(" class=\"form-floating mb-3\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// TemplElement
			err = FieldSelector("condition", view.SupportedConditions, rule.Condition, "Condition").Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"col\"")
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
			_, err = templBuffer.WriteString(" class=\"form-floating mb-3\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// Element (void)
			_, err = templBuffer.WriteString("<input")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" type=\"text\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" name=\"value\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" id=\"value\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" class=\"form-control form-control-md\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" placeholder=\"Value\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" value=")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(rule.Value))
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
			// Element (standard)
			_, err = templBuffer.WriteString("<label")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" for=\"value\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// Text
			var_13 := `Value`
			_, err = templBuffer.WriteString(var_13)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</label>")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"col-1 pt-3\"")
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
			_, err = templBuffer.WriteString(" hx-delete=")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(templ.EscapeString(fmt.Sprintf("/flows/new/rules/%d", i)))
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" hx-target=\"#content-root\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" class=\"button-secondary\"")
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
			_, err = templBuffer.WriteString(" class=\"fa fa-close\"")
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
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func FlowOptions(view Detail) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		templBuffer, templIsBuffer := w.(*bytes.Buffer)
		if !templIsBuffer {
			templBuffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templBuffer)
		}
		ctx = templ.InitializeContext(ctx)
		var_14 := templ.GetChildren(ctx)
		if var_14 == nil {
			var_14 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"form-floating pb-4\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<input")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" name=\"name\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" type=\"text\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" class=\"form-control\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" placeholder=\"Welcome Flow\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" value=")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString(view.Name))
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
		// Element (standard)
		_, err = templBuffer.WriteString("<label")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" for=\"name\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_15 := `Flow Name`
		_, err = templBuffer.WriteString(var_15)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"text-body-secondary pb-2\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_16 := `Trigger this flow when`
		_, err = templBuffer.WriteString(var_16)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"container pb-3\"")
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
		_, err = templBuffer.WriteString(" class=\"form-check\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// If
		if view.RequireAll {
			// Element (void)
			_, err = templBuffer.WriteString("<input")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"form-check-input\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" type=\"radio\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" name=\"requireAll\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" id=\"require-all-true\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" value=\"true\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" checked")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
		} else {
			// Element (void)
			_, err = templBuffer.WriteString("<input")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"form-check-input\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" type=\"radio\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" name=\"requireAll\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" id=\"require-all-true\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" value=\"true\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<label")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"form-check-label\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" for=\"require-all-true\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<em>")
		if err != nil {
			return err
		}
		// Text
		var_17 := `ALL`
		_, err = templBuffer.WriteString(var_17)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</em>")
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_18 := `rules match`
		_, err = templBuffer.WriteString(var_18)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"form-check\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// If
		if view.RequireAll {
			// Element (void)
			_, err = templBuffer.WriteString("<input")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"form-check-input\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" type=\"radio\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" name=\"requireAll\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" id=\"require-all-false\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" value=\"false\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
		} else {
			// Element (void)
			_, err = templBuffer.WriteString("<input")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"form-check-input\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" type=\"radio\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" name=\"requireAll\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" id=\"require-all-false\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" value=\"false\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(" checked")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<label")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"form-check-label\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" for=\"require-all-false\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<em>")
		if err != nil {
			return err
		}
		// Text
		var_19 := `ANY`
		_, err = templBuffer.WriteString(var_19)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</em>")
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_20 := `rule matches`
		_, err = templBuffer.WriteString(var_20)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"mb-3\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<label")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" for=\"prompt\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" class=\"form-label text-body-secondary\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_21 := `Prompt generated by this flow`
		_, err = templBuffer.WriteString(var_21)
		if err != nil {
			return err
		}
		// Element (void)
		_, err = templBuffer.WriteString("<br>")
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_22 := `(indicate replacements with `
		_, err = templBuffer.WriteString(var_22)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<code>")
		if err != nil {
			return err
		}
		// Text
		var_23 := `%s`
		_, err = templBuffer.WriteString(var_23)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</code>")
		if err != nil {
			return err
		}
		// Text
		var_24 := `)`
		_, err = templBuffer.WriteString(var_24)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</label>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<textarea")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"form-control\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" name=\"prompt\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" id=\"prompt\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" rows=\"5\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-put=\"/flows/new/prompt\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(" hx-trigger=\"keyup changed delay:1100ms\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// StringExpression
		var var_25 string = view.Prompt
		_, err = templBuffer.WriteString(templ.EscapeString(var_25))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</textarea>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"mb-3\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// If
		if len(view.PromptArgs) > 0 {
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"mb-2 text-body-secondary\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// Text
			var_26 := `Replacements:`
			_, err = templBuffer.WriteString(var_26)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		}
		// For
		for j, arg := range view.PromptArgs {
			// Element (standard)
			_, err = templBuffer.WriteString("<div")
			if err != nil {
				return err
			}
			// Element Attributes
			_, err = templBuffer.WriteString(" class=\"row mb-3\"")
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
			_, err = templBuffer.WriteString(" class=\"form-floating col\"")
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString(">")
			if err != nil {
				return err
			}
			// TemplElement
			err = FieldSelector("promptArgs", view.SupportedFields, arg.Source.String(), "Replacement "+fmt.Sprint(j)).Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
			// If
			if arg.Source == flow.FieldSourceFlow {
				// Element (standard)
				_, err = templBuffer.WriteString("<div")
				if err != nil {
					return err
				}
				// Element Attributes
				_, err = templBuffer.WriteString(" class=\"form-floating col\"")
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString(">")
				if err != nil {
					return err
				}
				// TemplElement
				err = FieldSelector("promptArgFlows", view.AvailableFlowsByID, arg.Value, "Flow").Render(ctx, templBuffer)
				if err != nil {
					return err
				}
				_, err = templBuffer.WriteString("</div>")
				if err != nil {
					return err
				}
			}
			_, err = templBuffer.WriteString("</div>")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}
