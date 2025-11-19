package app

import (
	"log/slog"
	"net/http"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/hesampakdaman/inventory-service/internal/adapters/kafka"
	"github.com/hesampakdaman/inventory-service/internal/adapters/repository"
	"github.com/hesampakdaman/inventory-service/internal/adapters/rest"
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
	"github.com/hesampakdaman/inventory-service/internal/ports"
	"github.com/hesampakdaman/inventory-service/internal/service"
)

type App struct {
	Logger   *slog.Logger
	Bus      *messagebus.Bus
	Repo     ports.Repository
	Service  *service.Service
	Router   http.Handler
	Consumer *kafka.Consumer
}

func New(logger *slog.Logger, session *gocql.Session, kafkaClient *kgo.Client) *App {
	bus := messagebus.New(kafka.NewProducer(logger, kafkaClient))

	repo := repository.NewWriterRepository(session)
	svc := service.New(logger, bus, repo)

	router := rest.NewRouter(bus, svc)

	consumer := kafka.NewConsumer(logger, bus, kafkaClient)

	return &App{
		Logger:   logger,
		Bus:      bus,
		Repo:     repo,
		Service:  svc,
		Router:   router,
		Consumer: consumer,
	}
}
