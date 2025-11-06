package event
import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GenericEvent struct {
	EventType   string `json:"eventType"` // order.created, payment.success
	UserID      string `json:"userId"`
	OrderID     string `json:"orderId"`
	Message     string `json:"message"`
	Reservation string `json:"reservationId,omitempty"`
}

// StartKafkaConsumer - listens to order and payment topics
func StartKafkaConsumer(broker, topic, group string, notifCollection *mongo.Collection) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: group,
	})

	fmt.Printf("📩 Notification consumer started on topic '%s'\n", topic)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("❌ Error reading message:", err)
			continue
		}

		var event GenericEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			fmt.Println("❌ Invalid message:", err)
			continue
		}

		fmt.Printf("📬 Received event: %+v\n", event)

		// Build a user-friendly message
		msg := buildMessage(event)

		// Send or simulate notification
		sendNotification(event.UserID, msg)

		// Save notification to DB
		notif := bson.M{
			"userId":    event.UserID,
			"orderId":   event.OrderID,
			"type":      event.EventType,
			"message":   msg,
			"status":    "SENT",
			"createdAt": time.Now(),
		}
		_, err = notifCollection.InsertOne(context.Background(), notif)
		if err != nil {
			fmt.Println("⚠️ Failed to insert notification:", err)
		}
	}
}

func buildMessage(event GenericEvent) string {
	switch event.EventType {
	case "order.created":
		return fmt.Sprintf("Your order #%s has been placed successfully!", event.OrderID)
	case "payment.success":
		return fmt.Sprintf("Payment for order #%s succeeded!", event.OrderID)
	case "payment.failed":
		return fmt.Sprintf("Payment for order #%s failed. Please retry.", event.OrderID)
	case "order.shipped":
		return fmt.Sprintf("Your order #%s has been shipped!", event.OrderID)
	case "order.delivered":
		return fmt.Sprintf("Your order #%s has been delivered!", event.OrderID)
	default:
		return "You have a new update on your order."
	}
}

func sendNotification(userID, message string) {
	// For now, just log
	fmt.Printf("📨 [Notify User %s] %s\n", userID, message)

	// Later, you can plug an email/SMS gateway here.
}
