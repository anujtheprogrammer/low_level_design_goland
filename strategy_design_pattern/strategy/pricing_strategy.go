package strategy

type PricingStrategy interface {
	CalculateFare(distance float64) float64
}
