package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type Repository interface {
	MarkRequestProcessed(
		ctx context.Context,
		productID uuid.UUID,
		requestID uuid.UUID,
	) (bool, error)

	AddStock(ctx context.Context, cmd commands.AddStock) error

	Reserve(ctx context.Context, cmd commands.ReserveProduct) (models.ReservationID, error)
	Commit(ctx context.Context, cmd commands.CommitReservation) error
	Cancel(ctx context.Context, cmd commands.CancelReservation) error
}
