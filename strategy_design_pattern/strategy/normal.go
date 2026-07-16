package strategy

type NormalPricing struct {
}

func (n NormalPricing) CalculateFare(distance float64) float64 {
	return distance * 10
}
