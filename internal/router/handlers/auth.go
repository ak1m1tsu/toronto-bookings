package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/romankravchuk/toronto-bookings/internal/config"
	"github.com/romankravchuk/toronto-bookings/internal/router/handlers/models"
	"github.com/romankravchuk/toronto-bookings/internal/service"
)

type body map[string]any

const (
	AccessTokenHeader  string = "access_token"
	RefreshTokenHeader string = "refresh_token"
)

var (
	ErrInvalidPassword error = errors.New("invalid password")
	ErrUserNotExists   error = errors.New("user does not exists")
)

type AuthenticationHandler struct {
	svc          service.UserServicer
	accessToken  config.Token
	refreshToken config.Token
}

func NewAuthenticationHandler(svc service.UserServicer, accessToken, refreshToken config.Token) *AuthenticationHandler {
	return &AuthenticationHandler{svc: svc, accessToken: accessToken, refreshToken: refreshToken}
}

func (h *AuthenticationHandler) HandleSignUp(writer http.ResponseWriter, request *http.Request) {
	resp := models.NewApiResponse(http.StatusBadRequest, body{})

	creds, err := models.NewCredentials(request.Body)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	user, err := h.svc.Insert(request.Context(), creds)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	resp = models.NewApiResponse(
		http.StatusOK,
		body{"user": user},
	)
	JSON(writer, resp.Status, resp)
}

func (h *AuthenticationHandler) HandleSignIn(writer http.ResponseWriter, request *http.Request) {
	resp := models.NewApiResponse(http.StatusBadRequest, body{})

	creds, err := models.NewCredentials(request.Body)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	user, err := h.svc.GetByEmail(request.Context(), creds.Email)
	if err != nil {
		resp.SetError(ErrUserNotExists)
		JSON(writer, resp.Status, resp)
		return
	}

	if !h.svc.ValidatePassword(user.ID, creds.Password) {
		resp.SetError(ErrInvalidPassword)
		JSON(writer, resp.Status, resp)
		return
	}

	resp.SetStatus(http.StatusInternalServerError)

	
	access_token, err := CreateToken(h.accessToken.ExpiresIn, user.ID, h.accessToken.PrivateKey)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	refresh_token, err := CreateToken(h.refreshToken.ExpiresIn, user.ID, h.refreshToken.PrivateKey)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     AccessTokenHeader,
		Value:    access_token,
		MaxAge:   h.accessToken.MaxAge * 60,
		HttpOnly: true,
	})
	http.SetCookie(writer, &http.Cookie{
		Name:     RefreshTokenHeader,
		Value:    refresh_token,
		MaxAge:   h.refreshToken.MaxAge * 60,
		HttpOnly: true,
	})

	resp = models.NewApiResponse(http.StatusOK, body{"access_token": access_token})
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
	resp := models.NewApiResponse(http.StatusOK, body{"success": true})
	JSON(writer, resp.Status, resp)
}

func (h *AuthenticationHandler) HandleRefreshToken(writer http.ResponseWriter, request *http.Request) {
	resp := models.NewApiResponse(http.StatusForbidden, body{})

	cookie, err := request.Cookie(RefreshTokenHeader)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	sub, err := ValidateToken(cookie.Value, h.refreshToken.PublicKey)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	user, err := h.svc.GetByID(request.Context(), fmt.Sprint(sub))
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	access_token, err := CreateToken(h.accessToken.ExpiresIn, user.ID, h.accessToken.PrivateKey)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:    AccessTokenHeader,
		Value:   access_token,
		Expires: time.Now().Add(h.accessToken.ExpiresIn),
	})

	resp = models.NewApiResponse(http.StatusOK, body{"access_token": access_token})
	JSON(writer, resp.Status, resp)
}
