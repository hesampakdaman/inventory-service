package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

type availabilityResponse struct {
	ProductID uuid.UUID `json:"product_id"`
	Available int       `json:"available"`
}

func (h *httpHandler) GetAvailability(w http.ResponseWriter, r *http.Request) {
	productParam := r.PathValue("id")
	if productParam == "" {
		http.Error(w, "missing product id", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(productParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	available, err := h.svc.GetAvailability(r.Context(), models.ProductID(productID))
	if err != nil {
		if errors.Is(err, models.ErrProductNotFound) {
			http.Error(w, models.ErrProductNotFound.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(availabilityResponse{
		ProductID: productID,
		Available: available,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
