package storage

import (
	"errors"
	"time"
)

var (
	ErrWrongType = errors.New("wrong type operation")
	ErrNotFound  = errors.New("key not found")
)

// Type checking utilities
// func (s *Store) checkType(key string, expected DataType) (*Value, error) {
// 	value, exists := s.data[key]
// 	if !exists {
// 		return nil, ErrNotFound
// 	}
// 	if value.Type != expected {
// 		return nil, ErrWrongType
// 	}
// 	return value, nil
// }

// Expiration utilities
func setExpiry(value *Value, ttl time.Duration) {
	expiry := time.Now().Add(ttl)
	switch value.Type {
	case StringType:
		value.String.Expiry = &expiry
	case ListType:
		value.List.Expiry = &expiry
	case HashType:
		value.Hash.Expiry = &expiry
	case SetType:
		value.Set.Expiry = &expiry
	}
}

func getExpiry(value *Value) *time.Time {
	switch value.Type {
	case StringType:
		return value.String.Expiry
	case ListType:
		return value.List.Expiry
	case HashType:
		return value.Hash.Expiry
	case SetType:
		return value.Set.Expiry
	}
	return nil
}
