package pubsub

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishGameLog(ch *amqp.Channel, exchange, key, username, message string) error {
	gL := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     message,
		Username:    username,
	}

	err := PublishGob(
		ch,
		exchange,
		key,
		gL,
	)
	if err != nil {
		return fmt.Errorf("failed to publish game log: %w", err)
	}

	return nil
}
