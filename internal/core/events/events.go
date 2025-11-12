package events

import (
	"github.com/google/uuid"

	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type StockAdded struct {
	ProductID models.ProductID
	RequestID uuid.UUID
	Qty       int
}

type StockAddFailed struct {
	ProductID models.ProductID
	RequestID uuid.UUID
	Reason    string
}

type ReservationCreated struct {
	ReservationID models.ReservationID
	ProductID     models.ProductID
	RequestID     models.RequestID
	Qty           int
}

type ReservationCommitted struct {
	ReservationID models.ReservationID
	ProductID     models.ProductID
	RequestID     uuid.UUID
}

type ReservationCancelled struct {
	ReservationID models.ReservationID
	ProductID     models.ProductID
	RequestID     uuid.UUID
}

type ReservationFailed struct {
	ProductID models.ProductID
	RequestID uuid.UUID
	Reason    string
}
