package api

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

// func HandleAPIError(ctx *weavebox.Context, err error) {
// 	fmt.Println("API server:", err)
// 	apiErr, ok := err.(*APIError)
// 	if !ok {
// 		_ = ctx.JSON(http.StatusBadRequest, APIError{Message: err.Error(), Status: http.StatusInternalServerError})
// 		return
// 	}
// 	_ = ctx.JSON(apiErr.Status, apiErr)
// }
