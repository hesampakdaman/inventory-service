package kafka

import (
	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/events"
)

var registry = map[string]func() any{
	"CreateProduct":  func() any { return new(commands.CreateProduct) },
	"ProductCreated": func() any { return new(events.ProductCreated) },
	"ReserveProduct": func() any { return new(events.ReservationCreated) },
}
