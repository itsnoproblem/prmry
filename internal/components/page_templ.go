// Code generated by templ@v0.2.334 DO NOT EDIT.

package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Page(cmp Component) templ.Component {
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
		_, err = templBuffer.WriteString("<html lang=\"en\" data-bs-theme=\"dark\"><head><title>")
		if err != nil {
			return err
		}
		var_2 := `PRMRY - AI Hub`
		_, err = templBuffer.WriteString(var_2)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</title><meta http-equiv=\"Content-Type\" content=\"text/html; charset=utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><link rel=\"icon\" type=\"image/x-png\" href=\"/static/favicon.png\"><link href=\"https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/css/bootstrap.min.css\" rel=\"stylesheet\" integrity=\"sha384-KK94CHFLLe+nY2dmCWGMq91rCGa5gtU4mk92HdvYe+M/SXH301p5ILy+dN9+nJOZ\" crossorigin=\"anonymous\"><link rel=\"preconnect\" href=\"https://fonts.googleapis.com\"><link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin><link href=\"https://fonts.googleapis.com/css2?family=Anonymous+Pro&amp;family=Cormorant+Garamond:wght@300;400&amp;family=Public+Sans:wght@200;300;500&amp;display=swap\" rel=\"stylesheet\"><link rel=\"stylesheet\" href=\"/static/main.css\"><script src=\"/static/main.js\">")
		if err != nil {
			return err
		}
		var_3 := ``
		_, err = templBuffer.WriteString(var_3)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/htmx.org@1.8.5\" integrity=\"sha384-7aHh9lqPYGYZ7sTHvzP1t3BAfLhYSTy9ArHdP3Xsr9/3TlGurYgcPBoFmXX2TX/w\" crossorigin=\"anonymous\">")
		if err != nil {
			return err
		}
		var_4 := ``
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/htmx.org/dist/ext/debug.js\">")
		if err != nil {
			return err
		}
		var_5 := ``
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/htmx.org/dist/ext/disable-element.js\">")
		if err != nil {
			return err
		}
		var_6 := ``
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/idiomorph/dist/idiomorph-ext.min.js\">")
		if err != nil {
			return err
		}
		var_7 := ``
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/htmx.org/dist/ext/json-enc.js\">")
		if err != nil {
			return err
		}
		var_8 := ``
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/htmx.org/dist/ext/remove-me.js\">")
		if err != nil {
			return err
		}
		var_9 := ``
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://unpkg.com/hyperscript.org@0.9.11\">")
		if err != nil {
			return err
		}
		var_10 := ``
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></head><body><div id=\"layout\" class=\"container\">")
		if err != nil {
			return err
		}
		err = TopNavigation(cmp).Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("<div id=\"notices\" class=\"fade-out\"></div><div id=\"content-root\"><div>")
		if err != nil {
			return err
		}
		err = var_1.Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</div></div></div><script src=\"https://kit.fontawesome.com/5a8d7a3cdf.js\" crossorigin=\"anonymous\">")
		if err != nil {
			return err
		}
		var_11 := ``
		_, err = templBuffer.WriteString(var_11)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script><script src=\"https://cdn.jsdelivr.net/npm/bootstrap@5.3.0-alpha3/dist/js/bootstrap.bundle.min.js\" integrity=\"sha384-ENjdO4Dr2bkBIFxQpeoTz1HIcje39Wm4jDKdf19U8gI4ddQ3GYNS7NTKfAdVQSZe\" crossorigin=\"anonymous\">")
		if err != nil {
			return err
		}
		var_12 := ``
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</script></body></html>")
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = templBuffer.WriteTo(w)
		}
		return err
	})
}
