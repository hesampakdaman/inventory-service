package events

import (
	"github.com/google/uuid"

	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type StockAdded struct {
	ProductID models.ProductID `json:"product_id"`
	RequestID uuid.UUID        `json:"request_id"`
	Qty       int              `json:"qty"`
}

type StockAddFailed struct {
	ProductID models.ProductID `json:"product_id"`
	RequestID uuid.UUID        `json:"request_id"`
	Reason    string           `json:"reason"`
}

type ReservationCreated struct {
	ReservationID models.ReservationID `json:"reservation_id"`
	ProductID     models.ProductID     `json:"product_id"`
	RequestID     models.RequestID     `json:"request_id"`
	Qty           int                  `json:"qty"`
}

type ReservationCommitted struct {
	ReservationID models.ReservationID `json:"reservation_id"`
	ProductID     models.ProductID     `json:"product_id"`
	RequestID     uuid.UUID            `json:"request_id"`
}

type ReservationCancelled struct {
	ReservationID models.ReservationID `json:"reservation_id"`
	ProductID     models.ProductID     `json:"product_id"`
	RequestID     uuid.UUID            `json:"request_id"`
}

type ReservationFailed struct {
	ProductID models.ProductID `json:"product_id"`
	RequestID uuid.UUID        `json:"request_id"`
	Reason    string           `json:"reason"`
}
