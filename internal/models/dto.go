package models

import "time"

type ReservationDto struct {
	RoomID    string    `json:"roomId"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
