package service

import (
	"log/slog"

	"github.com/hesampakdaman/inventory-service/internal/ports"
)

type Service struct {
	repo   ports.Repository
	logger *slog.Logger
}

func New(repo ports.Repository, logger *slog.Logger) *Service {
	logger = logger.With("component", "Service")
	return &Service{repo: repo, logger: logger}
}
