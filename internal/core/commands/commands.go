package commands

import (
	"github.com/google/uuid"

	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type AddStock struct {
	ProductID models.ProductID
	RequestID uuid.UUID
	Qty       int
}

type ReserveProduct struct {
	ProductID models.ProductID
	RequestID models.RequestID
	Qty       int
}

type CommitReservation struct {
	ReservationID models.ReservationID
	ProductID     models.ProductID
	RequestID     uuid.UUID
}

type CancelReservation struct {
	ReservationID models.ReservationID
	ProductID     models.ProductID
	RequestID     uuid.UUID
}
