package main

import (
	"fmt"
	"github.com/Ferluci/fast-realip"
	"github.com/technically-functional/heartbeat/templates"
	"github.com/valyala/fasthttp"
	"log"
)

func errorPageHandler(ctx *fasthttp.RequestCtx, code int, message string, plaintext bool) {
	ctx.SetStatusCode(code)
	log.Printf("- Returned %v to %s - tried to connect with %s to %s",
		code, realip.FromRequest(ctx), ctx.Method(), ctx.Path())

	if plaintext {
		ctx.Response.Header.Set(fasthttp.HeaderContentType, "text/plain; charset=utf-8")
		_, _ = fmt.Fprintf(ctx, "%s\n", message)
	} else {
		p := &templates.ErrorPage{
			Message: message,
			Path:    ctx.Path(),
			Method:  ctx.Method(),
		}
		templates.WritePageTemplate(ctx, p)
	}
}

func ErrorBadRequest(ctx *fasthttp.RequestCtx) {
	errorPageHandler(ctx, fasthttp.StatusBadRequest, "400 Bad Request", true)
}

func ErrorForbidden(ctx *fasthttp.RequestCtx) {
	errorPageHandler(ctx, fasthttp.StatusForbidden, "403 Forbidden", true)
}

func ErrorNotFound(ctx *fasthttp.RequestCtx) {
	errorPageHandler(ctx, fasthttp.StatusNotFound, "404 Not Found", false)
}

func HandleInternalErr(ctx *fasthttp.RequestCtx, message string, err error) {
	errorPageHandler(ctx, fasthttp.StatusInternalServerError, "500 "+message+": "+err.Error(), true)
}
