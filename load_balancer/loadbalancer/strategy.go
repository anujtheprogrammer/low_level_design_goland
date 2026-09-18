package loadbalancer

import "sync"

type IStrategy interface {
	NextBackend(backends []*Backend) *Backend
}

// implementation of interface (round robin)
// just find the reminder and use that
type RoundRobinStrategy struct {
	counter uint64
	mu      sync.Mutex
}

func (rrs *RoundRobinStrategy) NextBackend(bs []*Backend) *Backend {
	if len(bs) == 0 {
		return nil
	}
	rrs.mu.Lock()
	defer rrs.mu.Unlock()
	idx := rrs.counter % uint64(len(bs))
	rrs.counter++
	return bs[idx]
}

// least connections
// pick the backend with fewest in-flight req. better than round robin when req cost varies a lot
type LeastConnectionStrategy struct{}

func (lcs *LeastConnectionStrategy) NexyBackend(bs []*Backend) *Backend {
	if len(bs) == 0 {
		return nil
	}

	best := bs[0]
	for _, back := range bs {
		if back.activeConns < best.activeConns {
			best = back
		}
	}
	return best
}

// weighted round robin
// smooth (Nginx uses the same) so a heavily weighted backend get it extra share spread evenly across the cycle instead of aarriving in  a big burst
type WeightedRoundRobinStrategy struct {
	mu      sync.Mutex
	current map[string]int
}

func NewWeightedRoundRobinStrategy() *WeightedRoundRobinStrategy {
	return &WeightedRoundRobinStrategy{current: make(map[string]int)}
}

func (wrr *WeightedRoundRobinStrategy) NextBackend(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	wrr.mu.Lock()
	defer wrr.mu.Unlock()

	total := 0
	var chosen *Backend
	for _, b := range backends {
		key := b.URL.String()
		wrr.current[key] += b.Weight
		total += b.Weight
		if chosen == nil || wrr.current[key] > wrr.current[chosen.URL.String()] {
			chosen = b
		}
	}
	wrr.current[chosen.URL.String()] -= total
	return chosen
}
