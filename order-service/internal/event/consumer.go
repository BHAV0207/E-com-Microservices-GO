package event

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type PaymentEvent struct {
	OrderID string  `json:"orderId"`
	UserID  string  `json:"userId"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
	Status  string  `json:"status"`
}
type OrderConsumer struct {
	Collection *mongo.Collection
	Kafka      *KafkaConsumer
}

func (oc *OrderConsumer) ListenForPayments(ctx context.Context) {
	for {
		msg, err := oc.Kafka.Reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println("❌ Error reading Kafka message:", err)
			continue
		}

		var event PaymentEvent
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

		filter := bson.M{"_id": event.OrderID}
		update := bson.M{"$set": bson.M{"status": status}}
		_, err = oc.Collection.UpdateOne(ctx, filter, update)
		if err != nil {
			fmt.Println("❌ Failed to update order status:", err)
		} else {
			fmt.Printf("✅ Updated order %s to status %s\n", event.OrderID, status)
		}
	}
}
