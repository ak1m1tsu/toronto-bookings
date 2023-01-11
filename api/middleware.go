package api

import (
	"errors"
	"os"

	"github.com/anthdm/weavebox"
	"github.com/golang-jwt/jwt/v4"
	"github.com/romankravchuk/toronto-bookings/types"
)

var (
	ErrUnAuthenticated = errors.New("unauthenticated")
	ErrUnAuthorized    = errors.New("unauthorized")
)

type AdminAuthMiddleware struct{}

func (mw *AdminAuthMiddleware) Authenticate(ctx *weavebox.Context) error {
	cookie, err := ctx.Request().Cookie(tokenKeyName)
	if err != nil {
		return ErrUnAuthenticated
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
		return ErrUnAuthenticated
	}

	return nil
}

func parseJWT(tokenString string, claims *types.Claims) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}
