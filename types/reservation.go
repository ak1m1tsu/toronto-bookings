package types

import "fmt"

const minReservationFirstNameLen = 3

type Reservation struct {
	ID          string `bson:"_id,omitempty" json:"id"`
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name" `
	Email       string `bson:"email" json:"email"`
	PhoneNumber string `bson:"phone_number" json:"phone_number"`
}

type CreateReservationRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type GetReservationRequest struct {
	ID string `json:"id"`
}

func NewProductFromReq(req *CreateReservationRequest) (*Reservation, error) {
	if err := ValidateCreateReservationRequest(req); err != nil {
		return nil, err
	}
	return &Reservation{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}, nil
}

func ValidateCreateReservationRequest(req *CreateReservationRequest) error {
	if len(req.FirstName) < minReservationFirstNameLen {
		return fmt.Errorf("the first name is to short. min length is %d", minReservationFirstNameLen)
	}

	return nil
}
