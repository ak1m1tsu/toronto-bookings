package api

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/anthdm/weavebox"
	"github.com/golang-jwt/jwt/v4"
	"github.com/romankravchuk/toronto-bookings/store"
	"github.com/romankravchuk/toronto-bookings/types"
)

var (
	tokenKeyName     = "x-api-token"
	contextClaimsKey = "claims"
	defaultExpTime   = time.Now().Add(time.Minute * 5)
	jwtSecret        = os.Getenv("JWT_SECRET")
)

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

	resp := &types.ApiResponse{
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
		return ErrBadRequest
	}

	dbUser, err := h.store.GetByEmail(ctx.Context, creds.Email)
	if err != nil {
		return ErrInternalServer
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

	resp := &types.ApiResponse{
		Status: http.StatusOK,
		Body:   map[string]any{"message": "authorized"},
	}
	return ctx.JSON(resp.Status, resp)
}

func (h *AuthenticationHandler) HandleRefreshToken(ctx *weavebox.Context) error {
	claims, ok := ctx.Context.Value(contextClaimsKey).(*types.Claims)
	if !ok {
		log.Printf("%+v\n", ctx.Context.Value(contextClaimsKey))
		return ErrInternalServer
	}

	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return ErrBadRequest
	}

	claims.ExpiresAt = jwt.NewNumericDate(defaultExpTime)
	token, err := createJWT(claims)
	if err != nil {
		log.Println(err)
		return ErrInternalServer
	}

	ctx.Request().AddCookie(&http.Cookie{
		Name:    tokenKeyName,
		Value:   token,
		Expires: claims.ExpiresAt.Local(),
	})

	resp := &types.ApiResponse{
		Status: http.StatusOK,
		Body:   map[string]any{"success": true},
	}
	return ctx.JSON(resp.Status, resp)
}

func (h *AuthenticationHandler) HandleLogout(ctx *weavebox.Context) error {
	ctx.Request().AddCookie(&http.Cookie{
		Name:    tokenKeyName,
		Expires: time.Now(),
	})
	resp := &types.ApiResponse{
		Status: http.StatusOK,
		Body:   map[string]any{"success": true},
	}
	return ctx.JSON(resp.Status, resp)
}

func createJWT(claims *types.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
