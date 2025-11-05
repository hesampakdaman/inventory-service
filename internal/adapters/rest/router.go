package rest

import (
	"net/http"

	"github.com/hesampakdaman/inventory-service/internal/adapters/rest/handlers"
	"github.com/hesampakdaman/inventory-service/internal/core/bus"
)

func NewRouter(bus *bus.Bus) *http.ServeMux {
	handler := handlers.NewHTTPHandler(bus)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /products/{id}/reserve", handler.Reserve)

	return mux
}
