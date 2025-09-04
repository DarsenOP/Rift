package storage

import "sync"

// Store is a simple, thread-safe in-memory string-to-string store.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// New returns an empty Store ready for use.
func New() *Store {
	return &Store{data: make(map[string]string)}
}

// Set writes or overwrites a key with the given value.
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get returns the value and true if the key exists; otherwise (“”, false).
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

// Del removes the key and returns true if it was present.
func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	delete(s.data, key)
	return ok
}

// Exists reports whether the key is present.
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[key]
	return ok
}
