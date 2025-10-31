package event

import "github.com/segmentio/kafka-go"

type KafkaProducer struct {
	Writer *kafka.Writer
}
