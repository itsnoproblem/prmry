// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package legal

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/itsnoproblem/prmry/internal/components"

func PrivacyPage(view components.Component) templ.Component {
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
		templ_7745c5c3_Var2 := templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
			if !templ_7745c5c3_IsBuffer {
				templ_7745c5c3_Buffer = templ.GetBuffer()
				defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
			}
			templ_7745c5c3_Err = PrivacyPolicy().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			if !templ_7745c5c3_IsBuffer {
				_, templ_7745c5c3_Err = io.Copy(templ_7745c5c3_W, templ_7745c5c3_Buffer)
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = components.Page(view).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func PrivacyPolicy() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container privacy-policy legal-page\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = BackArrow().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h1>Privacy Policy</h1><h4>Effective Date: June 27, 2023</h4><p>This Privacy Policy explains how we collect, use, disclose, and safeguard your information when you visit and use our app, which allows you to generate, catalog, organize AI prompts, and capture the AI responses. Please read this Privacy Policy carefully. If you do not agree with the terms of this privacy policy, please do not access the app.</p><h4>Information We Collect</h4><p>Personal Identification Information: During account creation, we may collect personally identifiable information such as your name, email address, and username.</p><p><b>Usage Information:</b> We collect information related to how you use the web app, such as the types of prompts you generate, catalog, organize, and how you interact with the app.</p><p><b>Device Information:</b> We may collect device information such as device name, operating system, browser type, and IP address when you use our service.</p><h4>How We Use Your Information</h4><p>We may use the information we collect from you for various purposes, including:<ul><li>To provide, operate, and maintain our web app</li><li>To improve, personalize, and expand our web app</li><li>To understand and analyze how you use our web app</li><li>To develop new products, services, features, and functionality</li><li>To communicate with you, either directly or through one of our partners, including for customer service, to provide you with updates and other information relating to the web app</li></ul></p><h4>Sharing of Information</h4><p>We do not sell, trade, or rent your personal identification information to others. We may share generic aggregated demographic information not linked to any personal identification information regarding visitors and users with our business partners, trusted affiliates, and advertisers for the purposes outlined above.</p><h4>Cookies and Tracking</h4><p>Our web app may use 'cookies' to enhance user experience. You may choose to set your web browser to refuse cookies, or to alert you when cookies are being sent. If you choose to reject cookies, some parts of the web app may not function properly.</p><h4>Security of Your Information</h4><p>We use administrative, technical, and physical security measures to help protect your personal information. While we have taken reasonable steps to secure the personal information you provide to us, please be aware that despite our efforts, no security measures are perfect or impenetrable, and no method of data transmission can be guaranteed against any interception or other types of misuse.</p><h4>Children's Privacy</h4><p>Our web app is not intended for use by anyone under the age of 13. If we learn that we have collected personal information from a child under 13 without parental consent, we will delete that information as quickly as possible.</p><h4>Changes to This Privacy Policy</h4><p>We reserve the right to modify this privacy policy at any time. If we make changes to this policy, we will notify you by updating the date of this Privacy Policy and posting it on the web app. We encourage users to frequently check this page for any changes to stay informed about how we are helping to protect the personal information we collect.</p><h4>Contact Us</h4><p>If you have questions or comments about this Privacy Policy, please contact us at help@prmry.io.</p></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
