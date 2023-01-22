package service

import (
	"context"

	"github.com/romankravchuk/toronto-bookings/types"
)

type UserServicer interface {
	Insert(context.Context, *types.User) error
	GetByID(context.Context, string) (*types.User, error)
	GetByEmail(context.Context, string) (*types.User, error)
	GetAll(context.Context) ([]*types.User, error)
}

type ReservationServicer interface {
	Insert(context.Context, *types.Reservation) error
	GetByID(context.Context, string) (*types.Reservation, error)
	GetAll(context.Context) ([]*types.Reservation, error)
}
