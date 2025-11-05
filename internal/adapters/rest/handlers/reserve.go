package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type ReserveRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	RequestID uuid.UUID `json:"request_id"`
	Qty       int       `json:"quantity"`
}

func (h *httpHandler) Reserve(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.PathValue("product_id") == "" {
		http.Error(w, "missing product_id", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(r.PathValue("product_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var req ReserveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Qty <= 0 {
		http.Error(w, "Qty must be greater than zero", http.StatusBadRequest)
		return
	}

	if err := h.bus.Handle(r.Context(), commands.ReserveProduct{
		ProductID: models.ProductID(productID),
		RequestID: req.RequestID,
		Qty:       req.Qty,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
