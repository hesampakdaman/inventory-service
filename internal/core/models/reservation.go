package models

import "github.com/google/uuid"

type ReservationState string

const (
	Created   ReservationState = "created"
	Committed ReservationState = "committed"
	Cancelled ReservationState = "cancelled"
)

type ReservationID uuid.UUID

func NewReservationID(prodID ProductID, reqID RequestID) ReservationID {
	return ReservationID(uuid.NewSHA1(nsInventory, []byte(prodID.String()+reqID.String())))
}

func (r ReservationID) String() string {
	return uuid.UUID(r).String()
}

type Reservation struct {
	ID        ReservationID
	ProductID ProductID
	Qty       int
	State     ReservationState
}
