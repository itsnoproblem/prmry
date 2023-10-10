// Code generated by templ@v0.2.334 DO NOT EDIT.

package flow

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func TabNav(tabID string, label string, isSelected bool) templ.Component {
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
		_, err = templBuffer.WriteString("<li class=\"nav-item border-bottom\"><a hx-put=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString("/flow-builder/prompt?tab=" + tabID))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\" hx-push-url=\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(templ.EscapeString("?tab=" + tabID))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("\"")
		if err != nil {
			return err
		}
		if isSelected {
			_, err = templBuffer.WriteString(" class=\"nav-link active\" aria-current=\"page\"")
			if err != nil {
				return err
			}
		} else {
			_, err = templBuffer.WriteString(" class=\"nav-link\"")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		var var_2 string = label
		_, err = templBuffer.WriteString(templ.EscapeString(var_2))
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a></li>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}

func TabPanel(isActive bool, cmp templ.Component) templ.Component {
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
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		if !isActive {
			_, err = templBuffer.WriteString(" class=\"d-none\"")
			if err != nil {
				return err
			}
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		err = cmp.Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
