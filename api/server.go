package api

import (
	"encoding/json"
	"github.com/anthdm/weavebox"
	storage "github.com/romankravchuk/toronto-bookings/store"
	"github.com/romankravchuk/toronto-bookings/types"
	"net/http"
)

type ReservationHandler struct {
	store storage.ReservationStorer
}

func NewReservationHandler(store storage.ReservationStorer) *ReservationHandler {
	return &ReservationHandler{
		store: store,
	}
}

func (h *ReservationHandler) HandlePostReservation(ctx *weavebox.Context) error {
	reservationReq := &types.CreateReservationRequest{}
	if err := json.NewDecoder(ctx.Request().Body).Decode(reservationReq); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ServerInternalError)
	}

	reservation, err := types.NewProductFromReq(reservationReq)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, APIError{ErrorMessage: err.Error(), Status: http.StatusBadRequest})
	}

	if err = h.store.Insert(ctx.Context, reservation); err != nil {
		return ctx.JSON(http.StatusInternalServerError, ServerInternalError)
	}

	return ctx.JSON(http.StatusOK, reservation)
}

func (h *ReservationHandler) HandleGetReservationByID(ctx *weavebox.Context) error {
	id := ctx.Param("id")
	reservationReq := &types.GetReservationRequest{ID: id}
	reservation, err := h.store.GetByID(ctx.Context, reservationReq.ID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, ReservationNotFoundError)
	}
	return ctx.JSON(http.StatusOK, reservation)
}
