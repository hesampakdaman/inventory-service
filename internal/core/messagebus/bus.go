package messagebus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/google/uuid"

	"github.com/hesampakdaman/inventory-service/internal/ports"
)

type (
	Handler   func(context.Context, any) error
	HandlerID = uuid.UUID
	Message   any
)

var ErrNoHandler = errors.New("no handler registered")

type Bus struct {
	mu       sync.RWMutex
	handlers map[reflect.Type]map[uuid.UUID]Handler
	producer ports.Producer
}

func New(producer ports.Producer) *Bus {
	return &Bus{
		handlers: make(map[reflect.Type]map[uuid.UUID]Handler),
		producer: producer,
	}
}

func Register[M Message](b *Bus, h func(context.Context, *M) error) HandlerID {
	T := reflect.TypeOf(new(M))
	handlerID := uuid.New()

	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.handlers[T]; !ok {
		b.handlers[T] = make(map[uuid.UUID]Handler)
	}

	b.handlers[T][handlerID] = func(ctx context.Context, v any) error {
		return h(ctx, v.(*M))
	}

	return handlerID
}

func Unregister[M Message](b *Bus, handlerID HandlerID) {
	T := reflect.TypeOf(new(M))

	b.mu.Lock()
	defer b.mu.Unlock()

	handlers, ok := b.handlers[T]
	if !ok {
		return
	}

	delete(handlers, handlerID)
	if len(handlers) == 0 {
		delete(b.handlers, T)
	}
}

func (b *Bus) Handle(ctx context.Context, msg Message) error {
	T := reflect.TypeOf(msg)
	b.mu.RLock()
	handlers, ok := b.handlers[T]
	b.mu.RUnlock()

	if !ok {
		return fmt.Errorf("Message type %s: %w", T, ErrNoHandler)
	}

	for id, handler := range handlers {
		if err := handler(ctx, msg); err != nil {
			return fmt.Errorf("handler %s: %w", id, err)
		}
	}

	return nil
}

func (b *Bus) Publish(ctx context.Context, key uuid.UUID, msg Message) error {
	return b.producer.Publish(ctx, key, msg)
}
