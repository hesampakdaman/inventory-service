package models

import (
	"fmt"

	"github.com/google/uuid"
)

type ProductID uuid.UUID

func (p ProductID) String() string {
	return uuid.UUID(p).String()
}

type Product struct {
	ID          ProductID
	Available   int
	Reserved    int
	Title       string
	Description string

	reservations map[ReservationID]Reservation
}

func (p Product) Reservations() []Reservation {
	out := make([]Reservation, 0, len(p.reservations))
	for _, r := range p.reservations {
		out = append(out, r)
	}
	return out
}

func (p *Product) Reserve(qty int) (ReservationID, error) {
	if p.Available < qty {
		return ReservationID{}, fmt.Errorf(
			"%w: desired %d, available %d",
			ErrInsufficientStock,
			qty,
			p.Available,
		)
	}

	p.Available -= qty
	p.Reserved += qty

	resID := ReservationID(uuid.New())
	p.reservations[resID] = Reservation{
		ID:    resID,
		Qty:   qty,
		State: Created,
	}

	return resID, nil
}

func (p *Product) Commit(resID ReservationID) error {
	res, ok := p.reservations[resID]
	if !ok {
		return fmt.Errorf("%w: product %s missing reservation", ErrReservationNotFound, p.ID)
	}

	if res.State != Created {
		return fmt.Errorf("%w: reservation is %s", ErrInvalidReservationState, res.State)
	}

	p.Reserved -= res.Qty
	res.State = Committed

	p.reservations[resID] = res

	return nil
}

func (p *Product) Cancel(resID ReservationID) error {
	res, ok := p.reservations[resID]
	if !ok {
		return fmt.Errorf("%w: product %s missing reservation", ErrReservationNotFound, p.ID)
	}

	if res.State != Created {
		return fmt.Errorf("%w: reservation is %s", ErrInvalidReservationState, res.State)
	}

	p.Available += res.Qty
	p.Reserved -= res.Qty
	res.State = Cancelled

	p.reservations[resID] = res

	return nil
}
