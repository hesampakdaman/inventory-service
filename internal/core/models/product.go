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

	Res map[ReservationID]Reservation
}

func (p Product) Reservations() []Reservation {
	out := make([]Reservation, 0, len(p.Res))
	for _, r := range p.Res {
		out = append(out, r)
	}
	return out
}

func (p Product) Reservation(id ReservationID) (Reservation, error) {
	res, ok := p.Res[id]
	if !ok {
		return Reservation{}, ErrReservationNotFound
	}

	return res, nil
}

func (p *Product) Reserve(qty int, resID ReservationID) error {
	if p.Available < qty {
		return fmt.Errorf(
			"%w: desired %d, available %d",
			ErrInsufficientStock,
			qty,
			p.Available,
		)
	}

	p.Available -= qty
	p.Reserved += qty

	p.Res[resID] = Reservation{
		ID:    resID,
		Qty:   qty,
		State: Created,
	}

	return nil
}

func (p *Product) Commit(rID ReservationID) error {
	res, ok := p.Res[rID]
	if !ok {
		return fmt.Errorf("%w: product %s missing reservation", ErrReservationNotFound, p.ID)
	}

	if res.State != Created {
		return fmt.Errorf("%w: reservation is %s", ErrInvalidReservationState, res.State)
	}

	p.Reserved -= res.Qty
	res.State = Committed

	p.Res[rID] = res

	return nil
}

func (p *Product) Cancel(rID ReservationID) error {
	res, ok := p.Res[rID]
	if !ok {
		return fmt.Errorf("%w: product %s missing reservation", ErrReservationNotFound, p.ID)
	}

	if res.State != Created {
		return fmt.Errorf("%w: reservation is %s", ErrInvalidReservationState, res.State)
	}

	p.Available += res.Qty
	p.Reserved -= res.Qty
	res.State = Cancelled

	p.Res[rID] = res

	return nil
}
