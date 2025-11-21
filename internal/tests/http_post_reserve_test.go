package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hesampakdaman/inventory-service/internal/adapters/rest/handlers"
	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func TestHTTPPostReserveInvalidProductID(t *testing.T) {
	t.Parallel()

	//Arrange
	fx := NewFixture(t)
	payload := handlers.ReserveRequest{
		ProductID: uuid.New(),
		RequestID: uuid.New(),
		Qty:       1,
	}

	//Act
	resp, data, err := fx.DoJSONRaw(
		http.MethodPost,
		fmt.Sprintf("/products/%s/reserve", "not-a-uuid"),
		payload,
	)
	require.NoError(t, err)
	body := strings.TrimSpace(string(data))

	//Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Contains(t, body, "invalid UUID")
}

func TestHTTPPostReserveInvalidQuantity(t *testing.T) {
	t.Parallel()

	//Arrange
	fx := NewFixture(t)
	productID := uuid.New()
	payload := handlers.ReserveRequest{
		ProductID: productID,
		RequestID: uuid.New(),
		Qty:       0,
	}

	//Act
	resp, data, err := fx.DoJSONRaw(
		http.MethodPost,
		fmt.Sprintf("/products/%s/reserve", productID.String()),
		payload,
	)
	require.NoError(t, err)
	body := strings.TrimSpace(string(data))

	//Assert
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "Quantity must be greater than zero", body)
}

func TestHTTPPostReserveSuccess(t *testing.T) {
	t.Parallel()

	//Arrange
	fx := NewFixture(t)
	productID := models.ProductID(uuid.New())
	createCmd := commands.CreateProduct{
		ProductID:   productID,
		RequestID:   models.RequestID(uuid.New()),
		Title:       "test reserve product",
		Description: "created for HTTP reserve test",
		Qty:         10,
	}
	require.NoError(t, fx.Bus.Handle(t.Context(), &createCmd))
	request := handlers.ReserveRequest{
		ProductID: productID.UUID(),
		RequestID: uuid.New(),
		Qty:       3,
	}

	//Act
	resp, data, err := fx.DoJSONRaw(
		http.MethodPost,
		fmt.Sprintf("/products/%s/reserve", productID.String()),
		request,
	)
	require.NoError(t, err)
	body := strings.TrimSpace(string(data))

	//Assert
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Empty(t, body)
}
