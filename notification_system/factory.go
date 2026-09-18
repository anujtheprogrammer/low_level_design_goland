package main

import "sync"

// ChannelFactory maps a ChannelType to its Notifier implementation.
// New channels are added via Register() - NotificationService never
// needs a switch/case on channel type (Open/Closed Principle).
type ChannelFactory struct {
	mu       sync.RWMutex
	registry map[ChannelType]Notifier
}

// NewChannelFactory builds an empty factory. Call Register for
// each supported channel at application start-up (wiring/DI code).
func NewChannelFactory() *ChannelFactory {
	return &ChannelFactory{registry: make(map[ChannelType]Notifier)}
}

// register add/overides the notifier given for a channel
func (f *ChannelFactory) Register(n Notifier) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	f.registry[n.Name()] = n
}

// GetNotifier returns the Notifier for a channel, or ErrUnknownChannel
// if nothing was registered for it.
func (f *ChannelFactory) GetNotifier(t ChannelType) (Notifier, error) {
	f.mu.RLock()
	defer f.mu.Unlock()
	n, ok := f.registry[t]
	if !ok {
		return nil, ErrUnknownChannel
	}
	return n, nil
}
