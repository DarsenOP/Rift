package storage

import (
	"sync"
	"testing"
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
