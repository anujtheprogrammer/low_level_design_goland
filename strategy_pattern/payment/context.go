package payment

import "context"

type Processor struct {
	strategy PaymentStrategy
}

func NewProcessor(strategy PaymentStrategy) *Processor {
	return &Processor{strategy: strategy}
}

// this will allow switching of strategy at runtime
func (p *Processor) SetStrategy(strategy PaymentStrategy) {
	p.strategy = strategy
}

func (p *Processor) Process(ctx context.Context, req PaymentRequest) (PaymentResult, error) {
	err := p.strategy.Validate(req)
	if err != nil {
		return PaymentResult{}, nil
	}
	return p.strategy.Pay(ctx, req)
}
