package handlers

import (
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
	"github.com/hesampakdaman/inventory-service/internal/service"
)

type httpHandler struct {
	bus *messagebus.Bus
	svc *service.Service
}

func NewHTTPHandler(bus *messagebus.Bus, svc *service.Service) *httpHandler {
	return &httpHandler{
		bus: bus,
		svc: svc,
	}
}
