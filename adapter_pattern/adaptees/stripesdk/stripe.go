package stripesdk

type StripeClient struct {
	APIKey string
}

type chargeResult struct {
	ChargeId string
	Status   string
}

func (c *StripeClient) ChargeCard(money int64, currency, description string) (*chargeResult, error) {
	return &chargeResult{ChargeId: "stripe123", Status: "succsess"}, nil
}
