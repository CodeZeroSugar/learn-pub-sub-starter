package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	gamelogic.PrintServerHelp()
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ server: %s", err)
	} else {
		fmt.Println("Connection to RabbitMQ server was successful")
	}
	defer connection.Close()

	ch, err := connection.Channel()
	if err != nil {
		log.Printf("Failed to create channel from connection: %s", err)
	}

	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.DURABLE,
	)
	if err != nil {
		log.Printf("server failed to declare and bind: %s", err)
	}

	err = pubsub.SubscribeGob(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.DURABLE,
		handlerLogs(),
	)
	if err != nil {
		log.Printf("and error happened trying to receive gamelogs: %s", err)
	}

Loop:
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		arg := input[0]
		switch arg {
		case "pause":
			log.Printf("sending a pause message...")
			data := routing.PlayingState{
				IsPaused: true,
			}
			if err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, data); err != nil {
				log.Printf("Failed to PublishJSON: %s", err)
			}
		case "resume":
			log.Printf("sending a resume message...")
			data := routing.PlayingState{
				IsPaused: false,
			}
			if err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, string(routing.PauseKey), data); err != nil {
				log.Printf("Failed to PublishJSON: %s", err)
			}
		case "quit":
			log.Printf("exiting... ")
			break Loop

		default:
			log.Printf("invalid command")
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	fmt.Println("\nReceived signal:", sig)
	fmt.Println("Exiting program.")
}
