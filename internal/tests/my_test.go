package tests

import (
	"fmt"
	"testing"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func TestLeTest(t *testing.T) {
	fx := NewFixture(t)

	batch := fx.Session.Batch(gocql.LoggedBatch)

	batch.Query(`
    INSERT INTO products (product_id, available, reserved, title, description)
    VALUES (?, ?, ?, ?, ?);`,
		"ec1d1cb2-8d12-4686-9588-bb807e65aea7", 2, 3, "Title", "Desc",
	)

	product := models.Product{
		ID: models.ProductID(uuid.MustParse("ec1d1cb2-8d12-4686-9588-bb807e65aea7")),
	}
	err := fx.App.Bus.Publish(t.Context(), product.ID.UUID(), commands.ReserveProduct{
		ProductID: product.ID,
		RequestID: models.RequestID{},
		Qty:       1,
	})
	require.NoError(t, err)

	if err := batch.ExecContext(t.Context()); err != nil {
		fmt.Println(err)
		t.Fatalf("%s", err)
	}
}
