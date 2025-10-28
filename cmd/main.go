package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/hesampakdaman/inventory-service/internal/adapters/rest"
	"github.com/hesampakdaman/inventory-service/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(log.Writer(), nil))

	svc := service.New(nil, logger)

	// Initialize http server
	mux := rest.NewRouter(svc)
	loggedMux := rest.LoggingMiddleware(mux, logger)
	server := &http.Server{
		Addr:    ":8080",
		Handler: loggedMux,
	}

	logger.Info("Starting inventory-service on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
