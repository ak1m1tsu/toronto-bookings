package api

import (
	"encoding/json"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/romankravchuk/toronto-bookings/storage"
	"github.com/romankravchuk/toronto-bookings/types"
)

type ReservationHandler struct {
	store storage.ReservationStorage
}

func NewReservationHandler(store storage.ReservationStorage) *ReservationHandler {
	return &ReservationHandler{
		store: store,
	}
}

func (h *ReservationHandler) HandleGetReservationById(writer http.ResponseWriter, request *http.Request) {
	reservationID := chi.URLParam(request, "id")
	reservation, err := h.store.GetByID(request.Context(), reservationID)
	if err != nil {
		JSON(writer, http.StatusNotFound, body{"error": err.Error()})
		return
	}
	resp := &types.ApiResponse{
		Status: http.StatusOK,
		Body:   body{"reservation": reservation},
	}
	JSON(writer, resp.Status, resp)
}

func (h *ReservationHandler) HandleGetReservations(writer http.ResponseWriter, request *http.Request) {
	resp := &types.ApiResponse{Status: http.StatusInternalServerError}

	reservations, err := h.store.GetAll(request.Context())
	if err != nil {
		resp.Body = body{"error": err}
		JSON(writer, resp.Status, resp)
		return
	}

	resp = &types.ApiResponse{
		Status: http.StatusOK,
		Body:   body{"reservations": reservations},
	}
	JSON(writer, resp.Status, resp)
}

func (h *ReservationHandler) HandlePostReservation(writer http.ResponseWriter, request *http.Request) {
	resp := &types.ApiResponse{Status: http.StatusBadRequest}
	reservationReq := &types.CreateReservationRequest{}
	if err := json.NewDecoder(request.Body).Decode(reservationReq); err != nil {
		resp.Body = body{"error": err.Error()}
		JSON(writer, resp.Status, resp)
		return
	}

	reservation, err := types.NewReservationFromRequest(reservationReq)
	if err != nil {
		resp.Body = body{"error": err.Error()}
		JSON(writer, resp.Status, resp)
		return
	}

	if err = h.store.Insert(request.Context(), reservation); err != nil {
		resp.Body = body{"error": err.Error()}
		JSON(writer, resp.Status, resp)
		return
	}
	resp = &types.ApiResponse{
		Status: http.StatusOK,
		Body:   body{"reservation": reservation},
	}
	JSON(writer, resp.Status, resp)
}
