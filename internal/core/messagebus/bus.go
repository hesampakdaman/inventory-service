package messagebus

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hesampakdaman/inventory-service/internal/ports"
)

type (
	Handler func(context.Context, any) error
	Message any
)

var ErrNoHandler = errors.New("no handler registered")

type Bus struct {
	handlers  map[reflect.Type]Handler
	publisher ports.Publisher
}

func New() *Bus {
	return &Bus{
		handlers: make(map[reflect.Type]Handler),
	}
}

func Register[M Message](b *Bus, h func(context.Context, M) error) {
	T := reflect.TypeOf(*new(M))
	b.handlers[T] = func(ctx context.Context, v any) error {
		return h(ctx, v.(M))
	}
}

func (b *Bus) Handle(ctx context.Context, msg Message) error {
	T := reflect.TypeOf(msg)
	h, ok := b.handlers[T]
	if !ok {
		return fmt.Errorf("Message type %s: %w", T, ErrNoHandler)
	}
	return h(ctx, msg)
}

func (b Bus) Publish(ctx context.Context, msg Message) error {
	return b.publisher.Publish(ctx, msg)
}
