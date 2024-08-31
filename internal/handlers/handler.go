package handlers

import (
	"encoding/json"
	"github.com/askaroe/reservationAPI/internal/models"
	"github.com/askaroe/reservationAPI/internal/services"
	"github.com/askaroe/reservationAPI/pkg/response"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type ReservationHandler struct {
	s services.ReservationService
}

func NewReservationHandler(service services.ReservationService) *ReservationHandler {
	return &ReservationHandler{s: service}
}

func (h *ReservationHandler) GetReservationsByRoomId(w http.ResponseWriter, r *http.Request) {
	roomIdStr := chi.URLParam(r, "roomID")
	roomID, err := strconv.Atoi(roomIdStr)

	if err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	reservations, err := h.s.GetReservationsByRoomID(r.Context(), roomID)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	response.OK(w, r, reservations)
}

func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation models.Reservation

	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		response.BadRequest(w, r, err, "invalid request")
	}
}
