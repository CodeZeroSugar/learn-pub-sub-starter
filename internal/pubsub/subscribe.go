package pubsub

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func(amqp.Delivery) (T, error),
) error {
	ch, _, err := DeclareAndBind(
		conn,
		exchange,
		queueName,
		key,
		simpleQueueType,
	)
	if err != nil {
		return fmt.Errorf("failed to declare and bind: %w", err)
	}

	if err = ch.Qos(
		10,
		0,
		false,
	); err != nil {
		return fmt.Errorf("qos failed for channel: %w", err)
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

			body, err := unmarshaller(delivery)
			if err != nil {
				log.Fatal("failed to unmarshal delivery:", err)
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
