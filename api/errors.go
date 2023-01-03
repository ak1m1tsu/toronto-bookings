package api

import "net/http"

type APIError struct {
	ErrorMessage string `json:"error_message"`
	Status       int    `json:"status"`
}

func (e *APIError) Error() string {
	return e.ErrorMessage
}

var (
	ReservationNotFoundError = APIError{ErrorMessage: "reservation not found", Status: http.StatusNotFound}
	ServerInternalError      = APIError{ErrorMessage: "internal server error", Status: http.StatusInternalServerError}
)
