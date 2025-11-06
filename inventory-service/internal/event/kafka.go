package event

import "github.com/segmentio/kafka-go"

type KafkaConsumer struct {
	Reader *kafka.Reader
}

func KafkaReader(brokerUrl, topic, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{brokerUrl},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}
