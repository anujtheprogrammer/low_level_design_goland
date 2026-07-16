package ride

import "strategy_design_pattern/strategy"

type Ride struct {
	Distance float64
	Pricing  strategy.PricingStrategy
}

func (r Ride) CalculateFare() float64 {
	return r.Pricing.CalculateFare(r.Distance)
}
