package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(body gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		outcome, winner, loser := gs.HandleWar(body)
		switch outcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NACK_REQUEUE
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NACK_DISCARD
		case gamelogic.WarOutcomeOpponentWon:
			outcomeString := fmt.Sprintf("%s won a war against %s", winner, loser)
			err := pubsub.PublishGameLog(
				ch,
				routing.ExchangePerilTopic,
				routing.GameLogSlug+".#",
				gs.Player.Username,
				outcomeString,
			)
			if err != nil {
				return pubsub.NACK_REQUEUE
			}
			return pubsub.ACK
		case gamelogic.WarOutcomeYouWon:
			outcomeString := fmt.Sprintf("%s won a war against %s", winner, loser)
			err := pubsub.PublishGameLog(
				ch,
				routing.ExchangePerilTopic,
				routing.GameLogSlug+"."+gs.Player.Username,
				gs.Player.Username,
				outcomeString,
			)
			if err != nil {
				return pubsub.NACK_REQUEUE
			}
			return pubsub.ACK
		case gamelogic.WarOutcomeDraw:
			outcomeString := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			err := pubsub.PublishGameLog(
				ch,
				routing.ExchangePerilTopic,
				routing.GameLogSlug+"."+gs.Player.Username,
				gs.Player.Username,
				outcomeString,
			)
			if err != nil {
				return pubsub.NACK_REQUEUE
			}
			return pubsub.ACK
		default:
			fmt.Println("Error: war outcome was not defined")
			return pubsub.NACK_DISCARD
		}
	}
}
