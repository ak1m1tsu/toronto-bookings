package service

import (
	"context"

	"github.com/romankravchuk/toronto-bookings/storage"
	"github.com/romankravchuk/toronto-bookings/types"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	store storage.UserStorage
}

func NewUserService(store storage.UserStorage) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Insert(ctx context.Context, user *types.User) error {
	return s.store.Insert(ctx, user)
}
func (s *UserService) GetByID(ctx context.Context, id string) (*types.User, error) {
	return s.store.GetByID(ctx, id)
}
func (s *UserService) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	return s.store.GetByEmail(ctx, email)
}
func (s *UserService) GetAll(ctx context.Context) ([]*types.User, error) {
	return s.store.GetAll(ctx)
}

func (s *UserService) ValidatePassword(u *types.User, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(pwd)) == nil
}
