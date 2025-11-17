package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/google/uuid"

	"github.com/hesampakdaman/inventory-service/internal/core/commands"
	"github.com/hesampakdaman/inventory-service/internal/core/events"

	"github.com/twmb/franz-go/pkg/kgo"
)

type Topic string

var ErrTopicNotFound = errors.New("Topic not found")

var InventoryTopic Topic = "inventory"

type Producer struct {
	logger   *slog.Logger
	client   *kgo.Client
	topicMap map[reflect.Type]Topic
}

func NewProducer(logger *slog.Logger, client *kgo.Client) *Producer {
	topicMap := map[reflect.Type]Topic{
		reflect.TypeOf(commands.ReserveProduct{}):   InventoryTopic,
		reflect.TypeOf(events.ProductCreated{}):     InventoryTopic,
		reflect.TypeOf(events.ReservationCreated{}): InventoryTopic,
	}
	return &Producer{logger: logger, client: client, topicMap: topicMap}
}

func (p Producer) Publish(ctx context.Context, key uuid.UUID, msg any) error {
	T := reflect.TypeOf(msg)
	topic, ok := p.topicMap[T]
	if !ok {
		return fmt.Errorf("%s: %w", T, ErrTopicNotFound)
	}

	keyBytes, err := key.MarshalBinary()
	if err != nil {
		return err
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	value, err := json.Marshal(Envelope{
		Type: T.Name(),
		Data: payload,
	})
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Key:   keyBytes,
		Value: value,
		Topic: string(topic),
	}
	return p.client.ProduceSync(ctx, record).FirstErr()
}
