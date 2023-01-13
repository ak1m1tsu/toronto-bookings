package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/anthdm/weavebox"
	"github.com/golang-jwt/jwt/v4"
	"github.com/romankravchuk/toronto-bookings/types"
)

var (
	ErrUnAuthenticated = errors.New("unauthenticated")
	ErrUnAuthorized    = errors.New("unauthorized")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInternalServer  = errors.New("internal server error")
	ErrBadRequest      = errors.New("bad request")
)

type AdminAuthMiddleware struct{}

func (mw *AdminAuthMiddleware) Authenticate(ctx *weavebox.Context) error {
	cookie, err := ctx.Request().Cookie(tokenKeyName)
	if err != nil {
		if err == http.ErrNoCookie {
			return ErrUnAuthenticated
		}
		return ErrBadRequest
	}

	tokenString := cookie.Value
	if len(tokenString) == 0 {
		return ErrUnAuthenticated
	}

	claims := &types.Claims{}
	token, err := parseJWT(tokenString, claims)
	if err != nil {
		return ErrUnAuthenticated
	}
	if !token.Valid {
		return ErrInvalidToken
	}

	if _, ok := ctx.Context.Value(contextClaimsKey).(*types.Claims); !ok {
		ctx.Context = context.WithValue(ctx.Context, contextClaimsKey, claims)
	}
	return nil
}

func parseJWT(tokenString string, claims *types.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
}
