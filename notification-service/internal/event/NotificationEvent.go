package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (c *Consumer) StartConsuming() {
	fmt.Printf("📩 [%s] Kafka consumer started on topic: %s\n", c.ServiceName, c.Kafka.Reader.Config().Topic)

	for {
		msg, err := c.Kafka.Reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("❌ Kafka read error:", err)
			continue
		}

		var event GenericEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			fmt.Println("⚠️ Invalid Kafka event:", err)
			continue
		}

		fmt.Printf("📬 [%s] Received: %+v\n", c.ServiceName, event)

		// Build notification message
		userMsg := buildMessage(event)
		sendNotification(event.UserID, userMsg)

		// Store notification in DB
		notif := bson.M{
			"userId":    event.UserID,
			"orderId":   event.OrderID,
			"type":      event.EventType,
			"message":   userMsg,
			"status":    "SENT",
			"createdAt": time.Now(),
		}
		_, err = c.Collection.InsertOne(context.Background(), notif)
		if err != nil {
			fmt.Println("⚠️ DB insert error:", err)
		}
	}
}

// Helper for message content
func buildMessage(event GenericEvent) string {
	switch event.EventType {
	case "user-creted":
		return fmt.Sprintf("👋 Welcome aboard, User %s!", event.UserID)
	case "user-deleted":
		return fmt.Sprintf("👋 Goodbye, User %s! We're sad to see you go.", event.UserID)
	case "order-created":
		return fmt.Sprintf("✅ Order #%s placed successfully!", event.OrderID)
	case "payment-success":
		return fmt.Sprintf("💰 Payment for order #%s succeeded!", event.OrderID)
	case "payment-failed":
		return fmt.Sprintf("⚠️ Payment for order #%s failed. Please retry.", event.OrderID)
	default:
		return fmt.Sprintf("🔔 Update on your order #%s", event.OrderID)
	}
}

func sendNotification(userID, message string) {
	// For now, log the notification
	fmt.Printf("📨 Sending notification to user %s: %s\n", userID, message)
}
