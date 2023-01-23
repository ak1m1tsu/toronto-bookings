package models

import "regexp"

const (
	minReservationFirstNameLen int    = 3
	minReservationLastNameLen  int    = 3
	phoneNumberRegexPattern    string = `^(\+7|7|8)?[\s\-]?\(?[489][0-9]{2}\)?[\s\-]?[0-9]{3}[\s\-]?[0-9]{2}[\s\-]?[0-9]{2}$`
	emailRegexPattern          string = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

type CreateReservationRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func NewCreateReservationRequest() *CreateReservationRequest {
	return &CreateReservationRequest{}
}

type GetReservationRequest struct {
	ID string `json:"id"`
}

type ReservationResponse struct {
	ID          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func NewReservationResponse(id, firstName, lastName, email, phoneNumber string) *ReservationResponse {
	return &ReservationResponse{id, firstName, lastName, email, phoneNumber}
}

func ValidateRequestData(req *CreateReservationRequest) error {
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

// func NewReservationResponse(req *CreateReservationRequest) (*ReservationResponse, error) {
// 	if err := ValidateRequestData(req); err != nil {
// 		return nil, err
// 	}
// }

func isFieldValid(field, pattern string) bool {
	return regexp.MustCompile(pattern).Match([]byte(field))
}

func normalizePhoneNumber(phoneNumber string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(phoneNumber, "")
}
