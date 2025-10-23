package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/BHAV0207/cart-service/internal/service"
	"github.com/BHAV0207/cart-service/pkg"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserCreatedEvent struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
}

type Consumer struct {
	KafkaConsumer *KafkaConsumer
	Collection    *mongo.Collection
}

func NewConsumer(brokerURL, topics, groupID string, collection *mongo.Collection) *Consumer {
	return &Consumer{
		KafkaConsumer: KafkaReader(brokerURL, topics, groupID),
		Collection:    collection,
	} 
}

func (c *Consumer) Consume() {
	log.Println("🚀 Cart Service Kafka Consumer started and listening for user-created events...")

	for {
		msg, err := c.KafkaConsumer.Reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("❌ Error reading message from Kafka: %v", err)
			continue
		}

		log.Printf("📩 Received message on topic %s: %s", msg.Topic, string(msg.Value))

		var event UserCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("❌ Error decoding Kafka message: %v", err)
			continue
		}

		// Convert userId from string to ObjectID
		userID, err := primitive.ObjectIDFromHex(event.UserId)
		if err != nil {
			log.Printf("⚠️ Invalid userId format: %v", err)
			continue
		}

		// Create an empty cart for the new user
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cart := pkg.Cart{
			UserId: userID,
			Items:  []pkg.CartItem{},
		}

		err = service.CreateCart(ctx, c.Collection, cart)
		if err != nil {
			log.Printf("❌ Failed to create cart for user %s: %v", event.UserId, err)
			continue
		}

		log.Printf("🛒 Successfully created cart for new user: %s", event.UserId)
	}
}
