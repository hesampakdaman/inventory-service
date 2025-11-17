package tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func TestLeTest(t *testing.T) {
	fx := NewFixture(t)

	cmd := commands.CreateProduct{
		ProductID:   models.ProductID(uuid.MustParse("ec1d1cb2-8d12-4686-9588-bb807e65aea7")),
		RequestID:   models.RequestID{},
		Title:       "title",
		Description: "description",
		Qty:         5,
	}
	require.NoError(t, fx.App.Bus.Publish(t.Context(), cmd.ProductID.UUID(), cmd))

	product := models.Product{
		ID: models.ProductID(uuid.MustParse("ec1d1cb2-8d12-4686-9588-bb807e65aea7")),
	}
	err := fx.App.Bus.Publish(t.Context(), product.ID.UUID(), commands.ReserveProduct{
		ProductID: product.ID,
		RequestID: models.RequestID{},
		Qty:       1,
	})
	require.NoError(t, err)
}
