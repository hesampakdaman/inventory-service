package service

import (
	"context"

	"github.com/hesampakdaman/inventory-service/internal/core/models"
)

func (s Service) GetAvailability(ctx context.Context, productID models.ProductID) (int, error) {
	product, err := s.repo.Get(ctx, productID)
	if err != nil {
		return 0, err
	}

	return product.Available, nil
}
