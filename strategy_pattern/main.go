package main

import (
	"context"
	"fmt"
	"log"
	"strategy_pattern/payment"
)

func main() {
	fmt.Println("we are studing strategy design pattern")

	registry := payment.NewStrategyRegistery()
	//registry.Register(&payment.CreditcardStrategy{Gateway: realCardGateway()})

	// Selection happens at runtime, e.g. based on a field in the incoming HTTP request.
	incomingMethod := "upi"

	strategy, err := registry.Resolve(incomingMethod)

	if err != nil {
		log.Fatal(err)
	}

	processor := payment.NewProcessor(strategy)

	result, err := processor.Process(context.Background(), payment.PaymentRequest{
		OrderID:  "ORD-1001",
		Amount:   499.00,
		Currency: "INR",
		Metadata: map[string]string{"vpa": "anuj@okhdfcbank"},
	})

	if err != nil {
		log.Fatal(err)
	}
	log.Printf("payment result: %+v", result)

}
