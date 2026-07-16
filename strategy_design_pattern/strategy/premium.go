package strategy

type PremiumPricing struct {
}

func (p PremiumPricing) CalculateFare(distance float64) float64 {
	return distance*10 + 121
}
