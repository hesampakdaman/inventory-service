package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hesampakdaman/inventory-service/internal/adapters/rest/handlers"
	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func TestHTTPGetAvailability(t *testing.T) {
	// Arrange
	t.Parallel()
	fx := NewFixture(t)
	expected := handlers.AvailabilityResponse{
		ProductID: uuid.New(),
		Available: 5,
	}
	cmd := commands.CreateProduct{
		ProductID:   models.ProductID(expected.ProductID),
		RequestID:   models.RequestID(uuid.New()),
		Title:       "test product",
		Description: "integration test product",
		Qty:         expected.Available,
	}
	require.NoError(t, fx.Bus.Handle(t.Context(), &cmd))

	// Act
	var respBody handlers.AvailabilityResponse
	resp, err := fx.DoJSON(
		http.MethodGet,
		fmt.Sprintf("/products/%s/available", expected.ProductID),
		nil,
		&respBody,
	)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expected, respBody)
}
