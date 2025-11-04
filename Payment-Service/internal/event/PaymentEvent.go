package event

import "time"

type PaymentEvent struct {
	OrderID   string    `json:"orderId"`
	UserID    string    `json:"userId"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"` // "success" or "failure"
	Reason    string    `json:"reason,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
