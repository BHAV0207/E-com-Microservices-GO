package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentEvent struct {
	OrderID       string `json:"orderId"`
	ReservationID string `json:"reservationId"`
	Status        string `json:"status"` // "success" or "failed"
}

func ConsumePaymentEvents(broker, topic, group string, inventoryCol, reservationCol *mongo.Collection) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: group,
	})

	fmt.Println("🚀 Inventory consumer started on topic:", topic)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("❌ Error reading message:", err)
			continue
		}

		var event PaymentEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			fmt.Println("❌ Invalid payment event:", err)
			continue
		}

		fmt.Printf("📥 Received payment event: %+v\n", event)

		switch event.Status {
		case "success":
			commitInventory(inventoryCol, reservationCol, event)
		case "failed":
			cancelReservation(inventoryCol, reservationCol, event)
		default:
			fmt.Println("⚠️ Unknown status:", event.Status)
		}
	}
}

func commitInventory(inventoryCol, reservationCol *mongo.Collection, event PaymentEvent) {
	ctx := context.Background()

	// Find the reservation
	var reservation struct {
		Items []struct {
			ProductID string `bson:"productId"`
			Quantity  int64  `bson:"quantity"`
		} `bson:"items"`
		Status string `bson:"status"`
	}
	err := reservationCol.FindOne(ctx, bson.M{"reservationId": event.ReservationID}).Decode(&reservation)
	if err != nil {
		fmt.Printf("❌ Reservation not found: %s, error: %v\n", event.ReservationID, err)
		return
	}

	// Check if already processed
	if reservation.Status == "COMMITTED" {
		fmt.Printf("⚠️  Reservation %s already committed, skipping\n", event.ReservationID)
		return
	}

	// Update inventory for each item
	for _, item := range reservation.Items {
		// Convert string productId to ObjectID
		productID, err := primitive.ObjectIDFromHex(item.ProductID)
		if err != nil {
			fmt.Printf("❌ Invalid productId format: %s, error: %v\n", item.ProductID, err)
			continue
		}

		// First, check if inventory exists and has sufficient quantity
		var inventory struct {
			ID        primitive.ObjectID `bson:"_id"`
			ProductId primitive.ObjectID `bson:"productId"`
			Inventory int64              `bson:"inventory"`
		}
		err = inventoryCol.FindOne(ctx, bson.M{"productId": productID}).Decode(&inventory)
		if err != nil {
			fmt.Printf("❌ Product inventory not found: %s, error: %v\n", item.ProductID, err)
			continue
		}

		// Check if sufficient inventory available
		if inventory.Inventory < item.Quantity {
			fmt.Printf("❌ Insufficient inventory for product %s: available %d, required %d\n",
				item.ProductID, inventory.Inventory, item.Quantity)
			continue
		}

		// Soft update: decrement inventory and update timestamp
		filter := bson.M{"productId": productID}
		update := bson.M{
			"$inc": bson.M{"inventory": -item.Quantity},
			"$set": bson.M{"updatedAt": time.Now()},
		}
		result, err := inventoryCol.UpdateOne(ctx, filter, update)
		if err != nil {
			fmt.Printf("❌ Error committing inventory for product %s: %v\n", item.ProductID, err)
			continue
		}
		if result.MatchedCount == 0 {
			fmt.Printf("⚠️  No inventory document matched for product %s\n", item.ProductID)
			continue
		}
		fmt.Printf("✅ Inventory decremented for product %s: -%d\n", item.ProductID, item.Quantity)
	}

	// Update reservation status to COMMITTED
	updateResult, err := reservationCol.UpdateOne(ctx,
		bson.M{"reservationId": event.ReservationID},
		bson.M{"$set": bson.M{"status": "COMMITTED"}},
	)
	if err != nil {
		fmt.Printf("❌ Error updating reservation status: %v\n", err)
		return
	}
	if updateResult.MatchedCount == 0 {
		fmt.Printf("⚠️  Reservation %s not found for status update\n", event.ReservationID)
		return
	}
	fmt.Printf("✅ Inventory committed for reservation: %s\n", event.ReservationID)
}

func cancelReservation(inventoryCol, reservationCol *mongo.Collection, event PaymentEvent) {
	ctx := context.Background()

	// Check if reservation exists and get its status
	var reservation struct {
		Status string `bson:"status"`
	}
	err := reservationCol.FindOne(ctx, bson.M{"reservationId": event.ReservationID}).Decode(&reservation)
	if err != nil {
		fmt.Printf("❌ Reservation not found for cancellation: %s, error: %v\n", event.ReservationID, err)
		return
	}

	// Check if already processed
	if reservation.Status == "CANCELLED" {
		fmt.Printf("⚠️  Reservation %s already cancelled, skipping\n", event.ReservationID)
		return
	}

	// Update reservation status to CANCELLED
	updateResult, err := reservationCol.UpdateOne(ctx,
		bson.M{"reservationId": event.ReservationID},
		bson.M{"$set": bson.M{"status": "CANCELLED"}},
	)
	if err != nil {
		fmt.Printf("❌ Error cancelling reservation: %v\n", err)
		return
	}
	if updateResult.MatchedCount == 0 {
		fmt.Printf("⚠️  Reservation %s not found for cancellation\n", event.ReservationID)
		return
	}
	fmt.Printf("🚫 Reservation cancelled: %s\n", event.ReservationID)
}
