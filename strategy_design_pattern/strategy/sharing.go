package strategy

type SharingStrategy struct {
}

func (s SharingStrategy) CalculateFare(distance float64) float64 {
	return distance * 7
}
