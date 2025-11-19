package rest

import (
	"net/http"

	"github.com/hesampakdaman/inventory-service/internal/adapters/rest/handlers"
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
	"github.com/hesampakdaman/inventory-service/internal/service"
)

func NewRouter(bus *messagebus.Bus, svc *service.Service) *http.ServeMux {
	handler := handlers.NewHTTPHandler(bus, svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /products/{id}/available", handler.GetAvailability)
	mux.HandleFunc("POST /products/{id}/reserve", handler.Reserve)

	return mux
}
