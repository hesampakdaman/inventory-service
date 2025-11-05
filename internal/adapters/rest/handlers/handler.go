package handlers

import (
	"github.com/hesampakdaman/inventory-service/internal/core/bus"
)

type httpHandler struct {
	bus *bus.Bus
}

func NewHTTPHandler(bus *bus.Bus) *httpHandler {
	return &httpHandler{bus}
}
