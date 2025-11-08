package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/BHAV0207/payment-service/internal/handler"
	"github.com/BHAV0207/payment-service/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		c.sendNotification(event.UserID, userMsg)

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

func (c *Consumer) sendNotification(userID, message string) {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		fmt.Println("❌ Invalid userID format:", objID)
		return
	}
	user, err := service.GetUserByID(userID)
	if err != nil {
		fmt.Printf("❌ Failed to fetch user %s: %v\n", userID, err)
		return
	}

	subject := "Notification from E-com Website"
	body := fmt.Sprintf("Hey %s,<br><br>%s<br><br>– Team E-com", user.Name, message)

	err = handler.SendEmail(user.Email, subject, body)
	if err != nil {
		fmt.Printf("❌ Failed to send email to %s: %v\n", user.Email, err)
		return
	}

	fmt.Printf("✅ Notification email sent to %s\n", user.Email)
}
