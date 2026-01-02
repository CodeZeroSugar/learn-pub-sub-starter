package pubsub

type AckType string

const (
	ACK          AckType = "Ack"
	NACK_REQUEUE AckType = "NackRequeue"
	NACK_DISCARD AckType = "NackDiscard"
)
