package service

import (
	"log/slog"

	"github.com/hesampakdaman/inventory-service/internal/core/bus"
	"github.com/hesampakdaman/inventory-service/internal/ports"
)

type Service struct {
	repo   ports.Repository
	logger *slog.Logger
	bus    *bus.Bus
}

func New(repo ports.Repository, logger *slog.Logger, bus *bus.Bus) *Service {
	logger = logger.With("component", "Service")
	return &Service{repo: repo, logger: logger, bus: bus}
}
