package repository

import (
	"context"
	"fmt"

	gocql "github.com/apache/cassandra-gocql-driver/v2"

	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type Writer struct {
	session *gocql.Session
}

func (w Writer) Save(ctx context.Context, p models.Product, req models.RequstID) error {
	applied, err := w.session.Query(`
        INSERT INTO inventory.requests (product_id, request_id, created_at)
        VALUES (?, ?, toTimestamp(now()))
        IF NOT EXISTS;`,
		p.ID, req,
	).ScanCASContext(ctx)
	if err != nil {
		return err
	}
	if !applied {
		return fmt.Errorf(
			"%w: request %s already handled for product %s",
			ErrDuplicateRequest,
			req,
			p.ID,
		)
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
