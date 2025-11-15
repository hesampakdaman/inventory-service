package kafka

import (
	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/events"
	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
)

var registry = map[string]func() messagebus.Message{
	"ReservationCreated": func() messagebus.Message { return new(events.ReservationCreated) },
	"ReserveProduct":     func() messagebus.Message { return new(commands.ReserveProduct) },
}
