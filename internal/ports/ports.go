package ports

import (
	"context"

	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type Repository interface {
	GetWithReservation(
		context.Context,
		models.ProductID,
		models.ReservationID,
	) (models.Product, error)
	Save(context.Context, models.RequstID, models.Product) error
}
