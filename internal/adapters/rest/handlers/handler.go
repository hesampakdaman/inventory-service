package handlers

import (
	"github.com/hesampakdaman/inventory-service/internal/core/bus"
	"github.com/hesampakdaman/inventory-service/internal/service"
)

type httpHandler struct {
	bus *bus.Bus
}

func NewHTTPHandler(service *service.Service) *httpHandler {
	return &httpHandler{}
}
