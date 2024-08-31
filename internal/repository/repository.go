package repository

import (
	"context"
	"fmt"
	"github.com/askaroe/reservationAPI/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReservationRepository interface {
	GetByRoomID(ctx context.Context, roomID string) ([]models.Reservation, error)
	Create(ctx context.Context, reservation models.Reservation) (models.Reservation, error)
}

type reservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) GetByRoomID(ctx context.Context, roomID string) ([]models.Reservation, error) {
	query := "SELECT id, created_at, updated_at, room_id, start_date, end_date FROM reservations WHERE room_id = $1"

	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, fmt.Errorf("error getting reservations by room_id: %w", err)
	}
	defer rows.Close()

	var reservations []models.Reservation

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(&res.ID, &res.CreatedAt, &res.UpdatedAt, &res.RoomID, &res.StartDate, &res.EndDate)
		if err != nil {
			return nil, fmt.Errorf("error iterating reservation rows: %w", err)
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading rows: %w", err)
	}

	return reservations, nil
}

func (r *reservationRepository) Create(ctx context.Context, reservation models.Reservation) (models.Reservation, error) {
	query := `
		INSERT INTO reservations (room_id, start_date, end_date) 
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	var createdReservation models.Reservation
	err := r.db.QueryRow(ctx, query, reservation.RoomID, reservation.StartDate, reservation.EndDate).Scan(
		&createdReservation.ID,
		&createdReservation.CreatedAt,
		&createdReservation.UpdatedAt,
	)
	if err != nil {
		return models.Reservation{}, fmt.Errorf("error executing query: %w", err)
	}

	createdReservation.RoomID = reservation.RoomID
	createdReservation.StartDate = reservation.StartDate
	createdReservation.EndDate = reservation.EndDate

	return createdReservation, nil
}
