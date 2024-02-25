// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package web

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func PageNotFound() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\" class=\"h-full\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Not Found</title><script src=\"https://unpkg.com/htmx.org@1.9.6\" integrity=\"sha384-FhXw7b6AlE/jyjlZH5iHa/tTe9EpJ1Y55RjcgPbjeWMskSxZt1v9qkxLJWNJaGni\" crossorigin=\"anonymous\"></script><script src=\"https://cdn.tailwindcss.com\"></script></head><body class=\"dark h-full\"><main class=\"grid min-h-full place-items-center bg-white px-6 py-24 sm:py-32 lg:px-8\"><div class=\"text-center\"><p class=\"text-base font-semibold text-indigo-600\">404</p><h1 class=\"mt-4 text-3xl font-bold tracking-tight text-gray-900 sm:text-5xl\">Page not found</h1><p class=\"mt-6 text-base leading-7 text-gray-600\">Sorry, we couldn’t find the page you’re looking for.</p><div class=\"mt-10 flex items-center justify-center gap-x-6\"><a href=\"/\" class=\"rounded-md bg-indigo-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600\">Go back home</a></div></div></main></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}