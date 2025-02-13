package rest

import (
	"net/http"

	"github.com/hesampakdaman/inventory-service/internal/adapters/rest/handlers"
	"github.com/hesampakdaman/inventory-service/internal/service"
)

func NewRouter(svc *service.Service) *http.ServeMux {
	_ = handlers.NewHTTPHandler(svc)

	mux := http.NewServeMux()

	return mux
}
