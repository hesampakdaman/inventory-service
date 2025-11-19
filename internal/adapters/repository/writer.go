package repository

import (
	"context"
	"errors"

	gocql "github.com/apache/cassandra-gocql-driver/v2"

	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type Writer struct {
	session *gocql.Session
}

func NewWriterRepository(sess *gocql.Session) Writer {
	return Writer{
		session: sess,
	}
}

func (w Writer) Get(ctx context.Context, id models.ProductID) (models.Product, error) {
	var (
		available   int
		reserved    int
		title       string
		description string
	)
	if err := w.session.Query(`
        SELECT available, reserved, title, description FROM products
        WHERE product_id = ?;`,
		id.String(),
	).ScanContext(ctx, &available, &reserved, &title, &description); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return models.Product{}, models.ErrProductNotFound
		}
		return models.Product{}, err
	}

	return models.Product{
		ID:          id,
		Available:   available,
		Reserved:    reserved,
		Title:       title,
		Description: description,
		Res:         map[models.ReservationID]models.Reservation{},
	}, nil
}

func (w Writer) GetWithReservation(
	ctx context.Context,
	productID models.ProductID,
	reservationID models.ReservationID,
) (models.Product, error) {
	var (
		available   int
		reserved    int
		title       string
		description string
	)
	if err := w.session.Query(`
        SELECT available, reserved, title, description FROM products
        WHERE product_id = ?;`,
		productID.String(),
	).ScanContext(ctx, &available, &reserved, &title, &description); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return models.Product{}, models.ErrProductNotFound
		}
		return models.Product{}, err
	}
	product := models.Product{
		ID:          productID,
		Available:   available,
		Reserved:    reserved,
		Title:       title,
		Description: description,
		Res:         map[models.ReservationID]models.Reservation{},
	}

	var (
		quantity int
		state    models.ReservationState
	)
	if err := w.session.Query(`
        SELECT quantity, state FROM reservations
        WHERE product_id = ? AND reservation_id = ?;`,
		productID.String(), reservationID.String(),
	).ScanContext(ctx, &quantity, &state); err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return product, nil // no reservation yet
		}
		return models.Product{}, err
	}
	product.Res[reservationID] = models.Reservation{
		ID:        reservationID,
		ProductID: productID,
		Qty:       quantity,
		State:     state,
	}

	return product, nil
}

func (w Writer) Save(ctx context.Context, p models.Product, req models.RequestID) error {
	m := make(map[string]any)
	applied, err := w.session.Query(`
        INSERT INTO requests (product_id, request_id, created_at)
        VALUES (?, ?, toTimestamp(now()))
        IF NOT EXISTS;`,
		p.ID.String(), req.String(),
	).MapScanCASContext(ctx, m)
	if err != nil {
		return err
	}
	if !applied {
		return nil
	}

	batch := w.session.Batch(gocql.LoggedBatch)

	batch.Query(`
        INSERT INTO products (product_id, available, reserved, title, description)
        VALUES (?, ?, ?, ?, ?);`,
		p.ID.String(), p.Available, p.Reserved, p.Title, p.Description,
	)
	for _, r := range p.Reservations() {
		batch.Query(`
        INSERT INTO reservations (product_id, reservation_id, qty, state)
        VALUES (?, ?, ?, ?);`,
			p.ID.String(), r.ID.String(), r.Qty, r.State,
		)
	}

	return batch.ExecContext(ctx)
}
