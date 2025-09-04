package storage

import (
	"sync"
	"testing"
	"time"
)

func TestBasicCRUD(t *testing.T) {
	s := New()

	// Initial state
	if v, ok := s.Get("k"); ok || v != "" {
		t.Fatalf("expected empty, got %q %v", v, ok)
	}
	if s.Exists("k") {
		t.Fatal("key should not exist")
	}

	s.Set("k", "v")
	if v, ok := s.Get("k"); !ok || v != "v" {
		t.Fatalf("expected v=true, got %q %v", v, ok)
	}
	if !s.Exists("k") {
		t.Fatal("key should exist")
	}

	if !s.Del("k") {
		t.Fatal("Del should return true")
	}
	if s.Exists("k") {
		t.Fatal("key should be gone")
	}
}

func TestExpirationBlackBox(t *testing.T) {
	s := New()
	defer s.Shutdown()

	// 1. key with 50 ms TTL
	s.Set("k", "v")
	s.Expire("k", 50*time.Millisecond)

	// immediately
	if v, ok := s.Get("k"); !ok || v != "v" {
		t.Errorf("expected v, got %q %v", v, ok)
	}
	if ttl, _ := s.TTL("k"); ttl <= 0 || ttl > 50*time.Millisecond {
		t.Errorf("bad initial ttl: %v", ttl)
	}

	// 2. wait > TTL so janitor can run at least once
	time.Sleep(150 * time.Millisecond)

	// key must be gone
	if _, ok := s.Get("k"); ok {
		t.Error("key should be expired and deleted")
	}
	if ttl, _ := s.TTL("k"); ttl != -2 {
		t.Errorf("expected TTL -2, got %v", ttl)
	}
}

func TestExpireOverwrite(t *testing.T) {
	s := New()
	defer s.Shutdown()

	s.Set("k", "v")
	s.Expire("k", 200*time.Millisecond)

	// shorten expiry
	s.Expire("k", 50*time.Millisecond)

	time.Sleep(100 * time.Millisecond)
	if _, ok := s.Get("k"); ok {
		t.Error("key should be gone after shortened expiry")
	}
}

func TestManyExpirations(t *testing.T) {
	s := New()
	defer s.Shutdown()

	// insert 26 keys with staggered TTL
	for r := 'a'; r <= 'z'; r++ {
		key := "k" + string(r)
		s.Set(key, "v")
		s.Expire(key, time.Duration(r-'a'+1)*20*time.Millisecond)
	}

	// sleep until half are expired
	time.Sleep(300 * time.Millisecond)

	count := 0
	for r := 'a'; r <= 'z'; r++ {
		if s.Exists("k" + string(r)) {
			count++
		}
	}
	// roughly 13 left; allow ±3 for scheduling jitter
	if count < 10 || count > 16 {
		t.Errorf("expected ~13 keys left, got %d", count)
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := New()
	const workers = 100
	const opsPerWorker = 1000

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				key := string(rune('A' + (id+j)%26))
				s.Set(key, key)
				s.Get(key)
				s.Exists(key)
				s.Del(key)
			}
		}(i)
	}
	wg.Wait()
}
