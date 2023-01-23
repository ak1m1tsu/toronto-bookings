package handlers

import (
	"net/http"
)

type APIError struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

var (
	ReservationNotFoundError = APIError{Error: "reservation not found", Status: http.StatusNotFound}
	ServerInternalError      = APIError{Error: "internal server error", Status: http.StatusInternalServerError}
)
