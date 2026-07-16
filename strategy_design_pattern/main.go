package main

import (
	"fmt"
	"strategy_design_pattern/factory"
	"strategy_design_pattern/ride"
)

func main() {
	fmt.Println(" we are reading startegy design pattern")

	pricing := factory.GetPricingStrategy("premium")

	r := ride.Ride{
		Distance: 20,
		Pricing:  pricing,
	}

	fmt.Println(r.CalculateFare())
}
