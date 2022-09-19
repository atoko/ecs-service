package controller

import (
	"context"
	"net/http"
)

type GolandContext struct {
	Log *GolandLoggers
}

type GolandContextMiddleware struct {
	Handler http.Handler
	Context *GolandContext
}

var GolandContextKey = "GOLAND_CONTEXT"

func (p *GolandContextMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	golandCtx := context.WithValue(r.Context(), GolandContextKey, p.Context)
	p.Handler.ServeHTTP(w, r.WithContext(golandCtx))
}

func GetGolandContext(ctx context.Context) *GolandContext {
	return ctx.Value(GolandContextKey).(*GolandContext)
}
