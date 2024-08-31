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

// GetReservationsByRoomId get all reservations of the room.
// @Summary Gets a list of reservations
// @Description Create a new reservation for a room
// @Tags reservations
// @Produce json
// @Param roomID path int true "Room ID"
// @Success 200 {array} models.Reservation
// @Failure 400 {object} response.Object
// @Router /reservations/room/{roomID} [get]
func (h *ReservationHandler) GetReservationsByRoomId(w http.ResponseWriter, r *http.Request) {
	roomId := chi.URLParam(r, "roomID")

	reservations, err := h.s.GetReservationsByRoomID(r.Context(), roomId)
	if err != nil {
		response.BadRequest(w, r, err, roomId)
		return
	}

	response.OK(w, r, reservations)
}

// CreateReservation creates a new reservation.
// @Summary Create a new reservation
// @Description Create a new reservation for a room
// @Tags reservations
// @Accept json
// @Produce json
// @Param reservation body models.ReservationDto true "Reservation details"
// @Success 201 {object} models.Reservation
// @Failure 400 {object} response.Object
// @Failure 409 {object} response.Object
// @Router /reservations [post]
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
