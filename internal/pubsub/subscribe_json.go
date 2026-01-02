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
	handler func(T) AckType,
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
			ack := handler(body)

			switch ack {
			case ACK:
				delivery.Ack(false)
				log.Printf("Ack occurred.")
			case NACK_REQUEUE:
				delivery.Nack(false, true)
				log.Printf("NackRequeue occurred.")
			case NACK_DISCARD:
				delivery.Nack(false, false)
				log.Printf("NackDiscard occurred.")
			}

		}
	}()

	return nil
}
