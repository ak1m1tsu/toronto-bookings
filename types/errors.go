package types

import "fmt"

type ValidationError struct {
	Message string `json:"error"`
}

func (e ValidationError) Error() string {
	return e.Message
}

var (
	FirstNameValidationError    = ValidationError{Message: fmt.Sprintf("the first name is to short. min length is %d", minReservationFirstNameLen)}
	LastNameValidationError     = ValidationError{Message: fmt.Sprintf("the last name is to short. min length is %d", minReservationLastNameLen)}
	EmailAddressValidationError = ValidationError{Message: "the email address is not valid"}
	PhoneNumberValidationError  = ValidationError{Message: "the phone number is not valid"}
)
