package main

import (
	"fmt"
	lb "load_balancer/loadbalancer"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("we are implementing load balancer")

	balancer, err := lb.NewLoadBalancer(lb.Config{
		Backends:            []string{"http://localhost:9001", "http://localhost:9002"},
		Strategy:            &lb.RoundRobinStrategy{},
		HealthCheckInterval: 5 * time.Second,
		HealthCheckTimeout:  2 * time.Second,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Println("load balancer  lostening on : 8080")
	log.Fatal(http.ListenAndServe(":8080", balancer))
}
