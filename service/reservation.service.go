package service

import (
	"context"

	"github.com/romankravchuk/toronto-bookings/storage"
	"github.com/romankravchuk/toronto-bookings/types"
)

type ReservationService struct {
	store storage.ReservationStorage
}

func NewReservationService(store storage.ReservationStorage) *ReservationService {
	return &ReservationService{store: store}
}

func (s *ReservationService) Insert(ctx context.Context, reservation *types.Reservation) error {
	return s.store.Insert(ctx, reservation)
}
func (s *ReservationService) GetByID(ctx context.Context, id string) (*types.Reservation, error) {
	return s.store.GetByID(ctx, id)
}
func (s *ReservationService) GetAll(ctx context.Context) ([]*types.Reservation, error) {
	return s.store.GetAll(ctx)
}
