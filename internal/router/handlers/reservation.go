package handlers

import (
	"encoding/json"
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/romankravchuk/toronto-bookings/internal/router/handlers/models"
	"github.com/romankravchuk/toronto-bookings/internal/service"
)

type ReservationHandler struct {
	svc service.ReservationServicer
}

func NewReservationHandler(svc service.ReservationServicer) *ReservationHandler {
	return &ReservationHandler{svc: svc}
}

func (h *ReservationHandler) HandleGetReservationById(writer http.ResponseWriter, request *http.Request) {
	resp := models.NewApiResponse(http.StatusNotFound, body{})

	reservationID := chi.URLParam(request, "id")
	reservation, err := h.svc.GetByID(request.Context(), reservationID)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	resp = models.NewApiResponse(http.StatusOK, body{"reservation": reservation})
	JSON(writer, resp.Status, resp)
}

func (h *ReservationHandler) HandleGetReservations(writer http.ResponseWriter, request *http.Request) {
	resp := &models.ApiResponse{Status: http.StatusInternalServerError}

	reservations, err := h.svc.GetAll(request.Context())
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	resp = models.NewApiResponse(http.StatusOK, body{"reservations": reservations})
	JSON(writer, resp.Status, resp)
}

func (h *ReservationHandler) HandlePostReservation(writer http.ResponseWriter, request *http.Request) {
	resp := &models.ApiResponse{Status: http.StatusBadRequest}

	reservationReq := &models.CreateReservationRequest{}
	if err := json.NewDecoder(request.Body).Decode(reservationReq); err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	if err := models.ValidateRequestData(reservationReq); err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	reservation, err := h.svc.Insert(request.Context(), reservationReq)
	if err != nil {
		resp.SetError(err)
		JSON(writer, resp.Status, resp)
		return
	}

	resp = models.NewApiResponse(http.StatusOK, body{"reservation": reservation})
	JSON(writer, resp.Status, resp)
}
