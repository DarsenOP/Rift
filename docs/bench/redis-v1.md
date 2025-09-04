# Redis-Benchmark Stress Test 

**Hardware:** 13th Gen Intel i9-13980HX, 32 cores  
**Go:** 1.24 linux/amd64  
**Tool:** redis-benchmark 7.2.4 (single-core, 50 clients, 100 k ops each)  
**Date:** 2025-09-04

## Summary

| Command | Ops/sec | p50 (ms) | p99 (ms) | max (ms) |
|---------|--------:|---------:|---------:|---------:|
| **SET** | 240 385 |    0.103 |    0.263 |    1.263 |
| **GET** | 265 957 |    0.095 |    0.231 |    1.479 |

## Observations

- Zero panics, zero goroutine leaks (`go test -race` clean)  
- Global RWMutex handles 250 k ops/s single-core without contention collapse  
- Latency distribution tight: 99 % of ops < 0.3 ms  
- Reads slightly faster (no write-lock) – expected  
- Baseline captured; future sharding or lock-striping will be measured against these numbers
