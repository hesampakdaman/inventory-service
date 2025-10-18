package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/hesampakdaman/banking-service/internal/domain"
)

type Postgres struct{}

func (r *Postgres) Get(ctx context.Context, id uuid.UUID) (domain.Account, error) {
	return domain.Account{}, nil
}

func (r *Postgres) Save(ctx context.Context, account domain.Account) error {
	return nil
}
