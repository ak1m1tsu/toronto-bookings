package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/romankravchuk/toronto-bookings/config"
	"github.com/romankravchuk/toronto-bookings/storage"
	"github.com/romankravchuk/toronto-bookings/types"
)

type body map[string]any

var (
	AccessTokenHeader  = "access_token"
	RefreshTokenHeader = "refresh_token"
)

type AuthenticationHandler struct {
	store storage.UserStorage
}

func NewAuthenticationHandler(store storage.UserStorage) *AuthenticationHandler {
	return &AuthenticationHandler{
		store: store,
	}
}

func (h *AuthenticationHandler) HandleSignUp(writer http.ResponseWriter, request *http.Request) {
	var resp *types.ApiResponse
	creds, err := types.NewCredentialsFromRequestBody(request.Body)
	if err != nil {
		resp = &types.ApiResponse{
			Status: http.StatusBadRequest,
			Body:   body{"error": err.Error()},
		}
		JSON(writer, resp.Status, resp)
		return
	}

	user, err := types.NewUserFromCredentials(creds)
	if err != nil {
		resp = &types.ApiResponse{
			Status: http.StatusBadRequest,
			Body:   body{"error": err.Error()},
		}
		JSON(writer, resp.Status, resp)
		return
	}

	if err = h.store.Insert(request.Context(), user); err != nil {
		resp = &types.ApiResponse{
			Status: http.StatusBadRequest,
			Body:   body{"error": err.Error()},
		}
		JSON(writer, resp.Status, resp)
		return
	}

	resp = &types.ApiResponse{
		Status: http.StatusOK,
		Body:   body{"user": user},
	}
	JSON(writer, resp.Status, resp)
}

func (h *AuthenticationHandler) HandleSignIn(writer http.ResponseWriter, request *http.Request) {
	var resp *types.ApiResponse

	creds, err := types.NewCredentialsFromRequestBody(request.Body)
	if err != nil {
		resp = &types.ApiResponse{
			Status: http.StatusBadRequest,
			Body:   body{"error": err.Error()},
		}
		JSON(writer, resp.Status, resp)
		return
	}

	dbUser, err := h.store.GetByEmail(request.Context(), creds.Email)
	if err != nil {
		resp = &types.ApiResponse{
			Status: http.StatusBadRequest,
			Body:   body{"error": err.Error()},
		}
		JSON(writer, resp.Status, resp)
		return
	}

	if !dbUser.ValidatePassword(creds.Password) {
		resp = &types.ApiResponse{
			Status: http.StatusBadRequest,
			Body:   body{"error": "invalid password"},
		}
		JSON(writer, resp.Status, resp)
		return
	}

	conf, _ := config.LoadConfig(".")

	access_token, err := CreateToken(conf.AccessTokenExpiresIn, dbUser.ID, conf.AccessTokenPrivateKey)
	if err != nil {
		resp = &types.ApiResponse{
			Status: http.StatusInternalServerError,
			Body:   body{"error": err.Error()},
		}
		JSON(writer, resp.Status, resp)
		return
	}

	refresh_token, err := CreateToken(conf.RefreshTokenExpiresIn, dbUser.ID, conf.RefreshTokenPrivateKey)
	if err != nil {
		resp = &types.ApiResponse{
			Status: http.StatusInternalServerError,
			Body:   body{"error": err.Error()},
		}
		JSON(writer, resp.Status, resp)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     AccessTokenHeader,
		Value:    access_token,
		MaxAge:   conf.AccessTokenMaxAge * 60,
		HttpOnly: true,
	})
	http.SetCookie(writer, &http.Cookie{
		Name:     RefreshTokenHeader,
		Value:    refresh_token,
		MaxAge:   conf.RefreshTokenMaxAge * 60,
		HttpOnly: true,
	})

	resp = &types.ApiResponse{
		Status: http.StatusOK,
		Body:   body{"access_token": access_token},
	}
	JSON(writer, resp.Status, resp)
}

func (h *AuthenticationHandler) HandleLogout(writer http.ResponseWriter, request *http.Request) {
	http.SetCookie(writer, &http.Cookie{
		Name:    AccessTokenHeader,
		Expires: time.Now().Add(time.Duration(-1)),
	})
	http.SetCookie(writer, &http.Cookie{
		Name:    RefreshTokenHeader,
		Expires: time.Now().Add(time.Duration(-1)),
	})
	resp := &types.ApiResponse{
		Status: http.StatusOK,
		Body:   body{"success": true},
	}
	JSON(writer, resp.Status, resp)
}

func (h *AuthenticationHandler) HandleRefreshToken(writer http.ResponseWriter, request *http.Request) {
	resp := &types.ApiResponse{Status: http.StatusForbidden}

	cookie, err := request.Cookie(RefreshTokenHeader)
	if err != nil {
		resp.Body = body{"error": err.Error()}
		JSON(writer, resp.Status, resp)
	}

	conf, _ := config.LoadConfig(".")

	sub, err := ValidateToken(cookie.Value, conf.RefreshTokenPublicKey)
	if err != nil {
		resp.Body = body{"error": err.Error()}
		JSON(writer, resp.Status, resp)
		return
	}

	user, err := h.store.GetByID(request.Context(), fmt.Sprint(sub))
	if err != nil {
		resp.Body = body{"error": err.Error()}
		JSON(writer, resp.Status, resp)
		return
	}

	access_token, err := CreateToken(conf.AccessTokenExpiresIn, user.ID, conf.AccessTokenPrivateKey)
	if err != nil {
		resp.Body = body{"error": err.Error()}
		JSON(writer, resp.Status, resp)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:    AccessTokenHeader,
		Value:   access_token,
		Expires: time.Now().Add(conf.AccessTokenExpiresIn),
	})

	resp = &types.ApiResponse{
		Status: http.StatusOK,
		Body:   body{"access_token": access_token},
	}
	JSON(writer, resp.Status, resp)
}
