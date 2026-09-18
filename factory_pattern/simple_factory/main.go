package main

import "fmt"

type Ride interface {
	EstimateFare(dist float64) float64
}

type Bike struct{}

func (Bike) EstimateFare(dist float64) float64 {
	return dist * 3.9
}

type Auto struct{}

func (Auto) EstimateFare(dist float64) float64 {
	return dist * 8.9
}

type Cab struct{}

func (Cab) EstimateFare(dist float64) float64 {
	return dist * 13.00
}

func NewRide(kind string) (Ride, error) {
	switch kind {
	case "bike":
		return Bike{}, nil
	case "Auto":
		return Auto{}, nil
	case "cab":
		return Cab{}, nil
	default:
		return nil, fmt.Errorf("Ride of unknown kind %q", kind)
	}
}

func main() {
	fmt.Println("this is simple factory pattern")
}
