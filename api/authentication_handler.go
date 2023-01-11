package api

import (
	"errors"
	"net/http"
	"os"

	"github.com/anthdm/weavebox"
	"github.com/golang-jwt/jwt/v4"
	"github.com/romankravchuk/toronto-bookings/store"
	"github.com/romankravchuk/toronto-bookings/types"
)

var tokenKeyName = "x-api-token"

type AuthenticationHandler struct {
	store store.UserStorer
}

func NewAuthenticationHandler(store store.UserStorer) *AuthenticationHandler {
	return &AuthenticationHandler{
		store: store,
	}
}

func (h *AuthenticationHandler) HandleSignUp(ctx *weavebox.Context) error {
	creds, err := types.NewCredentialsFromRequestBody(ctx.Request().Body)
	if err != nil {
		return err
	}

	user, err := types.NewUserFromCredentials(creds)
	if err != nil {
		return err
	}

	_, err = h.store.GetByEmail(ctx.Context, user.Email)
	if err == nil {
		return errors.New("user already exists")
	}

	if err = h.store.Insert(ctx.Context, user); err != nil {
		return err
	}

	resp := &types.AuthResponse{
		Status: http.StatusOK,
		Body: map[string]any{
			"message": "user created",
		},
	}
	return ctx.JSON(resp.Status, resp)
}

func (h *AuthenticationHandler) HandleSignIn(ctx *weavebox.Context) error {
	creds, err := types.NewCredentialsFromRequestBody(ctx.Request().Body)
	if err != nil {
		return err
	}

	dbUser, err := h.store.GetByEmail(ctx.Context, creds.Email)
	if err != nil {
		return err
	}

	if !dbUser.ValidatePassword(creds.Password) {
		return ErrUnAuthenticated
	}

	claims := types.NewClaims(dbUser.ID)
	tokenString, err := createJWT(claims)
	if err != nil {
		return err
	}

	ctx.Request().AddCookie(&http.Cookie{
		Name:    tokenKeyName,
		Value:   tokenString,
		Expires: claims.ExpiresAt.Local(),
	})

	resp := &types.AuthResponse{
		Status: http.StatusOK,
		Body:   map[string]any{"message": "authorized"},
	}
	return ctx.JSON(resp.Status, resp)
}

func createJWT(claims *types.Claims) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
