package httpapi

import "sync"

type MessageStore interface {
	Save(message string)
	List() []string
}

type InMemoryStore struct {
	mu       sync.RWMutex
	messages []string
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{messages: make([]string, 0, 64)}
}

func (s *InMemoryStore) Save(message string) {
	s.mu.Lock()
	s.messages = append(s.messages, message)
	s.mu.Unlock()
}

func (s *InMemoryStore) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.messages))
	copy(out, s.messages)
	return out
}
