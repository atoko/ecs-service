package controller

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
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
	bearer := r.FormValue("auth")
	ctx := context.WithValue(r.Context(), AuthTokenContextKey, bearer)

	token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	}, jwt.WithValidMethods([]string{"HS512"}))

	if err != nil {
		p.Log.Info.Printf("UNABLE_TO_PARSE_JWT")
		p.Log.Error.Printf(err.Error())
		return
	}

	if token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx = context.WithValue(ctx, AuthPrincipalContextKey, &Principal{
				ProfileId: fmt.Sprint(claims["id"]),
			})
		} else {
			p.Log.Info.Printf("TOKEN_INVALID")
			return
		}
	}

	p.Handler.ServeHTTP(w, r.WithContext(ctx))
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
