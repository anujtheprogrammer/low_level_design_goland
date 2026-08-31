package payment

import "context"

// this contains the data every strategy needed
type PaymentRequest struct {
	OrderID  string
	Amount   float64
	Currency string
	Metadata map[string]string
}

// this is returned by any strategy on success
type PaymentResult struct {
	TransactionID string
	Status        string
}

// this is the common intertface that every payment method must satisfy
type PaymentStrategy interface {
	Validate(req PaymentRequest) error
	Pay(ctx context.Context, req PaymentRequest) (PaymentResult, error)
	Name() string
}
