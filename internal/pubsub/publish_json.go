package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal json to bytes while publishing: %w", err)
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonBytes,
	}
	if err = ch.PublishWithContext(context.Background(), exchange, key, false, false, msg); err != nil {
		return fmt.Errorf("channel failed to publish: %s", err)
	}

	return nil
}
