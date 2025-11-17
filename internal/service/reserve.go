package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/events"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func (s Service) Reserve(ctx context.Context, cmd *commands.ReserveProduct) error {
	resID := models.NewReservationID(cmd.ProductID, cmd.RequestID)

	log := s.logger.With(
		slog.String("product_id", cmd.ProductID.String()),
		slog.String("reservation_id", resID.String()),
		slog.String("request_id", cmd.RequestID.String()),
		slog.Int("quantity", cmd.Qty),
	)

	log.DebugContext(ctx, "Starting reservation flow.")

	product, err := s.repo.GetWithReservation(ctx, cmd.ProductID, resID)
	if err != nil && !errors.Is(err, models.ErrProductNotFound) {
		log.ErrorContext(
			ctx,
			"Failed to read product/reservation from repository.",
			slog.String("error", err.Error()),
		)
		return err
	}

	if reservation, err := product.Reservation(resID); err == nil {
		log.DebugContext(ctx, "Reservation already exists. Republishing event.")
		return s.republishReservation(ctx, reservation, cmd.RequestID)
	}

	if err := product.Reserve(cmd.Qty, resID); err != nil {
		log.ErrorContext(ctx, "Reservation validation failed.", slog.String("error", err.Error()))
		return err
	}

	if err := s.repo.Save(ctx, product, cmd.RequestID); err != nil {
		log.ErrorContext(ctx, "Failed to persist reservation.", slog.String("error", err.Error()))
		return err
	}

	if err = s.bus.Publish(ctx, product.ID.UUID(), events.ReservationCreated{
		ReservationID: resID,
		ProductID:     product.ID,
		RequestID:     cmd.RequestID,
		Qty:           cmd.Qty,
	}); err != nil {
		log.ErrorContext(
			ctx,
			"Failed to publish ReservationCreated event.",
			slog.String("error", err.Error()),
		)
		return err
	}

	log.DebugContext(ctx, "Reservation successfully completed.")
	return nil
}

func (s Service) republishReservation(
	ctx context.Context,
	r models.Reservation,
	reqID models.RequestID,
) error {
	return s.bus.Publish(ctx, r.ProductID.UUID(), events.ReservationCreated{
		ReservationID: r.ID,
		ProductID:     r.ProductID,
		RequestID:     reqID,
		Qty:           r.Qty,
	})
}
