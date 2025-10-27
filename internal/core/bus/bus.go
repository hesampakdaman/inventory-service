package bus

import (
	"context"
	"fmt"
	"reflect"
)

type (
	Handler func(context.Context, any) error
	Message any
)

type Bus struct {
	handlers map[reflect.Type]Handler
}

func New() *Bus {
	return &Bus{
		handlers: make(map[reflect.Type]Handler),
	}
}

func Register[M Message](b *Bus, h func(context.Context, M) error) {
	t := reflect.TypeOf(*new(M))
	b.handlers[t] = func(ctx context.Context, v any) error {
		return h(ctx, v.(M))
	}
}

func (b *Bus) Handle(ctx context.Context, msg Message) error {
	t := reflect.TypeOf(msg)
	h, ok := b.handlers[t]
	if !ok {
		return fmt.Errorf("no handler for %T", msg)
	}
	return h(ctx, msg)
}
