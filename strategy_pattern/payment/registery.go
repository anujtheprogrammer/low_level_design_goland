package payment

import "fmt"

type StretegyRegistery struct {
	strategies map[string]PaymentStrategy
}

func NewStrategyRegistery() *StretegyRegistery {
	return &StretegyRegistery{strategies: make(map[string]PaymentStrategy)}
}

func (r *StretegyRegistery) Register(s PaymentStrategy) {
	r.strategies[s.Name()] = s
}

func (r *StretegyRegistery) Resolve(method string) (PaymentStrategy, error) {
	s, ok := r.strategies[method]
	if !ok {
		return nil, fmt.Errorf("no strategy forpayment method %q", method)
	}
	return s, nil
}
