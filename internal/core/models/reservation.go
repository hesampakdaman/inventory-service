package models

import "github.com/google/uuid"

type ReservationState string

const (
	Created   ReservationState = "created"
	Committed ReservationState = "committed"
	Cancelled ReservationState = "cancelled"
)

type ReservationID uuid.UUID

func (r ReservationID) String() string {
	return uuid.UUID(r).String()
}

type Reservation struct {
	ID      ReservationID
	Product Product
	Qty     int
	State   ReservationState
}
