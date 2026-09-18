package main

import "sync"

// Preference stores a user's opt-in/opt-out choices per channel.
type Preference struct {
	UserID        string
	OptedChannels map[ChannelType]bool
}

// PreferenceRepository is intentionally an interface at the service
// layer, backed here by an in-memory map for the interview; in
// production this would hit Redis (cache) + Postgres (source of truth).
type PreferenceRepository interface {
	Get(UserID string) (Preference, error)
	IsOptedIN(UserID string, ch ChannelType) bool
}

type InMemoryPreferenceRepo struct {
	mu    sync.RWMutex
	store map[string]Preference
}

func NewPreferenceRepository() *InMemoryPreferenceRepo {
	return &InMemoryPreferenceRepo{store: make(map[string]Preference)}
}

func (i *InMemoryPreferenceRepo) Get(UserID string) (Preference, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	if p, ok := i.store[UserID]; ok {
		return p, nil
	}
	// default: opted in to everything if no explicit preference saved
	return Preference{UserID: UserID, OptedChannels: map[ChannelType]bool{
		ChannelEmail: true, ChannelSMS: true, ChannelPush: true,
	}}, nil
}

func (r *InMemoryPreferenceRepo) IsOptedIn(userID string, ch ChannelType) bool {
	p, _ := r.Get(userID)
	return p.OptedChannels[ch]
}
