package handlers

import (
	"encoding/json"
	"github.com/askaroe/reservationAPI/internal/models"
	"github.com/askaroe/reservationAPI/internal/services"
	"github.com/askaroe/reservationAPI/pkg/response"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type ReservationHandler struct {
	s services.ReservationService
}

func NewReservationHandler(service services.ReservationService) *ReservationHandler {
	return &ReservationHandler{s: service}
}

func (h *ReservationHandler) GetReservationsByRoomId(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomID")

	reservations, err := h.s.GetReservationsByRoomID(r.Context(), roomId)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, reservations)
}

func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservationDto models.ReservationDto

	if err := json.NewDecoder(r.Body).Decode(&reservationDto); err != nil {
		response.BadRequest(w, r, err, "invalid request")
		return
	}

	createdReservation, err := h.s.CreateReservation(r.Context(), reservationDto)
	if err != nil {
		response.Conflict(w, r, err)
		return
	}

	response.Created(w, r, createdReservation)
}
