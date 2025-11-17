package ports

import (
	"context"

	"github.com/google/uuid"

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
}

type Producer interface {
	Publish(ctx context.Context, key uuid.UUID, msg any) error
}
