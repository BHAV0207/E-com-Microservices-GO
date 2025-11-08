package event

import "go.mongodb.org/mongo-driver/mongo"

type GenericEvent struct {
	EventType   string `json:"eventType"` // e.g., order.created, payment.success
	UserID      string `json:"userId"`
	OrderID     string `json:"orderId"`
	Message     string `json:"message"`
	Reservation string `json:"reservationId,omitempty"`
}

type Consumer struct {
	Kafka       *KafkaConsumer
	Collection  *mongo.Collection
	ServiceName string
}


func NewConsumer(brokerURL, topic, groupID, serviceName string, collection *mongo.Collection) *Consumer {
	return &Consumer{
		Kafka:       NewKafkaConsumer(brokerURL, topic, groupID),
		Collection:  collection,
		ServiceName: serviceName,
	}
}