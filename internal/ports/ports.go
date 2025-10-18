package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hesampakdaman/banking-service/internal/domain"
)

type Repository interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Account, error)
	Save(ctx context.Context, account domain.Account) error
}
