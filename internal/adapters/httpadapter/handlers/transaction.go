package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *httpHandler) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"transaction_id": ""}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *httpHandler) TransferHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromAccountID string  `json:"from_account_id"`
		ToAccountID   string  `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fromTxn, toTxn, err := h.service.Transfer(r.Context(), req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), domainErrToStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"withdrawal_transaction_id": fromTxn.ID,
		"deposit_transaction_id":    toTxn.ID,
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *httpHandler) ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("id")

	transactions := h.service.ListTransactions(r.Context(), accountID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
