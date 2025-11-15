package kafka

import (
	"github.com/hesampakdaman/inventory-service/internal/core/events"
)

var registry = map[string]func() any{
	"ReserveProduct": func() any { return &events.ReservationCreated{} },
}
