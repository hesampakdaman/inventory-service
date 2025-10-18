package service

import (
	"context"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

func (s *BankService) GetAccount(ctx context.Context, accountID string) (domain.Account, error) {
	logger := s.logger.With("account_id", accountID)
	logger.InfoContext(ctx, "Successfully retrieved account")
	return domain.Account{}, nil
}

func (s *BankService) ListAccounts(ctx context.Context) []domain.Account {
	s.logger.InfoContext(ctx, "Listing all accounts")
	s.logger.InfoContext(ctx, "Successfully listed accounts", "count", 0)
	return nil
}

func (s *BankService) ListTransactions(ctx context.Context, accountID string) []domain.Transaction {
	logger := s.logger.With("account_id", accountID)
	logger.InfoContext(ctx, "Listing all transactions for account")
	logger.InfoContext(ctx, "Successfully listed transactions for account", "count", 0)
	return nil
}
