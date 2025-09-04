package storage

import (
	"sync"
	"time"
)

// Store is a simple, thread-safe in-memory string-to-string store.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
	exps map[string]time.Time

	stopJanitor chan struct{}
}

// New returns an empty Store ready for use.
func New() *Store {
	s := &Store{
		data:        make(map[string]string),
		exps:        make(map[string]time.Time),
		stopJanitor: make(chan struct{}),
	}

	go s.janitor()
	return s
}

// Shutdown stops the background janitor goroutine.
// Idempotent; safe to call multiple times.
func (s *Store) Shutdown() {
	s.mu.Lock()

	select {
	case <-s.stopJanitor:
	default:
		close(s.stopJanitor)
	}

	s.mu.Unlock()
}

// janitor wakes every 100 ms and deletes expired keys.
func (s *Store) janitor() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.deleteExpired()
		case <-s.stopJanitor:
			return
		}
	}
}

func (s *Store) deleteExpired() {
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, expiry := range s.exps {
		if now.After(expiry) {
			delete(s.data, key)
			delete(s.exps, key)
		}
	}
}

// Set writes or overwrites a key with the given value.
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	// Also delete the expiry if any
	delete(s.exps, key)
}

// Get returns the value and true if the key exists; otherwise (“”, false).
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if expiry, ok := s.exps[key]; ok && time.Now().After(expiry) {
		return "", false
	}

	v, ok := s.data[key]

	return v, ok
}

// Del removes the key and returns true if it was present.
func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key]

	delete(s.data, key)
	delete(s.exps, key)

	return ok
}

// Exists reports whether the key is present.
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if expiry, ok := s.exps[key]; ok && time.Now().After(expiry) {
		return false
	}

	_, ok := s.data[key]

	return ok
}

// Expire sets a TTL on an existing key. Returns true if the key exists.
func (s *Store) Expire(key string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[key]; !ok {
		return false
	}

	s.exps[key] = time.Now().Add(ttl)
	return true
}

// TTL returns:
//
//	-2  key does not exist
//	-1  key exists but has no associated expiration
//	>=0 seconds left to live
func (s *Store) TTL(key string) (time.Duration, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.data[key]; !ok {
		return -2, true // Redis TTL semantic
	}

	exp, has := s.exps[key]
	if !has {
		return -1, true
	}

	left := time.Until(exp)
	if left < 0 {
		return 0, true
	}

	return left, true
}

func (s *Store) SetEX(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	s.exps[key] = time.Now().Add(ttl)
}
