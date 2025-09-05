# ADR 0001: Global RWMutex for Storage Locking

## Status

**Accepted**

## Context

We need to choose a locking strategy for our Redis-compatible in-memory store. 
Options considered:

- Single global RWMutex: Simple, thread-safe
- Sharded locks: Higher throughput but more complex
- Per-key mutexes: Maximum concurrency but significant overhead

Redis itself uses a single-threaded model, but our Go implementation requires thread safety for concurrent access.

## Decision

Use a single global RWMutex for the entire key-value store.

### Consequences

#### Positive

- ✅ Simple implementation and maintenance
- ✅ Thread-safe by default
- ✅ Good read scalability (multiple concurrent readers)
- ✅ Easy to implement and debug
- ✅ Minimal memory overhead

#### Negative

- ⚠️ Write throughput limited by single mutex
- ⚠️ Potential contention under heavy write loads

### Metrics (Current Baseline)

- SET: 240,385 ops/sec
- GET: 265,957 ops/sec
- p99 Latency: <0.3ms

This provides a solid foundation and we can easily shard later if needed. The simplicity aligns with our goal of building a maintainable Redis clone.
