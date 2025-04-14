package services

import (
	"context"
	"fmt"
	"reservationAPI/internal/models"
	"reservationAPI/internal/repository"
	"reservationAPI/pkg/jsonlog"
)

type ReservationService interface {
	GetReservationsByRoomID(ctx context.Context, roomID string) ([]models.Reservation, error)
	CreateReservation(ctx context.Context, reservation models.ReservationDto) (models.Reservation, error)
}

type reservationService struct {
	repo   repository.ReservationRepository
	logger *jsonlog.Logger
}

func NewReservationService(repo repository.ReservationRepository, logger *jsonlog.Logger) ReservationService {
	return &reservationService{repo: repo, logger: logger}
}

func (s *reservationService) GetReservationsByRoomID(ctx context.Context, roomID string) ([]models.Reservation, error) {
	reservations, err := s.repo.GetByRoomID(ctx, roomID)

	if err != nil {

		s.logger.PrintError(err, nil)

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
			s.logger.PrintError(err, nil)
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
		s.logger.PrintError(err, nil)
		return models.Reservation{}, fmt.Errorf("failed to create reservation: %w", err)
	}

	return createdReservation, nil
}
