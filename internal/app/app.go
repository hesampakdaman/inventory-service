package app

import (
	"log/slog"
	"net/http"

	gocql "github.com/apache/cassandra-gocql-driver/v2"

	"github.com/hesampakdaman/inventory-service/internal/adapters/repository"
	"github.com/hesampakdaman/inventory-service/internal/adapters/rest"
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
	"github.com/hesampakdaman/inventory-service/internal/ports"
	"github.com/hesampakdaman/inventory-service/internal/service"
)

type App struct {
	Logger  *slog.Logger
	Bus     *messagebus.Bus
	Repo    ports.Repository
	Service *service.Service
	Router  http.Handler
}

func New(logger *slog.Logger, session *gocql.Session) *App {
	bus := messagebus.New()

	repo := repository.NewWriterRepository(session)
	svc := service.New(logger, bus, repo)

	messagebus.Register(bus, svc.Reserve)
	router := rest.NewRouter(bus)

	return &App{Logger: logger, Bus: bus, Repo: repo, Service: svc, Router: router}
}
