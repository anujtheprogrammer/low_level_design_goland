package loadbalancer

import (
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

type Backend struct {
	URL         *url.URL // what is this url package
	Weight      int
	alive       bool
	activeConns int64
	mu          sync.RWMutex
	Proxy       *httputil.ReverseProxy // what is reverse proxy ?
}

func NewBackend(rawURL string, weight int) (*Backend, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return &Backend{
		URL:    u,
		Weight: weight,
		Proxy:  httputil.NewSingleHostReverseProxy(u),
		alive:  true,
	}, nil
}

// this will be called by health checker only
func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.alive
}

func (b *Backend) IncrementConn() {
	atomic.AndInt64(&b.activeConns, 1)
}

func (b *Backend) DecrementConn() {
	atomic.AddInt64(&b.activeConns, -1)
}

func (b *Backend) ActiveConnections() int64 {
	return atomic.LoadInt64(&b.activeConns)
}
