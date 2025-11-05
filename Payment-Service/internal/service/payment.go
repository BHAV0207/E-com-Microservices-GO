package service

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ValidateUser(ctx context.Context, id string) bool {
	uri := "http://user-service:8080/users/"

	userUrl := fmt.Sprintf("%s%s", uri, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request to user service:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("User service returned non-OK status:", resp.StatusCode)
		return false
	}
	return true
}

func ValidateOrder(ctx context.Context, id string) bool {
	uri := "http://order-service:7000/order/"

	orderUrl := fmt.Sprintf("%s%s", uri, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, orderUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return false
	}

	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request to order service:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Order service returned non-OK status:", resp.StatusCode)
		return false
	}
	return true
}

func ProcessPayment(ctx context.Context, orderId string, amount float64, method string) bool {
	// Simulate payment processing logic
	time.Sleep(2 * time.Second) // Simulating processing delay
	// In real scenario, integrate with payment gateway here
	return true // Assume payment is always successful for this example
}
