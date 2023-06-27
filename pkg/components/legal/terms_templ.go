// Code generated by templ@v0.2.282 DO NOT EDIT.

package legal

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

// GoExpression
import "github.com/itsnoproblem/prmry/pkg/components"
import "github.com/itsnoproblem/prmry/pkg/components/page"

func TermsPage(view components.Component) templ.Component {
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
			err = TermsOfService().Render(ctx, templBuffer)
			if err != nil {
				return err
			}
			if !templIsBuffer {
				_, err = io.Copy(w, templBuffer)
			}
			return err
		})
		err = page.Page(view).Render(templ.WithChildren(ctx, var_2), templBuffer)
		if err != nil {
			return err
		}
		if !templIsBuffer {
			_, err = io.Copy(w, templBuffer)
		}
		return err
	})
}

func TermsOfService() templ.Component {
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
		_, err = templBuffer.WriteString("<div")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" class=\"container terms-of-service legal-page\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// TemplElement
		err = BackArrow().Render(ctx, templBuffer)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h1>")
		if err != nil {
			return err
		}
		// Text
		var_4 := `Terms of Service`
		_, err = templBuffer.WriteString(var_4)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h1>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_5 := `1. General Agreement`
		_, err = templBuffer.WriteString(var_5)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_6 := `By using this web app, you agree to abide by these Terms of Service. If you do not agree to these terms, you`
		_, err = templBuffer.WriteString(var_6)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_7 := `should not use this service.`
		_, err = templBuffer.WriteString(var_7)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_8 := `2. Account Registration`
		_, err = templBuffer.WriteString(var_8)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_9 := `To access some features of this service, you may need to create an account. You must provide accurate, current,`
		_, err = templBuffer.WriteString(var_9)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_10 := `and complete information during the registration process. You are solely responsible for the activity that`
		_, err = templBuffer.WriteString(var_10)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_11 := `occurs on your account, and you must keep your account password secure.`
		_, err = templBuffer.WriteString(var_11)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_12 := `3. User Responsibilities`
		_, err = templBuffer.WriteString(var_12)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_13 := `You are responsible for all your activities while using this service. You agree to use this service for lawful`
		_, err = templBuffer.WriteString(var_13)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_14 := `purposes only and not engage in any activity that interrupts or harms the service.`
		_, err = templBuffer.WriteString(var_14)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_15 := `4. User-Generated Content`
		_, err = templBuffer.WriteString(var_15)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_16 := `You are solely responsible for the content and prompts you generate. The service provider does not claim`
		_, err = templBuffer.WriteString(var_16)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_17 := `ownership of your content, but you grant us a license to use, reproduce, distribute, and display your content`
		_, err = templBuffer.WriteString(var_17)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_18 := `within the service. You also agree that your content will not infringe the intellectual property rights or`
		_, err = templBuffer.WriteString(var_18)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_19 := `other rights of others.`
		_, err = templBuffer.WriteString(var_19)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_20 := `5. Data Privacy`
		_, err = templBuffer.WriteString(var_20)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_21 := `We value your privacy. Please refer to our Privacy Policy for information on how we`
		_, err = templBuffer.WriteString(var_21)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_22 := `collect, use, and disclose your personal data.`
		_, err = templBuffer.WriteString(var_22)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_23 := `6. Intellectual Property`
		_, err = templBuffer.WriteString(var_23)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_24 := `The service and its original content, features, and functionality are and will remain the exclusive property`
		_, err = templBuffer.WriteString(var_24)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_25 := `of the service provider and its licensors.`
		_, err = templBuffer.WriteString(var_25)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_26 := `7. Modification of Service`
		_, err = templBuffer.WriteString(var_26)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_27 := `We reserve the right to modify or discontinue, temporarily or permanently, the service with or without notice`
		_, err = templBuffer.WriteString(var_27)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_28 := `to you. You agree that we shall not be liable to you or any third party for any modification, suspension, or`
		_, err = templBuffer.WriteString(var_28)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_29 := `discontinuance of the service.`
		_, err = templBuffer.WriteString(var_29)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_30 := `8. Termination`
		_, err = templBuffer.WriteString(var_30)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_31 := `We may, in our sole discretion, terminate or suspend your account with or without notice if you breach any of`
		_, err = templBuffer.WriteString(var_31)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_32 := `the terms of this agreement. Upon termination, you continue to be bound by Sections 3 and 4.`
		_, err = templBuffer.WriteString(var_32)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_33 := `9. Dispute Resolution`
		_, err = templBuffer.WriteString(var_33)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_34 := `Any disputes arising out of or related to these Terms of Service or the use of the service will be governed`
		_, err = templBuffer.WriteString(var_34)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_35 := `by the laws of the jurisdiction in which the service provider is located.`
		_, err = templBuffer.WriteString(var_35)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_36 := `10. Changes to Terms of Service`
		_, err = templBuffer.WriteString(var_36)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_37 := `We reserve the right, in our sole discretion, to make changes to these Terms of Service at any time. Your`
		_, err = templBuffer.WriteString(var_37)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_38 := `continued use of the service after the date any such changes become effective constitutes your acceptance of`
		_, err = templBuffer.WriteString(var_38)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_39 := `the new Terms of Service.`
		_, err = templBuffer.WriteString(var_39)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<p>")
		if err != nil {
			return err
		}
		// Text
		var_40 := `Please review these Terms of Service carefully and ensure you understand them before you start using the`
		_, err = templBuffer.WriteString(var_40)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Text
		var_41 := `service. If you have any questions regarding these Terms of Service, please contact our customer service:`
		_, err = templBuffer.WriteString(var_41)
		if err != nil {
			return err
		}
		// Whitespace (normalised)
		_, err = templBuffer.WriteString(` `)
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<a")
		if err != nil {
			return err
		}
		// Element Attributes
		_, err = templBuffer.WriteString(" href=\"mailto:help@prmry.io\"")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString(">")
		if err != nil {
			return err
		}
		// Text
		var_42 := `help@prmry.io`
		_, err = templBuffer.WriteString(var_42)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</a>")
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</p>")
		if err != nil {
			return err
		}
		// Element (standard)
		_, err = templBuffer.WriteString("<h4>")
		if err != nil {
			return err
		}
		// Text
		var_43 := `Date of Last Revision: June 27, 2023`
		_, err = templBuffer.WriteString(var_43)
		if err != nil {
			return err
		}
		_, err = templBuffer.WriteString("</h4>")
		if err != nil {
			return err
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

