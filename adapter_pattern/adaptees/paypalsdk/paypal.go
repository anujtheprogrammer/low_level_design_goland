package paypalsdk

type PayPalClient struct {
	ClientId string
	Secret   string
}

type PaypalTxn struct {
	TXNId string
	State string
}

func (c *PayPalClient) MakePayment(amount int64, payerEmail string) (PaypalTxn, error) {
	return PaypalTxn{TXNId: "paypal 123", State: "success"}, nil
}
