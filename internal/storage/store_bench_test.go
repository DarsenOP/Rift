package storage

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkSet(b *testing.B) {
	s := New()
	defer s.Shutdown()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Set(fmt.Sprintf("key%d", i), "value")
	}
}

func BenchmarkGetHit(b *testing.B) {
	s := New()
	defer s.Shutdown()
	s.Set("k", "v")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Get("k")
	}
}

func BenchmarkGetMiss(b *testing.B) {
	s := New()
	defer s.Shutdown()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Get("missing")
	}
}

func BenchmarkExpire(b *testing.B) {
	s := New()
	defer s.Shutdown()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("k%d", i)
		s.Set(key, "v")
		s.Expire(key, 10*time.Second)
	}
}
