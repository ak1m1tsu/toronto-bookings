package handlers

import (
	"context"
	"fmt"

	"github.com/romankravchuk/toronto-bookings/internal/storage/models"
)

var (
	mockUsersData = map[string]*models.User{
		"812ae142-945e-496c-b658-7ecbbc4c81f6": {
			ID:                "812ae142-945e-496c-b658-7ecbbc4c81f6",
			Email:             "test-user-1@test.com",
			EncryptedPassword: "$2a$10$2R5Ph/g8updXAhzv0duTs.XHdf8ZmUWgdID.dVsK/30Ds/4jvuDCS",
		},
		"612ae142-945e-496c-b658-7ecbbc4c81f6": {
			ID:                "612ae142-945e-496c-b658-7ecbbc4c81f6",
			Email:             "test-user-2@test.com",
			EncryptedPassword: "$2a$10$2R5Ph/g8updXAhzv0duTs.XHdf8ZmUWgdID.dVsK/30Ds/4jvuDCS",
		},
	}
)

type MockMongoUserStore struct {
	Data map[string]*models.User
}

func NewMockMongoUserStore() *MockMongoUserStore {
	return &MockMongoUserStore{Data: mockUsersData}
}

func (s *MockMongoUserStore) Insert(ctx context.Context, u *models.User) error {
	if _, ok := s.Data[u.ID]; ok {
		return fmt.Errorf("user already exists")
	}
	s.Data[u.ID] = u
	return nil
}

func (s *MockMongoUserStore) GetByID(ctx context.Context, id string) (*models.User, error) {
	if u, ok := s.Data[id]; ok {
		return u, nil
	}
	return nil, fmt.Errorf("user does not exists")
}

func (s *MockMongoUserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, u := range s.Data {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user does not exists")
}

func (s *MockMongoUserStore) GetAll(ctx context.Context) ([]*models.User, error) {
	ret := []*models.User{}
	for _, u := range s.Data {
		ret = append(ret, u)
	}
	return ret, nil
}
