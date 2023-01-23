package service

import (
	"context"

	hmodels "github.com/romankravchuk/toronto-bookings/internal/router/handlers/models"
	"github.com/romankravchuk/toronto-bookings/internal/storage"
	dbmodels "github.com/romankravchuk/toronto-bookings/internal/storage/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store storage.UserStorage
}

func NewUserService(store storage.UserStorage) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Insert(ctx context.Context, creds *hmodels.Credentials) (*hmodels.UserResponse, error) {
	user, err := dbmodels.NewUser(creds.Email, creds.Password)
	if err != nil {
		return nil, err
	}
	err = s.store.Insert(ctx, user)
	if err != nil {
		return nil, err
	}
	resp := hmodels.NewUserResponse(user.ID, user.Email, user.IsAdmin)
	return resp, nil
}
func (s *UserService) GetByID(ctx context.Context, id string) (*hmodels.UserResponse, error) {
	user, err := s.store.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := hmodels.NewUserResponse(user.ID, user.Email, user.IsAdmin)
	return resp, nil
}
func (s *UserService) GetByEmail(ctx context.Context, email string) (*hmodels.UserResponse, error) {
	user, err := s.store.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	resp := hmodels.NewUserResponse(user.ID, user.Email, user.IsAdmin)
	return resp, nil
}
func (s *UserService) GetAll(ctx context.Context) ([]*hmodels.UserResponse, error) {
	users, err := s.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	var resps []*hmodels.UserResponse
	for _, u := range users {
		resps = append(resps, hmodels.NewUserResponse(u.ID, u.Email, u.IsAdmin))
	}
	return resps, nil
}

func (s *UserService) ValidatePassword(id string, pwd string) bool {
	user, err := s.store.GetByID(context.Background(), id)
	if err != nil {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(pwd)) == nil
}
