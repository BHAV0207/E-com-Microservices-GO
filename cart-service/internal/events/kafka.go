package events

import "github.com/segmentio/kafka-go"

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func KafkaReader(brokerURL, topic, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerURL},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}
