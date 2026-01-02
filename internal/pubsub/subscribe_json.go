package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {
	ch, _, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		queueType,
	)
	if err != nil {
		return fmt.Errorf("failed to declare and bind: %w", err)
	}

	deliveryChan, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	go func() {
		for delivery := range deliveryChan {
			data := delivery.Body
			var body T
			if err = json.Unmarshal(data, &body); err != nil {
				log.Printf("failed to unmarshal json:%s", err)
			}
			handler(body)

			delivery.Ack(false)
		}
	}()

	return nil
}
