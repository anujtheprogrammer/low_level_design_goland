package checkout

import (
	"adapter_pattern/payment"
	"context"
	"fmt"
)

type CheckoutService struct {
	processor payment.PaymentProcessor
}

func NewCheckoutService(p payment.PaymentProcessor) *CheckoutService {
	return &CheckoutService{processor: p}
}

func (s *CheckoutService) Checkout(ctx context.Context, orderId string, amount payment.Money) error {
	result, err := s.processor.Pay(ctx, orderId, amount)
	if err != nil {
		return fmt.Errorf("error in payemt %w", err)
	}
	fmt.Printf("Order succesfully executed transaction id is : %s", result)
	return nil
}
