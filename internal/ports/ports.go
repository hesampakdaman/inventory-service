package ports

import (
	"context"

	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type Repository interface {
	Get(context.Context, models.ProductID) (models.Product, error)
	GetWithReservation(
		context.Context,
		models.ProductID,
		models.ReservationID,
	) (models.Product, error)
	Save(context.Context, models.Product, models.RequestID) error
	GetReservation(context.Context, models.ReservationID) (models.Reservation, error)
}

type Publisher interface {
	Publish(ctx context.Context, msg any) error
}
