package services

import (
	"context"
	"fmt"
	"github.com/askaroe/reservationAPI/internal/models"
	"github.com/askaroe/reservationAPI/internal/repository"
)

type ReservationService interface {
	GetReservationsByRoomID(ctx context.Context, roomID string) ([]models.Reservation, error)
	CreateReservation(ctx context.Context, reservation models.ReservationDto) (models.Reservation, error)
}

type reservationService struct {
	repo repository.ReservationRepository
}

func NewReservationService(repo repository.ReservationRepository) ReservationService {
	return &reservationService{repo: repo}
}

func (s *reservationService) GetReservationsByRoomID(ctx context.Context, roomID string) ([]models.Reservation, error) {
	reservations, err := s.repo.GetByRoomID(ctx, roomID)

	if err != nil {
		return nil, fmt.Errorf("failed to get reservations from repository: %w", err)
	}

	return reservations, nil
}

func (s *reservationService) CreateReservation(ctx context.Context, reservationDto models.ReservationDto) (models.Reservation, error) {
	existingReservations, err := s.repo.GetByRoomID(ctx, reservationDto.RoomID)
	if err != nil {
		return models.Reservation{}, fmt.Errorf("failed to check existing reservations: %w", err)
	}

	for _, res := range existingReservations {
		if (reservationDto.StartDate.After(res.StartDate) || reservationDto.StartDate.Equal(res.StartDate)) && reservationDto.StartDate.Before(res.EndDate) ||
			(reservationDto.EndDate.After(res.StartDate) && (reservationDto.EndDate.Before(res.EndDate) || reservationDto.EndDate.Equal(res.EndDate))) {
			return models.Reservation{}, fmt.Errorf("room is already booked for the selected dates")
		}
	}

	reservation := models.Reservation{
		RoomID:    reservationDto.RoomID,
		StartDate: reservationDto.StartDate,
		EndDate:   reservationDto.EndDate,
	}

	createdReservation, err := s.repo.Create(ctx, reservation)
	if err != nil {
		return models.Reservation{}, fmt.Errorf("failed to create reservation: %w", err)
	}

	return createdReservation, nil
}
