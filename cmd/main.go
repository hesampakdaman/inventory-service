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

	"github.com/hesampakdaman/inventory-service/internal/adapters/rest"
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
)

func main() {
	logger := slog.New(slog.NewTextHandler(log.Writer(), nil))

	// Initialize http server
	mux := rest.NewRouter(messagebus.New())
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
