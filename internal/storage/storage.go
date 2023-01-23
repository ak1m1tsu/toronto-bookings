package storage

import (
	"context"

	"github.com/romankravchuk/toronto-bookings/internal/storage/models"
)

type ReservationStorage interface {
	Insert(context.Context, *models.Reservation) error
	GetByID(context.Context, string) (*models.Reservation, error)
	GetAll(context.Context) ([]*models.Reservation, error)
}

type UserStorage interface {
	Insert(context.Context, *models.User) error
	GetByID(context.Context, string) (*models.User, error)
	GetByEmail(context.Context, string) (*models.User, error)
	GetAll(context.Context) ([]*models.User, error)
}
