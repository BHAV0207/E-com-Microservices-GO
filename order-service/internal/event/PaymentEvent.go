package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

type PaymentCreationEvent struct {
	OrderID string  `json:"orderId"`
	UserID  string  `json:"userId"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
	Status  string  `json:"status"`
}

func (c *Consumer) ListenForPayments() {
	for {
		msg, err := c.KafkaConsumer.Reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("❌ Error reading Kafka message:", err)
			continue
		}

		log.Printf("📩 Received message on topic %s: %s", msg.Topic, string(msg.Value))

		var event PaymentCreationEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			fmt.Println("❌ Failed to unmarshal PaymentEvent:", err)
			continue
		}

		fmt.Printf("📩 Received payment event: %+v\n", event)

		// Update order status in DB
		status := "failed"
		if event.Status == "success" {
			status = "confirmed"
		}

		objID, err := primitive.ObjectIDFromHex(event.OrderID)
		if err != nil {
			fmt.Println("❌ Invalid ObjectID:", event.OrderID)
			continue
		}
		filter := bson.M{"_id": objID}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		update := bson.M{"$set": bson.M{"status": status}}
		_, err = c.Collection.UpdateOne(ctx, filter, update)
		cancel()

		if err != nil {
			fmt.Println("❌ Failed to update order status:", err)
		} else {
			fmt.Printf("✅ Updated order %s to status %s\n", event.OrderID, status)
		}
	}
}
