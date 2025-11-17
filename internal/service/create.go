package service

import (
	"context"
	"log/slog"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/events"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func (s Service) Create(ctx context.Context, cmd *commands.CreateProduct) error {
	log := s.logger.With(
		slog.String("product_id", cmd.ProductID.String()),
		slog.String("request_id", cmd.RequestID.String()),
		slog.Int("quantity", cmd.Qty),
	)

	log.DebugContext(ctx, "Creating new product.")

	product := models.Product{
		ID:          cmd.ProductID,
		Available:   cmd.Qty,
		Reserved:    0,
		Title:       cmd.Title,
		Description: cmd.Description,
	}

	if err := s.repo.Save(ctx, product, cmd.RequestID); err != nil {
		log.ErrorContext(ctx, "Failed to persist product.", slog.String("error", err.Error()))
		return err
	}

	if err := s.bus.Publish(ctx, product.ID.UUID(), events.ProductCreated{
		ProductID:   product.ID,
		RequestID:   cmd.RequestID,
		Available:   product.Available,
		Title:       product.Title,
		Description: product.Description,
	}); err != nil {
		log.ErrorContext(ctx, "Failed to publish ProductCreated", slog.String("error", err.Error()))
		return err
	}

	log.DebugContext(ctx, "Product successfully created.")
	return nil
}
