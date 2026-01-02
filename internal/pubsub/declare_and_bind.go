package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	DURABLE   SimpleQueueType = "durable"
	TRANSIENT SimpleQueueType = "transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	c, err := conn.Channel()
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, fmt.Errorf("failed to create channel from connection: %w", err)
	}
	q, err := c.QueueDeclare(
		queueName,
		queueType == DURABLE,
		queueType == TRANSIENT,
		queueType == TRANSIENT,
		false,
		nil)
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, fmt.Errorf("failed to create queue from connection: %w", err)
	}
	err = c.QueueBind(
		queueName,
		key,
		exchange,
		false,
		nil)
	if err != nil {
		return &amqp.Channel{}, amqp.Queue{}, fmt.Errorf("failed to create bind from connection: %w", err)
	}

	return c, q, nil
}
