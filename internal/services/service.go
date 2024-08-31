package services

import (
	"context"
	"fmt"
	"github.com/askaroe/reservationAPI/internal/models"
	"github.com/askaroe/reservationAPI/internal/repository"
)

type ReservationService interface {
	GetReservationsByRoomID(ctx context.Context, roomID int) ([]models.Reservation, error)
	CreateReservation(ctx context.Context, reservation models.Reservation) error
}

type reservationService struct {
	repo repository.ReservationRepository
}

func NewReservationService(repo repository.ReservationRepository) ReservationService {
	return &reservationService{repo: repo}
}

func (s *reservationService) GetReservationsByRoomID(ctx context.Context, roomID int) ([]models.Reservation, error) {
	reservations, err := s.repo.GetByRoomID(ctx, roomID)

	if err != nil {
		return nil, fmt.Errorf("failed to get reservations from repository: %w", err)
	}

	return reservations, nil
}

func (s *reservationService) CreateReservation(ctx context.Context, reservation models.Reservation) error {
	existingReservations, err := s.repo.GetByRoomID(ctx, reservation.RoomID)

	if err != nil {
		return fmt.Errorf("failed to check existing reservations: %w", err)
	}

	for _, res := range existingReservations {
		if (reservation.StartDate >= res.StartDate && reservation.StartDate < res.EndDate) || // Overlaps with existing reservation
			(reservation.EndDate > res.StartDate && reservation.EndDate <= res.EndDate) { // Overlaps with existing reservation
			return fmt.Errorf("room is already booked for the selected dates")
		}
	}

	if err := s.repo.Create(ctx, reservation); err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	return nil
}
