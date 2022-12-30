package discord

import "sync"

func DefaultTracker() *Tracker {
	return &Tracker{
		messagesIds: make(map[string]string),
	}
}

type Tracker struct {
	mu          sync.Mutex
	messagesIds map[string]string
}

func (t *Tracker) TrackMessageID(key string, id string) {
	t.mu.Lock()
	t.messagesIds[key] = id
	t.mu.Unlock()
}

func (t *Tracker) GetMessageID(key string) (string, bool) {
	t.mu.Lock()
	id, ok := t.messagesIds[key]
	t.mu.Unlock()
	return id, ok
}
