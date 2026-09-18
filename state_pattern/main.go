package main

import "fmt"

func main() {
	fmt.Printf("we are implementing state design pattern")

	order := NewOrder("Ord - 1001")
	fmt.Println("Start state:", order.CurrentState())

	if err := order.Pay(); err != nil {
		fmt.Println("error,", err)
	}
	if err := order.Ship(); err != nil {
		fmt.Println("error:", err)
	}
	if err := order.Deliver(); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("Final State : ", order.CurrentState())
	fmt.Println()

	fmt.Println("--- Now trying an ILLEGAL transition ---")
	// Try to cancel an already-delivered order — this must fail,
	// and DeliveredState is the only place that decides that.
	if err := order.Cancel(); err != nil {
		fmt.Println("error:", err)
	}
}
