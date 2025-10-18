package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/hesampakdaman/banking-service/internal/adapter/repository"
	"github.com/hesampakdaman/banking-service/internal/adapter/rest"
	"github.com/hesampakdaman/banking-service/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(log.Writer(), nil))

	// Initialize repository & service layer
	repo := repository.Postgres{}
	bankService := service.NewBankService(&repo, logger)

	// Initialize http server
	mux := rest.NewRouter(bankService)
	loggedMux := rest.LoggingMiddleware(mux, logger)
	server := &http.Server{
		Addr:    ":8080",
		Handler: loggedMux,
	}

	logger.Info("Starting banking-service on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
