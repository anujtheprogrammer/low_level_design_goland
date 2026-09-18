package loadbalancer

import (
	"net/http"
	"time"
)

// LoadBalancer wires everything together and implements http.Handler,
// so it drops straight into an http.Server. Notice how little logic
// lives here -- it delegates to the pool, the strategy, and the proxy.
type LoadBalancer struct {
	pool          *ServerPool
	strategy      IStrategy
	healthChecker *Healthchecker
}

// // ServeHTTP implements [http.Handler].
// func (lb *LoadBalancer) ServeHTTP(http.ResponseWriter, *http.Request) {
// 	panic("unimplemented")
// }

type Config struct {
	Backends            []string
	Strategy            IStrategy
	HealthCheckInterval time.Duration
	HealthCheckTimeout  time.Duration
}

func NewLoadBalancer(cfg Config) (*LoadBalancer, error) {
	pool := NewServerPool()
	for _, raw := range cfg.Backends {
		b, err := NewBackend(raw, 1)
		if err != nil {
			return nil, err
		}
		pool.AddBackend(b)
	}

	hc := NewHealthChecker(pool, cfg.HealthCheckInterval, cfg.HealthCheckTimeout)
	hc.Start()

	return &LoadBalancer{pool: pool, strategy: cfg.Strategy, healthChecker: hc}, nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	healthy := lb.pool.GetHealthyBackends()
	backends := lb.strategy.NextBackend(healthy)

	if backends == nil {
		http.Error(w, "no healthy backend available", http.StatusServiceUnavailable)
		return
	}

	backends.IncrementConn()
	defer backends.DecrementConn()
	backends.Proxy.ServeHTTP(w, r)
}
