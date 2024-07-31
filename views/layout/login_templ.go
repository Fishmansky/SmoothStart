// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package layout

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "smoothstart/models"

func LoginPage(data models.LoginData) templ.Component {
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
		templ_7745c5c3_Err = LoginBaseHTML().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-blue px-12 flex justify-start\"><h1 class=\"p-4 text-light text-2xl font-semibold\">SmoothStart</h1></div><div class=\"container mx-auto py-12\"><form class=\"flex flex-col mx-auto w-1/4 space-y-2 p-4 px-4 bg-blue\" action=\"/login\" method=\"POST\" hx-disabled-elt=\"input[type=&#39;text&#39;], input[type=&#39;password&#39;], button\"><label for=\"username\" class=\"mx-auto font-medium text-xl text-light\">Username</label> <input type=\"text\" id=\"username\" name=\"username\" required class=\"bg-mint h-8 font-medium text-xl text-dark text-center\"> <label for=\"password\" class=\"mx-auto font-medium text-xl text-light\">Password</label> <input type=\"password\" id=\"password\" name=\"password\" required class=\"bg-mint h-8 font-medium text-xl text-dark text-center\"><div class=\"flex justify-center\"><button class=\"items-center p-2 bg-cream p-2 text-mywhite font-medium text-xl w-1/2\" type=\"submit\">Login</button></div><div class=\"flex justify-center\"><a class=\"items-center p-2 bg-cream p-2 text-mywhite font-medium text-md text-center w-1/2\" href=\"/recovery\">Forgot password?</a></div></form></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
