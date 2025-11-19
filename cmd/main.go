package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/hesampakdaman/inventory-service/internal/adapters/rest"
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
	"github.com/hesampakdaman/inventory-service/internal/core/models"
	"github.com/hesampakdaman/inventory-service/internal/ports"
	"github.com/hesampakdaman/inventory-service/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(log.Writer(), nil))

	bus := messagebus.New(noopProducer{})
	svc := service.New(logger, bus, noopRepo{})

	// Initialize http server
	mux := rest.NewRouter(bus, svc)
	loggedMux := rest.LoggingMiddleware(mux, logger)
	server := &http.Server{
		Addr:    ":8080",
		Handler: loggedMux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("Starting inventory-service on :8080")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_ = server.Shutdown(ctx)
}

type noopProducer struct{}

func (noopProducer) Publish(context.Context, uuid.UUID, any) error {
	return nil
}

type noopRepo struct{}

func (noopRepo) Get(context.Context, models.ProductID) (models.Product, error) {
	return models.Product{}, models.ErrProductNotFound
}

func (noopRepo) GetWithReservation(
	context.Context,
	models.ProductID,
	models.ReservationID,
) (models.Product, error) {
	return models.Product{}, models.ErrProductNotFound
}

func (noopRepo) Save(context.Context, models.Product, models.RequestID) error {
	return nil
}

var (
	_ ports.Repository = (*noopRepo)(nil)
	_ ports.Producer   = (*noopProducer)(nil)
)
