package kafka

import (
	"context"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"

	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
)

type Consumer struct {
	logger *slog.Logger
	bus    *messagebus.Bus
	client *kgo.Client
}

func NewConsumer(logger *slog.Logger, bus *messagebus.Bus, client *kgo.Client) *Consumer {
	return &Consumer{
		logger: logger,
		bus:    bus,
		client: client,
	}
}

func (c Consumer) Consume(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		fetches := c.client.PollFetches(ctx)
		for _, rec := range fetches.Records() {
			msg, err := Decode(rec.Value)
			if err != nil {
				c.logger.Error("decode failed", slog.Any("err", err))
				continue
			}

			if err := c.bus.Handle(ctx, msg); err != nil {
				c.logger.Error("handler failed", slog.Any("err", err))
				continue
			}
		}
	}
}
