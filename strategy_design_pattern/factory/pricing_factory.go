package factory

import "strategy_design_pattern/strategy"

func GetPricingStrategy(rideType string) strategy.PricingStrategy {
	switch rideType {
	case "normal":
		return strategy.NormalPricing{}

	case "premium":
		return strategy.PremiumPricing{}

	case "sharing":
		return strategy.SharingStrategy{}

	default:
		return strategy.NormalPricing{}
	}
}
