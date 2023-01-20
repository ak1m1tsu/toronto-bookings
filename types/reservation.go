package types

import "regexp"

const (
	minReservationFirstNameLen int    = 3
	minReservationLastNameLen  int    = 3
	phoneNumberRegexPattern    string = `^(\+7|7|8)?[\s\-]?\(?[489][0-9]{2}\)?[\s\-]?[0-9]{3}[\s\-]?[0-9]{2}[\s\-]?[0-9]{2}$`
	emailRegexPattern          string = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

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

func NewReservationFromRequest(req *CreateReservationRequest) (*Reservation, error) {
	if err := ValidateCreateReservationRequest(req); err != nil {
		return nil, err
	}
	return &Reservation{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: normalizePhoneNumber(req.PhoneNumber),
	}, nil
}

func ValidateCreateReservationRequest(req *CreateReservationRequest) error {
	if len(req.FirstName) < minReservationFirstNameLen {
		return FirstNameValidationError
	}

	if len(req.LastName) < minReservationLastNameLen {
		return LastNameValidationError
	}

	if !isFieldValid(req.Email, emailRegexPattern) {
		return EmailAddressValidationError
	}

	if !isFieldValid(req.PhoneNumber, phoneNumberRegexPattern) {
		return PhoneNumberValidationError
	}

	return nil
}

func isFieldValid(field, pattern string) bool {
	return regexp.MustCompile(pattern).Match([]byte(field))
}

func normalizePhoneNumber(phoneNumber string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(phoneNumber, "")
}
