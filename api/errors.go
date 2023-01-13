package api

import (
	"fmt"
	"net/http"

	"github.com/anthdm/weavebox"
)

type APIError struct {
	Message string `json:"error"`
	Status  int    `json:"status"`
}

func (e *APIError) Error() string {
	return e.Message
}

var (
	ReservationNotFoundError = APIError{Message: "reservation not found", Status: http.StatusNotFound}
	ServerInternalError      = APIError{Message: "internal server error", Status: http.StatusInternalServerError}
)

func HandleAPIError(ctx *weavebox.Context, err error) {
	fmt.Println("API server:", err)
	apiErr, ok := err.(*APIError)
	if !ok {
		_ = ctx.JSON(http.StatusBadRequest, APIError{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}
	_ = ctx.JSON(apiErr.Status, apiErr)
}
