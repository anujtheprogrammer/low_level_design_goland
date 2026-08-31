package payment

import (
	"context"
	"errors"
	"fmt"
)

// this is credit card strategy
type CreditcardStrategy struct {
	Gateway CardGateway
}

type CardGateway interface {
	Charge(ctx context.Context, cardToken string, amount float64) (string, error)
}

func (c *CreditcardStrategy) Name() string {
	return "credit card is being used"
}

func (c *CreditcardStrategy) Validate(req PaymentRequest) error {
	if req.Metadata["Card_token"] == "" {
		return errors.New("missing card token")
	}
	if req.Amount <= 0 {
		return errors.New("value must be positive")
	}
	return nil
}

func (c *CreditcardStrategy) Pay(ctx context.Context, req PaymentRequest) (PaymentResult, error) {
	result, error := c.Gateway.Charge(ctx, "random_token", 232)
	if error != nil {
		return PaymentResult{}, fmt.Errorf("creadit card charging failed %w", error)
	}
	return PaymentResult{TransactionID: result, Status: "succesful"}, nil
}

// this is the UPI strategy
type UPIStrategy struct {
	provide UPIProvider
}

type UPIProvider interface {
	PayByUPI(ctx context.Context, upiid string, amount float64) (string, error)
}

func (u *UPIStrategy) Name() string {
	return "UPI"
}

func (u *UPIStrategy) Validate(req PaymentRequest) error {
	if req.Amount < 0 {
		return fmt.Errorf("amount cannot be less than zero")
	}
	return nil
}

func (u *UPIStrategy) Pay(ctx context.Context, req PaymentRequest) (PaymentResult, error) {
	txnID, err := u.provide.PayByUPI(ctx, req.Metadata["ID"], 3425)
	if err != nil {
		return PaymentResult{}, fmt.Errorf("new error came %w", err)
	}
	return PaymentResult{TransactionID: txnID, Status: "confirmation_pending"}, nil
}
