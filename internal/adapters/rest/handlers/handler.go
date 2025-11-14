package handlers

import (
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
)

type httpHandler struct {
	bus *messagebus.Bus
}

func NewHTTPHandler(bus *messagebus.Bus) *httpHandler {
	return &httpHandler{bus}
}
