package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")

		moveOutcome := gs.HandleMove(move)
		switch moveOutcome {
		case gamelogic.MoveOutcomeMakeWar:
			if err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routing.WarRecognitionsPrefix+"."+move.Player.Username,
				gamelogic.RecognitionOfWar{Attacker: move.Player, Defender: gs.GetPlayerSnap()},
			); err != nil {
				log.Printf("Error occurred while publishing message in move handler: %s", err)
				return pubsub.NACK_REQUEUE
			} else {
				return pubsub.ACK
			}

		case gamelogic.MoveOutComeSafe:
			return pubsub.ACK

		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NACK_DISCARD
		default:
			return pubsub.NACK_DISCARD
		}
	}
}
