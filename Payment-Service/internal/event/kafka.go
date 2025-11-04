package event

import "github.com/segmentio/kafka-go"

type KafkaProducer struct {
	Writer *kafka.Writer
}

func KafkaWriter(brokerUrl, topic string) *KafkaProducer {
	return &KafkaProducer{
		Writer: &kafka.Writer{
			Addr: kafka.TCP(brokerUrl),
			Topic: topic,
			Balancer: &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}
