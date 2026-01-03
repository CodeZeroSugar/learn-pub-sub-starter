package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(val)
	if err != nil {
		return fmt.Errorf("failed to encode val: %w", err)
	}
	msg := amqp.Publishing{
		ContentType: "application/gob",
		Body:        buff.Bytes(),
	}
	if err = ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		msg,
	); err != nil {
		return fmt.Errorf("failed to publish with context: %w", err)
	}

	return nil
}
