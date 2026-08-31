package payment

import "context"

type Money struct {
	Amount   int64
	Currency string
}

// target interface
type PaymentProcessor interface {
	Pay(ctx context.Context, orderId string, amount Money) (transactionId string, err error)
}
