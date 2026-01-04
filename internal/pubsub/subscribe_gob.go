package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
) error {
	err := subscribe(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
		handler,
		unmarshaller,
	)
	if err != nil {
		return fmt.Errorf("failure during subscribe gob: %w", err)
	}
	return nil
}
