package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"reflect"

	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
)

type Producer struct {
	logger *slog.Logger
	client *kgo.Client
}

func NewProducer(logger *slog.Logger, client *kgo.Client) *Producer {
	return &Producer{logger: logger, client: client}
}

func (p Producer) Publish(ctx context.Context, key uuid.UUID, topic string, msg any) error {
	T := reflect.TypeOf(msg).Name()

	keyBytes, err := key.MarshalBinary()
	if err != nil {
		return err
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	value, err := json.Marshal(Envelope{
		Type: T,
		Data: payload,
	})
	if err != nil {
		return err
	}

	record := &kgo.Record{
		Key:   keyBytes,
		Value: value,
		Topic: topic,
	}
	return p.client.ProduceSync(ctx, record).FirstErr()
}
