package service

import (
	"context"
	"errors"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/events"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func (s Service) Reserve(ctx context.Context, cmd commands.ReserveProduct) error {
	resID := models.NewReservationID(cmd.ProductID, cmd.RequestID)

	reservation, err := s.repo.GetReservation(ctx, resID)
	if err == nil {
		return s.republishReservation(ctx, reservation, cmd.RequestID)
	}
	if !errors.Is(err, models.ErrReservationNotFound) {
		return err
	}

	product, err := s.repo.Get(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	if err := product.Reserve(cmd.Qty, resID); err != nil {
		return err
	}

	if err := s.repo.Save(ctx, product, cmd.RequestID); err != nil {
		return err
	}

	return s.bus.Publish(ctx, events.ReservationCreated{
		ReservationID: resID,
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
	return s.bus.Publish(ctx, events.ReservationCreated{
		ReservationID: r.ID,
		ProductID:     r.ProductID,
		RequestID:     reqID,
		Qty:           r.Qty,
	})
}
