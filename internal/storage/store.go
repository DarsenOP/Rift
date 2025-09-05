package storage

import (
	"sync"
	"time"
)

// Store is a simple, thread-safe in-memory string-to-string store.
type Store struct {
	mu          sync.RWMutex
	data        map[string]*Value
	stopJanitor chan struct{}
}

// New returns an empty Store ready for use.
func New() *Store {
	s := &Store{
		data:        make(map[string]*Value),
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

	for key, value := range s.data {
		expiry := getExpiry(value)

		if expiry != nil && now.After(*expiry) {
			delete(s.data, key)
		}
	}
}

// Set writes or overwrites a key with the given value.
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = NewStringValue(value)
}

// Get returns the value and true if the key exists; otherwise (“”, false).
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	if !exists || value.Type != StringType {
		return "", false
	}

	// Check expiration
	if value.String.Expiry != nil && time.Now().After(*value.String.Expiry) {
		return "", false
	}

	return value.String.Value, true
}

// Del removes the key and returns true if it was present.
func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[key]

	delete(s.data, key)
	return exists
}

// Exists reports whether the key is present.
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	if !exists {
		return false
	}

	// Check expiration
	expiry := getExpiry(value)

	if expiry != nil && time.Now().After(*expiry) {
		return false
	}

	return true
}

// Expire sets a TTL on an existing key. Returns true if the key exists.
func (s *Store) Expire(key string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.data[key]
	if !exists {
		return false
	}

	setExpiry(value, ttl)
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

	value, exists := s.data[key]
	if !exists {
		return -2, true
	}

	expiry := getExpiry(value)
	if expiry == nil {
		return -1, true
	}

	left := time.Until(*expiry)
	if left < 0 {
		return 0, true
	}

	return left, true
}

func (s *Store) SetEX(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	strValue := NewStringValue(value)
	setExpiry(strValue, ttl)
	s.data[key] = strValue
}
