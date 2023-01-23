package service

import (
	"context"

	"github.com/romankravchuk/toronto-bookings/internal/router/handlers/models"
)

type UserServicer interface {
	Insert(context.Context, *models.Credentials) (*models.UserResponse, error)
	GetByID(context.Context, string) (*models.UserResponse, error)
	GetByEmail(context.Context, string) (*models.UserResponse, error)
	GetAll(context.Context) ([]*models.UserResponse, error)
	ValidatePassword(string, string) bool
}

type ReservationServicer interface {
	Insert(context.Context, *models.CreateReservationRequest) (*models.ReservationResponse, error)
	GetByID(context.Context, string) (*models.ReservationResponse, error)
	GetAll(context.Context) ([]*models.ReservationResponse, error)
}
