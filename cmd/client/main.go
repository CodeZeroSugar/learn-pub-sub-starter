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
	fmt.Println("Starting Peril client...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ server: %s", err)
	} else {
		fmt.Println("Connection to RabbitMQ server was successful")
	}
	defer connection.Close()

	_, err = connection.Channel()
	if err != nil {
		log.Printf("Failed to create channel from connection: %s", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Printf("Error getting username: %s:", err)
	}

	_, _, err = pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.TRANSIENT,
	)
	if err != nil {
		log.Printf("Failed to get channel and queue: %s", err)
	}

	gamestate := gamelogic.NewGameState(username)
	if err = pubsub.SubscribeJSON(
		connection,
		string(routing.ExchangePerilDirect),
		string(routing.PauseKey)+"."+username,
		string(routing.PauseKey),
		pubsub.TRANSIENT,
		handlerPause(gamestate),
	); err != nil {
		log.Printf("failed to subscribe json: %s", err)
	}

Loop:
	for {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		arg := input[0]
		switch arg {
		case "spawn":
			if err = gamestate.CommandSpawn(input); err != nil {
				log.Printf("CommandSpawn failed: %s", err)
			}
		case "move":
			_, err := gamestate.CommandMove(input)
			if err != nil {
				log.Printf("CommandMove failed: %s", err)
			} else {
				log.Printf("Move was successful!")
			}
		case "status":
			gamestate.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			log.Printf("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			break Loop
		default:
			log.Printf("invalid command...")
		}

	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	fmt.Println("\nReceived signal:", sig)
	fmt.Println("Exiting program.")
}
