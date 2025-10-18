package handlers

import (
	"github.com/hesampakdaman/banking-service/internal/service"
)

// httpHandler handles HTTP requests for banking operations.
type httpHandler struct {
	service *service.BankService
}

func NewHTTPHandler(service *service.BankService) *httpHandler {
	return &httpHandler{service: service}
}
