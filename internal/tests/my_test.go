package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/events"
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func TestCreateProduct(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)

	expected := &events.ProductCreated{
		ProductID:   models.ProductID(uuid.MustParse("ec1d1cb2-8d12-4686-9588-bb807e65aea7")),
		RequestID:   models.RequestID(uuid.MustParse("2b9a1a92-d280-4ca5-b420-6a2fa2413ca5")),
		Title:       "title",
		Description: "description",
		Available:   1,
	}

	var received *events.ProductCreated
	messagebus.Register(fx.Bus, func(ctx context.Context, e *events.ProductCreated) error {
		received = e
		return nil
	})

	// Act
	err := fx.Bus.Publish(t.Context(), expected.ProductID.UUID(), commands.CreateProduct{
		ProductID:   expected.ProductID,
		RequestID:   expected.RequestID,
		Title:       expected.Title,
		Description: expected.Description,
		Qty:         1,
	})
	require.NoError(t, err)

	// Assert
	require.Eventually(t,
		func() bool { return received != nil },
		time.Second,         // Wait for
		10*time.Millisecond, // Tick rate

	)
	assert.Equal(t, expected, received)
}
