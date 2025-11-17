package service

import (
	"log/slog"

	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
	"github.com/hesampakdaman/inventory-service/internal/ports"
)

type Service struct {
	repo   ports.Repository
	logger *slog.Logger
	bus    *messagebus.Bus
}

func New(logger *slog.Logger, bus *messagebus.Bus, repo ports.Repository) *Service {
	logger = logger.With("component", "Service")
	svc := Service{repo: repo, logger: logger, bus: bus}

	messagebus.Register(bus, svc.Reserve)
	messagebus.Register(bus, svc.Create)

	return &svc
}
