package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/romankravchuk/toronto-bookings/config"
	"github.com/romankravchuk/toronto-bookings/internal/router/handlers/models"
	"github.com/romankravchuk/toronto-bookings/internal/storage"
)

type key int

const (
	keyUser key = iota
	keyReservation
)

var (
	ErrUnAuthenticated = errors.New("unauthenticated")
	ErrUnAuthorized    = errors.New("unauthorized")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInternalServer  = errors.New("internal server error")
	ErrBadRequest      = errors.New("bad request")
)

type AuthMiddleware struct {
	store storage.UserStorage
}

func NewAuthMiddleware(store storage.UserStorage) *AuthMiddleware {
	return &AuthMiddleware{
		store: store,
	}
}

func (m *AuthMiddleware) JWTRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			resp         *models.ApiResponse = &models.ApiResponse{Status: http.StatusUnauthorized}
			access_token string
		)

		cookie, _ := r.Cookie(AccessTokenHeader)
		authHeader := r.Header.Get("Authorization")
		fields := strings.Fields(authHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else {
			access_token = cookie.Value
		}

		if access_token == "" {
			resp.Body = body{"error": "your are not logged in"}
			JSON(w, resp.Status, resp)
			return
		}

		conf, _ := config.LoadConfig(".")
		sub, err := ValidateToken(access_token, conf.AccessTokenPublicKey)
		if err != nil {
			resp.Body = body{"error": err.Error()}
			JSON(w, resp.Status, resp)
			return
		}

		user, err := m.store.GetByID(r.Context(), fmt.Sprint(sub))
		if err != nil {
			resp.Body = body{"error": "the user belonging to this token no logger exists"}
			JSON(w, resp.Status, resp)
			return
		}
		ctx := context.WithValue(r.Context(), keyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				err, ok := p.(error)
				if !ok {
					err = fmt.Errorf("%+v", p)
				}

				resp := &models.ApiResponse{Status: http.StatusInternalServerError, Body: body{"error": err.Error()}}
				JSON(w, resp.Status, resp)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

type ReservationMiddleware struct {
	store storage.ReservationStorage
}

func NewReservationMiddleware(store storage.ReservationStorage) *ReservationMiddleware {
	return &ReservationMiddleware{
		store: store,
	}
}

func (m *ReservationMiddleware) Context(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reservationID := chi.URLParam(r, "id")
		reservation, err := m.store.GetByID(r.Context(), reservationID)
		if err != nil {
			resp := &models.ApiResponse{
				Status: http.StatusNotFound,
				Body:   body{"error": err},
			}
			JSON(w, resp.Status, resp)
			return
		}
		ctx := context.WithValue(r.Context(), keyReservation, reservation)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
