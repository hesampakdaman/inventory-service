package handlers

import (
	"github.com/hesampakdaman/inventory-service/internal/service"
)

// httpHandler handles HTTP requests for banking operations.
type httpHandler struct {
	service *service.Service
}

func NewHTTPHandler(service *service.Service) *httpHandler {
	return &httpHandler{service: service}
}
