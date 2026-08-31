package payment

import (
	"adapter_pattern/adaptees/paypalsdk"
	"adapter_pattern/adaptees/stripesdk"
	"context"
	"errors"
	"fmt"
)

type StripeAdapter struct {
	client *stripesdk.StripeClient
}

func NewStripeAdapter(client *stripesdk.StripeClient) *StripeAdapter {
	return &StripeAdapter{client: client}
}

func (a *StripeAdapter) Pay(ctx context.Context, orderId string, amount Money) (string, error) {
	result, err := a.client.ChargeCard(amount.Amount, amount.Currency, "order id :"+orderId)
	if err != nil {
		return "", fmt.Errorf("stripe adapter : %w", err)
	}
	if result.Status != "success" {
		return "", errors.New("stripe adapter: charge not succeeded")
	}
	return result.ChargeId, nil
}

type PaypalAdapter struct {
	client     *paypalsdk.PayPalClient
	payerEmail string
}

func NewPaypalAdapter(client *paypalsdk.PayPalClient, payerEmail string) *PaypalAdapter {
	return &PaypalAdapter{client: client, payerEmail: payerEmail}
}

func (a *PaypalAdapter) Pay(ctx context.Context, orderId string, amount Money) (string, error) {
	result, err := a.client.MakePayment(amount.Amount, a.payerEmail)
	if err != nil {
		return "", fmt.Errorf("paypal adapter : %w", err)
	}
	if result.State != "success" {
		return "", errors.New("transaction error")
	}
	return result.TXNId, nil
}
