package main

import (
	"adapter_pattern/adaptees/stripesdk"
	"adapter_pattern/checkout"
	"adapter_pattern/payment"
	"context"
	"fmt"
)

func main() {
	fmt.Println("Hi this is adapter pattern")

	ctx := context.Background()
	amount := payment.Money{Amount: 13200, Currency: "INR"}

	var processor payment.PaymentProcessor = payment.NewStripeAdapter(
		&stripesdk.StripeClient{APIKey: "testing_api_api"},
	)

	service := checkout.NewCheckoutService(processor)
	_ = service.Checkout(ctx, "some random id", amount)
}
