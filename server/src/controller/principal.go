package controller

import (
	"context"
	"net/http"
)

type AuthHeaderParser struct {
	Handler http.Handler
	Log     *GolandLoggers
}

const AuthTokenContextKey = "AUTH_TOKEN"
const AuthPrincipalContextKey = "AUTH_PRINCIPAL"

type Principal struct {
	ProfileId string
}

func (p *AuthHeaderParser) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := s(r.FormValue("auth"))
	ctx := context.WithValue(r.Context(), AuthTokenContextKey, token)

	//TODO: Parse token lol
	if token != nil {
		ctx = context.WithValue(ctx, AuthPrincipalContextKey, &Principal{
			ProfileId: *token,
		})
	}

	p.Handler.ServeHTTP(w, r.WithContext(ctx))
}

func s(s string) *string {
	if s == "" {
		return nil
	} else {
		return &s
	}
}

func AuthPrincipalFromContext(ctx context.Context) *Principal {
	GetGolandContext(ctx).Log.Info.Printf("Checking auth token: %s", ctx.Value(AuthPrincipalContextKey))

	principal := ctx.Value(AuthPrincipalContextKey)

	if principal != nil {
		return principal.(*Principal)
	} else {
		return nil
	}
}
