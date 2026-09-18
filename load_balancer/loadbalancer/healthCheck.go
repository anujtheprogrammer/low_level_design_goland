package loadbalancer

import (
	"net/http"
	"time"
)

type Healthchecker struct {
	pool     *ServerPool
	interval time.Duration
	client   *http.Client
	stopCh   chan struct{}
}

func NewHealthChecker(pool *ServerPool, interval, timeout time.Duration) *Healthchecker {
	return &Healthchecker{
		pool:     pool,
		interval: interval,
		client:   &http.Client{Timeout: timeout},
		stopCh:   make(chan struct{}),
	}
}

// start run checks on a ticker in a background goroutine
func (h *Healthchecker) Start() {
	ticker := time.NewTicker(h.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				h.runChecks()
			case <-h.stopCh:
				return
			}
		}
	}()
}

// stop let us shut the checker down cleanly(e.g during test)
func (h *Healthchecker) Stop() {
	close(h.stopCh)
}

func (h *Healthchecker) runChecks() {
	for _, b := range h.pool.All() {
		go func(b *Backend) {
			b.SetAlive(h.checkBackend(b))
		}(b)
	}
}

func (h *Healthchecker) checkBackend(b *Backend) bool {
	resp, err := h.client.Get(b.URL.String() + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
