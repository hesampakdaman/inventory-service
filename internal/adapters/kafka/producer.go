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

const InventoryTopic Topic = "inventory"

type Producer struct {
	logger   *slog.Logger
	client   *kgo.Client
	topicMap map[reflect.Type]Topic
}

func NewProducer(logger *slog.Logger, client *kgo.Client) *Producer {
	return newProducer(logger, client, InventoryTopic)
}

func NewProducerWithTopic(logger *slog.Logger, client *kgo.Client, topic Topic) *Producer {
	return newProducer(logger, client, topic)
}

func newProducer(logger *slog.Logger, client *kgo.Client, topic Topic) *Producer {
	topicMap := map[reflect.Type]Topic{
		reflect.TypeOf(commands.CreateProduct{}):    topic,
		reflect.TypeOf(commands.ReserveProduct{}):   topic,
		reflect.TypeOf(events.ProductCreated{}):     topic,
		reflect.TypeOf(events.ReservationCreated{}): topic,
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
