package tests

import (
	"context"
	"sync"
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
		ProductID:   models.ProductID(uuid.New()),
		RequestID:   models.RequestID(uuid.New()),
		Title:       "title",
		Description: "description",
		Available:   1,
	}

	var received *events.ProductCreated
	_ = messagebus.Register(fx.Bus, func(ctx context.Context, e *events.ProductCreated) error {
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

func TestCreateProductIdempotent(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)

	expected := &events.ProductCreated{
		ProductID:   models.ProductID(uuid.New()),
		RequestID:   models.RequestID(uuid.New()),
		Title:       "title",
		Description: "description",
		Available:   1,
	}

	var (
		mu       sync.Mutex
		received = make([]*events.ProductCreated, 0, 2)
	)
	_ = messagebus.Register(fx.Bus, func(ctx context.Context, e *events.ProductCreated) error {
		mu.Lock()
		received = append(received, e)
		mu.Unlock()
		return nil
	})

	cmd := commands.CreateProduct{
		ProductID:   expected.ProductID,
		RequestID:   expected.RequestID,
		Title:       expected.Title,
		Description: expected.Description,
		Qty:         1,
	}

	// Act
	require.NoError(t, fx.Bus.Publish(t.Context(), cmd.ProductID.UUID(), cmd))
	require.NoError(t, fx.Bus.Publish(t.Context(), cmd.ProductID.UUID(), cmd))

	// Assert
	require.Eventually(t,
		func() bool {
			mu.Lock()
			defer mu.Unlock()
			return len(received) == 2
		},
		time.Second,         // Wait for
		10*time.Millisecond, // Tick rate
	)

	mu.Lock()
	defer mu.Unlock()

	if assert.Len(t, received, 2) {
		assert.Equal(t, expected, received[0])
		assert.Equal(t, expected, received[1])
	}
}
