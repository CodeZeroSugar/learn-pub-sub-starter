package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerLogs() func(gamelog routing.GameLog) pubsub.AckType {
	return func(gamelog routing.GameLog) pubsub.AckType {
		defer fmt.Print("> ")
		if gamelog.Message == "" || gamelog.Username == "" {
			return pubsub.NACK_DISCARD
		}
		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			log.Printf("failed to write to gamelog: %s", err)
			return pubsub.NACK_REQUEUE
		}
		return pubsub.ACK
	}
}
