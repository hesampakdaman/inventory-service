package models

import "github.com/google/uuid"

type ReservationState string

const (
	Pending   ReservationState = "pending"
	Committed ReservationState = "committed"
	Cancelled ReservationState = "cancelled"
)

type ReservationID uuid.UUID

type Reservation struct {
	ID        ReservationID
	ProductID ProductID
	Qty       int
	State     ReservationState
}
