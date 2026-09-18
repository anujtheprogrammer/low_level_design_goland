package loadbalancer

import "sync"

type ServerPool struct {
	backends []*Backend
	current  uint64
	mu       sync.RWMutex
}

func NewServerPool() *ServerPool {
	return &ServerPool{backends: make([]*Backend, 0)}
}

func (s *ServerPool) AddBackend(b *Backend) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.backends = append(s.backends, b)
}

func (s *ServerPool) RemoveBackend(rawURL string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, b := range s.backends {
		if b.URL.String() == rawURL {
			s.backends = append(s.backends[:i], s.backends[i+1:]...)
			return
		}
	}
}

func (s *ServerPool) GetHealthyBackends() []*Backend {
	s.mu.Lock()
	defer s.mu.Unlock()
	healthy := make([]*Backend, 0, len(s.backends))
	for _, b := range s.backends {
		if b.IsAlive() {
			healthy = append(healthy, b)
		}
	}
	return healthy
}

func (s *ServerPool) All() []*Backend {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]*Backend, len(s.backends))
	copy(out, s.backends)
	return out
}

func (s *ServerPool) MarkBackendStatus(rawURL string, alive bool) {
	for _, b := range s.All() {
		if b.URL.String() == rawURL {
			b.SetAlive(alive)
			return
		}
	}
}
