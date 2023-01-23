package models

type Reservation struct {
	ID          string `bson:"_id,omitempty"`
	FirstName   string `bson:"first_name"`
	LastName    string `bson:"last_name"`
	Email       string `bson:"email"`
	PhoneNumber string `bson:"phone_number"`
}

func NewReservation(firstName, lastName, email, phoneNumber string) *Reservation {
	return &Reservation{
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		PhoneNumber: phoneNumber,
	}
}
