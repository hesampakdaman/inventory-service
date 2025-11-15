package service

import (
	"context"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/events"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func (s Service) Reserve(ctx context.Context, cmd *commands.ReserveProduct) error {
	s.logger.Info("HELLO THERE")
	reservationID := models.NewReservationID(cmd.ProductID, cmd.RequestID)

	product, err := s.repo.GetWithReservation(ctx, cmd.ProductID, reservationID)
	if err != nil {
		return err
	}

	if reservation, err := product.Reservation(reservationID); err == nil {
		return s.republishReservation(ctx, reservation, cmd.RequestID)
	}

	if err := product.Reserve(cmd.Qty, reservationID); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, product, cmd.RequestID); err != nil {
		return err
	}

	return s.bus.Publish(ctx, product.ID.UUID(), "topic", events.ReservationCreated{
		ReservationID: reservationID,
		ProductID:     product.ID,
		RequestID:     cmd.RequestID,
		Qty:           cmd.Qty,
	})
}

func (s Service) republishReservation(
	ctx context.Context,
	r models.Reservation,
	reqID models.RequestID,
) error {
	return s.bus.Publish(ctx, r.ProductID.UUID(), "topic", events.ReservationCreated{
		ReservationID: r.ID,
		ProductID:     r.ProductID,
		RequestID:     reqID,
		Qty:           r.Qty,
	})
}
