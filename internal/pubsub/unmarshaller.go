package pubsub

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func unmarshaller[T any](msg amqp.Delivery) (T, error) {
	var body T
	if msg.ContentType == "application/json" {
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return body, fmt.Errorf("umarshaller failed to unmarshal json: %w", err)
		}
		return body, nil
	} else {
		buff := bytes.NewBuffer(msg.Body)
		dec := gob.NewDecoder(buff)
		err := dec.Decode(&body)
		if err != nil {
			return body, fmt.Errorf("unmarshaller failed to decode message: %w", err)
		}
		return body, nil
	}
}
