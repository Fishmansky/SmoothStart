// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package layout

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Header(isAdmin bool) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-mywhite px-12 flex justify-between\"><h1 class=\"p-4 text-dark text-2xl font-semibold\">SmoothStart</h1><div class=\"flex justify-end place-items-center\"><div class=\"p-2 cursor-pointer font-medium text-xl\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isAdmin {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-get=\"/admin/home\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-get=\"/user/home\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-target=\"#app-view\" hx-swap=\"innerHTML\" hx-push-url=\"true\">Home</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isAdmin {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"p-2 cursor-pointer font-medium text-xl\" hx-get=\"/admin/team\" hx-target=\"#app-view\" hx-swap=\"innerHTML\" hx-push-url=\"true\">Team</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"p-2 cursor-pointer font-medium text-xl\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if isAdmin {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-get=\"/admin/plans\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-get=\"/user/plans\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" hx-target=\"#app-view\" hx-swap=\"innerHTML\" hx-push-url=\"true\">Plan</div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}