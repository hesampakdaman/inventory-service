package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/hesampakdaman/inventory-service/internal/core/messagebus"
)

type Envelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func Decode(b []byte) (messagebus.Message, error) {
	var env Envelope
	if err := json.Unmarshal(b, &env); err != nil {
		return nil, err
	}

	newMsgFn, ok := registry[env.Type]
	if !ok {
		return nil, fmt.Errorf("unknown message type: %s", env.Type)
	}

	msg := newMsgFn()
	if err := json.Unmarshal(env.Data, &msg); err != nil {
		return nil, err
	}

	return msg, nil
}
