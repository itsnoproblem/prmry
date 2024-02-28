// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package flow

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"fmt"

	"github.com/itsnoproblem/prmry/internal/flow"
)

func RuleBuilder(view Detail) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"p-4 ps-2\"><div class=\"row\"><div class=\"col text-left\"><label for=\"require-all\" class=\"pt-1 pe-2 larger\">Execute this flow when</label> <select name=\"requireAll\" id=\"require-all\" class=\"rounded p-1\" placeholder=\"Choose...\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if view.RequireAll {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"true\" selected=\"selected\">All rules match</option> <option value=\"false\">Any rule matches</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"true\">All rules match</option> <option value=\"false\" selected=\"selected\">Any rule matches</option>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select> <button id=\"add-rule\" hx-post=\"/flow-builder/rules\" hx-target=\"#content-root\" hx-push-url=\"false\" class=\"btn btn-info btn-sm ms-3\">Add Rule</button></div></div></div><div id=\"rules-container\" class=\"container-lg ps-2\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(view.Rules) == 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"text-left\"><h2 class=\"pb-2 text-body-secondary\">Flow always executes</h2><div><em>Create a rule to add conditions.</em></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		for i, rule := range view.Rules {
			if i > 0 {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<hr class=\"text-secondary pb-3\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <div class=\"row flow-rule fade-in\"><div class=\"col\"><div class=\"form-floating mb-3\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = FieldSelector(view.SupportedFields, fmt.Sprintf("fieldName-%d", i), "fieldName", rule.Field.Source, "Source", "/flow-builder").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if rule.Field.Source == flow.FieldSourceFlow.String() {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"form-floating mb-3\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = FieldSelector(view.AvailableFlowsByID, fmt.Sprintf("selectedFlows-%d", i), "selectedFlows", rule.Field.Value, "Flow", "/flow-builder").Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			if rule.Field.Source == flow.FieldSourceInputArg.String() {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"form-floating mb-3\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = FieldSelector(view.InputParams.Map(), fmt.Sprintf("ruleInputParams-%d", i), "ruleInputParams", rule.Field.Value, "Param name", "/flow-builder").Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"col\"><div class=\"form-floating mb-3\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = FieldSelector(view.SupportedConditions, fmt.Sprintf("condition-%d", i), "condition", rule.Condition, "Condition", "/flow-builder").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div><div class=\"col\"><div class=\"form-floating mb-3\"><input type=\"text\" name=\"value\" id=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("value-%d", i)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" class=\"form-control form-control-md\" placeholder=\"Value\" value=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(rule.Value))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <label for=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("value-%d", i)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">Value</label></div></div><div class=\"col-1 pt-3\"><a hx-delete=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(fmt.Sprintf("/flow-builder/rules/%d", i)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target=\"#content-root\" class=\"button-secondary\"><i class=\"fa fa-close\"></i></a></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
