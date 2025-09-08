package storage

import (
	"errors"
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

// List operations

// LPush adds values to the left (head) of a list
func (s *Store) LPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.data[key]
	if exists && existing.Type != ListType {
		return 0, ErrWrongType
	}

	if !exists {
		s.data[key] = NewListValue(values)
		return len(values), nil
	}

	// Prepend new values to existing list
	existing.List.Values = append(values, existing.List.Values...)
	return len(existing.List.Values), nil
}

// RPush adds values to the right (tail) of a list
func (s *Store) RPush(key string, values ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.data[key]
	if exists && existing.Type != ListType {
		return 0, ErrWrongType
	}

	if !exists {
		s.data[key] = NewListValue(values)
		return len(values), nil
	}

	// Append new values to existing list
	existing.List.Values = append(existing.List.Values, values...)
	return len(existing.List.Values), nil
}

// LPop removes and returns the leftmost (head) element of a list
func (s *Store) LPop(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, err := s.checkType(key, ListType)
	if err != nil {
		return "", err
	}

	if len(value.List.Values) == 0 {
		return "", nil // Redis returns nil for empty list
	}

	// Remove and return first element
	popped := value.List.Values[0]
	value.List.Values = value.List.Values[1:]

	// Clean up empty list
	if len(value.List.Values) == 0 {
		delete(s.data, key)
	}

	return popped, nil
}

// RPop removes and returns the rightmost (tail) element of a list
func (s *Store) RPop(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, err := s.checkType(key, ListType)
	if err != nil {
		return "", err
	}

	if len(value.List.Values) == 0 {
		return "", nil
	}

	// Remove and return last element
	lastIndex := len(value.List.Values) - 1
	popped := value.List.Values[lastIndex]
	value.List.Values = value.List.Values[:lastIndex]

	// Clean up empty list
	if len(value.List.Values) == 0 {
		delete(s.data, key)
	}

	return popped, nil
}

// LRange returns a range of elements from the list
func (s *Store) LRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, err := s.checkType(key, ListType)
	if err != nil {
		return nil, err
	}

	list := value.List.Values
	length := len(list)

	// Handle negative indices (Redis behavior)
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	// Clamp indices to valid range
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop {
		return []string{}, nil
	}

	return list[start : stop+1], nil
}

// LLen returns the length of a list
func (s *Store) LLen(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, err := s.checkType(key, ListType)
	if err != nil {
		return 0, err
	}

	return len(value.List.Values), nil
}

// --- Hash operations --------------------------------------------------------

func (s *Store) HSet(key string, fieldVals ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(fieldVals)%2 != 0 {
		return 0, errors.New("wrong number of arguments for HSET")
	}

	// Create or load hash
	v, exists := s.data[key]
	if !exists {
		v = NewHashValue(nil)
		s.data[key] = v
	}
	if v.Type != HashType {
		return 0, ErrWrongType
	}
	if v.Hash.Fields == nil {
		v.Hash.Fields = make(map[string]string)
	}

	added := 0
	for i := 0; i < len(fieldVals); i += 2 {
		field, val := fieldVals[i], fieldVals[i+1]
		if _, ok := v.Hash.Fields[field]; !ok {
			added++
		}
		v.Hash.Fields[field] = val
	}
	return added, nil
}

func (s *Store) HGet(key, field string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, err := s.checkType(key, HashType)
	if err != nil {
		return "", err
	}
	val, ok := v.Hash.Fields[field]
	if !ok {
		return "", ErrNotFound
	}
	return val, nil
}

func (s *Store) HGetAll(key string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, err := s.checkType(key, HashType)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(v.Hash.Fields)*2)
	for f, val := range v.Hash.Fields {
		out = append(out, f, val)
	}
	return out, nil
}

func (s *Store) HDel(key string, fields ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, err := s.checkType(key, HashType)
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, f := range fields {
		if _, ok := v.Hash.Fields[f]; ok {
			delete(v.Hash.Fields, f)
			removed++
		}
	}
	return removed, nil
}

func (s *Store) HExists(key, field string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, err := s.checkType(key, HashType)
	if err != nil {
		return false, err
	}
	_, ok := v.Hash.Fields[field]
	return ok, nil
}

func (s *Store) HLen(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, err := s.checkType(key, HashType)
	if err != nil {
		return 0, err
	}
	return len(v.Hash.Fields), nil
}

// --- Set operations ---------------------------------------------------------

func (s *Store) SAdd(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, exists := s.data[key]
	if !exists {
		v = NewSetValue(nil)
		s.data[key] = v
	}
	if v.Type != SetType {
		return 0, ErrWrongType
	}
	if v.Set.Members == nil {
		v.Set.Members = make(map[string]struct{})
	}
	added := 0
	for _, m := range members {
		if _, ok := v.Set.Members[m]; !ok {
			v.Set.Members[m] = struct{}{}
			added++
		}
	}
	return added, nil
}

func (s *Store) SRem(key string, members ...string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, err := s.checkType(key, SetType)
	if err != nil {
		return 0, err
	}
	removed := 0
	for _, m := range members {
		if _, ok := v.Set.Members[m]; ok {
			delete(v.Set.Members, m)
			removed++
		}
	}
	return removed, nil
}

func (s *Store) SIsMember(key, member string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, err := s.checkType(key, SetType)
	if err != nil {
		return false, err
	}
	_, ok := v.Set.Members[member]
	return ok, nil
}

func (s *Store) SMembers(key string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, err := s.checkType(key, SetType)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(v.Set.Members))
	for m := range v.Set.Members {
		out = append(out, m)
	}
	return out, nil
}

func (s *Store) SCard(key string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, err := s.checkType(key, SetType)
	if err != nil {
		if err == ErrNotFound {
			return 0, nil
		}
		return 0, err
	}
	return len(v.Set.Members), nil
}

// SInter returns the intersection of the given sets.
// If any key does not exist it is considered empty.
func (s *Store) SInter(keys ...string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(keys) == 0 {
		return nil, nil
	}
	// Build a slice of maps (nil = empty set)
	maps := make([]map[string]struct{}, 0, len(keys))
	for _, k := range keys {
		v, ok := s.data[k]
		if !ok || v.Type != SetType {
			maps = append(maps, nil)
			continue
		}
		maps = append(maps, v.Set.Members)
	}
	// Find smallest map to iterate over it
	minIdx := 0
	for i, m := range maps {
		if m == nil {
			return []string{}, nil // intersection with empty set
		}
		if len(m) < len(maps[minIdx]) {
			minIdx = i
		}
	}
	// Check each member of the smallest map against all others
	out := make([]string, 0, len(maps[minIdx]))
	for member := range maps[minIdx] {
		ok := true
		for _, m := range maps {
			if _, exists := m[member]; !exists {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, member)
		}
	}
	return out, nil
}
