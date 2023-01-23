package service

import (
	"context"

	hmodels "github.com/romankravchuk/toronto-bookings/internal/router/handlers/models"
	"github.com/romankravchuk/toronto-bookings/internal/storage"
	dbmodels "github.com/romankravchuk/toronto-bookings/internal/storage/models"
)

type ReservationService struct {
	store storage.ReservationStorage
}

func NewReservationService(store storage.ReservationStorage) *ReservationService {
	return &ReservationService{store: store}
}

func (s *ReservationService) Insert(ctx context.Context, req *hmodels.CreateReservationRequest) (*hmodels.ReservationResponse, error) {
	res := dbmodels.NewReservation(req.FirstName, req.LastName, req.Email, req.PhoneNumber)
	err := s.store.Insert(ctx, res)
	if err != nil {
		return nil, err
	}
	resp := hmodels.NewReservationResponse(res.ID, res.FirstName, res.LastName, res.Email, res.PhoneNumber)
	return resp, err
}
func (s *ReservationService) GetByID(ctx context.Context, id string) (*hmodels.ReservationResponse, error) {
	r, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := hmodels.NewReservationResponse(r.ID, r.FirstName, r.LastName, r.Email, r.PhoneNumber)
	return resp, nil
}
func (s *ReservationService) GetAll(ctx context.Context) ([]*hmodels.ReservationResponse, error) {
	reservations, err := s.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var resps []*hmodels.ReservationResponse
	for _, r := range reservations {
		resps = append(resps, hmodels.NewReservationResponse(r.ID, r.FirstName, r.LastName, r.Email, r.PhoneNumber))
	}
	return resps, nil
}
