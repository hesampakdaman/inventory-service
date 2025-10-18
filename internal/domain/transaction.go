package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
)

type Transaction struct {
	ID        uuid.UUID
	AccountID string
	Type      TransactionType
	Amount    float64
	Timestamp time.Time
}

func NewTransaction(accountID string, txnType TransactionType, amount float64) (Transaction, error) {
	if accountID == "" {
		return Transaction{}, ErrInvalidAccountID
	}

	if txnType != Deposit && txnType != Withdrawal {
		return Transaction{}, ErrInvalidTransactionType
	}

	if amount <= 0 {
		return Transaction{}, ErrInvalidAmount
	}

	return Transaction{
		ID:        GetUUID(),
		AccountID: accountID,
		Type:      txnType,
		Amount:    amount,
		Timestamp: GetTimeNow(),
	}, nil
}
