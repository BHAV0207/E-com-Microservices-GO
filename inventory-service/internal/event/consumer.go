package event

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Consumer struct {
	KafkaConsumer *KafkaConsumer
	Collection    *mongo.Collection
	ResvcColl     *mongo.Collection
}

func NewConsumer(brokerURL, topics, groupID string, collection *mongo.Collection) *Consumer {
	return &Consumer{
		KafkaConsumer: KafkaReader(brokerURL, topics, groupID),
		Collection:    collection,
	}
}
