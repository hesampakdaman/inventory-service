package service

import (
	"context"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

func (s *BankService) Transfer(
	ctx context.Context,
	fromAccountID, toAccountID string,
	amount float64,
) (domain.Transaction, domain.Transaction, error) {
	logger := s.logger.With(
		"from_account_id",
		fromAccountID,
		"to_account_id",
		toAccountID,
		"amount",
		amount,
	)

	logger.InfoContext(ctx, "Processing transfer")

	return domain.Transaction{}, domain.Transaction{}, nil
}
